func (m *Cluster) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) < 1 {
		return ClusterValidationError{
			field:  "Name",
			reason: "value length must be at least 1 bytes",
		}
	}

	// no validation rules for AltStatName

	{
		tmp := m.GetEdsClusterConfig()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "EdsClusterConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	if true {
		dur := m.GetConnectTimeout()

		gt := time.Duration(0*time.Second + 0*time.Nanosecond)

		if dur <= gt {
			return ClusterValidationError{
				field:  "ConnectTimeout",
				reason: "value must be greater than 0s",
			}
		}

	}

	{
		tmp := m.GetPerConnectionBufferLimitBytes()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "PerConnectionBufferLimitBytes",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	if _, ok := Cluster_LbPolicy_name[int32(m.GetLbPolicy())]; !ok {
		return ClusterValidationError{
			field:  "LbPolicy",
			reason: "value must be one of the defined enum values",
		}
	}

	for idx, item := range m.GetHosts() {
		_, _ = idx, item

		{
			tmp := item

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  fmt.Sprintf("Hosts[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	}

	{
		tmp := m.GetLoadAssignment()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "LoadAssignment",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	for idx, item := range m.GetHealthChecks() {
		_, _ = idx, item

		{
			tmp := item

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  fmt.Sprintf("HealthChecks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	}

	{
		tmp := m.GetMaxRequestsPerConnection()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "MaxRequestsPerConnection",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetCircuitBreakers()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "CircuitBreakers",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetTlsContext()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "TlsContext",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetCommonHttpProtocolOptions()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "CommonHttpProtocolOptions",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetHttpProtocolOptions()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "HttpProtocolOptions",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetHttp2ProtocolOptions()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "Http2ProtocolOptions",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	// no validation rules for ExtensionProtocolOptions

	// no validation rules for TypedExtensionProtocolOptions

	if d := m.GetDnsRefreshRate(); d != nil {
		dur := *d

		gt := time.Duration(0*time.Second + 0*time.Nanosecond)

		if dur <= gt {
			return ClusterValidationError{
				field:  "DnsRefreshRate",
				reason: "value must be greater than 0s",
			}
		}

	}

	if _, ok := Cluster_DnsLookupFamily_name[int32(m.GetDnsLookupFamily())]; !ok {
		return ClusterValidationError{
			field:  "DnsLookupFamily",
			reason: "value must be one of the defined enum values",
		}
	}

	for idx, item := range m.GetDnsResolvers() {
		_, _ = idx, item

		{
			tmp := item

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  fmt.Sprintf("DnsResolvers[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	}

	{
		tmp := m.GetOutlierDetection()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "OutlierDetection",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	if d := m.GetCleanupInterval(); d != nil {
		dur := *d

		gt := time.Duration(0*time.Second + 0*time.Nanosecond)

		if dur <= gt {
			return ClusterValidationError{
				field:  "CleanupInterval",
				reason: "value must be greater than 0s",
			}
		}

	}

	{
		tmp := m.GetUpstreamBindConfig()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "UpstreamBindConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetLbSubsetConfig()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "LbSubsetConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetCommonLbConfig()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "CommonLbConfig",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetTransportSocket()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "TransportSocket",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	{
		tmp := m.GetMetadata()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	// no validation rules for ProtocolSelection

	{
		tmp := m.GetUpstreamConnectionOptions()

		if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

			if err := v.Validate(); err != nil {
				return ClusterValidationError{
					field:  "UpstreamConnectionOptions",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}
	}

	// no validation rules for CloseConnectionsOnHostHealthFailure

	// no validation rules for DrainConnectionsOnHostRemoval

	switch m.ClusterDiscoveryType.(type) {

	case *Cluster_Type:

		if _, ok := Cluster_DiscoveryType_name[int32(m.GetType())]; !ok {
			return ClusterValidationError{
				field:  "Type",
				reason: "value must be one of the defined enum values",
			}
		}

	case *Cluster_ClusterType:

		{
			tmp := m.GetClusterType()

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  "ClusterType",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	}

	switch m.LbConfig.(type) {

	case *Cluster_RingHashLbConfig_:

		{
			tmp := m.GetRingHashLbConfig()

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  "RingHashLbConfig",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	case *Cluster_OriginalDstLbConfig_:

		{
			tmp := m.GetOriginalDstLbConfig()

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  "OriginalDstLbConfig",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	case *Cluster_LeastRequestLbConfig_:

		{
			tmp := m.GetLeastRequestLbConfig()

			if v, ok := interface{}(tmp).(interface{ Validate() error }); ok {

				if err := v.Validate(); err != nil {
					return ClusterValidationError{
						field:  "LeastRequestLbConfig",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}
		}

	}

	return nil
}
