func (gce *Cloud) EnsureLoadBalancer(clusterName string, apiService *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("Cannot EnsureLoadBalancer() with no hosts")
	}

	hostNames := nodeNames(nodes)
	hosts, err := gce.getInstancesByNames(hostNames)
	if err != nil {
		return nil, err
	}

	loadBalancerName := cloudprovider.GetLoadBalancerName(apiService)
	loadBalancerIP := apiService.Spec.LoadBalancerIP
	ports := apiService.Spec.Ports
	portStr := []string{}
	for _, p := range apiService.Spec.Ports {
		portStr = append(portStr, fmt.Sprintf("%s/%d", p.Protocol, p.Port))
	}

	affinityType := apiService.Spec.SessionAffinity

	serviceName := types.NamespacedName{Namespace: apiService.Namespace, Name: apiService.Name}
	glog.V(2).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)", loadBalancerName, gce.region, loadBalancerIP, portStr, hostNames, serviceName, apiService.Annotations)

	// Check if the forwarding rule exists, and if so, what its IP is.
	fwdRuleExists, fwdRuleNeedsUpdate, fwdRuleIP, err := gce.forwardingRuleNeedsUpdate(loadBalancerName, gce.region, loadBalancerIP, ports)
	if err != nil {
		return nil, err
	}
	if !fwdRuleExists {
		glog.Infof("Forwarding rule %v for Service %v/%v doesn't exist", loadBalancerName, apiService.Namespace, apiService.Name)
	}

	// Make sure we know which IP address will be used and have properly reserved
	// it as static before moving forward with the rest of our operations.
	//
	// We use static IP addresses when updating a load balancer to ensure that we
	// can replace the load balancer's other components without changing the
	// address its service is reachable on. We do it this way rather than always
	// keeping the static IP around even though this is more complicated because
	// it makes it less likely that we'll run into quota issues. Only 7 static
	// IP addresses are allowed per region by default.
	//
	// We could let an IP be allocated for us when the forwarding rule is created,
	// but we need the IP to set up the firewall rule, and we want to keep the
	// forwarding rule creation as the last thing that needs to be done in this
	// function in order to maintain the invariant that "if the forwarding rule
	// exists, the LB has been fully created".
	ipAddress := ""

	// Through this process we try to keep track of whether it is safe to
	// release the IP that was allocated.  If the user specifically asked for
	// an IP, we assume they are managing it themselves.  Otherwise, we will
	// release the IP in case of early-terminating failure or upon successful
	// creating of the LB.
	// TODO(#36535): boil this logic down into a set of component functions
	// and key the flag values off of errors returned.
	isUserOwnedIP := false // if this is set, we never release the IP
	isSafeToReleaseIP := false
	defer func() {
		if isUserOwnedIP {
			return
		}
		if isSafeToReleaseIP {
			if err = gce.deleteStaticIP(loadBalancerName, gce.region); err != nil {
				glog.Errorf("failed to release static IP %s for load balancer (%v(%v), %v): %v", ipAddress, loadBalancerName, serviceName, gce.region, err)
			}
			glog.V(2).Infof("EnsureLoadBalancer(%v(%v)): released static IP %s", loadBalancerName, serviceName, ipAddress)
		} else {
			glog.Warningf("orphaning static IP %s during update of load balancer (%v(%v), %v): %v", ipAddress, loadBalancerName, serviceName, gce.region, err)
		}
	}()

	if loadBalancerIP != "" {
		// If a specific IP address has been requested, we have to respect the
		// user's request and use that IP. If the forwarding rule was already using
		// a different IP, it will be harmlessly abandoned because it was only an
		// ephemeral IP (or it was a different static IP owned by the user, in which
		// case we shouldn't delete it anyway).
		if isStatic, addrErr := gce.projectOwnsStaticIP(loadBalancerName, gce.region, loadBalancerIP); err != nil {
			return nil, fmt.Errorf("failed to test if this GCE project owns the static IP %s: %v", loadBalancerIP, addrErr)
		} else if isStatic {
			// The requested IP is a static IP, owned and managed by the user.
			isUserOwnedIP = true
			isSafeToReleaseIP = false
			ipAddress = loadBalancerIP
			glog.V(4).Infof("EnsureLoadBalancer(%v(%v)): using user-provided static IP %s", loadBalancerName, serviceName, ipAddress)
		} else if loadBalancerIP == fwdRuleIP {
			// The requested IP is not a static IP, but is currently assigned
			// to this forwarding rule, so we can keep it.
			isUserOwnedIP = false
			isSafeToReleaseIP = true
			ipAddress, _, err = gce.ensureStaticIP(loadBalancerName, serviceName.String(), gce.region, fwdRuleIP)
			if err != nil {
				return nil, fmt.Errorf("failed to ensure static IP %s: %v", fwdRuleIP, err)
			}
			glog.V(4).Infof("EnsureLoadBalancer(%v(%v)): using user-provided non-static IP %s", loadBalancerName, serviceName, ipAddress)
		} else {
			// The requested IP is not static and it is not assigned to the
			// current forwarding rule.  It might be attached to a different
			// rule or it might not be part of this project at all.  Either
			// way, we can't use it.
			return nil, fmt.Errorf("requested ip %s is neither static nor assigned to LB %s(%v): %v", loadBalancerIP, loadBalancerName, serviceName, err)
		}
	} else {
		// The user did not request a specific IP.
		isUserOwnedIP = false

		// This will either allocate a new static IP if the forwarding rule didn't
		// already have an IP, or it will promote the forwarding rule's current
		// IP from ephemeral to static, or it will just get the IP if it is
		// already static.
		existed := false
		ipAddress, existed, err = gce.ensureStaticIP(loadBalancerName, serviceName.String(), gce.region, fwdRuleIP)
		if err != nil {
			return nil, fmt.Errorf("failed to ensure static IP %s: %v", fwdRuleIP, err)
		}
		if existed {
			// If the IP was not specifically requested by the user, but it
			// already existed, it seems to be a failed update cycle.  We can
			// use this IP and try to run through the process again, but we
			// should not release the IP unless it is explicitly flagged as OK.
			isSafeToReleaseIP = false
			glog.V(4).Infof("EnsureLoadBalancer(%v(%v)): adopting static IP %s", loadBalancerName, serviceName, ipAddress)
		} else {
			// For total clarity.  The IP did not pre-exist and the user did
			// not ask for a particular one, so we can release the IP in case
			// of failure or success.
			isSafeToReleaseIP = true
			glog.V(4).Infof("EnsureLoadBalancer(%v(%v)): allocated static IP %s", loadBalancerName, serviceName, ipAddress)
		}
	}

	// Deal with the firewall next. The reason we do this here rather than last
	// is because the forwarding rule is used as the indicator that the load
	// balancer is fully created - it's what getLoadBalancer checks for.
	// Check if user specified the allow source range
	sourceRanges, err := apiservice.GetLoadBalancerSourceRanges(apiService)
	if err != nil {
		return nil, err
	}

	firewallExists, firewallNeedsUpdate, err := gce.firewallNeedsUpdate(loadBalancerName, serviceName.String(), gce.region, ipAddress, ports, sourceRanges)
	if err != nil {
		return nil, err
	}

	if firewallNeedsUpdate {
		desc := makeFirewallDescription(serviceName.String(), ipAddress)
		// Unlike forwarding rules and target pools, firewalls can be updated
		// without needing to be deleted and recreated.
		if firewallExists {
			glog.Infof("EnsureLoadBalancer(%v(%v)): updating firewall", loadBalancerName, serviceName)
			if err = gce.updateFirewall(loadBalancerName, gce.region, desc, sourceRanges, ports, hosts); err != nil {
				return nil, err
			}
			glog.Infof("EnsureLoadBalancer(%v(%v)): updated firewall", loadBalancerName, serviceName)
		} else {
			glog.Infof("EnsureLoadBalancer(%v(%v)): creating firewall", loadBalancerName, serviceName)
			if err = gce.createFirewall(loadBalancerName, gce.region, desc, sourceRanges, ports, hosts); err != nil {
				return nil, err
			}
			glog.Infof("EnsureLoadBalancer(%v(%v)): created firewall", loadBalancerName, serviceName)
		}
	}

	tpExists, tpNeedsUpdate, err := gce.targetPoolNeedsUpdate(loadBalancerName, gce.region, affinityType)
	if err != nil {
		return nil, err
	}
	if !tpExists {
		glog.Infof("Target pool %v for Service %v/%v doesn't exist", loadBalancerName, apiService.Namespace, apiService.Name)
	}

	// Ensure health checks are created for this target pool to pass to createTargetPool for health check links
	// Alternately, if the annotation on the service was removed, we need to recreate the target pool without
	// health checks. This needs to be prior to the forwarding rule deletion below otherwise it is not possible
	// to delete just the target pool or http health checks later.
	var hcToCreate *compute.HttpHealthCheck
	hcExisting, err := gce.GetHTTPHealthCheck(loadBalancerName)
	if err != nil && !isHTTPErrorCode(err, http.StatusNotFound) {
		return nil, fmt.Errorf("Error checking HTTP health check %s: %v", loadBalancerName, err)
	}
	if path, healthCheckNodePort := apiservice.GetServiceHealthCheckPathPort(apiService); path != "" {
		glog.V(4).Infof("service %v (%v) needs health checks on :%d%s)", apiService.Name, loadBalancerName, healthCheckNodePort, path)
		if err != nil {
			// This logic exists to detect a transition for a pre-existing service and turn on
			// the tpNeedsUpdate flag to delete/recreate fwdrule/tpool adding the health check
			// to the target pool.
			glog.V(2).Infof("Annotation external-traffic=OnlyLocal added to new or pre-existing service")
			tpNeedsUpdate = true
		}
		hcToCreate, err = gce.ensureHTTPHealthCheck(loadBalancerName, path, healthCheckNodePort)
		if err != nil {
			return nil, fmt.Errorf("Failed to ensure health check for localized service %v on node port %v: %v", loadBalancerName, healthCheckNodePort, err)
		}
	} else {
		glog.V(4).Infof("service %v does not need health checks", apiService.Name)
		if err == nil {
			glog.V(2).Infof("Deleting stale health checks for service %v LB %v", apiService.Name, loadBalancerName)
			tpNeedsUpdate = true
		}
	}
	// Now we get to some slightly more interesting logic.
	// First, neither target pools nor forwarding rules can be updated in place -
	// they have to be deleted and recreated.
	// Second, forwarding rules are layered on top of target pools in that you
	// can't delete a target pool that's currently in use by a forwarding rule.
	// Thus, we have to tear down the forwarding rule if either it or the target
	// pool needs to be updated.
	if fwdRuleExists && (fwdRuleNeedsUpdate || tpNeedsUpdate) {
		// Begin critical section. If we have to delete the forwarding rule,
		// and something should fail before we recreate it, don't release the
		// IP.  That way we can come back to it later.
		isSafeToReleaseIP = false
		if err := gce.deleteForwardingRule(loadBalancerName, gce.region); err != nil {
			return nil, fmt.Errorf("failed to delete existing forwarding rule %s for load balancer update: %v", loadBalancerName, err)
		}
		glog.Infof("EnsureLoadBalancer(%v(%v)): deleted forwarding rule", loadBalancerName, serviceName)
	}
	if tpExists && tpNeedsUpdate {
		// Generate the list of health checks for this target pool to pass to deleteTargetPool
		if path, _ := apiservice.GetServiceHealthCheckPathPort(apiService); path != "" {
			var err error
			hcExisting, err = gce.GetHTTPHealthCheck(loadBalancerName)
			if err != nil && !isHTTPErrorCode(err, http.StatusNotFound) {
				glog.Infof("Failed to retrieve health check %v:%v", loadBalancerName, err)
			}
		}

		// Pass healthchecks to deleteTargetPool to cleanup health checks prior to cleaning up the target pool itself.
		if err := gce.deleteTargetPool(loadBalancerName, gce.region, hcExisting); err != nil {
			return nil, fmt.Errorf("failed to delete existing target pool %s for load balancer update: %v", loadBalancerName, err)
		}
		glog.Infof("EnsureLoadBalancer(%v(%v)): deleted target pool", loadBalancerName, serviceName)
	}

	// Once we've deleted the resources (if necessary), build them back up (or for
	// the first time if they're new).
	if tpNeedsUpdate {
		createInstances := hosts
		if len(hosts) > maxTargetPoolCreateInstances {
			createInstances = createInstances[:maxTargetPoolCreateInstances]
		}
		// Pass healthchecks to createTargetPool which needs them as health check links in the target pool
		if err := gce.createTargetPool(loadBalancerName, serviceName.String(), gce.region, createInstances, affinityType, hcToCreate); err != nil {
			return nil, fmt.Errorf("failed to create target pool %s: %v", loadBalancerName, err)
		}
		if hcToCreate != nil {
			glog.Infof("EnsureLoadBalancer(%v(%v)): created health checks for target pool", loadBalancerName, serviceName)
		}
		if len(hosts) <= maxTargetPoolCreateInstances {
			glog.Infof("EnsureLoadBalancer(%v(%v)): created target pool", loadBalancerName, serviceName)
		} else {
			glog.Infof("EnsureLoadBalancer(%v(%v)): created initial target pool (now updating with %d hosts)", loadBalancerName, serviceName, len(hosts)-maxTargetPoolCreateInstances)

			created := sets.NewString()
			for _, host := range createInstances {
				created.Insert(host.makeComparableHostPath())
			}
			if err := gce.updateTargetPool(loadBalancerName, created, hosts); err != nil {
				return nil, fmt.Errorf("failed to update target pool %s: %v", loadBalancerName, err)
			}
			glog.Infof("EnsureLoadBalancer(%v(%v)): updated target pool (with %d hosts)", loadBalancerName, serviceName, len(hosts)-maxTargetPoolCreateInstances)
		}
	}
	if tpNeedsUpdate || fwdRuleNeedsUpdate {
		glog.Infof("EnsureLoadBalancer(%v(%v)): creating forwarding rule, IP %s", loadBalancerName, serviceName, ipAddress)
		if err := gce.createForwardingRule(loadBalancerName, serviceName.String(), gce.region, ipAddress, ports); err != nil {
			return nil, fmt.Errorf("failed to create forwarding rule %s: %v", loadBalancerName, err)
		}
		// End critical section.  It is safe to release the static IP (which
		// just demotes it to ephemeral) now that it is attached.  In the case
		// of a user-requested IP, the "is user-owned" flag will be set,
		// preventing it from actually being released.
		isSafeToReleaseIP = true
		glog.Infof("EnsureLoadBalancer(%v(%v)): created forwarding rule, IP %s", loadBalancerName, serviceName, ipAddress)
	}

	status := &v1.LoadBalancerStatus{}
	status.Ingress = []v1.LoadBalancerIngress{{IP: ipAddress}}

	return status, nil
}
