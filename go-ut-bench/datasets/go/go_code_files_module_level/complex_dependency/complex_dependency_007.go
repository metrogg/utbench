func registerCommands(r commandRegistry, ctx *cmd.Context) {
	// Creation commands.
	r.Register(newBootstrapCommand())
	r.Register(application.NewAddRelationCommand())

	// Cross model relations commands.
	r.Register(crossmodel.NewOfferCommand())
	r.Register(crossmodel.NewRemoveOfferCommand())
	r.Register(crossmodel.NewShowOfferedEndpointCommand())
	r.Register(crossmodel.NewListEndpointsCommand())
	r.Register(crossmodel.NewFindEndpointsCommand())
	r.Register(application.NewConsumeCommand())
	r.Register(application.NewSuspendRelationCommand())
	r.Register(application.NewResumeRelationCommand())

	// Firewall rule commands.
	r.Register(firewall.NewSetFirewallRuleCommand())
	r.Register(firewall.NewListFirewallRulesCommand())

	// Destruction commands.
	r.Register(application.NewRemoveRelationCommand())
	r.Register(application.NewRemoveApplicationCommand())
	r.Register(application.NewRemoveUnitCommand())
	r.Register(application.NewRemoveSaasCommand())

	// Reporting commands.
	r.Register(status.NewStatusCommand())
	r.Register(newSwitchCommand())
	r.Register(status.NewStatusHistoryCommand())

	// Error resolution and debugging commands.
	r.Register(newDefaultRunCommand(nil))
	r.Register(newSCPCommand(nil))
	r.Register(newSSHCommand(nil, nil))
	r.Register(application.NewResolvedCommand())
	r.Register(newDebugLogCommand(nil))
	r.Register(newDebugHooksCommand(nil))

	// Configuration commands.
	r.Register(model.NewModelGetConstraintsCommand())
	r.Register(model.NewModelSetConstraintsCommand())
	r.Register(newSyncToolsCommand())
	r.Register(newUpgradeJujuCommand())
	r.Register(newUpgradeControllerCommand())
	r.Register(application.NewUpgradeCharmCommand())
	r.Register(application.NewSetSeriesCommand())

	// Charm tool commands.
	r.Register(newHelpToolCommand())
	// TODO (anastasiamac 2017-08-1) This needs to be removed in Juju 3.x
	// lp#1707836
	r.Register(charmcmd.NewSuperCommand())

	// Manage backups.
	r.Register(backups.NewCreateCommand())
	r.Register(backups.NewDownloadCommand())
	r.Register(backups.NewShowCommand())
	r.Register(backups.NewListCommand())
	r.Register(backups.NewRemoveCommand())
	r.Register(backups.NewRestoreCommand())
	r.Register(backups.NewUploadCommand())

	// Manage authorized ssh keys.
	r.Register(NewAddKeysCommand())
	r.Register(NewRemoveKeysCommand())
	r.Register(NewImportKeysCommand())
	r.Register(NewListKeysCommand())

	// Manage users and access
	r.Register(user.NewAddCommand())
	r.Register(user.NewChangePasswordCommand())
	r.Register(user.NewShowUserCommand())
	r.Register(user.NewListCommand())
	r.Register(user.NewEnableCommand())
	r.Register(user.NewDisableCommand())
	r.Register(user.NewLoginCommand())
	r.Register(user.NewLogoutCommand())
	r.Register(user.NewRemoveCommand())
	r.Register(user.NewWhoAmICommand())

	// Manage cached images
	r.Register(cachedimages.NewRemoveCommand())
	r.Register(cachedimages.NewListCommand())

	// Manage machines
	r.Register(machine.NewAddCommand())
	r.Register(machine.NewRemoveCommand())
	r.Register(machine.NewListMachinesCommand())
	r.Register(machine.NewShowMachineCommand())
	r.Register(machine.NewUpgradeSeriesCommand())

	// Manage model
	r.Register(model.NewConfigCommand())
	r.Register(model.NewDefaultsCommand())
	r.Register(model.NewRetryProvisioningCommand())
	r.Register(model.NewDestroyCommand())
	r.Register(model.NewGrantCommand())
	r.Register(model.NewRevokeCommand())
	r.Register(model.NewShowCommand())
	r.Register(model.NewModelCredentialCommand())
	if featureflag.Enabled(feature.Generations) {
		r.Register(model.NewBranchCommand())
		r.Register(model.NewCommitCommand())
		r.Register(model.NewTrackBranchCommand())
		r.Register(model.NewCheckoutCommand())
		r.Register(model.NewDiffCommand())
	}

	r.Register(newMigrateCommand())
	r.Register(model.NewExportBundleCommand())

	if featureflag.Enabled(feature.DeveloperMode) {
		r.Register(model.NewDumpCommand())
		r.Register(model.NewDumpDBCommand())
	}

	// Manage and control actions
	r.Register(action.NewStatusCommand())
	r.Register(action.NewRunCommand())
	r.Register(action.NewShowOutputCommand())
	r.Register(action.NewListCommand())
	r.Register(action.NewCancelCommand())

	// Manage controller availability
	r.Register(newEnableHACommand())

	// Manage and control applications
	r.Register(application.NewAddUnitCommand())
	r.Register(application.NewConfigCommand())
	r.Register(application.NewDeployCommand())
	r.Register(application.NewExposeCommand())
	r.Register(application.NewUnexposeCommand())
	r.Register(application.NewApplicationGetConstraintsCommand())
	r.Register(application.NewApplicationSetConstraintsCommand())
	r.Register(application.NewBundleDiffCommand())
	r.Register(application.NewShowApplicationCommand())

	// Operation protection commands
	r.Register(block.NewDisableCommand())
	r.Register(block.NewListCommand())
	r.Register(block.NewEnableCommand())

	// Manage storage
	r.Register(storage.NewAddCommand())
	r.Register(storage.NewListCommand())
	r.Register(storage.NewPoolCreateCommand())
	r.Register(storage.NewPoolListCommand())
	r.Register(storage.NewPoolRemoveCommand())
	r.Register(storage.NewPoolUpdateCommand())
	r.Register(storage.NewShowCommand())
	r.Register(storage.NewRemoveStorageCommandWithAPI())
	r.Register(storage.NewDetachStorageCommandWithAPI())
	r.Register(storage.NewAttachStorageCommandWithAPI())
	r.Register(storage.NewImportFilesystemCommand(storage.NewStorageImporter, nil))

	// Manage spaces
	r.Register(space.NewAddCommand())
	r.Register(space.NewListCommand())
	r.Register(space.NewReloadCommand())
	if featureflag.Enabled(feature.PostNetCLIMVP) {
		r.Register(space.NewRemoveCommand())
		r.Register(space.NewUpdateCommand())
		r.Register(space.NewRenameCommand())
	}

	// Manage subnets
	r.Register(subnet.NewAddCommand())
	r.Register(subnet.NewListCommand())
	if featureflag.Enabled(feature.PostNetCLIMVP) {
		r.Register(subnet.NewCreateCommand())
		r.Register(subnet.NewRemoveCommand())
	}

	// Manage controllers
	r.Register(controller.NewAddModelCommand())
	r.Register(controller.NewDestroyCommand())
	r.Register(controller.NewListModelsCommand())
	r.Register(controller.NewKillCommand())
	r.Register(controller.NewListControllersCommand())
	r.Register(controller.NewRegisterCommand())
	r.Register(controller.NewUnregisterCommand(jujuclient.NewFileClientStore()))
	r.Register(controller.NewEnableDestroyControllerCommand())
	r.Register(controller.NewShowControllerCommand())
	r.Register(controller.NewConfigCommand())

	// Debug Metrics
	r.Register(metricsdebug.New())
	r.Register(metricsdebug.NewCollectMetricsCommand())
	r.Register(setmeterstatus.New())

	// Manage clouds and credentials
	r.Register(cloud.NewUpdateCloudCommand(&cloudToCommandAdapter{}))
	r.Register(cloud.NewUpdatePublicCloudsCommand())
	r.Register(cloud.NewListCloudsCommand())
	r.Register(cloud.NewListRegionsCommand())
	r.Register(cloud.NewShowCloudCommand())
	r.Register(cloud.NewAddCloudCommand(&cloudToCommandAdapter{}))
	r.Register(cloud.NewRemoveCloudCommand())
	r.Register(cloud.NewListCredentialsCommand())
	r.Register(cloud.NewDetectCredentialsCommand())
	r.Register(cloud.NewSetDefaultRegionCommand())
	r.Register(cloud.NewSetDefaultCredentialCommand())
	r.Register(cloud.NewAddCredentialCommand())
	r.Register(cloud.NewRemoveCredentialCommand())
	r.Register(cloud.NewUpdateCredentialCommand())
	r.Register(cloud.NewShowCredentialCommand())
	r.Register(model.NewGrantCloudCommand())
	r.Register(model.NewRevokeCloudCommand())

	// CAAS commands
	r.Register(caas.NewAddCAASCommand(&cloudToCommandAdapter{}))
	r.Register(caas.NewRemoveCAASCommand(&cloudToCommandAdapter{}))
	r.Register(application.NewScaleApplicationCommand())

	// Manage Application Credential Access
	r.Register(application.NewTrustCommand())

	// Juju GUI commands.
	r.Register(gui.NewGUICommand())
	r.Register(gui.NewUpgradeGUICommand())

	// Resource commands
	r.Register(resource.NewUploadCommand(resource.UploadDeps{
		NewClient: func(c *resource.UploadCommand) (resource.UploadClient, error) {
			apiRoot, err := c.NewAPIRoot()
			if err != nil {
				return nil, errors.Trace(err)
			}
			return resourceadapters.NewAPIClient(apiRoot)
		},
		OpenResource: func(s string) (resource.ReadSeekCloser, error) {
			return os.Open(s)
		},
	}))
	r.Register(resource.NewListCommand(resource.ListDeps{
		NewClient: func(c *resource.ListCommand) (resource.ListClient, error) {
			apiRoot, err := c.NewAPIRoot()
			if err != nil {
				return nil, errors.Trace(err)
			}
			return resourceadapters.NewAPIClient(apiRoot)
		},
	}))
	r.Register(resource.NewCharmResourcesCommand(nil))

	// Commands registered elsewhere.
	for _, newCommand := range registeredCommands {
		command := newCommand()
		r.Register(command)
	}
	for _, newCommand := range registeredEnvCommands {
		command := newCommand()
		r.Register(modelcmd.Wrap(command))
	}
	rcmd.RegisterAll(r)
}
