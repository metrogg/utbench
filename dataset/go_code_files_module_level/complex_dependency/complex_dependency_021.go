func ClusterRoles() []rbacv1.ClusterRole {
	roles := []rbacv1.ClusterRole{
		{
			// a "root" role which can do absolutely anything
			ObjectMeta: metav1.ObjectMeta{Name: "cluster-admin"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("*").Groups("*").Resources("*").RuleOrDie(),
				rbacv1helpers.NewRule("*").URLs("*").RuleOrDie(),
			},
		},
		{
			// a role which provides just enough power to determine if the server is ready and discover API versions for negotiation
			ObjectMeta: metav1.ObjectMeta{Name: "system:discovery"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("get").URLs(
					"/healthz", "/version", "/version/",
					"/openapi", "/openapi/*",
					"/api", "/api/*",
					"/apis", "/apis/*",
				).RuleOrDie(),
			},
		},
		{
			// a role which provides minimal resource access to allow a "normal" user to learn information about themselves
			ObjectMeta: metav1.ObjectMeta{Name: "system:basic-user"},
			Rules: []rbacv1.PolicyRule{
				// TODO add future selfsubjectrulesreview, project request APIs, project listing APIs
				rbacv1helpers.NewRule("create").Groups(authorizationGroup).Resources("selfsubjectaccessreviews", "selfsubjectrulesreviews").RuleOrDie(),
			},
		},
		{
			// a role which provides just enough power read insensitive cluster information
			ObjectMeta: metav1.ObjectMeta{Name: "system:public-info-viewer"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("get").URLs(
					"/healthz", "/version", "/version/",
				).RuleOrDie(),
			},
		},
		{
			// a role for a namespace level admin.  It is `edit` plus the power to grant permissions to other users.
			ObjectMeta: metav1.ObjectMeta{Name: "admin"},
			AggregationRule: &rbacv1.AggregationRule{
				ClusterRoleSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-admin": "true"}},
				},
			},
		},
		{
			// a role for a namespace level editor.  It grants access to all user level actions in a namespace.
			// It does not grant powers for "privileged" resources which are domain of the system: `/status`
			// subresources or `quota`/`limits` which are used to control namespaces
			ObjectMeta: metav1.ObjectMeta{Name: "edit", Labels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-admin": "true"}},
			AggregationRule: &rbacv1.AggregationRule{
				ClusterRoleSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-edit": "true"}},
				},
			},
		},
		{
			// a role for namespace level viewing.  It grants Read-only access to non-escalating resources in
			// a namespace.
			ObjectMeta: metav1.ObjectMeta{Name: "view", Labels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-edit": "true"}},
			AggregationRule: &rbacv1.AggregationRule{
				ClusterRoleSelectors: []metav1.LabelSelector{
					{MatchLabels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-view": "true"}},
				},
			},
		},
		{
			// a role for a namespace level admin.  It is `edit` plus the power to grant permissions to other users.
			ObjectMeta: metav1.ObjectMeta{Name: "system:aggregate-to-admin", Labels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-admin": "true"}},
			Rules: []rbacv1.PolicyRule{
				// additional admin powers
				rbacv1helpers.NewRule("create").Groups(authorizationGroup).Resources("localsubjectaccessreviews").RuleOrDie(),
				rbacv1helpers.NewRule(ReadWrite...).Groups(rbacGroup).Resources("roles", "rolebindings").RuleOrDie(),
			},
		},
		{
			// a role for a namespace level editor.  It grants access to all user level actions in a namespace.
			// It does not grant powers for "privileged" resources which are domain of the system: `/status`
			// subresources or `quota`/`limits` which are used to control namespaces
			ObjectMeta: metav1.ObjectMeta{Name: "system:aggregate-to-edit", Labels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-edit": "true"}},
			Rules: []rbacv1.PolicyRule{
				// Allow read on escalating resources
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("pods/attach", "pods/proxy", "pods/exec", "pods/portforward", "secrets", "services/proxy").RuleOrDie(),
				rbacv1helpers.NewRule("impersonate").Groups(legacyGroup).Resources("serviceaccounts").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(legacyGroup).Resources("pods", "pods/attach", "pods/proxy", "pods/exec", "pods/portforward").RuleOrDie(),
				rbacv1helpers.NewRule(Write...).Groups(legacyGroup).Resources("replicationcontrollers", "replicationcontrollers/scale", "serviceaccounts",
					"services", "services/proxy", "endpoints", "persistentvolumeclaims", "configmaps", "secrets").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(appsGroup).Resources(
					"statefulsets", "statefulsets/scale",
					"daemonsets",
					"deployments", "deployments/scale", "deployments/rollback",
					"replicasets", "replicasets/scale").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(autoscalingGroup).Resources("horizontalpodautoscalers").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(batchGroup).Resources("jobs", "cronjobs").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(extensionsGroup).Resources("daemonsets",
					"deployments", "deployments/scale", "deployments/rollback", "ingresses",
					"replicasets", "replicasets/scale", "replicationcontrollers/scale",
					"networkpolicies").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(policyGroup).Resources("poddisruptionbudgets").RuleOrDie(),

				rbacv1helpers.NewRule(Write...).Groups(networkingGroup).Resources("networkpolicies", "ingresses").RuleOrDie(),
			},
		},
		{
			// a role for namespace level viewing.  It grants Read-only access to non-escalating resources in
			// a namespace.
			ObjectMeta: metav1.ObjectMeta{Name: "system:aggregate-to-view", Labels: map[string]string{"rbac.authorization.k8s.io/aggregate-to-view": "true"}},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("pods", "replicationcontrollers", "replicationcontrollers/scale", "serviceaccounts",
					"services", "endpoints", "persistentvolumeclaims", "configmaps").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("limitranges", "resourcequotas", "bindings", "events",
					"pods/status", "resourcequotas/status", "namespaces/status", "replicationcontrollers/status", "pods/log").RuleOrDie(),
				// read access to namespaces at the namespace scope means you can read *this* namespace.  This can be used as an
				// indicator of which namespaces you have access to.
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("namespaces").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(appsGroup).Resources(
					"controllerrevisions",
					"statefulsets", "statefulsets/scale",
					"daemonsets",
					"deployments", "deployments/scale",
					"replicasets", "replicasets/scale").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(autoscalingGroup).Resources("horizontalpodautoscalers").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(batchGroup).Resources("jobs", "cronjobs").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(extensionsGroup).Resources("daemonsets", "deployments", "deployments/scale",
					"ingresses", "replicasets", "replicasets/scale", "replicationcontrollers/scale",
					"networkpolicies").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(policyGroup).Resources("poddisruptionbudgets").RuleOrDie(),

				rbacv1helpers.NewRule(Read...).Groups(networkingGroup).Resources("networkpolicies", "ingresses").RuleOrDie(),
			},
		},
		{
			// a role to use for heapster's connections back to the API server
			ObjectMeta: metav1.ObjectMeta{Name: "system:heapster"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("events", "pods", "nodes", "namespaces").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(extensionsGroup).Resources("deployments").RuleOrDie(),
			},
		},
		{
			// a role for nodes to use to have the access they need for running pods
			ObjectMeta: metav1.ObjectMeta{Name: "system:node"},
			Rules:      NodeRules(),
		},
		{
			// a role to use for node-problem-detector access.  It does not get bound to default location since
			// deployment locations can reasonably vary.
			ObjectMeta: metav1.ObjectMeta{Name: "system:node-problem-detector"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("get").Groups(legacyGroup).Resources("nodes").RuleOrDie(),
				rbacv1helpers.NewRule("patch").Groups(legacyGroup).Resources("nodes/status").RuleOrDie(),
				eventsRule(),
			},
		},
		{
			// a role to use for setting up a proxy
			ObjectMeta: metav1.ObjectMeta{Name: "system:node-proxier"},
			Rules: []rbacv1.PolicyRule{
				// Used to build serviceLister
				rbacv1helpers.NewRule("list", "watch").Groups(legacyGroup).Resources("services", "endpoints").RuleOrDie(),
				rbacv1helpers.NewRule("get").Groups(legacyGroup).Resources("nodes").RuleOrDie(),

				eventsRule(),
			},
		},
		{
			// a role to use for full access to the kubelet API
			ObjectMeta: metav1.ObjectMeta{Name: "system:kubelet-api-admin"},
			Rules: []rbacv1.PolicyRule{
				// Allow read-only access to the Node API objects
				rbacv1helpers.NewRule("get", "list", "watch").Groups(legacyGroup).Resources("nodes").RuleOrDie(),
				// Allow all API calls to the nodes
				rbacv1helpers.NewRule("proxy").Groups(legacyGroup).Resources("nodes").RuleOrDie(),
				rbacv1helpers.NewRule("*").Groups(legacyGroup).Resources("nodes/proxy", "nodes/metrics", "nodes/spec", "nodes/stats", "nodes/log").RuleOrDie(),
			},
		},
		{
			// a role to use for bootstrapping a node's client certificates
			ObjectMeta: metav1.ObjectMeta{Name: "system:node-bootstrapper"},
			Rules: []rbacv1.PolicyRule{
				// used to create a certificatesigningrequest for a node-specific client certificate, and watch for it to be signed
				rbacv1helpers.NewRule("create", "get", "list", "watch").Groups(certificatesGroup).Resources("certificatesigningrequests").RuleOrDie(),
			},
		},
		{
			// a role to use for allowing authentication and authorization delegation
			ObjectMeta: metav1.ObjectMeta{Name: "system:auth-delegator"},
			Rules: []rbacv1.PolicyRule{
				// These creates are non-mutating
				rbacv1helpers.NewRule("create").Groups(authenticationGroup).Resources("tokenreviews").RuleOrDie(),
				rbacv1helpers.NewRule("create").Groups(authorizationGroup).Resources("subjectaccessreviews").RuleOrDie(),
			},
		},
		{
			// a role to use for the API registry, summarization, and proxy handling
			ObjectMeta: metav1.ObjectMeta{Name: "system:kube-aggregator"},
			Rules: []rbacv1.PolicyRule{
				// it needs to see all services so that it knows whether the ones it points to exist or not
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("services", "endpoints").RuleOrDie(),
			},
		},
		{
			// a role to use for bootstrapping the kube-controller-manager so it can create the shared informers
			// service accounts, and secrets that we need to create separate identities for other controllers
			ObjectMeta: metav1.ObjectMeta{Name: "system:kube-controller-manager"},
			Rules: []rbacv1.PolicyRule{
				eventsRule(),
				rbacv1helpers.NewRule("create").Groups(legacyGroup).Resources("endpoints", "secrets", "serviceaccounts").RuleOrDie(),
				rbacv1helpers.NewRule("delete").Groups(legacyGroup).Resources("secrets").RuleOrDie(),
				rbacv1helpers.NewRule("get").Groups(legacyGroup).Resources("endpoints", "namespaces", "secrets", "serviceaccounts", "configmaps").RuleOrDie(),
				rbacv1helpers.NewRule("update").Groups(legacyGroup).Resources("endpoints", "secrets", "serviceaccounts").RuleOrDie(),
				// Needed to check API access.  These creates are non-mutating
				rbacv1helpers.NewRule("create").Groups(authenticationGroup).Resources("tokenreviews").RuleOrDie(),
				rbacv1helpers.NewRule("create").Groups(authorizationGroup).Resources("subjectaccessreviews").RuleOrDie(),
				// Needed for all shared informers
				rbacv1helpers.NewRule("list", "watch").Groups("*").Resources("*").RuleOrDie(),
				rbacv1helpers.NewRule("create").Groups(legacyGroup).Resources("serviceaccounts/token").RuleOrDie(),
			},
		},
		{
			// a role to use for the kube-scheduler
			ObjectMeta: metav1.ObjectMeta{Name: "system:kube-scheduler"},
			Rules: []rbacv1.PolicyRule{
				eventsRule(),

				// this is for leaderlease access
				// TODO: scope this to the kube-system namespace
				rbacv1helpers.NewRule("create").Groups(legacyGroup).Resources("endpoints").RuleOrDie(),
				rbacv1helpers.NewRule("get", "update", "patch", "delete").Groups(legacyGroup).Resources("endpoints").Names("kube-scheduler").RuleOrDie(),

				// fundamental resources
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("nodes").RuleOrDie(),
				rbacv1helpers.NewRule("get", "list", "watch", "delete").Groups(legacyGroup).Resources("pods").RuleOrDie(),
				rbacv1helpers.NewRule("create").Groups(legacyGroup).Resources("pods/binding", "bindings").RuleOrDie(),
				rbacv1helpers.NewRule("patch", "update").Groups(legacyGroup).Resources("pods/status").RuleOrDie(),
				// things that select pods
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("services", "replicationcontrollers").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(appsGroup, extensionsGroup).Resources("replicasets").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(appsGroup).Resources("statefulsets").RuleOrDie(),
				// things that pods use or applies to them
				rbacv1helpers.NewRule(Read...).Groups(policyGroup).Resources("poddisruptionbudgets").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(legacyGroup).Resources("persistentvolumeclaims", "persistentvolumes").RuleOrDie(),
				// Needed to check API access.  These creates are non-mutating
				rbacv1helpers.NewRule("create").Groups(authenticationGroup).Resources("tokenreviews").RuleOrDie(),
				rbacv1helpers.NewRule("create").Groups(authorizationGroup).Resources("subjectaccessreviews").RuleOrDie(),
			},
		},
		{
			// a role to use for the kube-dns pod
			ObjectMeta: metav1.ObjectMeta{Name: "system:kube-dns"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("list", "watch").Groups(legacyGroup).Resources("endpoints", "services").RuleOrDie(),
			},
		},
		{
			// a role for an external/out-of-tree persistent volume provisioner
			ObjectMeta: metav1.ObjectMeta{Name: "system:persistent-volume-provisioner"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("get", "list", "watch", "create", "delete").Groups(legacyGroup).Resources("persistentvolumes").RuleOrDie(),
				// update is needed in addition to read access for setting lock annotations on PVCs
				rbacv1helpers.NewRule("get", "list", "watch", "update").Groups(legacyGroup).Resources("persistentvolumeclaims").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(storageGroup).Resources("storageclasses").RuleOrDie(),

				// Needed for watching provisioning success and failure events
				rbacv1helpers.NewRule("watch").Groups(legacyGroup).Resources("events").RuleOrDie(),

				eventsRule(),
			},
		},
		{
			// a role for the csi external attacher
			ObjectMeta: metav1.ObjectMeta{Name: "system:csi-external-attacher"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("get", "list", "watch", "update", "patch").Groups(legacyGroup).Resources("persistentvolumes").RuleOrDie(),
				rbacv1helpers.NewRule("get", "list", "watch").Groups(legacyGroup).Resources("nodes").RuleOrDie(),
				rbacv1helpers.NewRule("get", "list", "watch", "update", "patch").Groups(storageGroup).Resources("volumeattachments").RuleOrDie(),
				rbacv1helpers.NewRule("get", "list", "watch", "create", "update", "patch").Groups(legacyGroup).Resources("events").RuleOrDie(),
			},
		},
		{
			// a role making the csrapprover controller approve a node client CSR
			ObjectMeta: metav1.ObjectMeta{Name: "system:certificates.k8s.io:certificatesigningrequests:nodeclient"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("create").Groups(certificatesGroup).Resources("certificatesigningrequests/nodeclient").RuleOrDie(),
			},
		},
		{
			// a role making the csrapprover controller approve a node client CSR requested by the node itself
			ObjectMeta: metav1.ObjectMeta{Name: "system:certificates.k8s.io:certificatesigningrequests:selfnodeclient"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule("create").Groups(certificatesGroup).Resources("certificatesigningrequests/selfnodeclient").RuleOrDie(),
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "system:volume-scheduler"},
			Rules: []rbacv1.PolicyRule{
				rbacv1helpers.NewRule(ReadUpdate...).Groups(legacyGroup).Resources("persistentvolumes").RuleOrDie(),
				rbacv1helpers.NewRule(Read...).Groups(storageGroup).Resources("storageclasses").RuleOrDie(),
				rbacv1helpers.NewRule(ReadUpdate...).Groups(legacyGroup).Resources("persistentvolumeclaims").RuleOrDie(),
			},
		},
	}

	externalProvisionerRules := []rbacv1.PolicyRule{
		rbacv1helpers.NewRule("create", "delete", "get", "list", "watch").Groups(legacyGroup).Resources("persistentvolumes").RuleOrDie(),
		rbacv1helpers.NewRule("get", "list", "watch", "update", "patch").Groups(legacyGroup).Resources("persistentvolumeclaims").RuleOrDie(),
		rbacv1helpers.NewRule("list", "watch").Groups(storageGroup).Resources("storageclasses").RuleOrDie(),
		rbacv1helpers.NewRule("get", "list", "watch", "create", "update", "patch").Groups(legacyGroup).Resources("events").RuleOrDie(),
		rbacv1helpers.NewRule("get", "list", "watch").Groups(legacyGroup).Resources("nodes").RuleOrDie(),
	}
	if utilfeature.DefaultFeatureGate.Enabled(features.CSINodeInfo) {
		externalProvisionerRules = append(externalProvisionerRules, rbacv1helpers.NewRule("get", "watch", "list").Groups("storage.k8s.io").Resources("csinodes").RuleOrDie())
	}
	roles = append(roles, rbacv1.ClusterRole{
		// a role for the csi external provisioner
		ObjectMeta: metav1.ObjectMeta{Name: "system:csi-external-provisioner"},
		Rules:      externalProvisionerRules,
	})

	addClusterRoleLabel(roles)
	return roles
}
