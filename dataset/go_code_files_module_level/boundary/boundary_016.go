func loadConfig() (*config, error) {
	defaultCfg := config{
		LndDir:         defaultLndDir,
		ConfigFile:     defaultConfigFile,
		DataDir:        defaultDataDir,
		DebugLevel:     defaultLogLevel,
		TLSCertPath:    defaultTLSCertPath,
		TLSKeyPath:     defaultTLSKeyPath,
		LogDir:         defaultLogDir,
		MaxLogFiles:    defaultMaxLogFiles,
		MaxLogFileSize: defaultMaxLogFileSize,
		Bitcoin: &chainConfig{
			MinHTLC:       defaultBitcoinMinHTLCMSat,
			BaseFee:       defaultBitcoinBaseFeeMSat,
			FeeRate:       defaultBitcoinFeeRate,
			TimeLockDelta: defaultBitcoinTimeLockDelta,
			Node:          "btcd",
		},
		BtcdMode: &btcdConfig{
			Dir:     defaultBtcdDir,
			RPCHost: defaultRPCHost,
			RPCCert: defaultBtcdRPCCertFile,
		},
		BitcoindMode: &bitcoindConfig{
			Dir:     defaultBitcoindDir,
			RPCHost: defaultRPCHost,
		},
		Litecoin: &chainConfig{
			MinHTLC:       defaultLitecoinMinHTLCMSat,
			BaseFee:       defaultLitecoinBaseFeeMSat,
			FeeRate:       defaultLitecoinFeeRate,
			TimeLockDelta: defaultLitecoinTimeLockDelta,
			Node:          "ltcd",
		},
		LtcdMode: &btcdConfig{
			Dir:     defaultLtcdDir,
			RPCHost: defaultRPCHost,
			RPCCert: defaultLtcdRPCCertFile,
		},
		LitecoindMode: &bitcoindConfig{
			Dir:     defaultLitecoindDir,
			RPCHost: defaultRPCHost,
		},
		MaxPendingChannels: defaultMaxPendingChannels,
		NoSeedBackup:       defaultNoSeedBackup,
		MinBackoff:         defaultMinBackoff,
		MaxBackoff:         defaultMaxBackoff,
		SubRPCServers: &subRPCServerConfigs{
			SignRPC: &signrpc.Config{},
		},
		Autopilot: &autoPilotConfig{
			MaxChannels:    5,
			Allocation:     0.6,
			MinChannelSize: int64(minChanFundingSize),
			MaxChannelSize: int64(maxFundingAmount),
			Heuristic: map[string]float64{
				"preferential": 1.0,
			},
		},
		TrickleDelay:             defaultTrickleDelay,
		ChanStatusSampleInterval: defaultChanStatusSampleInterval,
		ChanEnableTimeout:        defaultChanEnableTimeout,
		ChanDisableTimeout:       defaultChanDisableTimeout,
		Alias:                    defaultAlias,
		Color:                    defaultColor,
		MinChanSize:              int64(minChanFundingSize),
		NumGraphSyncPeers:        defaultMinPeers,
		HistoricalSyncInterval:   discovery.DefaultHistoricalSyncInterval,
		Tor: &torConfig{
			SOCKS:   defaultTorSOCKS,
			DNS:     defaultTorDNS,
			Control: defaultTorControl,
		},
		net: &tor.ClearNet{},
		Workers: &lncfg.Workers{
			Read:  lncfg.DefaultReadWorkers,
			Write: lncfg.DefaultWriteWorkers,
			Sig:   lncfg.DefaultSigWorkers,
		},
		Caches: &lncfg.Caches{
			RejectCacheSize:  channeldb.DefaultRejectCacheSize,
			ChannelCacheSize: channeldb.DefaultChannelCacheSize,
		},
	}

	// Pre-parse the command line options to pick up an alternative config
	// file.
	preCfg := defaultCfg
	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", build.Version())
		os.Exit(0)
	}

	// If the config file path has not been modified by the user, then we'll
	// use the default config file path. However, if the user has modified
	// their lnddir, then we should assume they intend to use the config
	// file within it.
	configFileDir := cleanAndExpandPath(preCfg.LndDir)
	configFilePath := cleanAndExpandPath(preCfg.ConfigFile)
	if configFileDir != defaultLndDir {
		if configFilePath == defaultConfigFile {
			configFilePath = filepath.Join(
				configFileDir, defaultConfigFilename,
			)
		}
	}

	// Next, load any additional configuration options from the file.
	var configFileError error
	cfg := preCfg
	if err := flags.IniParse(configFilePath, &cfg); err != nil {
		// If it's a parsing related error, then we'll return
		// immediately, otherwise we can proceed as possibly the config
		// file doesn't exist which is OK.
		if _, ok := err.(*flags.IniError); ok {
			return nil, err
		}

		configFileError = err
	}

	// Finally, parse the remaining command line options again to ensure
	// they take precedence.
	if _, err := flags.Parse(&cfg); err != nil {
		return nil, err
	}

	// If the provided lnd directory is not the default, we'll modify the
	// path to all of the files and directories that will live within it.
	lndDir := cleanAndExpandPath(cfg.LndDir)
	if lndDir != defaultLndDir {
		cfg.DataDir = filepath.Join(lndDir, defaultDataDirname)
		cfg.TLSCertPath = filepath.Join(lndDir, defaultTLSCertFilename)
		cfg.TLSKeyPath = filepath.Join(lndDir, defaultTLSKeyFilename)
		cfg.LogDir = filepath.Join(lndDir, defaultLogDirname)
	}

	// Create the lnd directory if it doesn't already exist.
	funcName := "loadConfig"
	if err := os.MkdirAll(lndDir, 0700); err != nil {
		// Show a nicer error message if it's because a symlink is
		// linked to a directory that does not exist (probably because
		// it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(err) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				err = fmt.Errorf(str, e.Path, link)
			}
		}

		str := "%s: Failed to create lnd directory: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// As soon as we're done parsing configuration options, ensure all paths
	// to directories and files are cleaned and expanded before attempting
	// to use them later on.
	cfg.DataDir = cleanAndExpandPath(cfg.DataDir)
	cfg.TLSCertPath = cleanAndExpandPath(cfg.TLSCertPath)
	cfg.TLSKeyPath = cleanAndExpandPath(cfg.TLSKeyPath)
	cfg.AdminMacPath = cleanAndExpandPath(cfg.AdminMacPath)
	cfg.ReadMacPath = cleanAndExpandPath(cfg.ReadMacPath)
	cfg.InvoiceMacPath = cleanAndExpandPath(cfg.InvoiceMacPath)
	cfg.LogDir = cleanAndExpandPath(cfg.LogDir)
	cfg.BtcdMode.Dir = cleanAndExpandPath(cfg.BtcdMode.Dir)
	cfg.LtcdMode.Dir = cleanAndExpandPath(cfg.LtcdMode.Dir)
	cfg.BitcoindMode.Dir = cleanAndExpandPath(cfg.BitcoindMode.Dir)
	cfg.LitecoindMode.Dir = cleanAndExpandPath(cfg.LitecoindMode.Dir)
	cfg.Tor.PrivateKeyPath = cleanAndExpandPath(cfg.Tor.PrivateKeyPath)

	// Ensure that the user didn't attempt to specify negative values for
	// any of the autopilot params.
	if cfg.Autopilot.MaxChannels < 0 {
		str := "%s: autopilot.maxchannels must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.Allocation < 0 {
		str := "%s: autopilot.allocation must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.MinChannelSize < 0 {
		str := "%s: autopilot.minchansize must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.MaxChannelSize < 0 {
		str := "%s: autopilot.maxchansize must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.MinConfs < 0 {
		str := "%s: autopilot.minconfs must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// Ensure that the specified values for the min and max channel size
	// don't are within the bounds of the normal chan size constraints.
	if cfg.Autopilot.MinChannelSize < int64(minChanFundingSize) {
		cfg.Autopilot.MinChannelSize = int64(minChanFundingSize)
	}
	if cfg.Autopilot.MaxChannelSize > int64(maxFundingAmount) {
		cfg.Autopilot.MaxChannelSize = int64(maxFundingAmount)
	}

	if _, err := validateAtplCfg(cfg.Autopilot); err != nil {
		return nil, err
	}

	// Validate the Tor config parameters.
	socks, err := lncfg.ParseAddressString(
		cfg.Tor.SOCKS, strconv.Itoa(defaultTorSOCKSPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}
	cfg.Tor.SOCKS = socks.String()

	// We'll only attempt to normalize and resolve the DNS host if it hasn't
	// changed, as it doesn't need to be done for the default.
	if cfg.Tor.DNS != defaultTorDNS {
		dns, err := lncfg.ParseAddressString(
			cfg.Tor.DNS, strconv.Itoa(defaultTorDNSPort),
			cfg.net.ResolveTCPAddr,
		)
		if err != nil {
			return nil, err
		}
		cfg.Tor.DNS = dns.String()
	}

	control, err := lncfg.ParseAddressString(
		cfg.Tor.Control, strconv.Itoa(defaultTorControlPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}
	cfg.Tor.Control = control.String()

	switch {
	case cfg.Tor.V2 && cfg.Tor.V3:
		return nil, errors.New("either tor.v2 or tor.v3 can be set, " +
			"but not both")
	case cfg.DisableListen && (cfg.Tor.V2 || cfg.Tor.V3):
		return nil, errors.New("listening must be enabled when " +
			"enabling inbound connections over Tor")
	}

	if cfg.Tor.PrivateKeyPath == "" {
		switch {
		case cfg.Tor.V2:
			cfg.Tor.PrivateKeyPath = filepath.Join(
				lndDir, defaultTorV2PrivateKeyFilename,
			)
		case cfg.Tor.V3:
			cfg.Tor.PrivateKeyPath = filepath.Join(
				lndDir, defaultTorV3PrivateKeyFilename,
			)
		}
	}

	// Set up the network-related functions that will be used throughout
	// the daemon. We use the standard Go "net" package functions by
	// default. If we should be proxying all traffic through Tor, then
	// we'll use the Tor proxy specific functions in order to avoid leaking
	// our real information.
	if cfg.Tor.Active {
		cfg.net = &tor.ProxyNet{
			SOCKS:           cfg.Tor.SOCKS,
			DNS:             cfg.Tor.DNS,
			StreamIsolation: cfg.Tor.StreamIsolation,
		}
	}

	if cfg.DisableListen && cfg.NAT {
		return nil, errors.New("NAT traversal cannot be used when " +
			"listening is disabled")
	}

	// Determine the active chain configuration and its parameters.
	switch {
	// At this moment, multiple active chains are not supported.
	case cfg.Litecoin.Active && cfg.Bitcoin.Active:
		str := "%s: Currently both Bitcoin and Litecoin cannot be " +
			"active together"
		return nil, fmt.Errorf(str, funcName)

	// Either Bitcoin must be active, or Litecoin must be active.
	// Otherwise, we don't know which chain we're on.
	case !cfg.Bitcoin.Active && !cfg.Litecoin.Active:
		return nil, fmt.Errorf("%s: either bitcoin.active or "+
			"litecoin.active must be set to 1 (true)", funcName)

	case cfg.Litecoin.Active:
		if cfg.Litecoin.TimeLockDelta < minTimeLockDelta {
			return nil, fmt.Errorf("timelockdelta must be at least %v",
				minTimeLockDelta)
		}
		// Multiple networks can't be selected simultaneously.  Count
		// number of network flags passed; assign active network params
		// while we're at it.
		numNets := 0
		var ltcParams litecoinNetParams
		if cfg.Litecoin.MainNet {
			numNets++
			ltcParams = litecoinMainNetParams
		}
		if cfg.Litecoin.TestNet3 {
			numNets++
			ltcParams = litecoinTestNetParams
		}
		if cfg.Litecoin.SimNet {
			numNets++
			ltcParams = litecoinSimNetParams
		}
		if cfg.Litecoin.RegTest {
			numNets++
			ltcParams = litecoinRegTestNetParams
		}
		if cfg.Litecoin.SimNet {
			numNets++
			ltcParams = litecoinSimNetParams
		}

		if numNets > 1 {
			str := "%s: The mainnet, testnet, and simnet params " +
				"can't be used together -- choose one of the " +
				"three"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		// The target network must be provided, otherwise, we won't
		// know how to initialize the daemon.
		if numNets == 0 {
			str := "%s: either --litecoin.mainnet, or " +
				"litecoin.testnet must be specified"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		if cfg.Litecoin.MainNet && cfg.DebugHTLC {
			str := "%s: debug-htlc mode cannot be used " +
				"on litecoin mainnet"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		// The litecoin chain is the current active chain. However
		// throughout the codebase we required chaincfg.Params. So as a
		// temporary hack, we'll mutate the default net params for
		// bitcoin with the litecoin specific information.
		applyLitecoinParams(&activeNetParams, &ltcParams)

		switch cfg.Litecoin.Node {
		case "ltcd":
			err := parseRPCParams(cfg.Litecoin, cfg.LtcdMode,
				litecoinChain, funcName)
			if err != nil {
				err := fmt.Errorf("unable to load RPC "+
					"credentials for ltcd: %v", err)
				return nil, err
			}
		case "litecoind":
			if cfg.Litecoin.SimNet {
				return nil, fmt.Errorf("%s: litecoind does not "+
					"support simnet", funcName)
			}
			err := parseRPCParams(cfg.Litecoin, cfg.LitecoindMode,
				litecoinChain, funcName)
			if err != nil {
				err := fmt.Errorf("unable to load RPC "+
					"credentials for litecoind: %v", err)
				return nil, err
			}
		default:
			str := "%s: only ltcd and litecoind mode supported for " +
				"litecoin at this time"
			return nil, fmt.Errorf(str, funcName)
		}

		cfg.Litecoin.ChainDir = filepath.Join(cfg.DataDir,
			defaultChainSubDirname,
			litecoinChain.String())

		// Finally we'll register the litecoin chain as our current
		// primary chain.
		registeredChains.RegisterPrimaryChain(litecoinChain)
		maxFundingAmount = maxLtcFundingAmount
		maxPaymentMSat = maxLtcPaymentMSat

	case cfg.Bitcoin.Active:
		// Multiple networks can't be selected simultaneously.  Count
		// number of network flags passed; assign active network params
		// while we're at it.
		numNets := 0
		if cfg.Bitcoin.MainNet {
			numNets++
			activeNetParams = bitcoinMainNetParams
		}
		if cfg.Bitcoin.TestNet3 {
			numNets++
			activeNetParams = bitcoinTestNetParams
		}
		if cfg.Bitcoin.RegTest {
			numNets++
			activeNetParams = bitcoinRegTestNetParams
		}
		if cfg.Bitcoin.SimNet {
			numNets++
			activeNetParams = bitcoinSimNetParams
		}
		if numNets > 1 {
			str := "%s: The mainnet, testnet, regtest, and " +
				"simnet params can't be used together -- " +
				"choose one of the four"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		// The target network must be provided, otherwise, we won't
		// know how to initialize the daemon.
		if numNets == 0 {
			str := "%s: either --bitcoin.mainnet, or " +
				"bitcoin.testnet, bitcoin.simnet, or bitcoin.regtest " +
				"must be specified"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		if cfg.Bitcoin.MainNet && cfg.DebugHTLC {
			str := "%s: debug-htlc mode cannot be used " +
				"on bitcoin mainnet"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		if cfg.Bitcoin.Node == "neutrino" && cfg.Bitcoin.MainNet {
			str := "%s: neutrino isn't yet supported for " +
				"bitcoin's mainnet"
			err := fmt.Errorf(str, funcName)
			return nil, err
		}

		if cfg.Bitcoin.TimeLockDelta < minTimeLockDelta {
			return nil, fmt.Errorf("timelockdelta must be at least %v",
				minTimeLockDelta)
		}

		switch cfg.Bitcoin.Node {
		case "btcd":
			err := parseRPCParams(
				cfg.Bitcoin, cfg.BtcdMode, bitcoinChain, funcName,
			)
			if err != nil {
				err := fmt.Errorf("unable to load RPC "+
					"credentials for btcd: %v", err)
				return nil, err
			}
		case "bitcoind":
			if cfg.Bitcoin.SimNet {
				return nil, fmt.Errorf("%s: bitcoind does not "+
					"support simnet", funcName)
			}

			err := parseRPCParams(
				cfg.Bitcoin, cfg.BitcoindMode, bitcoinChain, funcName,
			)
			if err != nil {
				err := fmt.Errorf("unable to load RPC "+
					"credentials for bitcoind: %v", err)
				return nil, err
			}
		case "neutrino":
			// No need to get RPC parameters.

		default:
			str := "%s: only btcd, bitcoind, and neutrino mode " +
				"supported for bitcoin at this time"
			return nil, fmt.Errorf(str, funcName)
		}

		cfg.Bitcoin.ChainDir = filepath.Join(cfg.DataDir,
			defaultChainSubDirname,
			bitcoinChain.String())

		// Finally we'll register the bitcoin chain as our current
		// primary chain.
		registeredChains.RegisterPrimaryChain(bitcoinChain)
	}

	// Ensure that the user didn't attempt to specify negative values for
	// any of the autopilot params.
	if cfg.Autopilot.MaxChannels < 0 {
		str := "%s: autopilot.maxchannels must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.Allocation < 0 {
		str := "%s: autopilot.allocation must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.MinChannelSize < 0 {
		str := "%s: autopilot.minchansize must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	if cfg.Autopilot.MaxChannelSize < 0 {
		str := "%s: autopilot.maxchansize must be non-negative"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// Ensure that the specified values for the min and max channel size
	// don't are within the bounds of the normal chan size constraints.
	if cfg.Autopilot.MinChannelSize < int64(minChanFundingSize) {
		cfg.Autopilot.MinChannelSize = int64(minChanFundingSize)
	}
	if cfg.Autopilot.MaxChannelSize > int64(maxFundingAmount) {
		cfg.Autopilot.MaxChannelSize = int64(maxFundingAmount)
	}

	// Validate profile port number.
	if cfg.Profile != "" {
		profilePort, err := strconv.Atoi(cfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, err
		}
	}

	// We'll now construct the network directory which will be where we
	// store all the data specific to this chain/network.
	networkDir = filepath.Join(
		cfg.DataDir, defaultChainSubDirname,
		registeredChains.PrimaryChain().String(),
		normalizeNetwork(activeNetParams.Name),
	)

	// If a custom macaroon directory wasn't specified and the data
	// directory has changed from the default path, then we'll also update
	// the path for the macaroons to be generated.
	if cfg.AdminMacPath == "" {
		cfg.AdminMacPath = filepath.Join(
			networkDir, defaultAdminMacFilename,
		)
	}
	if cfg.ReadMacPath == "" {
		cfg.ReadMacPath = filepath.Join(
			networkDir, defaultReadMacFilename,
		)
	}
	if cfg.InvoiceMacPath == "" {
		cfg.InvoiceMacPath = filepath.Join(
			networkDir, defaultInvoiceMacFilename,
		)
	}

	// Similarly, if a custom back up file path wasn't specified, then
	// we'll update the file location to match our set network directory.
	if cfg.BackupFilePath == "" {
		cfg.BackupFilePath = filepath.Join(
			networkDir, chanbackup.DefaultBackupFileName,
		)
	}

	// Append the network type to the log directory so it is "namespaced"
	// per network in the same fashion as the data directory.
	cfg.LogDir = filepath.Join(cfg.LogDir,
		registeredChains.PrimaryChain().String(),
		normalizeNetwork(activeNetParams.Name))

	// Special show command to list supported subsystems and exit.
	if cfg.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Initialize logging at the default logging level.
	initLogRotator(
		filepath.Join(cfg.LogDir, defaultLogFilename),
		cfg.MaxLogFileSize, cfg.MaxLogFiles,
	)

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		err := fmt.Errorf("%s: %v", funcName, err.Error())
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, err
	}

	// At least one RPCListener is required. So listen on localhost per
	// default.
	if len(cfg.RawRPCListeners) == 0 {
		addr := fmt.Sprintf("localhost:%d", defaultRPCPort)
		cfg.RawRPCListeners = append(cfg.RawRPCListeners, addr)
	}

	// Listen on localhost if no REST listeners were specified.
	if len(cfg.RawRESTListeners) == 0 {
		addr := fmt.Sprintf("localhost:%d", defaultRESTPort)
		cfg.RawRESTListeners = append(cfg.RawRESTListeners, addr)
	}

	// Listen on the default interface/port if no listeners were specified.
	// An empty address string means default interface/address, which on
	// most unix systems is the same as 0.0.0.0. If Tor is active, we
	// default to only listening on localhost for hidden service
	// connections.
	if len(cfg.RawListeners) == 0 {
		addr := fmt.Sprintf(":%d", defaultPeerPort)
		if cfg.Tor.Active {
			addr = fmt.Sprintf("localhost:%d", defaultPeerPort)
		}
		cfg.RawListeners = append(cfg.RawListeners, addr)
	}

	// Add default port to all RPC listener addresses if needed and remove
	// duplicate addresses.
	cfg.RPCListeners, err = lncfg.NormalizeAddresses(
		cfg.RawRPCListeners, strconv.Itoa(defaultRPCPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}

	// Add default port to all REST listener addresses if needed and remove
	// duplicate addresses.
	cfg.RESTListeners, err = lncfg.NormalizeAddresses(
		cfg.RawRESTListeners, strconv.Itoa(defaultRESTPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}

	// For each of the RPC listeners (REST+gRPC), we'll ensure that users
	// have specified a safe combo for authentication. If not, we'll bail
	// out with an error.
	err = lncfg.EnforceSafeAuthentication(
		cfg.RPCListeners, !cfg.NoMacaroons,
	)
	if err != nil {
		return nil, err
	}
	err = lncfg.EnforceSafeAuthentication(
		cfg.RESTListeners, !cfg.NoMacaroons,
	)
	if err != nil {
		return nil, err
	}

	// Remove the listening addresses specified if listening is disabled.
	if cfg.DisableListen {
		ltndLog.Infof("Listening on the p2p interface is disabled!")
		cfg.Listeners = nil
		cfg.ExternalIPs = nil
	} else {

		// Add default port to all listener addresses if needed and remove
		// duplicate addresses.
		cfg.Listeners, err = lncfg.NormalizeAddresses(
			cfg.RawListeners, strconv.Itoa(defaultPeerPort),
			cfg.net.ResolveTCPAddr,
		)
		if err != nil {
			return nil, err
		}

		// Add default port to all external IP addresses if needed and remove
		// duplicate addresses.
		cfg.ExternalIPs, err = lncfg.NormalizeAddresses(
			cfg.RawExternalIPs, strconv.Itoa(defaultPeerPort),
			cfg.net.ResolveTCPAddr,
		)
		if err != nil {
			return nil, err
		}

		// For the p2p port it makes no sense to listen to an Unix socket.
		// Also, we would need to refactor the brontide listener to support
		// that.
		for _, p2pListener := range cfg.Listeners {
			if lncfg.IsUnix(p2pListener) {
				err := fmt.Errorf("unix socket addresses cannot be "+
					"used for the p2p connection listener: %s",
					p2pListener)
				return nil, err
			}
		}
	}

	// Ensure that the specified minimum backoff is below or equal to the
	// maximum backoff.
	if cfg.MinBackoff > cfg.MaxBackoff {
		return nil, fmt.Errorf("maxbackoff must be greater than " +
			"minbackoff")
	}

	// Validate the subconfigs for workers and caches.
	err = lncfg.Validate(
		cfg.Workers,
		cfg.Caches,
	)
	if err != nil {
		return nil, err
	}

	// Finally, ensure that the user's color is correctly formatted,
	// otherwise the server will not be able to start after the unlocking
	// the wallet.
	_, err = parseHexColor(cfg.Color)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse node color: %v", err)
	}

	// Warn about missing config file only after all other configuration is
	// done.  This prevents the warning on help messages and invalid
	// options.  Note this should go directly before the return.
	if configFileError != nil {
		ltndLog.Warnf("%v", configFileError)
	}

	return &cfg, nil
}
