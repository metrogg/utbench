func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.AggregationRule)(nil), (*rbac.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_AggregationRule_To_rbac_AggregationRule(a.(*v1beta1.AggregationRule), b.(*rbac.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.AggregationRule)(nil), (*v1beta1.AggregationRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_AggregationRule_To_v1beta1_AggregationRule(a.(*rbac.AggregationRule), b.(*v1beta1.AggregationRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ClusterRole)(nil), (*rbac.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterRole_To_rbac_ClusterRole(a.(*v1beta1.ClusterRole), b.(*rbac.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRole)(nil), (*v1beta1.ClusterRole)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRole_To_v1beta1_ClusterRole(a.(*rbac.ClusterRole), b.(*v1beta1.ClusterRole), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ClusterRoleBinding)(nil), (*rbac.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterRoleBinding_To_rbac_ClusterRoleBinding(a.(*v1beta1.ClusterRoleBinding), b.(*rbac.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBinding)(nil), (*v1beta1.ClusterRoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBinding_To_v1beta1_ClusterRoleBinding(a.(*rbac.ClusterRoleBinding), b.(*v1beta1.ClusterRoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ClusterRoleBindingList)(nil), (*rbac.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList(a.(*v1beta1.ClusterRoleBindingList), b.(*rbac.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleBindingList)(nil), (*v1beta1.ClusterRoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleBindingList_To_v1beta1_ClusterRoleBindingList(a.(*rbac.ClusterRoleBindingList), b.(*v1beta1.ClusterRoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ClusterRoleList)(nil), (*rbac.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ClusterRoleList_To_rbac_ClusterRoleList(a.(*v1beta1.ClusterRoleList), b.(*rbac.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.ClusterRoleList)(nil), (*v1beta1.ClusterRoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_ClusterRoleList_To_v1beta1_ClusterRoleList(a.(*rbac.ClusterRoleList), b.(*v1beta1.ClusterRoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.PolicyRule)(nil), (*rbac.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_PolicyRule_To_rbac_PolicyRule(a.(*v1beta1.PolicyRule), b.(*rbac.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.PolicyRule)(nil), (*v1beta1.PolicyRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_PolicyRule_To_v1beta1_PolicyRule(a.(*rbac.PolicyRule), b.(*v1beta1.PolicyRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Role)(nil), (*rbac.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Role_To_rbac_Role(a.(*v1beta1.Role), b.(*rbac.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Role)(nil), (*v1beta1.Role)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Role_To_v1beta1_Role(a.(*rbac.Role), b.(*v1beta1.Role), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RoleBinding)(nil), (*rbac.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RoleBinding_To_rbac_RoleBinding(a.(*v1beta1.RoleBinding), b.(*rbac.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBinding)(nil), (*v1beta1.RoleBinding)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBinding_To_v1beta1_RoleBinding(a.(*rbac.RoleBinding), b.(*v1beta1.RoleBinding), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RoleBindingList)(nil), (*rbac.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RoleBindingList_To_rbac_RoleBindingList(a.(*v1beta1.RoleBindingList), b.(*rbac.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleBindingList)(nil), (*v1beta1.RoleBindingList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleBindingList_To_v1beta1_RoleBindingList(a.(*rbac.RoleBindingList), b.(*v1beta1.RoleBindingList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RoleList)(nil), (*rbac.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RoleList_To_rbac_RoleList(a.(*v1beta1.RoleList), b.(*rbac.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleList)(nil), (*v1beta1.RoleList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleList_To_v1beta1_RoleList(a.(*rbac.RoleList), b.(*v1beta1.RoleList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.RoleRef)(nil), (*rbac.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_RoleRef_To_rbac_RoleRef(a.(*v1beta1.RoleRef), b.(*rbac.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.RoleRef)(nil), (*v1beta1.RoleRef)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_RoleRef_To_v1beta1_RoleRef(a.(*rbac.RoleRef), b.(*v1beta1.RoleRef), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.Subject)(nil), (*rbac.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Subject_To_rbac_Subject(a.(*v1beta1.Subject), b.(*rbac.Subject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*rbac.Subject)(nil), (*v1beta1.Subject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_rbac_Subject_To_v1beta1_Subject(a.(*rbac.Subject), b.(*v1beta1.Subject), scope)
	}); err != nil {
		return err
	}
	return nil
}
