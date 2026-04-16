func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*CustomResourceColumnDefinition)(nil), (*apiextensions.CustomResourceColumnDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceColumnDefinition_To_apiextensions_CustomResourceColumnDefinition(a.(*CustomResourceColumnDefinition), b.(*apiextensions.CustomResourceColumnDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceColumnDefinition)(nil), (*CustomResourceColumnDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceColumnDefinition_To_v1beta1_CustomResourceColumnDefinition(a.(*apiextensions.CustomResourceColumnDefinition), b.(*CustomResourceColumnDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceConversion)(nil), (*apiextensions.CustomResourceConversion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceConversion_To_apiextensions_CustomResourceConversion(a.(*CustomResourceConversion), b.(*apiextensions.CustomResourceConversion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceConversion)(nil), (*CustomResourceConversion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceConversion_To_v1beta1_CustomResourceConversion(a.(*apiextensions.CustomResourceConversion), b.(*CustomResourceConversion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinition)(nil), (*apiextensions.CustomResourceDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(a.(*CustomResourceDefinition), b.(*apiextensions.CustomResourceDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinition)(nil), (*CustomResourceDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinition_To_v1beta1_CustomResourceDefinition(a.(*apiextensions.CustomResourceDefinition), b.(*CustomResourceDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionCondition)(nil), (*apiextensions.CustomResourceDefinitionCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionCondition_To_apiextensions_CustomResourceDefinitionCondition(a.(*CustomResourceDefinitionCondition), b.(*apiextensions.CustomResourceDefinitionCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionCondition)(nil), (*CustomResourceDefinitionCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionCondition_To_v1beta1_CustomResourceDefinitionCondition(a.(*apiextensions.CustomResourceDefinitionCondition), b.(*CustomResourceDefinitionCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionList)(nil), (*apiextensions.CustomResourceDefinitionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionList_To_apiextensions_CustomResourceDefinitionList(a.(*CustomResourceDefinitionList), b.(*apiextensions.CustomResourceDefinitionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionList)(nil), (*CustomResourceDefinitionList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionList_To_v1beta1_CustomResourceDefinitionList(a.(*apiextensions.CustomResourceDefinitionList), b.(*CustomResourceDefinitionList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionNames)(nil), (*apiextensions.CustomResourceDefinitionNames)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionNames_To_apiextensions_CustomResourceDefinitionNames(a.(*CustomResourceDefinitionNames), b.(*apiextensions.CustomResourceDefinitionNames), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionNames)(nil), (*CustomResourceDefinitionNames)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionNames_To_v1beta1_CustomResourceDefinitionNames(a.(*apiextensions.CustomResourceDefinitionNames), b.(*CustomResourceDefinitionNames), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionSpec)(nil), (*apiextensions.CustomResourceDefinitionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionSpec_To_apiextensions_CustomResourceDefinitionSpec(a.(*CustomResourceDefinitionSpec), b.(*apiextensions.CustomResourceDefinitionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionSpec)(nil), (*CustomResourceDefinitionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionSpec_To_v1beta1_CustomResourceDefinitionSpec(a.(*apiextensions.CustomResourceDefinitionSpec), b.(*CustomResourceDefinitionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionStatus)(nil), (*apiextensions.CustomResourceDefinitionStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionStatus_To_apiextensions_CustomResourceDefinitionStatus(a.(*CustomResourceDefinitionStatus), b.(*apiextensions.CustomResourceDefinitionStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionStatus)(nil), (*CustomResourceDefinitionStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionStatus_To_v1beta1_CustomResourceDefinitionStatus(a.(*apiextensions.CustomResourceDefinitionStatus), b.(*CustomResourceDefinitionStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceDefinitionVersion)(nil), (*apiextensions.CustomResourceDefinitionVersion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceDefinitionVersion_To_apiextensions_CustomResourceDefinitionVersion(a.(*CustomResourceDefinitionVersion), b.(*apiextensions.CustomResourceDefinitionVersion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceDefinitionVersion)(nil), (*CustomResourceDefinitionVersion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceDefinitionVersion_To_v1beta1_CustomResourceDefinitionVersion(a.(*apiextensions.CustomResourceDefinitionVersion), b.(*CustomResourceDefinitionVersion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceSubresourceScale)(nil), (*apiextensions.CustomResourceSubresourceScale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceSubresourceScale_To_apiextensions_CustomResourceSubresourceScale(a.(*CustomResourceSubresourceScale), b.(*apiextensions.CustomResourceSubresourceScale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceSubresourceScale)(nil), (*CustomResourceSubresourceScale)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceSubresourceScale_To_v1beta1_CustomResourceSubresourceScale(a.(*apiextensions.CustomResourceSubresourceScale), b.(*CustomResourceSubresourceScale), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceSubresourceStatus)(nil), (*apiextensions.CustomResourceSubresourceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceSubresourceStatus_To_apiextensions_CustomResourceSubresourceStatus(a.(*CustomResourceSubresourceStatus), b.(*apiextensions.CustomResourceSubresourceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceSubresourceStatus)(nil), (*CustomResourceSubresourceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceSubresourceStatus_To_v1beta1_CustomResourceSubresourceStatus(a.(*apiextensions.CustomResourceSubresourceStatus), b.(*CustomResourceSubresourceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceSubresources)(nil), (*apiextensions.CustomResourceSubresources)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceSubresources_To_apiextensions_CustomResourceSubresources(a.(*CustomResourceSubresources), b.(*apiextensions.CustomResourceSubresources), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceSubresources)(nil), (*CustomResourceSubresources)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceSubresources_To_v1beta1_CustomResourceSubresources(a.(*apiextensions.CustomResourceSubresources), b.(*CustomResourceSubresources), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CustomResourceValidation)(nil), (*apiextensions.CustomResourceValidation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_CustomResourceValidation_To_apiextensions_CustomResourceValidation(a.(*CustomResourceValidation), b.(*apiextensions.CustomResourceValidation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.CustomResourceValidation)(nil), (*CustomResourceValidation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_CustomResourceValidation_To_v1beta1_CustomResourceValidation(a.(*apiextensions.CustomResourceValidation), b.(*CustomResourceValidation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ExternalDocumentation)(nil), (*apiextensions.ExternalDocumentation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ExternalDocumentation_To_apiextensions_ExternalDocumentation(a.(*ExternalDocumentation), b.(*apiextensions.ExternalDocumentation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.ExternalDocumentation)(nil), (*ExternalDocumentation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_ExternalDocumentation_To_v1beta1_ExternalDocumentation(a.(*apiextensions.ExternalDocumentation), b.(*ExternalDocumentation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JSON)(nil), (*apiextensions.JSON)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSON_To_apiextensions_JSON(a.(*JSON), b.(*apiextensions.JSON), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.JSON)(nil), (*JSON)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSON_To_v1beta1_JSON(a.(*apiextensions.JSON), b.(*JSON), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JSONSchemaProps)(nil), (*apiextensions.JSONSchemaProps)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(a.(*JSONSchemaProps), b.(*apiextensions.JSONSchemaProps), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.JSONSchemaProps)(nil), (*JSONSchemaProps)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSONSchemaProps_To_v1beta1_JSONSchemaProps(a.(*apiextensions.JSONSchemaProps), b.(*JSONSchemaProps), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JSONSchemaPropsOrArray)(nil), (*apiextensions.JSONSchemaPropsOrArray)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSONSchemaPropsOrArray_To_apiextensions_JSONSchemaPropsOrArray(a.(*JSONSchemaPropsOrArray), b.(*apiextensions.JSONSchemaPropsOrArray), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.JSONSchemaPropsOrArray)(nil), (*JSONSchemaPropsOrArray)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSONSchemaPropsOrArray_To_v1beta1_JSONSchemaPropsOrArray(a.(*apiextensions.JSONSchemaPropsOrArray), b.(*JSONSchemaPropsOrArray), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JSONSchemaPropsOrBool)(nil), (*apiextensions.JSONSchemaPropsOrBool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSONSchemaPropsOrBool_To_apiextensions_JSONSchemaPropsOrBool(a.(*JSONSchemaPropsOrBool), b.(*apiextensions.JSONSchemaPropsOrBool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.JSONSchemaPropsOrBool)(nil), (*JSONSchemaPropsOrBool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSONSchemaPropsOrBool_To_v1beta1_JSONSchemaPropsOrBool(a.(*apiextensions.JSONSchemaPropsOrBool), b.(*JSONSchemaPropsOrBool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*JSONSchemaPropsOrStringArray)(nil), (*apiextensions.JSONSchemaPropsOrStringArray)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSONSchemaPropsOrStringArray_To_apiextensions_JSONSchemaPropsOrStringArray(a.(*JSONSchemaPropsOrStringArray), b.(*apiextensions.JSONSchemaPropsOrStringArray), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.JSONSchemaPropsOrStringArray)(nil), (*JSONSchemaPropsOrStringArray)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSONSchemaPropsOrStringArray_To_v1beta1_JSONSchemaPropsOrStringArray(a.(*apiextensions.JSONSchemaPropsOrStringArray), b.(*JSONSchemaPropsOrStringArray), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ServiceReference)(nil), (*apiextensions.ServiceReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ServiceReference_To_apiextensions_ServiceReference(a.(*ServiceReference), b.(*apiextensions.ServiceReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.ServiceReference)(nil), (*ServiceReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_ServiceReference_To_v1beta1_ServiceReference(a.(*apiextensions.ServiceReference), b.(*ServiceReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*WebhookClientConfig)(nil), (*apiextensions.WebhookClientConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_WebhookClientConfig_To_apiextensions_WebhookClientConfig(a.(*WebhookClientConfig), b.(*apiextensions.WebhookClientConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*apiextensions.WebhookClientConfig)(nil), (*WebhookClientConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_WebhookClientConfig_To_v1beta1_WebhookClientConfig(a.(*apiextensions.WebhookClientConfig), b.(*WebhookClientConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apiextensions.JSONSchemaProps)(nil), (*JSONSchemaProps)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSONSchemaProps_To_v1beta1_JSONSchemaProps(a.(*apiextensions.JSONSchemaProps), b.(*JSONSchemaProps), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apiextensions.JSON)(nil), (*JSON)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_apiextensions_JSON_To_v1beta1_JSON(a.(*apiextensions.JSON), b.(*JSON), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*JSON)(nil), (*apiextensions.JSON)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_JSON_To_apiextensions_JSON(a.(*JSON), b.(*apiextensions.JSON), scope)
	}); err != nil {
		return err
	}
	return nil
}
