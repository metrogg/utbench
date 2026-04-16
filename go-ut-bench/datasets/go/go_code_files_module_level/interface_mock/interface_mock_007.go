func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.ControllerRevision)(nil), (*apps.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ControllerRevision_To_apps_ControllerRevision(a.(*v1.ControllerRevision), b.(*apps.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevision)(nil), (*v1.ControllerRevision)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevision_To_v1_ControllerRevision(a.(*apps.ControllerRevision), b.(*v1.ControllerRevision), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ControllerRevisionList)(nil), (*apps.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ControllerRevisionList_To_apps_ControllerRevisionList(a.(*v1.ControllerRevisionList), b.(*apps.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ControllerRevisionList)(nil), (*v1.ControllerRevisionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ControllerRevisionList_To_v1_ControllerRevisionList(a.(*apps.ControllerRevisionList), b.(*v1.ControllerRevisionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSet)(nil), (*apps.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSet_To_apps_DaemonSet(a.(*v1.DaemonSet), b.(*apps.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSet)(nil), (*v1.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSet_To_v1_DaemonSet(a.(*apps.DaemonSet), b.(*v1.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSetCondition)(nil), (*apps.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetCondition_To_apps_DaemonSetCondition(a.(*v1.DaemonSetCondition), b.(*apps.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetCondition)(nil), (*v1.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetCondition_To_v1_DaemonSetCondition(a.(*apps.DaemonSetCondition), b.(*v1.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSetList)(nil), (*apps.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetList_To_apps_DaemonSetList(a.(*v1.DaemonSetList), b.(*apps.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetList)(nil), (*v1.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetList_To_v1_DaemonSetList(a.(*apps.DaemonSetList), b.(*v1.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSetSpec)(nil), (*apps.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetSpec_To_apps_DaemonSetSpec(a.(*v1.DaemonSetSpec), b.(*apps.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetSpec)(nil), (*v1.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetSpec_To_v1_DaemonSetSpec(a.(*apps.DaemonSetSpec), b.(*v1.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSetStatus)(nil), (*apps.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetStatus_To_apps_DaemonSetStatus(a.(*v1.DaemonSetStatus), b.(*apps.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetStatus)(nil), (*v1.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetStatus_To_v1_DaemonSetStatus(a.(*apps.DaemonSetStatus), b.(*v1.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonSetUpdateStrategy)(nil), (*apps.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetUpdateStrategy_To_apps_DaemonSetUpdateStrategy(a.(*v1.DaemonSetUpdateStrategy), b.(*apps.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetUpdateStrategy)(nil), (*v1.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetUpdateStrategy_To_v1_DaemonSetUpdateStrategy(a.(*apps.DaemonSetUpdateStrategy), b.(*v1.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Deployment)(nil), (*apps.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Deployment_To_apps_Deployment(a.(*v1.Deployment), b.(*apps.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.Deployment)(nil), (*v1.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_Deployment_To_v1_Deployment(a.(*apps.Deployment), b.(*v1.Deployment), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentList)(nil), (*apps.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentList_To_apps_DeploymentList(a.(*v1.DeploymentList), b.(*apps.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentList)(nil), (*v1.DeploymentList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentList_To_v1_DeploymentList(a.(*apps.DeploymentList), b.(*v1.DeploymentList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentSpec)(nil), (*v1.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DeploymentStatus)(nil), (*apps.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentStatus_To_apps_DeploymentStatus(a.(*v1.DeploymentStatus), b.(*apps.DeploymentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DeploymentStatus)(nil), (*v1.DeploymentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStatus_To_v1_DeploymentStatus(a.(*apps.DeploymentStatus), b.(*v1.DeploymentStatus), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1.ReplicaSet)(nil), (*apps.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSet_To_apps_ReplicaSet(a.(*v1.ReplicaSet), b.(*apps.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSet)(nil), (*v1.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSet_To_v1_ReplicaSet(a.(*apps.ReplicaSet), b.(*v1.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicaSetCondition)(nil), (*apps.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSetCondition_To_apps_ReplicaSetCondition(a.(*v1.ReplicaSetCondition), b.(*apps.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetCondition)(nil), (*v1.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetCondition_To_v1_ReplicaSetCondition(a.(*apps.ReplicaSetCondition), b.(*v1.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicaSetList)(nil), (*apps.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSetList_To_apps_ReplicaSetList(a.(*v1.ReplicaSetList), b.(*apps.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetList)(nil), (*v1.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetList_To_v1_ReplicaSetList(a.(*apps.ReplicaSetList), b.(*v1.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicaSetStatus)(nil), (*apps.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSetStatus_To_apps_ReplicaSetStatus(a.(*v1.ReplicaSetStatus), b.(*apps.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetStatus)(nil), (*v1.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetStatus_To_v1_ReplicaSetStatus(a.(*apps.ReplicaSetStatus), b.(*v1.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RollingUpdateStatefulSetStrategy)(nil), (*apps.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingUpdateStatefulSetStrategy_To_apps_RollingUpdateStatefulSetStrategy(a.(*v1.RollingUpdateStatefulSetStrategy), b.(*apps.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateStatefulSetStrategy)(nil), (*v1.RollingUpdateStatefulSetStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateStatefulSetStrategy_To_v1_RollingUpdateStatefulSetStrategy(a.(*apps.RollingUpdateStatefulSetStrategy), b.(*v1.RollingUpdateStatefulSetStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSet)(nil), (*apps.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSet_To_apps_StatefulSet(a.(*v1.StatefulSet), b.(*apps.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSet)(nil), (*v1.StatefulSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSet_To_v1_StatefulSet(a.(*apps.StatefulSet), b.(*v1.StatefulSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSetCondition)(nil), (*apps.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetCondition_To_apps_StatefulSetCondition(a.(*v1.StatefulSetCondition), b.(*apps.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetCondition)(nil), (*v1.StatefulSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetCondition_To_v1_StatefulSetCondition(a.(*apps.StatefulSetCondition), b.(*v1.StatefulSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSetList)(nil), (*apps.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetList_To_apps_StatefulSetList(a.(*v1.StatefulSetList), b.(*apps.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetList)(nil), (*v1.StatefulSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetList_To_v1_StatefulSetList(a.(*apps.StatefulSetList), b.(*v1.StatefulSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetSpec)(nil), (*v1.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSetStatus)(nil), (*apps.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetStatus_To_apps_StatefulSetStatus(a.(*v1.StatefulSetStatus), b.(*apps.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetStatus)(nil), (*v1.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetStatus_To_v1_StatefulSetStatus(a.(*apps.StatefulSetStatus), b.(*v1.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSetSpec)(nil), (*v1.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetSpec_To_v1_DaemonSetSpec(a.(*apps.DaemonSetSpec), b.(*v1.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSetUpdateStrategy)(nil), (*v1.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetUpdateStrategy_To_v1_DaemonSetUpdateStrategy(a.(*apps.DaemonSetUpdateStrategy), b.(*v1.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DaemonSet)(nil), (*v1.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSet_To_v1_DaemonSet(a.(*apps.DaemonSet), b.(*v1.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentSpec)(nil), (*v1.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentSpec_To_v1_DeploymentSpec(a.(*apps.DeploymentSpec), b.(*v1.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.DeploymentStrategy)(nil), (*v1.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DeploymentStrategy_To_v1_DeploymentStrategy(a.(*apps.DeploymentStrategy), b.(*v1.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.Deployment)(nil), (*v1.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_Deployment_To_v1_Deployment(a.(*apps.Deployment), b.(*v1.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetSpec)(nil), (*v1.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetSpec_To_v1_StatefulSetSpec(a.(*apps.StatefulSetSpec), b.(*v1.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetStatus)(nil), (*v1.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetStatus_To_v1_StatefulSetStatus(a.(*apps.StatefulSetStatus), b.(*v1.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.StatefulSetUpdateStrategy)(nil), (*v1.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_StatefulSetUpdateStrategy_To_v1_StatefulSetUpdateStrategy(a.(*apps.StatefulSetUpdateStrategy), b.(*v1.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DaemonSetSpec)(nil), (*apps.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetSpec_To_apps_DaemonSetSpec(a.(*v1.DaemonSetSpec), b.(*apps.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DaemonSetUpdateStrategy)(nil), (*apps.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSetUpdateStrategy_To_apps_DaemonSetUpdateStrategy(a.(*v1.DaemonSetUpdateStrategy), b.(*apps.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DaemonSet)(nil), (*apps.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonSet_To_apps_DaemonSet(a.(*v1.DaemonSet), b.(*apps.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DeploymentSpec)(nil), (*apps.DeploymentSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentSpec_To_apps_DeploymentSpec(a.(*v1.DeploymentSpec), b.(*apps.DeploymentSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.DeploymentStrategy)(nil), (*apps.DeploymentStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DeploymentStrategy_To_apps_DeploymentStrategy(a.(*v1.DeploymentStrategy), b.(*apps.DeploymentStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.Deployment)(nil), (*apps.Deployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Deployment_To_apps_Deployment(a.(*v1.Deployment), b.(*apps.Deployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.RollingUpdateDeployment)(nil), (*apps.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RollingUpdateDeployment_To_apps_RollingUpdateDeployment(a.(*v1.RollingUpdateDeployment), b.(*apps.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.StatefulSetSpec)(nil), (*apps.StatefulSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetSpec_To_apps_StatefulSetSpec(a.(*v1.StatefulSetSpec), b.(*apps.StatefulSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.StatefulSetStatus)(nil), (*apps.StatefulSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetStatus_To_apps_StatefulSetStatus(a.(*v1.StatefulSetStatus), b.(*apps.StatefulSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.StatefulSetUpdateStrategy)(nil), (*apps.StatefulSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(a.(*v1.StatefulSetUpdateStrategy), b.(*apps.StatefulSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	return nil
}
