func GetCheckerForBuiltinRole(clusterName string, clusterConfig services.ClusterConfig, role teleport.Role) (services.RoleSet, error) {
	switch role {
	case teleport.RoleAuth:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Rules: []services.Rule{
						services.NewRule(services.KindAuthServer, services.RW()),
					},
				},
			})
	case teleport.RoleProvisionToken:
		return services.FromSpec(role.String(), services.RoleSpecV3{})
	case teleport.RoleNode:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Rules: []services.Rule{
						services.NewRule(services.KindNode, services.RW()),
						services.NewRule(services.KindSSHSession, services.RW()),
						services.NewRule(services.KindEvent, services.RW()),
						services.NewRule(services.KindProxy, services.RO()),
						services.NewRule(services.KindCertAuthority, services.ReadNoSecrets()),
						services.NewRule(services.KindUser, services.RO()),
						services.NewRule(services.KindNamespace, services.RO()),
						services.NewRule(services.KindRole, services.RO()),
						services.NewRule(services.KindAuthServer, services.RO()),
						services.NewRule(services.KindReverseTunnel, services.RW()),
						services.NewRule(services.KindTunnelConnection, services.RO()),
						services.NewRule(services.KindClusterConfig, services.RO()),
					},
				},
			})
	case teleport.RoleProxy:
		// if in recording mode, return a different set of permissions than regular
		// mode. recording proxy needs to be able to generate host certificates.
		if clusterConfig.GetSessionRecording() == services.RecordAtProxy {
			return services.FromSpec(
				role.String(),
				services.RoleSpecV3{
					Allow: services.RoleConditions{
						Namespaces: []string{services.Wildcard},
						Rules: []services.Rule{
							services.NewRule(services.KindProxy, services.RW()),
							services.NewRule(services.KindOIDCRequest, services.RW()),
							services.NewRule(services.KindSSHSession, services.RW()),
							services.NewRule(services.KindSession, services.RO()),
							services.NewRule(services.KindEvent, services.RW()),
							services.NewRule(services.KindSAMLRequest, services.RW()),
							services.NewRule(services.KindOIDC, services.ReadNoSecrets()),
							services.NewRule(services.KindSAML, services.ReadNoSecrets()),
							services.NewRule(services.KindGithub, services.ReadNoSecrets()),
							services.NewRule(services.KindGithubRequest, services.RW()),
							services.NewRule(services.KindNamespace, services.RO()),
							services.NewRule(services.KindNode, services.RO()),
							services.NewRule(services.KindAuthServer, services.RO()),
							services.NewRule(services.KindReverseTunnel, services.RO()),
							services.NewRule(services.KindCertAuthority, services.ReadNoSecrets()),
							services.NewRule(services.KindUser, services.RO()),
							services.NewRule(services.KindRole, services.RO()),
							services.NewRule(services.KindClusterAuthPreference, services.RO()),
							services.NewRule(services.KindClusterConfig, services.RO()),
							services.NewRule(services.KindClusterName, services.RO()),
							services.NewRule(services.KindStaticTokens, services.RO()),
							services.NewRule(services.KindTunnelConnection, services.RW()),
							services.NewRule(services.KindHostCert, services.RW()),
							services.NewRule(services.KindRemoteCluster, services.RO()),
							// this rule allows local proxy to update the remote cluster's host certificate authorities
							// during certificates renewal
							{
								Resources: []string{services.KindCertAuthority},
								Verbs:     []string{services.VerbCreate, services.VerbRead, services.VerbUpdate},
								// allow administrative access to the host certificate authorities
								// matching any cluster name except local
								Where: builder.And(
									builder.Equals(services.CertAuthorityTypeExpr, builder.String(string(services.HostCA))),
									builder.Not(
										builder.Equals(
											services.ResourceNameExpr,
											builder.String(clusterName),
										),
									),
								).String(),
							},
						},
					},
				})
		}
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Rules: []services.Rule{
						services.NewRule(services.KindProxy, services.RW()),
						services.NewRule(services.KindOIDCRequest, services.RW()),
						services.NewRule(services.KindSSHSession, services.RW()),
						services.NewRule(services.KindSession, services.RO()),
						services.NewRule(services.KindEvent, services.RW()),
						services.NewRule(services.KindSAMLRequest, services.RW()),
						services.NewRule(services.KindOIDC, services.ReadNoSecrets()),
						services.NewRule(services.KindSAML, services.ReadNoSecrets()),
						services.NewRule(services.KindGithub, services.ReadNoSecrets()),
						services.NewRule(services.KindGithubRequest, services.RW()),
						services.NewRule(services.KindNamespace, services.RO()),
						services.NewRule(services.KindNode, services.RO()),
						services.NewRule(services.KindAuthServer, services.RO()),
						services.NewRule(services.KindReverseTunnel, services.RO()),
						services.NewRule(services.KindCertAuthority, services.ReadNoSecrets()),
						services.NewRule(services.KindUser, services.RO()),
						services.NewRule(services.KindRole, services.RO()),
						services.NewRule(services.KindClusterAuthPreference, services.RO()),
						services.NewRule(services.KindClusterConfig, services.RO()),
						services.NewRule(services.KindClusterName, services.RO()),
						services.NewRule(services.KindStaticTokens, services.RO()),
						services.NewRule(services.KindTunnelConnection, services.RW()),
						services.NewRule(services.KindRemoteCluster, services.RO()),
						// this rule allows local proxy to update the remote cluster's host certificate authorities
						// during certificates renewal
						{
							Resources: []string{services.KindCertAuthority},
							Verbs:     []string{services.VerbCreate, services.VerbRead, services.VerbUpdate},
							// allow administrative access to the certificate authority names
							// matching any cluster name except local
							Where: builder.And(
								builder.Equals(services.CertAuthorityTypeExpr, builder.String(string(services.HostCA))),
								builder.Not(
									builder.Equals(
										services.ResourceNameExpr,
										builder.String(clusterName),
									),
								),
							).String(),
						},
					},
				},
			})
	case teleport.RoleWeb:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Rules: []services.Rule{
						services.NewRule(services.KindWebSession, services.RW()),
						services.NewRule(services.KindSSHSession, services.RW()),
						services.NewRule(services.KindAuthServer, services.RO()),
						services.NewRule(services.KindUser, services.RO()),
						services.NewRule(services.KindRole, services.RO()),
						services.NewRule(services.KindNamespace, services.RO()),
						services.NewRule(services.KindTrustedCluster, services.RO()),
					},
				},
			})
	case teleport.RoleSignup:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Rules: []services.Rule{
						services.NewRule(services.KindAuthServer, services.RO()),
						services.NewRule(services.KindClusterAuthPreference, services.RO()),
					},
				},
			})
	case teleport.RoleAdmin:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Options: services.RoleOptions{
					MaxSessionTTL: services.MaxDuration(),
				},
				Allow: services.RoleConditions{
					Namespaces: []string{services.Wildcard},
					Logins:     []string{},
					NodeLabels: services.Labels{services.Wildcard: []string{services.Wildcard}},
					Rules: []services.Rule{
						services.NewRule(services.Wildcard, services.RW()),
					},
				},
			})
	case teleport.RoleNop:
		return services.FromSpec(
			role.String(),
			services.RoleSpecV3{
				Allow: services.RoleConditions{
					Namespaces: []string{},
					Rules:      []services.Rule{},
				},
			})
	}

	return nil, trace.NotFound("%v is not reconginzed", role.String())
}
