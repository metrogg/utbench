func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.AllowedCSIDriver)(nil), (*policy.AllowedCSIDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_AllowedCSIDriver_To_policy_AllowedCSIDriver(a.(*v1beta1.AllowedCSIDriver), b.(*policy.AllowedCSIDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.AllowedCSIDriver)(nil), (*v1beta1.AllowedCSIDriver)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_AllowedCSIDriver_To_v1beta1_AllowedCSIDriver(a.(*policy.AllowedCSIDriver), b.(*v1beta1.AllowedCSIDriver), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.AllowedFlexVolume)(nil), (*policy.AllowedFlexVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_AllowedFlexVolume_To_policy_AllowedFlexVolume(a.(*v1beta1.AllowedFlexVolume), b.(*policy.AllowedFlexVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.AllowedFlexVolume)(nil), (*v1beta1.AllowedFlexVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_AllowedFlexVolume_To_v1beta1_AllowedFlexVolume(a.(*policy.AllowedFlexVolume), b.(*v1beta1.AllowedFlexVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.AllowedHostPath)(nil), (*policy.AllowedHostPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_AllowedHostPath_To_policy_AllowedHostPath(a.(*v1beta1.AllowedHostPath), b.(*policy.AllowedHostPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.AllowedHostPath)(nil), (*v1beta1.AllowedHostPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_AllowedHostPath_To_v1beta1_AllowedHostPath(a.(*policy.AllowedHostPath), b.(*v1beta1.AllowedHostPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSet)(nil), (*apps.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSet_To_apps_DaemonSet(a.(*v1beta1.DaemonSet), b.(*apps.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSet)(nil), (*v1beta1.DaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSet_To_v1beta1_DaemonSet(a.(*apps.DaemonSet), b.(*v1beta1.DaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSetCondition)(nil), (*apps.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSetCondition_To_apps_DaemonSetCondition(a.(*v1beta1.DaemonSetCondition), b.(*apps.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetCondition)(nil), (*v1beta1.DaemonSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetCondition_To_v1beta1_DaemonSetCondition(a.(*apps.DaemonSetCondition), b.(*v1beta1.DaemonSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSetList)(nil), (*apps.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSetList_To_apps_DaemonSetList(a.(*v1beta1.DaemonSetList), b.(*apps.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetList)(nil), (*v1beta1.DaemonSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetList_To_v1beta1_DaemonSetList(a.(*apps.DaemonSetList), b.(*v1beta1.DaemonSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSetSpec)(nil), (*apps.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSetSpec_To_apps_DaemonSetSpec(a.(*v1beta1.DaemonSetSpec), b.(*apps.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetSpec)(nil), (*v1beta1.DaemonSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetSpec_To_v1beta1_DaemonSetSpec(a.(*apps.DaemonSetSpec), b.(*v1beta1.DaemonSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSetStatus)(nil), (*apps.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSetStatus_To_apps_DaemonSetStatus(a.(*v1beta1.DaemonSetStatus), b.(*apps.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetStatus)(nil), (*v1beta1.DaemonSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetStatus_To_v1beta1_DaemonSetStatus(a.(*apps.DaemonSetStatus), b.(*v1beta1.DaemonSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.DaemonSetUpdateStrategy)(nil), (*apps.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_DaemonSetUpdateStrategy_To_apps_DaemonSetUpdateStrategy(a.(*v1beta1.DaemonSetUpdateStrategy), b.(*apps.DaemonSetUpdateStrategy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.DaemonSetUpdateStrategy)(nil), (*v1beta1.DaemonSetUpdateStrategy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_DaemonSetUpdateStrategy_To_v1beta1_DaemonSetUpdateStrategy(a.(*apps.DaemonSetUpdateStrategy), b.(*v1beta1.DaemonSetUpdateStrategy), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1beta1.FSGroupStrategyOptions)(nil), (*policy.FSGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_FSGroupStrategyOptions_To_policy_FSGroupStrategyOptions(a.(*v1beta1.FSGroupStrategyOptions), b.(*policy.FSGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.FSGroupStrategyOptions)(nil), (*v1beta1.FSGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_FSGroupStrategyOptions_To_v1beta1_FSGroupStrategyOptions(a.(*policy.FSGroupStrategyOptions), b.(*v1beta1.FSGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.HTTPIngressPath)(nil), (*networking.HTTPIngressPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_HTTPIngressPath_To_networking_HTTPIngressPath(a.(*v1beta1.HTTPIngressPath), b.(*networking.HTTPIngressPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.HTTPIngressPath)(nil), (*v1beta1.HTTPIngressPath)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_HTTPIngressPath_To_v1beta1_HTTPIngressPath(a.(*networking.HTTPIngressPath), b.(*v1beta1.HTTPIngressPath), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.HTTPIngressRuleValue)(nil), (*networking.HTTPIngressRuleValue)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_HTTPIngressRuleValue_To_networking_HTTPIngressRuleValue(a.(*v1beta1.HTTPIngressRuleValue), b.(*networking.HTTPIngressRuleValue), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.HTTPIngressRuleValue)(nil), (*v1beta1.HTTPIngressRuleValue)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_HTTPIngressRuleValue_To_v1beta1_HTTPIngressRuleValue(a.(*networking.HTTPIngressRuleValue), b.(*v1beta1.HTTPIngressRuleValue), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.HostPortRange)(nil), (*policy.HostPortRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_HostPortRange_To_policy_HostPortRange(a.(*v1beta1.HostPortRange), b.(*policy.HostPortRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.HostPortRange)(nil), (*v1beta1.HostPortRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_HostPortRange_To_v1beta1_HostPortRange(a.(*policy.HostPortRange), b.(*v1beta1.HostPortRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IDRange)(nil), (*policy.IDRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IDRange_To_policy_IDRange(a.(*v1beta1.IDRange), b.(*policy.IDRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.IDRange)(nil), (*v1beta1.IDRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_IDRange_To_v1beta1_IDRange(a.(*policy.IDRange), b.(*v1beta1.IDRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IPBlock)(nil), (*networking.IPBlock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IPBlock_To_networking_IPBlock(a.(*v1beta1.IPBlock), b.(*networking.IPBlock), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IPBlock)(nil), (*v1beta1.IPBlock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IPBlock_To_v1beta1_IPBlock(a.(*networking.IPBlock), b.(*v1beta1.IPBlock), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Ingress)(nil), (*networking.Ingress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Ingress_To_networking_Ingress(a.(*v1beta1.Ingress), b.(*networking.Ingress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.Ingress)(nil), (*v1beta1.Ingress)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_Ingress_To_v1beta1_Ingress(a.(*networking.Ingress), b.(*v1beta1.Ingress), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressBackend)(nil), (*networking.IngressBackend)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressBackend_To_networking_IngressBackend(a.(*v1beta1.IngressBackend), b.(*networking.IngressBackend), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressBackend)(nil), (*v1beta1.IngressBackend)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressBackend_To_v1beta1_IngressBackend(a.(*networking.IngressBackend), b.(*v1beta1.IngressBackend), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressList)(nil), (*networking.IngressList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressList_To_networking_IngressList(a.(*v1beta1.IngressList), b.(*networking.IngressList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressList)(nil), (*v1beta1.IngressList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressList_To_v1beta1_IngressList(a.(*networking.IngressList), b.(*v1beta1.IngressList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressRule)(nil), (*networking.IngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressRule_To_networking_IngressRule(a.(*v1beta1.IngressRule), b.(*networking.IngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressRule)(nil), (*v1beta1.IngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressRule_To_v1beta1_IngressRule(a.(*networking.IngressRule), b.(*v1beta1.IngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressRuleValue)(nil), (*networking.IngressRuleValue)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressRuleValue_To_networking_IngressRuleValue(a.(*v1beta1.IngressRuleValue), b.(*networking.IngressRuleValue), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressRuleValue)(nil), (*v1beta1.IngressRuleValue)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressRuleValue_To_v1beta1_IngressRuleValue(a.(*networking.IngressRuleValue), b.(*v1beta1.IngressRuleValue), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressSpec)(nil), (*networking.IngressSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressSpec_To_networking_IngressSpec(a.(*v1beta1.IngressSpec), b.(*networking.IngressSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressSpec)(nil), (*v1beta1.IngressSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressSpec_To_v1beta1_IngressSpec(a.(*networking.IngressSpec), b.(*v1beta1.IngressSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressStatus)(nil), (*networking.IngressStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressStatus_To_networking_IngressStatus(a.(*v1beta1.IngressStatus), b.(*networking.IngressStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressStatus)(nil), (*v1beta1.IngressStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressStatus_To_v1beta1_IngressStatus(a.(*networking.IngressStatus), b.(*v1beta1.IngressStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.IngressTLS)(nil), (*networking.IngressTLS)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IngressTLS_To_networking_IngressTLS(a.(*v1beta1.IngressTLS), b.(*networking.IngressTLS), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.IngressTLS)(nil), (*v1beta1.IngressTLS)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IngressTLS_To_v1beta1_IngressTLS(a.(*networking.IngressTLS), b.(*v1beta1.IngressTLS), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicy)(nil), (*networking.NetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicy_To_networking_NetworkPolicy(a.(*v1beta1.NetworkPolicy), b.(*networking.NetworkPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicy)(nil), (*v1beta1.NetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicy_To_v1beta1_NetworkPolicy(a.(*networking.NetworkPolicy), b.(*v1beta1.NetworkPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicyEgressRule)(nil), (*networking.NetworkPolicyEgressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyEgressRule_To_networking_NetworkPolicyEgressRule(a.(*v1beta1.NetworkPolicyEgressRule), b.(*networking.NetworkPolicyEgressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicyEgressRule)(nil), (*v1beta1.NetworkPolicyEgressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyEgressRule_To_v1beta1_NetworkPolicyEgressRule(a.(*networking.NetworkPolicyEgressRule), b.(*v1beta1.NetworkPolicyEgressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicyIngressRule)(nil), (*networking.NetworkPolicyIngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyIngressRule_To_networking_NetworkPolicyIngressRule(a.(*v1beta1.NetworkPolicyIngressRule), b.(*networking.NetworkPolicyIngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicyIngressRule)(nil), (*v1beta1.NetworkPolicyIngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyIngressRule_To_v1beta1_NetworkPolicyIngressRule(a.(*networking.NetworkPolicyIngressRule), b.(*v1beta1.NetworkPolicyIngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicyList)(nil), (*networking.NetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyList_To_networking_NetworkPolicyList(a.(*v1beta1.NetworkPolicyList), b.(*networking.NetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicyList)(nil), (*v1beta1.NetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyList_To_v1beta1_NetworkPolicyList(a.(*networking.NetworkPolicyList), b.(*v1beta1.NetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicyPeer)(nil), (*networking.NetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(a.(*v1beta1.NetworkPolicyPeer), b.(*networking.NetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicyPeer)(nil), (*v1beta1.NetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyPeer_To_v1beta1_NetworkPolicyPeer(a.(*networking.NetworkPolicyPeer), b.(*v1beta1.NetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicyPort)(nil), (*networking.NetworkPolicyPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyPort_To_networking_NetworkPolicyPort(a.(*v1beta1.NetworkPolicyPort), b.(*networking.NetworkPolicyPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicyPort)(nil), (*v1beta1.NetworkPolicyPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyPort_To_v1beta1_NetworkPolicyPort(a.(*networking.NetworkPolicyPort), b.(*v1beta1.NetworkPolicyPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NetworkPolicySpec)(nil), (*networking.NetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicySpec_To_networking_NetworkPolicySpec(a.(*v1beta1.NetworkPolicySpec), b.(*networking.NetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*networking.NetworkPolicySpec)(nil), (*v1beta1.NetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicySpec_To_v1beta1_NetworkPolicySpec(a.(*networking.NetworkPolicySpec), b.(*v1beta1.NetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodSecurityPolicy)(nil), (*policy.PodSecurityPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodSecurityPolicy_To_policy_PodSecurityPolicy(a.(*v1beta1.PodSecurityPolicy), b.(*policy.PodSecurityPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodSecurityPolicy)(nil), (*v1beta1.PodSecurityPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodSecurityPolicy_To_v1beta1_PodSecurityPolicy(a.(*policy.PodSecurityPolicy), b.(*v1beta1.PodSecurityPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodSecurityPolicyList)(nil), (*policy.PodSecurityPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodSecurityPolicyList_To_policy_PodSecurityPolicyList(a.(*v1beta1.PodSecurityPolicyList), b.(*policy.PodSecurityPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodSecurityPolicyList)(nil), (*v1beta1.PodSecurityPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodSecurityPolicyList_To_v1beta1_PodSecurityPolicyList(a.(*policy.PodSecurityPolicyList), b.(*v1beta1.PodSecurityPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodSecurityPolicySpec)(nil), (*policy.PodSecurityPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodSecurityPolicySpec_To_policy_PodSecurityPolicySpec(a.(*v1beta1.PodSecurityPolicySpec), b.(*policy.PodSecurityPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodSecurityPolicySpec)(nil), (*v1beta1.PodSecurityPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodSecurityPolicySpec_To_v1beta1_PodSecurityPolicySpec(a.(*policy.PodSecurityPolicySpec), b.(*v1beta1.PodSecurityPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicaSet)(nil), (*apps.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSet_To_apps_ReplicaSet(a.(*v1beta1.ReplicaSet), b.(*apps.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSet)(nil), (*v1beta1.ReplicaSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSet_To_v1beta1_ReplicaSet(a.(*apps.ReplicaSet), b.(*v1beta1.ReplicaSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicaSetCondition)(nil), (*apps.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSetCondition_To_apps_ReplicaSetCondition(a.(*v1beta1.ReplicaSetCondition), b.(*apps.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetCondition)(nil), (*v1beta1.ReplicaSetCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetCondition_To_v1beta1_ReplicaSetCondition(a.(*apps.ReplicaSetCondition), b.(*v1beta1.ReplicaSetCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicaSetList)(nil), (*apps.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSetList_To_apps_ReplicaSetList(a.(*v1beta1.ReplicaSetList), b.(*apps.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetList)(nil), (*v1beta1.ReplicaSetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetList_To_v1beta1_ReplicaSetList(a.(*apps.ReplicaSetList), b.(*v1beta1.ReplicaSetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1beta1.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1beta1.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1beta1_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1beta1.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicaSetStatus)(nil), (*apps.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSetStatus_To_apps_ReplicaSetStatus(a.(*v1beta1.ReplicaSetStatus), b.(*apps.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.ReplicaSetStatus)(nil), (*v1beta1.ReplicaSetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetStatus_To_v1beta1_ReplicaSetStatus(a.(*apps.ReplicaSetStatus), b.(*v1beta1.ReplicaSetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ReplicationControllerDummy)(nil), (*extensions.ReplicationControllerDummy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicationControllerDummy_To_extensions_ReplicationControllerDummy(a.(*v1beta1.ReplicationControllerDummy), b.(*extensions.ReplicationControllerDummy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extensions.ReplicationControllerDummy)(nil), (*v1beta1.ReplicationControllerDummy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extensions_ReplicationControllerDummy_To_v1beta1_ReplicationControllerDummy(a.(*extensions.ReplicationControllerDummy), b.(*v1beta1.ReplicationControllerDummy), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1beta1.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1beta1.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1beta1.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1beta1_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1beta1.RollingUpdateDaemonSet), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1beta1.RunAsGroupStrategyOptions)(nil), (*policy.RunAsGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RunAsGroupStrategyOptions_To_policy_RunAsGroupStrategyOptions(a.(*v1beta1.RunAsGroupStrategyOptions), b.(*policy.RunAsGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.RunAsGroupStrategyOptions)(nil), (*v1beta1.RunAsGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_RunAsGroupStrategyOptions_To_v1beta1_RunAsGroupStrategyOptions(a.(*policy.RunAsGroupStrategyOptions), b.(*v1beta1.RunAsGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RunAsUserStrategyOptions)(nil), (*policy.RunAsUserStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RunAsUserStrategyOptions_To_policy_RunAsUserStrategyOptions(a.(*v1beta1.RunAsUserStrategyOptions), b.(*policy.RunAsUserStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.RunAsUserStrategyOptions)(nil), (*v1beta1.RunAsUserStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_RunAsUserStrategyOptions_To_v1beta1_RunAsUserStrategyOptions(a.(*policy.RunAsUserStrategyOptions), b.(*v1beta1.RunAsUserStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RuntimeClassStrategyOptions)(nil), (*policy.RuntimeClassStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RuntimeClassStrategyOptions_To_policy_RuntimeClassStrategyOptions(a.(*v1beta1.RuntimeClassStrategyOptions), b.(*policy.RuntimeClassStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.RuntimeClassStrategyOptions)(nil), (*v1beta1.RuntimeClassStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_RuntimeClassStrategyOptions_To_v1beta1_RuntimeClassStrategyOptions(a.(*policy.RuntimeClassStrategyOptions), b.(*v1beta1.RuntimeClassStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SELinuxStrategyOptions)(nil), (*policy.SELinuxStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SELinuxStrategyOptions_To_policy_SELinuxStrategyOptions(a.(*v1beta1.SELinuxStrategyOptions), b.(*policy.SELinuxStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.SELinuxStrategyOptions)(nil), (*v1beta1.SELinuxStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_SELinuxStrategyOptions_To_v1beta1_SELinuxStrategyOptions(a.(*policy.SELinuxStrategyOptions), b.(*v1beta1.SELinuxStrategyOptions), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1beta1.SupplementalGroupsStrategyOptions)(nil), (*policy.SupplementalGroupsStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SupplementalGroupsStrategyOptions_To_policy_SupplementalGroupsStrategyOptions(a.(*v1beta1.SupplementalGroupsStrategyOptions), b.(*policy.SupplementalGroupsStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.SupplementalGroupsStrategyOptions)(nil), (*v1beta1.SupplementalGroupsStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_SupplementalGroupsStrategyOptions_To_v1beta1_SupplementalGroupsStrategyOptions(a.(*policy.SupplementalGroupsStrategyOptions), b.(*v1beta1.SupplementalGroupsStrategyOptions), scope)
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
	if err := s.AddConversionFunc((*apps.ReplicaSetSpec)(nil), (*v1beta1.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_ReplicaSetSpec_To_v1beta1_ReplicaSetSpec(a.(*apps.ReplicaSetSpec), b.(*v1beta1.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDaemonSet)(nil), (*v1beta1.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDaemonSet_To_v1beta1_RollingUpdateDaemonSet(a.(*apps.RollingUpdateDaemonSet), b.(*v1beta1.RollingUpdateDaemonSet), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apps.RollingUpdateDeployment)(nil), (*v1beta1.RollingUpdateDeployment)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apps_RollingUpdateDeployment_To_v1beta1_RollingUpdateDeployment(a.(*apps.RollingUpdateDeployment), b.(*v1beta1.RollingUpdateDeployment), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*autoscaling.ScaleStatus)(nil), (*v1beta1.ScaleStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(a.(*autoscaling.ScaleStatus), b.(*v1beta1.ScaleStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.IPBlock)(nil), (*v1beta1.IPBlock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_IPBlock_To_v1beta1_IPBlock(a.(*networking.IPBlock), b.(*v1beta1.IPBlock), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicyEgressRule)(nil), (*v1beta1.NetworkPolicyEgressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyEgressRule_To_v1beta1_NetworkPolicyEgressRule(a.(*networking.NetworkPolicyEgressRule), b.(*v1beta1.NetworkPolicyEgressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicyIngressRule)(nil), (*v1beta1.NetworkPolicyIngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyIngressRule_To_v1beta1_NetworkPolicyIngressRule(a.(*networking.NetworkPolicyIngressRule), b.(*v1beta1.NetworkPolicyIngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicyList)(nil), (*v1beta1.NetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyList_To_v1beta1_NetworkPolicyList(a.(*networking.NetworkPolicyList), b.(*v1beta1.NetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicyPeer)(nil), (*v1beta1.NetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyPeer_To_v1beta1_NetworkPolicyPeer(a.(*networking.NetworkPolicyPeer), b.(*v1beta1.NetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicyPort)(nil), (*v1beta1.NetworkPolicyPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicyPort_To_v1beta1_NetworkPolicyPort(a.(*networking.NetworkPolicyPort), b.(*v1beta1.NetworkPolicyPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicySpec)(nil), (*v1beta1.NetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicySpec_To_v1beta1_NetworkPolicySpec(a.(*networking.NetworkPolicySpec), b.(*v1beta1.NetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*networking.NetworkPolicy)(nil), (*v1beta1.NetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_networking_NetworkPolicy_To_v1beta1_NetworkPolicy(a.(*networking.NetworkPolicy), b.(*v1beta1.NetworkPolicy), scope)
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
	if err := s.AddConversionFunc((*v1beta1.IPBlock)(nil), (*networking.IPBlock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_IPBlock_To_networking_IPBlock(a.(*v1beta1.IPBlock), b.(*networking.IPBlock), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicyEgressRule)(nil), (*networking.NetworkPolicyEgressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyEgressRule_To_networking_NetworkPolicyEgressRule(a.(*v1beta1.NetworkPolicyEgressRule), b.(*networking.NetworkPolicyEgressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicyIngressRule)(nil), (*networking.NetworkPolicyIngressRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyIngressRule_To_networking_NetworkPolicyIngressRule(a.(*v1beta1.NetworkPolicyIngressRule), b.(*networking.NetworkPolicyIngressRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicyList)(nil), (*networking.NetworkPolicyList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyList_To_networking_NetworkPolicyList(a.(*v1beta1.NetworkPolicyList), b.(*networking.NetworkPolicyList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicyPeer)(nil), (*networking.NetworkPolicyPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(a.(*v1beta1.NetworkPolicyPeer), b.(*networking.NetworkPolicyPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicyPort)(nil), (*networking.NetworkPolicyPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicyPort_To_networking_NetworkPolicyPort(a.(*v1beta1.NetworkPolicyPort), b.(*networking.NetworkPolicyPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicySpec)(nil), (*networking.NetworkPolicySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicySpec_To_networking_NetworkPolicySpec(a.(*v1beta1.NetworkPolicySpec), b.(*networking.NetworkPolicySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.NetworkPolicy)(nil), (*networking.NetworkPolicy)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NetworkPolicy_To_networking_NetworkPolicy(a.(*v1beta1.NetworkPolicy), b.(*networking.NetworkPolicy), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.ReplicaSetSpec)(nil), (*apps.ReplicaSetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ReplicaSetSpec_To_apps_ReplicaSetSpec(a.(*v1beta1.ReplicaSetSpec), b.(*apps.ReplicaSetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.RollingUpdateDaemonSet)(nil), (*apps.RollingUpdateDaemonSet)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RollingUpdateDaemonSet_To_apps_RollingUpdateDaemonSet(a.(*v1beta1.RollingUpdateDaemonSet), b.(*apps.RollingUpdateDaemonSet), scope)
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
	return nil
}
