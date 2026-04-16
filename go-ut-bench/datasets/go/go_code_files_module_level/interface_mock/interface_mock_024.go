func NewCiliumAPI(spec *loads.Document) *CiliumAPI {
	return &CiliumAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		EndpointDeleteEndpointIDHandler: endpoint.DeleteEndpointIDHandlerFunc(func(params endpoint.DeleteEndpointIDParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointDeleteEndpointID has not yet been implemented")
		}),
		PolicyDeleteFqdnCacheHandler: policy.DeleteFqdnCacheHandlerFunc(func(params policy.DeleteFqdnCacheParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyDeleteFqdnCache has not yet been implemented")
		}),
		IPAMDeleteIPAMIPHandler: ipam.DeleteIPAMIPHandlerFunc(func(params ipam.DeleteIPAMIPParams) middleware.Responder {
			return middleware.NotImplemented("operation IPAMDeleteIPAMIP has not yet been implemented")
		}),
		PolicyDeletePolicyHandler: policy.DeletePolicyHandlerFunc(func(params policy.DeletePolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyDeletePolicy has not yet been implemented")
		}),
		ServiceDeleteServiceIDHandler: service.DeleteServiceIDHandlerFunc(func(params service.DeleteServiceIDParams) middleware.Responder {
			return middleware.NotImplemented("operation ServiceDeleteServiceID has not yet been implemented")
		}),
		DaemonGetConfigHandler: daemon.GetConfigHandlerFunc(func(params daemon.GetConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonGetConfig has not yet been implemented")
		}),
		DaemonGetDebuginfoHandler: daemon.GetDebuginfoHandlerFunc(func(params daemon.GetDebuginfoParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonGetDebuginfo has not yet been implemented")
		}),
		EndpointGetEndpointHandler: endpoint.GetEndpointHandlerFunc(func(params endpoint.GetEndpointParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpoint has not yet been implemented")
		}),
		EndpointGetEndpointIDHandler: endpoint.GetEndpointIDHandlerFunc(func(params endpoint.GetEndpointIDParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpointID has not yet been implemented")
		}),
		EndpointGetEndpointIDConfigHandler: endpoint.GetEndpointIDConfigHandlerFunc(func(params endpoint.GetEndpointIDConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpointIDConfig has not yet been implemented")
		}),
		EndpointGetEndpointIDHealthzHandler: endpoint.GetEndpointIDHealthzHandlerFunc(func(params endpoint.GetEndpointIDHealthzParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpointIDHealthz has not yet been implemented")
		}),
		EndpointGetEndpointIDLabelsHandler: endpoint.GetEndpointIDLabelsHandlerFunc(func(params endpoint.GetEndpointIDLabelsParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpointIDLabels has not yet been implemented")
		}),
		EndpointGetEndpointIDLogHandler: endpoint.GetEndpointIDLogHandlerFunc(func(params endpoint.GetEndpointIDLogParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointGetEndpointIDLog has not yet been implemented")
		}),
		PolicyGetFqdnCacheHandler: policy.GetFqdnCacheHandlerFunc(func(params policy.GetFqdnCacheParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetFqdnCache has not yet been implemented")
		}),
		PolicyGetFqdnCacheIDHandler: policy.GetFqdnCacheIDHandlerFunc(func(params policy.GetFqdnCacheIDParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetFqdnCacheID has not yet been implemented")
		}),
		DaemonGetHealthzHandler: daemon.GetHealthzHandlerFunc(func(params daemon.GetHealthzParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonGetHealthz has not yet been implemented")
		}),
		PolicyGetIdentityHandler: policy.GetIdentityHandlerFunc(func(params policy.GetIdentityParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetIdentity has not yet been implemented")
		}),
		PolicyGetIdentityEndpointsHandler: policy.GetIdentityEndpointsHandlerFunc(func(params policy.GetIdentityEndpointsParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetIdentityEndpoints has not yet been implemented")
		}),
		PolicyGetIdentityIDHandler: policy.GetIdentityIDHandlerFunc(func(params policy.GetIdentityIDParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetIdentityID has not yet been implemented")
		}),
		DaemonGetMapHandler: daemon.GetMapHandlerFunc(func(params daemon.GetMapParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonGetMap has not yet been implemented")
		}),
		DaemonGetMapNameHandler: daemon.GetMapNameHandlerFunc(func(params daemon.GetMapNameParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonGetMapName has not yet been implemented")
		}),
		MetricsGetMetricsHandler: metrics.GetMetricsHandlerFunc(func(params metrics.GetMetricsParams) middleware.Responder {
			return middleware.NotImplemented("operation MetricsGetMetrics has not yet been implemented")
		}),
		PolicyGetPolicyHandler: policy.GetPolicyHandlerFunc(func(params policy.GetPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetPolicy has not yet been implemented")
		}),
		PolicyGetPolicyResolveHandler: policy.GetPolicyResolveHandlerFunc(func(params policy.GetPolicyResolveParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyGetPolicyResolve has not yet been implemented")
		}),
		PrefilterGetPrefilterHandler: prefilter.GetPrefilterHandlerFunc(func(params prefilter.GetPrefilterParams) middleware.Responder {
			return middleware.NotImplemented("operation PrefilterGetPrefilter has not yet been implemented")
		}),
		ServiceGetServiceHandler: service.GetServiceHandlerFunc(func(params service.GetServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation ServiceGetService has not yet been implemented")
		}),
		ServiceGetServiceIDHandler: service.GetServiceIDHandlerFunc(func(params service.GetServiceIDParams) middleware.Responder {
			return middleware.NotImplemented("operation ServiceGetServiceID has not yet been implemented")
		}),
		DaemonPatchConfigHandler: daemon.PatchConfigHandlerFunc(func(params daemon.PatchConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation DaemonPatchConfig has not yet been implemented")
		}),
		EndpointPatchEndpointIDHandler: endpoint.PatchEndpointIDHandlerFunc(func(params endpoint.PatchEndpointIDParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointPatchEndpointID has not yet been implemented")
		}),
		EndpointPatchEndpointIDConfigHandler: endpoint.PatchEndpointIDConfigHandlerFunc(func(params endpoint.PatchEndpointIDConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointPatchEndpointIDConfig has not yet been implemented")
		}),
		EndpointPatchEndpointIDLabelsHandler: endpoint.PatchEndpointIDLabelsHandlerFunc(func(params endpoint.PatchEndpointIDLabelsParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointPatchEndpointIDLabels has not yet been implemented")
		}),
		PrefilterPatchPrefilterHandler: prefilter.PatchPrefilterHandlerFunc(func(params prefilter.PatchPrefilterParams) middleware.Responder {
			return middleware.NotImplemented("operation PrefilterPatchPrefilter has not yet been implemented")
		}),
		IPAMPostIPAMHandler: ipam.PostIPAMHandlerFunc(func(params ipam.PostIPAMParams) middleware.Responder {
			return middleware.NotImplemented("operation IPAMPostIPAM has not yet been implemented")
		}),
		IPAMPostIPAMIPHandler: ipam.PostIPAMIPHandlerFunc(func(params ipam.PostIPAMIPParams) middleware.Responder {
			return middleware.NotImplemented("operation IPAMPostIPAMIP has not yet been implemented")
		}),
		EndpointPutEndpointIDHandler: endpoint.PutEndpointIDHandlerFunc(func(params endpoint.PutEndpointIDParams) middleware.Responder {
			return middleware.NotImplemented("operation EndpointPutEndpointID has not yet been implemented")
		}),
		PolicyPutPolicyHandler: policy.PutPolicyHandlerFunc(func(params policy.PutPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation PolicyPutPolicy has not yet been implemented")
		}),
		ServicePutServiceIDHandler: service.PutServiceIDHandlerFunc(func(params service.PutServiceIDParams) middleware.Responder {
			return middleware.NotImplemented("operation ServicePutServiceID has not yet been implemented")
		}),
	}
}
