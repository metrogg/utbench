func (s *HTTPServer) registerHandlers(enableDebug bool) {
	s.mux.HandleFunc("/v1/jobs", s.wrap(s.JobsRequest))
	s.mux.HandleFunc("/v1/jobs/parse", s.wrap(s.JobsParseRequest))
	s.mux.HandleFunc("/v1/job/", s.wrap(s.JobSpecificRequest))

	s.mux.HandleFunc("/v1/nodes", s.wrap(s.NodesRequest))
	s.mux.HandleFunc("/v1/node/", s.wrap(s.NodeSpecificRequest))

	s.mux.HandleFunc("/v1/allocations", s.wrap(s.AllocsRequest))
	s.mux.HandleFunc("/v1/allocation/", s.wrap(s.AllocSpecificRequest))

	s.mux.HandleFunc("/v1/evaluations", s.wrap(s.EvalsRequest))
	s.mux.HandleFunc("/v1/evaluation/", s.wrap(s.EvalSpecificRequest))

	s.mux.HandleFunc("/v1/deployments", s.wrap(s.DeploymentsRequest))
	s.mux.HandleFunc("/v1/deployment/", s.wrap(s.DeploymentSpecificRequest))

	s.mux.HandleFunc("/v1/acl/policies", s.wrap(s.ACLPoliciesRequest))
	s.mux.HandleFunc("/v1/acl/policy/", s.wrap(s.ACLPolicySpecificRequest))

	s.mux.HandleFunc("/v1/acl/bootstrap", s.wrap(s.ACLTokenBootstrap))
	s.mux.HandleFunc("/v1/acl/tokens", s.wrap(s.ACLTokensRequest))
	s.mux.HandleFunc("/v1/acl/token", s.wrap(s.ACLTokenSpecificRequest))
	s.mux.HandleFunc("/v1/acl/token/", s.wrap(s.ACLTokenSpecificRequest))

	s.mux.Handle("/v1/client/fs/", wrapCORS(s.wrap(s.FsRequest)))
	s.mux.HandleFunc("/v1/client/gc", s.wrap(s.ClientGCRequest))
	s.mux.Handle("/v1/client/stats", wrapCORS(s.wrap(s.ClientStatsRequest)))
	s.mux.Handle("/v1/client/allocation/", wrapCORS(s.wrap(s.ClientAllocRequest)))

	s.mux.HandleFunc("/v1/agent/self", s.wrap(s.AgentSelfRequest))
	s.mux.HandleFunc("/v1/agent/join", s.wrap(s.AgentJoinRequest))
	s.mux.HandleFunc("/v1/agent/members", s.wrap(s.AgentMembersRequest))
	s.mux.HandleFunc("/v1/agent/force-leave", s.wrap(s.AgentForceLeaveRequest))
	s.mux.HandleFunc("/v1/agent/servers", s.wrap(s.AgentServersRequest))
	s.mux.HandleFunc("/v1/agent/keyring/", s.wrap(s.KeyringOperationRequest))
	s.mux.HandleFunc("/v1/agent/health", s.wrap(s.HealthRequest))

	s.mux.HandleFunc("/v1/metrics", s.wrap(s.MetricsRequest))

	s.mux.HandleFunc("/v1/validate/job", s.wrap(s.ValidateJobRequest))

	s.mux.HandleFunc("/v1/regions", s.wrap(s.RegionListRequest))

	s.mux.HandleFunc("/v1/status/leader", s.wrap(s.StatusLeaderRequest))
	s.mux.HandleFunc("/v1/status/peers", s.wrap(s.StatusPeersRequest))

	s.mux.HandleFunc("/v1/search", s.wrap(s.SearchRequest))

	s.mux.HandleFunc("/v1/operator/raft/", s.wrap(s.OperatorRequest))
	s.mux.HandleFunc("/v1/operator/autopilot/configuration", s.wrap(s.OperatorAutopilotConfiguration))
	s.mux.HandleFunc("/v1/operator/autopilot/health", s.wrap(s.OperatorServerHealth))

	s.mux.HandleFunc("/v1/system/gc", s.wrap(s.GarbageCollectRequest))
	s.mux.HandleFunc("/v1/system/reconcile/summaries", s.wrap(s.ReconcileJobSummaries))

	s.mux.HandleFunc("/v1/operator/scheduler/configuration", s.wrap(s.OperatorSchedulerConfiguration))

	if uiEnabled {
		s.mux.Handle("/ui/", http.StripPrefix("/ui/", handleUI(http.FileServer(&UIAssetWrapper{FileSystem: assetFS()}))))
	} else {
		// Write the stubHTML
		s.mux.HandleFunc("/ui/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(stubHTML))
		})
	}
	s.mux.Handle("/", handleRootRedirect())

	if enableDebug {
		s.mux.HandleFunc("/debug/pprof/", pprof.Index)
		s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Register enterprise endpoints.
	s.registerEnterpriseHandlers()
}
