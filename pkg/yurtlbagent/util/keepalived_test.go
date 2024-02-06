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

// func Test_writeKeepalivedVIPs(t *testing.T) {
// 	type args struct {
// 		content string
// 		vips    []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "test-1",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 					}
// 				}`,
// 				vips: []string{"192.168.1.5", "192.168.1.6", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 						192.168.1.5
// 						192.168.1.6
// 						192.168.1.7
// 					}
// 				}`,
// 		},
// 		{
// 			name: "test-2",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 					}
// 				}`,
// 				vips: []string{"192.168.1.2", "192.168.1.3", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 						192.168.1.7
// 					}
// 				}`,
// 		},
// 		{
// 			name: "test-3",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				 }

// 				 vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 					}
// 				}`,
// 				vips: []string{"192.168.1.2", "192.168.1.2", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				 }

// 				 vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 						192.168.1.7
// 					}
// 				}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := WriteKeepalivedVIPs(tt.args.content, tt.args.vips); got != tt.want {
// 				t.Errorf("writeKeepalivedVIPs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_deleteKeepalivedVIPs(t *testing.T) {
// 	type args struct {
// 		content string
// 		vips    []string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{
// 			name: "test-1",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 						192.168.1.5
// 						192.168.1.6
// 						192.168.1.7
// 					}
// 				}`,
// 				vips: []string{"192.168.1.5", "192.168.1.6", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 					}
// 				}`,
// 		},
// 		{
// 			name: "test-2",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 						192.168.1.7
// 					}
// 				}`,
// 				vips: []string{"192.168.1.2", "192.168.1.3", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				}

// 				vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 					}
// 				}`,
// 		},
// 		{
// 			name: "test-3",
// 			args: args{
// 				content: `global_defs {
// 					router_id LVS_DEVEL
// 				 }

// 				 vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.2
// 						192.168.1.3
// 					}
// 				}`,
// 				vips: []string{"192.168.1.2", "192.168.1.2", "192.168.1.7"},
// 			},
// 			want: `global_defs {
// 					router_id LVS_DEVEL
// 				 }

// 				 vrrp_instance VI_1 {
// 					state MASTER
// 					interface eth0
// 					virtual_router_id 51
// 					priority 100
// 					advert_int 1
// 					authentication {
// 						auth_type PASS
// 						auth_pass 1111
// 					}
// 					virtual_ipaddress {
// 						192.168.1.1
// 						192.168.1.3
// 					}
// 				}`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := DeleteKeepalivedVIPs(tt.args.content, tt.args.vips); got != tt.want {
// 				t.Errorf("deleteKeepalivedVIPs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
