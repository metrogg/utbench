func (lbaas *LbaasV2) EnsureLoadBalancer(ctx context.Context, clusterName string, apiService *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	klog.V(4).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)", clusterName, apiService.Namespace, apiService.Name, apiService.Spec.LoadBalancerIP, apiService.Spec.Ports, nodes, apiService.Annotations)

	if len(nodes) == 0 {
		return nil, fmt.Errorf("there are no available nodes for LoadBalancer service %s/%s", apiService.Namespace, apiService.Name)
	}

	lbaas.opts.SubnetID = getStringFromServiceAnnotation(apiService, ServiceAnnotationLoadBalancerSubnetID, lbaas.opts.SubnetID)
	if len(lbaas.opts.SubnetID) == 0 {
		// Get SubnetID automatically.
		// The LB needs to be configured with instance addresses on the same subnet, so get SubnetID by one node.
		subnetID, err := getSubnetIDForLB(lbaas.compute, *nodes[0])
		if err != nil {
			klog.Warningf("Failed to find subnet-id for loadbalancer service %s/%s: %v", apiService.Namespace, apiService.Name, err)
			return nil, fmt.Errorf("no subnet-id for service %s/%s : subnet-id not set in cloud provider config, "+
				"and failed to find subnet-id from OpenStack: %v", apiService.Namespace, apiService.Name, err)
		}
		lbaas.opts.SubnetID = subnetID
	}

	ports := apiService.Spec.Ports
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports provided to openstack load balancer")
	}

	floatingPool := getStringFromServiceAnnotation(apiService, ServiceAnnotationLoadBalancerFloatingNetworkID, lbaas.opts.FloatingNetworkID)
	if len(floatingPool) == 0 {
		var err error
		floatingPool, err = getFloatingNetworkIDForLB(lbaas.network)
		if err != nil {
			klog.Warningf("Failed to find floating-network-id for loadbalancer service %s/%s: %v", apiService.Namespace, apiService.Name, err)
		}
	}

	var internalAnnotation bool
	internal := getStringFromServiceAnnotation(apiService, ServiceAnnotationLoadBalancerInternal, "false")
	switch internal {
	case "true":
		klog.V(4).Infof("Ensure an internal loadbalancer service.")
		internalAnnotation = true
	case "false":
		if len(floatingPool) != 0 {
			klog.V(4).Infof("Ensure an external loadbalancer service, using floatingPool: %v", floatingPool)
			internalAnnotation = false
		} else {
			return nil, fmt.Errorf("floating-network-id or loadbalancer.openstack.org/floating-network-id should be specified when ensuring an external loadbalancer service")
		}
	default:
		return nil, fmt.Errorf("unknown service.beta.kubernetes.io/openstack-internal-load-balancer annotation: %v, specify \"true\" or \"false\" ",
			internal)
	}

	// Check for TCP protocol on each port
	// TODO: Convert all error messages to use an event recorder
	for _, port := range ports {
		if port.Protocol != v1.ProtocolTCP {
			return nil, fmt.Errorf("only TCP LoadBalancer is supported for openstack load balancers")
		}
	}

	sourceRanges, err := servicehelpers.GetLoadBalancerSourceRanges(apiService)
	if err != nil {
		return nil, fmt.Errorf("failed to get source ranges for loadbalancer service %s/%s: %v", apiService.Namespace, apiService.Name, err)
	}

	if !servicehelpers.IsAllowAll(sourceRanges) && !lbaas.opts.ManageSecurityGroups {
		return nil, fmt.Errorf("source range restrictions are not supported for openstack load balancers without managing security groups")
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

	name := lbaas.GetLoadBalancerName(ctx, clusterName, apiService)
	loadbalancer, err := getLoadbalancerByName(lbaas.lb, name)
	if err != nil {
		if err != ErrNotFound {
			return nil, fmt.Errorf("error getting loadbalancer %s: %v", name, err)
		}
		klog.V(2).Infof("Creating loadbalancer %s", name)
		loadbalancer, err = lbaas.createLoadBalancer(apiService, name, internalAnnotation)
		if err != nil {
			// Unknown error, retry later
			return nil, fmt.Errorf("error creating loadbalancer %s: %v", name, err)
		}
	} else {
		klog.V(2).Infof("LoadBalancer %s already exists", name)
	}

	provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
	}

	lbmethod := v2pools.LBMethod(lbaas.opts.LBMethod)
	if lbmethod == "" {
		lbmethod = v2pools.LBMethodRoundRobin
	}

	oldListeners, err := getListenersByLoadBalancerID(lbaas.lb, loadbalancer.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting LB %s listeners: %v", name, err)
	}
	for portIndex, port := range ports {
		listener := getListenerForPort(oldListeners, port)
		if listener == nil {
			klog.V(4).Infof("Creating listener for port %d", int(port.Port))
			listener, err = listeners.Create(lbaas.lb, listeners.CreateOpts{
				Name:           fmt.Sprintf("listener_%s_%d", name, portIndex),
				Protocol:       listeners.Protocol(port.Protocol),
				ProtocolPort:   int(port.Port),
				LoadbalancerID: loadbalancer.ID,
			}).Extract()
			if err != nil {
				// Unknown error, retry later
				return nil, fmt.Errorf("error creating LB listener: %v", err)
			}
			provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
			}
		}

		klog.V(4).Infof("Listener for %s port %d: %s", string(port.Protocol), int(port.Port), listener.ID)

		// After all ports have been processed, remaining listeners are removed as obsolete.
		// Pop valid listeners.
		oldListeners = popListener(oldListeners, listener.ID)
		pool, err := getPoolByListenerID(lbaas.lb, loadbalancer.ID, listener.ID)
		if err != nil && err != ErrNotFound {
			// Unknown error, retry later
			return nil, fmt.Errorf("error getting pool for listener %s: %v", listener.ID, err)
		}
		if pool == nil {
			klog.V(4).Infof("Creating pool for listener %s", listener.ID)
			pool, err = v2pools.Create(lbaas.lb, v2pools.CreateOpts{
				Name:        fmt.Sprintf("pool_%s_%d", name, portIndex),
				Protocol:    v2pools.Protocol(port.Protocol),
				LBMethod:    lbmethod,
				ListenerID:  listener.ID,
				Persistence: persistence,
			}).Extract()
			if err != nil {
				// Unknown error, retry later
				return nil, fmt.Errorf("error creating pool for listener %s: %v", listener.ID, err)
			}
			provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
			}

		}

		klog.V(4).Infof("Pool for listener %s: %s", listener.ID, pool.ID)
		members, err := getMembersByPoolID(lbaas.lb, pool.ID)
		if err != nil && !isNotFound(err) {
			return nil, fmt.Errorf("error getting pool members %s: %v", pool.ID, err)
		}
		for _, node := range nodes {
			addr, err := nodeAddressForLB(node)
			if err != nil {
				if err == ErrNotFound {
					// Node failure, do not create member
					klog.Warningf("Failed to create LB pool member for node %s: %v", node.Name, err)
					continue
				} else {
					return nil, fmt.Errorf("error getting address for node %s: %v", node.Name, err)
				}
			}

			if !memberExists(members, addr, int(port.NodePort)) {
				klog.V(4).Infof("Creating member for pool %s", pool.ID)
				_, err := v2pools.CreateMember(lbaas.lb, pool.ID, v2pools.CreateMemberOpts{
					Name:         fmt.Sprintf("member_%s_%d_%s", name, portIndex, node.Name),
					ProtocolPort: int(port.NodePort),
					Address:      addr,
					SubnetID:     lbaas.opts.SubnetID,
				}).Extract()
				if err != nil {
					return nil, fmt.Errorf("error creating LB pool member for node: %s, %v", node.Name, err)
				}

				provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
				}
			} else {
				// After all members have been processed, remaining members are deleted as obsolete.
				members = popMember(members, addr, int(port.NodePort))
			}

			klog.V(4).Infof("Ensured pool %s has member for %s at %s", pool.ID, node.Name, addr)
		}

		// Delete obsolete members for this pool
		for _, member := range members {
			klog.V(4).Infof("Deleting obsolete member %s for pool %s address %s", member.ID, pool.ID, member.Address)
			err := v2pools.DeleteMember(lbaas.lb, pool.ID, member.ID).ExtractErr()
			if err != nil && !isNotFound(err) {
				return nil, fmt.Errorf("error deleting obsolete member %s for pool %s address %s: %v", member.ID, pool.ID, member.Address, err)
			}
			provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
			}
		}

		monitorID := pool.MonitorID
		if monitorID == "" && lbaas.opts.CreateMonitor {
			klog.V(4).Infof("Creating monitor for pool %s", pool.ID)
			monitor, err := v2monitors.Create(lbaas.lb, v2monitors.CreateOpts{
				Name:       fmt.Sprintf("monitor_%s_%d", name, portIndex),
				PoolID:     pool.ID,
				Type:       string(port.Protocol),
				Delay:      int(lbaas.opts.MonitorDelay.Duration.Seconds()),
				Timeout:    int(lbaas.opts.MonitorTimeout.Duration.Seconds()),
				MaxRetries: int(lbaas.opts.MonitorMaxRetries),
			}).Extract()
			if err != nil {
				return nil, fmt.Errorf("error creating LB pool healthmonitor: %v", err)
			}
			provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
			}
			monitorID = monitor.ID
		} else if lbaas.opts.CreateMonitor == false {
			klog.V(4).Infof("Do not create monitor for pool %s when create-monitor is false", pool.ID)
		}

		if monitorID != "" {
			klog.V(4).Infof("Monitor for pool %s: %s", pool.ID, monitorID)
		}
	}

	// All remaining listeners are obsolete, delete
	for _, listener := range oldListeners {
		klog.V(4).Infof("Deleting obsolete listener %s:", listener.ID)
		// get pool for listener
		pool, err := getPoolByListenerID(lbaas.lb, loadbalancer.ID, listener.ID)
		if err != nil && err != ErrNotFound {
			return nil, fmt.Errorf("error getting pool for obsolete listener %s: %v", listener.ID, err)
		}
		if pool != nil {
			// get and delete monitor
			monitorID := pool.MonitorID
			if monitorID != "" {
				klog.V(4).Infof("Deleting obsolete monitor %s for pool %s", monitorID, pool.ID)
				err = v2monitors.Delete(lbaas.lb, monitorID).ExtractErr()
				if err != nil && !isNotFound(err) {
					return nil, fmt.Errorf("error deleting obsolete monitor %s for pool %s: %v", monitorID, pool.ID, err)
				}
				provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
				}
			}
			// get and delete pool members
			members, err := getMembersByPoolID(lbaas.lb, pool.ID)
			if err != nil && !isNotFound(err) {
				return nil, fmt.Errorf("error getting members for pool %s: %v", pool.ID, err)
			}
			if members != nil {
				for _, member := range members {
					klog.V(4).Infof("Deleting obsolete member %s for pool %s address %s", member.ID, pool.ID, member.Address)
					err := v2pools.DeleteMember(lbaas.lb, pool.ID, member.ID).ExtractErr()
					if err != nil && !isNotFound(err) {
						return nil, fmt.Errorf("error deleting obsolete member %s for pool %s address %s: %v", member.ID, pool.ID, member.Address, err)
					}
					provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
					if err != nil {
						return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
					}
				}
			}
			klog.V(4).Infof("Deleting obsolete pool %s for listener %s", pool.ID, listener.ID)
			// delete pool
			err = v2pools.Delete(lbaas.lb, pool.ID).ExtractErr()
			if err != nil && !isNotFound(err) {
				return nil, fmt.Errorf("error deleting obsolete pool %s for listener %s: %v", pool.ID, listener.ID, err)
			}
			provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
			}
		}
		// delete listener
		err = listeners.Delete(lbaas.lb, listener.ID).ExtractErr()
		if err != nil && !isNotFound(err) {
			return nil, fmt.Errorf("error deleteting obsolete listener: %v", err)
		}
		provisioningStatus, err := waitLoadbalancerActiveProvisioningStatus(lbaas.lb, loadbalancer.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to loadbalance ACTIVE provisioning status %v: %v", provisioningStatus, err)
		}
		klog.V(2).Infof("Deleted obsolete listener: %s", listener.ID)
	}

	portID := loadbalancer.VipPortID
	floatIP, err := getFloatingIPByPortID(lbaas.network, portID)
	if err != nil && err != ErrNotFound {
		return nil, fmt.Errorf("error getting floating ip for port %s: %v", portID, err)
	}
	if floatIP == nil && floatingPool != "" && !internalAnnotation {
		klog.V(4).Infof("Creating floating ip for loadbalancer %s port %s", loadbalancer.ID, portID)
		floatIPOpts := floatingips.CreateOpts{
			FloatingNetworkID: floatingPool,
			PortID:            portID,
		}

		loadBalancerIP := apiService.Spec.LoadBalancerIP
		if loadBalancerIP != "" {
			floatIPOpts.FloatingIP = loadBalancerIP
		}

		floatIP, err = floatingips.Create(lbaas.network, floatIPOpts).Extract()
		if err != nil {
			return nil, fmt.Errorf("error creating LB floatingip %+v: %v", floatIPOpts, err)
		}
	}

	status := &v1.LoadBalancerStatus{}

	if floatIP != nil {
		status.Ingress = []v1.LoadBalancerIngress{{IP: floatIP.FloatingIP}}
	} else {
		status.Ingress = []v1.LoadBalancerIngress{{IP: loadbalancer.VipAddress}}
	}

	if lbaas.opts.ManageSecurityGroups {
		err := lbaas.ensureSecurityGroup(clusterName, apiService, nodes, loadbalancer)
		if err != nil {
			// cleanup what was created so far
			_ = lbaas.EnsureLoadBalancerDeleted(ctx, clusterName, apiService)
			return status, err
		}
	}

	return status, nil
}
