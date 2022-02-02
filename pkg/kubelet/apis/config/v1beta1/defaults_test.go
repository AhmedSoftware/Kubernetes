/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logsapi "k8s.io/component-base/logs/api/v1"
	"k8s.io/kubelet/config/v1beta1"
	"k8s.io/kubernetes/pkg/cluster/ports"
	"k8s.io/kubernetes/pkg/kubelet/qos"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	utilpointer "k8s.io/utils/pointer"
)

func TestSetDefaultsKubeletConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		config   *v1beta1.KubeletConfiguration
		expected *v1beta1.KubeletConfiguration
	}{
		{
			"empty config",
			&v1beta1.KubeletConfiguration{},
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(true),
				SyncFrequency:      metav1.Duration{Duration: 1 * time.Minute},
				FileCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				HTTPCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				Address:            "0.0.0.0",
				Port:               ports.KubeletPort,
				Authentication: v1beta1.KubeletAuthentication{
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(false)},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(true),
						CacheTTL: metav1.Duration{Duration: 2 * time.Minute},
					},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeWebhook,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   metav1.Duration{Duration: 5 * time.Minute},
						CacheUnauthorizedTTL: metav1.Duration{Duration: 30 * time.Second},
					},
				},
				RegistryPullQPS:                           utilpointer.Int32Ptr(5),
				RegistryBurst:                             10,
				EventRecordQPS:                            utilpointer.Int32Ptr(5),
				EventBurst:                                10,
				EnableDebuggingHandlers:                   utilpointer.BoolPtr(true),
				HealthzPort:                               utilpointer.Int32Ptr(10248),
				HealthzBindAddress:                        "127.0.0.1",
				OOMScoreAdj:                               utilpointer.Int32Ptr(int32(qos.KubeletOOMScoreAdj)),
				PLEGRelistThreshold:                       metav1.Duration{Duration: 3 * time.Minute},
				StreamingConnectionIdleTimeout:            metav1.Duration{Duration: 4 * time.Hour},
				NodeStatusUpdateFrequency:                 metav1.Duration{Duration: 10 * time.Second},
				NodeStatusReportFrequency:                 metav1.Duration{Duration: 5 * time.Minute},
				NodeLeaseDurationSeconds:                  40,
				ImageMinimumGCAge:                         metav1.Duration{Duration: 2 * time.Minute},
				ImageGCHighThresholdPercent:               utilpointer.Int32Ptr(85),
				ImageGCLowThresholdPercent:                utilpointer.Int32Ptr(80),
				VolumeStatsAggPeriod:                      metav1.Duration{Duration: time.Minute},
				CgroupsPerQOS:                             utilpointer.BoolPtr(true),
				CgroupDriver:                              "cgroupfs",
				CPUManagerPolicy:                          "none",
				CPUManagerReconcilePeriod:                 metav1.Duration{Duration: 10 * time.Second},
				MemoryManagerPolicy:                       v1beta1.NoneMemoryManagerPolicy,
				TopologyManagerPolicy:                     v1beta1.NoneTopologyManagerPolicy,
				TopologyManagerScope:                      v1beta1.ContainerTopologyManagerScope,
				RuntimeRequestTimeout:                     metav1.Duration{Duration: 2 * time.Minute},
				HairpinMode:                               v1beta1.PromiscuousBridge,
				MaxPods:                                   110,
				PodPidsLimit:                              utilpointer.Int64(-1),
				ResolverConfig:                            utilpointer.String(kubetypes.ResolvConfDefault),
				CPUCFSQuota:                               utilpointer.BoolPtr(true),
				CPUCFSQuotaPeriod:                         &metav1.Duration{Duration: 100 * time.Millisecond},
				NodeStatusMaxImages:                       utilpointer.Int32Ptr(50),
				MaxOpenFiles:                              1000000,
				ContentType:                               "application/vnd.kubernetes.protobuf",
				KubeAPIQPS:                                utilpointer.Int32Ptr(5),
				KubeAPIBurst:                              10,
				SerializeImagePulls:                       utilpointer.BoolPtr(true),
				EvictionHard:                              DefaultEvictionHard,
				EvictionPressureTransitionPeriod:          metav1.Duration{Duration: 5 * time.Minute},
				EnableControllerAttachDetach:              utilpointer.BoolPtr(true),
				MakeIPTablesUtilChains:                    utilpointer.BoolPtr(true),
				IPTablesMasqueradeBit:                     utilpointer.Int32Ptr(DefaultIPTablesMasqueradeBit),
				IPTablesDropBit:                           utilpointer.Int32Ptr(DefaultIPTablesDropBit),
				FailSwapOn:                                utilpointer.BoolPtr(true),
				ContainerLogMaxSize:                       "10Mi",
				ContainerLogMaxFiles:                      utilpointer.Int32Ptr(5),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.WatchChangeDetectionStrategy,
				EnforceNodeAllocatable:                    DefaultNodeAllocatableEnforcement,
				VolumePluginDir:                           DefaultVolumePluginDir,
				Logging: logsapi.LoggingConfiguration{
					Format:         "text",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:        utilpointer.BoolPtr(true),
				EnableProfilingHandler:        utilpointer.BoolPtr(true),
				EnableDebugFlagsHandler:       utilpointer.BoolPtr(true),
				SeccompDefault:                utilpointer.BoolPtr(false),
				MemoryThrottlingFactor:        utilpointer.Float64Ptr(DefaultMemoryThrottlingFactor),
				RegisterNode:                  utilpointer.BoolPtr(true),
				LocalStorageCapacityIsolation: utilpointer.BoolPtr(true),
			},
		},
		{
			"all negative",
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(false),
				StaticPodPath:      "",
				SyncFrequency:      zeroDuration,
				FileCheckFrequency: zeroDuration,
				HTTPCheckFrequency: zeroDuration,
				StaticPodURL:       "",
				StaticPodURLHeader: map[string][]string{},
				Address:            "",
				Port:               0,
				ReadOnlyPort:       0,
				TLSCertFile:        "",
				TLSPrivateKeyFile:  "",
				TLSCipherSuites:    []string{},
				TLSMinVersion:      "",
				RotateCertificates: false,
				ServerTLSBootstrap: false,
				Authentication: v1beta1.KubeletAuthentication{
					X509: v1beta1.KubeletX509Authentication{ClientCAFile: ""},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(false),
						CacheTTL: zeroDuration,
					},
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(false)},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeWebhook,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   zeroDuration,
						CacheUnauthorizedTTL: zeroDuration,
					},
				},
				RegistryPullQPS:                  utilpointer.Int32(0),
				RegistryBurst:                    0,
				EventRecordQPS:                   utilpointer.Int32(0),
				EventBurst:                       0,
				EnableDebuggingHandlers:          utilpointer.BoolPtr(false),
				EnableContentionProfiling:        false,
				HealthzPort:                      utilpointer.Int32(0),
				HealthzBindAddress:               "",
				OOMScoreAdj:                      utilpointer.Int32(0),
				ClusterDomain:                    "",
				ClusterDNS:                       []string{},
				PLEGRelistThreshold:              zeroDuration,
				StreamingConnectionIdleTimeout:   zeroDuration,
				NodeStatusUpdateFrequency:        zeroDuration,
				NodeStatusReportFrequency:        zeroDuration,
				NodeLeaseDurationSeconds:         0,
				ImageMinimumGCAge:                zeroDuration,
				ImageGCHighThresholdPercent:      utilpointer.Int32(0),
				ImageGCLowThresholdPercent:       utilpointer.Int32(0),
				VolumeStatsAggPeriod:             zeroDuration,
				KubeletCgroups:                   "",
				SystemCgroups:                    "",
				CgroupRoot:                       "",
				CgroupsPerQOS:                    utilpointer.BoolPtr(false),
				CgroupDriver:                     "",
				CPUManagerPolicy:                 "",
				CPUManagerPolicyOptions:          map[string]string{},
				CPUManagerReconcilePeriod:        zeroDuration,
				MemoryManagerPolicy:              "",
				TopologyManagerPolicy:            "",
				TopologyManagerScope:             "",
				QOSReserved:                      map[string]string{},
				RuntimeRequestTimeout:            zeroDuration,
				HairpinMode:                      "",
				MaxPods:                          0,
				PodCIDR:                          "",
				PodPidsLimit:                     utilpointer.Int64(0),
				ResolverConfig:                   utilpointer.String(""),
				RunOnce:                          false,
				CPUCFSQuota:                      utilpointer.Bool(false),
				CPUCFSQuotaPeriod:                &zeroDuration,
				NodeStatusMaxImages:              utilpointer.Int32(0),
				MaxOpenFiles:                     0,
				ContentType:                      "",
				KubeAPIQPS:                       utilpointer.Int32(0),
				KubeAPIBurst:                     0,
				SerializeImagePulls:              utilpointer.Bool(false),
				EvictionHard:                     map[string]string{},
				EvictionSoft:                     map[string]string{},
				EvictionSoftGracePeriod:          map[string]string{},
				EvictionPressureTransitionPeriod: zeroDuration,
				EvictionMaxPodGracePeriod:        0,
				EvictionMinimumReclaim:           map[string]string{},
				PodsPerCore:                      0,
				EnableControllerAttachDetach:     utilpointer.Bool(false),
				ProtectKernelDefaults:            false,
				MakeIPTablesUtilChains:           utilpointer.Bool(false),
				IPTablesMasqueradeBit:            utilpointer.Int32(0),
				IPTablesDropBit:                  utilpointer.Int32(0),
				FeatureGates:                     map[string]bool{},
				FailSwapOn:                       utilpointer.Bool(false),
				MemorySwap:                       v1beta1.MemorySwapConfiguration{SwapBehavior: ""},
				ContainerLogMaxSize:              "",
				ContainerLogMaxFiles:             utilpointer.Int32(0),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.WatchChangeDetectionStrategy,
				SystemReserved:              map[string]string{},
				KubeReserved:                map[string]string{},
				ReservedSystemCPUs:          "",
				ShowHiddenMetricsForVersion: "",
				SystemReservedCgroup:        "",
				KubeReservedCgroup:          "",
				EnforceNodeAllocatable:      []string{},
				AllowedUnsafeSysctls:        []string{},
				VolumePluginDir:             "",
				ProviderID:                  "",
				KernelMemcgNotification:     false,
				Logging: logsapi.LoggingConfiguration{
					Format:         "",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:          utilpointer.Bool(false),
				ShutdownGracePeriod:             zeroDuration,
				ShutdownGracePeriodCriticalPods: zeroDuration,
				ReservedMemory:                  []v1beta1.MemoryReservation{},
				EnableProfilingHandler:          utilpointer.Bool(false),
				EnableDebugFlagsHandler:         utilpointer.Bool(false),
				SeccompDefault:                  utilpointer.Bool(false),
				MemoryThrottlingFactor:          utilpointer.Float64(0),
				RegisterNode:                    utilpointer.BoolPtr(false),
				LocalStorageCapacityIsolation:   utilpointer.BoolPtr(false),
			},
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(false),
				SyncFrequency:      metav1.Duration{Duration: 1 * time.Minute},
				FileCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				HTTPCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				StaticPodURLHeader: map[string][]string{},
				Address:            "0.0.0.0",
				Port:               10250,
				TLSCipherSuites:    []string{},
				Authentication: v1beta1.KubeletAuthentication{
					X509: v1beta1.KubeletX509Authentication{ClientCAFile: ""},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(false),
						CacheTTL: metav1.Duration{Duration: 2 * time.Minute},
					},
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(false)},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeWebhook,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   metav1.Duration{Duration: 5 * time.Minute},
						CacheUnauthorizedTTL: metav1.Duration{Duration: 30 * time.Second},
					},
				},
				RegistryPullQPS:                  utilpointer.Int32(0),
				RegistryBurst:                    10,
				EventRecordQPS:                   utilpointer.Int32(0),
				EventBurst:                       10,
				EnableDebuggingHandlers:          utilpointer.BoolPtr(false),
				HealthzPort:                      utilpointer.Int32(0),
				HealthzBindAddress:               "127.0.0.1",
				OOMScoreAdj:                      utilpointer.Int32(0),
				ClusterDNS:                       []string{},
				PLEGRelistThreshold:              metav1.Duration{Duration: 3 * time.Minute},
				StreamingConnectionIdleTimeout:   metav1.Duration{Duration: 4 * time.Hour},
				NodeStatusUpdateFrequency:        metav1.Duration{Duration: 10 * time.Second},
				NodeStatusReportFrequency:        metav1.Duration{Duration: 5 * time.Minute},
				NodeLeaseDurationSeconds:         40,
				ImageMinimumGCAge:                metav1.Duration{Duration: 2 * time.Minute},
				ImageGCHighThresholdPercent:      utilpointer.Int32(0),
				ImageGCLowThresholdPercent:       utilpointer.Int32(0),
				VolumeStatsAggPeriod:             metav1.Duration{Duration: time.Minute},
				CgroupsPerQOS:                    utilpointer.BoolPtr(false),
				CgroupDriver:                     "cgroupfs",
				CPUManagerPolicy:                 "none",
				CPUManagerPolicyOptions:          map[string]string{},
				CPUManagerReconcilePeriod:        metav1.Duration{Duration: 10 * time.Second},
				MemoryManagerPolicy:              v1beta1.NoneMemoryManagerPolicy,
				TopologyManagerPolicy:            v1beta1.NoneTopologyManagerPolicy,
				TopologyManagerScope:             v1beta1.ContainerTopologyManagerScope,
				QOSReserved:                      map[string]string{},
				RuntimeRequestTimeout:            metav1.Duration{Duration: 2 * time.Minute},
				HairpinMode:                      v1beta1.PromiscuousBridge,
				MaxPods:                          110,
				PodPidsLimit:                     utilpointer.Int64(0),
				ResolverConfig:                   utilpointer.String(""),
				CPUCFSQuota:                      utilpointer.Bool(false),
				CPUCFSQuotaPeriod:                &zeroDuration,
				NodeStatusMaxImages:              utilpointer.Int32(0),
				MaxOpenFiles:                     1000000,
				ContentType:                      "application/vnd.kubernetes.protobuf",
				KubeAPIQPS:                       utilpointer.Int32(0),
				KubeAPIBurst:                     10,
				SerializeImagePulls:              utilpointer.Bool(false),
				EvictionHard:                     map[string]string{},
				EvictionSoft:                     map[string]string{},
				EvictionSoftGracePeriod:          map[string]string{},
				EvictionPressureTransitionPeriod: metav1.Duration{Duration: 5 * time.Minute},
				EvictionMinimumReclaim:           map[string]string{},
				EnableControllerAttachDetach:     utilpointer.Bool(false),
				MakeIPTablesUtilChains:           utilpointer.Bool(false),
				IPTablesMasqueradeBit:            utilpointer.Int32(0),
				IPTablesDropBit:                  utilpointer.Int32(0),
				FeatureGates:                     map[string]bool{},
				FailSwapOn:                       utilpointer.Bool(false),
				MemorySwap:                       v1beta1.MemorySwapConfiguration{SwapBehavior: ""},
				ContainerLogMaxSize:              "10Mi",
				ContainerLogMaxFiles:             utilpointer.Int32(0),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.WatchChangeDetectionStrategy,
				SystemReserved:         map[string]string{},
				KubeReserved:           map[string]string{},
				EnforceNodeAllocatable: []string{},
				AllowedUnsafeSysctls:   []string{},
				VolumePluginDir:        DefaultVolumePluginDir,
				Logging: logsapi.LoggingConfiguration{
					Format:         "text",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:        utilpointer.Bool(false),
				ReservedMemory:                []v1beta1.MemoryReservation{},
				EnableProfilingHandler:        utilpointer.Bool(false),
				EnableDebugFlagsHandler:       utilpointer.Bool(false),
				SeccompDefault:                utilpointer.Bool(false),
				MemoryThrottlingFactor:        utilpointer.Float64(0),
				RegisterNode:                  utilpointer.BoolPtr(false),
				LocalStorageCapacityIsolation: utilpointer.BoolPtr(false),
			},
		},
		{
			"all positive",
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(true),
				StaticPodPath:      "static/pod/path",
				SyncFrequency:      metav1.Duration{Duration: 60 * time.Second},
				FileCheckFrequency: metav1.Duration{Duration: 60 * time.Second},
				HTTPCheckFrequency: metav1.Duration{Duration: 60 * time.Second},
				StaticPodURL:       "static-pod-url",
				StaticPodURLHeader: map[string][]string{"Static-Pod-URL-Header": {"true"}},
				Address:            "192.168.1.2",
				Port:               10250,
				ReadOnlyPort:       10251,
				TLSCertFile:        "tls-cert-file",
				TLSPrivateKeyFile:  "tls-private-key-file",
				TLSCipherSuites:    []string{"TLS_AES_128_GCM_SHA256"},
				TLSMinVersion:      "1.3",
				RotateCertificates: true,
				ServerTLSBootstrap: true,
				Authentication: v1beta1.KubeletAuthentication{
					X509: v1beta1.KubeletX509Authentication{ClientCAFile: "client-ca-file"},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(true),
						CacheTTL: metav1.Duration{Duration: 60 * time.Second},
					},
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(true)},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeAlwaysAllow,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   metav1.Duration{Duration: 60 * time.Second},
						CacheUnauthorizedTTL: metav1.Duration{Duration: 60 * time.Second},
					},
				},
				RegistryPullQPS:                utilpointer.Int32(1),
				RegistryBurst:                  1,
				EventRecordQPS:                 utilpointer.Int32(1),
				EventBurst:                     1,
				EnableDebuggingHandlers:        utilpointer.BoolPtr(true),
				EnableContentionProfiling:      true,
				HealthzPort:                    utilpointer.Int32(1),
				HealthzBindAddress:             "127.0.0.2",
				OOMScoreAdj:                    utilpointer.Int32(1),
				ClusterDomain:                  "cluster-domain",
				ClusterDNS:                     []string{"192.168.1.3"},
				PLEGRelistThreshold:            metav1.Duration{Duration: 60 * time.Second},
				StreamingConnectionIdleTimeout: metav1.Duration{Duration: 60 * time.Second},
				NodeStatusUpdateFrequency:      metav1.Duration{Duration: 60 * time.Second},
				NodeStatusReportFrequency:      metav1.Duration{Duration: 60 * time.Second},
				NodeLeaseDurationSeconds:       1,
				ImageMinimumGCAge:              metav1.Duration{Duration: 60 * time.Second},
				ImageGCHighThresholdPercent:    utilpointer.Int32(1),
				ImageGCLowThresholdPercent:     utilpointer.Int32(1),
				VolumeStatsAggPeriod:           metav1.Duration{Duration: 60 * time.Second},
				KubeletCgroups:                 "kubelet-cgroup",
				SystemCgroups:                  "system-cgroup",
				CgroupRoot:                     "root-cgroup",
				CgroupsPerQOS:                  utilpointer.BoolPtr(true),
				CgroupDriver:                   "systemd",
				CPUManagerPolicy:               "cpu-manager-policy",
				CPUManagerPolicyOptions:        map[string]string{"key": "value"},
				CPUManagerReconcilePeriod:      metav1.Duration{Duration: 60 * time.Second},
				MemoryManagerPolicy:            v1beta1.StaticMemoryManagerPolicy,
				TopologyManagerPolicy:          v1beta1.RestrictedTopologyManagerPolicy,
				TopologyManagerScope:           v1beta1.PodTopologyManagerScope,
				QOSReserved:                    map[string]string{"memory": "10%"},
				RuntimeRequestTimeout:          metav1.Duration{Duration: 60 * time.Second},
				HairpinMode:                    v1beta1.HairpinVeth,
				MaxPods:                        1,
				PodCIDR:                        "192.168.1.0/24",
				PodPidsLimit:                   utilpointer.Int64(1),
				ResolverConfig:                 utilpointer.String("resolver-config"),
				RunOnce:                        true,
				CPUCFSQuota:                    utilpointer.Bool(true),
				CPUCFSQuotaPeriod:              &metav1.Duration{Duration: 60 * time.Second},
				NodeStatusMaxImages:            utilpointer.Int32(1),
				MaxOpenFiles:                   1,
				ContentType:                    "application/protobuf",
				KubeAPIQPS:                     utilpointer.Int32(1),
				KubeAPIBurst:                   1,
				SerializeImagePulls:            utilpointer.Bool(true),
				EvictionHard: map[string]string{
					"memory.available":  "1Mi",
					"nodefs.available":  "1%",
					"imagefs.available": "1%",
				},
				EvictionSoft: map[string]string{
					"memory.available":  "2Mi",
					"nodefs.available":  "2%",
					"imagefs.available": "2%",
				},
				EvictionSoftGracePeriod: map[string]string{
					"memory.available":  "60s",
					"nodefs.available":  "60s",
					"imagefs.available": "60s",
				},
				EvictionPressureTransitionPeriod: metav1.Duration{Duration: 60 * time.Second},
				EvictionMaxPodGracePeriod:        1,
				EvictionMinimumReclaim: map[string]string{
					"imagefs.available": "1Gi",
				},
				PodsPerCore:                               1,
				EnableControllerAttachDetach:              utilpointer.Bool(true),
				ProtectKernelDefaults:                     true,
				MakeIPTablesUtilChains:                    utilpointer.Bool(true),
				IPTablesMasqueradeBit:                     utilpointer.Int32(1),
				IPTablesDropBit:                           utilpointer.Int32(1),
				FailSwapOn:                                utilpointer.Bool(true),
				MemorySwap:                                v1beta1.MemorySwapConfiguration{SwapBehavior: "UnlimitedSwap"},
				ContainerLogMaxSize:                       "1Mi",
				ContainerLogMaxFiles:                      utilpointer.Int32(1),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.TTLCacheChangeDetectionStrategy,
				SystemReserved: map[string]string{
					"memory": "1Gi",
				},
				KubeReserved: map[string]string{
					"memory": "1Gi",
				},
				ReservedSystemCPUs:          "0,1",
				ShowHiddenMetricsForVersion: "1.16",
				SystemReservedCgroup:        "system-reserved-cgroup",
				KubeReservedCgroup:          "kube-reserved-cgroup",
				EnforceNodeAllocatable:      []string{"system-reserved"},
				AllowedUnsafeSysctls:        []string{"kernel.msg*"},
				VolumePluginDir:             "volume-plugin-dir",
				ProviderID:                  "provider-id",
				KernelMemcgNotification:     true,
				Logging: logsapi.LoggingConfiguration{
					Format:         "json",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:          utilpointer.Bool(true),
				ShutdownGracePeriod:             metav1.Duration{Duration: 60 * time.Second},
				ShutdownGracePeriodCriticalPods: metav1.Duration{Duration: 60 * time.Second},
				ReservedMemory: []v1beta1.MemoryReservation{
					{
						NumaNode: 1,
						Limits:   v1.ResourceList{v1.ResourceMemory: resource.MustParse("1Gi")},
					},
				},
				EnableProfilingHandler:        utilpointer.Bool(true),
				EnableDebugFlagsHandler:       utilpointer.Bool(true),
				SeccompDefault:                utilpointer.Bool(true),
				MemoryThrottlingFactor:        utilpointer.Float64(1),
				RegisterNode:                  utilpointer.BoolPtr(true),
				LocalStorageCapacityIsolation: utilpointer.BoolPtr(true),
			},
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(true),
				StaticPodPath:      "static/pod/path",
				SyncFrequency:      metav1.Duration{Duration: 60 * time.Second},
				FileCheckFrequency: metav1.Duration{Duration: 60 * time.Second},
				HTTPCheckFrequency: metav1.Duration{Duration: 60 * time.Second},
				StaticPodURL:       "static-pod-url",
				StaticPodURLHeader: map[string][]string{"Static-Pod-URL-Header": {"true"}},
				Address:            "192.168.1.2",
				Port:               10250,
				ReadOnlyPort:       10251,
				TLSCertFile:        "tls-cert-file",
				TLSPrivateKeyFile:  "tls-private-key-file",
				TLSCipherSuites:    []string{"TLS_AES_128_GCM_SHA256"},
				TLSMinVersion:      "1.3",
				RotateCertificates: true,
				ServerTLSBootstrap: true,
				Authentication: v1beta1.KubeletAuthentication{
					X509: v1beta1.KubeletX509Authentication{ClientCAFile: "client-ca-file"},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(true),
						CacheTTL: metav1.Duration{Duration: 60 * time.Second},
					},
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(true)},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeAlwaysAllow,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   metav1.Duration{Duration: 60 * time.Second},
						CacheUnauthorizedTTL: metav1.Duration{Duration: 60 * time.Second},
					},
				},
				RegistryPullQPS:                utilpointer.Int32(1),
				RegistryBurst:                  1,
				EventRecordQPS:                 utilpointer.Int32(1),
				EventBurst:                     1,
				EnableDebuggingHandlers:        utilpointer.BoolPtr(true),
				EnableContentionProfiling:      true,
				HealthzPort:                    utilpointer.Int32(1),
				HealthzBindAddress:             "127.0.0.2",
				OOMScoreAdj:                    utilpointer.Int32(1),
				ClusterDomain:                  "cluster-domain",
				ClusterDNS:                     []string{"192.168.1.3"},
				PLEGRelistThreshold:            metav1.Duration{Duration: 60 * time.Second},
				StreamingConnectionIdleTimeout: metav1.Duration{Duration: 60 * time.Second},
				NodeStatusUpdateFrequency:      metav1.Duration{Duration: 60 * time.Second},
				NodeStatusReportFrequency:      metav1.Duration{Duration: 60 * time.Second},
				NodeLeaseDurationSeconds:       1,
				ImageMinimumGCAge:              metav1.Duration{Duration: 60 * time.Second},
				ImageGCHighThresholdPercent:    utilpointer.Int32(1),
				ImageGCLowThresholdPercent:     utilpointer.Int32(1),
				VolumeStatsAggPeriod:           metav1.Duration{Duration: 60 * time.Second},
				KubeletCgroups:                 "kubelet-cgroup",
				SystemCgroups:                  "system-cgroup",
				CgroupRoot:                     "root-cgroup",
				CgroupsPerQOS:                  utilpointer.BoolPtr(true),
				CgroupDriver:                   "systemd",
				CPUManagerPolicy:               "cpu-manager-policy",
				CPUManagerPolicyOptions:        map[string]string{"key": "value"},
				CPUManagerReconcilePeriod:      metav1.Duration{Duration: 60 * time.Second},
				MemoryManagerPolicy:            v1beta1.StaticMemoryManagerPolicy,
				TopologyManagerPolicy:          v1beta1.RestrictedTopologyManagerPolicy,
				TopologyManagerScope:           v1beta1.PodTopologyManagerScope,
				QOSReserved:                    map[string]string{"memory": "10%"},
				RuntimeRequestTimeout:          metav1.Duration{Duration: 60 * time.Second},
				HairpinMode:                    v1beta1.HairpinVeth,
				MaxPods:                        1,
				PodCIDR:                        "192.168.1.0/24",
				PodPidsLimit:                   utilpointer.Int64(1),
				ResolverConfig:                 utilpointer.String("resolver-config"),
				RunOnce:                        true,
				CPUCFSQuota:                    utilpointer.Bool(true),
				CPUCFSQuotaPeriod:              &metav1.Duration{Duration: 60 * time.Second},
				NodeStatusMaxImages:            utilpointer.Int32(1),
				MaxOpenFiles:                   1,
				ContentType:                    "application/protobuf",
				KubeAPIQPS:                     utilpointer.Int32(1),
				KubeAPIBurst:                   1,
				SerializeImagePulls:            utilpointer.Bool(true),
				EvictionHard: map[string]string{
					"memory.available":  "1Mi",
					"nodefs.available":  "1%",
					"imagefs.available": "1%",
				},
				EvictionSoft: map[string]string{
					"memory.available":  "2Mi",
					"nodefs.available":  "2%",
					"imagefs.available": "2%",
				},
				EvictionSoftGracePeriod: map[string]string{
					"memory.available":  "60s",
					"nodefs.available":  "60s",
					"imagefs.available": "60s",
				},
				EvictionPressureTransitionPeriod: metav1.Duration{Duration: 60 * time.Second},
				EvictionMaxPodGracePeriod:        1,
				EvictionMinimumReclaim: map[string]string{
					"imagefs.available": "1Gi",
				},
				PodsPerCore:                               1,
				EnableControllerAttachDetach:              utilpointer.Bool(true),
				ProtectKernelDefaults:                     true,
				MakeIPTablesUtilChains:                    utilpointer.Bool(true),
				IPTablesMasqueradeBit:                     utilpointer.Int32(1),
				IPTablesDropBit:                           utilpointer.Int32(1),
				FailSwapOn:                                utilpointer.Bool(true),
				MemorySwap:                                v1beta1.MemorySwapConfiguration{SwapBehavior: "UnlimitedSwap"},
				ContainerLogMaxSize:                       "1Mi",
				ContainerLogMaxFiles:                      utilpointer.Int32(1),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.TTLCacheChangeDetectionStrategy,
				SystemReserved: map[string]string{
					"memory": "1Gi",
				},
				KubeReserved: map[string]string{
					"memory": "1Gi",
				},
				ReservedSystemCPUs:          "0,1",
				ShowHiddenMetricsForVersion: "1.16",
				SystemReservedCgroup:        "system-reserved-cgroup",
				KubeReservedCgroup:          "kube-reserved-cgroup",
				EnforceNodeAllocatable:      []string{"system-reserved"},
				AllowedUnsafeSysctls:        []string{"kernel.msg*"},
				VolumePluginDir:             "volume-plugin-dir",
				ProviderID:                  "provider-id",
				KernelMemcgNotification:     true,
				Logging: logsapi.LoggingConfiguration{
					Format:         "json",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:          utilpointer.Bool(true),
				ShutdownGracePeriod:             metav1.Duration{Duration: 60 * time.Second},
				ShutdownGracePeriodCriticalPods: metav1.Duration{Duration: 60 * time.Second},
				ReservedMemory: []v1beta1.MemoryReservation{
					{
						NumaNode: 1,
						Limits:   v1.ResourceList{v1.ResourceMemory: resource.MustParse("1Gi")},
					},
				},
				EnableProfilingHandler:        utilpointer.Bool(true),
				EnableDebugFlagsHandler:       utilpointer.Bool(true),
				SeccompDefault:                utilpointer.Bool(true),
				MemoryThrottlingFactor:        utilpointer.Float64(1),
				RegisterNode:                  utilpointer.BoolPtr(true),
				LocalStorageCapacityIsolation: utilpointer.BoolPtr(true),
			},
		},
		{
			"NodeStatusUpdateFrequency is not zero",
			&v1beta1.KubeletConfiguration{
				NodeStatusUpdateFrequency: metav1.Duration{Duration: 1 * time.Minute},
			},
			&v1beta1.KubeletConfiguration{
				EnableServer:       utilpointer.BoolPtr(true),
				SyncFrequency:      metav1.Duration{Duration: 1 * time.Minute},
				FileCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				HTTPCheckFrequency: metav1.Duration{Duration: 20 * time.Second},
				Address:            "0.0.0.0",
				Port:               ports.KubeletPort,
				Authentication: v1beta1.KubeletAuthentication{
					Anonymous: v1beta1.KubeletAnonymousAuthentication{Enabled: utilpointer.BoolPtr(false)},
					Webhook: v1beta1.KubeletWebhookAuthentication{
						Enabled:  utilpointer.BoolPtr(true),
						CacheTTL: metav1.Duration{Duration: 2 * time.Minute},
					},
				},
				Authorization: v1beta1.KubeletAuthorization{
					Mode: v1beta1.KubeletAuthorizationModeWebhook,
					Webhook: v1beta1.KubeletWebhookAuthorization{
						CacheAuthorizedTTL:   metav1.Duration{Duration: 5 * time.Minute},
						CacheUnauthorizedTTL: metav1.Duration{Duration: 30 * time.Second},
					},
				},
				RegistryPullQPS:                           utilpointer.Int32Ptr(5),
				RegistryBurst:                             10,
				EventRecordQPS:                            utilpointer.Int32Ptr(5),
				EventBurst:                                10,
				EnableDebuggingHandlers:                   utilpointer.BoolPtr(true),
				HealthzPort:                               utilpointer.Int32Ptr(10248),
				HealthzBindAddress:                        "127.0.0.1",
				OOMScoreAdj:                               utilpointer.Int32Ptr(int32(qos.KubeletOOMScoreAdj)),
				PLEGRelistThreshold:                       metav1.Duration{Duration: 3 * time.Minute},
				StreamingConnectionIdleTimeout:            metav1.Duration{Duration: 4 * time.Hour},
				NodeStatusUpdateFrequency:                 metav1.Duration{Duration: 1 * time.Minute},
				NodeStatusReportFrequency:                 metav1.Duration{Duration: 1 * time.Minute},
				NodeLeaseDurationSeconds:                  40,
				ImageMinimumGCAge:                         metav1.Duration{Duration: 2 * time.Minute},
				ImageGCHighThresholdPercent:               utilpointer.Int32Ptr(85),
				ImageGCLowThresholdPercent:                utilpointer.Int32Ptr(80),
				VolumeStatsAggPeriod:                      metav1.Duration{Duration: time.Minute},
				CgroupsPerQOS:                             utilpointer.BoolPtr(true),
				CgroupDriver:                              "cgroupfs",
				CPUManagerPolicy:                          "none",
				CPUManagerReconcilePeriod:                 metav1.Duration{Duration: 10 * time.Second},
				MemoryManagerPolicy:                       v1beta1.NoneMemoryManagerPolicy,
				TopologyManagerPolicy:                     v1beta1.NoneTopologyManagerPolicy,
				TopologyManagerScope:                      v1beta1.ContainerTopologyManagerScope,
				RuntimeRequestTimeout:                     metav1.Duration{Duration: 2 * time.Minute},
				HairpinMode:                               v1beta1.PromiscuousBridge,
				MaxPods:                                   110,
				PodPidsLimit:                              utilpointer.Int64(-1),
				ResolverConfig:                            utilpointer.String(kubetypes.ResolvConfDefault),
				CPUCFSQuota:                               utilpointer.BoolPtr(true),
				CPUCFSQuotaPeriod:                         &metav1.Duration{Duration: 100 * time.Millisecond},
				NodeStatusMaxImages:                       utilpointer.Int32Ptr(50),
				MaxOpenFiles:                              1000000,
				ContentType:                               "application/vnd.kubernetes.protobuf",
				KubeAPIQPS:                                utilpointer.Int32Ptr(5),
				KubeAPIBurst:                              10,
				SerializeImagePulls:                       utilpointer.BoolPtr(true),
				EvictionHard:                              DefaultEvictionHard,
				EvictionPressureTransitionPeriod:          metav1.Duration{Duration: 5 * time.Minute},
				EnableControllerAttachDetach:              utilpointer.BoolPtr(true),
				MakeIPTablesUtilChains:                    utilpointer.BoolPtr(true),
				IPTablesMasqueradeBit:                     utilpointer.Int32Ptr(DefaultIPTablesMasqueradeBit),
				IPTablesDropBit:                           utilpointer.Int32Ptr(DefaultIPTablesDropBit),
				FailSwapOn:                                utilpointer.BoolPtr(true),
				ContainerLogMaxSize:                       "10Mi",
				ContainerLogMaxFiles:                      utilpointer.Int32Ptr(5),
				ConfigMapAndSecretChangeDetectionStrategy: v1beta1.WatchChangeDetectionStrategy,
				EnforceNodeAllocatable:                    DefaultNodeAllocatableEnforcement,
				VolumePluginDir:                           DefaultVolumePluginDir,
				Logging: logsapi.LoggingConfiguration{
					Format:         "text",
					FlushFrequency: 5 * time.Second,
				},
				EnableSystemLogHandler:        utilpointer.BoolPtr(true),
				EnableProfilingHandler:        utilpointer.BoolPtr(true),
				EnableDebugFlagsHandler:       utilpointer.BoolPtr(true),
				SeccompDefault:                utilpointer.BoolPtr(false),
				MemoryThrottlingFactor:        utilpointer.Float64Ptr(DefaultMemoryThrottlingFactor),
				RegisterNode:                  utilpointer.BoolPtr(true),
				LocalStorageCapacityIsolation: utilpointer.BoolPtr(true),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			SetDefaults_KubeletConfiguration(tc.config)
			if diff := cmp.Diff(tc.expected, tc.config); diff != "" {
				t.Errorf("Got unexpected defaults (-want, +got):\n%s", diff)
			}
		})
	}
}
