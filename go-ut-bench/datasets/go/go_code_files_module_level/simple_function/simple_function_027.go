func (d *ProjectStatusDescriber) Describe(namespace, name string) (string, error) {
	var f formatter = namespacedFormatter{}

	g, forbiddenResources, err := d.MakeGraph(namespace)
	if err != nil {
		return "", err
	}

	allNamespaces := namespace == metav1.NamespaceAll
	var project *projectv1.Project
	if !allNamespaces {
		p, err := d.ProjectClient.Projects().Get(namespace, metav1.GetOptions{})
		if err != nil {
			// a forbidden error here (without a --namespace value) means that
			// the user has not created any projects, and is therefore using a
			// default namespace that they cannot list projects from.
			if kapierrors.IsForbidden(err) && len(d.RequestedNamespace) == 0 && len(d.CurrentNamespace) == 0 {
				return loginerrors.NoProjectsExistMessage(d.CanRequestProjects, d.CommandBaseName), nil
			}
			if !kapierrors.IsNotFound(err) {
				return "", err
			}
			p = &projectv1.Project{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
		}
		project = p
		f = namespacedFormatter{currentNamespace: namespace}
	}

	coveredNodes := graphview.IntSet{}

	allServices, coveredByServices := graphview.AllServiceGroups(g, coveredNodes)
	coveredNodes.Insert(coveredByServices.List()...)

	// services grouped by selector
	servicesBySelector := map[string][]graphview.ServiceGroup{}
	services := []graphview.ServiceGroup{}

	// group services with identical selectors
	for _, svc := range allServices {
		selector := createSelector(svc.Service.Spec.Selector)
		if _, seen := servicesBySelector[selector.String()]; seen {
			servicesBySelector[selector.String()] = append(servicesBySelector[selector.String()], svc)
			continue
		}

		services = append(services, svc)
		servicesBySelector[selector.String()] = []graphview.ServiceGroup{}
	}

	standaloneDCs, coveredByDCs := graphview.AllDeploymentConfigPipelines(g, coveredNodes)
	coveredNodes.Insert(coveredByDCs.List()...)

	standaloneDeployments, coveredByDeployments := graphview.AllDeployments(g, coveredNodes)
	coveredNodes.Insert(coveredByDeployments.List()...)

	standaloneStatefulSets, coveredByStatefulSets := graphview.AllStatefulSets(g, coveredNodes)
	coveredNodes.Insert(coveredByStatefulSets.List()...)

	standaloneRCs, coveredByRCs := graphview.AllReplicationControllers(g, coveredNodes)
	coveredNodes.Insert(coveredByRCs.List()...)

	standaloneRSs, coveredByRSs := graphview.AllReplicaSets(g, coveredNodes)
	coveredNodes.Insert(coveredByRSs.List()...)

	standaloneImages, coveredByImages := graphview.AllImagePipelinesFromBuildConfig(g, coveredNodes)
	coveredNodes.Insert(coveredByImages.List()...)

	standaloneDaemonSets, coveredByDaemonSets := graphview.AllDaemonSets(g, coveredNodes)
	coveredNodes.Insert(coveredByDaemonSets.List()...)

	standaloneJobs, coveredByJobs := graphview.AllJobs(g, coveredNodes)
	coveredNodes.Insert(coveredByJobs.List()...)

	standalonePods, coveredByPods := graphview.AllPods(g, coveredNodes)
	coveredNodes.Insert(coveredByPods.List()...)

	return tabbedString(func(out *tabwriter.Writer) error {
		indent := "  "
		if allNamespaces {
			fmt.Fprintf(out, describeAllProjectsOnServer(f, d.Server))
		} else {
			fmt.Fprintf(out, describeProjectAndServer(f, project, d.Server))
		}

		for _, service := range services {
			if !service.Service.Found() {
				continue
			}
			local := namespacedFormatter{currentNamespace: service.Service.Namespace}

			var exposes []string
			for _, routeNode := range service.ExposingRoutes {
				exposes = append(exposes, describeRouteInServiceGroup(local, routeNode)...)
			}
			sort.Sort(exposedRoutes(exposes))

			fmt.Fprintln(out)

			// print services that should be grouped with this service based on matching selectors
			selector := createSelector(service.Service.Spec.Selector)
			groupedServices := servicesBySelector[selector.String()]
			for _, groupedSvc := range groupedServices {
				if !groupedSvc.Service.Found() {
					continue
				}

				grouppedLocal := namespacedFormatter{currentNamespace: service.Service.Namespace}

				var grouppedExposes []string
				for _, routeNode := range groupedSvc.ExposingRoutes {
					grouppedExposes = append(grouppedExposes, describeRouteInServiceGroup(grouppedLocal, routeNode)...)
				}
				sort.Sort(exposedRoutes(grouppedExposes))

				printLines(out, "", 0, describeServiceInServiceGroup(f, groupedSvc, grouppedExposes...)...)
			}

			printLines(out, "", 0, describeServiceInServiceGroup(f, service, exposes...)...)

			for _, dcPipeline := range service.DeploymentConfigPipelines {
				printLines(out, indent, 1, describeDeploymentConfigInServiceGroup(local, dcPipeline, func(rc *kubegraph.ReplicationControllerNode) int32 {
					return graphview.MaxRecentContainerRestartsForRC(g, rc)
				})...)
			}

			for _, node := range service.StatefulSets {
				printLines(out, indent, 1, describeStatefulSetInServiceGroup(local, node)...)
			}

			for _, node := range service.Deployments {
				printLines(out, indent, 1, describeDeploymentInServiceGroup(local, node, func(rs *kubegraph.ReplicaSetNode) int32 {
					return graphview.MaxRecentContainerRestartsForRS(g, rs)
				})...)
			}

			for _, node := range service.DaemonSets {
				printLines(out, indent, 1, describeDaemonSetInServiceGroup(local, node)...)
			}

		rsNode:
			for _, rsNode := range service.FulfillingRSs {
				for _, coveredD := range service.FulfillingDeployments {
					if kubeedges.BelongsToDeployment(coveredD.Deployment, rsNode.ReplicaSet) {
						continue rsNode
					}
				}
				printLines(out, indent, 1, describeRSInServiceGroup(local, rsNode)...)
			}

		rcNode:
			for _, rcNode := range service.FulfillingRCs {
				for _, coveredDC := range service.FulfillingDCs {
					if appsedges.BelongsToDeploymentConfig(coveredDC.DeploymentConfig, rcNode.ReplicationController) {
						continue rcNode
					}
				}
				printLines(out, indent, 1, describeRCInServiceGroup(local, rcNode)...)
			}

		pod:
			for _, node := range service.FulfillingPods {
				// skip pods that have been displayed in a roll-up of RCs and DCs (by implicit usage of RCs)
				for _, coveredRC := range service.FulfillingRCs {
					if g.Edge(node, coveredRC) != nil {
						continue pod
					}
				}
				for _, coveredRS := range service.FulfillingRSs {
					if g.Edge(node, coveredRS) != nil {
						continue pod
					}
				}
				// TODO: collapse into FulfillingControllers
				for _, covered := range service.FulfillingStatefulSets {
					if g.Edge(node, covered) != nil {
						continue pod
					}
				}
				printLines(out, indent, 1, describePodInServiceGroup(local, node)...)
			}
		}

		for _, standaloneDC := range standaloneDCs {
			if !standaloneDC.DeploymentConfig.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeDeploymentConfigInServiceGroup(f, standaloneDC, func(rc *kubegraph.ReplicationControllerNode) int32 {
				return graphview.MaxRecentContainerRestartsForRC(g, rc)
			})...)
		}
		for _, standaloneDeployment := range standaloneDeployments {
			if !standaloneDeployment.Deployment.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeDeploymentInServiceGroup(f, standaloneDeployment, func(rs *kubegraph.ReplicaSetNode) int32 {
				return graphview.MaxRecentContainerRestartsForRS(g, rs)
			})...)
		}

		for _, standaloneStatefulSet := range standaloneStatefulSets {
			if !standaloneStatefulSet.StatefulSet.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeStatefulSetInServiceGroup(f, standaloneStatefulSet)...)
		}

		for _, standaloneImage := range standaloneImages {
			fmt.Fprintln(out)
			lines := describeStandaloneBuildGroup(f, standaloneImage, namespace)
			lines = append(lines, describeAdditionalBuildDetail(standaloneImage.Build, standaloneImage.LastSuccessfulBuild, standaloneImage.LastUnsuccessfulBuild, standaloneImage.ActiveBuilds, standaloneImage.DestinationResolved, true)...)
			printLines(out, indent, 0, lines...)
		}

		for _, standaloneRC := range standaloneRCs {
			if !standaloneRC.RC.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeRCInServiceGroup(f, standaloneRC.RC)...)
		}

		for _, standaloneRS := range standaloneRSs {
			if !standaloneRS.RS.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeRSInServiceGroup(f, standaloneRS.RS)...)
		}

		for _, standaloneDaemonSet := range standaloneDaemonSets {
			if !standaloneDaemonSet.DaemonSet.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeDaemonSetInServiceGroup(f, standaloneDaemonSet)...)
		}

		for _, standaloneJob := range standaloneJobs {
			if !standaloneJob.Job.Found() {
				continue
			}

			fmt.Fprintln(out)
			printLines(out, indent, 0, describeStandaloneJob(f, standaloneJob)...)
		}

		monopods, err := filterBoringPods(standalonePods)
		if err != nil {
			return err
		}
		for _, monopod := range monopods {
			fmt.Fprintln(out)
			printLines(out, indent, 0, describeMonopod(f, monopod.Pod)...)
		}

		allMarkers := osgraph.Markers{}
		allMarkers = append(allMarkers, createForbiddenMarkers(forbiddenResources)...)
		for _, scanner := range getMarkerScanners(d.LogsCommandName, d.SecurityPolicyCommandFormat, d.SetProbeCommandName, forbiddenResources) {
			allMarkers = append(allMarkers, scanner(g, f)...)
		}

		// TODO: Provide an option to chase these hidden markers.
		allMarkers = allMarkers.FilterByNamespace(namespace)

		fmt.Fprintln(out)

		sort.Stable(osgraph.ByKey(allMarkers))
		sort.Stable(osgraph.ByNodeID(allMarkers))

		errorMarkers := allMarkers.BySeverity(osgraph.ErrorSeverity)
		errorSuggestions := 0
		if len(errorMarkers) > 0 {
			fmt.Fprintln(out, "Errors:")
			errorSuggestions += printMarkerSuggestions(errorMarkers, d.Suggest, out, indent)
		}

		warningMarkers := allMarkers.BySeverity(osgraph.WarningSeverity)
		if len(warningMarkers) > 0 {
			if d.Suggest {
				// add linebreak between Errors list and Warnings list
				if len(errorMarkers) > 0 {
					fmt.Fprintln(out)
				}
				fmt.Fprintln(out, "Warnings:")
			}
			printMarkerSuggestions(warningMarkers, d.Suggest, out, indent)
		}

		infoMarkers := allMarkers.BySeverity(osgraph.InfoSeverity)
		if len(infoMarkers) > 0 {
			if d.Suggest {
				// add linebreak between Warnings list and Info List
				if len(warningMarkers) > 0 || len(errorMarkers) > 0 {
					fmt.Fprintln(out)
				}
				fmt.Fprintln(out, "Info:")
			}
			printMarkerSuggestions(infoMarkers, d.Suggest, out, indent)
		}

		// We print errors by default and warnings if --sugest is used. If we get none,
		// this would be an extra new line.
		if len(errorMarkers) != 0 || len(infoMarkers) != 0 || (d.Suggest && len(warningMarkers) != 0) {
			fmt.Fprintln(out)
		}

		errors, warnings, infos := "", "", ""
		if len(errorMarkers) == 1 {
			errors = "1 error"
		} else if len(errorMarkers) > 1 {
			errors = fmt.Sprintf("%d errors", len(errorMarkers))
		}
		if len(warningMarkers) == 1 {
			warnings = "1 warning"
		} else if len(warningMarkers) > 1 {
			warnings = fmt.Sprintf("%d warnings", len(warningMarkers))
		}
		if len(infoMarkers) == 1 {
			infos = "1 info"
		} else if len(infoMarkers) > 0 {
			infos = fmt.Sprintf("%d infos", len(infoMarkers))
		}

		markerStrings := []string{errors, warnings, infos}
		markerString := ""
		count := 0
		for _, m := range markerStrings {
			if len(m) > 0 {
				if count > 0 {
					markerString = fmt.Sprintf("%s, ", markerString)
				}
				markerString = fmt.Sprintf("%s%s", markerString, m)
				count++
			}
		}

		switch {
		case !d.Suggest && ((len(errorMarkers) > 0 && errorSuggestions > 0) || len(warningMarkers) > 0 || len(infoMarkers) > 0):
			fmt.Fprintf(out, "%s identified, use '%s status --suggest' to see details.\n", markerString, d.CommandBaseName)

		case (len(services) == 0) && (len(standaloneDCs) == 0) && (len(standaloneImages) == 0):
			fmt.Fprintln(out, "You have no services, deployment configs, or build configs.")
			fmt.Fprintf(out, "Run '%[1]s new-app' to create an application.\n", d.CommandBaseName)

		default:
			fmt.Fprintf(out, "View details with '%[1]s describe <resource>/<name>' or list everything with '%[1]s get all'.\n", d.CommandBaseName)
		}

		return nil
	})
}
