func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.ActiveDirectoryConfig)(nil), (*config.ActiveDirectoryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ActiveDirectoryConfig_To_config_ActiveDirectoryConfig(a.(*v1.ActiveDirectoryConfig), b.(*config.ActiveDirectoryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ActiveDirectoryConfig)(nil), (*v1.ActiveDirectoryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ActiveDirectoryConfig_To_v1_ActiveDirectoryConfig(a.(*config.ActiveDirectoryConfig), b.(*v1.ActiveDirectoryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AdmissionConfig)(nil), (*config.AdmissionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AdmissionConfig_To_config_AdmissionConfig(a.(*v1.AdmissionConfig), b.(*config.AdmissionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AdmissionConfig)(nil), (*v1.AdmissionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AdmissionConfig_To_v1_AdmissionConfig(a.(*config.AdmissionConfig), b.(*v1.AdmissionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AdmissionPluginConfig)(nil), (*config.AdmissionPluginConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AdmissionPluginConfig_To_config_AdmissionPluginConfig(a.(*v1.AdmissionPluginConfig), b.(*config.AdmissionPluginConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AdmissionPluginConfig)(nil), (*v1.AdmissionPluginConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AdmissionPluginConfig_To_v1_AdmissionPluginConfig(a.(*config.AdmissionPluginConfig), b.(*v1.AdmissionPluginConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AggregatorConfig)(nil), (*config.AggregatorConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AggregatorConfig_To_config_AggregatorConfig(a.(*v1.AggregatorConfig), b.(*config.AggregatorConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AggregatorConfig)(nil), (*v1.AggregatorConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AggregatorConfig_To_v1_AggregatorConfig(a.(*config.AggregatorConfig), b.(*v1.AggregatorConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AllowAllPasswordIdentityProvider)(nil), (*config.AllowAllPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AllowAllPasswordIdentityProvider_To_config_AllowAllPasswordIdentityProvider(a.(*v1.AllowAllPasswordIdentityProvider), b.(*config.AllowAllPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AllowAllPasswordIdentityProvider)(nil), (*v1.AllowAllPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AllowAllPasswordIdentityProvider_To_v1_AllowAllPasswordIdentityProvider(a.(*config.AllowAllPasswordIdentityProvider), b.(*v1.AllowAllPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AuditConfig)(nil), (*config.AuditConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AuditConfig_To_config_AuditConfig(a.(*v1.AuditConfig), b.(*config.AuditConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AuditConfig)(nil), (*v1.AuditConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AuditConfig_To_v1_AuditConfig(a.(*config.AuditConfig), b.(*v1.AuditConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AugmentedActiveDirectoryConfig)(nil), (*config.AugmentedActiveDirectoryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AugmentedActiveDirectoryConfig_To_config_AugmentedActiveDirectoryConfig(a.(*v1.AugmentedActiveDirectoryConfig), b.(*config.AugmentedActiveDirectoryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.AugmentedActiveDirectoryConfig)(nil), (*v1.AugmentedActiveDirectoryConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AugmentedActiveDirectoryConfig_To_v1_AugmentedActiveDirectoryConfig(a.(*config.AugmentedActiveDirectoryConfig), b.(*v1.AugmentedActiveDirectoryConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BasicAuthPasswordIdentityProvider)(nil), (*config.BasicAuthPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BasicAuthPasswordIdentityProvider_To_config_BasicAuthPasswordIdentityProvider(a.(*v1.BasicAuthPasswordIdentityProvider), b.(*config.BasicAuthPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.BasicAuthPasswordIdentityProvider)(nil), (*v1.BasicAuthPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_BasicAuthPasswordIdentityProvider_To_v1_BasicAuthPasswordIdentityProvider(a.(*config.BasicAuthPasswordIdentityProvider), b.(*v1.BasicAuthPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildDefaultsConfig)(nil), (*config.BuildDefaultsConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildDefaultsConfig_To_config_BuildDefaultsConfig(a.(*v1.BuildDefaultsConfig), b.(*config.BuildDefaultsConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.BuildDefaultsConfig)(nil), (*v1.BuildDefaultsConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_BuildDefaultsConfig_To_v1_BuildDefaultsConfig(a.(*config.BuildDefaultsConfig), b.(*v1.BuildDefaultsConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildOverridesConfig)(nil), (*config.BuildOverridesConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildOverridesConfig_To_config_BuildOverridesConfig(a.(*v1.BuildOverridesConfig), b.(*config.BuildOverridesConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.BuildOverridesConfig)(nil), (*v1.BuildOverridesConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_BuildOverridesConfig_To_v1_BuildOverridesConfig(a.(*config.BuildOverridesConfig), b.(*v1.BuildOverridesConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CertInfo)(nil), (*config.CertInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CertInfo_To_config_CertInfo(a.(*v1.CertInfo), b.(*config.CertInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.CertInfo)(nil), (*v1.CertInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_CertInfo_To_v1_CertInfo(a.(*config.CertInfo), b.(*v1.CertInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClientConnectionOverrides)(nil), (*config.ClientConnectionOverrides)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClientConnectionOverrides_To_config_ClientConnectionOverrides(a.(*v1.ClientConnectionOverrides), b.(*config.ClientConnectionOverrides), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ClientConnectionOverrides)(nil), (*v1.ClientConnectionOverrides)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ClientConnectionOverrides_To_v1_ClientConnectionOverrides(a.(*config.ClientConnectionOverrides), b.(*v1.ClientConnectionOverrides), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterNetworkEntry)(nil), (*config.ClusterNetworkEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterNetworkEntry_To_config_ClusterNetworkEntry(a.(*v1.ClusterNetworkEntry), b.(*config.ClusterNetworkEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ClusterNetworkEntry)(nil), (*v1.ClusterNetworkEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ClusterNetworkEntry_To_v1_ClusterNetworkEntry(a.(*config.ClusterNetworkEntry), b.(*v1.ClusterNetworkEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ControllerConfig)(nil), (*config.ControllerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ControllerConfig_To_config_ControllerConfig(a.(*v1.ControllerConfig), b.(*config.ControllerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ControllerConfig)(nil), (*v1.ControllerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ControllerConfig_To_v1_ControllerConfig(a.(*config.ControllerConfig), b.(*v1.ControllerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ControllerElectionConfig)(nil), (*config.ControllerElectionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ControllerElectionConfig_To_config_ControllerElectionConfig(a.(*v1.ControllerElectionConfig), b.(*config.ControllerElectionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ControllerElectionConfig)(nil), (*v1.ControllerElectionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ControllerElectionConfig_To_v1_ControllerElectionConfig(a.(*config.ControllerElectionConfig), b.(*v1.ControllerElectionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DNSConfig)(nil), (*config.DNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DNSConfig_To_config_DNSConfig(a.(*v1.DNSConfig), b.(*config.DNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.DNSConfig)(nil), (*v1.DNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_DNSConfig_To_v1_DNSConfig(a.(*config.DNSConfig), b.(*v1.DNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DefaultAdmissionConfig)(nil), (*config.DefaultAdmissionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DefaultAdmissionConfig_To_config_DefaultAdmissionConfig(a.(*v1.DefaultAdmissionConfig), b.(*config.DefaultAdmissionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.DefaultAdmissionConfig)(nil), (*v1.DefaultAdmissionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_DefaultAdmissionConfig_To_v1_DefaultAdmissionConfig(a.(*config.DefaultAdmissionConfig), b.(*v1.DefaultAdmissionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DenyAllPasswordIdentityProvider)(nil), (*config.DenyAllPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DenyAllPasswordIdentityProvider_To_config_DenyAllPasswordIdentityProvider(a.(*v1.DenyAllPasswordIdentityProvider), b.(*config.DenyAllPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.DenyAllPasswordIdentityProvider)(nil), (*v1.DenyAllPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_DenyAllPasswordIdentityProvider_To_v1_DenyAllPasswordIdentityProvider(a.(*config.DenyAllPasswordIdentityProvider), b.(*v1.DenyAllPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DockerConfig)(nil), (*config.DockerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DockerConfig_To_config_DockerConfig(a.(*v1.DockerConfig), b.(*config.DockerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.DockerConfig)(nil), (*v1.DockerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_DockerConfig_To_v1_DockerConfig(a.(*config.DockerConfig), b.(*v1.DockerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EtcdConfig)(nil), (*config.EtcdConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EtcdConfig_To_config_EtcdConfig(a.(*v1.EtcdConfig), b.(*config.EtcdConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.EtcdConfig)(nil), (*v1.EtcdConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_EtcdConfig_To_v1_EtcdConfig(a.(*config.EtcdConfig), b.(*v1.EtcdConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EtcdConnectionInfo)(nil), (*config.EtcdConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EtcdConnectionInfo_To_config_EtcdConnectionInfo(a.(*v1.EtcdConnectionInfo), b.(*config.EtcdConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.EtcdConnectionInfo)(nil), (*v1.EtcdConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_EtcdConnectionInfo_To_v1_EtcdConnectionInfo(a.(*config.EtcdConnectionInfo), b.(*v1.EtcdConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EtcdStorageConfig)(nil), (*config.EtcdStorageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EtcdStorageConfig_To_config_EtcdStorageConfig(a.(*v1.EtcdStorageConfig), b.(*config.EtcdStorageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.EtcdStorageConfig)(nil), (*v1.EtcdStorageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_EtcdStorageConfig_To_v1_EtcdStorageConfig(a.(*config.EtcdStorageConfig), b.(*v1.EtcdStorageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitHubIdentityProvider)(nil), (*config.GitHubIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitHubIdentityProvider_To_config_GitHubIdentityProvider(a.(*v1.GitHubIdentityProvider), b.(*config.GitHubIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GitHubIdentityProvider)(nil), (*v1.GitHubIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GitHubIdentityProvider_To_v1_GitHubIdentityProvider(a.(*config.GitHubIdentityProvider), b.(*v1.GitHubIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitLabIdentityProvider)(nil), (*config.GitLabIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitLabIdentityProvider_To_config_GitLabIdentityProvider(a.(*v1.GitLabIdentityProvider), b.(*config.GitLabIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GitLabIdentityProvider)(nil), (*v1.GitLabIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GitLabIdentityProvider_To_v1_GitLabIdentityProvider(a.(*config.GitLabIdentityProvider), b.(*v1.GitLabIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GoogleIdentityProvider)(nil), (*config.GoogleIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GoogleIdentityProvider_To_config_GoogleIdentityProvider(a.(*v1.GoogleIdentityProvider), b.(*config.GoogleIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GoogleIdentityProvider)(nil), (*v1.GoogleIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GoogleIdentityProvider_To_v1_GoogleIdentityProvider(a.(*config.GoogleIdentityProvider), b.(*v1.GoogleIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GrantConfig)(nil), (*config.GrantConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GrantConfig_To_config_GrantConfig(a.(*v1.GrantConfig), b.(*config.GrantConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GrantConfig)(nil), (*v1.GrantConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GrantConfig_To_v1_GrantConfig(a.(*config.GrantConfig), b.(*v1.GrantConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GroupResource)(nil), (*config.GroupResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GroupResource_To_config_GroupResource(a.(*v1.GroupResource), b.(*config.GroupResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.GroupResource)(nil), (*v1.GroupResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_GroupResource_To_v1_GroupResource(a.(*config.GroupResource), b.(*v1.GroupResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HTPasswdPasswordIdentityProvider)(nil), (*config.HTPasswdPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HTPasswdPasswordIdentityProvider_To_config_HTPasswdPasswordIdentityProvider(a.(*v1.HTPasswdPasswordIdentityProvider), b.(*config.HTPasswdPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.HTPasswdPasswordIdentityProvider)(nil), (*v1.HTPasswdPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_HTPasswdPasswordIdentityProvider_To_v1_HTPasswdPasswordIdentityProvider(a.(*config.HTPasswdPasswordIdentityProvider), b.(*v1.HTPasswdPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HTTPServingInfo)(nil), (*config.HTTPServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HTTPServingInfo_To_config_HTTPServingInfo(a.(*v1.HTTPServingInfo), b.(*config.HTTPServingInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.HTTPServingInfo)(nil), (*v1.HTTPServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_HTTPServingInfo_To_v1_HTTPServingInfo(a.(*config.HTTPServingInfo), b.(*v1.HTTPServingInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.IdentityProvider)(nil), (*config.IdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_IdentityProvider_To_config_IdentityProvider(a.(*v1.IdentityProvider), b.(*config.IdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.IdentityProvider)(nil), (*v1.IdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_IdentityProvider_To_v1_IdentityProvider(a.(*config.IdentityProvider), b.(*v1.IdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageConfig)(nil), (*config.ImageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageConfig_To_config_ImageConfig(a.(*v1.ImageConfig), b.(*config.ImageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ImageConfig)(nil), (*v1.ImageConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ImageConfig_To_v1_ImageConfig(a.(*config.ImageConfig), b.(*v1.ImageConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImagePolicyConfig)(nil), (*config.ImagePolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImagePolicyConfig_To_config_ImagePolicyConfig(a.(*v1.ImagePolicyConfig), b.(*config.ImagePolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ImagePolicyConfig)(nil), (*v1.ImagePolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ImagePolicyConfig_To_v1_ImagePolicyConfig(a.(*config.ImagePolicyConfig), b.(*v1.ImagePolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.JenkinsPipelineConfig)(nil), (*config.JenkinsPipelineConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_JenkinsPipelineConfig_To_config_JenkinsPipelineConfig(a.(*v1.JenkinsPipelineConfig), b.(*config.JenkinsPipelineConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.JenkinsPipelineConfig)(nil), (*v1.JenkinsPipelineConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_JenkinsPipelineConfig_To_v1_JenkinsPipelineConfig(a.(*config.JenkinsPipelineConfig), b.(*v1.JenkinsPipelineConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.KeystonePasswordIdentityProvider)(nil), (*config.KeystonePasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KeystonePasswordIdentityProvider_To_config_KeystonePasswordIdentityProvider(a.(*v1.KeystonePasswordIdentityProvider), b.(*config.KeystonePasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.KeystonePasswordIdentityProvider)(nil), (*v1.KeystonePasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KeystonePasswordIdentityProvider_To_v1_KeystonePasswordIdentityProvider(a.(*config.KeystonePasswordIdentityProvider), b.(*v1.KeystonePasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.KubeletConnectionInfo)(nil), (*config.KubeletConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KubeletConnectionInfo_To_config_KubeletConnectionInfo(a.(*v1.KubeletConnectionInfo), b.(*config.KubeletConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.KubeletConnectionInfo)(nil), (*v1.KubeletConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KubeletConnectionInfo_To_v1_KubeletConnectionInfo(a.(*config.KubeletConnectionInfo), b.(*v1.KubeletConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.KubernetesMasterConfig)(nil), (*config.KubernetesMasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KubernetesMasterConfig_To_config_KubernetesMasterConfig(a.(*v1.KubernetesMasterConfig), b.(*config.KubernetesMasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.KubernetesMasterConfig)(nil), (*v1.KubernetesMasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KubernetesMasterConfig_To_v1_KubernetesMasterConfig(a.(*config.KubernetesMasterConfig), b.(*v1.KubernetesMasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LDAPAttributeMapping)(nil), (*config.LDAPAttributeMapping)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LDAPAttributeMapping_To_config_LDAPAttributeMapping(a.(*v1.LDAPAttributeMapping), b.(*config.LDAPAttributeMapping), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LDAPAttributeMapping)(nil), (*v1.LDAPAttributeMapping)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LDAPAttributeMapping_To_v1_LDAPAttributeMapping(a.(*config.LDAPAttributeMapping), b.(*v1.LDAPAttributeMapping), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LDAPPasswordIdentityProvider)(nil), (*config.LDAPPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LDAPPasswordIdentityProvider_To_config_LDAPPasswordIdentityProvider(a.(*v1.LDAPPasswordIdentityProvider), b.(*config.LDAPPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LDAPPasswordIdentityProvider)(nil), (*v1.LDAPPasswordIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LDAPPasswordIdentityProvider_To_v1_LDAPPasswordIdentityProvider(a.(*config.LDAPPasswordIdentityProvider), b.(*v1.LDAPPasswordIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LDAPQuery)(nil), (*config.LDAPQuery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LDAPQuery_To_config_LDAPQuery(a.(*v1.LDAPQuery), b.(*config.LDAPQuery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LDAPQuery)(nil), (*v1.LDAPQuery)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LDAPQuery_To_v1_LDAPQuery(a.(*config.LDAPQuery), b.(*v1.LDAPQuery), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LDAPSyncConfig)(nil), (*config.LDAPSyncConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LDAPSyncConfig_To_config_LDAPSyncConfig(a.(*v1.LDAPSyncConfig), b.(*config.LDAPSyncConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LDAPSyncConfig)(nil), (*v1.LDAPSyncConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LDAPSyncConfig_To_v1_LDAPSyncConfig(a.(*config.LDAPSyncConfig), b.(*v1.LDAPSyncConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LocalQuota)(nil), (*config.LocalQuota)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalQuota_To_config_LocalQuota(a.(*v1.LocalQuota), b.(*config.LocalQuota), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.LocalQuota)(nil), (*v1.LocalQuota)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_LocalQuota_To_v1_LocalQuota(a.(*config.LocalQuota), b.(*v1.LocalQuota), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MasterAuthConfig)(nil), (*config.MasterAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterAuthConfig_To_config_MasterAuthConfig(a.(*v1.MasterAuthConfig), b.(*config.MasterAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MasterAuthConfig)(nil), (*v1.MasterAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterAuthConfig_To_v1_MasterAuthConfig(a.(*config.MasterAuthConfig), b.(*v1.MasterAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MasterClients)(nil), (*config.MasterClients)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterClients_To_config_MasterClients(a.(*v1.MasterClients), b.(*config.MasterClients), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MasterClients)(nil), (*v1.MasterClients)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterClients_To_v1_MasterClients(a.(*config.MasterClients), b.(*v1.MasterClients), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MasterConfig)(nil), (*config.MasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterConfig_To_config_MasterConfig(a.(*v1.MasterConfig), b.(*config.MasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MasterConfig)(nil), (*v1.MasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterConfig_To_v1_MasterConfig(a.(*config.MasterConfig), b.(*v1.MasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MasterNetworkConfig)(nil), (*config.MasterNetworkConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterNetworkConfig_To_config_MasterNetworkConfig(a.(*v1.MasterNetworkConfig), b.(*config.MasterNetworkConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MasterNetworkConfig)(nil), (*v1.MasterNetworkConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterNetworkConfig_To_v1_MasterNetworkConfig(a.(*config.MasterNetworkConfig), b.(*v1.MasterNetworkConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MasterVolumeConfig)(nil), (*config.MasterVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterVolumeConfig_To_config_MasterVolumeConfig(a.(*v1.MasterVolumeConfig), b.(*config.MasterVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.MasterVolumeConfig)(nil), (*v1.MasterVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterVolumeConfig_To_v1_MasterVolumeConfig(a.(*config.MasterVolumeConfig), b.(*v1.MasterVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NamedCertificate)(nil), (*config.NamedCertificate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NamedCertificate_To_config_NamedCertificate(a.(*v1.NamedCertificate), b.(*config.NamedCertificate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.NamedCertificate)(nil), (*v1.NamedCertificate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NamedCertificate_To_v1_NamedCertificate(a.(*config.NamedCertificate), b.(*v1.NamedCertificate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeAuthConfig)(nil), (*config.NodeAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeAuthConfig_To_config_NodeAuthConfig(a.(*v1.NodeAuthConfig), b.(*config.NodeAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.NodeAuthConfig)(nil), (*v1.NodeAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NodeAuthConfig_To_v1_NodeAuthConfig(a.(*config.NodeAuthConfig), b.(*v1.NodeAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeConfig)(nil), (*config.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeConfig_To_config_NodeConfig(a.(*v1.NodeConfig), b.(*config.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.NodeConfig)(nil), (*v1.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NodeConfig_To_v1_NodeConfig(a.(*config.NodeConfig), b.(*v1.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeNetworkConfig)(nil), (*config.NodeNetworkConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeNetworkConfig_To_config_NodeNetworkConfig(a.(*v1.NodeNetworkConfig), b.(*config.NodeNetworkConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.NodeNetworkConfig)(nil), (*v1.NodeNetworkConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NodeNetworkConfig_To_v1_NodeNetworkConfig(a.(*config.NodeNetworkConfig), b.(*v1.NodeNetworkConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeVolumeConfig)(nil), (*config.NodeVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeVolumeConfig_To_config_NodeVolumeConfig(a.(*v1.NodeVolumeConfig), b.(*config.NodeVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.NodeVolumeConfig)(nil), (*v1.NodeVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NodeVolumeConfig_To_v1_NodeVolumeConfig(a.(*config.NodeVolumeConfig), b.(*v1.NodeVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthConfig)(nil), (*config.OAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthConfig_To_config_OAuthConfig(a.(*v1.OAuthConfig), b.(*config.OAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.OAuthConfig)(nil), (*v1.OAuthConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_OAuthConfig_To_v1_OAuthConfig(a.(*config.OAuthConfig), b.(*v1.OAuthConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthTemplates)(nil), (*config.OAuthTemplates)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthTemplates_To_config_OAuthTemplates(a.(*v1.OAuthTemplates), b.(*config.OAuthTemplates), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.OAuthTemplates)(nil), (*v1.OAuthTemplates)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_OAuthTemplates_To_v1_OAuthTemplates(a.(*config.OAuthTemplates), b.(*v1.OAuthTemplates), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OpenIDClaims)(nil), (*config.OpenIDClaims)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OpenIDClaims_To_config_OpenIDClaims(a.(*v1.OpenIDClaims), b.(*config.OpenIDClaims), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.OpenIDClaims)(nil), (*v1.OpenIDClaims)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_OpenIDClaims_To_v1_OpenIDClaims(a.(*config.OpenIDClaims), b.(*v1.OpenIDClaims), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OpenIDIdentityProvider)(nil), (*config.OpenIDIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OpenIDIdentityProvider_To_config_OpenIDIdentityProvider(a.(*v1.OpenIDIdentityProvider), b.(*config.OpenIDIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.OpenIDIdentityProvider)(nil), (*v1.OpenIDIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_OpenIDIdentityProvider_To_v1_OpenIDIdentityProvider(a.(*config.OpenIDIdentityProvider), b.(*v1.OpenIDIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OpenIDURLs)(nil), (*config.OpenIDURLs)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OpenIDURLs_To_config_OpenIDURLs(a.(*v1.OpenIDURLs), b.(*config.OpenIDURLs), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.OpenIDURLs)(nil), (*v1.OpenIDURLs)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_OpenIDURLs_To_v1_OpenIDURLs(a.(*config.OpenIDURLs), b.(*v1.OpenIDURLs), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodManifestConfig)(nil), (*config.PodManifestConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodManifestConfig_To_config_PodManifestConfig(a.(*v1.PodManifestConfig), b.(*config.PodManifestConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.PodManifestConfig)(nil), (*v1.PodManifestConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_PodManifestConfig_To_v1_PodManifestConfig(a.(*config.PodManifestConfig), b.(*v1.PodManifestConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PolicyConfig)(nil), (*config.PolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PolicyConfig_To_config_PolicyConfig(a.(*v1.PolicyConfig), b.(*config.PolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.PolicyConfig)(nil), (*v1.PolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_PolicyConfig_To_v1_PolicyConfig(a.(*config.PolicyConfig), b.(*v1.PolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ProjectConfig)(nil), (*config.ProjectConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ProjectConfig_To_config_ProjectConfig(a.(*v1.ProjectConfig), b.(*config.ProjectConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ProjectConfig)(nil), (*v1.ProjectConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ProjectConfig_To_v1_ProjectConfig(a.(*config.ProjectConfig), b.(*v1.ProjectConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RFC2307Config)(nil), (*config.RFC2307Config)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RFC2307Config_To_config_RFC2307Config(a.(*v1.RFC2307Config), b.(*config.RFC2307Config), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RFC2307Config)(nil), (*v1.RFC2307Config)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RFC2307Config_To_v1_RFC2307Config(a.(*config.RFC2307Config), b.(*v1.RFC2307Config), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RegistryLocation)(nil), (*config.RegistryLocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RegistryLocation_To_config_RegistryLocation(a.(*v1.RegistryLocation), b.(*config.RegistryLocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RegistryLocation)(nil), (*v1.RegistryLocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RegistryLocation_To_v1_RegistryLocation(a.(*config.RegistryLocation), b.(*v1.RegistryLocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RemoteConnectionInfo)(nil), (*config.RemoteConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RemoteConnectionInfo_To_config_RemoteConnectionInfo(a.(*v1.RemoteConnectionInfo), b.(*config.RemoteConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RemoteConnectionInfo)(nil), (*v1.RemoteConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RemoteConnectionInfo_To_v1_RemoteConnectionInfo(a.(*config.RemoteConnectionInfo), b.(*v1.RemoteConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RequestHeaderAuthenticationOptions)(nil), (*config.RequestHeaderAuthenticationOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RequestHeaderAuthenticationOptions_To_config_RequestHeaderAuthenticationOptions(a.(*v1.RequestHeaderAuthenticationOptions), b.(*config.RequestHeaderAuthenticationOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RequestHeaderAuthenticationOptions)(nil), (*v1.RequestHeaderAuthenticationOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RequestHeaderAuthenticationOptions_To_v1_RequestHeaderAuthenticationOptions(a.(*config.RequestHeaderAuthenticationOptions), b.(*v1.RequestHeaderAuthenticationOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RequestHeaderIdentityProvider)(nil), (*config.RequestHeaderIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RequestHeaderIdentityProvider_To_config_RequestHeaderIdentityProvider(a.(*v1.RequestHeaderIdentityProvider), b.(*config.RequestHeaderIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RequestHeaderIdentityProvider)(nil), (*v1.RequestHeaderIdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RequestHeaderIdentityProvider_To_v1_RequestHeaderIdentityProvider(a.(*config.RequestHeaderIdentityProvider), b.(*v1.RequestHeaderIdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoutingConfig)(nil), (*config.RoutingConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoutingConfig_To_config_RoutingConfig(a.(*v1.RoutingConfig), b.(*config.RoutingConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.RoutingConfig)(nil), (*v1.RoutingConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RoutingConfig_To_v1_RoutingConfig(a.(*config.RoutingConfig), b.(*v1.RoutingConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecurityAllocator)(nil), (*config.SecurityAllocator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecurityAllocator_To_config_SecurityAllocator(a.(*v1.SecurityAllocator), b.(*config.SecurityAllocator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SecurityAllocator)(nil), (*v1.SecurityAllocator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SecurityAllocator_To_v1_SecurityAllocator(a.(*config.SecurityAllocator), b.(*v1.SecurityAllocator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountConfig)(nil), (*config.ServiceAccountConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountConfig_To_config_ServiceAccountConfig(a.(*v1.ServiceAccountConfig), b.(*config.ServiceAccountConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ServiceAccountConfig)(nil), (*v1.ServiceAccountConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ServiceAccountConfig_To_v1_ServiceAccountConfig(a.(*config.ServiceAccountConfig), b.(*v1.ServiceAccountConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceServingCert)(nil), (*config.ServiceServingCert)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceServingCert_To_config_ServiceServingCert(a.(*v1.ServiceServingCert), b.(*config.ServiceServingCert), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ServiceServingCert)(nil), (*v1.ServiceServingCert)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ServiceServingCert_To_v1_ServiceServingCert(a.(*config.ServiceServingCert), b.(*v1.ServiceServingCert), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServingInfo)(nil), (*config.ServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServingInfo_To_config_ServingInfo(a.(*v1.ServingInfo), b.(*config.ServingInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.ServingInfo)(nil), (*v1.ServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ServingInfo_To_v1_ServingInfo(a.(*config.ServingInfo), b.(*v1.ServingInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SessionConfig)(nil), (*config.SessionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SessionConfig_To_config_SessionConfig(a.(*v1.SessionConfig), b.(*config.SessionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SessionConfig)(nil), (*v1.SessionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SessionConfig_To_v1_SessionConfig(a.(*config.SessionConfig), b.(*v1.SessionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SessionSecret)(nil), (*config.SessionSecret)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SessionSecret_To_config_SessionSecret(a.(*v1.SessionSecret), b.(*config.SessionSecret), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SessionSecret)(nil), (*v1.SessionSecret)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SessionSecret_To_v1_SessionSecret(a.(*config.SessionSecret), b.(*v1.SessionSecret), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SessionSecrets)(nil), (*config.SessionSecrets)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SessionSecrets_To_config_SessionSecrets(a.(*v1.SessionSecrets), b.(*config.SessionSecrets), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SessionSecrets)(nil), (*v1.SessionSecrets)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SessionSecrets_To_v1_SessionSecrets(a.(*config.SessionSecrets), b.(*v1.SessionSecrets), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SourceStrategyDefaultsConfig)(nil), (*config.SourceStrategyDefaultsConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceStrategyDefaultsConfig_To_config_SourceStrategyDefaultsConfig(a.(*v1.SourceStrategyDefaultsConfig), b.(*config.SourceStrategyDefaultsConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.SourceStrategyDefaultsConfig)(nil), (*v1.SourceStrategyDefaultsConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_SourceStrategyDefaultsConfig_To_v1_SourceStrategyDefaultsConfig(a.(*config.SourceStrategyDefaultsConfig), b.(*v1.SourceStrategyDefaultsConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StringSource)(nil), (*config.StringSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StringSource_To_config_StringSource(a.(*v1.StringSource), b.(*config.StringSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.StringSource)(nil), (*v1.StringSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_StringSource_To_v1_StringSource(a.(*config.StringSource), b.(*v1.StringSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StringSourceSpec)(nil), (*config.StringSourceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StringSourceSpec_To_config_StringSourceSpec(a.(*v1.StringSourceSpec), b.(*config.StringSourceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.StringSourceSpec)(nil), (*v1.StringSourceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_StringSourceSpec_To_v1_StringSourceSpec(a.(*config.StringSourceSpec), b.(*v1.StringSourceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TokenConfig)(nil), (*config.TokenConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TokenConfig_To_config_TokenConfig(a.(*v1.TokenConfig), b.(*config.TokenConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.TokenConfig)(nil), (*v1.TokenConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_TokenConfig_To_v1_TokenConfig(a.(*config.TokenConfig), b.(*v1.TokenConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.UserAgentDenyRule)(nil), (*config.UserAgentDenyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_UserAgentDenyRule_To_config_UserAgentDenyRule(a.(*v1.UserAgentDenyRule), b.(*config.UserAgentDenyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.UserAgentDenyRule)(nil), (*v1.UserAgentDenyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_UserAgentDenyRule_To_v1_UserAgentDenyRule(a.(*config.UserAgentDenyRule), b.(*v1.UserAgentDenyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.UserAgentMatchRule)(nil), (*config.UserAgentMatchRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_UserAgentMatchRule_To_config_UserAgentMatchRule(a.(*v1.UserAgentMatchRule), b.(*config.UserAgentMatchRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.UserAgentMatchRule)(nil), (*v1.UserAgentMatchRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_UserAgentMatchRule_To_v1_UserAgentMatchRule(a.(*config.UserAgentMatchRule), b.(*v1.UserAgentMatchRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.UserAgentMatchingConfig)(nil), (*config.UserAgentMatchingConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_UserAgentMatchingConfig_To_config_UserAgentMatchingConfig(a.(*v1.UserAgentMatchingConfig), b.(*config.UserAgentMatchingConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.UserAgentMatchingConfig)(nil), (*v1.UserAgentMatchingConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_UserAgentMatchingConfig_To_v1_UserAgentMatchingConfig(a.(*config.UserAgentMatchingConfig), b.(*v1.UserAgentMatchingConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.WebhookTokenAuthenticator)(nil), (*config.WebhookTokenAuthenticator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_WebhookTokenAuthenticator_To_config_WebhookTokenAuthenticator(a.(*v1.WebhookTokenAuthenticator), b.(*config.WebhookTokenAuthenticator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*config.WebhookTokenAuthenticator)(nil), (*v1.WebhookTokenAuthenticator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_WebhookTokenAuthenticator_To_v1_WebhookTokenAuthenticator(a.(*config.WebhookTokenAuthenticator), b.(*v1.WebhookTokenAuthenticator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.AdmissionPluginConfig)(nil), (*v1.AdmissionPluginConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AdmissionPluginConfig_To_v1_AdmissionPluginConfig(a.(*config.AdmissionPluginConfig), b.(*v1.AdmissionPluginConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.AuditConfig)(nil), (*v1.AuditConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_AuditConfig_To_v1_AuditConfig(a.(*config.AuditConfig), b.(*v1.AuditConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.EtcdConnectionInfo)(nil), (*v1.EtcdConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_EtcdConnectionInfo_To_v1_EtcdConnectionInfo(a.(*config.EtcdConnectionInfo), b.(*v1.EtcdConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.IdentityProvider)(nil), (*v1.IdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_IdentityProvider_To_v1_IdentityProvider(a.(*config.IdentityProvider), b.(*v1.IdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.ImagePolicyConfig)(nil), (*v1.ImagePolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ImagePolicyConfig_To_v1_ImagePolicyConfig(a.(*config.ImagePolicyConfig), b.(*v1.ImagePolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.KubeletConnectionInfo)(nil), (*v1.KubeletConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KubeletConnectionInfo_To_v1_KubeletConnectionInfo(a.(*config.KubeletConnectionInfo), b.(*v1.KubeletConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.KubernetesMasterConfig)(nil), (*v1.KubernetesMasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_KubernetesMasterConfig_To_v1_KubernetesMasterConfig(a.(*config.KubernetesMasterConfig), b.(*v1.KubernetesMasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.MasterVolumeConfig)(nil), (*v1.MasterVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_MasterVolumeConfig_To_v1_MasterVolumeConfig(a.(*config.MasterVolumeConfig), b.(*v1.MasterVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.NodeConfig)(nil), (*v1.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_NodeConfig_To_v1_NodeConfig(a.(*config.NodeConfig), b.(*v1.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.RemoteConnectionInfo)(nil), (*v1.RemoteConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_RemoteConnectionInfo_To_v1_RemoteConnectionInfo(a.(*config.RemoteConnectionInfo), b.(*v1.RemoteConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*config.ServingInfo)(nil), (*v1.ServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_config_ServingInfo_To_v1_ServingInfo(a.(*config.ServingInfo), b.(*v1.ServingInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.AdmissionPluginConfig)(nil), (*config.AdmissionPluginConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AdmissionPluginConfig_To_config_AdmissionPluginConfig(a.(*v1.AdmissionPluginConfig), b.(*config.AdmissionPluginConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.AuditConfig)(nil), (*config.AuditConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AuditConfig_To_config_AuditConfig(a.(*v1.AuditConfig), b.(*config.AuditConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.EtcdConnectionInfo)(nil), (*config.EtcdConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EtcdConnectionInfo_To_config_EtcdConnectionInfo(a.(*v1.EtcdConnectionInfo), b.(*config.EtcdConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.IdentityProvider)(nil), (*config.IdentityProvider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_IdentityProvider_To_config_IdentityProvider(a.(*v1.IdentityProvider), b.(*config.IdentityProvider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ImagePolicyConfig)(nil), (*config.ImagePolicyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImagePolicyConfig_To_config_ImagePolicyConfig(a.(*v1.ImagePolicyConfig), b.(*config.ImagePolicyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.KubeletConnectionInfo)(nil), (*config.KubeletConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KubeletConnectionInfo_To_config_KubeletConnectionInfo(a.(*v1.KubeletConnectionInfo), b.(*config.KubeletConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.KubernetesMasterConfig)(nil), (*config.KubernetesMasterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KubernetesMasterConfig_To_config_KubernetesMasterConfig(a.(*v1.KubernetesMasterConfig), b.(*config.KubernetesMasterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.MasterNetworkConfig)(nil), (*config.MasterNetworkConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterNetworkConfig_To_config_MasterNetworkConfig(a.(*v1.MasterNetworkConfig), b.(*config.MasterNetworkConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.MasterVolumeConfig)(nil), (*config.MasterVolumeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MasterVolumeConfig_To_config_MasterVolumeConfig(a.(*v1.MasterVolumeConfig), b.(*config.MasterVolumeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.NodeConfig)(nil), (*config.NodeConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeConfig_To_config_NodeConfig(a.(*v1.NodeConfig), b.(*config.NodeConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.RemoteConnectionInfo)(nil), (*config.RemoteConnectionInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RemoteConnectionInfo_To_config_RemoteConnectionInfo(a.(*v1.RemoteConnectionInfo), b.(*config.RemoteConnectionInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ServingInfo)(nil), (*config.ServingInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServingInfo_To_config_ServingInfo(a.(*v1.ServingInfo), b.(*config.ServingInfo), scope)
	}); err != nil {
		return err
	}
	return nil
}
