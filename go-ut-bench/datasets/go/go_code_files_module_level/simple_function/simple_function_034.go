func NewDaemon(ctx context.Context, config *config.Config, pluginStore *plugin.Store) (daemon *Daemon, err error) {
	setDefaultMtu(config)

	registryService, err := registry.NewService(config.ServiceOptions)
	if err != nil {
		return nil, err
	}

	// Ensure that we have a correct root key limit for launching containers.
	if err := ModifyRootKeyLimit(); err != nil {
		logrus.Warnf("unable to modify root key limit, number of containers could be limited by this quota: %v", err)
	}

	// Ensure we have compatible and valid configuration options
	if err := verifyDaemonSettings(config); err != nil {
		return nil, err
	}

	// Do we have a disabled network?
	config.DisableBridge = isBridgeNetworkDisabled(config)

	// Setup the resolv.conf
	setupResolvConf(config)

	// Verify the platform is supported as a daemon
	if !platformSupported {
		return nil, errSystemNotSupported
	}

	// Validate platform-specific requirements
	if err := checkSystem(); err != nil {
		return nil, err
	}

	idMapping, err := setupRemappedRoot(config)
	if err != nil {
		return nil, err
	}
	rootIDs := idMapping.RootPair()
	if err := setupDaemonProcess(config); err != nil {
		return nil, err
	}

	// set up the tmpDir to use a canonical path
	tmp, err := prepareTempDir(config.Root, rootIDs)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the TempDir under %s: %s", config.Root, err)
	}
	realTmp, err := getRealPath(tmp)
	if err != nil {
		return nil, fmt.Errorf("Unable to get the full path to the TempDir (%s): %s", tmp, err)
	}
	if runtime.GOOS == "windows" {
		if _, err := os.Stat(realTmp); err != nil && os.IsNotExist(err) {
			if err := system.MkdirAll(realTmp, 0700, ""); err != nil {
				return nil, fmt.Errorf("Unable to create the TempDir (%s): %s", realTmp, err)
			}
		}
		os.Setenv("TEMP", realTmp)
		os.Setenv("TMP", realTmp)
	} else {
		os.Setenv("TMPDIR", realTmp)
	}

	d := &Daemon{
		configStore: config,
		PluginStore: pluginStore,
		startupDone: make(chan struct{}),
	}
	// Ensure the daemon is properly shutdown if there is a failure during
	// initialization
	defer func() {
		if err != nil {
			if err := d.Shutdown(); err != nil {
				logrus.Error(err)
			}
		}
	}()

	if err := d.setGenericResources(config); err != nil {
		return nil, err
	}
	// set up SIGUSR1 handler on Unix-like systems, or a Win32 global event
	// on Windows to dump Go routine stacks
	stackDumpDir := config.Root
	if execRoot := config.GetExecRoot(); execRoot != "" {
		stackDumpDir = execRoot
	}
	d.setupDumpStackTrap(stackDumpDir)

	if err := d.setupSeccompProfile(); err != nil {
		return nil, err
	}

	// Set the default isolation mode (only applicable on Windows)
	if err := d.setDefaultIsolation(); err != nil {
		return nil, fmt.Errorf("error setting default isolation mode: %v", err)
	}

	if err := configureMaxThreads(config); err != nil {
		logrus.Warnf("Failed to configure golang's threads limit: %v", err)
	}

	// ensureDefaultAppArmorProfile does nothing if apparmor is disabled
	if err := ensureDefaultAppArmorProfile(); err != nil {
		logrus.Errorf(err.Error())
	}

	daemonRepo := filepath.Join(config.Root, "containers")
	if err := idtools.MkdirAllAndChown(daemonRepo, 0700, rootIDs); err != nil {
		return nil, err
	}

	// Create the directory where we'll store the runtime scripts (i.e. in
	// order to support runtimeArgs)
	daemonRuntimes := filepath.Join(config.Root, "runtimes")
	if err := system.MkdirAll(daemonRuntimes, 0700, ""); err != nil {
		return nil, err
	}
	if err := d.loadRuntimes(); err != nil {
		return nil, err
	}

	if runtime.GOOS == "windows" {
		if err := system.MkdirAll(filepath.Join(config.Root, "credentialspecs"), 0, ""); err != nil {
			return nil, err
		}
	}

	// On Windows we don't support the environment variable, or a user supplied graphdriver
	// as Windows has no choice in terms of which graphdrivers to use. It's a case of
	// running Windows containers on Windows - windowsfilter, running Linux containers on Windows,
	// lcow. Unix platforms however run a single graphdriver for all containers, and it can
	// be set through an environment variable, a daemon start parameter, or chosen through
	// initialization of the layerstore through driver priority order for example.
	d.graphDrivers = make(map[string]string)
	layerStores := make(map[string]layer.Store)
	if runtime.GOOS == "windows" {
		d.graphDrivers[runtime.GOOS] = "windowsfilter"
		if system.LCOWSupported() {
			d.graphDrivers["linux"] = "lcow"
		}
	} else {
		driverName := os.Getenv("DOCKER_DRIVER")
		if driverName == "" {
			driverName = config.GraphDriver
		} else {
			logrus.Infof("Setting the storage driver from the $DOCKER_DRIVER environment variable (%s)", driverName)
		}
		d.graphDrivers[runtime.GOOS] = driverName // May still be empty. Layerstore init determines instead.
	}

	d.RegistryService = registryService
	logger.RegisterPluginGetter(d.PluginStore)

	metricsSockPath, err := d.listenMetricsSock()
	if err != nil {
		return nil, err
	}
	registerMetricsPluginCallback(d.PluginStore, metricsSockPath)

	gopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(3 * time.Second),
		grpc.WithDialer(dialer.Dialer),

		// TODO(stevvooe): We may need to allow configuration of this on the client.
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaults.DefaultMaxRecvMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaults.DefaultMaxSendMsgSize)),
	}
	if config.ContainerdAddr != "" {
		d.containerdCli, err = containerd.New(config.ContainerdAddr, containerd.WithDefaultNamespace(ContainersNamespace), containerd.WithDialOpts(gopts), containerd.WithTimeout(60*time.Second))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to dial %q", config.ContainerdAddr)
		}
	}

	createPluginExec := func(m *plugin.Manager) (plugin.Executor, error) {
		var pluginCli *containerd.Client

		// Windows is not currently using containerd, keep the
		// client as nil
		if config.ContainerdAddr != "" {
			pluginCli, err = containerd.New(config.ContainerdAddr, containerd.WithDefaultNamespace(pluginexec.PluginNamespace), containerd.WithDialOpts(gopts), containerd.WithTimeout(60*time.Second))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to dial %q", config.ContainerdAddr)
			}
		}

		return pluginexec.New(ctx, getPluginExecRoot(config.Root), pluginCli, m)
	}

	// Plugin system initialization should happen before restore. Do not change order.
	d.pluginManager, err = plugin.NewManager(plugin.ManagerConfig{
		Root:               filepath.Join(config.Root, "plugins"),
		ExecRoot:           getPluginExecRoot(config.Root),
		Store:              d.PluginStore,
		CreateExecutor:     createPluginExec,
		RegistryService:    registryService,
		LiveRestoreEnabled: config.LiveRestoreEnabled,
		LogPluginEvent:     d.LogPluginEvent, // todo: make private
		AuthzMiddleware:    config.AuthzMiddleware,
	})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create plugin manager")
	}

	if err := d.setupDefaultLogConfig(); err != nil {
		return nil, err
	}

	for operatingSystem, gd := range d.graphDrivers {
		layerStores[operatingSystem], err = layer.NewStoreFromOptions(layer.StoreOptions{
			Root:                      config.Root,
			MetadataStorePathTemplate: filepath.Join(config.Root, "image", "%s", "layerdb"),
			GraphDriver:               gd,
			GraphDriverOptions:        config.GraphOptions,
			IDMapping:                 idMapping,
			PluginGetter:              d.PluginStore,
			ExperimentalEnabled:       config.Experimental,
			OS:                        operatingSystem,
		})
		if err != nil {
			return nil, err
		}

		// As layerstore initialization may set the driver
		d.graphDrivers[operatingSystem] = layerStores[operatingSystem].DriverName()
	}

	// Configure and validate the kernels security support. Note this is a Linux/FreeBSD
	// operation only, so it is safe to pass *just* the runtime OS graphdriver.
	if err := configureKernelSecuritySupport(config, d.graphDrivers[runtime.GOOS]); err != nil {
		return nil, err
	}

	imageRoot := filepath.Join(config.Root, "image", d.graphDrivers[runtime.GOOS])
	ifs, err := image.NewFSStoreBackend(filepath.Join(imageRoot, "imagedb"))
	if err != nil {
		return nil, err
	}

	lgrMap := make(map[string]image.LayerGetReleaser)
	for os, ls := range layerStores {
		lgrMap[os] = ls
	}
	imageStore, err := image.NewImageStore(ifs, lgrMap)
	if err != nil {
		return nil, err
	}

	d.volumes, err = volumesservice.NewVolumeService(config.Root, d.PluginStore, rootIDs, d)
	if err != nil {
		return nil, err
	}

	uuid, err := loadOrCreateUUID(filepath.Join(config.Root, "engine_uuid"))
	if err != nil {
		return nil, err
	}

	trustDir := filepath.Join(config.Root, "trust")

	if err := system.MkdirAll(trustDir, 0700, ""); err != nil {
		return nil, err
	}

	// We have a single tag/reference store for the daemon globally. However, it's
	// stored under the graphdriver. On host platforms which only support a single
	// container OS, but multiple selectable graphdrivers, this means depending on which
	// graphdriver is chosen, the global reference store is under there. For
	// platforms which support multiple container operating systems, this is slightly
	// more problematic as where does the global ref store get located? Fortunately,
	// for Windows, which is currently the only daemon supporting multiple container
	// operating systems, the list of graphdrivers available isn't user configurable.
	// For backwards compatibility, we just put it under the windowsfilter
	// directory regardless.
	refStoreLocation := filepath.Join(imageRoot, `repositories.json`)
	rs, err := refstore.NewReferenceStore(refStoreLocation)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create reference store repository: %s", err)
	}

	distributionMetadataStore, err := dmetadata.NewFSMetadataStore(filepath.Join(imageRoot, "distribution"))
	if err != nil {
		return nil, err
	}

	// Discovery is only enabled when the daemon is launched with an address to advertise.  When
	// initialized, the daemon is registered and we can store the discovery backend as it's read-only
	if err := d.initDiscovery(config); err != nil {
		return nil, err
	}

	sysInfo := sysinfo.New(false)
	// Check if Devices cgroup is mounted, it is hard requirement for container security,
	// on Linux.
	if runtime.GOOS == "linux" && !sysInfo.CgroupDevicesEnabled {
		return nil, errors.New("Devices cgroup isn't mounted")
	}

	d.ID = uuid
	d.repository = daemonRepo
	d.containers = container.NewMemoryStore()
	if d.containersReplica, err = container.NewViewDB(); err != nil {
		return nil, err
	}
	d.execCommands = exec.NewStore()
	d.idIndex = truncindex.NewTruncIndex([]string{})
	d.statsCollector = d.newStatsCollector(1 * time.Second)

	d.EventsService = events.New()
	d.root = config.Root
	d.idMapping = idMapping
	d.seccompEnabled = sysInfo.Seccomp
	d.apparmorEnabled = sysInfo.AppArmor

	d.linkIndex = newLinkIndex()

	// TODO: imageStore, distributionMetadataStore, and ReferenceStore are only
	// used above to run migration. They could be initialized in ImageService
	// if migration is called from daemon/images. layerStore might move as well.
	d.imageService = images.NewImageService(images.ImageServiceConfig{
		ContainerStore:            d.containers,
		DistributionMetadataStore: distributionMetadataStore,
		EventsService:             d.EventsService,
		ImageStore:                imageStore,
		LayerStores:               layerStores,
		MaxConcurrentDownloads:    *config.MaxConcurrentDownloads,
		MaxConcurrentUploads:      *config.MaxConcurrentUploads,
		ReferenceStore:            rs,
		RegistryService:           registryService,
	})

	go d.execCommandGC()

	d.containerd, err = libcontainerd.NewClient(ctx, d.containerdCli, filepath.Join(config.ExecRoot, "containerd"), ContainersNamespace, d)
	if err != nil {
		return nil, err
	}

	if err := d.restore(); err != nil {
		return nil, err
	}
	close(d.startupDone)

	// FIXME: this method never returns an error
	info, _ := d.SystemInfo()

	engineInfo.WithValues(
		dockerversion.Version,
		dockerversion.GitCommit,
		info.Architecture,
		info.Driver,
		info.KernelVersion,
		info.OperatingSystem,
		info.OSType,
		info.ID,
	).Set(1)
	engineCpus.Set(float64(info.NCPU))
	engineMemory.Set(float64(info.MemTotal))

	gd := ""
	for os, driver := range d.graphDrivers {
		if len(gd) > 0 {
			gd += ", "
		}
		gd += driver
		if len(d.graphDrivers) > 1 {
			gd = fmt.Sprintf("%s (%s)", gd, os)
		}
	}
	logrus.WithFields(logrus.Fields{
		"version":        dockerversion.Version,
		"commit":         dockerversion.GitCommit,
		"graphdriver(s)": gd,
	}).Info("Docker daemon")

	return d, nil
}
