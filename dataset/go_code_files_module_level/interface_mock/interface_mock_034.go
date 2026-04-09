func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha1.AggregationRule)(nil), (*rbac.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_AggregationRule_To_rbac_AggregationRule(a.(*v1alpha1.AggregationRule), b.(*rbac.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.AggregationRule)(nil), (*v1alpha1.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_AggregationRule_To_v1alpha1_AggregationRule(a.(*rbac.AggregationRule), b.(*v1alpha1.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ClusterRole)(nil), (*rbac.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterRole_To_rbac_ClusterRole(a.(*v1alpha1.ClusterRole), b.(*rbac.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRole)(nil), (*v1alpha1.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRole_To_v1alpha1_ClusterRole(a.(*rbac.ClusterRole), b.(*v1alpha1.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ClusterRoleBinding)(nil), (*rbac.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterRoleBinding_To_rbac_ClusterRoleBinding(a.(*v1alpha1.ClusterRoleBinding), b.(*rbac.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBinding)(nil), (*v1alpha1.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBinding_To_v1alpha1_ClusterRoleBinding(a.(*rbac.ClusterRoleBinding), b.(*v1alpha1.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ClusterRoleBindingList)(nil), (*rbac.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList(a.(*v1alpha1.ClusterRoleBindingList), b.(*rbac.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBindingList)(nil), (*v1alpha1.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBindingList_To_v1alpha1_ClusterRoleBindingList(a.(*rbac.ClusterRoleBindingList), b.(*v1alpha1.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.ClusterRoleList)(nil), (*rbac.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterRoleList_To_rbac_ClusterRoleList(a.(*v1alpha1.ClusterRoleList), b.(*rbac.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleList)(nil), (*v1alpha1.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleList_To_v1alpha1_ClusterRoleList(a.(*rbac.ClusterRoleList), b.(*v1alpha1.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PolicyRule)(nil), (*rbac.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PolicyRule_To_rbac_PolicyRule(a.(*v1alpha1.PolicyRule), b.(*rbac.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.PolicyRule)(nil), (*v1alpha1.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_PolicyRule_To_v1alpha1_PolicyRule(a.(*rbac.PolicyRule), b.(*v1alpha1.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Role)(nil), (*rbac.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Role_To_rbac_Role(a.(*v1alpha1.Role), b.(*rbac.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Role)(nil), (*v1alpha1.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Role_To_v1alpha1_Role(a.(*rbac.Role), b.(*v1alpha1.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.RoleBinding)(nil), (*rbac.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RoleBinding_To_rbac_RoleBinding(a.(*v1alpha1.RoleBinding), b.(*rbac.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBinding)(nil), (*v1alpha1.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBinding_To_v1alpha1_RoleBinding(a.(*rbac.RoleBinding), b.(*v1alpha1.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.RoleBindingList)(nil), (*rbac.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RoleBindingList_To_rbac_RoleBindingList(a.(*v1alpha1.RoleBindingList), b.(*rbac.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBindingList)(nil), (*v1alpha1.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBindingList_To_v1alpha1_RoleBindingList(a.(*rbac.RoleBindingList), b.(*v1alpha1.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.RoleList)(nil), (*rbac.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RoleList_To_rbac_RoleList(a.(*v1alpha1.RoleList), b.(*rbac.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleList)(nil), (*v1alpha1.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleList_To_v1alpha1_RoleList(a.(*rbac.RoleList), b.(*v1alpha1.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.RoleRef)(nil), (*rbac.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RoleRef_To_rbac_RoleRef(a.(*v1alpha1.RoleRef), b.(*rbac.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleRef)(nil), (*v1alpha1.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleRef_To_v1alpha1_RoleRef(a.(*rbac.RoleRef), b.(*v1alpha1.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Subject)(nil), (*rbac.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Subject_To_rbac_Subject(a.(*v1alpha1.Subject), b.(*rbac.Subject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Subject)(nil), (*v1alpha1.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Subject_To_v1alpha1_Subject(a.(*rbac.Subject), b.(*v1alpha1.Subject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*rbac.Subject)(nil), (*v1alpha1.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Subject_To_v1alpha1_Subject(a.(*rbac.Subject), b.(*v1alpha1.Subject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1alpha1.Subject)(nil), (*rbac.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Subject_To_rbac_Subject(a.(*v1alpha1.Subject), b.(*rbac.Subject), scope)
	}); err != nil {
		return err
	}
	return nil
}
