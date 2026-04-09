func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta2.ControllerRevision)(nil), (*apps.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ControllerRevision_To_apps_ControllerRevision(a.(*v1beta2.ControllerRevision), b.(*apps.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevision)(nil), (*v1beta2.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevision_To_v1beta2_ControllerRevision(a.(*apps.ControllerRevision), b.(*v1beta2.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ControllerRevisionList)(nil), (*apps.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ControllerRevisionList_To_apps_ControllerRevisionList(a.(*v1beta2.ControllerRevisionList), b.(*apps.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevisionList)(nil), (*v1beta2.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevisionList_To_v1beta2_ControllerRevisionList(a.(*apps.ControllerRevisionList), b.(*v1beta2.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSet)(nil), (*apps.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSet_To_apps_DaemonSet(a.(*v1beta2.DaemonSet), b.(*apps.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSet)(nil), (*v1beta2.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSet_To_v1beta2_DaemonSet(a.(*apps.DaemonSet), b.(*v1beta2.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSetCondition)(nil), (*apps.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetCondition_To_apps_DaemonSetCondition(a.(*v1beta2.DaemonSetCondition), b.(*apps.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetCondition)(nil), (*v1beta2.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetCondition_To_v1beta2_DaemonSetCondition(a.(*apps.DaemonSetCondition), b.(*v1beta2.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSetList)(nil), (*apps.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetList_To_apps_DaemonSetList(a.(*v1beta2.DaemonSetList), b.(*apps.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetList)(nil), (*v1beta2.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetList_To_v1beta2_DaemonSetList(a.(*apps.DaemonSetList), b.(*v1beta2.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSetSpec)(nil), (*apps.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetSpec_To_apps_DaemonSetSpec(a.(*v1beta2.DaemonSetSpec), b.(*apps.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetSpec)(nil), (*v1beta2.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetSpec_To_v1beta2_DaemonSetSpec(a.(*apps.DaemonSetSpec), b.(*v1beta2.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSetStatus)(nil), (*apps.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetStatus_To_apps_DaemonSetStatus(a.(*v1beta2.DaemonSetStatus), b.(*apps.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetStatus)(nil), (*v1beta2.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetStatus_To_v1beta2_DaemonSetStatus(a.(*apps.DaemonSetStatus), b.(*v1beta2.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DaemonSetUpdateStrategy)(nil), (*apps.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetUpdateStrategy_To_apps_DaemonSetUpdateStrategy(a.(*v1beta2.DaemonSetUpdateStrategy), b.(*apps.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetUpdateStrategy)(nil), (*v1beta2.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetUpdateStrategy_To_v1beta2_DaemonSetUpdateStrategy(a.(*apps.DaemonSetUpdateStrategy), b.(*v1beta2.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.Deployment)(nil), (*apps.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Deployment_To_apps_Deployment(a.(*v1beta2.Deployment), b.(*apps.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.Deployment)(nil), (*v1beta2.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_Deployment_To_v1beta2_Deployment(a.(*apps.Deployment), b.(*v1beta2.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DeploymentCondition)(nil), (*apps.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentCondition_To_apps_DeploymentCondition(a.(*v1beta2.DeploymentCondition), b.(*apps.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentCondition)(nil), (*v1beta2.DeploymentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentCondition_To_v1beta2_DeploymentCondition(a.(*apps.DeploymentCondition), b.(*v1beta2.DeploymentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DeploymentList)(nil), (*apps.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentList_To_apps_DeploymentList(a.(*v1beta2.DeploymentList), b.(*apps.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentList)(nil), (*v1beta2.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentList_To_v1beta2_DeploymentList(a.(*apps.DeploymentList), b.(*v1beta2.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1beta2.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentSpec)(nil), (*v1beta2.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1beta2_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1beta2.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DeploymentStatus)(nil), (*apps.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentStatus_To_apps_DeploymentStatus(a.(*v1beta2.DeploymentStatus), b.(*apps.DeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStatus)(nil), (*v1beta2.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStatus_To_v1beta2_DeploymentStatus(a.(*apps.DeploymentStatus), b.(*v1beta2.DeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1beta2.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStrategy)(nil), (*v1beta2.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1beta2_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1beta2.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ReplicaSet)(nil), (*apps.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSet_To_apps_ReplicaSet(a.(*v1beta2.ReplicaSet), b.(*apps.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSet)(nil), (*v1beta2.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSet_To_v1beta2_ReplicaSet(a.(*apps.ReplicaSet), b.(*v1beta2.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ReplicaSetCondition)(nil), (*apps.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSetCondition_To_apps_ReplicaSetCondition(a.(*v1beta2.ReplicaSetCondition), b.(*apps.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetCondition)(nil), (*v1beta2.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetCondition_To_v1beta2_ReplicaSetCondition(a.(*apps.ReplicaSetCondition), b.(*v1beta2.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ReplicaSetList)(nil), (*apps.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSetList_To_apps_ReplicaSetList(a.(*v1beta2.ReplicaSetList), b.(*apps.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetList)(nil), (*v1beta2.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetList_To_v1beta2_ReplicaSetList(a.(*apps.ReplicaSetList), b.(*v1beta2.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1beta2.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1beta2.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1beta2_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1beta2.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ReplicaSetStatus)(nil), (*apps.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSetStatus_To_apps_ReplicaSetStatus(a.(*v1beta2.ReplicaSetStatus), b.(*apps.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetStatus)(nil), (*v1beta2.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetStatus_To_v1beta2_ReplicaSetStatus(a.(*apps.ReplicaSetStatus), b.(*v1beta2.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1beta2.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1beta2.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1beta2_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1beta2.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1beta2.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1beta2.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1beta2_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1beta2.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.RollingUpdateStatefulSetStrategy)(nil), (*apps.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(a.(*v1beta2.RollingUpdateStatefulSetStrategy), b.(*apps.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateStatefulSetStrategy)(nil), (*v1beta2.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateStatefulSetStrategy_To_v1beta2_RollingUpdateStatefulSetStrategy(a.(*apps.RollingUpdateStatefulSetStrategy), b.(*v1beta2.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.Scale)(nil), (*autoscaling.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Scale_To_autoscaling_Scale(a.(*v1beta2.Scale), b.(*autoscaling.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.Scale)(nil), (*v1beta2.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_Scale_To_v1beta2_Scale(a.(*autoscaling.Scale), b.(*v1beta2.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ScaleSpec)(nil), (*autoscaling.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ScaleSpec_To_autoscaling_ScaleSpec(a.(*v1beta2.ScaleSpec), b.(*autoscaling.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleSpec)(nil), (*v1beta2.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleSpec_To_v1beta2_ScaleSpec(a.(*autoscaling.ScaleSpec), b.(*v1beta2.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.ScaleStatus)(nil), (*autoscaling.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ScaleStatus_To_autoscaling_ScaleStatus(a.(*v1beta2.ScaleStatus), b.(*autoscaling.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1beta2.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1beta2_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1beta2.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSet)(nil), (*apps.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSet_To_apps_StatefulSet(a.(*v1beta2.StatefulSet), b.(*apps.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSet)(nil), (*v1beta2.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSet_To_v1beta2_StatefulSet(a.(*apps.StatefulSet), b.(*v1beta2.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSetCondition)(nil), (*apps.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetCondition_To_apps_StatefulSetCondition(a.(*v1beta2.StatefulSetCondition), b.(*apps.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetCondition)(nil), (*v1beta2.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetCondition_To_v1beta2_StatefulSetCondition(a.(*apps.StatefulSetCondition), b.(*v1beta2.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSetList)(nil), (*apps.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetList_To_apps_StatefulSetList(a.(*v1beta2.StatefulSetList), b.(*apps.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetList)(nil), (*v1beta2.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetList_To_v1beta2_StatefulSetList(a.(*apps.StatefulSetList), b.(*v1beta2.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1beta2.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetSpec)(nil), (*v1beta2.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1beta2_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1beta2.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSetStatus)(nil), (*apps.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetStatus_To_apps_StatefulSetStatus(a.(*v1beta2.StatefulSetStatus), b.(*apps.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetStatus)(nil), (*v1beta2.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetStatus_To_v1beta2_StatefulSetStatus(a.(*apps.StatefulSetStatus), b.(*v1beta2.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1beta2.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1beta2.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1beta2_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1beta2.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSetSpec)(nil), (*v1beta2.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetSpec_To_v1beta2_DaemonSetSpec(a.(*apps.DaemonSetSpec), b.(*v1beta2.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSetUpdateStrategy)(nil), (*v1beta2.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetUpdateStrategy_To_v1beta2_DaemonSetUpdateStrategy(a.(*apps.DaemonSetUpdateStrategy), b.(*v1beta2.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSet)(nil), (*v1beta2.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSet_To_v1beta2_DaemonSet(a.(*apps.DaemonSet), b.(*v1beta2.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentSpec)(nil), (*v1beta2.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1beta2_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1beta2.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentStrategy)(nil), (*v1beta2.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1beta2_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1beta2.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.Deployment)(nil), (*v1beta2.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_Deployment_To_v1beta2_Deployment(a.(*apps.Deployment), b.(*v1beta2.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1beta2.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1beta2_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1beta2.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1beta2.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1beta2_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1beta2.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1beta2.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1beta2_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1beta2.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetSpec)(nil), (*v1beta2.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1beta2_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1beta2.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetStatus)(nil), (*v1beta2.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetStatus_To_v1beta2_StatefulSetStatus(a.(*apps.StatefulSetStatus), b.(*v1beta2.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1beta2.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1beta2_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1beta2.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1beta2.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1beta2_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1beta2.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.DaemonSetSpec)(nil), (*apps.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetSpec_To_apps_DaemonSetSpec(a.(*v1beta2.DaemonSetSpec), b.(*apps.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.DaemonSetUpdateStrategy)(nil), (*apps.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSetUpdateStrategy_To_apps_DaemonSetUpdateStrategy(a.(*v1beta2.DaemonSetUpdateStrategy), b.(*apps.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.DaemonSet)(nil), (*apps.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DaemonSet_To_apps_DaemonSet(a.(*v1beta2.DaemonSet), b.(*apps.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1beta2.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1beta2.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.Deployment)(nil), (*apps.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_Deployment_To_apps_Deployment(a.(*v1beta2.Deployment), b.(*apps.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1beta2.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1beta2.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1beta2.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.ScaleStatus)(nil), (*autoscaling.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_ScaleStatus_To_autoscaling_ScaleStatus(a.(*v1beta2.ScaleStatus), b.(*autoscaling.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1beta2.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.StatefulSetStatus)(nil), (*apps.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetStatus_To_apps_StatefulSetStatus(a.(*v1beta2.StatefulSetStatus), b.(*apps.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1beta2.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	return nil
}
