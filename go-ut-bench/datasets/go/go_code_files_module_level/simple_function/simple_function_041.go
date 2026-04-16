func (b *Builder) Validate(rt RuntimeConfig) error {
	// reDatacenter defines a regexp for a valid datacenter name
	var reDatacenter = regexp.MustCompile("^[a-z0-9_-]+$")

	// ----------------------------------------------------------------
	// check required params we cannot recover from first
	//

	if rt.Datacenter == "" {
		return fmt.Errorf("datacenter cannot be empty")
	}
	if !reDatacenter.MatchString(rt.Datacenter) {
		return fmt.Errorf("datacenter cannot be %q. Please use only [a-z0-9-_].", rt.Datacenter)
	}
	if rt.DataDir == "" && !rt.DevMode {
		return fmt.Errorf("data_dir cannot be empty")
	}
	if !rt.DevMode {
		fi, err := os.Stat(rt.DataDir)
		switch {
		case err != nil && !os.IsNotExist(err):
			return fmt.Errorf("Error getting info on data_dir: %s", err)
		case err == nil && !fi.IsDir():
			return fmt.Errorf("data_dir %q is not a directory", rt.DataDir)
		}
	}
	if rt.NodeName == "" {
		return fmt.Errorf("node_name cannot be empty")
	}
	if ipaddr.IsAny(rt.AdvertiseAddrLAN.IP) {
		return fmt.Errorf("Advertise address cannot be 0.0.0.0, :: or [::]")
	}
	if ipaddr.IsAny(rt.AdvertiseAddrWAN.IP) {
		return fmt.Errorf("Advertise WAN address cannot be 0.0.0.0, :: or [::]")
	}
	if err := b.validateSegments(rt); err != nil {
		return err
	}
	for _, a := range rt.DNSAddrs {
		if _, ok := a.(*net.UnixAddr); ok {
			return fmt.Errorf("DNS address cannot be a unix socket")
		}
	}
	for _, a := range rt.DNSRecursors {
		if ipaddr.IsAny(a) {
			return fmt.Errorf("DNS recursor address cannot be 0.0.0.0, :: or [::]")
		}
	}
	if rt.Bootstrap && !rt.ServerMode {
		return fmt.Errorf("'bootstrap = true' requires 'server = true'")
	}
	if rt.BootstrapExpect < 0 {
		return fmt.Errorf("bootstrap_expect cannot be %d. Must be greater than or equal to zero", rt.BootstrapExpect)
	}
	if rt.BootstrapExpect > 0 && !rt.ServerMode {
		return fmt.Errorf("'bootstrap_expect > 0' requires 'server = true'")
	}
	if rt.BootstrapExpect > 0 && rt.DevMode {
		return fmt.Errorf("'bootstrap_expect > 0' not allowed in dev mode")
	}
	if rt.BootstrapExpect > 0 && rt.Bootstrap {
		return fmt.Errorf("'bootstrap_expect > 0' and 'bootstrap = true' are mutually exclusive")
	}
	if rt.AEInterval <= 0 {
		return fmt.Errorf("ae_interval cannot be %s. Must be positive", rt.AEInterval)
	}
	if rt.AutopilotMaxTrailingLogs < 0 {
		return fmt.Errorf("autopilot.max_trailing_logs cannot be %d. Must be greater than or equal to zero", rt.AutopilotMaxTrailingLogs)
	}
	if rt.ACLDatacenter != "" && !reDatacenter.MatchString(rt.ACLDatacenter) {
		return fmt.Errorf("acl_datacenter cannot be %q. Please use only [a-z0-9-_].", rt.ACLDatacenter)
	}
	if rt.EnableUI && rt.UIDir != "" {
		return fmt.Errorf(
			"Both the ui and ui-dir flags were specified, please provide only one.\n" +
				"If trying to use your own web UI resources, use the ui-dir flag.\n" +
				"If using Consul version 0.7.0 or later, the web UI is included in the binary so use ui to enable it")
	}
	if rt.DNSUDPAnswerLimit < 0 {
		return fmt.Errorf("dns_config.udp_answer_limit cannot be %d. Must be greater than or equal to zero", rt.DNSUDPAnswerLimit)
	}
	if rt.DNSARecordLimit < 0 {
		return fmt.Errorf("dns_config.a_record_limit cannot be %d. Must be greater than or equal to zero", rt.DNSARecordLimit)
	}
	if err := structs.ValidateMetadata(rt.NodeMeta, false); err != nil {
		return fmt.Errorf("node_meta invalid: %v", err)
	}
	if rt.EncryptKey != "" {
		if _, err := decodeBytes(rt.EncryptKey); err != nil {
			return fmt.Errorf("encrypt has invalid key: %s", err)
		}
		keyfileLAN := filepath.Join(rt.DataDir, SerfLANKeyring)
		if _, err := os.Stat(keyfileLAN); err == nil {
			b.warn("WARNING: LAN keyring exists but -encrypt given, using keyring")
		}
		if rt.ServerMode {
			keyfileWAN := filepath.Join(rt.DataDir, SerfWANKeyring)
			if _, err := os.Stat(keyfileWAN); err == nil {
				b.warn("WARNING: WAN keyring exists but -encrypt given, using keyring")
			}
		}
	}

	// Check the data dir for signs of an un-migrated Consul 0.5.x or older
	// server. Consul refuses to start if this is present to protect a server
	// with existing data from starting on a fresh data set.
	if rt.ServerMode {
		mdbPath := filepath.Join(rt.DataDir, "mdb")
		if _, err := os.Stat(mdbPath); !os.IsNotExist(err) {
			if os.IsPermission(err) {
				return fmt.Errorf(
					"CRITICAL: Permission denied for data folder at %q!\n"+
						"Consul will refuse to boot without access to this directory.\n"+
						"Please correct permissions and try starting again.", mdbPath)
			}
			return fmt.Errorf("CRITICAL: Deprecated data folder found at %q!\n"+
				"Consul will refuse to boot with this directory present.\n"+
				"See https://www.consul.io/docs/upgrade-specific.html for more information.", mdbPath)
		}
	}

	inuse := map[string]string{}
	if err := addrsUnique(inuse, "DNS", rt.DNSAddrs); err != nil {
		// cannot happen since this is the first address
		// we leave this for consistency
		return err
	}
	if err := addrsUnique(inuse, "HTTP", rt.HTTPAddrs); err != nil {
		return err
	}
	if err := addrsUnique(inuse, "HTTPS", rt.HTTPSAddrs); err != nil {
		return err
	}
	if err := addrUnique(inuse, "RPC Advertise", rt.RPCAdvertiseAddr); err != nil {
		return err
	}
	if err := addrUnique(inuse, "Serf Advertise LAN", rt.SerfAdvertiseAddrLAN); err != nil {
		return err
	}
	// Validate serf WAN advertise address only when its set
	if rt.SerfAdvertiseAddrWAN != nil {
		if err := addrUnique(inuse, "Serf Advertise WAN", rt.SerfAdvertiseAddrWAN); err != nil {
			return err
		}
	}
	if b.err != nil {
		return b.err
	}

	// Check for errors in the service definitions
	for _, s := range rt.Services {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("service %q: %s", s.Name, err)
		}
	}

	// Validate the given Connect CA provider config
	validCAProviders := map[string]bool{
		"":                       true,
		structs.ConsulCAProvider: true,
		structs.VaultCAProvider:  true,
	}
	if _, ok := validCAProviders[rt.ConnectCAProvider]; !ok {
		return fmt.Errorf("%s is not a valid CA provider", rt.ConnectCAProvider)
	} else {
		switch rt.ConnectCAProvider {
		case structs.ConsulCAProvider:
			if _, err := ca.ParseConsulCAConfig(rt.ConnectCAConfig); err != nil {
				return err
			}
		case structs.VaultCAProvider:
			if _, err := ca.ParseVaultCAConfig(rt.ConnectCAConfig); err != nil {
				return err
			}
		}
	}

	// ----------------------------------------------------------------
	// warnings
	//

	if rt.ServerMode && !rt.DevMode && !rt.Bootstrap && rt.BootstrapExpect == 2 {
		b.warn(`bootstrap_expect = 2: A cluster with 2 servers will provide no failure tolerance. See https://www.consul.io/docs/internals/consensus.html#deployment-table`)
	}

	if rt.ServerMode && !rt.Bootstrap && rt.BootstrapExpect > 2 && rt.BootstrapExpect%2 == 0 {
		b.warn(`bootstrap_expect is even number: A cluster with an even number of servers does not achieve optimum fault tolerance. See https://www.consul.io/docs/internals/consensus.html#deployment-table`)
	}

	if rt.ServerMode && rt.Bootstrap && rt.BootstrapExpect == 0 {
		b.warn(`bootstrap = true: do not enable unless necessary`)
	}

	if rt.ServerMode && !rt.DevMode && !rt.Bootstrap && rt.BootstrapExpect > 1 {
		b.warn("bootstrap_expect > 0: expecting %d servers", rt.BootstrapExpect)
	}

	return nil
}
