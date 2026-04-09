func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.BrokerTemplateInstance)(nil), (*template.BrokerTemplateInstance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BrokerTemplateInstance_To_template_BrokerTemplateInstance(a.(*v1.BrokerTemplateInstance), b.(*template.BrokerTemplateInstance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.BrokerTemplateInstance)(nil), (*v1.BrokerTemplateInstance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_BrokerTemplateInstance_To_v1_BrokerTemplateInstance(a.(*template.BrokerTemplateInstance), b.(*v1.BrokerTemplateInstance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BrokerTemplateInstanceList)(nil), (*template.BrokerTemplateInstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BrokerTemplateInstanceList_To_template_BrokerTemplateInstanceList(a.(*v1.BrokerTemplateInstanceList), b.(*template.BrokerTemplateInstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.BrokerTemplateInstanceList)(nil), (*v1.BrokerTemplateInstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_BrokerTemplateInstanceList_To_v1_BrokerTemplateInstanceList(a.(*template.BrokerTemplateInstanceList), b.(*v1.BrokerTemplateInstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.BrokerTemplateInstanceSpec)(nil), (*template.BrokerTemplateInstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_BrokerTemplateInstanceSpec_To_template_BrokerTemplateInstanceSpec(a.(*v1.BrokerTemplateInstanceSpec), b.(*template.BrokerTemplateInstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.BrokerTemplateInstanceSpec)(nil), (*v1.BrokerTemplateInstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_BrokerTemplateInstanceSpec_To_v1_BrokerTemplateInstanceSpec(a.(*template.BrokerTemplateInstanceSpec), b.(*v1.BrokerTemplateInstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Parameter)(nil), (*template.Parameter)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Parameter_To_template_Parameter(a.(*v1.Parameter), b.(*template.Parameter), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.Parameter)(nil), (*v1.Parameter)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_Parameter_To_v1_Parameter(a.(*template.Parameter), b.(*v1.Parameter), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.Template)(nil), (*template.Template)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Template_To_template_Template(a.(*v1.Template), b.(*template.Template), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.Template)(nil), (*v1.Template)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_Template_To_v1_Template(a.(*template.Template), b.(*v1.Template), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstance)(nil), (*template.TemplateInstance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstance_To_template_TemplateInstance(a.(*v1.TemplateInstance), b.(*template.TemplateInstance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstance)(nil), (*v1.TemplateInstance)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstance_To_v1_TemplateInstance(a.(*template.TemplateInstance), b.(*v1.TemplateInstance), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceCondition)(nil), (*template.TemplateInstanceCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceCondition_To_template_TemplateInstanceCondition(a.(*v1.TemplateInstanceCondition), b.(*template.TemplateInstanceCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceCondition)(nil), (*v1.TemplateInstanceCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceCondition_To_v1_TemplateInstanceCondition(a.(*template.TemplateInstanceCondition), b.(*v1.TemplateInstanceCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceList)(nil), (*template.TemplateInstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceList_To_template_TemplateInstanceList(a.(*v1.TemplateInstanceList), b.(*template.TemplateInstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceList)(nil), (*v1.TemplateInstanceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceList_To_v1_TemplateInstanceList(a.(*template.TemplateInstanceList), b.(*v1.TemplateInstanceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceObject)(nil), (*template.TemplateInstanceObject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceObject_To_template_TemplateInstanceObject(a.(*v1.TemplateInstanceObject), b.(*template.TemplateInstanceObject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceObject)(nil), (*v1.TemplateInstanceObject)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceObject_To_v1_TemplateInstanceObject(a.(*template.TemplateInstanceObject), b.(*v1.TemplateInstanceObject), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceRequester)(nil), (*template.TemplateInstanceRequester)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceRequester_To_template_TemplateInstanceRequester(a.(*v1.TemplateInstanceRequester), b.(*template.TemplateInstanceRequester), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceRequester)(nil), (*v1.TemplateInstanceRequester)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceRequester_To_v1_TemplateInstanceRequester(a.(*template.TemplateInstanceRequester), b.(*v1.TemplateInstanceRequester), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceSpec)(nil), (*template.TemplateInstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceSpec_To_template_TemplateInstanceSpec(a.(*v1.TemplateInstanceSpec), b.(*template.TemplateInstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceSpec)(nil), (*v1.TemplateInstanceSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceSpec_To_v1_TemplateInstanceSpec(a.(*template.TemplateInstanceSpec), b.(*v1.TemplateInstanceSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateInstanceStatus)(nil), (*template.TemplateInstanceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateInstanceStatus_To_template_TemplateInstanceStatus(a.(*v1.TemplateInstanceStatus), b.(*template.TemplateInstanceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateInstanceStatus)(nil), (*v1.TemplateInstanceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateInstanceStatus_To_v1_TemplateInstanceStatus(a.(*template.TemplateInstanceStatus), b.(*v1.TemplateInstanceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.TemplateList)(nil), (*template.TemplateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_TemplateList_To_template_TemplateList(a.(*v1.TemplateList), b.(*template.TemplateList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*template.TemplateList)(nil), (*v1.TemplateList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_template_TemplateList_To_v1_TemplateList(a.(*template.TemplateList), b.(*v1.TemplateList), scope)
	}); err != nil {
		return err
	}
	return nil
}
