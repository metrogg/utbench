func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.LocalSubjectAccessReview)(nil), (*authorization.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_LocalSubjectAccessReview_To_authorization_LocalSubjectAccessReview(a.(*v1beta1.LocalSubjectAccessReview), b.(*authorization.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.LocalSubjectAccessReview)(nil), (*v1beta1.LocalSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_LocalSubjectAccessReview_To_v1beta1_LocalSubjectAccessReview(a.(*authorization.LocalSubjectAccessReview), b.(*v1beta1.LocalSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NonResourceAttributes)(nil), (*authorization.NonResourceAttributes)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NonResourceAttributes_To_authorization_NonResourceAttributes(a.(*v1beta1.NonResourceAttributes), b.(*authorization.NonResourceAttributes), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.NonResourceAttributes)(nil), (*v1beta1.NonResourceAttributes)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_NonResourceAttributes_To_v1beta1_NonResourceAttributes(a.(*authorization.NonResourceAttributes), b.(*v1beta1.NonResourceAttributes), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.NonResourceRule)(nil), (*authorization.NonResourceRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_NonResourceRule_To_authorization_NonResourceRule(a.(*v1beta1.NonResourceRule), b.(*authorization.NonResourceRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.NonResourceRule)(nil), (*v1beta1.NonResourceRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_NonResourceRule_To_v1beta1_NonResourceRule(a.(*authorization.NonResourceRule), b.(*v1beta1.NonResourceRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ResourceAttributes)(nil), (*authorization.ResourceAttributes)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ResourceAttributes_To_authorization_ResourceAttributes(a.(*v1beta1.ResourceAttributes), b.(*authorization.ResourceAttributes), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ResourceAttributes)(nil), (*v1beta1.ResourceAttributes)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ResourceAttributes_To_v1beta1_ResourceAttributes(a.(*authorization.ResourceAttributes), b.(*v1beta1.ResourceAttributes), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ResourceRule)(nil), (*authorization.ResourceRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ResourceRule_To_authorization_ResourceRule(a.(*v1beta1.ResourceRule), b.(*authorization.ResourceRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.ResourceRule)(nil), (*v1beta1.ResourceRule)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_ResourceRule_To_v1beta1_ResourceRule(a.(*authorization.ResourceRule), b.(*v1beta1.ResourceRule), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SelfSubjectAccessReview)(nil), (*authorization.SelfSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SelfSubjectAccessReview_To_authorization_SelfSubjectAccessReview(a.(*v1beta1.SelfSubjectAccessReview), b.(*authorization.SelfSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectAccessReview)(nil), (*v1beta1.SelfSubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectAccessReview_To_v1beta1_SelfSubjectAccessReview(a.(*authorization.SelfSubjectAccessReview), b.(*v1beta1.SelfSubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SelfSubjectAccessReviewSpec)(nil), (*authorization.SelfSubjectAccessReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SelfSubjectAccessReviewSpec_To_authorization_SelfSubjectAccessReviewSpec(a.(*v1beta1.SelfSubjectAccessReviewSpec), b.(*authorization.SelfSubjectAccessReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectAccessReviewSpec)(nil), (*v1beta1.SelfSubjectAccessReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectAccessReviewSpec_To_v1beta1_SelfSubjectAccessReviewSpec(a.(*authorization.SelfSubjectAccessReviewSpec), b.(*v1beta1.SelfSubjectAccessReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SelfSubjectRulesReview)(nil), (*authorization.SelfSubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SelfSubjectRulesReview_To_authorization_SelfSubjectRulesReview(a.(*v1beta1.SelfSubjectRulesReview), b.(*authorization.SelfSubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectRulesReview)(nil), (*v1beta1.SelfSubjectRulesReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectRulesReview_To_v1beta1_SelfSubjectRulesReview(a.(*authorization.SelfSubjectRulesReview), b.(*v1beta1.SelfSubjectRulesReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SelfSubjectRulesReviewSpec)(nil), (*authorization.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SelfSubjectRulesReviewSpec_To_authorization_SelfSubjectRulesReviewSpec(a.(*v1beta1.SelfSubjectRulesReviewSpec), b.(*authorization.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SelfSubjectRulesReviewSpec)(nil), (*v1beta1.SelfSubjectRulesReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SelfSubjectRulesReviewSpec_To_v1beta1_SelfSubjectRulesReviewSpec(a.(*authorization.SelfSubjectRulesReviewSpec), b.(*v1beta1.SelfSubjectRulesReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SubjectAccessReview)(nil), (*authorization.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SubjectAccessReview_To_authorization_SubjectAccessReview(a.(*v1beta1.SubjectAccessReview), b.(*authorization.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectAccessReview)(nil), (*v1beta1.SubjectAccessReview)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReview_To_v1beta1_SubjectAccessReview(a.(*authorization.SubjectAccessReview), b.(*v1beta1.SubjectAccessReview), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SubjectAccessReviewSpec)(nil), (*authorization.SubjectAccessReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SubjectAccessReviewSpec_To_authorization_SubjectAccessReviewSpec(a.(*v1beta1.SubjectAccessReviewSpec), b.(*authorization.SubjectAccessReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectAccessReviewSpec)(nil), (*v1beta1.SubjectAccessReviewSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReviewSpec_To_v1beta1_SubjectAccessReviewSpec(a.(*authorization.SubjectAccessReviewSpec), b.(*v1beta1.SubjectAccessReviewSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SubjectAccessReviewStatus)(nil), (*authorization.SubjectAccessReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SubjectAccessReviewStatus_To_authorization_SubjectAccessReviewStatus(a.(*v1beta1.SubjectAccessReviewStatus), b.(*authorization.SubjectAccessReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectAccessReviewStatus)(nil), (*v1beta1.SubjectAccessReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectAccessReviewStatus_To_v1beta1_SubjectAccessReviewStatus(a.(*authorization.SubjectAccessReviewStatus), b.(*v1beta1.SubjectAccessReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.SubjectRulesReviewStatus)(nil), (*authorization.SubjectRulesReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_SubjectRulesReviewStatus_To_authorization_SubjectRulesReviewStatus(a.(*v1beta1.SubjectRulesReviewStatus), b.(*authorization.SubjectRulesReviewStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*authorization.SubjectRulesReviewStatus)(nil), (*v1beta1.SubjectRulesReviewStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_authorization_SubjectRulesReviewStatus_To_v1beta1_SubjectRulesReviewStatus(a.(*authorization.SubjectRulesReviewStatus), b.(*v1beta1.SubjectRulesReviewStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}
