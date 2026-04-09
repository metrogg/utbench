func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.ControllerRevision)(nil), (*apps.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision(a.(*v1beta1.ControllerRevision), b.(*apps.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevision)(nil), (*v1beta1.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevision_To_v1beta1_ControllerRevision(a.(*apps.ControllerRevision), b.(*v1beta1.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ControllerRevisionList)(nil), (*apps.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList(a.(*v1beta1.ControllerRevisionList), b.(*apps.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevisionList)(nil), (*v1beta1.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevisionList_To_v1beta1_ControllerRevisionList(a.(*apps.ControllerRevisionList), b.(*v1beta1.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Deployment)(nil), (*apps.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Deployment_To_apps_Deployment(a.(*v1beta1.Deployment), b.(*apps.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.Deployment)(nil), (*v1beta1.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_Deployment_To_v1beta1_Deployment(a.(*apps.Deployment), b.(*v1beta1.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentCondition)(nil), (*apps.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentCondition_To_apps_DeploymentCondition(a.(*v1beta1.DeploymentCondition), b.(*apps.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentCondition)(nil), (*v1beta1.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentCondition_To_v1beta1_DeploymentCondition(a.(*apps.DeploymentCondition), b.(*v1beta1.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentList)(nil), (*apps.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentList_To_apps_DeploymentList(a.(*v1beta1.DeploymentList), b.(*apps.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentList)(nil), (*v1beta1.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentList_To_v1beta1_DeploymentList(a.(*apps.DeploymentList), b.(*v1beta1.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentRollback)(nil), (*apps.DeploymentRollback)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentRollback_To_apps_DeploymentRollback(a.(*v1beta1.DeploymentRollback), b.(*apps.DeploymentRollback), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentRollback)(nil), (*v1beta1.DeploymentRollback)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentRollback_To_v1beta1_DeploymentRollback(a.(*apps.DeploymentRollback), b.(*v1beta1.DeploymentRollback), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1beta1.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentSpec)(nil), (*v1beta1.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1beta1_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1beta1.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentStatus)(nil), (*apps.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentStatus_To_apps_DeploymentStatus(a.(*v1beta1.DeploymentStatus), b.(*apps.DeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStatus)(nil), (*v1beta1.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStatus_To_v1beta1_DeploymentStatus(a.(*apps.DeploymentStatus), b.(*v1beta1.DeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1beta1.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStrategy)(nil), (*v1beta1.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1beta1_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1beta1.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RollbackConfig)(nil), (*apps.RollbackConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollbackConfig_To_apps_RollbackConfig(a.(*v1beta1.RollbackConfig), b.(*apps.RollbackConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollbackConfig)(nil), (*v1beta1.RollbackConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollbackConfig_To_v1beta1_RollbackConfig(a.(*apps.RollbackConfig), b.(*v1beta1.RollbackConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1beta1.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1beta1.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1beta1.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RollingUpdateStatefulSetStrategy)(nil), (*apps.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(a.(*v1beta1.RollingUpdateStatefulSetStrategy), b.(*apps.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateStatefulSetStrategy)(nil), (*v1beta1.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta1_RollingUpdateStatefulSetStrategy(a.(*apps.RollingUpdateStatefulSetStrategy), b.(*v1beta1.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Scale)(nil), (*autoscaling.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Scale_To_autoscaling_Scale(a.(*v1beta1.Scale), b.(*autoscaling.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.Scale)(nil), (*v1beta1.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_Scale_To_v1beta1_Scale(a.(*autoscaling.Scale), b.(*v1beta1.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ScaleSpec)(nil), (*autoscaling.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ScaleSpec_To_autoscaling_ScaleSpec(a.(*v1beta1.ScaleSpec), b.(*autoscaling.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleSpec)(nil), (*v1beta1.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleSpec_To_v1beta1_ScaleSpec(a.(*autoscaling.ScaleSpec), b.(*v1beta1.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ScaleStatus)(nil), (*autoscaling.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus(a.(*v1beta1.ScaleStatus), b.(*autoscaling.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1beta1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1beta1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSet)(nil), (*apps.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSet_To_apps_StatefulSet(a.(*v1beta1.StatefulSet), b.(*apps.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSet)(nil), (*v1beta1.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSet_To_v1beta1_StatefulSet(a.(*apps.StatefulSet), b.(*v1beta1.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSetCondition)(nil), (*apps.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetCondition_To_apps_StatefulSetCondition(a.(*v1beta1.StatefulSetCondition), b.(*apps.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetCondition)(nil), (*v1beta1.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetCondition_To_v1beta1_StatefulSetCondition(a.(*apps.StatefulSetCondition), b.(*v1beta1.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSetList)(nil), (*apps.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetList_To_apps_StatefulSetList(a.(*v1beta1.StatefulSetList), b.(*apps.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetList)(nil), (*v1beta1.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetList_To_v1beta1_StatefulSetList(a.(*apps.StatefulSetList), b.(*v1beta1.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1beta1.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetSpec)(nil), (*v1beta1.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1beta1_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1beta1.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSetStatus)(nil), (*apps.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetStatus_To_apps_StatefulSetStatus(a.(*v1beta1.StatefulSetStatus), b.(*apps.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetStatus)(nil), (*v1beta1.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetStatus_To_v1beta1_StatefulSetStatus(a.(*apps.StatefulSetStatus), b.(*v1beta1.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1beta1.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1beta1.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1beta1_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1beta1.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentSpec)(nil), (*v1beta1.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1beta1_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1beta1.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentStrategy)(nil), (*v1beta1.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1beta1_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1beta1.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1beta1.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1beta1.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetSpec)(nil), (*v1beta1.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1beta1_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1beta1.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1beta1.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1beta1_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1beta1.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1beta1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1beta1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1beta1.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1beta1.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1beta1.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.ScaleStatus)(nil), (*autoscaling.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus(a.(*v1beta1.ScaleStatus), b.(*autoscaling.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1beta1.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1beta1.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	return nil
}
