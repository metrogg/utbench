func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.CrossVersionObjectReference)(nil), (*autoscaling.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CrossVersionObjectReference_To_autoscaling_CrossVersionObjectReference(a.(*v1.CrossVersionObjectReference), b.(*autoscaling.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.CrossVersionObjectReference)(nil), (*v1.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_CrossVersionObjectReference_To_v1_CrossVersionObjectReference(a.(*autoscaling.CrossVersionObjectReference), b.(*v1.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ExternalMetricSource)(nil), (*autoscaling.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExternalMetricSource_To_autoscaling_ExternalMetricSource(a.(*v1.ExternalMetricSource), b.(*autoscaling.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ExternalMetricSource)(nil), (*v1.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricSource_To_v1_ExternalMetricSource(a.(*autoscaling.ExternalMetricSource), b.(*v1.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ExternalMetricStatus)(nil), (*autoscaling.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExternalMetricStatus_To_autoscaling_ExternalMetricStatus(a.(*v1.ExternalMetricStatus), b.(*autoscaling.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ExternalMetricStatus)(nil), (*v1.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricStatus_To_v1_ExternalMetricStatus(a.(*autoscaling.ExternalMetricStatus), b.(*v1.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HorizontalPodAutoscaler)(nil), (*autoscaling.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler(a.(*v1.HorizontalPodAutoscaler), b.(*autoscaling.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscaler)(nil), (*v1.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscaler_To_v1_HorizontalPodAutoscaler(a.(*autoscaling.HorizontalPodAutoscaler), b.(*v1.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HorizontalPodAutoscalerCondition)(nil), (*autoscaling.HorizontalPodAutoscalerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerCondition_To_autoscaling_HorizontalPodAutoscalerCondition(a.(*v1.HorizontalPodAutoscalerCondition), b.(*autoscaling.HorizontalPodAutoscalerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerCondition)(nil), (*v1.HorizontalPodAutoscalerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerCondition_To_v1_HorizontalPodAutoscalerCondition(a.(*autoscaling.HorizontalPodAutoscalerCondition), b.(*v1.HorizontalPodAutoscalerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HorizontalPodAutoscalerList)(nil), (*autoscaling.HorizontalPodAutoscalerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList(a.(*v1.HorizontalPodAutoscalerList), b.(*autoscaling.HorizontalPodAutoscalerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerList)(nil), (*v1.HorizontalPodAutoscalerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerList_To_v1_HorizontalPodAutoscalerList(a.(*autoscaling.HorizontalPodAutoscalerList), b.(*v1.HorizontalPodAutoscalerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HorizontalPodAutoscalerSpec)(nil), (*autoscaling.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerSpec_To_autoscaling_HorizontalPodAutoscalerSpec(a.(*v1.HorizontalPodAutoscalerSpec), b.(*autoscaling.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerSpec)(nil), (*v1.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerSpec_To_v1_HorizontalPodAutoscalerSpec(a.(*autoscaling.HorizontalPodAutoscalerSpec), b.(*v1.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.HorizontalPodAutoscalerStatus)(nil), (*autoscaling.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerStatus_To_autoscaling_HorizontalPodAutoscalerStatus(a.(*v1.HorizontalPodAutoscalerStatus), b.(*autoscaling.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerStatus)(nil), (*v1.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerStatus_To_v1_HorizontalPodAutoscalerStatus(a.(*autoscaling.HorizontalPodAutoscalerStatus), b.(*v1.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MetricSpec)(nil), (*autoscaling.MetricSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MetricSpec_To_autoscaling_MetricSpec(a.(*v1.MetricSpec), b.(*autoscaling.MetricSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.MetricSpec)(nil), (*v1.MetricSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricSpec_To_v1_MetricSpec(a.(*autoscaling.MetricSpec), b.(*v1.MetricSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.MetricStatus)(nil), (*autoscaling.MetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_MetricStatus_To_autoscaling_MetricStatus(a.(*v1.MetricStatus), b.(*autoscaling.MetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.MetricStatus)(nil), (*v1.MetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricStatus_To_v1_MetricStatus(a.(*autoscaling.MetricStatus), b.(*v1.MetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ObjectMetricSource)(nil), (*autoscaling.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectMetricSource_To_autoscaling_ObjectMetricSource(a.(*v1.ObjectMetricSource), b.(*autoscaling.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ObjectMetricSource)(nil), (*v1.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricSource_To_v1_ObjectMetricSource(a.(*autoscaling.ObjectMetricSource), b.(*v1.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ObjectMetricStatus)(nil), (*autoscaling.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectMetricStatus_To_autoscaling_ObjectMetricStatus(a.(*v1.ObjectMetricStatus), b.(*autoscaling.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ObjectMetricStatus)(nil), (*v1.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricStatus_To_v1_ObjectMetricStatus(a.(*autoscaling.ObjectMetricStatus), b.(*v1.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodsMetricSource)(nil), (*autoscaling.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodsMetricSource_To_autoscaling_PodsMetricSource(a.(*v1.PodsMetricSource), b.(*autoscaling.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.PodsMetricSource)(nil), (*v1.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricSource_To_v1_PodsMetricSource(a.(*autoscaling.PodsMetricSource), b.(*v1.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodsMetricStatus)(nil), (*autoscaling.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodsMetricStatus_To_autoscaling_PodsMetricStatus(a.(*v1.PodsMetricStatus), b.(*autoscaling.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.PodsMetricStatus)(nil), (*v1.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricStatus_To_v1_PodsMetricStatus(a.(*autoscaling.PodsMetricStatus), b.(*v1.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceMetricSource)(nil), (*autoscaling.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceMetricSource_To_autoscaling_ResourceMetricSource(a.(*v1.ResourceMetricSource), b.(*autoscaling.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ResourceMetricSource)(nil), (*v1.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricSource_To_v1_ResourceMetricSource(a.(*autoscaling.ResourceMetricSource), b.(*v1.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceMetricStatus)(nil), (*autoscaling.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceMetricStatus_To_autoscaling_ResourceMetricStatus(a.(*v1.ResourceMetricStatus), b.(*autoscaling.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ResourceMetricStatus)(nil), (*v1.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricStatus_To_v1_ResourceMetricStatus(a.(*autoscaling.ResourceMetricStatus), b.(*v1.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Scale)(nil), (*autoscaling.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Scale_To_autoscaling_Scale(a.(*v1.Scale), b.(*autoscaling.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.Scale)(nil), (*v1.Scale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_Scale_To_v1_Scale(a.(*autoscaling.Scale), b.(*v1.Scale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleSpec)(nil), (*autoscaling.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleSpec_To_autoscaling_ScaleSpec(a.(*v1.ScaleSpec), b.(*autoscaling.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleSpec)(nil), (*v1.ScaleSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleSpec_To_v1_ScaleSpec(a.(*autoscaling.ScaleSpec), b.(*v1.ScaleSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScaleStatus)(nil), (*autoscaling.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScaleStatus_To_autoscaling_ScaleStatus(a.(*v1.ScaleStatus), b.(*autoscaling.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ExternalMetricSource)(nil), (*v1.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricSource_To_v1_ExternalMetricSource(a.(*autoscaling.ExternalMetricSource), b.(*v1.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ExternalMetricStatus)(nil), (*v1.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricStatus_To_v1_ExternalMetricStatus(a.(*autoscaling.ExternalMetricStatus), b.(*v1.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.HorizontalPodAutoscalerSpec)(nil), (*v1.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerSpec_To_v1_HorizontalPodAutoscalerSpec(a.(*autoscaling.HorizontalPodAutoscalerSpec), b.(*v1.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.HorizontalPodAutoscalerStatus)(nil), (*v1.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerStatus_To_v1_HorizontalPodAutoscalerStatus(a.(*autoscaling.HorizontalPodAutoscalerStatus), b.(*v1.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.HorizontalPodAutoscaler)(nil), (*v1.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscaler_To_v1_HorizontalPodAutoscaler(a.(*autoscaling.HorizontalPodAutoscaler), b.(*v1.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.MetricTarget)(nil), (*v1.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricTarget_To_v1_CrossVersionObjectReference(a.(*autoscaling.MetricTarget), b.(*v1.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ObjectMetricSource)(nil), (*v1.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricSource_To_v1_ObjectMetricSource(a.(*autoscaling.ObjectMetricSource), b.(*v1.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ObjectMetricStatus)(nil), (*v1.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricStatus_To_v1_ObjectMetricStatus(a.(*autoscaling.ObjectMetricStatus), b.(*v1.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.PodsMetricSource)(nil), (*v1.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricSource_To_v1_PodsMetricSource(a.(*autoscaling.PodsMetricSource), b.(*v1.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.PodsMetricStatus)(nil), (*v1.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricStatus_To_v1_PodsMetricStatus(a.(*autoscaling.PodsMetricStatus), b.(*v1.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ResourceMetricSource)(nil), (*v1.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricSource_To_v1_ResourceMetricSource(a.(*autoscaling.ResourceMetricSource), b.(*v1.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ResourceMetricStatus)(nil), (*v1.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricStatus_To_v1_ResourceMetricStatus(a.(*autoscaling.ResourceMetricStatus), b.(*v1.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.CrossVersionObjectReference)(nil), (*autoscaling.MetricTarget)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_CrossVersionObjectReference_To_autoscaling_MetricTarget(a.(*v1.CrossVersionObjectReference), b.(*autoscaling.MetricTarget), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ExternalMetricSource)(nil), (*autoscaling.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExternalMetricSource_To_autoscaling_ExternalMetricSource(a.(*v1.ExternalMetricSource), b.(*autoscaling.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ExternalMetricStatus)(nil), (*autoscaling.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ExternalMetricStatus_To_autoscaling_ExternalMetricStatus(a.(*v1.ExternalMetricStatus), b.(*autoscaling.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.HorizontalPodAutoscalerSpec)(nil), (*autoscaling.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerSpec_To_autoscaling_HorizontalPodAutoscalerSpec(a.(*v1.HorizontalPodAutoscalerSpec), b.(*autoscaling.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.HorizontalPodAutoscalerStatus)(nil), (*autoscaling.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscalerStatus_To_autoscaling_HorizontalPodAutoscalerStatus(a.(*v1.HorizontalPodAutoscalerStatus), b.(*autoscaling.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.HorizontalPodAutoscaler)(nil), (*autoscaling.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler(a.(*v1.HorizontalPodAutoscaler), b.(*autoscaling.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ObjectMetricSource)(nil), (*autoscaling.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectMetricSource_To_autoscaling_ObjectMetricSource(a.(*v1.ObjectMetricSource), b.(*autoscaling.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ObjectMetricStatus)(nil), (*autoscaling.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ObjectMetricStatus_To_autoscaling_ObjectMetricStatus(a.(*v1.ObjectMetricStatus), b.(*autoscaling.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.PodsMetricSource)(nil), (*autoscaling.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodsMetricSource_To_autoscaling_PodsMetricSource(a.(*v1.PodsMetricSource), b.(*autoscaling.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.PodsMetricStatus)(nil), (*autoscaling.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodsMetricStatus_To_autoscaling_PodsMetricStatus(a.(*v1.PodsMetricStatus), b.(*autoscaling.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ResourceMetricSource)(nil), (*autoscaling.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceMetricSource_To_autoscaling_ResourceMetricSource(a.(*v1.ResourceMetricSource), b.(*autoscaling.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ResourceMetricStatus)(nil), (*autoscaling.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceMetricStatus_To_autoscaling_ResourceMetricStatus(a.(*v1.ResourceMetricStatus), b.(*autoscaling.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}
