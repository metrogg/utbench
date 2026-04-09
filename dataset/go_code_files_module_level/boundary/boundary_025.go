func loadConfig() (*config, []string, error) {
	// Default config.
	cfg := config{
		ConfigFile:           defaultConfigFile,
		DebugLevel:           defaultLogLevel,
		MaxPeers:             defaultMaxPeers,
		BanDuration:          defaultBanDuration,
		BanThreshold:         defaultBanThreshold,
		RPCMaxClients:        defaultMaxRPCClients,
		RPCMaxWebsockets:     defaultMaxRPCWebsockets,
		RPCMaxConcurrentReqs: defaultMaxRPCConcurrentReqs,
		DataDir:              defaultDataDir,
		LogDir:               defaultLogDir,
		DbType:               defaultDbType,
		RPCKey:               defaultRPCKeyFile,
		RPCCert:              defaultRPCCertFile,
		MinRelayTxFee:        mempool.DefaultMinRelayTxFee.ToBTC(),
		FreeTxRelayLimit:     defaultFreeTxRelayLimit,
		TrickleInterval:      defaultTrickleInterval,
		BlockMinSize:         defaultBlockMinSize,
		BlockMaxSize:         defaultBlockMaxSize,
		BlockMinWeight:       defaultBlockMinWeight,
		BlockMaxWeight:       defaultBlockMaxWeight,
		BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
		MaxOrphanTxs:         defaultMaxOrphanTransactions,
		SigCacheMaxSize:      defaultSigCacheMaxSize,
		Generate:             defaultGenerate,
		TxIndex:              defaultTxIndex,
		AddrIndex:            defaultAddrIndex,
	}

	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.  Any errors aside from the
	// help message error can be ignored here since they will be caught by
	// the final parse below.
	preCfg := cfg
	preParser := newConfigParser(&preCfg, &serviceOpts, flags.HelpFlag)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return nil, nil, err
		}
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", version())
		os.Exit(0)
	}

	// Perform service command and exit if specified.  Invalid service
	// commands show an appropriate error.  Only runs on Windows since
	// the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		err := runServiceCommand(serviceOpts.ServiceCommand)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(0)
	}

	// Load additional config from file.
	var configFileError error
	parser := newConfigParser(&cfg, &serviceOpts, flags.Default)
	if !(preCfg.RegressionTest || preCfg.SimNet) || preCfg.ConfigFile !=
		defaultConfigFile {

		if _, err := os.Stat(preCfg.ConfigFile); os.IsNotExist(err) {
			err := createDefaultConfigFile(preCfg.ConfigFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating a "+
					"default config file: %v\n", err)
			}
		}

		err := flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				fmt.Fprintf(os.Stderr, "Error parsing config "+
					"file: %v\n", err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, err
			}
			configFileError = err
		}
	}

	// Don't add peers from the config file when in regression test mode.
	if preCfg.RegressionTest && len(cfg.AddPeers) > 0 {
		cfg.AddPeers = nil
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			fmt.Fprintln(os.Stderr, usageMessage)
		}
		return nil, nil, err
	}

	// Create the home directory if it doesn't already exist.
	funcName := "loadConfig"
	err = os.MkdirAll(defaultHomeDir, 0700)
	if err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "%s: Failed to create home directory: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, nil, err
	}

	// Multiple networks can't be selected simultaneously.
	numNets := 0
	// Count number of network flags passed; assign active network params
	// while we're at it
	if cfg.TestNet3 {
		numNets++
		activeNetParams = &testNet3Params
	}
	if cfg.RegressionTest {
		numNets++
		activeNetParams = &regressionNetParams
	}
	if cfg.SimNet {
		numNets++
		// Also disable dns seeding on the simulation test network.
		activeNetParams = &simNetParams
		cfg.DisableDNSSeed = true
	}
	if numNets > 1 {
		str := "%s: The testnet, regtest, segnet, and simnet params " +
			"can't be used together -- choose one of the four"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Set the default policy for relaying non-standard transactions
	// according to the default of the active network. The set
	// configuration value takes precedence over the default value for the
	// selected network.
	relayNonStd := activeNetParams.RelayNonStdTxs
	switch {
	case cfg.RelayNonStd && cfg.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used " +
			"together -- choose only one"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	case cfg.RejectNonStd:
		relayNonStd = false
	case cfg.RelayNonStd:
		relayNonStd = true
	}
	cfg.RelayNonStd = relayNonStd

	// Append the network type to the data directory so it is "namespaced"
	// per network.  In addition to the block database, there are other
	// pieces of data that are saved to disk such as address manager state.
	// All data is specific to a network, so namespacing the data directory
	// means each individual piece of serialized data does not have to
	// worry about changing names per network and such.
	cfg.DataDir = cleanAndExpandPath(cfg.DataDir)
	cfg.DataDir = filepath.Join(cfg.DataDir, netName(activeNetParams))

	// Append the network type to the log directory so it is "namespaced"
	// per network in the same fashion as the data directory.
	cfg.LogDir = cleanAndExpandPath(cfg.LogDir)
	cfg.LogDir = filepath.Join(cfg.LogDir, netName(activeNetParams))

	// Special show command to list supported subsystems and exit.
	if cfg.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename))

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		err := fmt.Errorf("%s: %v", funcName, err.Error())
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Validate database type.
	if !validDbType(cfg.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, cfg.DbType, knownDbTypes)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Validate profile port number
	if cfg.Profile != "" {
		profilePort, err := strconv.Atoi(cfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}
	}

	// Don't allow ban durations that are too short.
	if cfg.BanDuration < time.Second {
		str := "%s: The banduration option may not be less than 1s -- parsed [%v]"
		err := fmt.Errorf(str, funcName, cfg.BanDuration)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Validate any given whitelisted IP addresses and networks.
	if len(cfg.Whitelists) > 0 {
		var ip net.IP
		cfg.whitelists = make([]*net.IPNet, 0, len(cfg.Whitelists))

		for _, addr := range cfg.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				ip = net.ParseIP(addr)
				if ip == nil {
					str := "%s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					fmt.Fprintln(os.Stderr, err)
					fmt.Fprintln(os.Stderr, usageMessage)
					return nil, nil, err
				}
				var bits int
				if ip.To4() == nil {
					// IPv6
					bits = 128
				} else {
					bits = 32
				}
				ipnet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(bits, bits),
				}
			}
			cfg.whitelists = append(cfg.whitelists, ipnet)
		}
	}

	// --addPeer and --connect do not mix.
	if len(cfg.AddPeers) > 0 && len(cfg.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// --proxy or --connect without --listen disables listening.
	if (cfg.Proxy != "" || len(cfg.ConnectPeers) > 0) &&
		len(cfg.Listeners) == 0 {
		cfg.DisableListen = true
	}

	// Connect means no DNS seeding.
	if len(cfg.ConnectPeers) > 0 {
		cfg.DisableDNSSeed = true
	}

	// Add the default listener if none were specified. The default
	// listener is all addresses on the listen port for the network
	// we are to connect to.
	if len(cfg.Listeners) == 0 {
		cfg.Listeners = []string{
			net.JoinHostPort("", activeNetParams.DefaultPort),
		}
	}

	// Check to make sure limited and admin users don't have the same username
	if cfg.RPCUser == cfg.RPCLimitUser && cfg.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the " +
			"same username"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Check to make sure limited and admin users don't have the same password
	if cfg.RPCPass == cfg.RPCLimitPass && cfg.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// The RPC server is disabled if no username or password is provided.
	if (cfg.RPCUser == "" || cfg.RPCPass == "") &&
		(cfg.RPCLimitUser == "" || cfg.RPCLimitPass == "") {
		cfg.DisableRPC = true
	}

	if cfg.DisableRPC {
		btcdLog.Infof("RPC service is disabled")
	}

	// Default RPC to listen on localhost only.
	if !cfg.DisableRPC && len(cfg.RPCListeners) == 0 {
		addrs, err := net.LookupHost("localhost")
		if err != nil {
			return nil, nil, err
		}
		cfg.RPCListeners = make([]string, 0, len(addrs))
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, activeNetParams.rpcPort)
			cfg.RPCListeners = append(cfg.RPCListeners, addr)
		}
	}

	if cfg.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may " +
			"not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, cfg.RPCMaxConcurrentReqs)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Validate the the minrelaytxfee.
	cfg.minRelayTxFee, err = btcutil.NewAmount(cfg.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Limit the max block size to a sane value.
	if cfg.BlockMaxSize < blockMaxSizeMin || cfg.BlockMaxSize >
		blockMaxSizeMax {

		str := "%s: The blockmaxsize option must be in between %d " +
			"and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, blockMaxSizeMin,
			blockMaxSizeMax, cfg.BlockMaxSize)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Limit the max block weight to a sane value.
	if cfg.BlockMaxWeight < blockMaxWeightMin ||
		cfg.BlockMaxWeight > blockMaxWeightMax {

		str := "%s: The blockmaxweight option must be in between %d " +
			"and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, blockMaxWeightMin,
			blockMaxWeightMax, cfg.BlockMaxWeight)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Limit the max orphan count to a sane vlue.
	if cfg.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 " +
			"-- parsed [%d]"
		err := fmt.Errorf(str, funcName, cfg.MaxOrphanTxs)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Limit the block priority and minimum block sizes to max block size.
	cfg.BlockPrioritySize = minUint32(cfg.BlockPrioritySize, cfg.BlockMaxSize)
	cfg.BlockMinSize = minUint32(cfg.BlockMinSize, cfg.BlockMaxSize)
	cfg.BlockMinWeight = minUint32(cfg.BlockMinWeight, cfg.BlockMaxWeight)

	switch {
	// If the max block size isn't set, but the max weight is, then we'll
	// set the limit for the max block size to a safe limit so weight takes
	// precedence.
	case cfg.BlockMaxSize == defaultBlockMaxSize &&
		cfg.BlockMaxWeight != defaultBlockMaxWeight:

		cfg.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000

	// If the max block weight isn't set, but the block size is, then we'll
	// scale the set weight accordingly based on the max block size value.
	case cfg.BlockMaxSize != defaultBlockMaxSize &&
		cfg.BlockMaxWeight == defaultBlockMaxWeight:

		cfg.BlockMaxWeight = cfg.BlockMaxSize * blockchain.WitnessScaleFactor
	}

	// Look for illegal characters in the user agent comments.
	for _, uaComment := range cfg.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}
	}

	// --txindex and --droptxindex do not mix.
	if cfg.TxIndex && cfg.DropTxIndex {
		err := fmt.Errorf("%s: the --txindex and --droptxindex "+
			"options may  not be activated at the same time",
			funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// --addrindex and --dropaddrindex do not mix.
	if cfg.AddrIndex && cfg.DropAddrIndex {
		err := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// --addrindex and --droptxindex do not mix.
	if cfg.AddrIndex && cfg.DropTxIndex {
		err := fmt.Errorf("%s: the --addrindex and --droptxindex "+
			"options may not be activated at the same time "+
			"because the address index relies on the transaction "+
			"index",
			funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Check mining addresses are valid and saved parsed versions.
	cfg.miningAddrs = make([]btcutil.Address, 0, len(cfg.MiningAddrs))
	for _, strAddr := range cfg.MiningAddrs {
		addr, err := btcutil.DecodeAddress(strAddr, activeNetParams.Params)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}
		if !addr.IsForNet(activeNetParams.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}
		cfg.miningAddrs = append(cfg.miningAddrs, addr)
	}

	// Ensure there is at least one mining address when the generate flag is
	// set.
	if cfg.Generate && len(cfg.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining " +
			"addresses specified "
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Add default port to all listener addresses if needed and remove
	// duplicate addresses.
	cfg.Listeners = normalizeAddresses(cfg.Listeners,
		activeNetParams.DefaultPort)

	// Add default port to all rpc listener addresses if needed and remove
	// duplicate addresses.
	cfg.RPCListeners = normalizeAddresses(cfg.RPCListeners,
		activeNetParams.rpcPort)

	// Only allow TLS to be disabled if the RPC is bound to localhost
	// addresses.
	if !cfg.DisableRPC && cfg.DisableTLS {
		allowedTLSListeners := map[string]struct{}{
			"localhost": {},
			"127.0.0.1": {},
			"::1":       {},
		}
		for _, addr := range cfg.RPCListeners {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				str := "%s: RPC listen interface '%s' is " +
					"invalid: %v"
				err := fmt.Errorf(str, funcName, addr, err)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, err
			}
			if _, ok := allowedTLSListeners[host]; !ok {
				str := "%s: the --notls option may not be used " +
					"when binding RPC to non localhost " +
					"addresses: %s"
				err := fmt.Errorf(str, funcName, addr)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, err
			}
		}
	}

	// Add default port to all added peer addresses if needed and remove
	// duplicate addresses.
	cfg.AddPeers = normalizeAddresses(cfg.AddPeers,
		activeNetParams.DefaultPort)
	cfg.ConnectPeers = normalizeAddresses(cfg.ConnectPeers,
		activeNetParams.DefaultPort)

	// --noonion and --onion do not mix.
	if cfg.NoOnion && cfg.OnionProxy != "" {
		err := fmt.Errorf("%s: the --noonion and --onion options may "+
			"not be activated at the same time", funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Check the checkpoints for syntax errors.
	cfg.addCheckpoints, err = parseCheckpoints(cfg.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Tor stream isolation requires either proxy or onion proxy to be set.
	if cfg.TorIsolation && cfg.Proxy == "" && cfg.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or " +
			"onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, err
	}

	// Setup dial and DNS resolution (lookup) functions depending on the
	// specified options.  The default is to use the standard
	// net.DialTimeout function as well as the system DNS resolver.  When a
	// proxy is specified, the dial function is set to the proxy specific
	// dial function and the lookup is set to use tor (unless --noonion is
	// specified in which case the system DNS resolver is used).
	cfg.dial = net.DialTimeout
	cfg.lookup = net.LookupIP
	if cfg.Proxy != "" {
		_, _, err := net.SplitHostPort(cfg.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, cfg.Proxy, err)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}

		// Tor isolation flag means proxy credentials will be overridden
		// unless there is also an onion proxy configured in which case
		// that one will be overridden.
		torIsolation := false
		if cfg.TorIsolation && cfg.OnionProxy == "" &&
			(cfg.ProxyUser != "" || cfg.ProxyPass != "") {

			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}

		proxy := &socks.Proxy{
			Addr:         cfg.Proxy,
			Username:     cfg.ProxyUser,
			Password:     cfg.ProxyPass,
			TorIsolation: torIsolation,
		}
		cfg.dial = proxy.DialTimeout

		// Treat the proxy as tor and perform DNS resolution through it
		// unless the --noonion flag is set or there is an
		// onion-specific proxy configured.
		if !cfg.NoOnion && cfg.OnionProxy == "" {
			cfg.lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.Proxy)
			}
		}
	}

	// Setup onion address dial function depending on the specified options.
	// The default is to use the same dial function selected above.  However,
	// when an onion-specific proxy is specified, the onion address dial
	// function is set to use the onion-specific proxy while leaving the
	// normal dial function as selected above.  This allows .onion address
	// traffic to be routed through a different proxy than normal traffic.
	if cfg.OnionProxy != "" {
		_, _, err := net.SplitHostPort(cfg.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, cfg.OnionProxy, err)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, err
		}

		// Tor isolation flag means onion proxy credentials will be
		// overridden.
		if cfg.TorIsolation &&
			(cfg.OnionProxyUser != "" || cfg.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}

		cfg.oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         cfg.OnionProxy,
				Username:     cfg.OnionProxyUser,
				Password:     cfg.OnionProxyPass,
				TorIsolation: cfg.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}

		// When configured in bridge mode (both --onion and --proxy are
		// configured), it means that the proxy configured by --proxy is
		// not a tor proxy, so override the DNS resolution to use the
		// onion-specific proxy.
		if cfg.Proxy != "" {
			cfg.lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.OnionProxy)
			}
		}
	} else {
		cfg.oniondial = cfg.dial
	}

	// Specifying --noonion means the onion address dial function results in
	// an error.
	if cfg.NoOnion {
		cfg.oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}

	// Warn about missing config file only after all other configuration is
	// done.  This prevents the warning on help messages and invalid
	// options.  Note this should go directly before the return.
	if configFileError != nil {
		btcdLog.Warnf("%v", configFileError)
	}

	return &cfg, remainingArgs, nil
}
