func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*AccessSpec)(nil), (*kops.AccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AccessSpec_To_kops_AccessSpec(a.(*AccessSpec), b.(*kops.AccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AccessSpec)(nil), (*AccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AccessSpec_To_v1alpha2_AccessSpec(a.(*kops.AccessSpec), b.(*AccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AddonSpec)(nil), (*kops.AddonSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AddonSpec_To_kops_AddonSpec(a.(*AddonSpec), b.(*kops.AddonSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AddonSpec)(nil), (*AddonSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AddonSpec_To_v1alpha2_AddonSpec(a.(*kops.AddonSpec), b.(*AddonSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AlwaysAllowAuthorizationSpec)(nil), (*kops.AlwaysAllowAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AlwaysAllowAuthorizationSpec_To_kops_AlwaysAllowAuthorizationSpec(a.(*AlwaysAllowAuthorizationSpec), b.(*kops.AlwaysAllowAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AlwaysAllowAuthorizationSpec)(nil), (*AlwaysAllowAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AlwaysAllowAuthorizationSpec_To_v1alpha2_AlwaysAllowAuthorizationSpec(a.(*kops.AlwaysAllowAuthorizationSpec), b.(*AlwaysAllowAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AmazonVPCNetworkingSpec)(nil), (*kops.AmazonVPCNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AmazonVPCNetworkingSpec_To_kops_AmazonVPCNetworkingSpec(a.(*AmazonVPCNetworkingSpec), b.(*kops.AmazonVPCNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AmazonVPCNetworkingSpec)(nil), (*AmazonVPCNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AmazonVPCNetworkingSpec_To_v1alpha2_AmazonVPCNetworkingSpec(a.(*kops.AmazonVPCNetworkingSpec), b.(*AmazonVPCNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Assets)(nil), (*kops.Assets)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Assets_To_kops_Assets(a.(*Assets), b.(*kops.Assets), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.Assets)(nil), (*Assets)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_Assets_To_v1alpha2_Assets(a.(*kops.Assets), b.(*Assets), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AuthenticationSpec)(nil), (*kops.AuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AuthenticationSpec_To_kops_AuthenticationSpec(a.(*AuthenticationSpec), b.(*kops.AuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AuthenticationSpec)(nil), (*AuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AuthenticationSpec_To_v1alpha2_AuthenticationSpec(a.(*kops.AuthenticationSpec), b.(*AuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AuthorizationSpec)(nil), (*kops.AuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AuthorizationSpec_To_kops_AuthorizationSpec(a.(*AuthorizationSpec), b.(*kops.AuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AuthorizationSpec)(nil), (*AuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AuthorizationSpec_To_v1alpha2_AuthorizationSpec(a.(*kops.AuthorizationSpec), b.(*AuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*AwsAuthenticationSpec)(nil), (*kops.AwsAuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_AwsAuthenticationSpec_To_kops_AwsAuthenticationSpec(a.(*AwsAuthenticationSpec), b.(*kops.AwsAuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.AwsAuthenticationSpec)(nil), (*AwsAuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_AwsAuthenticationSpec_To_v1alpha2_AwsAuthenticationSpec(a.(*kops.AwsAuthenticationSpec), b.(*AwsAuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*BastionSpec)(nil), (*kops.BastionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_BastionSpec_To_kops_BastionSpec(a.(*BastionSpec), b.(*kops.BastionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.BastionSpec)(nil), (*BastionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_BastionSpec_To_v1alpha2_BastionSpec(a.(*kops.BastionSpec), b.(*BastionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CNINetworkingSpec)(nil), (*kops.CNINetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CNINetworkingSpec_To_kops_CNINetworkingSpec(a.(*CNINetworkingSpec), b.(*kops.CNINetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CNINetworkingSpec)(nil), (*CNINetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CNINetworkingSpec_To_v1alpha2_CNINetworkingSpec(a.(*kops.CNINetworkingSpec), b.(*CNINetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CalicoNetworkingSpec)(nil), (*kops.CalicoNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CalicoNetworkingSpec_To_kops_CalicoNetworkingSpec(a.(*CalicoNetworkingSpec), b.(*kops.CalicoNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CalicoNetworkingSpec)(nil), (*CalicoNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CalicoNetworkingSpec_To_v1alpha2_CalicoNetworkingSpec(a.(*kops.CalicoNetworkingSpec), b.(*CalicoNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CanalNetworkingSpec)(nil), (*kops.CanalNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CanalNetworkingSpec_To_kops_CanalNetworkingSpec(a.(*CanalNetworkingSpec), b.(*kops.CanalNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CanalNetworkingSpec)(nil), (*CanalNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CanalNetworkingSpec_To_v1alpha2_CanalNetworkingSpec(a.(*kops.CanalNetworkingSpec), b.(*CanalNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CiliumNetworkingSpec)(nil), (*kops.CiliumNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CiliumNetworkingSpec_To_kops_CiliumNetworkingSpec(a.(*CiliumNetworkingSpec), b.(*kops.CiliumNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CiliumNetworkingSpec)(nil), (*CiliumNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CiliumNetworkingSpec_To_v1alpha2_CiliumNetworkingSpec(a.(*kops.CiliumNetworkingSpec), b.(*CiliumNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClassicNetworkingSpec)(nil), (*kops.ClassicNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ClassicNetworkingSpec_To_kops_ClassicNetworkingSpec(a.(*ClassicNetworkingSpec), b.(*kops.ClassicNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ClassicNetworkingSpec)(nil), (*ClassicNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ClassicNetworkingSpec_To_v1alpha2_ClassicNetworkingSpec(a.(*kops.ClassicNetworkingSpec), b.(*ClassicNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CloudConfiguration)(nil), (*kops.CloudConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CloudConfiguration_To_kops_CloudConfiguration(a.(*CloudConfiguration), b.(*kops.CloudConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CloudConfiguration)(nil), (*CloudConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CloudConfiguration_To_v1alpha2_CloudConfiguration(a.(*kops.CloudConfiguration), b.(*CloudConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CloudControllerManagerConfig)(nil), (*kops.CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_CloudControllerManagerConfig_To_kops_CloudControllerManagerConfig(a.(*CloudControllerManagerConfig), b.(*kops.CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.CloudControllerManagerConfig)(nil), (*CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_CloudControllerManagerConfig_To_v1alpha2_CloudControllerManagerConfig(a.(*kops.CloudControllerManagerConfig), b.(*CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Cluster)(nil), (*kops.Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Cluster_To_kops_Cluster(a.(*Cluster), b.(*kops.Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.Cluster)(nil), (*Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_Cluster_To_v1alpha2_Cluster(a.(*kops.Cluster), b.(*Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterList)(nil), (*kops.ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ClusterList_To_kops_ClusterList(a.(*ClusterList), b.(*kops.ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ClusterList)(nil), (*ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ClusterList_To_v1alpha2_ClusterList(a.(*kops.ClusterList), b.(*ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterSpec)(nil), (*kops.ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ClusterSpec_To_kops_ClusterSpec(a.(*ClusterSpec), b.(*kops.ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ClusterSpec)(nil), (*ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ClusterSpec_To_v1alpha2_ClusterSpec(a.(*kops.ClusterSpec), b.(*ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterSubnetSpec)(nil), (*kops.ClusterSubnetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ClusterSubnetSpec_To_kops_ClusterSubnetSpec(a.(*ClusterSubnetSpec), b.(*kops.ClusterSubnetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ClusterSubnetSpec)(nil), (*ClusterSubnetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ClusterSubnetSpec_To_v1alpha2_ClusterSubnetSpec(a.(*kops.ClusterSubnetSpec), b.(*ClusterSubnetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DNSAccessSpec)(nil), (*kops.DNSAccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_DNSAccessSpec_To_kops_DNSAccessSpec(a.(*DNSAccessSpec), b.(*kops.DNSAccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.DNSAccessSpec)(nil), (*DNSAccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_DNSAccessSpec_To_v1alpha2_DNSAccessSpec(a.(*kops.DNSAccessSpec), b.(*DNSAccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DNSSpec)(nil), (*kops.DNSSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_DNSSpec_To_kops_DNSSpec(a.(*DNSSpec), b.(*kops.DNSSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.DNSSpec)(nil), (*DNSSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_DNSSpec_To_v1alpha2_DNSSpec(a.(*kops.DNSSpec), b.(*DNSSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*DockerConfig)(nil), (*kops.DockerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_DockerConfig_To_kops_DockerConfig(a.(*DockerConfig), b.(*kops.DockerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.DockerConfig)(nil), (*DockerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_DockerConfig_To_v1alpha2_DockerConfig(a.(*kops.DockerConfig), b.(*DockerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EgressProxySpec)(nil), (*kops.EgressProxySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_EgressProxySpec_To_kops_EgressProxySpec(a.(*EgressProxySpec), b.(*kops.EgressProxySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.EgressProxySpec)(nil), (*EgressProxySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_EgressProxySpec_To_v1alpha2_EgressProxySpec(a.(*kops.EgressProxySpec), b.(*EgressProxySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EtcdBackupSpec)(nil), (*kops.EtcdBackupSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_EtcdBackupSpec_To_kops_EtcdBackupSpec(a.(*EtcdBackupSpec), b.(*kops.EtcdBackupSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.EtcdBackupSpec)(nil), (*EtcdBackupSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_EtcdBackupSpec_To_v1alpha2_EtcdBackupSpec(a.(*kops.EtcdBackupSpec), b.(*EtcdBackupSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EtcdClusterSpec)(nil), (*kops.EtcdClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_EtcdClusterSpec_To_kops_EtcdClusterSpec(a.(*EtcdClusterSpec), b.(*kops.EtcdClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.EtcdClusterSpec)(nil), (*EtcdClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_EtcdClusterSpec_To_v1alpha2_EtcdClusterSpec(a.(*kops.EtcdClusterSpec), b.(*EtcdClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EtcdManagerSpec)(nil), (*kops.EtcdManagerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_EtcdManagerSpec_To_kops_EtcdManagerSpec(a.(*EtcdManagerSpec), b.(*kops.EtcdManagerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.EtcdManagerSpec)(nil), (*EtcdManagerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_EtcdManagerSpec_To_v1alpha2_EtcdManagerSpec(a.(*kops.EtcdManagerSpec), b.(*EtcdManagerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EtcdMemberSpec)(nil), (*kops.EtcdMemberSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_EtcdMemberSpec_To_kops_EtcdMemberSpec(a.(*EtcdMemberSpec), b.(*kops.EtcdMemberSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.EtcdMemberSpec)(nil), (*EtcdMemberSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_EtcdMemberSpec_To_v1alpha2_EtcdMemberSpec(a.(*kops.EtcdMemberSpec), b.(*EtcdMemberSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExecContainerAction)(nil), (*kops.ExecContainerAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ExecContainerAction_To_kops_ExecContainerAction(a.(*ExecContainerAction), b.(*kops.ExecContainerAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ExecContainerAction)(nil), (*ExecContainerAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ExecContainerAction_To_v1alpha2_ExecContainerAction(a.(*kops.ExecContainerAction), b.(*ExecContainerAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExternalDNSConfig)(nil), (*kops.ExternalDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ExternalDNSConfig_To_kops_ExternalDNSConfig(a.(*ExternalDNSConfig), b.(*kops.ExternalDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ExternalDNSConfig)(nil), (*ExternalDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ExternalDNSConfig_To_v1alpha2_ExternalDNSConfig(a.(*kops.ExternalDNSConfig), b.(*ExternalDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExternalNetworkingSpec)(nil), (*kops.ExternalNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_ExternalNetworkingSpec_To_kops_ExternalNetworkingSpec(a.(*ExternalNetworkingSpec), b.(*kops.ExternalNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.ExternalNetworkingSpec)(nil), (*ExternalNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_ExternalNetworkingSpec_To_v1alpha2_ExternalNetworkingSpec(a.(*kops.ExternalNetworkingSpec), b.(*ExternalNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FileAssetSpec)(nil), (*kops.FileAssetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_FileAssetSpec_To_kops_FileAssetSpec(a.(*FileAssetSpec), b.(*kops.FileAssetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.FileAssetSpec)(nil), (*FileAssetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_FileAssetSpec_To_v1alpha2_FileAssetSpec(a.(*kops.FileAssetSpec), b.(*FileAssetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FlannelNetworkingSpec)(nil), (*kops.FlannelNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_FlannelNetworkingSpec_To_kops_FlannelNetworkingSpec(a.(*FlannelNetworkingSpec), b.(*kops.FlannelNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.FlannelNetworkingSpec)(nil), (*FlannelNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_FlannelNetworkingSpec_To_v1alpha2_FlannelNetworkingSpec(a.(*kops.FlannelNetworkingSpec), b.(*FlannelNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HTTPProxy)(nil), (*kops.HTTPProxy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_HTTPProxy_To_kops_HTTPProxy(a.(*HTTPProxy), b.(*kops.HTTPProxy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.HTTPProxy)(nil), (*HTTPProxy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_HTTPProxy_To_v1alpha2_HTTPProxy(a.(*kops.HTTPProxy), b.(*HTTPProxy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*HookSpec)(nil), (*kops.HookSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_HookSpec_To_kops_HookSpec(a.(*HookSpec), b.(*kops.HookSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.HookSpec)(nil), (*HookSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_HookSpec_To_v1alpha2_HookSpec(a.(*kops.HookSpec), b.(*HookSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*IAMProfileSpec)(nil), (*kops.IAMProfileSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_IAMProfileSpec_To_kops_IAMProfileSpec(a.(*IAMProfileSpec), b.(*kops.IAMProfileSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.IAMProfileSpec)(nil), (*IAMProfileSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_IAMProfileSpec_To_v1alpha2_IAMProfileSpec(a.(*kops.IAMProfileSpec), b.(*IAMProfileSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*IAMSpec)(nil), (*kops.IAMSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_IAMSpec_To_kops_IAMSpec(a.(*IAMSpec), b.(*kops.IAMSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.IAMSpec)(nil), (*IAMSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_IAMSpec_To_v1alpha2_IAMSpec(a.(*kops.IAMSpec), b.(*IAMSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceGroup)(nil), (*kops.InstanceGroup)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_InstanceGroup_To_kops_InstanceGroup(a.(*InstanceGroup), b.(*kops.InstanceGroup), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.InstanceGroup)(nil), (*InstanceGroup)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_InstanceGroup_To_v1alpha2_InstanceGroup(a.(*kops.InstanceGroup), b.(*InstanceGroup), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceGroupList)(nil), (*kops.InstanceGroupList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_InstanceGroupList_To_kops_InstanceGroupList(a.(*InstanceGroupList), b.(*kops.InstanceGroupList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.InstanceGroupList)(nil), (*InstanceGroupList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_InstanceGroupList_To_v1alpha2_InstanceGroupList(a.(*kops.InstanceGroupList), b.(*InstanceGroupList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstanceGroupSpec)(nil), (*kops.InstanceGroupSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_InstanceGroupSpec_To_kops_InstanceGroupSpec(a.(*InstanceGroupSpec), b.(*kops.InstanceGroupSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.InstanceGroupSpec)(nil), (*InstanceGroupSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_InstanceGroupSpec_To_v1alpha2_InstanceGroupSpec(a.(*kops.InstanceGroupSpec), b.(*InstanceGroupSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Keyset)(nil), (*kops.Keyset)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_Keyset_To_kops_Keyset(a.(*Keyset), b.(*kops.Keyset), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.Keyset)(nil), (*Keyset)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_Keyset_To_v1alpha2_Keyset(a.(*kops.Keyset), b.(*Keyset), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KeysetItem)(nil), (*kops.KeysetItem)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KeysetItem_To_kops_KeysetItem(a.(*KeysetItem), b.(*kops.KeysetItem), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KeysetItem)(nil), (*KeysetItem)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KeysetItem_To_v1alpha2_KeysetItem(a.(*kops.KeysetItem), b.(*KeysetItem), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KeysetList)(nil), (*kops.KeysetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KeysetList_To_kops_KeysetList(a.(*KeysetList), b.(*kops.KeysetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KeysetList)(nil), (*KeysetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KeysetList_To_v1alpha2_KeysetList(a.(*kops.KeysetList), b.(*KeysetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KeysetSpec)(nil), (*kops.KeysetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KeysetSpec_To_kops_KeysetSpec(a.(*KeysetSpec), b.(*kops.KeysetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KeysetSpec)(nil), (*KeysetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KeysetSpec_To_v1alpha2_KeysetSpec(a.(*kops.KeysetSpec), b.(*KeysetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KopeioAuthenticationSpec)(nil), (*kops.KopeioAuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KopeioAuthenticationSpec_To_kops_KopeioAuthenticationSpec(a.(*KopeioAuthenticationSpec), b.(*kops.KopeioAuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KopeioAuthenticationSpec)(nil), (*KopeioAuthenticationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KopeioAuthenticationSpec_To_v1alpha2_KopeioAuthenticationSpec(a.(*kops.KopeioAuthenticationSpec), b.(*KopeioAuthenticationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KopeioNetworkingSpec)(nil), (*kops.KopeioNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KopeioNetworkingSpec_To_kops_KopeioNetworkingSpec(a.(*KopeioNetworkingSpec), b.(*kops.KopeioNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KopeioNetworkingSpec)(nil), (*KopeioNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KopeioNetworkingSpec_To_v1alpha2_KopeioNetworkingSpec(a.(*kops.KopeioNetworkingSpec), b.(*KopeioNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeAPIServerConfig)(nil), (*kops.KubeAPIServerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeAPIServerConfig_To_kops_KubeAPIServerConfig(a.(*KubeAPIServerConfig), b.(*kops.KubeAPIServerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeAPIServerConfig)(nil), (*KubeAPIServerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeAPIServerConfig_To_v1alpha2_KubeAPIServerConfig(a.(*kops.KubeAPIServerConfig), b.(*KubeAPIServerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeControllerManagerConfig)(nil), (*kops.KubeControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeControllerManagerConfig_To_kops_KubeControllerManagerConfig(a.(*KubeControllerManagerConfig), b.(*kops.KubeControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeControllerManagerConfig)(nil), (*KubeControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeControllerManagerConfig_To_v1alpha2_KubeControllerManagerConfig(a.(*kops.KubeControllerManagerConfig), b.(*KubeControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeDNSConfig)(nil), (*kops.KubeDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeDNSConfig_To_kops_KubeDNSConfig(a.(*KubeDNSConfig), b.(*kops.KubeDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeDNSConfig)(nil), (*KubeDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeDNSConfig_To_v1alpha2_KubeDNSConfig(a.(*kops.KubeDNSConfig), b.(*KubeDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeProxyConfig)(nil), (*kops.KubeProxyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeProxyConfig_To_kops_KubeProxyConfig(a.(*KubeProxyConfig), b.(*kops.KubeProxyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeProxyConfig)(nil), (*KubeProxyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeProxyConfig_To_v1alpha2_KubeProxyConfig(a.(*kops.KubeProxyConfig), b.(*KubeProxyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeSchedulerConfig)(nil), (*kops.KubeSchedulerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeSchedulerConfig_To_kops_KubeSchedulerConfig(a.(*KubeSchedulerConfig), b.(*kops.KubeSchedulerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeSchedulerConfig)(nil), (*KubeSchedulerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeSchedulerConfig_To_v1alpha2_KubeSchedulerConfig(a.(*kops.KubeSchedulerConfig), b.(*KubeSchedulerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubeletConfigSpec)(nil), (*kops.KubeletConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubeletConfigSpec_To_kops_KubeletConfigSpec(a.(*KubeletConfigSpec), b.(*kops.KubeletConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubeletConfigSpec)(nil), (*KubeletConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubeletConfigSpec_To_v1alpha2_KubeletConfigSpec(a.(*kops.KubeletConfigSpec), b.(*KubeletConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KubenetNetworkingSpec)(nil), (*kops.KubenetNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KubenetNetworkingSpec_To_kops_KubenetNetworkingSpec(a.(*KubenetNetworkingSpec), b.(*kops.KubenetNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KubenetNetworkingSpec)(nil), (*KubenetNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KubenetNetworkingSpec_To_v1alpha2_KubenetNetworkingSpec(a.(*kops.KubenetNetworkingSpec), b.(*KubenetNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*KuberouterNetworkingSpec)(nil), (*kops.KuberouterNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_KuberouterNetworkingSpec_To_kops_KuberouterNetworkingSpec(a.(*KuberouterNetworkingSpec), b.(*kops.KuberouterNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.KuberouterNetworkingSpec)(nil), (*KuberouterNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_KuberouterNetworkingSpec_To_v1alpha2_KuberouterNetworkingSpec(a.(*kops.KuberouterNetworkingSpec), b.(*KuberouterNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LeaderElectionConfiguration)(nil), (*kops.LeaderElectionConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_LeaderElectionConfiguration_To_kops_LeaderElectionConfiguration(a.(*LeaderElectionConfiguration), b.(*kops.LeaderElectionConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.LeaderElectionConfiguration)(nil), (*LeaderElectionConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_LeaderElectionConfiguration_To_v1alpha2_LeaderElectionConfiguration(a.(*kops.LeaderElectionConfiguration), b.(*LeaderElectionConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LoadBalancer)(nil), (*kops.LoadBalancer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_LoadBalancer_To_kops_LoadBalancer(a.(*LoadBalancer), b.(*kops.LoadBalancer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.LoadBalancer)(nil), (*LoadBalancer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_LoadBalancer_To_v1alpha2_LoadBalancer(a.(*kops.LoadBalancer), b.(*LoadBalancer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LoadBalancerAccessSpec)(nil), (*kops.LoadBalancerAccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_LoadBalancerAccessSpec_To_kops_LoadBalancerAccessSpec(a.(*LoadBalancerAccessSpec), b.(*kops.LoadBalancerAccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.LoadBalancerAccessSpec)(nil), (*LoadBalancerAccessSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_LoadBalancerAccessSpec_To_v1alpha2_LoadBalancerAccessSpec(a.(*kops.LoadBalancerAccessSpec), b.(*LoadBalancerAccessSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LyftVPCNetworkingSpec)(nil), (*kops.LyftVPCNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_LyftVPCNetworkingSpec_To_kops_LyftVPCNetworkingSpec(a.(*LyftVPCNetworkingSpec), b.(*kops.LyftVPCNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.LyftVPCNetworkingSpec)(nil), (*LyftVPCNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_LyftVPCNetworkingSpec_To_v1alpha2_LyftVPCNetworkingSpec(a.(*kops.LyftVPCNetworkingSpec), b.(*LyftVPCNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MixedInstancesPolicySpec)(nil), (*kops.MixedInstancesPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_MixedInstancesPolicySpec_To_kops_MixedInstancesPolicySpec(a.(*MixedInstancesPolicySpec), b.(*kops.MixedInstancesPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.MixedInstancesPolicySpec)(nil), (*MixedInstancesPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_MixedInstancesPolicySpec_To_v1alpha2_MixedInstancesPolicySpec(a.(*kops.MixedInstancesPolicySpec), b.(*MixedInstancesPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NetworkingSpec)(nil), (*kops.NetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_NetworkingSpec_To_kops_NetworkingSpec(a.(*NetworkingSpec), b.(*kops.NetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.NetworkingSpec)(nil), (*NetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_NetworkingSpec_To_v1alpha2_NetworkingSpec(a.(*kops.NetworkingSpec), b.(*NetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NodeAuthorizationSpec)(nil), (*kops.NodeAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_NodeAuthorizationSpec_To_kops_NodeAuthorizationSpec(a.(*NodeAuthorizationSpec), b.(*kops.NodeAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.NodeAuthorizationSpec)(nil), (*NodeAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_NodeAuthorizationSpec_To_v1alpha2_NodeAuthorizationSpec(a.(*kops.NodeAuthorizationSpec), b.(*NodeAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NodeAuthorizerSpec)(nil), (*kops.NodeAuthorizerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_NodeAuthorizerSpec_To_kops_NodeAuthorizerSpec(a.(*NodeAuthorizerSpec), b.(*kops.NodeAuthorizerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.NodeAuthorizerSpec)(nil), (*NodeAuthorizerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_NodeAuthorizerSpec_To_v1alpha2_NodeAuthorizerSpec(a.(*kops.NodeAuthorizerSpec), b.(*NodeAuthorizerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OpenstackBlockStorageConfig)(nil), (*kops.OpenstackBlockStorageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_OpenstackBlockStorageConfig_To_kops_OpenstackBlockStorageConfig(a.(*OpenstackBlockStorageConfig), b.(*kops.OpenstackBlockStorageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.OpenstackBlockStorageConfig)(nil), (*OpenstackBlockStorageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_OpenstackBlockStorageConfig_To_v1alpha2_OpenstackBlockStorageConfig(a.(*kops.OpenstackBlockStorageConfig), b.(*OpenstackBlockStorageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OpenstackConfiguration)(nil), (*kops.OpenstackConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_OpenstackConfiguration_To_kops_OpenstackConfiguration(a.(*OpenstackConfiguration), b.(*kops.OpenstackConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.OpenstackConfiguration)(nil), (*OpenstackConfiguration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_OpenstackConfiguration_To_v1alpha2_OpenstackConfiguration(a.(*kops.OpenstackConfiguration), b.(*OpenstackConfiguration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OpenstackLoadbalancerConfig)(nil), (*kops.OpenstackLoadbalancerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_OpenstackLoadbalancerConfig_To_kops_OpenstackLoadbalancerConfig(a.(*OpenstackLoadbalancerConfig), b.(*kops.OpenstackLoadbalancerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.OpenstackLoadbalancerConfig)(nil), (*OpenstackLoadbalancerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_OpenstackLoadbalancerConfig_To_v1alpha2_OpenstackLoadbalancerConfig(a.(*kops.OpenstackLoadbalancerConfig), b.(*OpenstackLoadbalancerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OpenstackMonitor)(nil), (*kops.OpenstackMonitor)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_OpenstackMonitor_To_kops_OpenstackMonitor(a.(*OpenstackMonitor), b.(*kops.OpenstackMonitor), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.OpenstackMonitor)(nil), (*OpenstackMonitor)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_OpenstackMonitor_To_v1alpha2_OpenstackMonitor(a.(*kops.OpenstackMonitor), b.(*OpenstackMonitor), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OpenstackRouter)(nil), (*kops.OpenstackRouter)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_OpenstackRouter_To_kops_OpenstackRouter(a.(*OpenstackRouter), b.(*kops.OpenstackRouter), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.OpenstackRouter)(nil), (*OpenstackRouter)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_OpenstackRouter_To_v1alpha2_OpenstackRouter(a.(*kops.OpenstackRouter), b.(*OpenstackRouter), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RBACAuthorizationSpec)(nil), (*kops.RBACAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_RBACAuthorizationSpec_To_kops_RBACAuthorizationSpec(a.(*RBACAuthorizationSpec), b.(*kops.RBACAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.RBACAuthorizationSpec)(nil), (*RBACAuthorizationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_RBACAuthorizationSpec_To_v1alpha2_RBACAuthorizationSpec(a.(*kops.RBACAuthorizationSpec), b.(*RBACAuthorizationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RomanaNetworkingSpec)(nil), (*kops.RomanaNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_RomanaNetworkingSpec_To_kops_RomanaNetworkingSpec(a.(*RomanaNetworkingSpec), b.(*kops.RomanaNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.RomanaNetworkingSpec)(nil), (*RomanaNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_RomanaNetworkingSpec_To_v1alpha2_RomanaNetworkingSpec(a.(*kops.RomanaNetworkingSpec), b.(*RomanaNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SSHCredential)(nil), (*kops.SSHCredential)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_SSHCredential_To_kops_SSHCredential(a.(*SSHCredential), b.(*kops.SSHCredential), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.SSHCredential)(nil), (*SSHCredential)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_SSHCredential_To_v1alpha2_SSHCredential(a.(*kops.SSHCredential), b.(*SSHCredential), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SSHCredentialList)(nil), (*kops.SSHCredentialList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_SSHCredentialList_To_kops_SSHCredentialList(a.(*SSHCredentialList), b.(*kops.SSHCredentialList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.SSHCredentialList)(nil), (*SSHCredentialList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_SSHCredentialList_To_v1alpha2_SSHCredentialList(a.(*kops.SSHCredentialList), b.(*SSHCredentialList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SSHCredentialSpec)(nil), (*kops.SSHCredentialSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_SSHCredentialSpec_To_kops_SSHCredentialSpec(a.(*SSHCredentialSpec), b.(*kops.SSHCredentialSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.SSHCredentialSpec)(nil), (*SSHCredentialSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_SSHCredentialSpec_To_v1alpha2_SSHCredentialSpec(a.(*kops.SSHCredentialSpec), b.(*SSHCredentialSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TargetSpec)(nil), (*kops.TargetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_TargetSpec_To_kops_TargetSpec(a.(*TargetSpec), b.(*kops.TargetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.TargetSpec)(nil), (*TargetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_TargetSpec_To_v1alpha2_TargetSpec(a.(*kops.TargetSpec), b.(*TargetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TerraformSpec)(nil), (*kops.TerraformSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_TerraformSpec_To_kops_TerraformSpec(a.(*TerraformSpec), b.(*kops.TerraformSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.TerraformSpec)(nil), (*TerraformSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_TerraformSpec_To_v1alpha2_TerraformSpec(a.(*kops.TerraformSpec), b.(*TerraformSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*TopologySpec)(nil), (*kops.TopologySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_TopologySpec_To_kops_TopologySpec(a.(*TopologySpec), b.(*kops.TopologySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.TopologySpec)(nil), (*TopologySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_TopologySpec_To_v1alpha2_TopologySpec(a.(*kops.TopologySpec), b.(*TopologySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*UserData)(nil), (*kops.UserData)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_UserData_To_kops_UserData(a.(*UserData), b.(*kops.UserData), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.UserData)(nil), (*UserData)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_UserData_To_v1alpha2_UserData(a.(*kops.UserData), b.(*UserData), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VolumeMountSpec)(nil), (*kops.VolumeMountSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_VolumeMountSpec_To_kops_VolumeMountSpec(a.(*VolumeMountSpec), b.(*kops.VolumeMountSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.VolumeMountSpec)(nil), (*VolumeMountSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_VolumeMountSpec_To_v1alpha2_VolumeMountSpec(a.(*kops.VolumeMountSpec), b.(*VolumeMountSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VolumeSpec)(nil), (*kops.VolumeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_VolumeSpec_To_kops_VolumeSpec(a.(*VolumeSpec), b.(*kops.VolumeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.VolumeSpec)(nil), (*VolumeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_VolumeSpec_To_v1alpha2_VolumeSpec(a.(*kops.VolumeSpec), b.(*VolumeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*WeaveNetworkingSpec)(nil), (*kops.WeaveNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha2_WeaveNetworkingSpec_To_kops_WeaveNetworkingSpec(a.(*WeaveNetworkingSpec), b.(*kops.WeaveNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*kops.WeaveNetworkingSpec)(nil), (*WeaveNetworkingSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_kops_WeaveNetworkingSpec_To_v1alpha2_WeaveNetworkingSpec(a.(*kops.WeaveNetworkingSpec), b.(*WeaveNetworkingSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}
