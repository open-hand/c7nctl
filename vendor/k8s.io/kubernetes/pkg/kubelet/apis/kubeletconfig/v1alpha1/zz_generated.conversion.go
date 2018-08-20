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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/kubeletconfig"
	unsafe "unsafe"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication,
		Convert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication,
		Convert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication,
		Convert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication,
		Convert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization,
		Convert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization,
		Convert_v1alpha1_KubeletConfiguration_To_kubeletconfig_KubeletConfiguration,
		Convert_kubeletconfig_KubeletConfiguration_To_v1alpha1_KubeletConfiguration,
		Convert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication,
		Convert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication,
		Convert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization,
		Convert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization,
		Convert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication,
		Convert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication,
	)
}

func autoConvert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication(in *KubeletAnonymousAuthentication, out *kubeletconfig.KubeletAnonymousAuthentication, s conversion.Scope) error {
	if err := v1.Convert_Pointer_bool_To_bool(&in.Enabled, &out.Enabled, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication is an autogenerated conversion function.
func Convert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication(in *KubeletAnonymousAuthentication, out *kubeletconfig.KubeletAnonymousAuthentication, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication(in, out, s)
}

func autoConvert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication(in *kubeletconfig.KubeletAnonymousAuthentication, out *KubeletAnonymousAuthentication, s conversion.Scope) error {
	if err := v1.Convert_bool_To_Pointer_bool(&in.Enabled, &out.Enabled, s); err != nil {
		return err
	}
	return nil
}

// Convert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication(in *kubeletconfig.KubeletAnonymousAuthentication, out *KubeletAnonymousAuthentication, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication(in, out, s)
}

func autoConvert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication(in *KubeletAuthentication, out *kubeletconfig.KubeletAuthentication, s conversion.Scope) error {
	if err := Convert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication(&in.X509, &out.X509, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication(&in.Webhook, &out.Webhook, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_KubeletAnonymousAuthentication_To_kubeletconfig_KubeletAnonymousAuthentication(&in.Anonymous, &out.Anonymous, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication is an autogenerated conversion function.
func Convert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication(in *KubeletAuthentication, out *kubeletconfig.KubeletAuthentication, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication(in, out, s)
}

func autoConvert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication(in *kubeletconfig.KubeletAuthentication, out *KubeletAuthentication, s conversion.Scope) error {
	if err := Convert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication(&in.X509, &out.X509, s); err != nil {
		return err
	}
	if err := Convert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication(&in.Webhook, &out.Webhook, s); err != nil {
		return err
	}
	if err := Convert_kubeletconfig_KubeletAnonymousAuthentication_To_v1alpha1_KubeletAnonymousAuthentication(&in.Anonymous, &out.Anonymous, s); err != nil {
		return err
	}
	return nil
}

// Convert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication(in *kubeletconfig.KubeletAuthentication, out *KubeletAuthentication, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication(in, out, s)
}

func autoConvert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization(in *KubeletAuthorization, out *kubeletconfig.KubeletAuthorization, s conversion.Scope) error {
	out.Mode = kubeletconfig.KubeletAuthorizationMode(in.Mode)
	if err := Convert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization(&in.Webhook, &out.Webhook, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization is an autogenerated conversion function.
func Convert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization(in *KubeletAuthorization, out *kubeletconfig.KubeletAuthorization, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization(in, out, s)
}

func autoConvert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization(in *kubeletconfig.KubeletAuthorization, out *KubeletAuthorization, s conversion.Scope) error {
	out.Mode = KubeletAuthorizationMode(in.Mode)
	if err := Convert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization(&in.Webhook, &out.Webhook, s); err != nil {
		return err
	}
	return nil
}

// Convert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization(in *kubeletconfig.KubeletAuthorization, out *KubeletAuthorization, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization(in, out, s)
}

func autoConvert_v1alpha1_KubeletConfiguration_To_kubeletconfig_KubeletConfiguration(in *KubeletConfiguration, out *kubeletconfig.KubeletConfiguration, s conversion.Scope) error {
	out.ConfigTrialDuration = (*v1.Duration)(unsafe.Pointer(in.ConfigTrialDuration))
	out.PodManifestPath = in.PodManifestPath
	out.SyncFrequency = in.SyncFrequency
	out.FileCheckFrequency = in.FileCheckFrequency
	out.HTTPCheckFrequency = in.HTTPCheckFrequency
	out.ManifestURL = in.ManifestURL
	out.ManifestURLHeader = *(*map[string][]string)(unsafe.Pointer(&in.ManifestURLHeader))
	if err := v1.Convert_Pointer_bool_To_bool(&in.EnableServer, &out.EnableServer, s); err != nil {
		return err
	}
	out.Address = in.Address
	out.Port = in.Port
	if err := v1.Convert_Pointer_int32_To_int32(&in.ReadOnlyPort, &out.ReadOnlyPort, s); err != nil {
		return err
	}
	out.TLSCertFile = in.TLSCertFile
	out.TLSPrivateKeyFile = in.TLSPrivateKeyFile
	if err := Convert_v1alpha1_KubeletAuthentication_To_kubeletconfig_KubeletAuthentication(&in.Authentication, &out.Authentication, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_KubeletAuthorization_To_kubeletconfig_KubeletAuthorization(&in.Authorization, &out.Authorization, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_bool_To_bool(&in.AllowPrivileged, &out.AllowPrivileged, s); err != nil {
		return err
	}
	out.HostNetworkSources = *(*[]string)(unsafe.Pointer(&in.HostNetworkSources))
	out.HostPIDSources = *(*[]string)(unsafe.Pointer(&in.HostPIDSources))
	out.HostIPCSources = *(*[]string)(unsafe.Pointer(&in.HostIPCSources))
	if err := v1.Convert_Pointer_int32_To_int32(&in.RegistryPullQPS, &out.RegistryPullQPS, s); err != nil {
		return err
	}
	out.RegistryBurst = in.RegistryBurst
	if err := v1.Convert_Pointer_int32_To_int32(&in.EventRecordQPS, &out.EventRecordQPS, s); err != nil {
		return err
	}
	out.EventBurst = in.EventBurst
	if err := v1.Convert_Pointer_bool_To_bool(&in.EnableDebuggingHandlers, &out.EnableDebuggingHandlers, s); err != nil {
		return err
	}
	out.EnableContentionProfiling = in.EnableContentionProfiling
	if err := v1.Convert_Pointer_int32_To_int32(&in.CAdvisorPort, &out.CAdvisorPort, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_int32_To_int32(&in.HealthzPort, &out.HealthzPort, s); err != nil {
		return err
	}
	out.HealthzBindAddress = in.HealthzBindAddress
	if err := v1.Convert_Pointer_int32_To_int32(&in.OOMScoreAdj, &out.OOMScoreAdj, s); err != nil {
		return err
	}
	out.ClusterDomain = in.ClusterDomain
	out.ClusterDNS = *(*[]string)(unsafe.Pointer(&in.ClusterDNS))
	out.StreamingConnectionIdleTimeout = in.StreamingConnectionIdleTimeout
	out.NodeStatusUpdateFrequency = in.NodeStatusUpdateFrequency
	out.ImageMinimumGCAge = in.ImageMinimumGCAge
	if err := v1.Convert_Pointer_int32_To_int32(&in.ImageGCHighThresholdPercent, &out.ImageGCHighThresholdPercent, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_int32_To_int32(&in.ImageGCLowThresholdPercent, &out.ImageGCLowThresholdPercent, s); err != nil {
		return err
	}
	out.VolumeStatsAggPeriod = in.VolumeStatsAggPeriod
	out.KubeletCgroups = in.KubeletCgroups
	out.SystemCgroups = in.SystemCgroups
	out.CgroupRoot = in.CgroupRoot
	if err := v1.Convert_Pointer_bool_To_bool(&in.CgroupsPerQOS, &out.CgroupsPerQOS, s); err != nil {
		return err
	}
	out.CgroupDriver = in.CgroupDriver
	out.CPUManagerPolicy = in.CPUManagerPolicy
	out.CPUManagerReconcilePeriod = in.CPUManagerReconcilePeriod
	out.RuntimeRequestTimeout = in.RuntimeRequestTimeout
	out.HairpinMode = in.HairpinMode
	out.MaxPods = in.MaxPods
	out.PodCIDR = in.PodCIDR
	out.ResolverConfig = in.ResolverConfig
	if err := v1.Convert_Pointer_bool_To_bool(&in.CPUCFSQuota, &out.CPUCFSQuota, s); err != nil {
		return err
	}
	out.MaxOpenFiles = in.MaxOpenFiles
	out.ContentType = in.ContentType
	if err := v1.Convert_Pointer_int32_To_int32(&in.KubeAPIQPS, &out.KubeAPIQPS, s); err != nil {
		return err
	}
	out.KubeAPIBurst = in.KubeAPIBurst
	if err := v1.Convert_Pointer_bool_To_bool(&in.SerializeImagePulls, &out.SerializeImagePulls, s); err != nil {
		return err
	}
	out.EvictionHard = *(*map[string]string)(unsafe.Pointer(&in.EvictionHard))
	out.EvictionSoft = *(*map[string]string)(unsafe.Pointer(&in.EvictionSoft))
	out.EvictionSoftGracePeriod = *(*map[string]string)(unsafe.Pointer(&in.EvictionSoftGracePeriod))
	out.EvictionPressureTransitionPeriod = in.EvictionPressureTransitionPeriod
	out.EvictionMaxPodGracePeriod = in.EvictionMaxPodGracePeriod
	out.EvictionMinimumReclaim = *(*map[string]string)(unsafe.Pointer(&in.EvictionMinimumReclaim))
	out.PodsPerCore = in.PodsPerCore
	if err := v1.Convert_Pointer_bool_To_bool(&in.EnableControllerAttachDetach, &out.EnableControllerAttachDetach, s); err != nil {
		return err
	}
	out.ProtectKernelDefaults = in.ProtectKernelDefaults
	if err := v1.Convert_Pointer_bool_To_bool(&in.MakeIPTablesUtilChains, &out.MakeIPTablesUtilChains, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_int32_To_int32(&in.IPTablesMasqueradeBit, &out.IPTablesMasqueradeBit, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_int32_To_int32(&in.IPTablesDropBit, &out.IPTablesDropBit, s); err != nil {
		return err
	}
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.FailSwapOn = in.FailSwapOn
	out.SystemReserved = *(*map[string]string)(unsafe.Pointer(&in.SystemReserved))
	out.KubeReserved = *(*map[string]string)(unsafe.Pointer(&in.KubeReserved))
	out.SystemReservedCgroup = in.SystemReservedCgroup
	out.KubeReservedCgroup = in.KubeReservedCgroup
	out.EnforceNodeAllocatable = *(*[]string)(unsafe.Pointer(&in.EnforceNodeAllocatable))
	return nil
}

// Convert_v1alpha1_KubeletConfiguration_To_kubeletconfig_KubeletConfiguration is an autogenerated conversion function.
func Convert_v1alpha1_KubeletConfiguration_To_kubeletconfig_KubeletConfiguration(in *KubeletConfiguration, out *kubeletconfig.KubeletConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletConfiguration_To_kubeletconfig_KubeletConfiguration(in, out, s)
}

func autoConvert_kubeletconfig_KubeletConfiguration_To_v1alpha1_KubeletConfiguration(in *kubeletconfig.KubeletConfiguration, out *KubeletConfiguration, s conversion.Scope) error {
	out.ConfigTrialDuration = (*v1.Duration)(unsafe.Pointer(in.ConfigTrialDuration))
	out.PodManifestPath = in.PodManifestPath
	out.SyncFrequency = in.SyncFrequency
	out.FileCheckFrequency = in.FileCheckFrequency
	out.HTTPCheckFrequency = in.HTTPCheckFrequency
	out.ManifestURL = in.ManifestURL
	out.ManifestURLHeader = *(*map[string][]string)(unsafe.Pointer(&in.ManifestURLHeader))
	if err := v1.Convert_bool_To_Pointer_bool(&in.EnableServer, &out.EnableServer, s); err != nil {
		return err
	}
	out.Address = in.Address
	out.Port = in.Port
	if err := v1.Convert_int32_To_Pointer_int32(&in.ReadOnlyPort, &out.ReadOnlyPort, s); err != nil {
		return err
	}
	out.TLSCertFile = in.TLSCertFile
	out.TLSPrivateKeyFile = in.TLSPrivateKeyFile
	if err := Convert_kubeletconfig_KubeletAuthentication_To_v1alpha1_KubeletAuthentication(&in.Authentication, &out.Authentication, s); err != nil {
		return err
	}
	if err := Convert_kubeletconfig_KubeletAuthorization_To_v1alpha1_KubeletAuthorization(&in.Authorization, &out.Authorization, s); err != nil {
		return err
	}
	if err := v1.Convert_bool_To_Pointer_bool(&in.AllowPrivileged, &out.AllowPrivileged, s); err != nil {
		return err
	}
	out.HostNetworkSources = *(*[]string)(unsafe.Pointer(&in.HostNetworkSources))
	out.HostPIDSources = *(*[]string)(unsafe.Pointer(&in.HostPIDSources))
	out.HostIPCSources = *(*[]string)(unsafe.Pointer(&in.HostIPCSources))
	if err := v1.Convert_int32_To_Pointer_int32(&in.RegistryPullQPS, &out.RegistryPullQPS, s); err != nil {
		return err
	}
	out.RegistryBurst = in.RegistryBurst
	if err := v1.Convert_int32_To_Pointer_int32(&in.EventRecordQPS, &out.EventRecordQPS, s); err != nil {
		return err
	}
	out.EventBurst = in.EventBurst
	if err := v1.Convert_bool_To_Pointer_bool(&in.EnableDebuggingHandlers, &out.EnableDebuggingHandlers, s); err != nil {
		return err
	}
	out.EnableContentionProfiling = in.EnableContentionProfiling
	if err := v1.Convert_int32_To_Pointer_int32(&in.CAdvisorPort, &out.CAdvisorPort, s); err != nil {
		return err
	}
	if err := v1.Convert_int32_To_Pointer_int32(&in.HealthzPort, &out.HealthzPort, s); err != nil {
		return err
	}
	out.HealthzBindAddress = in.HealthzBindAddress
	if err := v1.Convert_int32_To_Pointer_int32(&in.OOMScoreAdj, &out.OOMScoreAdj, s); err != nil {
		return err
	}
	out.ClusterDomain = in.ClusterDomain
	out.ClusterDNS = *(*[]string)(unsafe.Pointer(&in.ClusterDNS))
	out.StreamingConnectionIdleTimeout = in.StreamingConnectionIdleTimeout
	out.NodeStatusUpdateFrequency = in.NodeStatusUpdateFrequency
	out.ImageMinimumGCAge = in.ImageMinimumGCAge
	if err := v1.Convert_int32_To_Pointer_int32(&in.ImageGCHighThresholdPercent, &out.ImageGCHighThresholdPercent, s); err != nil {
		return err
	}
	if err := v1.Convert_int32_To_Pointer_int32(&in.ImageGCLowThresholdPercent, &out.ImageGCLowThresholdPercent, s); err != nil {
		return err
	}
	out.VolumeStatsAggPeriod = in.VolumeStatsAggPeriod
	out.KubeletCgroups = in.KubeletCgroups
	if err := v1.Convert_bool_To_Pointer_bool(&in.CgroupsPerQOS, &out.CgroupsPerQOS, s); err != nil {
		return err
	}
	out.CgroupDriver = in.CgroupDriver
	out.SystemCgroups = in.SystemCgroups
	out.CgroupRoot = in.CgroupRoot
	out.CPUManagerPolicy = in.CPUManagerPolicy
	out.CPUManagerReconcilePeriod = in.CPUManagerReconcilePeriod
	out.RuntimeRequestTimeout = in.RuntimeRequestTimeout
	out.HairpinMode = in.HairpinMode
	out.MaxPods = in.MaxPods
	out.PodCIDR = in.PodCIDR
	out.ResolverConfig = in.ResolverConfig
	if err := v1.Convert_bool_To_Pointer_bool(&in.CPUCFSQuota, &out.CPUCFSQuota, s); err != nil {
		return err
	}
	out.MaxOpenFiles = in.MaxOpenFiles
	out.ContentType = in.ContentType
	if err := v1.Convert_int32_To_Pointer_int32(&in.KubeAPIQPS, &out.KubeAPIQPS, s); err != nil {
		return err
	}
	out.KubeAPIBurst = in.KubeAPIBurst
	if err := v1.Convert_bool_To_Pointer_bool(&in.SerializeImagePulls, &out.SerializeImagePulls, s); err != nil {
		return err
	}
	out.EvictionHard = *(*map[string]string)(unsafe.Pointer(&in.EvictionHard))
	out.EvictionSoft = *(*map[string]string)(unsafe.Pointer(&in.EvictionSoft))
	out.EvictionSoftGracePeriod = *(*map[string]string)(unsafe.Pointer(&in.EvictionSoftGracePeriod))
	out.EvictionPressureTransitionPeriod = in.EvictionPressureTransitionPeriod
	out.EvictionMaxPodGracePeriod = in.EvictionMaxPodGracePeriod
	out.EvictionMinimumReclaim = *(*map[string]string)(unsafe.Pointer(&in.EvictionMinimumReclaim))
	out.PodsPerCore = in.PodsPerCore
	if err := v1.Convert_bool_To_Pointer_bool(&in.EnableControllerAttachDetach, &out.EnableControllerAttachDetach, s); err != nil {
		return err
	}
	out.ProtectKernelDefaults = in.ProtectKernelDefaults
	if err := v1.Convert_bool_To_Pointer_bool(&in.MakeIPTablesUtilChains, &out.MakeIPTablesUtilChains, s); err != nil {
		return err
	}
	if err := v1.Convert_int32_To_Pointer_int32(&in.IPTablesMasqueradeBit, &out.IPTablesMasqueradeBit, s); err != nil {
		return err
	}
	if err := v1.Convert_int32_To_Pointer_int32(&in.IPTablesDropBit, &out.IPTablesDropBit, s); err != nil {
		return err
	}
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.FailSwapOn = in.FailSwapOn
	out.SystemReserved = *(*map[string]string)(unsafe.Pointer(&in.SystemReserved))
	out.KubeReserved = *(*map[string]string)(unsafe.Pointer(&in.KubeReserved))
	out.SystemReservedCgroup = in.SystemReservedCgroup
	out.KubeReservedCgroup = in.KubeReservedCgroup
	out.EnforceNodeAllocatable = *(*[]string)(unsafe.Pointer(&in.EnforceNodeAllocatable))
	return nil
}

// Convert_kubeletconfig_KubeletConfiguration_To_v1alpha1_KubeletConfiguration is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletConfiguration_To_v1alpha1_KubeletConfiguration(in *kubeletconfig.KubeletConfiguration, out *KubeletConfiguration, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletConfiguration_To_v1alpha1_KubeletConfiguration(in, out, s)
}

func autoConvert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication(in *KubeletWebhookAuthentication, out *kubeletconfig.KubeletWebhookAuthentication, s conversion.Scope) error {
	if err := v1.Convert_Pointer_bool_To_bool(&in.Enabled, &out.Enabled, s); err != nil {
		return err
	}
	out.CacheTTL = in.CacheTTL
	return nil
}

// Convert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication is an autogenerated conversion function.
func Convert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication(in *KubeletWebhookAuthentication, out *kubeletconfig.KubeletWebhookAuthentication, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletWebhookAuthentication_To_kubeletconfig_KubeletWebhookAuthentication(in, out, s)
}

func autoConvert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication(in *kubeletconfig.KubeletWebhookAuthentication, out *KubeletWebhookAuthentication, s conversion.Scope) error {
	if err := v1.Convert_bool_To_Pointer_bool(&in.Enabled, &out.Enabled, s); err != nil {
		return err
	}
	out.CacheTTL = in.CacheTTL
	return nil
}

// Convert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication(in *kubeletconfig.KubeletWebhookAuthentication, out *KubeletWebhookAuthentication, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletWebhookAuthentication_To_v1alpha1_KubeletWebhookAuthentication(in, out, s)
}

func autoConvert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization(in *KubeletWebhookAuthorization, out *kubeletconfig.KubeletWebhookAuthorization, s conversion.Scope) error {
	out.CacheAuthorizedTTL = in.CacheAuthorizedTTL
	out.CacheUnauthorizedTTL = in.CacheUnauthorizedTTL
	return nil
}

// Convert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization is an autogenerated conversion function.
func Convert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization(in *KubeletWebhookAuthorization, out *kubeletconfig.KubeletWebhookAuthorization, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletWebhookAuthorization_To_kubeletconfig_KubeletWebhookAuthorization(in, out, s)
}

func autoConvert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization(in *kubeletconfig.KubeletWebhookAuthorization, out *KubeletWebhookAuthorization, s conversion.Scope) error {
	out.CacheAuthorizedTTL = in.CacheAuthorizedTTL
	out.CacheUnauthorizedTTL = in.CacheUnauthorizedTTL
	return nil
}

// Convert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization(in *kubeletconfig.KubeletWebhookAuthorization, out *KubeletWebhookAuthorization, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletWebhookAuthorization_To_v1alpha1_KubeletWebhookAuthorization(in, out, s)
}

func autoConvert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication(in *KubeletX509Authentication, out *kubeletconfig.KubeletX509Authentication, s conversion.Scope) error {
	out.ClientCAFile = in.ClientCAFile
	return nil
}

// Convert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication is an autogenerated conversion function.
func Convert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication(in *KubeletX509Authentication, out *kubeletconfig.KubeletX509Authentication, s conversion.Scope) error {
	return autoConvert_v1alpha1_KubeletX509Authentication_To_kubeletconfig_KubeletX509Authentication(in, out, s)
}

func autoConvert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication(in *kubeletconfig.KubeletX509Authentication, out *KubeletX509Authentication, s conversion.Scope) error {
	out.ClientCAFile = in.ClientCAFile
	return nil
}

// Convert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication is an autogenerated conversion function.
func Convert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication(in *kubeletconfig.KubeletX509Authentication, out *KubeletX509Authentication, s conversion.Scope) error {
	return autoConvert_kubeletconfig_KubeletX509Authentication_To_v1alpha1_KubeletX509Authentication(in, out, s)
}
