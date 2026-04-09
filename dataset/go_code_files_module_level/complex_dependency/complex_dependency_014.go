func commonManifolds(config ManifoldsConfig) dependency.Manifolds {

	// connectFilter exists:
	//  1) to let us retry api connections immediately on password change,
	//     rather than causing the dependency engine to wait for a while;
	//  2) to ensure that certain connection failures correctly trigger
	//     complete agent removal. (It's not safe to let any agent other
	//     than the machine mess around with SetCanUninstall).
	connectFilter := func(err error) error {
		cause := errors.Cause(err)
		if cause == apicaller.ErrConnectImpossible {
			err2 := coreagent.SetCanUninstall(config.Agent)
			if err2 != nil {
				return errors.Trace(err2)
			}
			return jworker.ErrTerminateAgent
		} else if cause == apicaller.ErrChangedPassword {
			return dependency.ErrBounce
		}
		return err
	}

	newExternalControllerWatcherClient := func(apiInfo *api.Info) (
		externalcontrollerupdater.ExternalControllerWatcherClientCloser, error,
	) {
		conn, err := apicaller.NewExternalControllerConnection(apiInfo)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return crosscontroller.NewClient(conn), nil
	}

	agentConfig := config.Agent.CurrentConfig()
	machineTag := agentConfig.Tag().(names.MachineTag)
	controllerTag := agentConfig.Controller()

	leaseFSM := raftlease.NewFSM()

	manifolds := dependency.Manifolds{
		// The agent manifold references the enclosing agent, and is the
		// foundation stone on which most other manifolds ultimately depend.
		agentName: agent.Manifold(config.Agent),

		// The termination worker returns ErrTerminateAgent if a
		// termination signal is received by the process it's running
		// in. It has no inputs and its only output is the error it
		// returns. It depends on the uninstall file having been
		// written *by the manual provider* at install time; it would
		// be Very Wrong Indeed to use SetCanUninstall in conjunction
		// with this code.
		terminationName: terminationworker.Manifold(),

		clockName: clockManifold(config.Clock),

		// Each machine agent has a flag manifold/worker which
		// reports whether or not the agent is a controller.
		isControllerFlagName: isControllerFlagManifold(),

		// The stateconfigwatcher manifold watches the machine agent's
		// configuration and reports if state serving info is
		// present. It will bounce itself if state serving info is
		// added or removed. It is intended as a dependency just for
		// the state manifold.
		stateConfigWatcherName: stateconfigwatcher.Manifold(stateconfigwatcher.ManifoldConfig{
			AgentName:          agentName,
			AgentConfigChanged: config.AgentConfigChanged,
		}),

		// The centralhub manifold watches the state config to make sure it
		// only starts for machines that are api servers. Currently the hub is
		// passed in as config, but when the apiserver and peergrouper are
		// updated to use the dependency engine, the centralhub manifold
		// should also take the agentName so the worker can get the machine ID
		// for the creation of the hub.
		centralHubName: centralhub.Manifold(centralhub.ManifoldConfig{
			StateConfigWatcherName: stateConfigWatcherName,
			Hub:                    config.CentralHub,
		}),

		// The pubsub manifold gets the APIInfo from the agent config,
		// and uses this as a basis to talk to the other API servers.
		// The worker subscribes to the messages sent by the peergrouper
		// that defines the set of machines that are the API servers.
		// All non-local messages that originate from the machine that
		// is running the worker get forwarded to the other API servers.
		// This worker does not run in non-API server machines through
		// the hub dependency, as that is only available if the machine
		// is an API server.
		pubSubName: psworker.Manifold(psworker.ManifoldConfig{
			AgentName:      agentName,
			CentralHubName: centralHubName,
			Clock:          config.Clock,
			Logger:         loggo.GetLogger("juju.worker.pubsub"),
			NewWorker:      psworker.NewWorker,
			Reporter:       config.PubSubReporter,
		}),

		// The presence manifold listens to pubsub messages about the pubsub
		// forwarding connections and api connection and disconnections to
		// establish a view on which agents are "alive".
		presenceName: prworker.Manifold(prworker.ManifoldConfig{
			AgentName:              agentName,
			CentralHubName:         centralHubName,
			StateConfigWatcherName: stateConfigWatcherName,
			Recorder:               config.PresenceRecorder,
			Logger:                 loggo.GetLogger("juju.worker.presence"),
			NewWorker:              prworker.NewWorker,
		}),

		/* TODO(menn0) - this is currently unused, pending further
		 * refactoring in the state package.

			// The controller manifold creates a *state.Controller and
			// makes it available to other manifolds. It pings the MongoDB
			// session regularly and will die if pings fail.
			controllerName: workercontroller.Manifold(workercontroller.ManifoldConfig{
				AgentName:              agentName,
				StateConfigWatcherName: stateConfigWatcherName,
				OpenController:         config.OpenController,
			}),
		*/

		// The state manifold creates a *state.State and makes it
		// available to other manifolds. It pings the mongodb session
		// regularly and will die if pings fail.
		stateName: workerstate.Manifold(workerstate.ManifoldConfig{
			AgentName:              agentName,
			StateConfigWatcherName: stateConfigWatcherName,
			OpenStatePool:          config.OpenStatePool,
			SetStatePool:           config.SetStatePool,
		}),

		// The modelcache manifold creates a cache.Controller and keeps
		// it up to date using an all model watcher. The controller is then
		// used by the apiserver.
		modelCacheName: modelcache.Manifold(modelcache.ManifoldConfig{
			StateName:            stateName,
			Logger:               loggo.GetLogger("juju.worker.modelcache"),
			PrometheusRegisterer: config.PrometheusRegisterer,
			NewWorker:            modelcache.NewWorker,
		}),

		// The api-config-watcher manifold monitors the API server
		// addresses in the agent config and bounces when they
		// change. It's required as part of model migrations.
		apiConfigWatcherName: apiconfigwatcher.Manifold(apiconfigwatcher.ManifoldConfig{
			AgentName:          agentName,
			AgentConfigChanged: config.AgentConfigChanged,
		}),

		// The certificate-watcher manifold monitors the API server
		// certificate in the agent config for changes, and parses
		// and offers the result to other manifolds. This is only
		// run by state servers.
		certificateWatcherName: ifController(apiservercertwatcher.Manifold(apiservercertwatcher.ManifoldConfig{
			AgentName:          agentName,
			AgentConfigChanged: config.AgentConfigChanged,
		})),

		// The api caller is a thin concurrent wrapper around a connection
		// to some API server. It's used by many other manifolds, which all
		// select their own desired facades. It will be interesting to see
		// how this works when we consolidate the agents; might be best to
		// handle the auth changes server-side..?
		apiCallerName: apicaller.Manifold(apicaller.ManifoldConfig{
			AgentName:            agentName,
			APIConfigWatcherName: apiConfigWatcherName,
			APIOpen:              api.Open,
			NewConnection:        apicaller.ScaryConnect,
			Filter:               connectFilter,
		}),

		// The upgrade steps gate is used to coordinate workers which
		// shouldn't do anything until the upgrade-steps worker has
		// finished running any required upgrade steps. The flag of
		// similar name is used to implement the isFullyUpgraded func
		// that keeps upgrade concerns out of unrelated manifolds.
		upgradeStepsGateName: gate.ManifoldEx(config.UpgradeStepsLock),
		upgradeStepsFlagName: gate.FlagManifold(gate.FlagManifoldConfig{
			GateName:  upgradeStepsGateName,
			NewWorker: gate.NewFlagWorker,
		}),

		// The upgrade check gate is used to coordinate workers which
		// shouldn't do anything until the upgrader worker has
		// completed its first check for a new tools version to
		// upgrade to. The flag of similar name is used to implement
		// the isFullyUpgraded func that keeps upgrade concerns out of
		// unrelated manifolds.
		upgradeCheckGateName: gate.ManifoldEx(config.UpgradeCheckLock),
		upgradeCheckFlagName: gate.FlagManifold(gate.FlagManifoldConfig{
			GateName:  upgradeCheckGateName,
			NewWorker: gate.NewFlagWorker,
		}),

		// The upgradesteps worker runs soon after the machine agent
		// starts and runs any steps required to upgrade to the
		// running jujud version. Once upgrade steps have run, the
		// upgradesteps gate is unlocked and the worker exits.
		upgradeStepsName: upgradesteps.Manifold(upgradesteps.ManifoldConfig{
			AgentName:            agentName,
			APICallerName:        apiCallerName,
			UpgradeStepsGateName: upgradeStepsGateName,
			OpenStateForUpgrade:  config.OpenStateForUpgrade,
			PreUpgradeSteps:      config.PreUpgradeSteps,
			NewAgentStatusSetter: config.NewAgentStatusSetter,
		}),

		// The migration workers collaborate to run migrations;
		// and to create a mechanism for running other workers
		// so they can't accidentally interfere with a migration
		// in progress. Such a manifold should (1) depend on the
		// migration-inactive flag, to know when to start or die;
		// and (2) occupy the migration-fortress, so as to avoid
		// possible interference with the minion (which will not
		// take action until it's gained sole control of the
		// fortress).
		//
		// Note that the fortress itself will not be created
		// until the upgrade process is complete; this frees all
		// its dependencies from upgrade concerns.
		migrationFortressName: ifFullyUpgraded(fortress.Manifold()),
		migrationInactiveFlagName: migrationflag.Manifold(migrationflag.ManifoldConfig{
			APICallerName: apiCallerName,
			Check:         migrationflag.IsTerminal,
			NewFacade:     migrationflag.NewFacade,
			NewWorker:     migrationflag.NewWorker,
		}),
		migrationMinionName: migrationminion.Manifold(migrationminion.ManifoldConfig{
			AgentName:         agentName,
			APICallerName:     apiCallerName,
			FortressName:      migrationFortressName,
			APIOpen:           api.Open,
			ValidateMigration: config.ValidateMigration,
			NewFacade:         migrationminion.NewFacade,
			NewWorker:         migrationminion.NewWorker,
		}),

		// We run clock updaters for every controller machine to
		// ensure the lease clock is updated monotonically and at a
		// rate no faster than real time.
		//
		// If the legacy-leases feature flag is set the global clock
		// updater updates the lease clock in the database.  .
		globalClockUpdaterName: ifLegacyLeasesEnabled(globalclockupdater.Manifold(globalclockupdater.ManifoldConfig{
			ClockName:      clockName,
			StateName:      stateName,
			NewWorker:      globalclockupdater.NewWorker,
			UpdateInterval: globalClockUpdaterUpdateInterval,
			BackoffDelay:   globalClockUpdaterBackoffDelay,
			Logger:         loggo.GetLogger("juju.worker.globalclockupdater.mongo"),
		})),
		// We also run another clock updater to feed time updates into
		// the lease FSM.
		leaseClockUpdaterName: globalclockupdater.Manifold(globalclockupdater.ManifoldConfig{
			ClockName:        clockName,
			LeaseManagerName: leaseManagerName,
			RaftName:         raftForwarderName,
			NewWorker:        globalclockupdater.NewWorker,
			UpdateInterval:   globalClockUpdaterUpdateInterval,
			BackoffDelay:     globalClockUpdaterBackoffDelay,
			Logger:           loggo.GetLogger("juju.worker.globalclockupdater.raft"),
		}),

		// Each controller machine runs a singular worker which will
		// attempt to claim responsibility for running certain workers
		// that must not be run concurrently by multiple agents.
		isPrimaryControllerFlagName: ifController(singular.Manifold(singular.ManifoldConfig{
			ClockName:     clockName,
			APICallerName: apiCallerName,
			Duration:      config.ControllerLeaseDuration,
			Claimant:      machineTag,
			Entity:        controllerTag,
			NewFacade:     singular.NewFacade,
			NewWorker:     singular.NewWorker,
		})),

		// The agent-config-updater manifold sets the state serving info from
		// the API connection and writes it to the agent config.
		agentConfigUpdaterName: ifNotMigrating(agentconfigupdater.Manifold(agentconfigupdater.ManifoldConfig{
			AgentName:      agentName,
			APICallerName:  apiCallerName,
			CentralHubName: centralHubName,
			Logger:         loggo.GetLogger("juju.worker.agentconfigupdater"),
		})),

		// The apiworkers manifold starts workers which rely on the
		// machine agent's API connection but have not been converted
		// to work directly under the dependency engine. It waits for
		// upgrades to be finished before starting these workers.
		apiWorkersName: ifNotMigrating(APIWorkersManifold(APIWorkersConfig{
			APICallerName:   apiCallerName,
			StartAPIWorkers: config.StartAPIWorkers,
		})),

		// The logging config updater is a leaf worker that indirectly
		// controls the messages sent via the log sender or rsyslog,
		// according to changes in environment config. We should only need
		// one of these in a consolidated agent.
		loggingConfigUpdaterName: ifNotMigrating(logger.Manifold(logger.ManifoldConfig{
			AgentName:       agentName,
			APICallerName:   apiCallerName,
			UpdateAgentFunc: config.UpdateLoggerConfig,
		})),

		// The diskmanager worker periodically lists block devices on the
		// machine it runs on. This worker will be run on all Juju-managed
		// machines (one per machine agent).
		diskManagerName: ifNotMigrating(diskmanager.Manifold(diskmanager.ManifoldConfig{
			AgentName:     agentName,
			APICallerName: apiCallerName,
		})),

		// The api address updater is a leaf worker that rewrites agent config
		// as the state server addresses change. We should only need one of
		// these in a consolidated agent.
		apiAddressUpdaterName: ifNotMigrating(apiaddressupdater.Manifold(apiaddressupdater.ManifoldConfig{
			AgentName:     agentName,
			APICallerName: apiCallerName,
		})),

		// The machiner Worker will wait for the identified machine to become
		// Dying and make it Dead; or until the machine becomes Dead by other
		// means. This worker needs to be launched after fanconfigurer
		// so that it reports interfaces created by it.
		machinerName: ifNotMigrating(machiner.Manifold(machiner.ManifoldConfig{
			AgentName:         agentName,
			APICallerName:     apiCallerName,
			FanConfigurerName: fanConfigurerName,
		})),

		// The log sender is a leaf worker that sends log messages to some
		// API server, when configured so to do. We should only need one of
		// these in a consolidated agent.
		//
		// NOTE: the LogSource will buffer a large number of messages as an upgrade
		// runs; it currently seems better to fill the buffer and send when stable,
		// optimising for stable controller upgrades rather than up-to-the-moment
		// observable normal-machine upgrades.
		logSenderName: ifNotMigrating(logsender.Manifold(logsender.ManifoldConfig{
			APICallerName: apiCallerName,
			LogSource:     config.LogSource,
		})),

		resumerName: ifNotMigrating(resumer.Manifold(resumer.ManifoldConfig{
			AgentName:     agentName,
			APICallerName: apiCallerName,
			Clock:         config.Clock,
			Interval:      time.Minute,
			NewFacade:     resumer.NewFacade,
			NewWorker:     resumer.NewWorker,
		})),

		identityFileWriterName: ifNotMigrating(identityfilewriter.Manifold(identityfilewriter.ManifoldConfig{
			AgentName:     agentName,
			APICallerName: apiCallerName,
		})),

		machineActionName: ifNotMigrating(machineactions.Manifold(machineactions.ManifoldConfig{
			AgentName:     agentName,
			APICallerName: apiCallerName,
			NewFacade:     machineactions.NewFacade,
			NewWorker:     machineactions.NewMachineActionsWorker,
		})),

		externalControllerUpdaterName: ifNotMigrating(ifPrimaryController(externalcontrollerupdater.Manifold(
			externalcontrollerupdater.ManifoldConfig{
				APICallerName:                      apiCallerName,
				NewExternalControllerWatcherClient: newExternalControllerWatcherClient,
			},
		))),

		logPrunerName: ifNotMigrating(ifPrimaryController(dblogpruner.Manifold(
			dblogpruner.ManifoldConfig{
				ClockName:     clockName,
				StateName:     stateName,
				PruneInterval: config.LogPruneInterval,
				NewWorker:     dblogpruner.NewWorker,
			},
		))),

		txnPrunerName: ifNotMigrating(ifPrimaryController(txnpruner.Manifold(
			txnpruner.ManifoldConfig{
				ClockName:     clockName,
				StateName:     stateName,
				PruneInterval: config.TransactionPruneInterval,
				NewWorker:     txnpruner.New,
			},
		))),

		httpServerArgsName: httpserverargs.Manifold(httpserverargs.ManifoldConfig{
			ClockName:             clockName,
			ControllerPortName:    controllerPortName,
			StateName:             stateName,
			NewStateAuthenticator: httpserverargs.NewStateAuthenticator,
		}),

		// TODO Juju 3.0: the controller port worker is only needed while
		// the controller port is a mutable controller config value.
		// When we hit 3.0 we should make controller-port a required
		// and unmutable value.
		controllerPortName: controllerport.Manifold(controllerport.ManifoldConfig{
			AgentName:               agentName,
			HubName:                 centralHubName,
			StateName:               stateName,
			Logger:                  loggo.GetLogger("juju.worker.controllerport"),
			UpdateControllerAPIPort: config.UpdateControllerAPIPort,
			GetControllerConfig:     controllerport.GetControllerConfig,
			NewWorker:               controllerport.NewWorker,
		}),

		httpServerName: httpserver.Manifold(httpserver.ManifoldConfig{
			CertWatcherName:      certificateWatcherName,
			HubName:              centralHubName,
			StateName:            stateName,
			MuxName:              httpServerArgsName,
			APIServerName:        apiServerName,
			RaftTransportName:    raftTransportName,
			PrometheusRegisterer: config.PrometheusRegisterer,
			AgentName:            config.AgentName,
			Clock:                config.Clock,
			MuxShutdownWait:      config.MuxShutdownWait,
			LogDir:               agentConfig.LogDir(),
			GetControllerConfig:  httpserver.GetControllerConfig,
			NewTLSConfig:         httpserver.NewTLSConfig,
			NewWorker:            httpserver.NewWorkerShim,
		}),

		apiServerName: apiserver.Manifold(apiserver.ManifoldConfig{
			AgentName:              agentName,
			AuthenticatorName:      httpServerArgsName,
			ClockName:              clockName,
			StateName:              stateName,
			ModelCacheName:         modelCacheName,
			MuxName:                httpServerArgsName,
			LeaseManagerName:       leaseManagerName,
			UpgradeGateName:        upgradeStepsGateName,
			RestoreStatusName:      restoreWatcherName,
			AuditConfigUpdaterName: auditConfigUpdaterName,
			// Synthetic dependency - if raft-transport bounces we
			// need to bounce api-server too, otherwise http-server
			// can't shutdown properly.
			RaftTransportName: raftTransportName,

			PrometheusRegisterer:              config.PrometheusRegisterer,
			RegisterIntrospectionHTTPHandlers: config.RegisterIntrospectionHTTPHandlers,
			Hub:                               config.CentralHub,
			Presence:                          config.PresenceRecorder,
			NewWorker:                         apiserver.NewWorker,
			NewMetricsCollector:               apiserver.NewMetricsCollector,
		}),

		modelWorkerManagerName: ifFullyUpgraded(modelworkermanager.Manifold(modelworkermanager.ManifoldConfig{
			StateName:      stateName,
			NewWorker:      modelworkermanager.New,
			NewModelWorker: config.NewModelWorker,
		})),

		peergrouperName: ifFullyUpgraded(peergrouper.Manifold(peergrouper.ManifoldConfig{
			AgentName:          agentName,
			ClockName:          clockName,
			ControllerPortName: controllerPortName,
			StateName:          stateName,
			Hub:                config.CentralHub,
			NewWorker:          peergrouper.New,
		})),

		restoreWatcherName: restorewatcher.Manifold(restorewatcher.ManifoldConfig{
			StateName: stateName,
			NewWorker: restorewatcher.NewWorker,
		}),

		auditConfigUpdaterName: ifController(auditconfigupdater.Manifold(auditconfigupdater.ManifoldConfig{
			AgentName: agentName,
			StateName: stateName,
			NewWorker: auditconfigupdater.New,
		})),

		raftTransportName: ifController(rafttransport.Manifold(rafttransport.ManifoldConfig{
			ClockName:         clockName,
			AgentName:         agentName,
			AuthenticatorName: httpServerArgsName,
			HubName:           centralHubName,
			MuxName:           httpServerArgsName,
			DialConn:          rafttransport.DialConn,
			NewWorker:         rafttransport.NewWorker,
			Path:              "/raft",
		})),

		raftName: ifFullyUpgraded(raft.Manifold(raft.ManifoldConfig{
			ClockName:            clockName,
			AgentName:            agentName,
			TransportName:        raftTransportName,
			FSM:                  leaseFSM,
			Logger:               loggo.GetLogger("juju.worker.raft"),
			PrometheusRegisterer: config.PrometheusRegisterer,
			NewWorker:            raft.NewWorker,
		})),

		raftFlagName: raftflag.Manifold(raftflag.ManifoldConfig{
			RaftName:  raftName,
			NewWorker: raftflag.NewWorker,
		}),

		// The raft clusterer can only run on the raft leader, since
		// it makes configuration updates based on changes in API
		// server details.
		raftClustererName: ifRaftLeader(raftclusterer.Manifold(raftclusterer.ManifoldConfig{
			RaftName:       raftName,
			CentralHubName: centralHubName,
			NewWorker:      raftclusterer.NewWorker,
		})),

		raftBackstopName: raftbackstop.Manifold(raftbackstop.ManifoldConfig{
			RaftName:       raftName,
			CentralHubName: centralHubName,
			AgentName:      agentName,
			Logger:         loggo.GetLogger("juju.worker.raft.raftbackstop"),
			NewWorker:      raftbackstop.NewWorker,
		}),

		// The raft forwarder accepts FSM commands from the hub and
		// applies them to the raft leader.
		raftForwarderName: ifRaftLeader(raftforwarder.Manifold(raftforwarder.ManifoldConfig{
			AgentName:      agentName,
			RaftName:       raftName,
			StateName:      stateName,
			CentralHubName: centralHubName,
			RequestTopic:   leaseRequestTopic,
			Logger:         loggo.GetLogger("juju.worker.raft.raftforwarder"),
			NewWorker:      raftforwarder.NewWorker,
			NewTarget:      raftforwarder.NewTarget,
		})),

		// The global lease manager tracks lease information in the raft
		// cluster rather than in mongo.
		leaseManagerName: ifController(leasemanager.Manifold(leasemanager.ManifoldConfig{
			AgentName:            agentName,
			ClockName:            clockName,
			CentralHubName:       centralHubName,
			StateName:            stateName,
			FSM:                  leaseFSM,
			RequestTopic:         leaseRequestTopic,
			Logger:               loggo.GetLogger("juju.worker.lease.raft"),
			PrometheusRegisterer: config.PrometheusRegisterer,
			NewWorker:            leasemanager.NewWorker,
			NewStore:             leasemanager.NewStore,
		})),

		validCredentialFlagName: credentialvalidator.Manifold(credentialvalidator.ManifoldConfig{
			APICallerName: apiCallerName,
			NewFacade:     credentialvalidator.NewFacade,
			NewWorker:     credentialvalidator.NewWorker,
		}),

		legacyLeasesFlagName: ifController(featureflag.Manifold(featureflag.ManifoldConfig{
			StateName: stateName,
			FlagName:  feature.LegacyLeases,
			Logger:    loggo.GetLogger("juju.worker.legacyleasesenabled"),
			NewWorker: featureflag.NewWorker,
		})),

		certificateUpdaterName: ifFullyUpgraded(certupdater.Manifold(certupdater.ManifoldConfig{
			AgentName:                agentName,
			StateName:                stateName,
			NewWorker:                certupdater.NewCertificateUpdater,
			NewMachineAddressWatcher: certupdater.NewMachineAddressWatcher,
		})),

		fanConfigurerName: ifNotMigrating(fanconfigurer.Manifold(fanconfigurer.ManifoldConfig{
			APICallerName: apiCallerName,
			Clock:         config.Clock,
		})),
	}

	return manifolds
}
