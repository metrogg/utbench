func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v2beta1.CrossVersionObjectReference)(nil), (*autoscaling.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_CrossVersionObjectReference_To_autoscaling_CrossVersionObjectReference(a.(*v2beta1.CrossVersionObjectReference), b.(*autoscaling.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.CrossVersionObjectReference)(nil), (*v2beta1.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_CrossVersionObjectReference_To_v2beta1_CrossVersionObjectReference(a.(*autoscaling.CrossVersionObjectReference), b.(*v2beta1.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ExternalMetricSource)(nil), (*autoscaling.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ExternalMetricSource_To_autoscaling_ExternalMetricSource(a.(*v2beta1.ExternalMetricSource), b.(*autoscaling.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ExternalMetricSource)(nil), (*v2beta1.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricSource_To_v2beta1_ExternalMetricSource(a.(*autoscaling.ExternalMetricSource), b.(*v2beta1.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ExternalMetricStatus)(nil), (*autoscaling.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ExternalMetricStatus_To_autoscaling_ExternalMetricStatus(a.(*v2beta1.ExternalMetricStatus), b.(*autoscaling.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ExternalMetricStatus)(nil), (*v2beta1.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricStatus_To_v2beta1_ExternalMetricStatus(a.(*autoscaling.ExternalMetricStatus), b.(*v2beta1.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.HorizontalPodAutoscaler)(nil), (*autoscaling.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler(a.(*v2beta1.HorizontalPodAutoscaler), b.(*autoscaling.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscaler)(nil), (*v2beta1.HorizontalPodAutoscaler)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscaler_To_v2beta1_HorizontalPodAutoscaler(a.(*autoscaling.HorizontalPodAutoscaler), b.(*v2beta1.HorizontalPodAutoscaler), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.HorizontalPodAutoscalerCondition)(nil), (*autoscaling.HorizontalPodAutoscalerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_HorizontalPodAutoscalerCondition_To_autoscaling_HorizontalPodAutoscalerCondition(a.(*v2beta1.HorizontalPodAutoscalerCondition), b.(*autoscaling.HorizontalPodAutoscalerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerCondition)(nil), (*v2beta1.HorizontalPodAutoscalerCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerCondition_To_v2beta1_HorizontalPodAutoscalerCondition(a.(*autoscaling.HorizontalPodAutoscalerCondition), b.(*v2beta1.HorizontalPodAutoscalerCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.HorizontalPodAutoscalerList)(nil), (*autoscaling.HorizontalPodAutoscalerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList(a.(*v2beta1.HorizontalPodAutoscalerList), b.(*autoscaling.HorizontalPodAutoscalerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerList)(nil), (*v2beta1.HorizontalPodAutoscalerList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerList_To_v2beta1_HorizontalPodAutoscalerList(a.(*autoscaling.HorizontalPodAutoscalerList), b.(*v2beta1.HorizontalPodAutoscalerList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.HorizontalPodAutoscalerSpec)(nil), (*autoscaling.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_HorizontalPodAutoscalerSpec_To_autoscaling_HorizontalPodAutoscalerSpec(a.(*v2beta1.HorizontalPodAutoscalerSpec), b.(*autoscaling.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerSpec)(nil), (*v2beta1.HorizontalPodAutoscalerSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerSpec_To_v2beta1_HorizontalPodAutoscalerSpec(a.(*autoscaling.HorizontalPodAutoscalerSpec), b.(*v2beta1.HorizontalPodAutoscalerSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.HorizontalPodAutoscalerStatus)(nil), (*autoscaling.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_HorizontalPodAutoscalerStatus_To_autoscaling_HorizontalPodAutoscalerStatus(a.(*v2beta1.HorizontalPodAutoscalerStatus), b.(*autoscaling.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.HorizontalPodAutoscalerStatus)(nil), (*v2beta1.HorizontalPodAutoscalerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_HorizontalPodAutoscalerStatus_To_v2beta1_HorizontalPodAutoscalerStatus(a.(*autoscaling.HorizontalPodAutoscalerStatus), b.(*v2beta1.HorizontalPodAutoscalerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.MetricSpec)(nil), (*autoscaling.MetricSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_MetricSpec_To_autoscaling_MetricSpec(a.(*v2beta1.MetricSpec), b.(*autoscaling.MetricSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.MetricSpec)(nil), (*v2beta1.MetricSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricSpec_To_v2beta1_MetricSpec(a.(*autoscaling.MetricSpec), b.(*v2beta1.MetricSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.MetricStatus)(nil), (*autoscaling.MetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_MetricStatus_To_autoscaling_MetricStatus(a.(*v2beta1.MetricStatus), b.(*autoscaling.MetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.MetricStatus)(nil), (*v2beta1.MetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricStatus_To_v2beta1_MetricStatus(a.(*autoscaling.MetricStatus), b.(*v2beta1.MetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ObjectMetricSource)(nil), (*autoscaling.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ObjectMetricSource_To_autoscaling_ObjectMetricSource(a.(*v2beta1.ObjectMetricSource), b.(*autoscaling.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ObjectMetricSource)(nil), (*v2beta1.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricSource_To_v2beta1_ObjectMetricSource(a.(*autoscaling.ObjectMetricSource), b.(*v2beta1.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ObjectMetricStatus)(nil), (*autoscaling.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ObjectMetricStatus_To_autoscaling_ObjectMetricStatus(a.(*v2beta1.ObjectMetricStatus), b.(*autoscaling.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ObjectMetricStatus)(nil), (*v2beta1.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricStatus_To_v2beta1_ObjectMetricStatus(a.(*autoscaling.ObjectMetricStatus), b.(*v2beta1.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.PodsMetricSource)(nil), (*autoscaling.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_PodsMetricSource_To_autoscaling_PodsMetricSource(a.(*v2beta1.PodsMetricSource), b.(*autoscaling.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.PodsMetricSource)(nil), (*v2beta1.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricSource_To_v2beta1_PodsMetricSource(a.(*autoscaling.PodsMetricSource), b.(*v2beta1.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.PodsMetricStatus)(nil), (*autoscaling.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_PodsMetricStatus_To_autoscaling_PodsMetricStatus(a.(*v2beta1.PodsMetricStatus), b.(*autoscaling.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.PodsMetricStatus)(nil), (*v2beta1.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricStatus_To_v2beta1_PodsMetricStatus(a.(*autoscaling.PodsMetricStatus), b.(*v2beta1.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ResourceMetricSource)(nil), (*autoscaling.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ResourceMetricSource_To_autoscaling_ResourceMetricSource(a.(*v2beta1.ResourceMetricSource), b.(*autoscaling.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ResourceMetricSource)(nil), (*v2beta1.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricSource_To_v2beta1_ResourceMetricSource(a.(*autoscaling.ResourceMetricSource), b.(*v2beta1.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v2beta1.ResourceMetricStatus)(nil), (*autoscaling.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ResourceMetricStatus_To_autoscaling_ResourceMetricStatus(a.(*v2beta1.ResourceMetricStatus), b.(*autoscaling.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*autoscaling.ResourceMetricStatus)(nil), (*v2beta1.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricStatus_To_v2beta1_ResourceMetricStatus(a.(*autoscaling.ResourceMetricStatus), b.(*v2beta1.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ExternalMetricSource)(nil), (*v2beta1.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricSource_To_v2beta1_ExternalMetricSource(a.(*autoscaling.ExternalMetricSource), b.(*v2beta1.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ExternalMetricStatus)(nil), (*v2beta1.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ExternalMetricStatus_To_v2beta1_ExternalMetricStatus(a.(*autoscaling.ExternalMetricStatus), b.(*v2beta1.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.MetricTarget)(nil), (*v2beta1.CrossVersionObjectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_MetricTarget_To_v2beta1_CrossVersionObjectReference(a.(*autoscaling.MetricTarget), b.(*v2beta1.CrossVersionObjectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ObjectMetricSource)(nil), (*v2beta1.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricSource_To_v2beta1_ObjectMetricSource(a.(*autoscaling.ObjectMetricSource), b.(*v2beta1.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ObjectMetricStatus)(nil), (*v2beta1.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ObjectMetricStatus_To_v2beta1_ObjectMetricStatus(a.(*autoscaling.ObjectMetricStatus), b.(*v2beta1.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.PodsMetricSource)(nil), (*v2beta1.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricSource_To_v2beta1_PodsMetricSource(a.(*autoscaling.PodsMetricSource), b.(*v2beta1.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.PodsMetricStatus)(nil), (*v2beta1.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_PodsMetricStatus_To_v2beta1_PodsMetricStatus(a.(*autoscaling.PodsMetricStatus), b.(*v2beta1.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ResourceMetricSource)(nil), (*v2beta1.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricSource_To_v2beta1_ResourceMetricSource(a.(*autoscaling.ResourceMetricSource), b.(*v2beta1.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ResourceMetricStatus)(nil), (*v2beta1.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ResourceMetricStatus_To_v2beta1_ResourceMetricStatus(a.(*autoscaling.ResourceMetricStatus), b.(*v2beta1.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.CrossVersionObjectReference)(nil), (*autoscaling.MetricTarget)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_CrossVersionObjectReference_To_autoscaling_MetricTarget(a.(*v2beta1.CrossVersionObjectReference), b.(*autoscaling.MetricTarget), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ExternalMetricSource)(nil), (*autoscaling.ExternalMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ExternalMetricSource_To_autoscaling_ExternalMetricSource(a.(*v2beta1.ExternalMetricSource), b.(*autoscaling.ExternalMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ExternalMetricStatus)(nil), (*autoscaling.ExternalMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ExternalMetricStatus_To_autoscaling_ExternalMetricStatus(a.(*v2beta1.ExternalMetricStatus), b.(*autoscaling.ExternalMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ObjectMetricSource)(nil), (*autoscaling.ObjectMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ObjectMetricSource_To_autoscaling_ObjectMetricSource(a.(*v2beta1.ObjectMetricSource), b.(*autoscaling.ObjectMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ObjectMetricStatus)(nil), (*autoscaling.ObjectMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ObjectMetricStatus_To_autoscaling_ObjectMetricStatus(a.(*v2beta1.ObjectMetricStatus), b.(*autoscaling.ObjectMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.PodsMetricSource)(nil), (*autoscaling.PodsMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_PodsMetricSource_To_autoscaling_PodsMetricSource(a.(*v2beta1.PodsMetricSource), b.(*autoscaling.PodsMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.PodsMetricStatus)(nil), (*autoscaling.PodsMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_PodsMetricStatus_To_autoscaling_PodsMetricStatus(a.(*v2beta1.PodsMetricStatus), b.(*autoscaling.PodsMetricStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ResourceMetricSource)(nil), (*autoscaling.ResourceMetricSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ResourceMetricSource_To_autoscaling_ResourceMetricSource(a.(*v2beta1.ResourceMetricSource), b.(*autoscaling.ResourceMetricSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v2beta1.ResourceMetricStatus)(nil), (*autoscaling.ResourceMetricStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v2beta1_ResourceMetricStatus_To_autoscaling_ResourceMetricStatus(a.(*v2beta1.ResourceMetricStatus), b.(*autoscaling.ResourceMetricStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}
