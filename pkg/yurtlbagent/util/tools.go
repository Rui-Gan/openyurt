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

package util

import (
	"context"

	"github.com/openyurtio/openyurt/pkg/projectinfo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetNodePool(cfg *rest.Config, nodeName string) (string, error) {
	var nodePool string
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nodePool, err
	}

	node, err := client.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Warningf("could not get node(%s) in yurtlb agent, %v", nodeName, err)
		return "", err
	}
	nodePool = node.Labels[projectinfo.GetNodePoolLabel()]
	return nodePool, nil
}
