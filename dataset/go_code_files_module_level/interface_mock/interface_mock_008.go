func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.Action)(nil), (*authorization.Action)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Action_To_authorization_Action(a.(*v1.Action), b.(*authorization.Action), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.Action)(nil), (*v1.Action)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_Action_To_v1_Action(a.(*authorization.Action), b.(*v1.Action), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRole)(nil), (*authorization.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRole_To_authorization_ClusterRole(a.(*v1.ClusterRole), b.(*authorization.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ClusterRole)(nil), (*v1.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ClusterRole_To_v1_ClusterRole(a.(*authorization.ClusterRole), b.(*v1.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleBinding)(nil), (*authorization.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleBinding_To_authorization_ClusterRoleBinding(a.(*v1.ClusterRoleBinding), b.(*authorization.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ClusterRoleBinding)(nil), (*v1.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ClusterRoleBinding_To_v1_ClusterRoleBinding(a.(*authorization.ClusterRoleBinding), b.(*v1.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleBindingList)(nil), (*authorization.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleBindingList_To_authorization_ClusterRoleBindingList(a.(*v1.ClusterRoleBindingList), b.(*authorization.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ClusterRoleBindingList)(nil), (*v1.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ClusterRoleBindingList_To_v1_ClusterRoleBindingList(a.(*authorization.ClusterRoleBindingList), b.(*v1.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleList)(nil), (*authorization.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleList_To_authorization_ClusterRoleList(a.(*v1.ClusterRoleList), b.(*authorization.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ClusterRoleList)(nil), (*v1.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ClusterRoleList_To_v1_ClusterRoleList(a.(*authorization.ClusterRoleList), b.(*v1.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.GroupRestriction)(nil), (*authorization.GroupRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_GroupRestriction_To_authorization_GroupRestriction(a.(*v1.GroupRestriction), b.(*authorization.GroupRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.GroupRestriction)(nil), (*v1.GroupRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_GroupRestriction_To_v1_GroupRestriction(a.(*authorization.GroupRestriction), b.(*v1.GroupRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.IsPersonalSubjectAccessReview)(nil), (*authorization.IsPersonalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_IsPersonalSubjectAccessReview_To_authorization_IsPersonalSubjectAccessReview(a.(*v1.IsPersonalSubjectAccessReview), b.(*authorization.IsPersonalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.IsPersonalSubjectAccessReview)(nil), (*v1.IsPersonalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_IsPersonalSubjectAccessReview_To_v1_IsPersonalSubjectAccessReview(a.(*authorization.IsPersonalSubjectAccessReview), b.(*v1.IsPersonalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LocalResourceAccessReview)(nil), (*authorization.LocalResourceAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalResourceAccessReview_To_authorization_LocalResourceAccessReview(a.(*v1.LocalResourceAccessReview), b.(*authorization.LocalResourceAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.LocalResourceAccessReview)(nil), (*v1.LocalResourceAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_LocalResourceAccessReview_To_v1_LocalResourceAccessReview(a.(*authorization.LocalResourceAccessReview), b.(*v1.LocalResourceAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.LocalSubjectAccessReview)(nil), (*authorization.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalSubjectAccessReview_To_authorization_LocalSubjectAccessReview(a.(*v1.LocalSubjectAccessReview), b.(*authorization.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.LocalSubjectAccessReview)(nil), (*v1.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_LocalSubjectAccessReview_To_v1_LocalSubjectAccessReview(a.(*authorization.LocalSubjectAccessReview), b.(*v1.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PolicyRule)(nil), (*authorization.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PolicyRule_To_authorization_PolicyRule(a.(*v1.PolicyRule), b.(*authorization.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.PolicyRule)(nil), (*v1.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_PolicyRule_To_v1_PolicyRule(a.(*authorization.PolicyRule), b.(*v1.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceAccessReview)(nil), (*authorization.ResourceAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceAccessReview_To_authorization_ResourceAccessReview(a.(*v1.ResourceAccessReview), b.(*authorization.ResourceAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ResourceAccessReview)(nil), (*v1.ResourceAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ResourceAccessReview_To_v1_ResourceAccessReview(a.(*authorization.ResourceAccessReview), b.(*v1.ResourceAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ResourceAccessReviewResponse)(nil), (*authorization.ResourceAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceAccessReviewResponse_To_authorization_ResourceAccessReviewResponse(a.(*v1.ResourceAccessReviewResponse), b.(*authorization.ResourceAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ResourceAccessReviewResponse)(nil), (*v1.ResourceAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ResourceAccessReviewResponse_To_v1_ResourceAccessReviewResponse(a.(*authorization.ResourceAccessReviewResponse), b.(*v1.ResourceAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Role)(nil), (*authorization.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Role_To_authorization_Role(a.(*v1.Role), b.(*authorization.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.Role)(nil), (*v1.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_Role_To_v1_Role(a.(*authorization.Role), b.(*v1.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBinding)(nil), (*authorization.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBinding_To_authorization_RoleBinding(a.(*v1.RoleBinding), b.(*authorization.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleBinding)(nil), (*v1.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBinding_To_v1_RoleBinding(a.(*authorization.RoleBinding), b.(*v1.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBindingList)(nil), (*authorization.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBindingList_To_authorization_RoleBindingList(a.(*v1.RoleBindingList), b.(*authorization.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleBindingList)(nil), (*v1.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBindingList_To_v1_RoleBindingList(a.(*authorization.RoleBindingList), b.(*v1.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBindingRestriction)(nil), (*authorization.RoleBindingRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBindingRestriction_To_authorization_RoleBindingRestriction(a.(*v1.RoleBindingRestriction), b.(*authorization.RoleBindingRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleBindingRestriction)(nil), (*v1.RoleBindingRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBindingRestriction_To_v1_RoleBindingRestriction(a.(*authorization.RoleBindingRestriction), b.(*v1.RoleBindingRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBindingRestrictionList)(nil), (*authorization.RoleBindingRestrictionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBindingRestrictionList_To_authorization_RoleBindingRestrictionList(a.(*v1.RoleBindingRestrictionList), b.(*authorization.RoleBindingRestrictionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleBindingRestrictionList)(nil), (*v1.RoleBindingRestrictionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBindingRestrictionList_To_v1_RoleBindingRestrictionList(a.(*authorization.RoleBindingRestrictionList), b.(*v1.RoleBindingRestrictionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBindingRestrictionSpec)(nil), (*authorization.RoleBindingRestrictionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBindingRestrictionSpec_To_authorization_RoleBindingRestrictionSpec(a.(*v1.RoleBindingRestrictionSpec), b.(*authorization.RoleBindingRestrictionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleBindingRestrictionSpec)(nil), (*v1.RoleBindingRestrictionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBindingRestrictionSpec_To_v1_RoleBindingRestrictionSpec(a.(*authorization.RoleBindingRestrictionSpec), b.(*v1.RoleBindingRestrictionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleList)(nil), (*authorization.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleList_To_authorization_RoleList(a.(*v1.RoleList), b.(*authorization.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.RoleList)(nil), (*v1.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleList_To_v1_RoleList(a.(*authorization.RoleList), b.(*v1.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SelfSubjectRulesReview)(nil), (*authorization.SelfSubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SelfSubjectRulesReview_To_authorization_SelfSubjectRulesReview(a.(*v1.SelfSubjectRulesReview), b.(*authorization.SelfSubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectRulesReview)(nil), (*v1.SelfSubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectRulesReview_To_v1_SelfSubjectRulesReview(a.(*authorization.SelfSubjectRulesReview), b.(*v1.SelfSubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SelfSubjectRulesReviewSpec)(nil), (*authorization.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SelfSubjectRulesReviewSpec_To_authorization_SelfSubjectRulesReviewSpec(a.(*v1.SelfSubjectRulesReviewSpec), b.(*authorization.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectRulesReviewSpec)(nil), (*v1.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectRulesReviewSpec_To_v1_SelfSubjectRulesReviewSpec(a.(*authorization.SelfSubjectRulesReviewSpec), b.(*v1.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountReference)(nil), (*authorization.ServiceAccountReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountReference_To_authorization_ServiceAccountReference(a.(*v1.ServiceAccountReference), b.(*authorization.ServiceAccountReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ServiceAccountReference)(nil), (*v1.ServiceAccountReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ServiceAccountReference_To_v1_ServiceAccountReference(a.(*authorization.ServiceAccountReference), b.(*v1.ServiceAccountReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ServiceAccountRestriction)(nil), (*authorization.ServiceAccountRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ServiceAccountRestriction_To_authorization_ServiceAccountRestriction(a.(*v1.ServiceAccountRestriction), b.(*authorization.ServiceAccountRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ServiceAccountRestriction)(nil), (*v1.ServiceAccountRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ServiceAccountRestriction_To_v1_ServiceAccountRestriction(a.(*authorization.ServiceAccountRestriction), b.(*v1.ServiceAccountRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SubjectAccessReview)(nil), (*authorization.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectAccessReview_To_authorization_SubjectAccessReview(a.(*v1.SubjectAccessReview), b.(*authorization.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectAccessReview)(nil), (*v1.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReview_To_v1_SubjectAccessReview(a.(*authorization.SubjectAccessReview), b.(*v1.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SubjectAccessReviewResponse)(nil), (*authorization.SubjectAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectAccessReviewResponse_To_authorization_SubjectAccessReviewResponse(a.(*v1.SubjectAccessReviewResponse), b.(*authorization.SubjectAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectAccessReviewResponse)(nil), (*v1.SubjectAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReviewResponse_To_v1_SubjectAccessReviewResponse(a.(*authorization.SubjectAccessReviewResponse), b.(*v1.SubjectAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SubjectRulesReview)(nil), (*authorization.SubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectRulesReview_To_authorization_SubjectRulesReview(a.(*v1.SubjectRulesReview), b.(*authorization.SubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectRulesReview)(nil), (*v1.SubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectRulesReview_To_v1_SubjectRulesReview(a.(*authorization.SubjectRulesReview), b.(*v1.SubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SubjectRulesReviewSpec)(nil), (*authorization.SubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectRulesReviewSpec_To_authorization_SubjectRulesReviewSpec(a.(*v1.SubjectRulesReviewSpec), b.(*authorization.SubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectRulesReviewSpec)(nil), (*v1.SubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectRulesReviewSpec_To_v1_SubjectRulesReviewSpec(a.(*authorization.SubjectRulesReviewSpec), b.(*v1.SubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.SubjectRulesReviewStatus)(nil), (*authorization.SubjectRulesReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectRulesReviewStatus_To_authorization_SubjectRulesReviewStatus(a.(*v1.SubjectRulesReviewStatus), b.(*authorization.SubjectRulesReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectRulesReviewStatus)(nil), (*v1.SubjectRulesReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectRulesReviewStatus_To_v1_SubjectRulesReviewStatus(a.(*authorization.SubjectRulesReviewStatus), b.(*v1.SubjectRulesReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.UserRestriction)(nil), (*authorization.UserRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_UserRestriction_To_authorization_UserRestriction(a.(*v1.UserRestriction), b.(*authorization.UserRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.UserRestriction)(nil), (*v1.UserRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_UserRestriction_To_v1_UserRestriction(a.(*authorization.UserRestriction), b.(*v1.UserRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.ClusterRoleBinding)(nil), (*v1.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ClusterRoleBinding_To_v1_ClusterRoleBinding(a.(*authorization.ClusterRoleBinding), b.(*v1.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.LocalSubjectAccessReview)(nil), (*v1.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_LocalSubjectAccessReview_To_v1_LocalSubjectAccessReview(a.(*authorization.LocalSubjectAccessReview), b.(*v1.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.PolicyRule)(nil), (*v1.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_PolicyRule_To_v1_PolicyRule(a.(*authorization.PolicyRule), b.(*v1.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.ResourceAccessReviewResponse)(nil), (*v1.ResourceAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ResourceAccessReviewResponse_To_v1_ResourceAccessReviewResponse(a.(*authorization.ResourceAccessReviewResponse), b.(*v1.ResourceAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.RoleBinding)(nil), (*v1.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_RoleBinding_To_v1_RoleBinding(a.(*authorization.RoleBinding), b.(*v1.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.SelfSubjectRulesReviewSpec)(nil), (*v1.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectRulesReviewSpec_To_v1_SelfSubjectRulesReviewSpec(a.(*authorization.SelfSubjectRulesReviewSpec), b.(*v1.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*authorization.SubjectAccessReview)(nil), (*v1.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReview_To_v1_SubjectAccessReview(a.(*authorization.SubjectAccessReview), b.(*v1.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ClusterRoleBinding)(nil), (*authorization.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleBinding_To_authorization_ClusterRoleBinding(a.(*v1.ClusterRoleBinding), b.(*authorization.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.LocalSubjectAccessReview)(nil), (*authorization.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_LocalSubjectAccessReview_To_authorization_LocalSubjectAccessReview(a.(*v1.LocalSubjectAccessReview), b.(*authorization.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.PolicyRule)(nil), (*authorization.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PolicyRule_To_authorization_PolicyRule(a.(*v1.PolicyRule), b.(*authorization.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.ResourceAccessReviewResponse)(nil), (*authorization.ResourceAccessReviewResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ResourceAccessReviewResponse_To_authorization_ResourceAccessReviewResponse(a.(*v1.ResourceAccessReviewResponse), b.(*authorization.ResourceAccessReviewResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.RoleBinding)(nil), (*authorization.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBinding_To_authorization_RoleBinding(a.(*v1.RoleBinding), b.(*authorization.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.SelfSubjectRulesReviewSpec)(nil), (*authorization.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SelfSubjectRulesReviewSpec_To_authorization_SelfSubjectRulesReviewSpec(a.(*v1.SelfSubjectRulesReviewSpec), b.(*authorization.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.SubjectAccessReview)(nil), (*authorization.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_SubjectAccessReview_To_authorization_SubjectAccessReview(a.(*v1.SubjectAccessReview), b.(*authorization.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	return nil
}
