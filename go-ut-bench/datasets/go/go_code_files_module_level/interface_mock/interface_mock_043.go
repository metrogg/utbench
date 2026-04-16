func (service *Service) AddRoutes(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(service.HomePage)
	//Registration form
	router.Methods("GET").Path("/register").HandlerFunc(service.ShowRegistrationForm)
	router.Methods("POST").Path("/register").HandlerFunc(service.ProcessRegistrationForm)
	router.Methods("GET").Path("/phonevalidation").HandlerFunc(service.PhonenumberValidation)
	router.Methods("GET").Path("/phoneregistrationvalidation").HandlerFunc(service.PhonenumberRegistrationValidation)
	router.Methods("GET").Path("/pvl").HandlerFunc(service.PhonenumberValidationAndLogin)
	router.Methods("GET").Path("/emailvalidation").HandlerFunc(service.EmailValidation)
	router.Methods("GET").Path("/emailregistrationvalidation").HandlerFunc(service.EmailRegistrationValidation)
	router.Methods("GET").Path("/register/smsconfirmed").HandlerFunc(service.CheckRegistrationSMSConfirmation)
	router.Methods("GET").Path("/register/emailconfirmed").HandlerFunc(service.CheckRegistrationEmailConfirmation)
	router.Methods("POST").Path("/register/smsconfirmation").HandlerFunc(service.ProcessPhonenumberConfirmationForm)
	router.Methods("POST").Path("/register/validation").HandlerFunc(service.ValidateInfo)
	router.Handle("/register/resendvalidation", alice.New(middleware.RateLimit(middleware.DefaultRateLimitPeriod, middleware.DefaultRateLimit).Handler).Then(http.HandlerFunc(service.ResendValidationInfo))).Methods("POST")
	//Enable us to "forget" users in case we are not in production
	router.Methods("GET").Path("/register/delete").HandlerFunc(service.ServeForgetAccountPage)
	router.Methods("POST").Path("/register/delete").HandlerFunc(service.ForgetAccountHandler)
	//Login forms
	router.Methods("GET").Path("/login").HandlerFunc(service.ShowLoginForm)
	router.Methods("POST").Path("/login").HandlerFunc(service.ProcessLoginForm)
	router.Methods("GET").Path("/login/twofamethods").HandlerFunc(service.GetTwoFactorAuthenticationMethods)
	router.Methods("POST").Path("/login/totpconfirmation").HandlerFunc(service.ProcessTOTPConfirmation)
	router.Handle("/login/smscode/{phoneLabel}", alice.New(middleware.RateLimit(middleware.DefaultRateLimitPeriod, middleware.DefaultRateLimit).Handler).Then(http.HandlerFunc(service.GetSmsCode))).Methods("POST")
	router.Methods("POST").Path("/login/smsconfirmation").HandlerFunc(service.Process2FASMSConfirmation)
	router.Handle("/login/resendsms", alice.New(middleware.RateLimit(middleware.DefaultRateLimitPeriod, middleware.DefaultRateLimit).Handler).Then(http.HandlerFunc(service.LoginResendPhonenumberConfirmation))).Methods("POST")
	router.Methods("GET").Path("/sc").HandlerFunc(service.MobileSMSConfirmation)
	router.Methods("GET").Path("/login/smsconfirmed").HandlerFunc(service.Check2FASMSConfirmation)
	router.Methods("POST").Path("/login/validateemail").HandlerFunc(service.ValidateEmail)
	router.Methods("POST").Path("/login/forgotpassword").HandlerFunc(service.ForgotPassword)
	router.Methods("POST").Path("/login/resetpassword").HandlerFunc(service.ResetPassword)
	router.Methods("GET").Path("/login/organizationinvitation/{code}").HandlerFunc(service.GetOrganizationInvitation)
	//Authorize form
	router.Methods("GET").Path("/authorize").HandlerFunc(service.ShowAuthorizeForm)
	//Facebook callback
	router.Methods("GET").Path("/facebook_callback").HandlerFunc(service.FacebookCallback)
	//Github callback
	router.Methods("GET").Path("/github_callback").HandlerFunc(service.GithubCallback)
	//Logout link
	router.Methods("GET").Path("/logout").HandlerFunc(service.Logout)
	//Error page
	router.Methods("GET").Path("/error").HandlerFunc(service.ErrorPage)
	router.Methods("GET").Path("/error{errornumber}").HandlerFunc(service.ErrorPage)
	router.Methods("GET").Path("/config").HandlerFunc(service.GetConfig)

	router.Methods("GET").Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Version string
		}{
			Version: service.version,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&response)
	})

	router.Methods("GET").Path("/location").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the header from the cloudflare IP Geolocation service. Value is the
		// country code in ISO 3166-1 Alpha 2 format.
		location := r.Header.Get("CF-IPCountry")
		log.Debug("request location: ", location)
		response := struct {
			Location string `json:"location"`
		}{
			Location: location,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&response)
	})

	//host the assets used in the htmlpages
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(
		&assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assets.AssetInfo})))
	router.PathPrefix("/thirdpartyassets/").Handler(http.StripPrefix("/thirdpartyassets/", http.FileServer(
		&assetfs.AssetFS{Asset: thirdpartyassets.Asset, AssetDir: thirdpartyassets.AssetDir, AssetInfo: thirdpartyassets.AssetInfo})))
	router.PathPrefix("/components/").Handler(http.StripPrefix("/components/", http.FileServer(
		&assetfs.AssetFS{Asset: components.Asset, AssetDir: components.AssetDir, AssetInfo: components.AssetInfo})))

	//host the apidocumentation
	router.Methods("GET").Path("/apidocumentation").HandlerFunc(service.APIDocs)
	router.PathPrefix("/apidocumentation/raml/").Handler(http.StripPrefix("/apidocumentation/raml", http.FileServer(
		&assetfs.AssetFS{Asset: specifications.Asset, AssetDir: specifications.AssetDir, AssetInfo: specifications.AssetInfo})))
	router.PathPrefix("/apidocumentation/").Handler(http.StripPrefix("/apidocumentation/", http.FileServer(
		&assetfs.AssetFS{Asset: apiconsole.Asset, AssetDir: apiconsole.AssetDir, AssetInfo: apiconsole.AssetInfo})))

}
