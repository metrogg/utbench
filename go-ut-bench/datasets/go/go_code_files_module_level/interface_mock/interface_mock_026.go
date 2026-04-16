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
	if err := s.AddGeneratedConversionFunc((*v1beta1.Eviction)(nil), (*policy.Eviction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Eviction_To_policy_Eviction(a.(*v1beta1.Eviction), b.(*policy.Eviction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.Eviction)(nil), (*v1beta1.Eviction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_Eviction_To_v1beta1_Eviction(a.(*policy.Eviction), b.(*v1beta1.Eviction), scope)
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
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodDisruptionBudget)(nil), (*policy.PodDisruptionBudget)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodDisruptionBudget_To_policy_PodDisruptionBudget(a.(*v1beta1.PodDisruptionBudget), b.(*policy.PodDisruptionBudget), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodDisruptionBudget)(nil), (*v1beta1.PodDisruptionBudget)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodDisruptionBudget_To_v1beta1_PodDisruptionBudget(a.(*policy.PodDisruptionBudget), b.(*v1beta1.PodDisruptionBudget), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodDisruptionBudgetList)(nil), (*policy.PodDisruptionBudgetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodDisruptionBudgetList_To_policy_PodDisruptionBudgetList(a.(*v1beta1.PodDisruptionBudgetList), b.(*policy.PodDisruptionBudgetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodDisruptionBudgetList)(nil), (*v1beta1.PodDisruptionBudgetList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodDisruptionBudgetList_To_v1beta1_PodDisruptionBudgetList(a.(*policy.PodDisruptionBudgetList), b.(*v1beta1.PodDisruptionBudgetList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodDisruptionBudgetSpec)(nil), (*policy.PodDisruptionBudgetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodDisruptionBudgetSpec_To_policy_PodDisruptionBudgetSpec(a.(*v1beta1.PodDisruptionBudgetSpec), b.(*policy.PodDisruptionBudgetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodDisruptionBudgetSpec)(nil), (*v1beta1.PodDisruptionBudgetSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodDisruptionBudgetSpec_To_v1beta1_PodDisruptionBudgetSpec(a.(*policy.PodDisruptionBudgetSpec), b.(*v1beta1.PodDisruptionBudgetSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PodDisruptionBudgetStatus)(nil), (*policy.PodDisruptionBudgetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PodDisruptionBudgetStatus_To_policy_PodDisruptionBudgetStatus(a.(*v1beta1.PodDisruptionBudgetStatus), b.(*policy.PodDisruptionBudgetStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*policy.PodDisruptionBudgetStatus)(nil), (*v1beta1.PodDisruptionBudgetStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_policy_PodDisruptionBudgetStatus_To_v1beta1_PodDisruptionBudgetStatus(a.(*policy.PodDisruptionBudgetStatus), b.(*v1beta1.PodDisruptionBudgetStatus), scope)
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
	return nil
}
