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

package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// YurtIoTDockOptions is the main settings for the yurt-iot-dock
type YurtLBAgentOptions struct {
	Node              string
	Namespace         string
	Version           string
	Nodepool          string
	KeepalivedPath    string
	LoadBalancerClass string
	VIPs              map[string]struct{}
}

func NewYurtLBAgentOptions() *YurtLBAgentOptions {
	return &YurtLBAgentOptions{
		Node:              "",
		Namespace:         "default",
		Version:           "",
		Nodepool:          "",
		KeepalivedPath:    "/etc/keepalived/keepalived.conf",
		LoadBalancerClass: "yurtlb-agent",
		VIPs:              map[string]struct{}{},
	}
}

func ValidateOptions(options *YurtLBAgentOptions) error {
	if len(options.Node) == 0 {
		return fmt.Errorf("node name is empty")
	}
	if len(options.KeepalivedPath) == 0 {
		return fmt.Errorf("keepalived config path is empty")
	}
	return nil
}

func (o *YurtLBAgentOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Node, "node", "", "The node YurtLB Agent is deployed in.(just for debugging)")
	fs.StringVar(&o.Namespace, "namespace", "", "The namespace YurtLB Agent is deployed in.(just for debugging)")
	fs.StringVar(&o.Version, "version", "", "The version of edge resources deploymenet.")
	fs.StringVar(&o.Nodepool, "nodepool", "", "The nodePool YurtLB Agent is deployed in.(just for debugging)")
	fs.StringVar(&o.KeepalivedPath, "keepalived_path", "", "")
	fs.StringVar(&o.LoadBalancerClass, "LoadBalancerClass", "", "")
}
