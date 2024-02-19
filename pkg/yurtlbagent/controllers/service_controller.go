/*
Copyright 2024 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"net"

	"github.com/openyurtio/openyurt/cmd/yurt-lb-agent/app/options"
	"github.com/openyurtio/openyurt/pkg/yurtlbagent/util"
	v1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	keepalivedConf *util.Keepalived
)

func Format(format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s: %s", "yurtlb-agent: controller ServiceReconciler", s)
}

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	Node              string
	Nodepool          string
	Namespace         string
	KeepalivedPath    string
	KeepalivedPid     string
	LoadBalancerClass string
	VIPs              map[string]struct{}
}

// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch

func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof(Format("Reconcile Service %s/%s", req.Namespace, req.Name))

	// Fetch the service
	service := &v1.Service{}
	if err := r.Get(ctx, req.NamespacedName, service); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		klog.Errorf(Format("Get Service %s/%s error %v", req.Namespace, req.Name, err))
		return ctrl.Result{}, err
	}
	klog.Infof(Format("Get Service  %s/%s", req.Namespace, req.Name))

	// filter loadbalancerclass
	if filterByLoadBalancerClass(service, r.LoadBalancerClass) {
		klog.Infof(Format("filter Service %s/%s by LoadBalancerClass %s", req.Namespace, req.Name, r.LoadBalancerClass))
		return ctrl.Result{}, nil
	}

	klog.Infof(Format("not filter Service  %s/%sby LoadBalancerClass %s", req.Namespace, req.Name, r.LoadBalancerClass))

	if service.DeletionTimestamp != nil {
		return r.reconcileDelete(ctx, service)
	}

	return r.reconcileNormal(ctx, service)
}

func (r *ServiceReconciler) reconcileDelete(ctx context.Context, svc *v1.Service) (ctrl.Result, error) {
	// delete vip
	err := updateVIPs(r.VIPs, map[string]struct{}{}, r.KeepalivedPath)
	if err != nil {
		klog.Errorf(Format("update vip err: %v", err))
		return ctrl.Result{}, err
	}
	r.VIPs = map[string]struct{}{}

	// delete ipvs

	return ctrl.Result{}, nil
}

func (r *ServiceReconciler) reconcileNormal(ctx context.Context, svc *v1.Service) (ctrl.Result, error) {
	// 只有 有对应pod的node 才参与竞选！！
	newVIPs := getVIPsByNodepool(svc.Status.LoadBalancer.Ingress, r.Nodepool)
	err := updateVIPs(r.VIPs, newVIPs, r.KeepalivedPath)
	if err != nil {
		klog.Errorf(Format("update vip err: %v", err))
		return ctrl.Result{}, err
	}
	r.VIPs = newVIPs

	// bind ipvs

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager, opts *options.YurtLBAgentOptions) error {
	r.Node = opts.Node
	r.Nodepool = opts.Nodepool
	r.Namespace = opts.Namespace
	r.KeepalivedPath = opts.KeepalivedPath
	r.LoadBalancerClass = opts.LoadBalancerClass
	r.VIPs = opts.VIPs

	var err error

	keepalivedConf = &util.Keepalived{
		Iface:    "eth0",
		Priority: 100,
		VIPs:     []string{},
		Vrid:     20,
	}

	err = keepalivedConf.LoadTemplates()
	if err != nil {
		return fmt.Errorf("load keepalived template err, %v", err)
	}
	err = keepalivedConf.WriteCfg(r.KeepalivedPath, r.VIPs)
	if err != nil {
		return fmt.Errorf("write keepalived config err, %v", err)
	}
	err = keepalivedConf.Start()
	if err != nil {
		return fmt.Errorf("start keepalived err, %v", err)
	}
	// err = keepalivedConf.ReloadKeepalived()
	// if err != nil {
	// 	return fmt.Errorf("reload keepalived err, %v", err)
	// }

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Service{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				oldService, ok1 := e.ObjectOld.(*v1.Service)
				newService, ok2 := e.ObjectNew.(*v1.Service)
				if !ok1 || !ok2 {
					return false
				}
				return !equality.Semantic.DeepEqual(oldService.Status, newService.Status)
			},
			CreateFunc: func(ce event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(de event.DeleteEvent) bool {
				return true
			},
			GenericFunc: func(ge event.GenericEvent) bool {
				return true
			},
		}).
		Watches(&source.Kind{Type: &discovery.EndpointSlice{}},
			handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
				epSlice, ok := obj.(*discovery.EndpointSlice)
				if !ok {
					klog.Errorf(Format("received an object that is not epslice"))
					return []reconcile.Request{}
				}
				serviceName, err := ServiceKeyForSlice(epSlice)
				if err != nil {
					klog.Errorf(Format("failed to get serviceName for slice %s, %v", epSlice.Name, err))
					return []reconcile.Request{}
				}
				klog.Infof(Format("enqueueing %s for slice %s", serviceName, epSlice.Name))
				return []reconcile.Request{{NamespacedName: serviceName}}
			})).
		Complete(r)
}

func ServiceKeyForSlice(endpointSlice *discovery.EndpointSlice) (types.NamespacedName, error) {
	if endpointSlice == nil {
		return types.NamespacedName{}, fmt.Errorf("nil EndpointSlice")
	}
	serviceName, err := serviceNameForSlice(endpointSlice)
	if err != nil {
		return types.NamespacedName{}, err
	}

	return types.NamespacedName{Namespace: endpointSlice.Namespace, Name: serviceName}, nil
}

func serviceNameForSlice(endpointSlice *discovery.EndpointSlice) (string, error) {
	serviceName, ok := endpointSlice.Labels[discovery.LabelServiceName]
	if !ok || serviceName == "" {
		return "", fmt.Errorf("endpointSlice missing the %s label", discovery.LabelServiceName)
	}
	return serviceName, nil
}

func filterByLoadBalancerClass(service *v1.Service, loadBalancerClass string) bool {
	if service == nil {
		return false
	}
	if service.Spec.LoadBalancerClass == nil && loadBalancerClass != "" {
		return true
	}
	if service.Spec.LoadBalancerClass == nil && loadBalancerClass == "" {
		return false
	}
	if *service.Spec.LoadBalancerClass != loadBalancerClass {
		return true
	}
	return false
}

func getVIPsByNodepool(ingresses []v1.LoadBalancerIngress, nodepool string) map[string]struct{} {
	ips := []net.IP{}
	for _, ingress := range ingresses {
		if ingress.Hostname == nodepool {
			ip := net.ParseIP(ingress.IP)
			if ip == nil {
				klog.Errorf(Format("invalid LoadBalancer IP %s", ingress.IP))
			}
			ips = append(ips, ip)
		}
	}
	vips := map[string]struct{}{}
	for _, ip := range ips {
		vips[ip.String()] = struct{}{}
	}
	return vips
}

func updateVIPs(oldVIPs map[string]struct{}, newVIPs map[string]struct{}, keepalivedPath string) error {
	deleteVIPs := []string{}
	for old := range oldVIPs {
		if _, ok := newVIPs[old]; !ok {
			deleteVIPs = append(deleteVIPs, old)
		}
	}

	addVIPs := []string{}
	for new := range newVIPs {
		if _, ok := oldVIPs[new]; !ok {
			addVIPs = append(addVIPs, new)
		}
	}

	if len(deleteVIPs) == 0 && len(addVIPs) == 0 {
		return nil
	}

	err := keepalivedConf.WriteCfg(keepalivedPath, newVIPs)
	if err != nil {
		return fmt.Errorf("write keepalived config err, %v", err)
	}
	err = keepalivedConf.ReloadKeepalived()
	if err != nil {
		return fmt.Errorf("reload keepalived err, %v", err)
	}
	return nil
}
