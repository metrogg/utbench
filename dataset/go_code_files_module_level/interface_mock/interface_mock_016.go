func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.CustomDeploymentStrategyParams)(nil), (*apps.CustomDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CustomDeploymentStrategyParams_To_apps_CustomDeploymentStrategyParams(a.(*v1.CustomDeploymentStrategyParams), b.(*apps.CustomDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.CustomDeploymentStrategyParams)(nil), (*v1.CustomDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_CustomDeploymentStrategyParams_To_v1_CustomDeploymentStrategyParams(a.(*apps.CustomDeploymentStrategyParams), b.(*v1.CustomDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentCause)(nil), (*apps.DeploymentCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentCause_To_apps_DeploymentCause(a.(*v1.DeploymentCause), b.(*apps.DeploymentCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentCause)(nil), (*v1.DeploymentCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentCause_To_v1_DeploymentCause(a.(*apps.DeploymentCause), b.(*v1.DeploymentCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentCauseImageTrigger)(nil), (*apps.DeploymentCauseImageTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentCauseImageTrigger_To_apps_DeploymentCauseImageTrigger(a.(*v1.DeploymentCauseImageTrigger), b.(*apps.DeploymentCauseImageTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentCauseImageTrigger)(nil), (*v1.DeploymentCauseImageTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentCauseImageTrigger_To_v1_DeploymentCauseImageTrigger(a.(*apps.DeploymentCauseImageTrigger), b.(*v1.DeploymentCauseImageTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentCondition)(nil), (*apps.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentCondition_To_apps_DeploymentCondition(a.(*v1.DeploymentCondition), b.(*apps.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentCondition)(nil), (*v1.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentCondition_To_v1_DeploymentCondition(a.(*apps.DeploymentCondition), b.(*v1.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfig)(nil), (*apps.DeploymentConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfig_To_apps_DeploymentConfig(a.(*v1.DeploymentConfig), b.(*apps.DeploymentConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfig)(nil), (*v1.DeploymentConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfig_To_v1_DeploymentConfig(a.(*apps.DeploymentConfig), b.(*v1.DeploymentConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfigList)(nil), (*apps.DeploymentConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfigList_To_apps_DeploymentConfigList(a.(*v1.DeploymentConfigList), b.(*apps.DeploymentConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfigList)(nil), (*v1.DeploymentConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfigList_To_v1_DeploymentConfigList(a.(*apps.DeploymentConfigList), b.(*v1.DeploymentConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfigRollback)(nil), (*apps.DeploymentConfigRollback)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfigRollback_To_apps_DeploymentConfigRollback(a.(*v1.DeploymentConfigRollback), b.(*apps.DeploymentConfigRollback), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfigRollback)(nil), (*v1.DeploymentConfigRollback)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfigRollback_To_v1_DeploymentConfigRollback(a.(*apps.DeploymentConfigRollback), b.(*v1.DeploymentConfigRollback), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfigRollbackSpec)(nil), (*apps.DeploymentConfigRollbackSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfigRollbackSpec_To_apps_DeploymentConfigRollbackSpec(a.(*v1.DeploymentConfigRollbackSpec), b.(*apps.DeploymentConfigRollbackSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfigRollbackSpec)(nil), (*v1.DeploymentConfigRollbackSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfigRollbackSpec_To_v1_DeploymentConfigRollbackSpec(a.(*apps.DeploymentConfigRollbackSpec), b.(*v1.DeploymentConfigRollbackSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfigSpec)(nil), (*apps.DeploymentConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfigSpec_To_apps_DeploymentConfigSpec(a.(*v1.DeploymentConfigSpec), b.(*apps.DeploymentConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfigSpec)(nil), (*v1.DeploymentConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfigSpec_To_v1_DeploymentConfigSpec(a.(*apps.DeploymentConfigSpec), b.(*v1.DeploymentConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentConfigStatus)(nil), (*apps.DeploymentConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentConfigStatus_To_apps_DeploymentConfigStatus(a.(*v1.DeploymentConfigStatus), b.(*apps.DeploymentConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentConfigStatus)(nil), (*v1.DeploymentConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentConfigStatus_To_v1_DeploymentConfigStatus(a.(*apps.DeploymentConfigStatus), b.(*v1.DeploymentConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentDetails)(nil), (*apps.DeploymentDetails)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentDetails_To_apps_DeploymentDetails(a.(*v1.DeploymentDetails), b.(*apps.DeploymentDetails), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentDetails)(nil), (*v1.DeploymentDetails)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentDetails_To_v1_DeploymentDetails(a.(*apps.DeploymentDetails), b.(*v1.DeploymentDetails), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentLog)(nil), (*apps.DeploymentLog)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentLog_To_apps_DeploymentLog(a.(*v1.DeploymentLog), b.(*apps.DeploymentLog), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentLog)(nil), (*v1.DeploymentLog)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentLog_To_v1_DeploymentLog(a.(*apps.DeploymentLog), b.(*v1.DeploymentLog), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentLogOptions)(nil), (*apps.DeploymentLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentLogOptions_To_apps_DeploymentLogOptions(a.(*v1.DeploymentLogOptions), b.(*apps.DeploymentLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentLogOptions)(nil), (*v1.DeploymentLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentLogOptions_To_v1_DeploymentLogOptions(a.(*apps.DeploymentLogOptions), b.(*v1.DeploymentLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentRequest)(nil), (*apps.DeploymentRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentRequest_To_apps_DeploymentRequest(a.(*v1.DeploymentRequest), b.(*apps.DeploymentRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentRequest)(nil), (*v1.DeploymentRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentRequest_To_v1_DeploymentRequest(a.(*apps.DeploymentRequest), b.(*v1.DeploymentRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStrategy)(nil), (*v1.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentTriggerImageChangeParams)(nil), (*apps.DeploymentTriggerImageChangeParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentTriggerImageChangeParams_To_apps_DeploymentTriggerImageChangeParams(a.(*v1.DeploymentTriggerImageChangeParams), b.(*apps.DeploymentTriggerImageChangeParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentTriggerImageChangeParams)(nil), (*v1.DeploymentTriggerImageChangeParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentTriggerImageChangeParams_To_v1_DeploymentTriggerImageChangeParams(a.(*apps.DeploymentTriggerImageChangeParams), b.(*v1.DeploymentTriggerImageChangeParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentTriggerPolicy)(nil), (*apps.DeploymentTriggerPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentTriggerPolicy_To_apps_DeploymentTriggerPolicy(a.(*v1.DeploymentTriggerPolicy), b.(*apps.DeploymentTriggerPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentTriggerPolicy)(nil), (*v1.DeploymentTriggerPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentTriggerPolicy_To_v1_DeploymentTriggerPolicy(a.(*apps.DeploymentTriggerPolicy), b.(*v1.DeploymentTriggerPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ExecNewPodHook)(nil), (*apps.ExecNewPodHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExecNewPodHook_To_apps_ExecNewPodHook(a.(*v1.ExecNewPodHook), b.(*apps.ExecNewPodHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ExecNewPodHook)(nil), (*v1.ExecNewPodHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ExecNewPodHook_To_v1_ExecNewPodHook(a.(*apps.ExecNewPodHook), b.(*v1.ExecNewPodHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LifecycleHook)(nil), (*apps.LifecycleHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LifecycleHook_To_apps_LifecycleHook(a.(*v1.LifecycleHook), b.(*apps.LifecycleHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.LifecycleHook)(nil), (*v1.LifecycleHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_LifecycleHook_To_v1_LifecycleHook(a.(*apps.LifecycleHook), b.(*v1.LifecycleHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RecreateDeploymentStrategyParams)(nil), (*apps.RecreateDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RecreateDeploymentStrategyParams_To_apps_RecreateDeploymentStrategyParams(a.(*v1.RecreateDeploymentStrategyParams), b.(*apps.RecreateDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RecreateDeploymentStrategyParams)(nil), (*v1.RecreateDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RecreateDeploymentStrategyParams_To_v1_RecreateDeploymentStrategyParams(a.(*apps.RecreateDeploymentStrategyParams), b.(*v1.RecreateDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RollingDeploymentStrategyParams)(nil), (*apps.RollingDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingDeploymentStrategyParams_To_apps_RollingDeploymentStrategyParams(a.(*v1.RollingDeploymentStrategyParams), b.(*apps.RollingDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingDeploymentStrategyParams)(nil), (*v1.RollingDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingDeploymentStrategyParams_To_v1_RollingDeploymentStrategyParams(a.(*apps.RollingDeploymentStrategyParams), b.(*v1.RollingDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TagImageHook)(nil), (*apps.TagImageHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TagImageHook_To_apps_TagImageHook(a.(*v1.TagImageHook), b.(*apps.TagImageHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.TagImageHook)(nil), (*v1.TagImageHook)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_TagImageHook_To_v1_TagImageHook(a.(*apps.TagImageHook), b.(*v1.TagImageHook), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentTriggerImageChangeParams)(nil), (*v1.DeploymentTriggerImageChangeParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentTriggerImageChangeParams_To_v1_DeploymentTriggerImageChangeParams(a.(*apps.DeploymentTriggerImageChangeParams), b.(*v1.DeploymentTriggerImageChangeParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingDeploymentStrategyParams)(nil), (*v1.RollingDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingDeploymentStrategyParams_To_v1_RollingDeploymentStrategyParams(a.(*apps.RollingDeploymentStrategyParams), b.(*v1.RollingDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DeploymentTriggerImageChangeParams)(nil), (*apps.DeploymentTriggerImageChangeParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentTriggerImageChangeParams_To_apps_DeploymentTriggerImageChangeParams(a.(*v1.DeploymentTriggerImageChangeParams), b.(*apps.DeploymentTriggerImageChangeParams), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.RollingDeploymentStrategyParams)(nil), (*apps.RollingDeploymentStrategyParams)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingDeploymentStrategyParams_To_apps_RollingDeploymentStrategyParams(a.(*v1.RollingDeploymentStrategyParams), b.(*apps.RollingDeploymentStrategyParams), scope)
	}); err != nil {
		return err
	}
	return nil
}
