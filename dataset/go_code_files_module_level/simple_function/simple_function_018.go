func (lbaas *LbaasV2) EnsureLoadBalancer(clusterName string, apiService *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	glog.V(4).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)", clusterName, apiService.Namespace, apiService.Name, apiService.Spec.LoadBalancerIP, apiService.Spec.Ports, nodes, apiService.Annotations)

	ports := apiService.Spec.Ports
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports provided to openstack load balancer")
	}

	// Check for TCP protocol on each port
	// TODO: Convert all error messages to use an event recorder
	for _, port := range ports {
		if port.Protocol != v1.ProtocolTCP {
			return nil, fmt.Errorf("Only TCP LoadBalancer is supported for openstack load balancers")
		}
	}

	sourceRanges, err := service.GetLoadBalancerSourceRanges(apiService)
	if err != nil {
		return nil, err
	}

	if !service.IsAllowAll(sourceRanges) && !lbaas.opts.ManageSecurityGroups {
		return nil, fmt.Errorf("Source range restrictions are not supported for openstack load balancers without managing security groups")
	}

	affinity := apiService.Spec.SessionAffinity
	var persistence *v2pools.SessionPersistence
	switch affinity {
	case v1.ServiceAffinityNone:
		persistence = nil
	case v1.ServiceAffinityClientIP:
		persistence = &v2pools.SessionPersistence{Type: "SOURCE_IP"}
	default:
		return nil, fmt.Errorf("unsupported load balancer affinity: %v", affinity)
	}

	name := cloudprovider.GetLoadBalancerName(apiService)
	loadbalancer, err := getLoadbalancerByName(lbaas.network, name)
	if err != nil {
		if err != ErrNotFound {
			return nil, fmt.Errorf("Error getting loadbalancer %s: %v", name, err)
		}
		glog.V(2).Infof("Creating loadbalancer %s", name)
		loadbalancer, err = lbaas.createLoadBalancer(apiService, name)
		if err != nil {
			// Unknown error, retry later
			return nil, fmt.Errorf("Error creating loadbalancer %s: %v", name, err)
		}
	} else {
		glog.V(2).Infof("LoadBalancer %s already exists", name)
	}

	waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)

	lbmethod := v2pools.LBMethod(lbaas.opts.LBMethod)
	if lbmethod == "" {
		lbmethod = v2pools.LBMethodRoundRobin
	}

	oldListeners, err := getListenersByLoadBalancerID(lbaas.network, loadbalancer.ID)
	if err != nil {
		return nil, fmt.Errorf("Error getting LB %s listeners: %v", name, err)
	}
	for portIndex, port := range ports {
		listener := getListenerForPort(oldListeners, port)
		if listener == nil {
			glog.V(4).Infof("Creating listener for port %d", int(port.Port))
			listener, err = listeners.Create(lbaas.network, listeners.CreateOpts{
				Name:           fmt.Sprintf("listener_%s_%d", name, portIndex),
				Protocol:       listeners.Protocol(port.Protocol),
				ProtocolPort:   int(port.Port),
				LoadbalancerID: loadbalancer.ID,
			}).Extract()
			if err != nil {
				// Unknown error, retry later
				return nil, fmt.Errorf("Error creating LB listener: %v", err)
			}
			waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
		}

		glog.V(4).Infof("Listener for %s port %d: %s", string(port.Protocol), int(port.Port), listener.ID)

		// After all ports have been processed, remaining listeners are removed as obsolete.
		// Pop valid listeners.
		oldListeners = popListener(oldListeners, listener.ID)
		pool, err := getPoolByListenerID(lbaas.network, loadbalancer.ID, listener.ID)
		if err != nil && err != ErrNotFound {
			// Unknown error, retry later
			return nil, fmt.Errorf("Error getting pool for listener %s: %v", listener.ID, err)
		}
		if pool == nil {
			glog.V(4).Infof("Creating pool for listener %s", listener.ID)
			pool, err = v2pools.Create(lbaas.network, v2pools.CreateOpts{
				Name:        fmt.Sprintf("pool_%s_%d", name, portIndex),
				Protocol:    v2pools.Protocol(port.Protocol),
				LBMethod:    lbmethod,
				ListenerID:  listener.ID,
				Persistence: persistence,
			}).Extract()
			if err != nil {
				// Unknown error, retry later
				return nil, fmt.Errorf("Error creating pool for listener %s: %v", listener.ID, err)
			}
			waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
		}

		glog.V(4).Infof("Pool for listener %s: %s", listener.ID, pool.ID)
		members, err := getMembersByPoolID(lbaas.network, pool.ID)
		if err != nil && !isNotFound(err) {
			return nil, fmt.Errorf("Error getting pool members %s: %v", pool.ID, err)
		}
		for _, node := range nodes {
			addr, err := nodeAddressForLB(node)
			if err != nil {
				if err == ErrNotFound {
					// Node failure, do not create member
					glog.Warningf("Failed to create LB pool member for node %s: %v", node.Name, err)
					continue
				} else {
					return nil, fmt.Errorf("Error getting address for node %s: %v", node.Name, err)
				}
			}

			if !memberExists(members, addr, int(port.NodePort)) {
				glog.V(4).Infof("Creating member for pool %s", pool.ID)
				_, err := v2pools.CreateMember(lbaas.network, pool.ID, v2pools.CreateMemberOpts{
					ProtocolPort: int(port.NodePort),
					Address:      addr,
					SubnetID:     lbaas.opts.SubnetID,
				}).Extract()
				if err != nil {
					return nil, fmt.Errorf("Error creating LB pool member for node: %s, %v", node.Name, err)
				}

				waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
			} else {
				// After all members have been processed, remaining members are deleted as obsolete.
				members = popMember(members, addr, int(port.NodePort))
			}

			glog.V(4).Infof("Ensured pool %s has member for %s at %s", pool.ID, node.Name, addr)
		}

		// Delete obsolete members for this pool
		for _, member := range members {
			glog.V(4).Infof("Deleting obsolete member %s for pool %s address %s", member.ID, pool.ID, member.Address)
			err := v2pools.DeleteMember(lbaas.network, pool.ID, member.ID).ExtractErr()
			if err != nil && !isNotFound(err) {
				return nil, fmt.Errorf("Error deleting obsolete member %s for pool %s address %s: %v", member.ID, pool.ID, member.Address, err)
			}
			waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
		}

		monitorID := pool.MonitorID
		if monitorID == "" && lbaas.opts.CreateMonitor {
			glog.V(4).Infof("Creating monitor for pool %s", pool.ID)
			monitor, err := v2monitors.Create(lbaas.network, v2monitors.CreateOpts{
				PoolID:     pool.ID,
				Type:       string(port.Protocol),
				Delay:      int(lbaas.opts.MonitorDelay.Duration.Seconds()),
				Timeout:    int(lbaas.opts.MonitorTimeout.Duration.Seconds()),
				MaxRetries: int(lbaas.opts.MonitorMaxRetries),
			}).Extract()
			if err != nil {
				return nil, fmt.Errorf("Error creating LB pool healthmonitor: %v", err)
			}
			waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
			monitorID = monitor.ID
		}

		glog.V(4).Infof("Monitor for pool %s: %s", pool.ID, monitorID)
	}

	// All remaining listeners are obsolete, delete
	for _, listener := range oldListeners {
		glog.V(4).Infof("Deleting obsolete listener %s:", listener.ID)
		// get pool for listener
		pool, err := getPoolByListenerID(lbaas.network, loadbalancer.ID, listener.ID)
		if err != nil && err != ErrNotFound {
			return nil, fmt.Errorf("Error getting pool for obsolete listener %s: %v", listener.ID, err)
		}
		if pool != nil {
			// get and delete monitor
			monitorID := pool.MonitorID
			if monitorID != "" {
				glog.V(4).Infof("Deleting obsolete monitor %s for pool %s", monitorID, pool.ID)
				err = v2monitors.Delete(lbaas.network, monitorID).ExtractErr()
				if err != nil && !isNotFound(err) {
					return nil, fmt.Errorf("Error deleting obsolete monitor %s for pool %s: %v", monitorID, pool.ID, err)
				}
				waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
			}
			// get and delete pool members
			members, poolErr := getMembersByPoolID(lbaas.network, pool.ID)
			if poolErr != nil && !isNotFound(poolErr) {
				return nil, fmt.Errorf("Error getting members for pool %s: %v", pool.ID, poolErr)
			}
			if members != nil {
				for _, member := range members {
					glog.V(4).Infof("Deleting obsolete member %s for pool %s address %s", member.ID, pool.ID, member.Address)
					err := v2pools.DeleteMember(lbaas.network, pool.ID, member.ID).ExtractErr()
					if err != nil && !isNotFound(err) {
						return nil, fmt.Errorf("Error deleting obsolete member %s for pool %s address %s: %v", member.ID, pool.ID, member.Address, err)
					}
					waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
				}
			}
			glog.V(4).Infof("Deleting obsolete pool %s for listener %s", pool.ID, listener.ID)
			// delete pool
			err = v2pools.Delete(lbaas.network, pool.ID).ExtractErr()
			if err != nil && !isNotFound(err) {
				return nil, fmt.Errorf("Error deleting obsolete pool %s for listener %s: %v", pool.ID, listener.ID, err)
			}
			waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
		}
		// delete listener
		err = listeners.Delete(lbaas.network, listener.ID).ExtractErr()
		if err != nil && !isNotFound(err) {
			return nil, fmt.Errorf("Error deleteting obsolete listener: %v", err)
		}
		waitLoadbalancerActiveProvisioningStatus(lbaas.network, loadbalancer.ID)
		glog.V(2).Infof("Deleted obsolete listener: %s", listener.ID)
	}

	status := &v1.LoadBalancerStatus{}

	status.Ingress = []v1.LoadBalancerIngress{{IP: loadbalancer.VipAddress}}

	portID := loadbalancer.VipPortID
	floatIP, err := getFloatingIPByPortID(lbaas.network, portID)
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("Error getting floating ip for port %s: %v", portID, err)
	}
	if floatIP == nil && lbaas.opts.FloatingNetworkID != "" {
		glog.V(4).Infof("Creating floating ip for loadbalancer %s port %s", loadbalancer.ID, portID)
		floatIPOpts := floatingips.CreateOpts{
			FloatingNetworkID: lbaas.opts.FloatingNetworkID,
			PortID:            portID,
		}
		floatIP, err = floatingips.Create(lbaas.network, floatIPOpts).Extract()
		if err != nil {
			return nil, fmt.Errorf("Error creating LB floatingip %+v: %v", floatIPOpts, err)
		}
	}
	if floatIP != nil {
		status.Ingress = append(status.Ingress, v1.LoadBalancerIngress{IP: floatIP.FloatingIP})
	}

	if lbaas.opts.ManageSecurityGroups {
		lbSecGroupCreateOpts := groups.CreateOpts{
			Name:        getSecurityGroupName(clusterName, apiService),
			Description: fmt.Sprintf("Securty Group for %v Service LoadBalancer", apiService.Name),
		}

		lbSecGroup, err := groups.Create(lbaas.network, lbSecGroupCreateOpts).Extract()

		if err != nil {
			// cleanup what was created so far
			_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
			return nil, err
		}

		for _, port := range ports {

			for _, sourceRange := range sourceRanges.StringSlice() {
				ethertype := rules.EtherType4
				network, _, err := net.ParseCIDR(sourceRange)

				if err != nil {
					// cleanup what was created so far
					glog.Errorf("Error parsing source range %s as a CIDR", sourceRange)
					_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
					return nil, err
				}

				if network.To4() == nil {
					ethertype = rules.EtherType6
				}

				lbSecGroupRuleCreateOpts := rules.CreateOpts{
					Direction:      rules.DirIngress,
					PortRangeMax:   int(port.Port),
					PortRangeMin:   int(port.Port),
					Protocol:       toRuleProtocol(port.Protocol),
					RemoteIPPrefix: sourceRange,
					SecGroupID:     lbSecGroup.ID,
					EtherType:      ethertype,
				}

				_, err = rules.Create(lbaas.network, lbSecGroupRuleCreateOpts).Extract()

				if err != nil {
					// cleanup what was created so far
					_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
					return nil, err
				}
			}

			err := createNodeSecurityGroup(lbaas.network, lbaas.opts.NodeSecurityGroupID, int(port.NodePort), port.Protocol, lbSecGroup.ID)
			if err != nil {
				glog.Errorf("Error occurred creating security group for loadbalancer %s:", loadbalancer.ID)
				_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
				return nil, err
			}
		}

		lbSecGroupRuleCreateOpts := rules.CreateOpts{
			Direction:      rules.DirIngress,
			PortRangeMax:   4, // ICMP: Code -  Values for ICMP  "Destination Unreachable: Fragmentation Needed and Don't Fragment was Set"
			PortRangeMin:   3, // ICMP: Type
			Protocol:       rules.ProtocolICMP,
			RemoteIPPrefix: "0.0.0.0/0", // The Fragmentation packet can come from anywhere along the path back to the sourceRange - we need to all this from all
			SecGroupID:     lbSecGroup.ID,
			EtherType:      rules.EtherType4,
		}

		_, err = rules.Create(lbaas.network, lbSecGroupRuleCreateOpts).Extract()

		if err != nil {
			// cleanup what was created so far
			_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
			return nil, err
		}

		lbSecGroupRuleCreateOpts = rules.CreateOpts{
			Direction:      rules.DirIngress,
			PortRangeMax:   0, // ICMP: Code - Values for ICMP "Packet Too Big"
			PortRangeMin:   2, // ICMP: Type
			Protocol:       rules.ProtocolICMP,
			RemoteIPPrefix: "::/0", // The Fragmentation packet can come from anywhere along the path back to the sourceRange - we need to all this from all
			SecGroupID:     lbSecGroup.ID,
			EtherType:      rules.EtherType6,
		}

		_, err = rules.Create(lbaas.network, lbSecGroupRuleCreateOpts).Extract()

		if err != nil {
			// cleanup what was created so far
			_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
			return nil, err
		}

		portID := loadbalancer.VipPortID
		newOpts := []string{lbSecGroup.ID}
		updateOpts := neutronports.UpdateOpts{SecurityGroups: &newOpts}
		res := neutronports.Update(lbaas.network, portID, updateOpts)
		if res.Err != nil {
			glog.Errorf("Error occurred updating port: %s", portID)
			// cleanup what was created so far
			_ = lbaas.EnsureLoadBalancerDeleted(clusterName, apiService)
			return nil, res.Err
		}
	}

	return status, nil
}
