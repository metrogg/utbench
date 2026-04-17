func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.BinaryBuildRequestOptions)(nil), (*build.BinaryBuildRequestOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BinaryBuildRequestOptions_To_build_BinaryBuildRequestOptions(a.(*v1.BinaryBuildRequestOptions), b.(*build.BinaryBuildRequestOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BinaryBuildRequestOptions)(nil), (*v1.BinaryBuildRequestOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BinaryBuildRequestOptions_To_v1_BinaryBuildRequestOptions(a.(*build.BinaryBuildRequestOptions), b.(*v1.BinaryBuildRequestOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BinaryBuildSource)(nil), (*build.BinaryBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BinaryBuildSource_To_build_BinaryBuildSource(a.(*v1.BinaryBuildSource), b.(*build.BinaryBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BinaryBuildSource)(nil), (*v1.BinaryBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BinaryBuildSource_To_v1_BinaryBuildSource(a.(*build.BinaryBuildSource), b.(*v1.BinaryBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BitbucketWebHookCause)(nil), (*build.BitbucketWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BitbucketWebHookCause_To_build_BitbucketWebHookCause(a.(*v1.BitbucketWebHookCause), b.(*build.BitbucketWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BitbucketWebHookCause)(nil), (*v1.BitbucketWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BitbucketWebHookCause_To_v1_BitbucketWebHookCause(a.(*build.BitbucketWebHookCause), b.(*v1.BitbucketWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Build)(nil), (*build.Build)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Build_To_build_Build(a.(*v1.Build), b.(*build.Build), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.Build)(nil), (*v1.Build)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_Build_To_v1_Build(a.(*build.Build), b.(*v1.Build), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildConfig)(nil), (*build.BuildConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildConfig_To_build_BuildConfig(a.(*v1.BuildConfig), b.(*build.BuildConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildConfig)(nil), (*v1.BuildConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildConfig_To_v1_BuildConfig(a.(*build.BuildConfig), b.(*v1.BuildConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildConfigList)(nil), (*build.BuildConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildConfigList_To_build_BuildConfigList(a.(*v1.BuildConfigList), b.(*build.BuildConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildConfigList)(nil), (*v1.BuildConfigList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildConfigList_To_v1_BuildConfigList(a.(*build.BuildConfigList), b.(*v1.BuildConfigList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildConfigSpec)(nil), (*build.BuildConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildConfigSpec_To_build_BuildConfigSpec(a.(*v1.BuildConfigSpec), b.(*build.BuildConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildConfigSpec)(nil), (*v1.BuildConfigSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildConfigSpec_To_v1_BuildConfigSpec(a.(*build.BuildConfigSpec), b.(*v1.BuildConfigSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildConfigStatus)(nil), (*build.BuildConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildConfigStatus_To_build_BuildConfigStatus(a.(*v1.BuildConfigStatus), b.(*build.BuildConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildConfigStatus)(nil), (*v1.BuildConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildConfigStatus_To_v1_BuildConfigStatus(a.(*build.BuildConfigStatus), b.(*v1.BuildConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildList)(nil), (*build.BuildList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildList_To_build_BuildList(a.(*v1.BuildList), b.(*build.BuildList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildList)(nil), (*v1.BuildList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildList_To_v1_BuildList(a.(*build.BuildList), b.(*v1.BuildList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildLog)(nil), (*build.BuildLog)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildLog_To_build_BuildLog(a.(*v1.BuildLog), b.(*build.BuildLog), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildLog)(nil), (*v1.BuildLog)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildLog_To_v1_BuildLog(a.(*build.BuildLog), b.(*v1.BuildLog), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildLogOptions)(nil), (*build.BuildLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildLogOptions_To_build_BuildLogOptions(a.(*v1.BuildLogOptions), b.(*build.BuildLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildLogOptions)(nil), (*v1.BuildLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildLogOptions_To_v1_BuildLogOptions(a.(*build.BuildLogOptions), b.(*v1.BuildLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildOutput)(nil), (*build.BuildOutput)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildOutput_To_build_BuildOutput(a.(*v1.BuildOutput), b.(*build.BuildOutput), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildOutput)(nil), (*v1.BuildOutput)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildOutput_To_v1_BuildOutput(a.(*build.BuildOutput), b.(*v1.BuildOutput), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildPostCommitSpec)(nil), (*build.BuildPostCommitSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildPostCommitSpec_To_build_BuildPostCommitSpec(a.(*v1.BuildPostCommitSpec), b.(*build.BuildPostCommitSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildPostCommitSpec)(nil), (*v1.BuildPostCommitSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildPostCommitSpec_To_v1_BuildPostCommitSpec(a.(*build.BuildPostCommitSpec), b.(*v1.BuildPostCommitSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildRequest)(nil), (*build.BuildRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildRequest_To_build_BuildRequest(a.(*v1.BuildRequest), b.(*build.BuildRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildRequest)(nil), (*v1.BuildRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildRequest_To_v1_BuildRequest(a.(*build.BuildRequest), b.(*v1.BuildRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildSource)(nil), (*build.BuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildSource_To_build_BuildSource(a.(*v1.BuildSource), b.(*build.BuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildSource)(nil), (*v1.BuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildSource_To_v1_BuildSource(a.(*build.BuildSource), b.(*v1.BuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildSpec)(nil), (*build.BuildSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildSpec_To_build_BuildSpec(a.(*v1.BuildSpec), b.(*build.BuildSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildSpec)(nil), (*v1.BuildSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildSpec_To_v1_BuildSpec(a.(*build.BuildSpec), b.(*v1.BuildSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildStatus)(nil), (*build.BuildStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildStatus_To_build_BuildStatus(a.(*v1.BuildStatus), b.(*build.BuildStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildStatus)(nil), (*v1.BuildStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildStatus_To_v1_BuildStatus(a.(*build.BuildStatus), b.(*v1.BuildStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildStatusOutput)(nil), (*build.BuildStatusOutput)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildStatusOutput_To_build_BuildStatusOutput(a.(*v1.BuildStatusOutput), b.(*build.BuildStatusOutput), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildStatusOutput)(nil), (*v1.BuildStatusOutput)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildStatusOutput_To_v1_BuildStatusOutput(a.(*build.BuildStatusOutput), b.(*v1.BuildStatusOutput), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildStatusOutputTo)(nil), (*build.BuildStatusOutputTo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildStatusOutputTo_To_build_BuildStatusOutputTo(a.(*v1.BuildStatusOutputTo), b.(*build.BuildStatusOutputTo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildStatusOutputTo)(nil), (*v1.BuildStatusOutputTo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildStatusOutputTo_To_v1_BuildStatusOutputTo(a.(*build.BuildStatusOutputTo), b.(*v1.BuildStatusOutputTo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildStrategy)(nil), (*build.BuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildStrategy_To_build_BuildStrategy(a.(*v1.BuildStrategy), b.(*build.BuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildStrategy)(nil), (*v1.BuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildStrategy_To_v1_BuildStrategy(a.(*build.BuildStrategy), b.(*v1.BuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildTriggerCause)(nil), (*build.BuildTriggerCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildTriggerCause_To_build_BuildTriggerCause(a.(*v1.BuildTriggerCause), b.(*build.BuildTriggerCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildTriggerCause)(nil), (*v1.BuildTriggerCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildTriggerCause_To_v1_BuildTriggerCause(a.(*build.BuildTriggerCause), b.(*v1.BuildTriggerCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BuildTriggerPolicy)(nil), (*build.BuildTriggerPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildTriggerPolicy_To_build_BuildTriggerPolicy(a.(*v1.BuildTriggerPolicy), b.(*build.BuildTriggerPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.BuildTriggerPolicy)(nil), (*v1.BuildTriggerPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildTriggerPolicy_To_v1_BuildTriggerPolicy(a.(*build.BuildTriggerPolicy), b.(*v1.BuildTriggerPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CommonSpec)(nil), (*build.CommonSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CommonSpec_To_build_CommonSpec(a.(*v1.CommonSpec), b.(*build.CommonSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.CommonSpec)(nil), (*v1.CommonSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_CommonSpec_To_v1_CommonSpec(a.(*build.CommonSpec), b.(*v1.CommonSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CommonWebHookCause)(nil), (*build.CommonWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CommonWebHookCause_To_build_CommonWebHookCause(a.(*v1.CommonWebHookCause), b.(*build.CommonWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.CommonWebHookCause)(nil), (*v1.CommonWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_CommonWebHookCause_To_v1_CommonWebHookCause(a.(*build.CommonWebHookCause), b.(*v1.CommonWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapBuildSource)(nil), (*build.ConfigMapBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapBuildSource_To_build_ConfigMapBuildSource(a.(*v1.ConfigMapBuildSource), b.(*build.ConfigMapBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ConfigMapBuildSource)(nil), (*v1.ConfigMapBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ConfigMapBuildSource_To_v1_ConfigMapBuildSource(a.(*build.ConfigMapBuildSource), b.(*v1.ConfigMapBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CustomBuildStrategy)(nil), (*build.CustomBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CustomBuildStrategy_To_build_CustomBuildStrategy(a.(*v1.CustomBuildStrategy), b.(*build.CustomBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.CustomBuildStrategy)(nil), (*v1.CustomBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_CustomBuildStrategy_To_v1_CustomBuildStrategy(a.(*build.CustomBuildStrategy), b.(*v1.CustomBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DockerBuildStrategy)(nil), (*build.DockerBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DockerBuildStrategy_To_build_DockerBuildStrategy(a.(*v1.DockerBuildStrategy), b.(*build.DockerBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.DockerBuildStrategy)(nil), (*v1.DockerBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_DockerBuildStrategy_To_v1_DockerBuildStrategy(a.(*build.DockerBuildStrategy), b.(*v1.DockerBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DockerStrategyOptions)(nil), (*build.DockerStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DockerStrategyOptions_To_build_DockerStrategyOptions(a.(*v1.DockerStrategyOptions), b.(*build.DockerStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.DockerStrategyOptions)(nil), (*v1.DockerStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_DockerStrategyOptions_To_v1_DockerStrategyOptions(a.(*build.DockerStrategyOptions), b.(*v1.DockerStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GenericWebHookCause)(nil), (*build.GenericWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GenericWebHookCause_To_build_GenericWebHookCause(a.(*v1.GenericWebHookCause), b.(*build.GenericWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GenericWebHookCause)(nil), (*v1.GenericWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GenericWebHookCause_To_v1_GenericWebHookCause(a.(*build.GenericWebHookCause), b.(*v1.GenericWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GenericWebHookEvent)(nil), (*build.GenericWebHookEvent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GenericWebHookEvent_To_build_GenericWebHookEvent(a.(*v1.GenericWebHookEvent), b.(*build.GenericWebHookEvent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GenericWebHookEvent)(nil), (*v1.GenericWebHookEvent)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GenericWebHookEvent_To_v1_GenericWebHookEvent(a.(*build.GenericWebHookEvent), b.(*v1.GenericWebHookEvent), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitBuildSource)(nil), (*build.GitBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitBuildSource_To_build_GitBuildSource(a.(*v1.GitBuildSource), b.(*build.GitBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitBuildSource)(nil), (*v1.GitBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitBuildSource_To_v1_GitBuildSource(a.(*build.GitBuildSource), b.(*v1.GitBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitHubWebHookCause)(nil), (*build.GitHubWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitHubWebHookCause_To_build_GitHubWebHookCause(a.(*v1.GitHubWebHookCause), b.(*build.GitHubWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitHubWebHookCause)(nil), (*v1.GitHubWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitHubWebHookCause_To_v1_GitHubWebHookCause(a.(*build.GitHubWebHookCause), b.(*v1.GitHubWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitInfo)(nil), (*build.GitInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitInfo_To_build_GitInfo(a.(*v1.GitInfo), b.(*build.GitInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitInfo)(nil), (*v1.GitInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitInfo_To_v1_GitInfo(a.(*build.GitInfo), b.(*v1.GitInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitLabWebHookCause)(nil), (*build.GitLabWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitLabWebHookCause_To_build_GitLabWebHookCause(a.(*v1.GitLabWebHookCause), b.(*build.GitLabWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitLabWebHookCause)(nil), (*v1.GitLabWebHookCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitLabWebHookCause_To_v1_GitLabWebHookCause(a.(*build.GitLabWebHookCause), b.(*v1.GitLabWebHookCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitRefInfo)(nil), (*build.GitRefInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitRefInfo_To_build_GitRefInfo(a.(*v1.GitRefInfo), b.(*build.GitRefInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitRefInfo)(nil), (*v1.GitRefInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitRefInfo_To_v1_GitRefInfo(a.(*build.GitRefInfo), b.(*v1.GitRefInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitSourceRevision)(nil), (*build.GitSourceRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitSourceRevision_To_build_GitSourceRevision(a.(*v1.GitSourceRevision), b.(*build.GitSourceRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.GitSourceRevision)(nil), (*v1.GitSourceRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_GitSourceRevision_To_v1_GitSourceRevision(a.(*build.GitSourceRevision), b.(*v1.GitSourceRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageChangeCause)(nil), (*build.ImageChangeCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageChangeCause_To_build_ImageChangeCause(a.(*v1.ImageChangeCause), b.(*build.ImageChangeCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ImageChangeCause)(nil), (*v1.ImageChangeCause)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ImageChangeCause_To_v1_ImageChangeCause(a.(*build.ImageChangeCause), b.(*v1.ImageChangeCause), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageChangeTrigger)(nil), (*build.ImageChangeTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageChangeTrigger_To_build_ImageChangeTrigger(a.(*v1.ImageChangeTrigger), b.(*build.ImageChangeTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ImageChangeTrigger)(nil), (*v1.ImageChangeTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ImageChangeTrigger_To_v1_ImageChangeTrigger(a.(*build.ImageChangeTrigger), b.(*v1.ImageChangeTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageLabel)(nil), (*build.ImageLabel)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageLabel_To_build_ImageLabel(a.(*v1.ImageLabel), b.(*build.ImageLabel), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ImageLabel)(nil), (*v1.ImageLabel)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ImageLabel_To_v1_ImageLabel(a.(*build.ImageLabel), b.(*v1.ImageLabel), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageSource)(nil), (*build.ImageSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageSource_To_build_ImageSource(a.(*v1.ImageSource), b.(*build.ImageSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ImageSource)(nil), (*v1.ImageSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ImageSource_To_v1_ImageSource(a.(*build.ImageSource), b.(*v1.ImageSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ImageSourcePath)(nil), (*build.ImageSourcePath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ImageSourcePath_To_build_ImageSourcePath(a.(*v1.ImageSourcePath), b.(*build.ImageSourcePath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ImageSourcePath)(nil), (*v1.ImageSourcePath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ImageSourcePath_To_v1_ImageSourcePath(a.(*build.ImageSourcePath), b.(*v1.ImageSourcePath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.JenkinsPipelineBuildStrategy)(nil), (*build.JenkinsPipelineBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_JenkinsPipelineBuildStrategy_To_build_JenkinsPipelineBuildStrategy(a.(*v1.JenkinsPipelineBuildStrategy), b.(*build.JenkinsPipelineBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.JenkinsPipelineBuildStrategy)(nil), (*v1.JenkinsPipelineBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_JenkinsPipelineBuildStrategy_To_v1_JenkinsPipelineBuildStrategy(a.(*build.JenkinsPipelineBuildStrategy), b.(*v1.JenkinsPipelineBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ProxyConfig)(nil), (*build.ProxyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ProxyConfig_To_build_ProxyConfig(a.(*v1.ProxyConfig), b.(*build.ProxyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.ProxyConfig)(nil), (*v1.ProxyConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_ProxyConfig_To_v1_ProxyConfig(a.(*build.ProxyConfig), b.(*v1.ProxyConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretBuildSource)(nil), (*build.SecretBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretBuildSource_To_build_SecretBuildSource(a.(*v1.SecretBuildSource), b.(*build.SecretBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SecretBuildSource)(nil), (*v1.SecretBuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SecretBuildSource_To_v1_SecretBuildSource(a.(*build.SecretBuildSource), b.(*v1.SecretBuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretLocalReference)(nil), (*build.SecretLocalReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretLocalReference_To_build_SecretLocalReference(a.(*v1.SecretLocalReference), b.(*build.SecretLocalReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SecretLocalReference)(nil), (*v1.SecretLocalReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SecretLocalReference_To_v1_SecretLocalReference(a.(*build.SecretLocalReference), b.(*v1.SecretLocalReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretSpec)(nil), (*build.SecretSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretSpec_To_build_SecretSpec(a.(*v1.SecretSpec), b.(*build.SecretSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SecretSpec)(nil), (*v1.SecretSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SecretSpec_To_v1_SecretSpec(a.(*build.SecretSpec), b.(*v1.SecretSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SourceBuildStrategy)(nil), (*build.SourceBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceBuildStrategy_To_build_SourceBuildStrategy(a.(*v1.SourceBuildStrategy), b.(*build.SourceBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SourceBuildStrategy)(nil), (*v1.SourceBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SourceBuildStrategy_To_v1_SourceBuildStrategy(a.(*build.SourceBuildStrategy), b.(*v1.SourceBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SourceControlUser)(nil), (*build.SourceControlUser)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceControlUser_To_build_SourceControlUser(a.(*v1.SourceControlUser), b.(*build.SourceControlUser), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SourceControlUser)(nil), (*v1.SourceControlUser)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SourceControlUser_To_v1_SourceControlUser(a.(*build.SourceControlUser), b.(*v1.SourceControlUser), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SourceRevision)(nil), (*build.SourceRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceRevision_To_build_SourceRevision(a.(*v1.SourceRevision), b.(*build.SourceRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SourceRevision)(nil), (*v1.SourceRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SourceRevision_To_v1_SourceRevision(a.(*build.SourceRevision), b.(*v1.SourceRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SourceStrategyOptions)(nil), (*build.SourceStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceStrategyOptions_To_build_SourceStrategyOptions(a.(*v1.SourceStrategyOptions), b.(*build.SourceStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.SourceStrategyOptions)(nil), (*v1.SourceStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SourceStrategyOptions_To_v1_SourceStrategyOptions(a.(*build.SourceStrategyOptions), b.(*v1.SourceStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StageInfo)(nil), (*build.StageInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StageInfo_To_build_StageInfo(a.(*v1.StageInfo), b.(*build.StageInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.StageInfo)(nil), (*v1.StageInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_StageInfo_To_v1_StageInfo(a.(*build.StageInfo), b.(*v1.StageInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StepInfo)(nil), (*build.StepInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StepInfo_To_build_StepInfo(a.(*v1.StepInfo), b.(*build.StepInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.StepInfo)(nil), (*v1.StepInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_StepInfo_To_v1_StepInfo(a.(*build.StepInfo), b.(*v1.StepInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.WebHookTrigger)(nil), (*build.WebHookTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_WebHookTrigger_To_build_WebHookTrigger(a.(*v1.WebHookTrigger), b.(*build.WebHookTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*build.WebHookTrigger)(nil), (*v1.WebHookTrigger)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_WebHookTrigger_To_v1_WebHookTrigger(a.(*build.WebHookTrigger), b.(*v1.WebHookTrigger), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*build.BuildSource)(nil), (*v1.BuildSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildSource_To_v1_BuildSource(a.(*build.BuildSource), b.(*v1.BuildSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*build.BuildStrategy)(nil), (*v1.BuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_BuildStrategy_To_v1_BuildStrategy(a.(*build.BuildStrategy), b.(*v1.BuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*build.SourceRevision)(nil), (*v1.SourceRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_build_SourceRevision_To_v1_SourceRevision(a.(*build.SourceRevision), b.(*v1.SourceRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.BuildConfig)(nil), (*build.BuildConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildConfig_To_build_BuildConfig(a.(*v1.BuildConfig), b.(*build.BuildConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.BuildOutput)(nil), (*build.BuildOutput)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildOutput_To_build_BuildOutput(a.(*v1.BuildOutput), b.(*build.BuildOutput), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.BuildTriggerPolicy)(nil), (*build.BuildTriggerPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BuildTriggerPolicy_To_build_BuildTriggerPolicy(a.(*v1.BuildTriggerPolicy), b.(*build.BuildTriggerPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.CustomBuildStrategy)(nil), (*build.CustomBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CustomBuildStrategy_To_build_CustomBuildStrategy(a.(*v1.CustomBuildStrategy), b.(*build.CustomBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DockerBuildStrategy)(nil), (*build.DockerBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DockerBuildStrategy_To_build_DockerBuildStrategy(a.(*v1.DockerBuildStrategy), b.(*build.DockerBuildStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.SourceBuildStrategy)(nil), (*build.SourceBuildStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SourceBuildStrategy_To_build_SourceBuildStrategy(a.(*v1.SourceBuildStrategy), b.(*build.SourceBuildStrategy), scope)
	}); err != nil {
		return err
	}
	return nil
}
