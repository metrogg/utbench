func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1.ClusterRoleScopeRestriction)(nil), (*oauth.ClusterRoleScopeRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ClusterRoleScopeRestriction_To_oauth_ClusterRoleScopeRestriction(a.(*v1.ClusterRoleScopeRestriction), b.(*oauth.ClusterRoleScopeRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.ClusterRoleScopeRestriction)(nil), (*v1.ClusterRoleScopeRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_ClusterRoleScopeRestriction_To_v1_ClusterRoleScopeRestriction(a.(*oauth.ClusterRoleScopeRestriction), b.(*v1.ClusterRoleScopeRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthAccessToken)(nil), (*oauth.OAuthAccessToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthAccessToken_To_oauth_OAuthAccessToken(a.(*v1.OAuthAccessToken), b.(*oauth.OAuthAccessToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthAccessToken)(nil), (*v1.OAuthAccessToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthAccessToken_To_v1_OAuthAccessToken(a.(*oauth.OAuthAccessToken), b.(*v1.OAuthAccessToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthAccessTokenList)(nil), (*oauth.OAuthAccessTokenList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthAccessTokenList_To_oauth_OAuthAccessTokenList(a.(*v1.OAuthAccessTokenList), b.(*oauth.OAuthAccessTokenList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthAccessTokenList)(nil), (*v1.OAuthAccessTokenList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthAccessTokenList_To_v1_OAuthAccessTokenList(a.(*oauth.OAuthAccessTokenList), b.(*v1.OAuthAccessTokenList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthAuthorizeToken)(nil), (*oauth.OAuthAuthorizeToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthAuthorizeToken_To_oauth_OAuthAuthorizeToken(a.(*v1.OAuthAuthorizeToken), b.(*oauth.OAuthAuthorizeToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthAuthorizeToken)(nil), (*v1.OAuthAuthorizeToken)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthAuthorizeToken_To_v1_OAuthAuthorizeToken(a.(*oauth.OAuthAuthorizeToken), b.(*v1.OAuthAuthorizeToken), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthAuthorizeTokenList)(nil), (*oauth.OAuthAuthorizeTokenList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthAuthorizeTokenList_To_oauth_OAuthAuthorizeTokenList(a.(*v1.OAuthAuthorizeTokenList), b.(*oauth.OAuthAuthorizeTokenList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthAuthorizeTokenList)(nil), (*v1.OAuthAuthorizeTokenList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthAuthorizeTokenList_To_v1_OAuthAuthorizeTokenList(a.(*oauth.OAuthAuthorizeTokenList), b.(*v1.OAuthAuthorizeTokenList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthClient)(nil), (*oauth.OAuthClient)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthClient_To_oauth_OAuthClient(a.(*v1.OAuthClient), b.(*oauth.OAuthClient), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthClient)(nil), (*v1.OAuthClient)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthClient_To_v1_OAuthClient(a.(*oauth.OAuthClient), b.(*v1.OAuthClient), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthClientAuthorization)(nil), (*oauth.OAuthClientAuthorization)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthClientAuthorization_To_oauth_OAuthClientAuthorization(a.(*v1.OAuthClientAuthorization), b.(*oauth.OAuthClientAuthorization), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthClientAuthorization)(nil), (*v1.OAuthClientAuthorization)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthClientAuthorization_To_v1_OAuthClientAuthorization(a.(*oauth.OAuthClientAuthorization), b.(*v1.OAuthClientAuthorization), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthClientAuthorizationList)(nil), (*oauth.OAuthClientAuthorizationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthClientAuthorizationList_To_oauth_OAuthClientAuthorizationList(a.(*v1.OAuthClientAuthorizationList), b.(*oauth.OAuthClientAuthorizationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthClientAuthorizationList)(nil), (*v1.OAuthClientAuthorizationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthClientAuthorizationList_To_v1_OAuthClientAuthorizationList(a.(*oauth.OAuthClientAuthorizationList), b.(*v1.OAuthClientAuthorizationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthClientList)(nil), (*oauth.OAuthClientList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthClientList_To_oauth_OAuthClientList(a.(*v1.OAuthClientList), b.(*oauth.OAuthClientList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthClientList)(nil), (*v1.OAuthClientList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthClientList_To_v1_OAuthClientList(a.(*oauth.OAuthClientList), b.(*v1.OAuthClientList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.OAuthRedirectReference)(nil), (*oauth.OAuthRedirectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_OAuthRedirectReference_To_oauth_OAuthRedirectReference(a.(*v1.OAuthRedirectReference), b.(*oauth.OAuthRedirectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.OAuthRedirectReference)(nil), (*v1.OAuthRedirectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_OAuthRedirectReference_To_v1_OAuthRedirectReference(a.(*oauth.OAuthRedirectReference), b.(*v1.OAuthRedirectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.RedirectReference)(nil), (*oauth.RedirectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_RedirectReference_To_oauth_RedirectReference(a.(*v1.RedirectReference), b.(*oauth.RedirectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.RedirectReference)(nil), (*v1.RedirectReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_RedirectReference_To_v1_RedirectReference(a.(*oauth.RedirectReference), b.(*v1.RedirectReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1.ScopeRestriction)(nil), (*oauth.ScopeRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_ScopeRestriction_To_oauth_ScopeRestriction(a.(*v1.ScopeRestriction), b.(*oauth.ScopeRestriction), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*oauth.ScopeRestriction)(nil), (*v1.ScopeRestriction)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_oauth_ScopeRestriction_To_v1_ScopeRestriction(a.(*oauth.ScopeRestriction), b.(*v1.ScopeRestriction), scope)
	}); err != nil {
		return err
	}
	return nil
}
