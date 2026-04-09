func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.AllowedFlexVolume)(nil), (*security.AllowedFlexVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AllowedFlexVolume_To_security_AllowedFlexVolume(a.(*v1.AllowedFlexVolume), b.(*security.AllowedFlexVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.AllowedFlexVolume)(nil), (*v1.AllowedFlexVolume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_AllowedFlexVolume_To_v1_AllowedFlexVolume(a.(*security.AllowedFlexVolume), b.(*v1.AllowedFlexVolume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.FSGroupStrategyOptions)(nil), (*security.FSGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_FSGroupStrategyOptions_To_security_FSGroupStrategyOptions(a.(*v1.FSGroupStrategyOptions), b.(*security.FSGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.FSGroupStrategyOptions)(nil), (*v1.FSGroupStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_FSGroupStrategyOptions_To_v1_FSGroupStrategyOptions(a.(*security.FSGroupStrategyOptions), b.(*v1.FSGroupStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.IDRange)(nil), (*security.IDRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_IDRange_To_security_IDRange(a.(*v1.IDRange), b.(*security.IDRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.IDRange)(nil), (*v1.IDRange)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_IDRange_To_v1_IDRange(a.(*security.IDRange), b.(*v1.IDRange), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicyReview)(nil), (*security.PodSecurityPolicyReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicyReview_To_security_PodSecurityPolicyReview(a.(*v1.PodSecurityPolicyReview), b.(*security.PodSecurityPolicyReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicyReview)(nil), (*v1.PodSecurityPolicyReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicyReview_To_v1_PodSecurityPolicyReview(a.(*security.PodSecurityPolicyReview), b.(*v1.PodSecurityPolicyReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicyReviewSpec)(nil), (*security.PodSecurityPolicyReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicyReviewSpec_To_security_PodSecurityPolicyReviewSpec(a.(*v1.PodSecurityPolicyReviewSpec), b.(*security.PodSecurityPolicyReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicyReviewSpec)(nil), (*v1.PodSecurityPolicyReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicyReviewSpec_To_v1_PodSecurityPolicyReviewSpec(a.(*security.PodSecurityPolicyReviewSpec), b.(*v1.PodSecurityPolicyReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicyReviewStatus)(nil), (*security.PodSecurityPolicyReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicyReviewStatus_To_security_PodSecurityPolicyReviewStatus(a.(*v1.PodSecurityPolicyReviewStatus), b.(*security.PodSecurityPolicyReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicyReviewStatus)(nil), (*v1.PodSecurityPolicyReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicyReviewStatus_To_v1_PodSecurityPolicyReviewStatus(a.(*security.PodSecurityPolicyReviewStatus), b.(*v1.PodSecurityPolicyReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicySelfSubjectReview)(nil), (*security.PodSecurityPolicySelfSubjectReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicySelfSubjectReview_To_security_PodSecurityPolicySelfSubjectReview(a.(*v1.PodSecurityPolicySelfSubjectReview), b.(*security.PodSecurityPolicySelfSubjectReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicySelfSubjectReview)(nil), (*v1.PodSecurityPolicySelfSubjectReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicySelfSubjectReview_To_v1_PodSecurityPolicySelfSubjectReview(a.(*security.PodSecurityPolicySelfSubjectReview), b.(*v1.PodSecurityPolicySelfSubjectReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicySelfSubjectReviewSpec)(nil), (*security.PodSecurityPolicySelfSubjectReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicySelfSubjectReviewSpec_To_security_PodSecurityPolicySelfSubjectReviewSpec(a.(*v1.PodSecurityPolicySelfSubjectReviewSpec), b.(*security.PodSecurityPolicySelfSubjectReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicySelfSubjectReviewSpec)(nil), (*v1.PodSecurityPolicySelfSubjectReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicySelfSubjectReviewSpec_To_v1_PodSecurityPolicySelfSubjectReviewSpec(a.(*security.PodSecurityPolicySelfSubjectReviewSpec), b.(*v1.PodSecurityPolicySelfSubjectReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicySubjectReview)(nil), (*security.PodSecurityPolicySubjectReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicySubjectReview_To_security_PodSecurityPolicySubjectReview(a.(*v1.PodSecurityPolicySubjectReview), b.(*security.PodSecurityPolicySubjectReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicySubjectReview)(nil), (*v1.PodSecurityPolicySubjectReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicySubjectReview_To_v1_PodSecurityPolicySubjectReview(a.(*security.PodSecurityPolicySubjectReview), b.(*v1.PodSecurityPolicySubjectReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicySubjectReviewSpec)(nil), (*security.PodSecurityPolicySubjectReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicySubjectReviewSpec_To_security_PodSecurityPolicySubjectReviewSpec(a.(*v1.PodSecurityPolicySubjectReviewSpec), b.(*security.PodSecurityPolicySubjectReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicySubjectReviewSpec)(nil), (*v1.PodSecurityPolicySubjectReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicySubjectReviewSpec_To_v1_PodSecurityPolicySubjectReviewSpec(a.(*security.PodSecurityPolicySubjectReviewSpec), b.(*v1.PodSecurityPolicySubjectReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PodSecurityPolicySubjectReviewStatus)(nil), (*security.PodSecurityPolicySubjectReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PodSecurityPolicySubjectReviewStatus_To_security_PodSecurityPolicySubjectReviewStatus(a.(*v1.PodSecurityPolicySubjectReviewStatus), b.(*security.PodSecurityPolicySubjectReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.PodSecurityPolicySubjectReviewStatus)(nil), (*v1.PodSecurityPolicySubjectReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_PodSecurityPolicySubjectReviewStatus_To_v1_PodSecurityPolicySubjectReviewStatus(a.(*security.PodSecurityPolicySubjectReviewStatus), b.(*v1.PodSecurityPolicySubjectReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RangeAllocation)(nil), (*security.RangeAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RangeAllocation_To_security_RangeAllocation(a.(*v1.RangeAllocation), b.(*security.RangeAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.RangeAllocation)(nil), (*v1.RangeAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_RangeAllocation_To_v1_RangeAllocation(a.(*security.RangeAllocation), b.(*v1.RangeAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RangeAllocationList)(nil), (*security.RangeAllocationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RangeAllocationList_To_security_RangeAllocationList(a.(*v1.RangeAllocationList), b.(*security.RangeAllocationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.RangeAllocationList)(nil), (*v1.RangeAllocationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_RangeAllocationList_To_v1_RangeAllocationList(a.(*security.RangeAllocationList), b.(*v1.RangeAllocationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RunAsUserStrategyOptions)(nil), (*security.RunAsUserStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RunAsUserStrategyOptions_To_security_RunAsUserStrategyOptions(a.(*v1.RunAsUserStrategyOptions), b.(*security.RunAsUserStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.RunAsUserStrategyOptions)(nil), (*v1.RunAsUserStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_RunAsUserStrategyOptions_To_v1_RunAsUserStrategyOptions(a.(*security.RunAsUserStrategyOptions), b.(*v1.RunAsUserStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SELinuxContextStrategyOptions)(nil), (*security.SELinuxContextStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SELinuxContextStrategyOptions_To_security_SELinuxContextStrategyOptions(a.(*v1.SELinuxContextStrategyOptions), b.(*security.SELinuxContextStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.SELinuxContextStrategyOptions)(nil), (*v1.SELinuxContextStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_SELinuxContextStrategyOptions_To_v1_SELinuxContextStrategyOptions(a.(*security.SELinuxContextStrategyOptions), b.(*v1.SELinuxContextStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecurityContextConstraints)(nil), (*security.SecurityContextConstraints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecurityContextConstraints_To_security_SecurityContextConstraints(a.(*v1.SecurityContextConstraints), b.(*security.SecurityContextConstraints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.SecurityContextConstraints)(nil), (*v1.SecurityContextConstraints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_SecurityContextConstraints_To_v1_SecurityContextConstraints(a.(*security.SecurityContextConstraints), b.(*v1.SecurityContextConstraints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SecurityContextConstraintsList)(nil), (*security.SecurityContextConstraintsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecurityContextConstraintsList_To_security_SecurityContextConstraintsList(a.(*v1.SecurityContextConstraintsList), b.(*security.SecurityContextConstraintsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.SecurityContextConstraintsList)(nil), (*v1.SecurityContextConstraintsList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_SecurityContextConstraintsList_To_v1_SecurityContextConstraintsList(a.(*security.SecurityContextConstraintsList), b.(*v1.SecurityContextConstraintsList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountPodSecurityPolicyReviewStatus)(nil), (*security.ServiceAccountPodSecurityPolicyReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountPodSecurityPolicyReviewStatus_To_security_ServiceAccountPodSecurityPolicyReviewStatus(a.(*v1.ServiceAccountPodSecurityPolicyReviewStatus), b.(*security.ServiceAccountPodSecurityPolicyReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.ServiceAccountPodSecurityPolicyReviewStatus)(nil), (*v1.ServiceAccountPodSecurityPolicyReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_ServiceAccountPodSecurityPolicyReviewStatus_To_v1_ServiceAccountPodSecurityPolicyReviewStatus(a.(*security.ServiceAccountPodSecurityPolicyReviewStatus), b.(*v1.ServiceAccountPodSecurityPolicyReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SupplementalGroupsStrategyOptions)(nil), (*security.SupplementalGroupsStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SupplementalGroupsStrategyOptions_To_security_SupplementalGroupsStrategyOptions(a.(*v1.SupplementalGroupsStrategyOptions), b.(*security.SupplementalGroupsStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*security.SupplementalGroupsStrategyOptions)(nil), (*v1.SupplementalGroupsStrategyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_SupplementalGroupsStrategyOptions_To_v1_SupplementalGroupsStrategyOptions(a.(*security.SupplementalGroupsStrategyOptions), b.(*v1.SupplementalGroupsStrategyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*security.SecurityContextConstraints)(nil), (*v1.SecurityContextConstraints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_security_SecurityContextConstraints_To_v1_SecurityContextConstraints(a.(*security.SecurityContextConstraints), b.(*v1.SecurityContextConstraints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.SecurityContextConstraints)(nil), (*security.SecurityContextConstraints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SecurityContextConstraints_To_security_SecurityContextConstraints(a.(*v1.SecurityContextConstraints), b.(*security.SecurityContextConstraints), scope)
	}); err != nil {
		return err
	}
	return nil
}
