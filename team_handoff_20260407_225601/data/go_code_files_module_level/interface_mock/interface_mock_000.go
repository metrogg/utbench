func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.AWSElasticBlockStoreVolumeSource)(nil), (*core.AWSElasticBlockStoreVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AWSElasticBlockStoreVolumeSource_To_core_AWSElasticBlockStoreVolumeSource(a.(*v1.AWSElasticBlockStoreVolumeSource), b.(*core.AWSElasticBlockStoreVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AWSElasticBlockStoreVolumeSource)(nil), (*v1.AWSElasticBlockStoreVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AWSElasticBlockStoreVolumeSource_To_v1_AWSElasticBlockStoreVolumeSource(a.(*core.AWSElasticBlockStoreVolumeSource), b.(*v1.AWSElasticBlockStoreVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Affinity)(nil), (*core.Affinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Affinity_To_core_Affinity(a.(*v1.Affinity), b.(*core.Affinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Affinity)(nil), (*v1.Affinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Affinity_To_v1_Affinity(a.(*core.Affinity), b.(*v1.Affinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AttachedVolume)(nil), (*core.AttachedVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AttachedVolume_To_core_AttachedVolume(a.(*v1.AttachedVolume), b.(*core.AttachedVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AttachedVolume)(nil), (*v1.AttachedVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AttachedVolume_To_v1_AttachedVolume(a.(*core.AttachedVolume), b.(*v1.AttachedVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AvoidPods)(nil), (*core.AvoidPods)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AvoidPods_To_core_AvoidPods(a.(*v1.AvoidPods), b.(*core.AvoidPods), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AvoidPods)(nil), (*v1.AvoidPods)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AvoidPods_To_v1_AvoidPods(a.(*core.AvoidPods), b.(*v1.AvoidPods), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AzureDiskVolumeSource)(nil), (*core.AzureDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AzureDiskVolumeSource_To_core_AzureDiskVolumeSource(a.(*v1.AzureDiskVolumeSource), b.(*core.AzureDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AzureDiskVolumeSource)(nil), (*v1.AzureDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AzureDiskVolumeSource_To_v1_AzureDiskVolumeSource(a.(*core.AzureDiskVolumeSource), b.(*v1.AzureDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AzureFilePersistentVolumeSource)(nil), (*core.AzureFilePersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AzureFilePersistentVolumeSource_To_core_AzureFilePersistentVolumeSource(a.(*v1.AzureFilePersistentVolumeSource), b.(*core.AzureFilePersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AzureFilePersistentVolumeSource)(nil), (*v1.AzureFilePersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AzureFilePersistentVolumeSource_To_v1_AzureFilePersistentVolumeSource(a.(*core.AzureFilePersistentVolumeSource), b.(*v1.AzureFilePersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.AzureFileVolumeSource)(nil), (*core.AzureFileVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AzureFileVolumeSource_To_core_AzureFileVolumeSource(a.(*v1.AzureFileVolumeSource), b.(*core.AzureFileVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.AzureFileVolumeSource)(nil), (*v1.AzureFileVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_AzureFileVolumeSource_To_v1_AzureFileVolumeSource(a.(*core.AzureFileVolumeSource), b.(*v1.AzureFileVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Binding)(nil), (*core.Binding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Binding_To_core_Binding(a.(*v1.Binding), b.(*core.Binding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Binding)(nil), (*v1.Binding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Binding_To_v1_Binding(a.(*core.Binding), b.(*v1.Binding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CSIPersistentVolumeSource)(nil), (*core.CSIPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CSIPersistentVolumeSource_To_core_CSIPersistentVolumeSource(a.(*v1.CSIPersistentVolumeSource), b.(*core.CSIPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CSIPersistentVolumeSource)(nil), (*v1.CSIPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CSIPersistentVolumeSource_To_v1_CSIPersistentVolumeSource(a.(*core.CSIPersistentVolumeSource), b.(*v1.CSIPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CSIVolumeSource)(nil), (*core.CSIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CSIVolumeSource_To_core_CSIVolumeSource(a.(*v1.CSIVolumeSource), b.(*core.CSIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CSIVolumeSource)(nil), (*v1.CSIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CSIVolumeSource_To_v1_CSIVolumeSource(a.(*core.CSIVolumeSource), b.(*v1.CSIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Capabilities)(nil), (*core.Capabilities)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Capabilities_To_core_Capabilities(a.(*v1.Capabilities), b.(*core.Capabilities), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Capabilities)(nil), (*v1.Capabilities)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Capabilities_To_v1_Capabilities(a.(*core.Capabilities), b.(*v1.Capabilities), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CephFSPersistentVolumeSource)(nil), (*core.CephFSPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CephFSPersistentVolumeSource_To_core_CephFSPersistentVolumeSource(a.(*v1.CephFSPersistentVolumeSource), b.(*core.CephFSPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CephFSPersistentVolumeSource)(nil), (*v1.CephFSPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CephFSPersistentVolumeSource_To_v1_CephFSPersistentVolumeSource(a.(*core.CephFSPersistentVolumeSource), b.(*v1.CephFSPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CephFSVolumeSource)(nil), (*core.CephFSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CephFSVolumeSource_To_core_CephFSVolumeSource(a.(*v1.CephFSVolumeSource), b.(*core.CephFSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CephFSVolumeSource)(nil), (*v1.CephFSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CephFSVolumeSource_To_v1_CephFSVolumeSource(a.(*core.CephFSVolumeSource), b.(*v1.CephFSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CinderPersistentVolumeSource)(nil), (*core.CinderPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CinderPersistentVolumeSource_To_core_CinderPersistentVolumeSource(a.(*v1.CinderPersistentVolumeSource), b.(*core.CinderPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CinderPersistentVolumeSource)(nil), (*v1.CinderPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CinderPersistentVolumeSource_To_v1_CinderPersistentVolumeSource(a.(*core.CinderPersistentVolumeSource), b.(*v1.CinderPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.CinderVolumeSource)(nil), (*core.CinderVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CinderVolumeSource_To_core_CinderVolumeSource(a.(*v1.CinderVolumeSource), b.(*core.CinderVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.CinderVolumeSource)(nil), (*v1.CinderVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_CinderVolumeSource_To_v1_CinderVolumeSource(a.(*core.CinderVolumeSource), b.(*v1.CinderVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClientIPConfig)(nil), (*core.ClientIPConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClientIPConfig_To_core_ClientIPConfig(a.(*v1.ClientIPConfig), b.(*core.ClientIPConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ClientIPConfig)(nil), (*v1.ClientIPConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ClientIPConfig_To_v1_ClientIPConfig(a.(*core.ClientIPConfig), b.(*v1.ClientIPConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ComponentCondition)(nil), (*core.ComponentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ComponentCondition_To_core_ComponentCondition(a.(*v1.ComponentCondition), b.(*core.ComponentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ComponentCondition)(nil), (*v1.ComponentCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ComponentCondition_To_v1_ComponentCondition(a.(*core.ComponentCondition), b.(*v1.ComponentCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ComponentStatus)(nil), (*core.ComponentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ComponentStatus_To_core_ComponentStatus(a.(*v1.ComponentStatus), b.(*core.ComponentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ComponentStatus)(nil), (*v1.ComponentStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ComponentStatus_To_v1_ComponentStatus(a.(*core.ComponentStatus), b.(*v1.ComponentStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ComponentStatusList)(nil), (*core.ComponentStatusList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ComponentStatusList_To_core_ComponentStatusList(a.(*v1.ComponentStatusList), b.(*core.ComponentStatusList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ComponentStatusList)(nil), (*v1.ComponentStatusList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ComponentStatusList_To_v1_ComponentStatusList(a.(*core.ComponentStatusList), b.(*v1.ComponentStatusList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMap)(nil), (*core.ConfigMap)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMap_To_core_ConfigMap(a.(*v1.ConfigMap), b.(*core.ConfigMap), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMap)(nil), (*v1.ConfigMap)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMap_To_v1_ConfigMap(a.(*core.ConfigMap), b.(*v1.ConfigMap), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapEnvSource)(nil), (*core.ConfigMapEnvSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapEnvSource_To_core_ConfigMapEnvSource(a.(*v1.ConfigMapEnvSource), b.(*core.ConfigMapEnvSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapEnvSource)(nil), (*v1.ConfigMapEnvSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapEnvSource_To_v1_ConfigMapEnvSource(a.(*core.ConfigMapEnvSource), b.(*v1.ConfigMapEnvSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapKeySelector)(nil), (*core.ConfigMapKeySelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapKeySelector_To_core_ConfigMapKeySelector(a.(*v1.ConfigMapKeySelector), b.(*core.ConfigMapKeySelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapKeySelector)(nil), (*v1.ConfigMapKeySelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapKeySelector_To_v1_ConfigMapKeySelector(a.(*core.ConfigMapKeySelector), b.(*v1.ConfigMapKeySelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapList)(nil), (*core.ConfigMapList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapList_To_core_ConfigMapList(a.(*v1.ConfigMapList), b.(*core.ConfigMapList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapList)(nil), (*v1.ConfigMapList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapList_To_v1_ConfigMapList(a.(*core.ConfigMapList), b.(*v1.ConfigMapList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapNodeConfigSource)(nil), (*core.ConfigMapNodeConfigSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapNodeConfigSource_To_core_ConfigMapNodeConfigSource(a.(*v1.ConfigMapNodeConfigSource), b.(*core.ConfigMapNodeConfigSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapNodeConfigSource)(nil), (*v1.ConfigMapNodeConfigSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapNodeConfigSource_To_v1_ConfigMapNodeConfigSource(a.(*core.ConfigMapNodeConfigSource), b.(*v1.ConfigMapNodeConfigSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapProjection)(nil), (*core.ConfigMapProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapProjection_To_core_ConfigMapProjection(a.(*v1.ConfigMapProjection), b.(*core.ConfigMapProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapProjection)(nil), (*v1.ConfigMapProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapProjection_To_v1_ConfigMapProjection(a.(*core.ConfigMapProjection), b.(*v1.ConfigMapProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ConfigMapVolumeSource)(nil), (*core.ConfigMapVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ConfigMapVolumeSource_To_core_ConfigMapVolumeSource(a.(*v1.ConfigMapVolumeSource), b.(*core.ConfigMapVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ConfigMapVolumeSource)(nil), (*v1.ConfigMapVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ConfigMapVolumeSource_To_v1_ConfigMapVolumeSource(a.(*core.ConfigMapVolumeSource), b.(*v1.ConfigMapVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Container)(nil), (*core.Container)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Container_To_core_Container(a.(*v1.Container), b.(*core.Container), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Container)(nil), (*v1.Container)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Container_To_v1_Container(a.(*core.Container), b.(*v1.Container), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerImage)(nil), (*core.ContainerImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerImage_To_core_ContainerImage(a.(*v1.ContainerImage), b.(*core.ContainerImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerImage)(nil), (*v1.ContainerImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerImage_To_v1_ContainerImage(a.(*core.ContainerImage), b.(*v1.ContainerImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerPort)(nil), (*core.ContainerPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerPort_To_core_ContainerPort(a.(*v1.ContainerPort), b.(*core.ContainerPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerPort)(nil), (*v1.ContainerPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerPort_To_v1_ContainerPort(a.(*core.ContainerPort), b.(*v1.ContainerPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerState)(nil), (*core.ContainerState)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerState_To_core_ContainerState(a.(*v1.ContainerState), b.(*core.ContainerState), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerState)(nil), (*v1.ContainerState)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerState_To_v1_ContainerState(a.(*core.ContainerState), b.(*v1.ContainerState), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerStateRunning)(nil), (*core.ContainerStateRunning)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerStateRunning_To_core_ContainerStateRunning(a.(*v1.ContainerStateRunning), b.(*core.ContainerStateRunning), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerStateRunning)(nil), (*v1.ContainerStateRunning)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerStateRunning_To_v1_ContainerStateRunning(a.(*core.ContainerStateRunning), b.(*v1.ContainerStateRunning), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerStateTerminated)(nil), (*core.ContainerStateTerminated)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerStateTerminated_To_core_ContainerStateTerminated(a.(*v1.ContainerStateTerminated), b.(*core.ContainerStateTerminated), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerStateTerminated)(nil), (*v1.ContainerStateTerminated)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerStateTerminated_To_v1_ContainerStateTerminated(a.(*core.ContainerStateTerminated), b.(*v1.ContainerStateTerminated), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerStateWaiting)(nil), (*core.ContainerStateWaiting)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerStateWaiting_To_core_ContainerStateWaiting(a.(*v1.ContainerStateWaiting), b.(*core.ContainerStateWaiting), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerStateWaiting)(nil), (*v1.ContainerStateWaiting)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerStateWaiting_To_v1_ContainerStateWaiting(a.(*core.ContainerStateWaiting), b.(*v1.ContainerStateWaiting), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ContainerStatus)(nil), (*core.ContainerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ContainerStatus_To_core_ContainerStatus(a.(*v1.ContainerStatus), b.(*core.ContainerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ContainerStatus)(nil), (*v1.ContainerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ContainerStatus_To_v1_ContainerStatus(a.(*core.ContainerStatus), b.(*v1.ContainerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DaemonEndpoint)(nil), (*core.DaemonEndpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DaemonEndpoint_To_core_DaemonEndpoint(a.(*v1.DaemonEndpoint), b.(*core.DaemonEndpoint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.DaemonEndpoint)(nil), (*v1.DaemonEndpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_DaemonEndpoint_To_v1_DaemonEndpoint(a.(*core.DaemonEndpoint), b.(*v1.DaemonEndpoint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DownwardAPIProjection)(nil), (*core.DownwardAPIProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DownwardAPIProjection_To_core_DownwardAPIProjection(a.(*v1.DownwardAPIProjection), b.(*core.DownwardAPIProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.DownwardAPIProjection)(nil), (*v1.DownwardAPIProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_DownwardAPIProjection_To_v1_DownwardAPIProjection(a.(*core.DownwardAPIProjection), b.(*v1.DownwardAPIProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DownwardAPIVolumeFile)(nil), (*core.DownwardAPIVolumeFile)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DownwardAPIVolumeFile_To_core_DownwardAPIVolumeFile(a.(*v1.DownwardAPIVolumeFile), b.(*core.DownwardAPIVolumeFile), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.DownwardAPIVolumeFile)(nil), (*v1.DownwardAPIVolumeFile)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_DownwardAPIVolumeFile_To_v1_DownwardAPIVolumeFile(a.(*core.DownwardAPIVolumeFile), b.(*v1.DownwardAPIVolumeFile), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.DownwardAPIVolumeSource)(nil), (*core.DownwardAPIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_DownwardAPIVolumeSource_To_core_DownwardAPIVolumeSource(a.(*v1.DownwardAPIVolumeSource), b.(*core.DownwardAPIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.DownwardAPIVolumeSource)(nil), (*v1.DownwardAPIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_DownwardAPIVolumeSource_To_v1_DownwardAPIVolumeSource(a.(*core.DownwardAPIVolumeSource), b.(*v1.DownwardAPIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EmptyDirVolumeSource)(nil), (*core.EmptyDirVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EmptyDirVolumeSource_To_core_EmptyDirVolumeSource(a.(*v1.EmptyDirVolumeSource), b.(*core.EmptyDirVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EmptyDirVolumeSource)(nil), (*v1.EmptyDirVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EmptyDirVolumeSource_To_v1_EmptyDirVolumeSource(a.(*core.EmptyDirVolumeSource), b.(*v1.EmptyDirVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EndpointAddress)(nil), (*core.EndpointAddress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EndpointAddress_To_core_EndpointAddress(a.(*v1.EndpointAddress), b.(*core.EndpointAddress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EndpointAddress)(nil), (*v1.EndpointAddress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EndpointAddress_To_v1_EndpointAddress(a.(*core.EndpointAddress), b.(*v1.EndpointAddress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EndpointPort)(nil), (*core.EndpointPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EndpointPort_To_core_EndpointPort(a.(*v1.EndpointPort), b.(*core.EndpointPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EndpointPort)(nil), (*v1.EndpointPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EndpointPort_To_v1_EndpointPort(a.(*core.EndpointPort), b.(*v1.EndpointPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EndpointSubset)(nil), (*core.EndpointSubset)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EndpointSubset_To_core_EndpointSubset(a.(*v1.EndpointSubset), b.(*core.EndpointSubset), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EndpointSubset)(nil), (*v1.EndpointSubset)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EndpointSubset_To_v1_EndpointSubset(a.(*core.EndpointSubset), b.(*v1.EndpointSubset), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Endpoints)(nil), (*core.Endpoints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Endpoints_To_core_Endpoints(a.(*v1.Endpoints), b.(*core.Endpoints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Endpoints)(nil), (*v1.Endpoints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Endpoints_To_v1_Endpoints(a.(*core.Endpoints), b.(*v1.Endpoints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EndpointsList)(nil), (*core.EndpointsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EndpointsList_To_core_EndpointsList(a.(*v1.EndpointsList), b.(*core.EndpointsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EndpointsList)(nil), (*v1.EndpointsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EndpointsList_To_v1_EndpointsList(a.(*core.EndpointsList), b.(*v1.EndpointsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EnvFromSource)(nil), (*core.EnvFromSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EnvFromSource_To_core_EnvFromSource(a.(*v1.EnvFromSource), b.(*core.EnvFromSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EnvFromSource)(nil), (*v1.EnvFromSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EnvFromSource_To_v1_EnvFromSource(a.(*core.EnvFromSource), b.(*v1.EnvFromSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EnvVar)(nil), (*core.EnvVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EnvVar_To_core_EnvVar(a.(*v1.EnvVar), b.(*core.EnvVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EnvVar)(nil), (*v1.EnvVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EnvVar_To_v1_EnvVar(a.(*core.EnvVar), b.(*v1.EnvVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EnvVarSource)(nil), (*core.EnvVarSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EnvVarSource_To_core_EnvVarSource(a.(*v1.EnvVarSource), b.(*core.EnvVarSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EnvVarSource)(nil), (*v1.EnvVarSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EnvVarSource_To_v1_EnvVarSource(a.(*core.EnvVarSource), b.(*v1.EnvVarSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Event)(nil), (*core.Event)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Event_To_core_Event(a.(*v1.Event), b.(*core.Event), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Event)(nil), (*v1.Event)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Event_To_v1_Event(a.(*core.Event), b.(*v1.Event), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EventList)(nil), (*core.EventList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EventList_To_core_EventList(a.(*v1.EventList), b.(*core.EventList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EventList)(nil), (*v1.EventList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EventList_To_v1_EventList(a.(*core.EventList), b.(*v1.EventList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EventSeries)(nil), (*core.EventSeries)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EventSeries_To_core_EventSeries(a.(*v1.EventSeries), b.(*core.EventSeries), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EventSeries)(nil), (*v1.EventSeries)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EventSeries_To_v1_EventSeries(a.(*core.EventSeries), b.(*v1.EventSeries), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.EventSource)(nil), (*core.EventSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_EventSource_To_core_EventSource(a.(*v1.EventSource), b.(*core.EventSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.EventSource)(nil), (*v1.EventSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_EventSource_To_v1_EventSource(a.(*core.EventSource), b.(*v1.EventSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ExecAction)(nil), (*core.ExecAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExecAction_To_core_ExecAction(a.(*v1.ExecAction), b.(*core.ExecAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ExecAction)(nil), (*v1.ExecAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ExecAction_To_v1_ExecAction(a.(*core.ExecAction), b.(*v1.ExecAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.FCVolumeSource)(nil), (*core.FCVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_FCVolumeSource_To_core_FCVolumeSource(a.(*v1.FCVolumeSource), b.(*core.FCVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.FCVolumeSource)(nil), (*v1.FCVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_FCVolumeSource_To_v1_FCVolumeSource(a.(*core.FCVolumeSource), b.(*v1.FCVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.FlexPersistentVolumeSource)(nil), (*core.FlexPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_FlexPersistentVolumeSource_To_core_FlexPersistentVolumeSource(a.(*v1.FlexPersistentVolumeSource), b.(*core.FlexPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.FlexPersistentVolumeSource)(nil), (*v1.FlexPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_FlexPersistentVolumeSource_To_v1_FlexPersistentVolumeSource(a.(*core.FlexPersistentVolumeSource), b.(*v1.FlexPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.FlexVolumeSource)(nil), (*core.FlexVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_FlexVolumeSource_To_core_FlexVolumeSource(a.(*v1.FlexVolumeSource), b.(*core.FlexVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.FlexVolumeSource)(nil), (*v1.FlexVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_FlexVolumeSource_To_v1_FlexVolumeSource(a.(*core.FlexVolumeSource), b.(*v1.FlexVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.FlockerVolumeSource)(nil), (*core.FlockerVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_FlockerVolumeSource_To_core_FlockerVolumeSource(a.(*v1.FlockerVolumeSource), b.(*core.FlockerVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.FlockerVolumeSource)(nil), (*v1.FlockerVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_FlockerVolumeSource_To_v1_FlockerVolumeSource(a.(*core.FlockerVolumeSource), b.(*v1.FlockerVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GCEPersistentDiskVolumeSource)(nil), (*core.GCEPersistentDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GCEPersistentDiskVolumeSource_To_core_GCEPersistentDiskVolumeSource(a.(*v1.GCEPersistentDiskVolumeSource), b.(*core.GCEPersistentDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.GCEPersistentDiskVolumeSource)(nil), (*v1.GCEPersistentDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_GCEPersistentDiskVolumeSource_To_v1_GCEPersistentDiskVolumeSource(a.(*core.GCEPersistentDiskVolumeSource), b.(*v1.GCEPersistentDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GitRepoVolumeSource)(nil), (*core.GitRepoVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GitRepoVolumeSource_To_core_GitRepoVolumeSource(a.(*v1.GitRepoVolumeSource), b.(*core.GitRepoVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.GitRepoVolumeSource)(nil), (*v1.GitRepoVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_GitRepoVolumeSource_To_v1_GitRepoVolumeSource(a.(*core.GitRepoVolumeSource), b.(*v1.GitRepoVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GlusterfsPersistentVolumeSource)(nil), (*core.GlusterfsPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GlusterfsPersistentVolumeSource_To_core_GlusterfsPersistentVolumeSource(a.(*v1.GlusterfsPersistentVolumeSource), b.(*core.GlusterfsPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.GlusterfsPersistentVolumeSource)(nil), (*v1.GlusterfsPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_GlusterfsPersistentVolumeSource_To_v1_GlusterfsPersistentVolumeSource(a.(*core.GlusterfsPersistentVolumeSource), b.(*v1.GlusterfsPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GlusterfsVolumeSource)(nil), (*core.GlusterfsVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GlusterfsVolumeSource_To_core_GlusterfsVolumeSource(a.(*v1.GlusterfsVolumeSource), b.(*core.GlusterfsVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.GlusterfsVolumeSource)(nil), (*v1.GlusterfsVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_GlusterfsVolumeSource_To_v1_GlusterfsVolumeSource(a.(*core.GlusterfsVolumeSource), b.(*v1.GlusterfsVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HTTPGetAction)(nil), (*core.HTTPGetAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HTTPGetAction_To_core_HTTPGetAction(a.(*v1.HTTPGetAction), b.(*core.HTTPGetAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.HTTPGetAction)(nil), (*v1.HTTPGetAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_HTTPGetAction_To_v1_HTTPGetAction(a.(*core.HTTPGetAction), b.(*v1.HTTPGetAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HTTPHeader)(nil), (*core.HTTPHeader)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HTTPHeader_To_core_HTTPHeader(a.(*v1.HTTPHeader), b.(*core.HTTPHeader), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.HTTPHeader)(nil), (*v1.HTTPHeader)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_HTTPHeader_To_v1_HTTPHeader(a.(*core.HTTPHeader), b.(*v1.HTTPHeader), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Handler)(nil), (*core.Handler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Handler_To_core_Handler(a.(*v1.Handler), b.(*core.Handler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Handler)(nil), (*v1.Handler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Handler_To_v1_Handler(a.(*core.Handler), b.(*v1.Handler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HostAlias)(nil), (*core.HostAlias)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HostAlias_To_core_HostAlias(a.(*v1.HostAlias), b.(*core.HostAlias), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.HostAlias)(nil), (*v1.HostAlias)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_HostAlias_To_v1_HostAlias(a.(*core.HostAlias), b.(*v1.HostAlias), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HostPathVolumeSource)(nil), (*core.HostPathVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HostPathVolumeSource_To_core_HostPathVolumeSource(a.(*v1.HostPathVolumeSource), b.(*core.HostPathVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.HostPathVolumeSource)(nil), (*v1.HostPathVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_HostPathVolumeSource_To_v1_HostPathVolumeSource(a.(*core.HostPathVolumeSource), b.(*v1.HostPathVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ISCSIPersistentVolumeSource)(nil), (*core.ISCSIPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ISCSIPersistentVolumeSource_To_core_ISCSIPersistentVolumeSource(a.(*v1.ISCSIPersistentVolumeSource), b.(*core.ISCSIPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ISCSIPersistentVolumeSource)(nil), (*v1.ISCSIPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ISCSIPersistentVolumeSource_To_v1_ISCSIPersistentVolumeSource(a.(*core.ISCSIPersistentVolumeSource), b.(*v1.ISCSIPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ISCSIVolumeSource)(nil), (*core.ISCSIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ISCSIVolumeSource_To_core_ISCSIVolumeSource(a.(*v1.ISCSIVolumeSource), b.(*core.ISCSIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ISCSIVolumeSource)(nil), (*v1.ISCSIVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ISCSIVolumeSource_To_v1_ISCSIVolumeSource(a.(*core.ISCSIVolumeSource), b.(*v1.ISCSIVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.KeyToPath)(nil), (*core.KeyToPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_KeyToPath_To_core_KeyToPath(a.(*v1.KeyToPath), b.(*core.KeyToPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.KeyToPath)(nil), (*v1.KeyToPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_KeyToPath_To_v1_KeyToPath(a.(*core.KeyToPath), b.(*v1.KeyToPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Lifecycle)(nil), (*core.Lifecycle)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Lifecycle_To_core_Lifecycle(a.(*v1.Lifecycle), b.(*core.Lifecycle), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Lifecycle)(nil), (*v1.Lifecycle)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Lifecycle_To_v1_Lifecycle(a.(*core.Lifecycle), b.(*v1.Lifecycle), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LimitRange)(nil), (*core.LimitRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LimitRange_To_core_LimitRange(a.(*v1.LimitRange), b.(*core.LimitRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LimitRange)(nil), (*v1.LimitRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LimitRange_To_v1_LimitRange(a.(*core.LimitRange), b.(*v1.LimitRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LimitRangeItem)(nil), (*core.LimitRangeItem)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LimitRangeItem_To_core_LimitRangeItem(a.(*v1.LimitRangeItem), b.(*core.LimitRangeItem), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LimitRangeItem)(nil), (*v1.LimitRangeItem)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LimitRangeItem_To_v1_LimitRangeItem(a.(*core.LimitRangeItem), b.(*v1.LimitRangeItem), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LimitRangeList)(nil), (*core.LimitRangeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LimitRangeList_To_core_LimitRangeList(a.(*v1.LimitRangeList), b.(*core.LimitRangeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LimitRangeList)(nil), (*v1.LimitRangeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LimitRangeList_To_v1_LimitRangeList(a.(*core.LimitRangeList), b.(*v1.LimitRangeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LimitRangeSpec)(nil), (*core.LimitRangeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LimitRangeSpec_To_core_LimitRangeSpec(a.(*v1.LimitRangeSpec), b.(*core.LimitRangeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LimitRangeSpec)(nil), (*v1.LimitRangeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LimitRangeSpec_To_v1_LimitRangeSpec(a.(*core.LimitRangeSpec), b.(*v1.LimitRangeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.List)(nil), (*core.List)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_List_To_core_List(a.(*v1.List), b.(*core.List), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.List)(nil), (*v1.List)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_List_To_v1_List(a.(*core.List), b.(*v1.List), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LoadBalancerIngress)(nil), (*core.LoadBalancerIngress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LoadBalancerIngress_To_core_LoadBalancerIngress(a.(*v1.LoadBalancerIngress), b.(*core.LoadBalancerIngress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LoadBalancerIngress)(nil), (*v1.LoadBalancerIngress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LoadBalancerIngress_To_v1_LoadBalancerIngress(a.(*core.LoadBalancerIngress), b.(*v1.LoadBalancerIngress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LoadBalancerStatus)(nil), (*core.LoadBalancerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LoadBalancerStatus_To_core_LoadBalancerStatus(a.(*v1.LoadBalancerStatus), b.(*core.LoadBalancerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LoadBalancerStatus)(nil), (*v1.LoadBalancerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LoadBalancerStatus_To_v1_LoadBalancerStatus(a.(*core.LoadBalancerStatus), b.(*v1.LoadBalancerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LocalObjectReference)(nil), (*core.LocalObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalObjectReference_To_core_LocalObjectReference(a.(*v1.LocalObjectReference), b.(*core.LocalObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LocalObjectReference)(nil), (*v1.LocalObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LocalObjectReference_To_v1_LocalObjectReference(a.(*core.LocalObjectReference), b.(*v1.LocalObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LocalVolumeSource)(nil), (*core.LocalVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalVolumeSource_To_core_LocalVolumeSource(a.(*v1.LocalVolumeSource), b.(*core.LocalVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.LocalVolumeSource)(nil), (*v1.LocalVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_LocalVolumeSource_To_v1_LocalVolumeSource(a.(*core.LocalVolumeSource), b.(*v1.LocalVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NFSVolumeSource)(nil), (*core.NFSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NFSVolumeSource_To_core_NFSVolumeSource(a.(*v1.NFSVolumeSource), b.(*core.NFSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NFSVolumeSource)(nil), (*v1.NFSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NFSVolumeSource_To_v1_NFSVolumeSource(a.(*core.NFSVolumeSource), b.(*v1.NFSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Namespace)(nil), (*core.Namespace)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Namespace_To_core_Namespace(a.(*v1.Namespace), b.(*core.Namespace), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Namespace)(nil), (*v1.Namespace)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Namespace_To_v1_Namespace(a.(*core.Namespace), b.(*v1.Namespace), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NamespaceList)(nil), (*core.NamespaceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NamespaceList_To_core_NamespaceList(a.(*v1.NamespaceList), b.(*core.NamespaceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NamespaceList)(nil), (*v1.NamespaceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NamespaceList_To_v1_NamespaceList(a.(*core.NamespaceList), b.(*v1.NamespaceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NamespaceSpec)(nil), (*core.NamespaceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NamespaceSpec_To_core_NamespaceSpec(a.(*v1.NamespaceSpec), b.(*core.NamespaceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NamespaceSpec)(nil), (*v1.NamespaceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NamespaceSpec_To_v1_NamespaceSpec(a.(*core.NamespaceSpec), b.(*v1.NamespaceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NamespaceStatus)(nil), (*core.NamespaceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NamespaceStatus_To_core_NamespaceStatus(a.(*v1.NamespaceStatus), b.(*core.NamespaceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NamespaceStatus)(nil), (*v1.NamespaceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NamespaceStatus_To_v1_NamespaceStatus(a.(*core.NamespaceStatus), b.(*v1.NamespaceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Node)(nil), (*core.Node)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Node_To_core_Node(a.(*v1.Node), b.(*core.Node), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Node)(nil), (*v1.Node)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Node_To_v1_Node(a.(*core.Node), b.(*v1.Node), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeAddress)(nil), (*core.NodeAddress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeAddress_To_core_NodeAddress(a.(*v1.NodeAddress), b.(*core.NodeAddress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeAddress)(nil), (*v1.NodeAddress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeAddress_To_v1_NodeAddress(a.(*core.NodeAddress), b.(*v1.NodeAddress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeAffinity)(nil), (*core.NodeAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeAffinity_To_core_NodeAffinity(a.(*v1.NodeAffinity), b.(*core.NodeAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeAffinity)(nil), (*v1.NodeAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeAffinity_To_v1_NodeAffinity(a.(*core.NodeAffinity), b.(*v1.NodeAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeCondition)(nil), (*core.NodeCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeCondition_To_core_NodeCondition(a.(*v1.NodeCondition), b.(*core.NodeCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeCondition)(nil), (*v1.NodeCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeCondition_To_v1_NodeCondition(a.(*core.NodeCondition), b.(*v1.NodeCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeConfigSource)(nil), (*core.NodeConfigSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeConfigSource_To_core_NodeConfigSource(a.(*v1.NodeConfigSource), b.(*core.NodeConfigSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeConfigSource)(nil), (*v1.NodeConfigSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeConfigSource_To_v1_NodeConfigSource(a.(*core.NodeConfigSource), b.(*v1.NodeConfigSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeConfigStatus)(nil), (*core.NodeConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeConfigStatus_To_core_NodeConfigStatus(a.(*v1.NodeConfigStatus), b.(*core.NodeConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeConfigStatus)(nil), (*v1.NodeConfigStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeConfigStatus_To_v1_NodeConfigStatus(a.(*core.NodeConfigStatus), b.(*v1.NodeConfigStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeDaemonEndpoints)(nil), (*core.NodeDaemonEndpoints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeDaemonEndpoints_To_core_NodeDaemonEndpoints(a.(*v1.NodeDaemonEndpoints), b.(*core.NodeDaemonEndpoints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeDaemonEndpoints)(nil), (*v1.NodeDaemonEndpoints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeDaemonEndpoints_To_v1_NodeDaemonEndpoints(a.(*core.NodeDaemonEndpoints), b.(*v1.NodeDaemonEndpoints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeList)(nil), (*core.NodeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeList_To_core_NodeList(a.(*v1.NodeList), b.(*core.NodeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeList)(nil), (*v1.NodeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeList_To_v1_NodeList(a.(*core.NodeList), b.(*v1.NodeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeProxyOptions)(nil), (*core.NodeProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeProxyOptions_To_core_NodeProxyOptions(a.(*v1.NodeProxyOptions), b.(*core.NodeProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeProxyOptions)(nil), (*v1.NodeProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeProxyOptions_To_v1_NodeProxyOptions(a.(*core.NodeProxyOptions), b.(*v1.NodeProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeResources)(nil), (*core.NodeResources)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeResources_To_core_NodeResources(a.(*v1.NodeResources), b.(*core.NodeResources), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeResources)(nil), (*v1.NodeResources)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeResources_To_v1_NodeResources(a.(*core.NodeResources), b.(*v1.NodeResources), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeSelector)(nil), (*core.NodeSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeSelector_To_core_NodeSelector(a.(*v1.NodeSelector), b.(*core.NodeSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeSelector)(nil), (*v1.NodeSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeSelector_To_v1_NodeSelector(a.(*core.NodeSelector), b.(*v1.NodeSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeSelectorRequirement)(nil), (*core.NodeSelectorRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeSelectorRequirement_To_core_NodeSelectorRequirement(a.(*v1.NodeSelectorRequirement), b.(*core.NodeSelectorRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeSelectorRequirement)(nil), (*v1.NodeSelectorRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeSelectorRequirement_To_v1_NodeSelectorRequirement(a.(*core.NodeSelectorRequirement), b.(*v1.NodeSelectorRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeSelectorTerm)(nil), (*core.NodeSelectorTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeSelectorTerm_To_core_NodeSelectorTerm(a.(*v1.NodeSelectorTerm), b.(*core.NodeSelectorTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeSelectorTerm)(nil), (*v1.NodeSelectorTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeSelectorTerm_To_v1_NodeSelectorTerm(a.(*core.NodeSelectorTerm), b.(*v1.NodeSelectorTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeSpec)(nil), (*core.NodeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeSpec_To_core_NodeSpec(a.(*v1.NodeSpec), b.(*core.NodeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeSpec)(nil), (*v1.NodeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeSpec_To_v1_NodeSpec(a.(*core.NodeSpec), b.(*v1.NodeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeStatus)(nil), (*core.NodeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeStatus_To_core_NodeStatus(a.(*v1.NodeStatus), b.(*core.NodeStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeStatus)(nil), (*v1.NodeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeStatus_To_v1_NodeStatus(a.(*core.NodeStatus), b.(*v1.NodeStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.NodeSystemInfo)(nil), (*core.NodeSystemInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_NodeSystemInfo_To_core_NodeSystemInfo(a.(*v1.NodeSystemInfo), b.(*core.NodeSystemInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.NodeSystemInfo)(nil), (*v1.NodeSystemInfo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_NodeSystemInfo_To_v1_NodeSystemInfo(a.(*core.NodeSystemInfo), b.(*v1.NodeSystemInfo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ObjectFieldSelector)(nil), (*core.ObjectFieldSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectFieldSelector_To_core_ObjectFieldSelector(a.(*v1.ObjectFieldSelector), b.(*core.ObjectFieldSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ObjectFieldSelector)(nil), (*v1.ObjectFieldSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ObjectFieldSelector_To_v1_ObjectFieldSelector(a.(*core.ObjectFieldSelector), b.(*v1.ObjectFieldSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ObjectReference)(nil), (*core.ObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectReference_To_core_ObjectReference(a.(*v1.ObjectReference), b.(*core.ObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ObjectReference)(nil), (*v1.ObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ObjectReference_To_v1_ObjectReference(a.(*core.ObjectReference), b.(*v1.ObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolume)(nil), (*core.PersistentVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolume_To_core_PersistentVolume(a.(*v1.PersistentVolume), b.(*core.PersistentVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolume)(nil), (*v1.PersistentVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolume_To_v1_PersistentVolume(a.(*core.PersistentVolume), b.(*v1.PersistentVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaim)(nil), (*core.PersistentVolumeClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaim_To_core_PersistentVolumeClaim(a.(*v1.PersistentVolumeClaim), b.(*core.PersistentVolumeClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaim)(nil), (*v1.PersistentVolumeClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaim_To_v1_PersistentVolumeClaim(a.(*core.PersistentVolumeClaim), b.(*v1.PersistentVolumeClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaimCondition)(nil), (*core.PersistentVolumeClaimCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaimCondition_To_core_PersistentVolumeClaimCondition(a.(*v1.PersistentVolumeClaimCondition), b.(*core.PersistentVolumeClaimCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaimCondition)(nil), (*v1.PersistentVolumeClaimCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaimCondition_To_v1_PersistentVolumeClaimCondition(a.(*core.PersistentVolumeClaimCondition), b.(*v1.PersistentVolumeClaimCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaimList)(nil), (*core.PersistentVolumeClaimList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaimList_To_core_PersistentVolumeClaimList(a.(*v1.PersistentVolumeClaimList), b.(*core.PersistentVolumeClaimList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaimList)(nil), (*v1.PersistentVolumeClaimList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaimList_To_v1_PersistentVolumeClaimList(a.(*core.PersistentVolumeClaimList), b.(*v1.PersistentVolumeClaimList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaimSpec)(nil), (*core.PersistentVolumeClaimSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaimSpec_To_core_PersistentVolumeClaimSpec(a.(*v1.PersistentVolumeClaimSpec), b.(*core.PersistentVolumeClaimSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaimSpec)(nil), (*v1.PersistentVolumeClaimSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaimSpec_To_v1_PersistentVolumeClaimSpec(a.(*core.PersistentVolumeClaimSpec), b.(*v1.PersistentVolumeClaimSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaimStatus)(nil), (*core.PersistentVolumeClaimStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaimStatus_To_core_PersistentVolumeClaimStatus(a.(*v1.PersistentVolumeClaimStatus), b.(*core.PersistentVolumeClaimStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaimStatus)(nil), (*v1.PersistentVolumeClaimStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaimStatus_To_v1_PersistentVolumeClaimStatus(a.(*core.PersistentVolumeClaimStatus), b.(*v1.PersistentVolumeClaimStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeClaimVolumeSource)(nil), (*core.PersistentVolumeClaimVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeClaimVolumeSource_To_core_PersistentVolumeClaimVolumeSource(a.(*v1.PersistentVolumeClaimVolumeSource), b.(*core.PersistentVolumeClaimVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeClaimVolumeSource)(nil), (*v1.PersistentVolumeClaimVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeClaimVolumeSource_To_v1_PersistentVolumeClaimVolumeSource(a.(*core.PersistentVolumeClaimVolumeSource), b.(*v1.PersistentVolumeClaimVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeList)(nil), (*core.PersistentVolumeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeList_To_core_PersistentVolumeList(a.(*v1.PersistentVolumeList), b.(*core.PersistentVolumeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeList)(nil), (*v1.PersistentVolumeList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeList_To_v1_PersistentVolumeList(a.(*core.PersistentVolumeList), b.(*v1.PersistentVolumeList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeSource)(nil), (*core.PersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeSource_To_core_PersistentVolumeSource(a.(*v1.PersistentVolumeSource), b.(*core.PersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeSource)(nil), (*v1.PersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeSource_To_v1_PersistentVolumeSource(a.(*core.PersistentVolumeSource), b.(*v1.PersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeSpec)(nil), (*core.PersistentVolumeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeSpec_To_core_PersistentVolumeSpec(a.(*v1.PersistentVolumeSpec), b.(*core.PersistentVolumeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeSpec)(nil), (*v1.PersistentVolumeSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeSpec_To_v1_PersistentVolumeSpec(a.(*core.PersistentVolumeSpec), b.(*v1.PersistentVolumeSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PersistentVolumeStatus)(nil), (*core.PersistentVolumeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PersistentVolumeStatus_To_core_PersistentVolumeStatus(a.(*v1.PersistentVolumeStatus), b.(*core.PersistentVolumeStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PersistentVolumeStatus)(nil), (*v1.PersistentVolumeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PersistentVolumeStatus_To_v1_PersistentVolumeStatus(a.(*core.PersistentVolumeStatus), b.(*v1.PersistentVolumeStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PhotonPersistentDiskVolumeSource)(nil), (*core.PhotonPersistentDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PhotonPersistentDiskVolumeSource_To_core_PhotonPersistentDiskVolumeSource(a.(*v1.PhotonPersistentDiskVolumeSource), b.(*core.PhotonPersistentDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PhotonPersistentDiskVolumeSource)(nil), (*v1.PhotonPersistentDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PhotonPersistentDiskVolumeSource_To_v1_PhotonPersistentDiskVolumeSource(a.(*core.PhotonPersistentDiskVolumeSource), b.(*v1.PhotonPersistentDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Pod)(nil), (*core.Pod)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Pod_To_core_Pod(a.(*v1.Pod), b.(*core.Pod), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Pod)(nil), (*v1.Pod)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Pod_To_v1_Pod(a.(*core.Pod), b.(*v1.Pod), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodAffinity)(nil), (*core.PodAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodAffinity_To_core_PodAffinity(a.(*v1.PodAffinity), b.(*core.PodAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodAffinity)(nil), (*v1.PodAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodAffinity_To_v1_PodAffinity(a.(*core.PodAffinity), b.(*v1.PodAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodAffinityTerm)(nil), (*core.PodAffinityTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodAffinityTerm_To_core_PodAffinityTerm(a.(*v1.PodAffinityTerm), b.(*core.PodAffinityTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodAffinityTerm)(nil), (*v1.PodAffinityTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodAffinityTerm_To_v1_PodAffinityTerm(a.(*core.PodAffinityTerm), b.(*v1.PodAffinityTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodAntiAffinity)(nil), (*core.PodAntiAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodAntiAffinity_To_core_PodAntiAffinity(a.(*v1.PodAntiAffinity), b.(*core.PodAntiAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodAntiAffinity)(nil), (*v1.PodAntiAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodAntiAffinity_To_v1_PodAntiAffinity(a.(*core.PodAntiAffinity), b.(*v1.PodAntiAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodAttachOptions)(nil), (*core.PodAttachOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodAttachOptions_To_core_PodAttachOptions(a.(*v1.PodAttachOptions), b.(*core.PodAttachOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodAttachOptions)(nil), (*v1.PodAttachOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodAttachOptions_To_v1_PodAttachOptions(a.(*core.PodAttachOptions), b.(*v1.PodAttachOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodCondition)(nil), (*core.PodCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodCondition_To_core_PodCondition(a.(*v1.PodCondition), b.(*core.PodCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodCondition)(nil), (*v1.PodCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodCondition_To_v1_PodCondition(a.(*core.PodCondition), b.(*v1.PodCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodDNSConfig)(nil), (*core.PodDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodDNSConfig_To_core_PodDNSConfig(a.(*v1.PodDNSConfig), b.(*core.PodDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodDNSConfig)(nil), (*v1.PodDNSConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodDNSConfig_To_v1_PodDNSConfig(a.(*core.PodDNSConfig), b.(*v1.PodDNSConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodDNSConfigOption)(nil), (*core.PodDNSConfigOption)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodDNSConfigOption_To_core_PodDNSConfigOption(a.(*v1.PodDNSConfigOption), b.(*core.PodDNSConfigOption), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodDNSConfigOption)(nil), (*v1.PodDNSConfigOption)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodDNSConfigOption_To_v1_PodDNSConfigOption(a.(*core.PodDNSConfigOption), b.(*v1.PodDNSConfigOption), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodExecOptions)(nil), (*core.PodExecOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodExecOptions_To_core_PodExecOptions(a.(*v1.PodExecOptions), b.(*core.PodExecOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodExecOptions)(nil), (*v1.PodExecOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodExecOptions_To_v1_PodExecOptions(a.(*core.PodExecOptions), b.(*v1.PodExecOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodList)(nil), (*core.PodList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodList_To_core_PodList(a.(*v1.PodList), b.(*core.PodList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodList)(nil), (*v1.PodList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodList_To_v1_PodList(a.(*core.PodList), b.(*v1.PodList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodLogOptions)(nil), (*core.PodLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodLogOptions_To_core_PodLogOptions(a.(*v1.PodLogOptions), b.(*core.PodLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodLogOptions)(nil), (*v1.PodLogOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodLogOptions_To_v1_PodLogOptions(a.(*core.PodLogOptions), b.(*v1.PodLogOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodPortForwardOptions)(nil), (*core.PodPortForwardOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodPortForwardOptions_To_core_PodPortForwardOptions(a.(*v1.PodPortForwardOptions), b.(*core.PodPortForwardOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodPortForwardOptions)(nil), (*v1.PodPortForwardOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodPortForwardOptions_To_v1_PodPortForwardOptions(a.(*core.PodPortForwardOptions), b.(*v1.PodPortForwardOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodProxyOptions)(nil), (*core.PodProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodProxyOptions_To_core_PodProxyOptions(a.(*v1.PodProxyOptions), b.(*core.PodProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodProxyOptions)(nil), (*v1.PodProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodProxyOptions_To_v1_PodProxyOptions(a.(*core.PodProxyOptions), b.(*v1.PodProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodReadinessGate)(nil), (*core.PodReadinessGate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodReadinessGate_To_core_PodReadinessGate(a.(*v1.PodReadinessGate), b.(*core.PodReadinessGate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodReadinessGate)(nil), (*v1.PodReadinessGate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodReadinessGate_To_v1_PodReadinessGate(a.(*core.PodReadinessGate), b.(*v1.PodReadinessGate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityContext)(nil), (*core.PodSecurityContext)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityContext_To_core_PodSecurityContext(a.(*v1.PodSecurityContext), b.(*core.PodSecurityContext), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodSecurityContext)(nil), (*v1.PodSecurityContext)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodSecurityContext_To_v1_PodSecurityContext(a.(*core.PodSecurityContext), b.(*v1.PodSecurityContext), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSignature)(nil), (*core.PodSignature)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSignature_To_core_PodSignature(a.(*v1.PodSignature), b.(*core.PodSignature), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodSignature)(nil), (*v1.PodSignature)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodSignature_To_v1_PodSignature(a.(*core.PodSignature), b.(*v1.PodSignature), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSpec)(nil), (*core.PodSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSpec_To_core_PodSpec(a.(*v1.PodSpec), b.(*core.PodSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodSpec)(nil), (*v1.PodSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodSpec_To_v1_PodSpec(a.(*core.PodSpec), b.(*v1.PodSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodStatus)(nil), (*core.PodStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodStatus_To_core_PodStatus(a.(*v1.PodStatus), b.(*core.PodStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodStatus)(nil), (*v1.PodStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodStatus_To_v1_PodStatus(a.(*core.PodStatus), b.(*v1.PodStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodStatusResult)(nil), (*core.PodStatusResult)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodStatusResult_To_core_PodStatusResult(a.(*v1.PodStatusResult), b.(*core.PodStatusResult), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodStatusResult)(nil), (*v1.PodStatusResult)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodStatusResult_To_v1_PodStatusResult(a.(*core.PodStatusResult), b.(*v1.PodStatusResult), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodTemplate)(nil), (*core.PodTemplate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodTemplate_To_core_PodTemplate(a.(*v1.PodTemplate), b.(*core.PodTemplate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodTemplate)(nil), (*v1.PodTemplate)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodTemplate_To_v1_PodTemplate(a.(*core.PodTemplate), b.(*v1.PodTemplate), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodTemplateList)(nil), (*core.PodTemplateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodTemplateList_To_core_PodTemplateList(a.(*v1.PodTemplateList), b.(*core.PodTemplateList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodTemplateList)(nil), (*v1.PodTemplateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodTemplateList_To_v1_PodTemplateList(a.(*core.PodTemplateList), b.(*v1.PodTemplateList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodTemplateSpec)(nil), (*core.PodTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodTemplateSpec_To_core_PodTemplateSpec(a.(*v1.PodTemplateSpec), b.(*core.PodTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PodTemplateSpec)(nil), (*v1.PodTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodTemplateSpec_To_v1_PodTemplateSpec(a.(*core.PodTemplateSpec), b.(*v1.PodTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PortworxVolumeSource)(nil), (*core.PortworxVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PortworxVolumeSource_To_core_PortworxVolumeSource(a.(*v1.PortworxVolumeSource), b.(*core.PortworxVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PortworxVolumeSource)(nil), (*v1.PortworxVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PortworxVolumeSource_To_v1_PortworxVolumeSource(a.(*core.PortworxVolumeSource), b.(*v1.PortworxVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Preconditions)(nil), (*core.Preconditions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Preconditions_To_core_Preconditions(a.(*v1.Preconditions), b.(*core.Preconditions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Preconditions)(nil), (*v1.Preconditions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Preconditions_To_v1_Preconditions(a.(*core.Preconditions), b.(*v1.Preconditions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PreferAvoidPodsEntry)(nil), (*core.PreferAvoidPodsEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PreferAvoidPodsEntry_To_core_PreferAvoidPodsEntry(a.(*v1.PreferAvoidPodsEntry), b.(*core.PreferAvoidPodsEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PreferAvoidPodsEntry)(nil), (*v1.PreferAvoidPodsEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PreferAvoidPodsEntry_To_v1_PreferAvoidPodsEntry(a.(*core.PreferAvoidPodsEntry), b.(*v1.PreferAvoidPodsEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PreferredSchedulingTerm)(nil), (*core.PreferredSchedulingTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PreferredSchedulingTerm_To_core_PreferredSchedulingTerm(a.(*v1.PreferredSchedulingTerm), b.(*core.PreferredSchedulingTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.PreferredSchedulingTerm)(nil), (*v1.PreferredSchedulingTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PreferredSchedulingTerm_To_v1_PreferredSchedulingTerm(a.(*core.PreferredSchedulingTerm), b.(*v1.PreferredSchedulingTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Probe)(nil), (*core.Probe)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Probe_To_core_Probe(a.(*v1.Probe), b.(*core.Probe), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Probe)(nil), (*v1.Probe)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Probe_To_v1_Probe(a.(*core.Probe), b.(*v1.Probe), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ProjectedVolumeSource)(nil), (*core.ProjectedVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ProjectedVolumeSource_To_core_ProjectedVolumeSource(a.(*v1.ProjectedVolumeSource), b.(*core.ProjectedVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ProjectedVolumeSource)(nil), (*v1.ProjectedVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ProjectedVolumeSource_To_v1_ProjectedVolumeSource(a.(*core.ProjectedVolumeSource), b.(*v1.ProjectedVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.QuobyteVolumeSource)(nil), (*core.QuobyteVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_QuobyteVolumeSource_To_core_QuobyteVolumeSource(a.(*v1.QuobyteVolumeSource), b.(*core.QuobyteVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.QuobyteVolumeSource)(nil), (*v1.QuobyteVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_QuobyteVolumeSource_To_v1_QuobyteVolumeSource(a.(*core.QuobyteVolumeSource), b.(*v1.QuobyteVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RBDPersistentVolumeSource)(nil), (*core.RBDPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RBDPersistentVolumeSource_To_core_RBDPersistentVolumeSource(a.(*v1.RBDPersistentVolumeSource), b.(*core.RBDPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.RBDPersistentVolumeSource)(nil), (*v1.RBDPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_RBDPersistentVolumeSource_To_v1_RBDPersistentVolumeSource(a.(*core.RBDPersistentVolumeSource), b.(*v1.RBDPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RBDVolumeSource)(nil), (*core.RBDVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RBDVolumeSource_To_core_RBDVolumeSource(a.(*v1.RBDVolumeSource), b.(*core.RBDVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.RBDVolumeSource)(nil), (*v1.RBDVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_RBDVolumeSource_To_v1_RBDVolumeSource(a.(*core.RBDVolumeSource), b.(*v1.RBDVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RangeAllocation)(nil), (*core.RangeAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RangeAllocation_To_core_RangeAllocation(a.(*v1.RangeAllocation), b.(*core.RangeAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.RangeAllocation)(nil), (*v1.RangeAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_RangeAllocation_To_v1_RangeAllocation(a.(*core.RangeAllocation), b.(*v1.RangeAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicationController)(nil), (*core.ReplicationController)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationController_To_core_ReplicationController(a.(*v1.ReplicationController), b.(*core.ReplicationController), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ReplicationController)(nil), (*v1.ReplicationController)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationController_To_v1_ReplicationController(a.(*core.ReplicationController), b.(*v1.ReplicationController), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicationControllerCondition)(nil), (*core.ReplicationControllerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerCondition_To_core_ReplicationControllerCondition(a.(*v1.ReplicationControllerCondition), b.(*core.ReplicationControllerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ReplicationControllerCondition)(nil), (*v1.ReplicationControllerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationControllerCondition_To_v1_ReplicationControllerCondition(a.(*core.ReplicationControllerCondition), b.(*v1.ReplicationControllerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicationControllerList)(nil), (*core.ReplicationControllerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerList_To_core_ReplicationControllerList(a.(*v1.ReplicationControllerList), b.(*core.ReplicationControllerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ReplicationControllerList)(nil), (*v1.ReplicationControllerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationControllerList_To_v1_ReplicationControllerList(a.(*core.ReplicationControllerList), b.(*v1.ReplicationControllerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicationControllerSpec)(nil), (*core.ReplicationControllerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerSpec_To_core_ReplicationControllerSpec(a.(*v1.ReplicationControllerSpec), b.(*core.ReplicationControllerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ReplicationControllerSpec)(nil), (*v1.ReplicationControllerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationControllerSpec_To_v1_ReplicationControllerSpec(a.(*core.ReplicationControllerSpec), b.(*v1.ReplicationControllerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ReplicationControllerStatus)(nil), (*core.ReplicationControllerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerStatus_To_core_ReplicationControllerStatus(a.(*v1.ReplicationControllerStatus), b.(*core.ReplicationControllerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ReplicationControllerStatus)(nil), (*v1.ReplicationControllerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationControllerStatus_To_v1_ReplicationControllerStatus(a.(*core.ReplicationControllerStatus), b.(*v1.ReplicationControllerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceFieldSelector)(nil), (*core.ResourceFieldSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceFieldSelector_To_core_ResourceFieldSelector(a.(*v1.ResourceFieldSelector), b.(*core.ResourceFieldSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceFieldSelector)(nil), (*v1.ResourceFieldSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceFieldSelector_To_v1_ResourceFieldSelector(a.(*core.ResourceFieldSelector), b.(*v1.ResourceFieldSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceQuota)(nil), (*core.ResourceQuota)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceQuota_To_core_ResourceQuota(a.(*v1.ResourceQuota), b.(*core.ResourceQuota), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceQuota)(nil), (*v1.ResourceQuota)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceQuota_To_v1_ResourceQuota(a.(*core.ResourceQuota), b.(*v1.ResourceQuota), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceQuotaList)(nil), (*core.ResourceQuotaList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceQuotaList_To_core_ResourceQuotaList(a.(*v1.ResourceQuotaList), b.(*core.ResourceQuotaList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceQuotaList)(nil), (*v1.ResourceQuotaList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceQuotaList_To_v1_ResourceQuotaList(a.(*core.ResourceQuotaList), b.(*v1.ResourceQuotaList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceQuotaSpec)(nil), (*core.ResourceQuotaSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceQuotaSpec_To_core_ResourceQuotaSpec(a.(*v1.ResourceQuotaSpec), b.(*core.ResourceQuotaSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceQuotaSpec)(nil), (*v1.ResourceQuotaSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceQuotaSpec_To_v1_ResourceQuotaSpec(a.(*core.ResourceQuotaSpec), b.(*v1.ResourceQuotaSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceQuotaStatus)(nil), (*core.ResourceQuotaStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceQuotaStatus_To_core_ResourceQuotaStatus(a.(*v1.ResourceQuotaStatus), b.(*core.ResourceQuotaStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceQuotaStatus)(nil), (*v1.ResourceQuotaStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceQuotaStatus_To_v1_ResourceQuotaStatus(a.(*core.ResourceQuotaStatus), b.(*v1.ResourceQuotaStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceRequirements)(nil), (*core.ResourceRequirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceRequirements_To_core_ResourceRequirements(a.(*v1.ResourceRequirements), b.(*core.ResourceRequirements), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ResourceRequirements)(nil), (*v1.ResourceRequirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ResourceRequirements_To_v1_ResourceRequirements(a.(*core.ResourceRequirements), b.(*v1.ResourceRequirements), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SELinuxOptions)(nil), (*core.SELinuxOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SELinuxOptions_To_core_SELinuxOptions(a.(*v1.SELinuxOptions), b.(*core.SELinuxOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SELinuxOptions)(nil), (*v1.SELinuxOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SELinuxOptions_To_v1_SELinuxOptions(a.(*core.SELinuxOptions), b.(*v1.SELinuxOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleIOPersistentVolumeSource)(nil), (*core.ScaleIOPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleIOPersistentVolumeSource_To_core_ScaleIOPersistentVolumeSource(a.(*v1.ScaleIOPersistentVolumeSource), b.(*core.ScaleIOPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ScaleIOPersistentVolumeSource)(nil), (*v1.ScaleIOPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ScaleIOPersistentVolumeSource_To_v1_ScaleIOPersistentVolumeSource(a.(*core.ScaleIOPersistentVolumeSource), b.(*v1.ScaleIOPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleIOVolumeSource)(nil), (*core.ScaleIOVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleIOVolumeSource_To_core_ScaleIOVolumeSource(a.(*v1.ScaleIOVolumeSource), b.(*core.ScaleIOVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ScaleIOVolumeSource)(nil), (*v1.ScaleIOVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ScaleIOVolumeSource_To_v1_ScaleIOVolumeSource(a.(*core.ScaleIOVolumeSource), b.(*v1.ScaleIOVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScopeSelector)(nil), (*core.ScopeSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScopeSelector_To_core_ScopeSelector(a.(*v1.ScopeSelector), b.(*core.ScopeSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ScopeSelector)(nil), (*v1.ScopeSelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ScopeSelector_To_v1_ScopeSelector(a.(*core.ScopeSelector), b.(*v1.ScopeSelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScopedResourceSelectorRequirement)(nil), (*core.ScopedResourceSelectorRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScopedResourceSelectorRequirement_To_core_ScopedResourceSelectorRequirement(a.(*v1.ScopedResourceSelectorRequirement), b.(*core.ScopedResourceSelectorRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ScopedResourceSelectorRequirement)(nil), (*v1.ScopedResourceSelectorRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ScopedResourceSelectorRequirement_To_v1_ScopedResourceSelectorRequirement(a.(*core.ScopedResourceSelectorRequirement), b.(*v1.ScopedResourceSelectorRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Secret)(nil), (*core.Secret)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Secret_To_core_Secret(a.(*v1.Secret), b.(*core.Secret), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Secret)(nil), (*v1.Secret)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Secret_To_v1_Secret(a.(*core.Secret), b.(*v1.Secret), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretEnvSource)(nil), (*core.SecretEnvSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretEnvSource_To_core_SecretEnvSource(a.(*v1.SecretEnvSource), b.(*core.SecretEnvSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretEnvSource)(nil), (*v1.SecretEnvSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretEnvSource_To_v1_SecretEnvSource(a.(*core.SecretEnvSource), b.(*v1.SecretEnvSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretKeySelector)(nil), (*core.SecretKeySelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretKeySelector_To_core_SecretKeySelector(a.(*v1.SecretKeySelector), b.(*core.SecretKeySelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretKeySelector)(nil), (*v1.SecretKeySelector)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretKeySelector_To_v1_SecretKeySelector(a.(*core.SecretKeySelector), b.(*v1.SecretKeySelector), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretList)(nil), (*core.SecretList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretList_To_core_SecretList(a.(*v1.SecretList), b.(*core.SecretList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretList)(nil), (*v1.SecretList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretList_To_v1_SecretList(a.(*core.SecretList), b.(*v1.SecretList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretProjection)(nil), (*core.SecretProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretProjection_To_core_SecretProjection(a.(*v1.SecretProjection), b.(*core.SecretProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretProjection)(nil), (*v1.SecretProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretProjection_To_v1_SecretProjection(a.(*core.SecretProjection), b.(*v1.SecretProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretReference)(nil), (*core.SecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretReference_To_core_SecretReference(a.(*v1.SecretReference), b.(*core.SecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretReference)(nil), (*v1.SecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretReference_To_v1_SecretReference(a.(*core.SecretReference), b.(*v1.SecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecretVolumeSource)(nil), (*core.SecretVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecretVolumeSource_To_core_SecretVolumeSource(a.(*v1.SecretVolumeSource), b.(*core.SecretVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecretVolumeSource)(nil), (*v1.SecretVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecretVolumeSource_To_v1_SecretVolumeSource(a.(*core.SecretVolumeSource), b.(*v1.SecretVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecurityContext)(nil), (*core.SecurityContext)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecurityContext_To_core_SecurityContext(a.(*v1.SecurityContext), b.(*core.SecurityContext), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SecurityContext)(nil), (*v1.SecurityContext)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SecurityContext_To_v1_SecurityContext(a.(*core.SecurityContext), b.(*v1.SecurityContext), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SerializedReference)(nil), (*core.SerializedReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SerializedReference_To_core_SerializedReference(a.(*v1.SerializedReference), b.(*core.SerializedReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SerializedReference)(nil), (*v1.SerializedReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SerializedReference_To_v1_SerializedReference(a.(*core.SerializedReference), b.(*v1.SerializedReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Service)(nil), (*core.Service)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Service_To_core_Service(a.(*v1.Service), b.(*core.Service), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Service)(nil), (*v1.Service)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Service_To_v1_Service(a.(*core.Service), b.(*v1.Service), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccount)(nil), (*core.ServiceAccount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccount_To_core_ServiceAccount(a.(*v1.ServiceAccount), b.(*core.ServiceAccount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceAccount)(nil), (*v1.ServiceAccount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceAccount_To_v1_ServiceAccount(a.(*core.ServiceAccount), b.(*v1.ServiceAccount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountList)(nil), (*core.ServiceAccountList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountList_To_core_ServiceAccountList(a.(*v1.ServiceAccountList), b.(*core.ServiceAccountList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceAccountList)(nil), (*v1.ServiceAccountList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceAccountList_To_v1_ServiceAccountList(a.(*core.ServiceAccountList), b.(*v1.ServiceAccountList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountTokenProjection)(nil), (*core.ServiceAccountTokenProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountTokenProjection_To_core_ServiceAccountTokenProjection(a.(*v1.ServiceAccountTokenProjection), b.(*core.ServiceAccountTokenProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceAccountTokenProjection)(nil), (*v1.ServiceAccountTokenProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceAccountTokenProjection_To_v1_ServiceAccountTokenProjection(a.(*core.ServiceAccountTokenProjection), b.(*v1.ServiceAccountTokenProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceList)(nil), (*core.ServiceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceList_To_core_ServiceList(a.(*v1.ServiceList), b.(*core.ServiceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceList)(nil), (*v1.ServiceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceList_To_v1_ServiceList(a.(*core.ServiceList), b.(*v1.ServiceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServicePort)(nil), (*core.ServicePort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServicePort_To_core_ServicePort(a.(*v1.ServicePort), b.(*core.ServicePort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServicePort)(nil), (*v1.ServicePort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServicePort_To_v1_ServicePort(a.(*core.ServicePort), b.(*v1.ServicePort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceProxyOptions)(nil), (*core.ServiceProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceProxyOptions_To_core_ServiceProxyOptions(a.(*v1.ServiceProxyOptions), b.(*core.ServiceProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceProxyOptions)(nil), (*v1.ServiceProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceProxyOptions_To_v1_ServiceProxyOptions(a.(*core.ServiceProxyOptions), b.(*v1.ServiceProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceSpec)(nil), (*core.ServiceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceSpec_To_core_ServiceSpec(a.(*v1.ServiceSpec), b.(*core.ServiceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceSpec)(nil), (*v1.ServiceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceSpec_To_v1_ServiceSpec(a.(*core.ServiceSpec), b.(*v1.ServiceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceStatus)(nil), (*core.ServiceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceStatus_To_core_ServiceStatus(a.(*v1.ServiceStatus), b.(*core.ServiceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.ServiceStatus)(nil), (*v1.ServiceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ServiceStatus_To_v1_ServiceStatus(a.(*core.ServiceStatus), b.(*v1.ServiceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SessionAffinityConfig)(nil), (*core.SessionAffinityConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SessionAffinityConfig_To_core_SessionAffinityConfig(a.(*v1.SessionAffinityConfig), b.(*core.SessionAffinityConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.SessionAffinityConfig)(nil), (*v1.SessionAffinityConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_SessionAffinityConfig_To_v1_SessionAffinityConfig(a.(*core.SessionAffinityConfig), b.(*v1.SessionAffinityConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StorageOSPersistentVolumeSource)(nil), (*core.StorageOSPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StorageOSPersistentVolumeSource_To_core_StorageOSPersistentVolumeSource(a.(*v1.StorageOSPersistentVolumeSource), b.(*core.StorageOSPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.StorageOSPersistentVolumeSource)(nil), (*v1.StorageOSPersistentVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_StorageOSPersistentVolumeSource_To_v1_StorageOSPersistentVolumeSource(a.(*core.StorageOSPersistentVolumeSource), b.(*v1.StorageOSPersistentVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.StorageOSVolumeSource)(nil), (*core.StorageOSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_StorageOSVolumeSource_To_core_StorageOSVolumeSource(a.(*v1.StorageOSVolumeSource), b.(*core.StorageOSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.StorageOSVolumeSource)(nil), (*v1.StorageOSVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_StorageOSVolumeSource_To_v1_StorageOSVolumeSource(a.(*core.StorageOSVolumeSource), b.(*v1.StorageOSVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Sysctl)(nil), (*core.Sysctl)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Sysctl_To_core_Sysctl(a.(*v1.Sysctl), b.(*core.Sysctl), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Sysctl)(nil), (*v1.Sysctl)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Sysctl_To_v1_Sysctl(a.(*core.Sysctl), b.(*v1.Sysctl), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TCPSocketAction)(nil), (*core.TCPSocketAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TCPSocketAction_To_core_TCPSocketAction(a.(*v1.TCPSocketAction), b.(*core.TCPSocketAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.TCPSocketAction)(nil), (*v1.TCPSocketAction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_TCPSocketAction_To_v1_TCPSocketAction(a.(*core.TCPSocketAction), b.(*v1.TCPSocketAction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Taint)(nil), (*core.Taint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Taint_To_core_Taint(a.(*v1.Taint), b.(*core.Taint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Taint)(nil), (*v1.Taint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Taint_To_v1_Taint(a.(*core.Taint), b.(*v1.Taint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Toleration)(nil), (*core.Toleration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Toleration_To_core_Toleration(a.(*v1.Toleration), b.(*core.Toleration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Toleration)(nil), (*v1.Toleration)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Toleration_To_v1_Toleration(a.(*core.Toleration), b.(*v1.Toleration), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TopologySelectorLabelRequirement)(nil), (*core.TopologySelectorLabelRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TopologySelectorLabelRequirement_To_core_TopologySelectorLabelRequirement(a.(*v1.TopologySelectorLabelRequirement), b.(*core.TopologySelectorLabelRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.TopologySelectorLabelRequirement)(nil), (*v1.TopologySelectorLabelRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_TopologySelectorLabelRequirement_To_v1_TopologySelectorLabelRequirement(a.(*core.TopologySelectorLabelRequirement), b.(*v1.TopologySelectorLabelRequirement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TopologySelectorTerm)(nil), (*core.TopologySelectorTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TopologySelectorTerm_To_core_TopologySelectorTerm(a.(*v1.TopologySelectorTerm), b.(*core.TopologySelectorTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.TopologySelectorTerm)(nil), (*v1.TopologySelectorTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_TopologySelectorTerm_To_v1_TopologySelectorTerm(a.(*core.TopologySelectorTerm), b.(*v1.TopologySelectorTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TypedLocalObjectReference)(nil), (*core.TypedLocalObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TypedLocalObjectReference_To_core_TypedLocalObjectReference(a.(*v1.TypedLocalObjectReference), b.(*core.TypedLocalObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.TypedLocalObjectReference)(nil), (*v1.TypedLocalObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_TypedLocalObjectReference_To_v1_TypedLocalObjectReference(a.(*core.TypedLocalObjectReference), b.(*v1.TypedLocalObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Volume)(nil), (*core.Volume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Volume_To_core_Volume(a.(*v1.Volume), b.(*core.Volume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.Volume)(nil), (*v1.Volume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Volume_To_v1_Volume(a.(*core.Volume), b.(*v1.Volume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VolumeDevice)(nil), (*core.VolumeDevice)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VolumeDevice_To_core_VolumeDevice(a.(*v1.VolumeDevice), b.(*core.VolumeDevice), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VolumeDevice)(nil), (*v1.VolumeDevice)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VolumeDevice_To_v1_VolumeDevice(a.(*core.VolumeDevice), b.(*v1.VolumeDevice), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VolumeMount)(nil), (*core.VolumeMount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VolumeMount_To_core_VolumeMount(a.(*v1.VolumeMount), b.(*core.VolumeMount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VolumeMount)(nil), (*v1.VolumeMount)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VolumeMount_To_v1_VolumeMount(a.(*core.VolumeMount), b.(*v1.VolumeMount), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VolumeNodeAffinity)(nil), (*core.VolumeNodeAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VolumeNodeAffinity_To_core_VolumeNodeAffinity(a.(*v1.VolumeNodeAffinity), b.(*core.VolumeNodeAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VolumeNodeAffinity)(nil), (*v1.VolumeNodeAffinity)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VolumeNodeAffinity_To_v1_VolumeNodeAffinity(a.(*core.VolumeNodeAffinity), b.(*v1.VolumeNodeAffinity), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VolumeProjection)(nil), (*core.VolumeProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VolumeProjection_To_core_VolumeProjection(a.(*v1.VolumeProjection), b.(*core.VolumeProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VolumeProjection)(nil), (*v1.VolumeProjection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VolumeProjection_To_v1_VolumeProjection(a.(*core.VolumeProjection), b.(*v1.VolumeProjection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VolumeSource)(nil), (*core.VolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VolumeSource_To_core_VolumeSource(a.(*v1.VolumeSource), b.(*core.VolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VolumeSource)(nil), (*v1.VolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VolumeSource_To_v1_VolumeSource(a.(*core.VolumeSource), b.(*v1.VolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.VsphereVirtualDiskVolumeSource)(nil), (*core.VsphereVirtualDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_VsphereVirtualDiskVolumeSource_To_core_VsphereVirtualDiskVolumeSource(a.(*v1.VsphereVirtualDiskVolumeSource), b.(*core.VsphereVirtualDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.VsphereVirtualDiskVolumeSource)(nil), (*v1.VsphereVirtualDiskVolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_VsphereVirtualDiskVolumeSource_To_v1_VsphereVirtualDiskVolumeSource(a.(*core.VsphereVirtualDiskVolumeSource), b.(*v1.VsphereVirtualDiskVolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.WeightedPodAffinityTerm)(nil), (*core.WeightedPodAffinityTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_WeightedPodAffinityTerm_To_core_WeightedPodAffinityTerm(a.(*v1.WeightedPodAffinityTerm), b.(*core.WeightedPodAffinityTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.WeightedPodAffinityTerm)(nil), (*v1.WeightedPodAffinityTerm)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_WeightedPodAffinityTerm_To_v1_WeightedPodAffinityTerm(a.(*core.WeightedPodAffinityTerm), b.(*v1.WeightedPodAffinityTerm), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.WindowsSecurityContextOptions)(nil), (*core.WindowsSecurityContextOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_WindowsSecurityContextOptions_To_core_WindowsSecurityContextOptions(a.(*v1.WindowsSecurityContextOptions), b.(*core.WindowsSecurityContextOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*core.WindowsSecurityContextOptions)(nil), (*v1.WindowsSecurityContextOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_WindowsSecurityContextOptions_To_v1_WindowsSecurityContextOptions(a.(*core.WindowsSecurityContextOptions), b.(*v1.WindowsSecurityContextOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1.ReplicationControllerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1_ReplicationControllerSpec(a.(*apps.ReplicaSetSpec), b.(*v1.ReplicationControllerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.ReplicaSetStatus)(nil), (*v1.ReplicationControllerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetStatus_To_v1_ReplicationControllerStatus(a.(*apps.ReplicaSetStatus), b.(*v1.ReplicationControllerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.ReplicaSet)(nil), (*v1.ReplicationController)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSet_To_v1_ReplicationController(a.(*apps.ReplicaSet), b.(*v1.ReplicationController), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*core.PodSpec)(nil), (*v1.PodSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodSpec_To_v1_PodSpec(a.(*core.PodSpec), b.(*v1.PodSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*core.PodTemplateSpec)(nil), (*v1.PodTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_PodTemplateSpec_To_v1_PodTemplateSpec(a.(*core.PodTemplateSpec), b.(*v1.PodTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*core.Pod)(nil), (*v1.Pod)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_Pod_To_v1_Pod(a.(*core.Pod), b.(*v1.Pod), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*core.ReplicationControllerSpec)(nil), (*v1.ReplicationControllerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_core_ReplicationControllerSpec_To_v1_ReplicationControllerSpec(a.(*core.ReplicationControllerSpec), b.(*v1.ReplicationControllerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.PodSpec)(nil), (*core.PodSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSpec_To_core_PodSpec(a.(*v1.PodSpec), b.(*core.PodSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.PodTemplateSpec)(nil), (*core.PodTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodTemplateSpec_To_core_PodTemplateSpec(a.(*v1.PodTemplateSpec), b.(*core.PodTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.Pod)(nil), (*core.Pod)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Pod_To_core_Pod(a.(*v1.Pod), b.(*core.Pod), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ReplicationControllerSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerSpec_To_apps_ReplicaSetSpec(a.(*v1.ReplicationControllerSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ReplicationControllerSpec)(nil), (*core.ReplicationControllerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerSpec_To_core_ReplicationControllerSpec(a.(*v1.ReplicationControllerSpec), b.(*core.ReplicationControllerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ReplicationControllerStatus)(nil), (*apps.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationControllerStatus_To_apps_ReplicaSetStatus(a.(*v1.ReplicationControllerStatus), b.(*apps.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ReplicationController)(nil), (*apps.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ReplicationController_To_apps_ReplicaSet(a.(*v1.ReplicationController), b.(*apps.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ResourceList)(nil), (*core.ResourceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceList_To_core_ResourceList(a.(*v1.ResourceList), b.(*core.ResourceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.Secret)(nil), (*core.Secret)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Secret_To_core_Secret(a.(*v1.Secret), b.(*core.Secret), scope)
	}); err != nil {
		return err
	}
	return nil
}
