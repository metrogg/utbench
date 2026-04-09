func (proxier *Proxier) syncProxyRules() {
	proxier.mu.Lock()
	defer proxier.mu.Unlock()

	start := time.Now()
	defer func() {
		metrics.SyncProxyRulesLatency.Observe(metrics.SinceInSeconds(start))
		metrics.DeprecatedSyncProxyRulesLatency.Observe(metrics.SinceInMicroseconds(start))
		klog.V(4).Infof("syncProxyRules took %v", time.Since(start))
	}()
	// don't sync rules till we've received services and endpoints
	if !proxier.endpointsSynced || !proxier.servicesSynced {
		klog.V(2).Info("Not syncing iptables until Services and Endpoints have been received from master")
		return
	}

	// We assume that if this was called, we really want to sync them,
	// even if nothing changed in the meantime. In other words, callers are
	// responsible for detecting no-op changes and not calling this function.
	serviceUpdateResult := proxy.UpdateServiceMap(proxier.serviceMap, proxier.serviceChanges)
	endpointUpdateResult := proxy.UpdateEndpointsMap(proxier.endpointsMap, proxier.endpointsChanges)

	staleServices := serviceUpdateResult.UDPStaleClusterIP
	// merge stale services gathered from updateEndpointsMap
	for _, svcPortName := range endpointUpdateResult.StaleServiceNames {
		if svcInfo, ok := proxier.serviceMap[svcPortName]; ok && svcInfo != nil && svcInfo.GetProtocol() == v1.ProtocolUDP {
			klog.V(2).Infof("Stale udp service %v -> %s", svcPortName, svcInfo.ClusterIPString())
			staleServices.Insert(svcInfo.ClusterIPString())
			for _, extIP := range svcInfo.ExternalIPStrings() {
				staleServices.Insert(extIP)
			}
		}
	}

	klog.V(3).Info("Syncing iptables rules")

	// Create and link the kube chains.
	for _, jump := range iptablesJumpChains {
		if _, err := proxier.iptables.EnsureChain(jump.table, jump.dstChain); err != nil {
			klog.Errorf("Failed to ensure that %s chain %s exists: %v", jump.table, jump.dstChain, err)
			return
		}
		args := append(jump.extraArgs,
			"-m", "comment", "--comment", jump.comment,
			"-j", string(jump.dstChain),
		)
		if _, err := proxier.iptables.EnsureRule(utiliptables.Prepend, jump.table, jump.srcChain, args...); err != nil {
			klog.Errorf("Failed to ensure that %s chain %s jumps to %s: %v", jump.table, jump.srcChain, jump.dstChain, err)
			return
		}
	}

	//
	// Below this point we will not return until we try to write the iptables rules.
	//

	// Get iptables-save output so we can check for existing chains and rules.
	// This will be a map of chain name to chain with rules as stored in iptables-save/iptables-restore
	existingFilterChains := make(map[utiliptables.Chain][]byte)
	proxier.existingFilterChainsData.Reset()
	err := proxier.iptables.SaveInto(utiliptables.TableFilter, proxier.existingFilterChainsData)
	if err != nil { // if we failed to get any rules
		klog.Errorf("Failed to execute iptables-save, syncing all rules: %v", err)
	} else { // otherwise parse the output
		existingFilterChains = utiliptables.GetChainLines(utiliptables.TableFilter, proxier.existingFilterChainsData.Bytes())
	}

	// IMPORTANT: existingNATChains may share memory with proxier.iptablesData.
	existingNATChains := make(map[utiliptables.Chain][]byte)
	proxier.iptablesData.Reset()
	err = proxier.iptables.SaveInto(utiliptables.TableNAT, proxier.iptablesData)
	if err != nil { // if we failed to get any rules
		klog.Errorf("Failed to execute iptables-save, syncing all rules: %v", err)
	} else { // otherwise parse the output
		existingNATChains = utiliptables.GetChainLines(utiliptables.TableNAT, proxier.iptablesData.Bytes())
	}

	// Reset all buffers used later.
	// This is to avoid memory reallocations and thus improve performance.
	proxier.filterChains.Reset()
	proxier.filterRules.Reset()
	proxier.natChains.Reset()
	proxier.natRules.Reset()

	// Write table headers.
	writeLine(proxier.filterChains, "*filter")
	writeLine(proxier.natChains, "*nat")

	// Make sure we keep stats for the top-level chains, if they existed
	// (which most should have because we created them above).
	for _, chainName := range []utiliptables.Chain{kubeServicesChain, kubeExternalServicesChain, kubeForwardChain} {
		if chain, ok := existingFilterChains[chainName]; ok {
			writeBytesLine(proxier.filterChains, chain)
		} else {
			writeLine(proxier.filterChains, utiliptables.MakeChainLine(chainName))
		}
	}
	for _, chainName := range []utiliptables.Chain{kubeServicesChain, kubeNodePortsChain, kubePostroutingChain, KubeMarkMasqChain} {
		if chain, ok := existingNATChains[chainName]; ok {
			writeBytesLine(proxier.natChains, chain)
		} else {
			writeLine(proxier.natChains, utiliptables.MakeChainLine(chainName))
		}
	}

	// Install the kubernetes-specific postrouting rules. We use a whole chain for
	// this so that it is easier to flush and change, for example if the mark
	// value should ever change.
	writeLine(proxier.natRules, []string{
		"-A", string(kubePostroutingChain),
		"-m", "comment", "--comment", `"kubernetes service traffic requiring SNAT"`,
		"-m", "mark", "--mark", proxier.masqueradeMark,
		"-j", "MASQUERADE",
	}...)

	// Install the kubernetes-specific masquerade mark rule. We use a whole chain for
	// this so that it is easier to flush and change, for example if the mark
	// value should ever change.
	writeLine(proxier.natRules, []string{
		"-A", string(KubeMarkMasqChain),
		"-j", "MARK", "--set-xmark", proxier.masqueradeMark,
	}...)

	// Accumulate NAT chains to keep.
	activeNATChains := map[utiliptables.Chain]bool{} // use a map as a set

	// Accumulate the set of local ports that we will be holding open once this update is complete
	replacementPortsMap := map[utilproxy.LocalPort]utilproxy.Closeable{}

	// We are creating those slices ones here to avoid memory reallocations
	// in every loop. Note that reuse the memory, instead of doing:
	//   slice = <some new slice>
	// you should always do one of the below:
	//   slice = slice[:0] // and then append to it
	//   slice = append(slice[:0], ...)
	endpoints := make([]*endpointsInfo, 0)
	endpointChains := make([]utiliptables.Chain, 0)
	// To avoid growing this slice, we arbitrarily set its size to 64,
	// there is never more than that many arguments for a single line.
	// Note that even if we go over 64, it will still be correct - it
	// is just for efficiency, not correctness.
	args := make([]string, 64)

	// Compute total number of endpoint chains across all services.
	proxier.endpointChainsNumber = 0
	for svcName := range proxier.serviceMap {
		proxier.endpointChainsNumber += len(proxier.endpointsMap[svcName])
	}

	// Build rules for each service.
	for svcName, svc := range proxier.serviceMap {
		svcInfo, ok := svc.(*serviceInfo)
		if !ok {
			klog.Errorf("Failed to cast serviceInfo %q", svcName.String())
			continue
		}
		isIPv6 := utilnet.IsIPv6(svcInfo.ClusterIP)
		protocol := strings.ToLower(string(svcInfo.Protocol))
		svcNameString := svcInfo.serviceNameString
		hasEndpoints := len(proxier.endpointsMap[svcName]) > 0

		svcChain := svcInfo.servicePortChainName
		if hasEndpoints {
			// Create the per-service chain, retaining counters if possible.
			if chain, ok := existingNATChains[svcChain]; ok {
				writeBytesLine(proxier.natChains, chain)
			} else {
				writeLine(proxier.natChains, utiliptables.MakeChainLine(svcChain))
			}
			activeNATChains[svcChain] = true
		}

		svcXlbChain := svcInfo.serviceLBChainName
		if svcInfo.OnlyNodeLocalEndpoints {
			// Only for services request OnlyLocal traffic
			// create the per-service LB chain, retaining counters if possible.
			if lbChain, ok := existingNATChains[svcXlbChain]; ok {
				writeBytesLine(proxier.natChains, lbChain)
			} else {
				writeLine(proxier.natChains, utiliptables.MakeChainLine(svcXlbChain))
			}
			activeNATChains[svcXlbChain] = true
		}

		// Capture the clusterIP.
		if hasEndpoints {
			args = append(args[:0],
				"-A", string(kubeServicesChain),
				"-m", "comment", "--comment", fmt.Sprintf(`"%s cluster IP"`, svcNameString),
				"-m", protocol, "-p", protocol,
				"-d", utilproxy.ToCIDR(svcInfo.ClusterIP),
				"--dport", strconv.Itoa(svcInfo.Port),
			)
			if proxier.masqueradeAll {
				writeLine(proxier.natRules, append(args, "-j", string(KubeMarkMasqChain))...)
			} else if len(proxier.clusterCIDR) > 0 {
				// This masquerades off-cluster traffic to a service VIP.  The idea
				// is that you can establish a static route for your Service range,
				// routing to any node, and that node will bridge into the Service
				// for you.  Since that might bounce off-node, we masquerade here.
				// If/when we support "Local" policy for VIPs, we should update this.
				writeLine(proxier.natRules, append(args, "! -s", proxier.clusterCIDR, "-j", string(KubeMarkMasqChain))...)
			}
			writeLine(proxier.natRules, append(args, "-j", string(svcChain))...)
		} else {
			// No endpoints.
			writeLine(proxier.filterRules,
				"-A", string(kubeServicesChain),
				"-m", "comment", "--comment", fmt.Sprintf(`"%s has no endpoints"`, svcNameString),
				"-m", protocol, "-p", protocol,
				"-d", utilproxy.ToCIDR(svcInfo.ClusterIP),
				"--dport", strconv.Itoa(svcInfo.Port),
				"-j", "REJECT",
			)
		}

		// Capture externalIPs.
		for _, externalIP := range svcInfo.ExternalIPs {
			// If the "external" IP happens to be an IP that is local to this
			// machine, hold the local port open so no other process can open it
			// (because the socket might open but it would never work).
			if local, err := utilproxy.IsLocalIP(externalIP); err != nil {
				klog.Errorf("can't determine if IP is local, assuming not: %v", err)
			} else if local && (svcInfo.GetProtocol() != v1.ProtocolSCTP) {
				lp := utilproxy.LocalPort{
					Description: "externalIP for " + svcNameString,
					IP:          externalIP,
					Port:        svcInfo.Port,
					Protocol:    protocol,
				}
				if proxier.portsMap[lp] != nil {
					klog.V(4).Infof("Port %s was open before and is still needed", lp.String())
					replacementPortsMap[lp] = proxier.portsMap[lp]
				} else {
					socket, err := proxier.portMapper.OpenLocalPort(&lp)
					if err != nil {
						msg := fmt.Sprintf("can't open %s, skipping this externalIP: %v", lp.String(), err)

						proxier.recorder.Eventf(
							&v1.ObjectReference{
								Kind:      "Node",
								Name:      proxier.hostname,
								UID:       types.UID(proxier.hostname),
								Namespace: "",
							}, v1.EventTypeWarning, err.Error(), msg)
						klog.Error(msg)
						continue
					}
					replacementPortsMap[lp] = socket
				}
			}

			if hasEndpoints {
				args = append(args[:0],
					"-A", string(kubeServicesChain),
					"-m", "comment", "--comment", fmt.Sprintf(`"%s external IP"`, svcNameString),
					"-m", protocol, "-p", protocol,
					"-d", utilproxy.ToCIDR(net.ParseIP(externalIP)),
					"--dport", strconv.Itoa(svcInfo.Port),
				)
				// We have to SNAT packets to external IPs.
				writeLine(proxier.natRules, append(args, "-j", string(KubeMarkMasqChain))...)

				// Allow traffic for external IPs that does not come from a bridge (i.e. not from a container)
				// nor from a local process to be forwarded to the service.
				// This rule roughly translates to "all traffic from off-machine".
				// This is imperfect in the face of network plugins that might not use a bridge, but we can revisit that later.
				externalTrafficOnlyArgs := append(args,
					"-m", "physdev", "!", "--physdev-is-in",
					"-m", "addrtype", "!", "--src-type", "LOCAL")
				writeLine(proxier.natRules, append(externalTrafficOnlyArgs, "-j", string(svcChain))...)
				dstLocalOnlyArgs := append(args, "-m", "addrtype", "--dst-type", "LOCAL")
				// Allow traffic bound for external IPs that happen to be recognized as local IPs to stay local.
				// This covers cases like GCE load-balancers which get added to the local routing table.
				writeLine(proxier.natRules, append(dstLocalOnlyArgs, "-j", string(svcChain))...)
			} else {
				// No endpoints.
				writeLine(proxier.filterRules,
					"-A", string(kubeExternalServicesChain),
					"-m", "comment", "--comment", fmt.Sprintf(`"%s has no endpoints"`, svcNameString),
					"-m", protocol, "-p", protocol,
					"-d", utilproxy.ToCIDR(net.ParseIP(externalIP)),
					"--dport", strconv.Itoa(svcInfo.Port),
					"-j", "REJECT",
				)
			}
		}

		// Capture load-balancer ingress.
		fwChain := svcInfo.serviceFirewallChainName
		for _, ingress := range svcInfo.LoadBalancerStatus.Ingress {
			if ingress.IP != "" {
				if hasEndpoints {
					// create service firewall chain
					if chain, ok := existingNATChains[fwChain]; ok {
						writeBytesLine(proxier.natChains, chain)
					} else {
						writeLine(proxier.natChains, utiliptables.MakeChainLine(fwChain))
					}
					activeNATChains[fwChain] = true
					// The service firewall rules are created based on ServiceSpec.loadBalancerSourceRanges field.
					// This currently works for loadbalancers that preserves source ips.
					// For loadbalancers which direct traffic to service NodePort, the firewall rules will not apply.

					args = append(args[:0],
						"-A", string(kubeServicesChain),
						"-m", "comment", "--comment", fmt.Sprintf(`"%s loadbalancer IP"`, svcNameString),
						"-m", protocol, "-p", protocol,
						"-d", utilproxy.ToCIDR(net.ParseIP(ingress.IP)),
						"--dport", strconv.Itoa(svcInfo.Port),
					)
					// jump to service firewall chain
					writeLine(proxier.natRules, append(args, "-j", string(fwChain))...)

					args = append(args[:0],
						"-A", string(fwChain),
						"-m", "comment", "--comment", fmt.Sprintf(`"%s loadbalancer IP"`, svcNameString),
					)

					// Each source match rule in the FW chain may jump to either the SVC or the XLB chain
					chosenChain := svcXlbChain
					// If we are proxying globally, we need to masquerade in case we cross nodes.
					// If we are proxying only locally, we can retain the source IP.
					if !svcInfo.OnlyNodeLocalEndpoints {
						writeLine(proxier.natRules, append(args, "-j", string(KubeMarkMasqChain))...)
						chosenChain = svcChain
					}

					if len(svcInfo.LoadBalancerSourceRanges) == 0 {
						// allow all sources, so jump directly to the KUBE-SVC or KUBE-XLB chain
						writeLine(proxier.natRules, append(args, "-j", string(chosenChain))...)
					} else {
						// firewall filter based on each source range
						allowFromNode := false
						for _, src := range svcInfo.LoadBalancerSourceRanges {
							writeLine(proxier.natRules, append(args, "-s", src, "-j", string(chosenChain))...)
							// ignore error because it has been validated
							_, cidr, _ := net.ParseCIDR(src)
							if cidr.Contains(proxier.nodeIP) {
								allowFromNode = true
							}
						}
						// generally, ip route rule was added to intercept request to loadbalancer vip from the
						// loadbalancer's backend hosts. In this case, request will not hit the loadbalancer but loop back directly.
						// Need to add the following rule to allow request on host.
						if allowFromNode {
							writeLine(proxier.natRules, append(args, "-s", utilproxy.ToCIDR(net.ParseIP(ingress.IP)), "-j", string(chosenChain))...)
						}
					}

					// If the packet was able to reach the end of firewall chain, then it did not get DNATed.
					// It means the packet cannot go thru the firewall, then mark it for DROP
					writeLine(proxier.natRules, append(args, "-j", string(KubeMarkDropChain))...)
				} else {
					// No endpoints.
					writeLine(proxier.filterRules,
						"-A", string(kubeServicesChain),
						"-m", "comment", "--comment", fmt.Sprintf(`"%s has no endpoints"`, svcNameString),
						"-m", protocol, "-p", protocol,
						"-d", utilproxy.ToCIDR(net.ParseIP(ingress.IP)),
						"--dport", strconv.Itoa(svcInfo.Port),
						"-j", "REJECT",
					)
				}
			}
		}

		// Capture nodeports.  If we had more than 2 rules it might be
		// worthwhile to make a new per-service chain for nodeport rules, but
		// with just 2 rules it ends up being a waste and a cognitive burden.
		if svcInfo.NodePort != 0 {
			// Hold the local port open so no other process can open it
			// (because the socket might open but it would never work).
			addresses, err := utilproxy.GetNodeAddresses(proxier.nodePortAddresses, proxier.networkInterfacer)
			if err != nil {
				klog.Errorf("Failed to get node ip address matching nodeport cidr: %v", err)
				continue
			}

			lps := make([]utilproxy.LocalPort, 0)
			for address := range addresses {
				lp := utilproxy.LocalPort{
					Description: "nodePort for " + svcNameString,
					IP:          address,
					Port:        svcInfo.NodePort,
					Protocol:    protocol,
				}
				if utilproxy.IsZeroCIDR(address) {
					// Empty IP address means all
					lp.IP = ""
					lps = append(lps, lp)
					// If we encounter a zero CIDR, then there is no point in processing the rest of the addresses.
					break
				}
				lps = append(lps, lp)
			}

			// For ports on node IPs, open the actual port and hold it.
			for _, lp := range lps {
				if proxier.portsMap[lp] != nil {
					klog.V(4).Infof("Port %s was open before and is still needed", lp.String())
					replacementPortsMap[lp] = proxier.portsMap[lp]
				} else if svcInfo.GetProtocol() != v1.ProtocolSCTP {
					socket, err := proxier.portMapper.OpenLocalPort(&lp)
					if err != nil {
						klog.Errorf("can't open %s, skipping this nodePort: %v", lp.String(), err)
						continue
					}
					if lp.Protocol == "udp" {
						// TODO: We might have multiple services using the same port, and this will clear conntrack for all of them.
						// This is very low impact. The NodePort range is intentionally obscure, and unlikely to actually collide with real Services.
						// This only affects UDP connections, which are not common.
						// See issue: https://github.com/kubernetes/kubernetes/issues/49881
						err := conntrack.ClearEntriesForPort(proxier.exec, lp.Port, isIPv6, v1.ProtocolUDP)
						if err != nil {
							klog.Errorf("Failed to clear udp conntrack for port %d, error: %v", lp.Port, err)
						}
					}
					replacementPortsMap[lp] = socket
				}
			}

			if hasEndpoints {
				args = append(args[:0],
					"-A", string(kubeNodePortsChain),
					"-m", "comment", "--comment", svcNameString,
					"-m", protocol, "-p", protocol,
					"--dport", strconv.Itoa(svcInfo.NodePort),
				)
				if !svcInfo.OnlyNodeLocalEndpoints {
					// Nodeports need SNAT, unless they're local.
					writeLine(proxier.natRules, append(args, "-j", string(KubeMarkMasqChain))...)
					// Jump to the service chain.
					writeLine(proxier.natRules, append(args, "-j", string(svcChain))...)
				} else {
					// TODO: Make all nodePorts jump to the firewall chain.
					// Currently we only create it for loadbalancers (#33586).

					// Fix localhost martian source error
					loopback := "127.0.0.0/8"
					if isIPv6 {
						loopback = "::1/128"
					}
					writeLine(proxier.natRules, append(args, "-s", loopback, "-j", string(KubeMarkMasqChain))...)
					writeLine(proxier.natRules, append(args, "-j", string(svcXlbChain))...)
				}
			} else {
				// No endpoints.
				writeLine(proxier.filterRules,
					"-A", string(kubeExternalServicesChain),
					"-m", "comment", "--comment", fmt.Sprintf(`"%s has no endpoints"`, svcNameString),
					"-m", "addrtype", "--dst-type", "LOCAL",
					"-m", protocol, "-p", protocol,
					"--dport", strconv.Itoa(svcInfo.NodePort),
					"-j", "REJECT",
				)
			}
		}

		if !hasEndpoints {
			continue
		}

		// Generate the per-endpoint chains.  We do this in multiple passes so we
		// can group rules together.
		// These two slices parallel each other - keep in sync
		endpoints = endpoints[:0]
		endpointChains = endpointChains[:0]
		var endpointChain utiliptables.Chain
		for _, ep := range proxier.endpointsMap[svcName] {
			epInfo, ok := ep.(*endpointsInfo)
			if !ok {
				klog.Errorf("Failed to cast endpointsInfo %q", ep.String())
				continue
			}
			endpoints = append(endpoints, epInfo)
			endpointChain = epInfo.endpointChain(svcNameString, protocol)
			endpointChains = append(endpointChains, endpointChain)

			// Create the endpoint chain, retaining counters if possible.
			if chain, ok := existingNATChains[utiliptables.Chain(endpointChain)]; ok {
				writeBytesLine(proxier.natChains, chain)
			} else {
				writeLine(proxier.natChains, utiliptables.MakeChainLine(endpointChain))
			}
			activeNATChains[endpointChain] = true
		}

		// First write session affinity rules, if applicable.
		if svcInfo.SessionAffinityType == v1.ServiceAffinityClientIP {
			for _, endpointChain := range endpointChains {
				args = append(args[:0],
					"-A", string(svcChain),
				)
				proxier.appendServiceCommentLocked(args, svcNameString)
				args = append(args,
					"-m", "recent", "--name", string(endpointChain),
					"--rcheck", "--seconds", strconv.Itoa(svcInfo.StickyMaxAgeSeconds), "--reap",
					"-j", string(endpointChain),
				)
				writeLine(proxier.natRules, args...)
			}
		}

		// Now write loadbalancing & DNAT rules.
		n := len(endpointChains)
		localEndpoints := make([]*endpointsInfo, 0)
		localEndpointChains := make([]utiliptables.Chain, 0)
		for i, endpointChain := range endpointChains {
			// Write ingress loadbalancing & DNAT rules only for services that request OnlyLocal traffic.
			if svcInfo.OnlyNodeLocalEndpoints && endpoints[i].IsLocal {
				// These slices parallel each other; must be kept in sync
				localEndpoints = append(localEndpoints, endpoints[i])
				localEndpointChains = append(localEndpointChains, endpointChains[i])
			}

			epIP := endpoints[i].IP()
			if epIP == "" {
				// Error parsing this endpoint has been logged. Skip to next endpoint.
				continue
			}
			// Balancing rules in the per-service chain.
			args = append(args[:0], "-A", string(svcChain))
			proxier.appendServiceCommentLocked(args, svcNameString)
			if i < (n - 1) {
				// Each rule is a probabilistic match.
				args = append(args,
					"-m", "statistic",
					"--mode", "random",
					"--probability", proxier.probability(n-i))
			}
			// The final (or only if n == 1) rule is a guaranteed match.
			args = append(args, "-j", string(endpointChain))
			writeLine(proxier.natRules, args...)

			// Rules in the per-endpoint chain.
			args = append(args[:0], "-A", string(endpointChain))
			proxier.appendServiceCommentLocked(args, svcNameString)
			// Handle traffic that loops back to the originator with SNAT.
			writeLine(proxier.natRules, append(args,
				"-s", utilproxy.ToCIDR(net.ParseIP(epIP)),
				"-j", string(KubeMarkMasqChain))...)
			// Update client-affinity lists.
			if svcInfo.SessionAffinityType == v1.ServiceAffinityClientIP {
				args = append(args, "-m", "recent", "--name", string(endpointChain), "--set")
			}
			// DNAT to final destination.
			args = append(args, "-m", protocol, "-p", protocol, "-j", "DNAT", "--to-destination", endpoints[i].Endpoint)
			writeLine(proxier.natRules, args...)
		}

		// The logic below this applies only if this service is marked as OnlyLocal
		if !svcInfo.OnlyNodeLocalEndpoints {
			continue
		}

		// First rule in the chain redirects all pod -> external VIP traffic to the
		// Service's ClusterIP instead. This happens whether or not we have local
		// endpoints; only if clusterCIDR is specified
		if len(proxier.clusterCIDR) > 0 {
			args = append(args[:0],
				"-A", string(svcXlbChain),
				"-m", "comment", "--comment",
				`"Redirect pods trying to reach external loadbalancer VIP to clusterIP"`,
				"-s", proxier.clusterCIDR,
				"-j", string(svcChain),
			)
			writeLine(proxier.natRules, args...)
		}

		numLocalEndpoints := len(localEndpointChains)
		if numLocalEndpoints == 0 {
			// Blackhole all traffic since there are no local endpoints
			args = append(args[:0],
				"-A", string(svcXlbChain),
				"-m", "comment", "--comment",
				fmt.Sprintf(`"%s has no local endpoints"`, svcNameString),
				"-j",
				string(KubeMarkDropChain),
			)
			writeLine(proxier.natRules, args...)
		} else {
			// First write session affinity rules only over local endpoints, if applicable.
			if svcInfo.SessionAffinityType == v1.ServiceAffinityClientIP {
				for _, endpointChain := range localEndpointChains {
					writeLine(proxier.natRules,
						"-A", string(svcXlbChain),
						"-m", "comment", "--comment", svcNameString,
						"-m", "recent", "--name", string(endpointChain),
						"--rcheck", "--seconds", strconv.Itoa(svcInfo.StickyMaxAgeSeconds), "--reap",
						"-j", string(endpointChain))
				}
			}

			// Setup probability filter rules only over local endpoints
			for i, endpointChain := range localEndpointChains {
				// Balancing rules in the per-service chain.
				args = append(args[:0],
					"-A", string(svcXlbChain),
					"-m", "comment", "--comment",
					fmt.Sprintf(`"Balancing rule %d for %s"`, i, svcNameString),
				)
				if i < (numLocalEndpoints - 1) {
					// Each rule is a probabilistic match.
					args = append(args,
						"-m", "statistic",
						"--mode", "random",
						"--probability", proxier.probability(numLocalEndpoints-i))
				}
				// The final (or only if n == 1) rule is a guaranteed match.
				args = append(args, "-j", string(endpointChain))
				writeLine(proxier.natRules, args...)
			}
		}
	}

	// Delete chains no longer in use.
	for chain := range existingNATChains {
		if !activeNATChains[chain] {
			chainString := string(chain)
			if !strings.HasPrefix(chainString, "KUBE-SVC-") && !strings.HasPrefix(chainString, "KUBE-SEP-") && !strings.HasPrefix(chainString, "KUBE-FW-") && !strings.HasPrefix(chainString, "KUBE-XLB-") {
				// Ignore chains that aren't ours.
				continue
			}
			// We must (as per iptables) write a chain-line for it, which has
			// the nice effect of flushing the chain.  Then we can remove the
			// chain.
			writeBytesLine(proxier.natChains, existingNATChains[chain])
			writeLine(proxier.natRules, "-X", chainString)
		}
	}

	// Finally, tail-call to the nodeports chain.  This needs to be after all
	// other service portal rules.
	addresses, err := utilproxy.GetNodeAddresses(proxier.nodePortAddresses, proxier.networkInterfacer)
	if err != nil {
		klog.Errorf("Failed to get node ip address matching nodeport cidr")
	} else {
		isIPv6 := proxier.iptables.IsIpv6()
		for address := range addresses {
			// TODO(thockin, m1093782566): If/when we have dual-stack support we will want to distinguish v4 from v6 zero-CIDRs.
			if utilproxy.IsZeroCIDR(address) {
				args = append(args[:0],
					"-A", string(kubeServicesChain),
					"-m", "comment", "--comment", `"kubernetes service nodeports; NOTE: this must be the last rule in this chain"`,
					"-m", "addrtype", "--dst-type", "LOCAL",
					"-j", string(kubeNodePortsChain))
				writeLine(proxier.natRules, args...)
				// Nothing else matters after the zero CIDR.
				break
			}
			// Ignore IP addresses with incorrect version
			if isIPv6 && !utilnet.IsIPv6String(address) || !isIPv6 && utilnet.IsIPv6String(address) {
				klog.Errorf("IP address %s has incorrect IP version", address)
				continue
			}
			// create nodeport rules for each IP one by one
			args = append(args[:0],
				"-A", string(kubeServicesChain),
				"-m", "comment", "--comment", `"kubernetes service nodeports; NOTE: this must be the last rule in this chain"`,
				"-d", address,
				"-j", string(kubeNodePortsChain))
			writeLine(proxier.natRules, args...)
		}
	}

	// Drop the packets in INVALID state, which would potentially cause
	// unexpected connection reset.
	// https://github.com/kubernetes/kubernetes/issues/74839
	writeLine(proxier.filterRules,
		"-A", string(kubeForwardChain),
		"-m", "conntrack",
		"--ctstate", "INVALID",
		"-j", "DROP",
	)

	// If the masqueradeMark has been added then we want to forward that same
	// traffic, this allows NodePort traffic to be forwarded even if the default
	// FORWARD policy is not accept.
	writeLine(proxier.filterRules,
		"-A", string(kubeForwardChain),
		"-m", "comment", "--comment", `"kubernetes forwarding rules"`,
		"-m", "mark", "--mark", proxier.masqueradeMark,
		"-j", "ACCEPT",
	)

	// The following rules can only be set if clusterCIDR has been defined.
	if len(proxier.clusterCIDR) != 0 {
		// The following two rules ensure the traffic after the initial packet
		// accepted by the "kubernetes forwarding rules" rule above will be
		// accepted, to be as specific as possible the traffic must be sourced
		// or destined to the clusterCIDR (to/from a pod).
		writeLine(proxier.filterRules,
			"-A", string(kubeForwardChain),
			"-s", proxier.clusterCIDR,
			"-m", "comment", "--comment", `"kubernetes forwarding conntrack pod source rule"`,
			"-m", "conntrack",
			"--ctstate", "RELATED,ESTABLISHED",
			"-j", "ACCEPT",
		)
		writeLine(proxier.filterRules,
			"-A", string(kubeForwardChain),
			"-m", "comment", "--comment", `"kubernetes forwarding conntrack pod destination rule"`,
			"-d", proxier.clusterCIDR,
			"-m", "conntrack",
			"--ctstate", "RELATED,ESTABLISHED",
			"-j", "ACCEPT",
		)
	}

	// Write the end-of-table markers.
	writeLine(proxier.filterRules, "COMMIT")
	writeLine(proxier.natRules, "COMMIT")

	// Sync rules.
	// NOTE: NoFlushTables is used so we don't flush non-kubernetes chains in the table
	proxier.iptablesData.Reset()
	proxier.iptablesData.Write(proxier.filterChains.Bytes())
	proxier.iptablesData.Write(proxier.filterRules.Bytes())
	proxier.iptablesData.Write(proxier.natChains.Bytes())
	proxier.iptablesData.Write(proxier.natRules.Bytes())

	klog.V(5).Infof("Restoring iptables rules: %s", proxier.iptablesData.String())
	err = proxier.iptables.RestoreAll(proxier.iptablesData.Bytes(), utiliptables.NoFlushTables, utiliptables.RestoreCounters)
	if err != nil {
		klog.Errorf("Failed to execute iptables-restore: %v\nfailed payload:\n%s", err, proxier.iptablesData.String())
		// Revert new local ports.
		klog.V(2).Infof("Closing local ports after iptables-restore failure")
		utilproxy.RevertPorts(replacementPortsMap, proxier.portsMap)
		return
	}
	for _, lastChangeTriggerTime := range endpointUpdateResult.LastChangeTriggerTimes {
		latency := metrics.SinceInSeconds(lastChangeTriggerTime)
		metrics.NetworkProgrammingLatency.Observe(latency)
		klog.V(4).Infof("Network programming took %f seconds", latency)
	}

	// Close old local ports and save new ones.
	for k, v := range proxier.portsMap {
		if replacementPortsMap[k] == nil {
			v.Close()
		}
	}
	proxier.portsMap = replacementPortsMap

	// Update healthz timestamp.
	if proxier.healthzServer != nil {
		proxier.healthzServer.UpdateTimestamp()
	}

	// Update healthchecks.  The endpoints list might include services that are
	// not "OnlyLocal", but the services list will not, and the healthChecker
	// will just drop those endpoints.
	if err := proxier.healthChecker.SyncServices(serviceUpdateResult.HCServiceNodePorts); err != nil {
		klog.Errorf("Error syncing healthcheck services: %v", err)
	}
	if err := proxier.healthChecker.SyncEndpoints(endpointUpdateResult.HCEndpointsLocalIPSize); err != nil {
		klog.Errorf("Error syncing healthcheck endpoints: %v", err)
	}

	// Finish housekeeping.
	// TODO: these could be made more consistent.
	for _, svcIP := range staleServices.UnsortedList() {
		if err := conntrack.ClearEntriesForIP(proxier.exec, svcIP, v1.ProtocolUDP); err != nil {
			klog.Errorf("Failed to delete stale service IP %s connections, error: %v", svcIP, err)
		}
	}
	proxier.deleteEndpointConnections(endpointUpdateResult.StaleEndpoints)
}
