func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.AggregationRule)(nil), (*rbac.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_AggregationRule_To_rbac_AggregationRule(a.(*v1.AggregationRule), b.(*rbac.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.AggregationRule)(nil), (*v1.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_AggregationRule_To_v1_AggregationRule(a.(*rbac.AggregationRule), b.(*v1.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRole)(nil), (*rbac.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRole_To_rbac_ClusterRole(a.(*v1.ClusterRole), b.(*rbac.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRole)(nil), (*v1.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRole_To_v1_ClusterRole(a.(*rbac.ClusterRole), b.(*v1.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleBinding)(nil), (*rbac.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleBinding_To_rbac_ClusterRoleBinding(a.(*v1.ClusterRoleBinding), b.(*rbac.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBinding)(nil), (*v1.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBinding_To_v1_ClusterRoleBinding(a.(*rbac.ClusterRoleBinding), b.(*v1.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleBindingList)(nil), (*rbac.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList(a.(*v1.ClusterRoleBindingList), b.(*rbac.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBindingList)(nil), (*v1.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBindingList_To_v1_ClusterRoleBindingList(a.(*rbac.ClusterRoleBindingList), b.(*v1.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleList)(nil), (*rbac.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleList_To_rbac_ClusterRoleList(a.(*v1.ClusterRoleList), b.(*rbac.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleList)(nil), (*v1.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleList_To_v1_ClusterRoleList(a.(*rbac.ClusterRoleList), b.(*v1.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.PolicyRule)(nil), (*rbac.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_PolicyRule_To_rbac_PolicyRule(a.(*v1.PolicyRule), b.(*rbac.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.PolicyRule)(nil), (*v1.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_PolicyRule_To_v1_PolicyRule(a.(*rbac.PolicyRule), b.(*v1.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Role)(nil), (*rbac.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Role_To_rbac_Role(a.(*v1.Role), b.(*rbac.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Role)(nil), (*v1.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Role_To_v1_Role(a.(*rbac.Role), b.(*v1.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBinding)(nil), (*rbac.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBinding_To_rbac_RoleBinding(a.(*v1.RoleBinding), b.(*rbac.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBinding)(nil), (*v1.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBinding_To_v1_RoleBinding(a.(*rbac.RoleBinding), b.(*v1.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleBindingList)(nil), (*rbac.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleBindingList_To_rbac_RoleBindingList(a.(*v1.RoleBindingList), b.(*rbac.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBindingList)(nil), (*v1.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBindingList_To_v1_RoleBindingList(a.(*rbac.RoleBindingList), b.(*v1.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleList)(nil), (*rbac.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleList_To_rbac_RoleList(a.(*v1.RoleList), b.(*rbac.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleList)(nil), (*v1.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleList_To_v1_RoleList(a.(*rbac.RoleList), b.(*v1.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RoleRef)(nil), (*rbac.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RoleRef_To_rbac_RoleRef(a.(*v1.RoleRef), b.(*rbac.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleRef)(nil), (*v1.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleRef_To_v1_RoleRef(a.(*rbac.RoleRef), b.(*v1.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Subject)(nil), (*rbac.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Subject_To_rbac_Subject(a.(*v1.Subject), b.(*rbac.Subject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Subject)(nil), (*v1.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Subject_To_v1_Subject(a.(*rbac.Subject), b.(*v1.Subject), scope)
	}); err != nil {
		return err
	}
	return nil
}
