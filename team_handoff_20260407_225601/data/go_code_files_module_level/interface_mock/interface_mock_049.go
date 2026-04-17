func RegisterConversions(s *runtime.Scheme) error {
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
	return nil
}
