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

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openyurtio/openyurt/cmd/yurt-lb-agent/app/options"
	"github.com/openyurtio/openyurt/pkg/apis"
	"github.com/openyurtio/openyurt/pkg/yurtlbagent/controllers"
	"github.com/openyurtio/openyurt/pkg/yurtlbagent/util"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = apis.AddToScheme(clientgoscheme.Scheme)
	_ = apis.AddToScheme(scheme)

	// +kubebuilder:scaffold:scheme
}

func NewCmdYurtLBAgent(stopCh <-chan struct{}) *cobra.Command {
	yurtLBAgentOptions := options.NewYurtLBAgentOptions()
	cmd := &cobra.Command{
		Use:   "yurt-lb-agent",
		Short: "Launch yurt-lb-agent",
		Long:  "Launch yurt-lb-agent",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if err := options.ValidateOptions(yurtLBAgentOptions); err != nil {
				klog.Fatalf("validate options: %v", err)
			}
			Run(yurtLBAgentOptions, stopCh)
		},
	}

	yurtLBAgentOptions.AddFlags(cmd.Flags())
	return cmd
}

func Run(opts *options.YurtLBAgentOptions, stopCh <-chan struct{}) {
	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           scheme,
		LeaderElectionID: "yurt-lb-agent",
		Namespace:        opts.Namespace,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// perform preflight check
	setupLog.Info("[preflight] Running pre-flight checks")
	if err := preflightCheck(mgr, opts); err != nil {
		setupLog.Error(err, "could not run pre-flight checks")
		os.Exit(1)
	}

	if opts.Nodepool == "" {
		opts.Nodepool, err = util.GetNodePool(mgr.GetConfig(), opts.Node)
		if err != nil {
			setupLog.Error(err, "could not get the nodepool where yurt-lb-agent run")
			os.Exit(1)
		}
	}

	// setup the Service Reconciler
	if err = (&controllers.ServiceReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, opts); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Yurtlb-agent")
		os.Exit(1)
	}

	// if err := mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
	// 	setupLog.Error(err, "unable to set up health check")
	// 	os.Exit(1)
	// }
	// if err := mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
	// 	setupLog.Error(err, "unable to set up ready check")
	// 	os.Exit(1)
	// }

	setupLog.Info("[run controllers] Starting manager, acting on " + fmt.Sprintf("[Node: %s, Namespace: %s]", opts.Node, opts.Namespace))
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "could not running manager")
		os.Exit(1)
	}
}

func preflightCheck(mgr ctrl.Manager, opts *options.YurtLBAgentOptions) error {
	client, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return err
	}
	if _, err := client.CoreV1().Namespaces().Get(context.TODO(), opts.Namespace, metav1.GetOptions{}); err != nil {
		return err
	}
	return nil
}
