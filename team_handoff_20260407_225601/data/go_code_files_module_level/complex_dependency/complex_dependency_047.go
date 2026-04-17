func (in *KubeAPIServerConfig) DeepCopyInto(out *KubeAPIServerConfig) {
	*out = *in
	if in.EnableBootstrapAuthToken != nil {
		in, out := &in.EnableBootstrapAuthToken, &out.EnableBootstrapAuthToken
		*out = new(bool)
		**out = **in
	}
	if in.EnableAggregatorRouting != nil {
		in, out := &in.EnableAggregatorRouting, &out.EnableAggregatorRouting
		*out = new(bool)
		**out = **in
	}
	if in.AdmissionControl != nil {
		in, out := &in.AdmissionControl, &out.AdmissionControl
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EnableAdmissionPlugins != nil {
		in, out := &in.EnableAdmissionPlugins, &out.EnableAdmissionPlugins
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DisableAdmissionPlugins != nil {
		in, out := &in.DisableAdmissionPlugins, &out.DisableAdmissionPlugins
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EtcdServers != nil {
		in, out := &in.EtcdServers, &out.EtcdServers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.EtcdServersOverrides != nil {
		in, out := &in.EtcdServersOverrides, &out.EtcdServersOverrides
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.TLSCipherSuites != nil {
		in, out := &in.TLSCipherSuites, &out.TLSCipherSuites
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowPrivileged != nil {
		in, out := &in.AllowPrivileged, &out.AllowPrivileged
		*out = new(bool)
		**out = **in
	}
	if in.APIServerCount != nil {
		in, out := &in.APIServerCount, &out.APIServerCount
		*out = new(int32)
		**out = **in
	}
	if in.RuntimeConfig != nil {
		in, out := &in.RuntimeConfig, &out.RuntimeConfig
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.AnonymousAuth != nil {
		in, out := &in.AnonymousAuth, &out.AnonymousAuth
		*out = new(bool)
		**out = **in
	}
	if in.KubeletPreferredAddressTypes != nil {
		in, out := &in.KubeletPreferredAddressTypes, &out.KubeletPreferredAddressTypes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.StorageBackend != nil {
		in, out := &in.StorageBackend, &out.StorageBackend
		*out = new(string)
		**out = **in
	}
	if in.OIDCUsernameClaim != nil {
		in, out := &in.OIDCUsernameClaim, &out.OIDCUsernameClaim
		*out = new(string)
		**out = **in
	}
	if in.OIDCUsernamePrefix != nil {
		in, out := &in.OIDCUsernamePrefix, &out.OIDCUsernamePrefix
		*out = new(string)
		**out = **in
	}
	if in.OIDCGroupsClaim != nil {
		in, out := &in.OIDCGroupsClaim, &out.OIDCGroupsClaim
		*out = new(string)
		**out = **in
	}
	if in.OIDCGroupsPrefix != nil {
		in, out := &in.OIDCGroupsPrefix, &out.OIDCGroupsPrefix
		*out = new(string)
		**out = **in
	}
	if in.OIDCIssuerURL != nil {
		in, out := &in.OIDCIssuerURL, &out.OIDCIssuerURL
		*out = new(string)
		**out = **in
	}
	if in.OIDCClientID != nil {
		in, out := &in.OIDCClientID, &out.OIDCClientID
		*out = new(string)
		**out = **in
	}
	if in.OIDCRequiredClaim != nil {
		in, out := &in.OIDCRequiredClaim, &out.OIDCRequiredClaim
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.OIDCCAFile != nil {
		in, out := &in.OIDCCAFile, &out.OIDCCAFile
		*out = new(string)
		**out = **in
	}
	if in.ProxyClientCertFile != nil {
		in, out := &in.ProxyClientCertFile, &out.ProxyClientCertFile
		*out = new(string)
		**out = **in
	}
	if in.ProxyClientKeyFile != nil {
		in, out := &in.ProxyClientKeyFile, &out.ProxyClientKeyFile
		*out = new(string)
		**out = **in
	}
	if in.AuditLogFormat != nil {
		in, out := &in.AuditLogFormat, &out.AuditLogFormat
		*out = new(string)
		**out = **in
	}
	if in.AuditLogPath != nil {
		in, out := &in.AuditLogPath, &out.AuditLogPath
		*out = new(string)
		**out = **in
	}
	if in.AuditLogMaxAge != nil {
		in, out := &in.AuditLogMaxAge, &out.AuditLogMaxAge
		*out = new(int32)
		**out = **in
	}
	if in.AuditLogMaxBackups != nil {
		in, out := &in.AuditLogMaxBackups, &out.AuditLogMaxBackups
		*out = new(int32)
		**out = **in
	}
	if in.AuditLogMaxSize != nil {
		in, out := &in.AuditLogMaxSize, &out.AuditLogMaxSize
		*out = new(int32)
		**out = **in
	}
	if in.AuditWebhookBatchBufferSize != nil {
		in, out := &in.AuditWebhookBatchBufferSize, &out.AuditWebhookBatchBufferSize
		*out = new(int32)
		**out = **in
	}
	if in.AuditWebhookBatchMaxSize != nil {
		in, out := &in.AuditWebhookBatchMaxSize, &out.AuditWebhookBatchMaxSize
		*out = new(int32)
		**out = **in
	}
	if in.AuditWebhookBatchMaxWait != nil {
		in, out := &in.AuditWebhookBatchMaxWait, &out.AuditWebhookBatchMaxWait
		*out = new(v1.Duration)
		**out = **in
	}
	if in.AuditWebhookBatchThrottleBurst != nil {
		in, out := &in.AuditWebhookBatchThrottleBurst, &out.AuditWebhookBatchThrottleBurst
		*out = new(int32)
		**out = **in
	}
	if in.AuditWebhookBatchThrottleEnable != nil {
		in, out := &in.AuditWebhookBatchThrottleEnable, &out.AuditWebhookBatchThrottleEnable
		*out = new(bool)
		**out = **in
	}
	if in.AuditWebhookBatchThrottleQps != nil {
		in, out := &in.AuditWebhookBatchThrottleQps, &out.AuditWebhookBatchThrottleQps
		*out = new(float32)
		**out = **in
	}
	if in.AuditWebhookInitialBackoff != nil {
		in, out := &in.AuditWebhookInitialBackoff, &out.AuditWebhookInitialBackoff
		*out = new(v1.Duration)
		**out = **in
	}
	if in.AuthenticationTokenWebhookConfigFile != nil {
		in, out := &in.AuthenticationTokenWebhookConfigFile, &out.AuthenticationTokenWebhookConfigFile
		*out = new(string)
		**out = **in
	}
	if in.AuthenticationTokenWebhookCacheTTL != nil {
		in, out := &in.AuthenticationTokenWebhookCacheTTL, &out.AuthenticationTokenWebhookCacheTTL
		*out = new(v1.Duration)
		**out = **in
	}
	if in.AuthorizationMode != nil {
		in, out := &in.AuthorizationMode, &out.AuthorizationMode
		*out = new(string)
		**out = **in
	}
	if in.AuthorizationRBACSuperUser != nil {
		in, out := &in.AuthorizationRBACSuperUser, &out.AuthorizationRBACSuperUser
		*out = new(string)
		**out = **in
	}
	if in.ExperimentalEncryptionProviderConfig != nil {
		in, out := &in.ExperimentalEncryptionProviderConfig, &out.ExperimentalEncryptionProviderConfig
		*out = new(string)
		**out = **in
	}
	if in.RequestheaderUsernameHeaders != nil {
		in, out := &in.RequestheaderUsernameHeaders, &out.RequestheaderUsernameHeaders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequestheaderGroupHeaders != nil {
		in, out := &in.RequestheaderGroupHeaders, &out.RequestheaderGroupHeaders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequestheaderExtraHeaderPrefixes != nil {
		in, out := &in.RequestheaderExtraHeaderPrefixes, &out.RequestheaderExtraHeaderPrefixes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequestheaderAllowedNames != nil {
		in, out := &in.RequestheaderAllowedNames, &out.RequestheaderAllowedNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.HTTP2MaxStreamsPerConnection != nil {
		in, out := &in.HTTP2MaxStreamsPerConnection, &out.HTTP2MaxStreamsPerConnection
		*out = new(int32)
		**out = **in
	}
	if in.EtcdQuorumRead != nil {
		in, out := &in.EtcdQuorumRead, &out.EtcdQuorumRead
		*out = new(bool)
		**out = **in
	}
	if in.MinRequestTimeout != nil {
		in, out := &in.MinRequestTimeout, &out.MinRequestTimeout
		*out = new(int32)
		**out = **in
	}
	if in.ServiceAccountKeyFile != nil {
		in, out := &in.ServiceAccountKeyFile, &out.ServiceAccountKeyFile
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}
