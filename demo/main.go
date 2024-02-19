package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := "kubeconfig"

	// 使用 kubeconfig 文件获取配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// 在集群内部时使用 InClusterConfig
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// 创建客户端集
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取指定的 Service
	serviceClient := clientset.CoreV1().Services("default")

	service, err := serviceClient.Get(context.TODO(), "nginx2", v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 修改 Service 状态
	service.Status.LoadBalancer = corev1.LoadBalancerStatus{
		Ingress: []corev1.LoadBalancerIngress{
			{
				Hostname: "beijing",
				IP:       "172.20.0.6",
			},
			{
				Hostname: "hangzhou",
				IP:       "172.20.0.2",
			},
		},
	}

	// 更新 Service 状态
	_, err = serviceClient.UpdateStatus(context.TODO(), service, v1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Service status updated")
}
