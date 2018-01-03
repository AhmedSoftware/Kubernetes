// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAnonymousAuthentication) DeepCopyInto(out *KubeletAnonymousAuthentication) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAnonymousAuthentication.
func (in *KubeletAnonymousAuthentication) DeepCopy() *KubeletAnonymousAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletAnonymousAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAuthentication) DeepCopyInto(out *KubeletAuthentication) {
	*out = *in
	out.X509 = in.X509
	in.Webhook.DeepCopyInto(&out.Webhook)
	in.Anonymous.DeepCopyInto(&out.Anonymous)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAuthentication.
func (in *KubeletAuthentication) DeepCopy() *KubeletAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletAuthorization) DeepCopyInto(out *KubeletAuthorization) {
	*out = *in
	out.Webhook = in.Webhook
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletAuthorization.
func (in *KubeletAuthorization) DeepCopy() *KubeletAuthorization {
	if in == nil {
		return nil
	}
	out := new(KubeletAuthorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletConfiguration) DeepCopyInto(out *KubeletConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.ConfigTrialDuration != nil {
		in, out := &in.ConfigTrialDuration, &out.ConfigTrialDuration
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.Duration)
			**out = **in
		}
	}
	out.SyncFrequency = in.SyncFrequency
	out.FileCheckFrequency = in.FileCheckFrequency
	out.HTTPCheckFrequency = in.HTTPCheckFrequency
	if in.ManifestURLHeader != nil {
		in, out := &in.ManifestURLHeader, &out.ManifestURLHeader
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			if val == nil {
				(*out)[key] = nil
			} else {
				(*out)[key] = make([]string, len(val))
				copy((*out)[key], val)
			}
		}
	}
	if in.EnableServer != nil {
		in, out := &in.EnableServer, &out.EnableServer
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.ReadOnlyPort != nil {
		in, out := &in.ReadOnlyPort, &out.ReadOnlyPort
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	in.Authentication.DeepCopyInto(&out.Authentication)
	out.Authorization = in.Authorization
	if in.AllowPrivileged != nil {
		in, out := &in.AllowPrivileged, &out.AllowPrivileged
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.HostNetworkSources != nil {
		in, out := &in.HostNetworkSources, &out.HostNetworkSources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.HostPIDSources != nil {
		in, out := &in.HostPIDSources, &out.HostPIDSources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.HostIPCSources != nil {
		in, out := &in.HostIPCSources, &out.HostIPCSources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RegistryPullQPS != nil {
		in, out := &in.RegistryPullQPS, &out.RegistryPullQPS
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.EventRecordQPS != nil {
		in, out := &in.EventRecordQPS, &out.EventRecordQPS
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.EnableDebuggingHandlers != nil {
		in, out := &in.EnableDebuggingHandlers, &out.EnableDebuggingHandlers
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.CAdvisorAllowDynamicHousekeeping != nil {
		in, out := &in.CAdvisorAllowDynamicHousekeeping, &out.CAdvisorAllowDynamicHousekeeping
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.CAdvisorPort != nil {
		in, out := &in.CAdvisorPort, &out.CAdvisorPort
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.HealthzPort != nil {
		in, out := &in.HealthzPort, &out.HealthzPort
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.OOMScoreAdj != nil {
		in, out := &in.OOMScoreAdj, &out.OOMScoreAdj
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.ClusterDNS != nil {
		in, out := &in.ClusterDNS, &out.ClusterDNS
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.StreamingConnectionIdleTimeout = in.StreamingConnectionIdleTimeout
	out.NodeStatusUpdateFrequency = in.NodeStatusUpdateFrequency
	out.ImageMinimumGCAge = in.ImageMinimumGCAge
	if in.ImageGCHighThresholdPercent != nil {
		in, out := &in.ImageGCHighThresholdPercent, &out.ImageGCHighThresholdPercent
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.ImageGCLowThresholdPercent != nil {
		in, out := &in.ImageGCLowThresholdPercent, &out.ImageGCLowThresholdPercent
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	out.VolumeStatsAggPeriod = in.VolumeStatsAggPeriod
	if in.CgroupsPerQOS != nil {
		in, out := &in.CgroupsPerQOS, &out.CgroupsPerQOS
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	out.CPUManagerReconcilePeriod = in.CPUManagerReconcilePeriod
	out.RuntimeRequestTimeout = in.RuntimeRequestTimeout
	if in.CPUCFSQuota != nil {
		in, out := &in.CPUCFSQuota, &out.CPUCFSQuota
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.KubeAPIQPS != nil {
		in, out := &in.KubeAPIQPS, &out.KubeAPIQPS
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.SerializeImagePulls != nil {
		in, out := &in.SerializeImagePulls, &out.SerializeImagePulls
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.EvictionHard != nil {
		in, out := &in.EvictionHard, &out.EvictionHard
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EvictionSoft != nil {
		in, out := &in.EvictionSoft, &out.EvictionSoft
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EvictionSoftGracePeriod != nil {
		in, out := &in.EvictionSoftGracePeriod, &out.EvictionSoftGracePeriod
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	out.EvictionPressureTransitionPeriod = in.EvictionPressureTransitionPeriod
	if in.EvictionMinimumReclaim != nil {
		in, out := &in.EvictionMinimumReclaim, &out.EvictionMinimumReclaim
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EnableControllerAttachDetach != nil {
		in, out := &in.EnableControllerAttachDetach, &out.EnableControllerAttachDetach
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.MakeIPTablesUtilChains != nil {
		in, out := &in.MakeIPTablesUtilChains, &out.MakeIPTablesUtilChains
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	if in.IPTablesMasqueradeBit != nil {
		in, out := &in.IPTablesMasqueradeBit, &out.IPTablesMasqueradeBit
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.IPTablesDropBit != nil {
		in, out := &in.IPTablesDropBit, &out.IPTablesDropBit
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SystemReserved != nil {
		in, out := &in.SystemReserved, &out.SystemReserved
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.KubeReserved != nil {
		in, out := &in.KubeReserved, &out.KubeReserved
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.EnforceNodeAllocatable != nil {
		in, out := &in.EnforceNodeAllocatable, &out.EnforceNodeAllocatable
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletConfiguration.
func (in *KubeletConfiguration) DeepCopy() *KubeletConfiguration {
	if in == nil {
		return nil
	}
	out := new(KubeletConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeletConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletWebhookAuthentication) DeepCopyInto(out *KubeletWebhookAuthentication) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	out.CacheTTL = in.CacheTTL
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletWebhookAuthentication.
func (in *KubeletWebhookAuthentication) DeepCopy() *KubeletWebhookAuthentication {
	if in == nil {
		return nil
	}
	out := new(KubeletWebhookAuthentication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletWebhookAuthorization) DeepCopyInto(out *KubeletWebhookAuthorization) {
	*out = *in
	out.CacheAuthorizedTTL = in.CacheAuthorizedTTL
	out.CacheUnauthorizedTTL = in.CacheUnauthorizedTTL
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletWebhookAuthorization.
func (in *KubeletWebhookAuthorization) DeepCopy() *KubeletWebhookAuthorization {
	if in == nil {
		return nil
	}
	out := new(KubeletWebhookAuthorization)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeletX509Authentication) DeepCopyInto(out *KubeletX509Authentication) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeletX509Authentication.
func (in *KubeletX509Authentication) DeepCopy() *KubeletX509Authentication {
	if in == nil {
		return nil
	}
	out := new(KubeletX509Authentication)
	in.DeepCopyInto(out)
	return out
}
