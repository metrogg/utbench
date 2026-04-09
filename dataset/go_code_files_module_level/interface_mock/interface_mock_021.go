func AddRoutes(router *mux.Router) {
	var route, listRoute, inspectRoute string

	// Register aciGw
	route = "/api/v1/aciGws/{key}/"
	listRoute = "/api/v1/aciGws/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListAciGws))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetAciGw))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateAciGw))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateAciGw))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteAciGw))

	inspectRoute = "/api/v1/inspect/aciGws/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectAciGw))

	// Register appProfile
	route = "/api/v1/appProfiles/{key}/"
	listRoute = "/api/v1/appProfiles/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListAppProfiles))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetAppProfile))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateAppProfile))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateAppProfile))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteAppProfile))

	inspectRoute = "/api/v1/inspect/appProfiles/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectAppProfile))

	// Register Bgp
	route = "/api/v1/Bgps/{key}/"
	listRoute = "/api/v1/Bgps/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListBgps))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetBgp))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateBgp))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateBgp))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteBgp))

	inspectRoute = "/api/v1/inspect/Bgps/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectBgp))

	inspectRoute = "/api/v1/inspect/endpoints/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectEndpoint))

	// Register endpointGroup
	route = "/api/v1/endpointGroups/{key}/"
	listRoute = "/api/v1/endpointGroups/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListEndpointGroups))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetEndpointGroup))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateEndpointGroup))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateEndpointGroup))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteEndpointGroup))

	inspectRoute = "/api/v1/inspect/endpointGroups/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectEndpointGroup))

	// Register extContractsGroup
	route = "/api/v1/extContractsGroups/{key}/"
	listRoute = "/api/v1/extContractsGroups/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListExtContractsGroups))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetExtContractsGroup))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateExtContractsGroup))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateExtContractsGroup))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteExtContractsGroup))

	inspectRoute = "/api/v1/inspect/extContractsGroups/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectExtContractsGroup))

	// Register global
	route = "/api/v1/globals/{key}/"
	listRoute = "/api/v1/globals/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListGlobals))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetGlobal))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateGlobal))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateGlobal))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteGlobal))

	inspectRoute = "/api/v1/inspect/globals/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectGlobal))

	// Register netprofile
	route = "/api/v1/netprofiles/{key}/"
	listRoute = "/api/v1/netprofiles/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListNetprofiles))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetNetprofile))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateNetprofile))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateNetprofile))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteNetprofile))

	inspectRoute = "/api/v1/inspect/netprofiles/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectNetprofile))

	// Register network
	route = "/api/v1/networks/{key}/"
	listRoute = "/api/v1/networks/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListNetworks))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetNetwork))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateNetwork))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateNetwork))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteNetwork))

	inspectRoute = "/api/v1/inspect/networks/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectNetwork))

	// Register policy
	route = "/api/v1/policys/{key}/"
	listRoute = "/api/v1/policys/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListPolicys))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetPolicy))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreatePolicy))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreatePolicy))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeletePolicy))

	inspectRoute = "/api/v1/inspect/policys/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectPolicy))

	// Register rule
	route = "/api/v1/rules/{key}/"
	listRoute = "/api/v1/rules/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListRules))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetRule))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateRule))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateRule))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteRule))

	inspectRoute = "/api/v1/inspect/rules/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectRule))

	// Register serviceLB
	route = "/api/v1/serviceLBs/{key}/"
	listRoute = "/api/v1/serviceLBs/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListServiceLBs))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetServiceLB))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateServiceLB))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateServiceLB))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteServiceLB))

	inspectRoute = "/api/v1/inspect/serviceLBs/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectServiceLB))

	// Register tenant
	route = "/api/v1/tenants/{key}/"
	listRoute = "/api/v1/tenants/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListTenants))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetTenant))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateTenant))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateTenant))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteTenant))

	inspectRoute = "/api/v1/inspect/tenants/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectTenant))

	// Register volume
	route = "/api/v1/volumes/{key}/"
	listRoute = "/api/v1/volumes/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListVolumes))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetVolume))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateVolume))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateVolume))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteVolume))

	inspectRoute = "/api/v1/inspect/volumes/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectVolume))

	// Register volumeProfile
	route = "/api/v1/volumeProfiles/{key}/"
	listRoute = "/api/v1/volumeProfiles/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpListVolumeProfiles))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGetVolumeProfile))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreateVolumeProfile))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreateVolumeProfile))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDeleteVolumeProfile))

	inspectRoute = "/api/v1/inspect/volumeProfiles/{key}/"
	router.Path(inspectRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpInspectVolumeProfile))

}
