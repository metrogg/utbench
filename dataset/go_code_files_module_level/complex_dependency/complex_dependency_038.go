func NewCore(conf *CoreConfig) (*Core, error) {
	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		if conf.RedirectAddr == "" {
			return nil, fmt.Errorf("missing API address, please set in configuration or via environment")
		}
	}

	if conf.DefaultLeaseTTL == 0 {
		conf.DefaultLeaseTTL = defaultLeaseTTL
	}
	if conf.MaxLeaseTTL == 0 {
		conf.MaxLeaseTTL = maxLeaseTTL
	}
	if conf.DefaultLeaseTTL > conf.MaxLeaseTTL {
		return nil, fmt.Errorf("cannot have DefaultLeaseTTL larger than MaxLeaseTTL")
	}

	// Validate the advertise addr if its given to us
	if conf.RedirectAddr != "" {
		u, err := url.Parse(conf.RedirectAddr)
		if err != nil {
			return nil, errwrap.Wrapf("redirect address is not valid url: {{err}}", err)
		}

		if u.Scheme == "" {
			return nil, fmt.Errorf("redirect address must include scheme (ex. 'http')")
		}
	}

	// Make a default logger if not provided
	if conf.Logger == nil {
		conf.Logger = logging.NewVaultLogger(log.Trace)
	}

	syncInterval := conf.CounterSyncInterval
	if syncInterval.Nanoseconds() == 0 {
		syncInterval = 30 * time.Second
	}

	// Setup the core
	c := &Core{
		entCore:                      entCore{},
		devToken:                     conf.DevToken,
		physical:                     conf.Physical,
		redirectAddr:                 conf.RedirectAddr,
		clusterAddr:                  conf.ClusterAddr,
		seal:                         conf.Seal,
		router:                       NewRouter(),
		sealed:                       new(uint32),
		standby:                      true,
		baseLogger:                   conf.Logger,
		logger:                       conf.Logger.Named("core"),
		defaultLeaseTTL:              conf.DefaultLeaseTTL,
		maxLeaseTTL:                  conf.MaxLeaseTTL,
		cachingDisabled:              conf.DisableCache,
		clusterName:                  conf.ClusterName,
		clusterPeerClusterAddrsCache: cache.New(3*cluster.HeartbeatInterval, time.Second),
		enableMlock:                  !conf.DisableMlock,
		rawEnabled:                   conf.EnableRaw,
		replicationState:             new(uint32),
		atomicPrimaryClusterAddrs:    new(atomic.Value),
		atomicPrimaryFailoverAddrs:   new(atomic.Value),
		localClusterPrivateKey:       new(atomic.Value),
		localClusterCert:             new(atomic.Value),
		localClusterParsedCert:       new(atomic.Value),
		activeNodeReplicationState:   new(uint32),
		keepHALockOnStepDown:         new(uint32),
		replicationFailure:           new(uint32),
		disablePerfStandby:           true,
		activeContextCancelFunc:      new(atomic.Value),
		allLoggers:                   conf.AllLoggers,
		builtinRegistry:              conf.BuiltinRegistry,
		neverBecomeActive:            new(uint32),
		clusterLeaderParams:          new(atomic.Value),
		metricsHelper:                conf.MetricsHelper,
		counters: counters{
			requests:     new(uint64),
			syncInterval: syncInterval,
		},
	}

	atomic.StoreUint32(c.sealed, 1)
	c.allLoggers = append(c.allLoggers, c.logger)

	atomic.StoreUint32(c.replicationState, uint32(consts.ReplicationDRDisabled|consts.ReplicationPerformanceDisabled))
	c.localClusterCert.Store(([]byte)(nil))
	c.localClusterParsedCert.Store((*x509.Certificate)(nil))
	c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

	c.clusterLeaderParams.Store((*ClusterLeaderParams)(nil))

	c.activeContextCancelFunc.Store((context.CancelFunc)(nil))

	if conf.ClusterCipherSuites != "" {
		suites, err := tlsutil.ParseCiphers(conf.ClusterCipherSuites)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing cluster cipher suites: {{err}}", err)
		}
		c.clusterCipherSuites = suites
	}

	// Load CORS config and provide a value for the core field.
	c.corsConfig = &CORSConfig{
		core:    c,
		Enabled: new(uint32),
	}

	if c.seal == nil {
		c.seal = NewDefaultSeal()
	}
	c.seal.SetCore(c)

	if err := coreInit(c, conf); err != nil {
		return nil, err
	}

	if !conf.DisableMlock {
		// Ensure our memory usage is locked into physical RAM
		if err := mlock.LockMemory(); err != nil {
			return nil, fmt.Errorf(
				"Failed to lock memory: %v\n\n"+
					"This usually means that the mlock syscall is not available.\n"+
					"Vault uses mlock to prevent memory from being swapped to\n"+
					"disk. This requires root privileges as well as a machine\n"+
					"that supports mlock. Please enable mlock on your system or\n"+
					"disable Vault from using it. To disable Vault from using it,\n"+
					"set the `disable_mlock` configuration option in your configuration\n"+
					"file.",
				err)
		}
	}

	var err error

	if conf.PluginDirectory != "" {
		c.pluginDirectory, err = filepath.Abs(conf.PluginDirectory)
		if err != nil {
			return nil, errwrap.Wrapf("core setup failed, could not verify plugin directory: {{err}}", err)
		}
	}

	// Construct a new AES-GCM barrier
	c.barrier, err = NewAESGCMBarrier(c.physical)
	if err != nil {
		return nil, errwrap.Wrapf("barrier setup failed: {{err}}", err)
	}

	createSecondaries(c, conf)

	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		c.ha = conf.HAPhysical
	}

	// We create the funcs here, then populate the given config with it so that
	// the caller can share state
	conf.ReloadFuncsLock = &c.reloadFuncsLock
	c.reloadFuncsLock.Lock()
	c.reloadFuncs = make(map[string][]reload.ReloadFunc)
	c.reloadFuncsLock.Unlock()
	conf.ReloadFuncs = &c.reloadFuncs

	logicalBackends := make(map[string]logical.Factory)
	for k, f := range conf.LogicalBackends {
		logicalBackends[k] = f
	}
	_, ok := logicalBackends["kv"]
	if !ok {
		logicalBackends["kv"] = PassthroughBackendFactory
	}

	logicalBackends["cubbyhole"] = CubbyholeBackendFactory
	logicalBackends[systemMountType] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		sysBackendLogger := conf.Logger.Named("system")
		c.AddLogger(sysBackendLogger)
		b := NewSystemBackend(c, sysBackendLogger)
		if err := b.Setup(ctx, config); err != nil {
			return nil, err
		}
		return b, nil
	}
	logicalBackends["identity"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		identityLogger := conf.Logger.Named("identity")
		c.AddLogger(identityLogger)
		return NewIdentityStore(ctx, c, config, identityLogger)
	}
	addExtraLogicalBackends(c, logicalBackends)
	c.logicalBackends = logicalBackends

	credentialBackends := make(map[string]logical.Factory)
	for k, f := range conf.CredentialBackends {
		credentialBackends[k] = f
	}
	credentialBackends["token"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		tsLogger := conf.Logger.Named("token")
		c.AddLogger(tsLogger)
		return NewTokenStore(ctx, tsLogger, c, config)
	}
	addExtraCredentialBackends(c, credentialBackends)
	c.credentialBackends = credentialBackends

	auditBackends := make(map[string]audit.Factory)
	for k, f := range conf.AuditBackends {
		auditBackends[k] = f
	}
	c.auditBackends = auditBackends

	uiStoragePrefix := systemBarrierPrefix + "ui"
	c.uiConfig = NewUIConfig(conf.EnableUI, physical.NewView(c.physical, uiStoragePrefix), NewBarrierView(c.barrier, uiStoragePrefix))

	return c, nil
}
