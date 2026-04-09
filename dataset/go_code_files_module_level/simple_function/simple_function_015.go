func (c *Cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, apiService *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	annotations := apiService.Annotations
	klog.V(2).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)",
		clusterName, apiService.Namespace, apiService.Name, c.region, apiService.Spec.LoadBalancerIP, apiService.Spec.Ports, annotations)

	if apiService.Spec.SessionAffinity != v1.ServiceAffinityNone {
		// ELB supports sticky sessions, but only when configured for HTTP/HTTPS
		return nil, fmt.Errorf("unsupported load balancer affinity: %v", apiService.Spec.SessionAffinity)
	}

	if len(apiService.Spec.Ports) == 0 {
		return nil, fmt.Errorf("requested load balancer with no ports")
	}

	// Figure out what mappings we want on the load balancer
	listeners := []*elb.Listener{}
	v2Mappings := []nlbPortMapping{}

	sslPorts := getPortSets(annotations[ServiceAnnotationLoadBalancerSSLPorts])
	for _, port := range apiService.Spec.Ports {
		if port.Protocol != v1.ProtocolTCP {
			return nil, fmt.Errorf("Only TCP LoadBalancer is supported for AWS ELB")
		}
		if port.NodePort == 0 {
			klog.Errorf("Ignoring port without NodePort defined: %v", port)
			continue
		}

		if isNLB(annotations) {
			portMapping := nlbPortMapping{
				FrontendPort:     int64(port.Port),
				FrontendProtocol: string(port.Protocol),
				TrafficPort:      int64(port.NodePort),
				TrafficProtocol:  string(port.Protocol),

				// if externalTrafficPolicy == "Local", we'll override the
				// health check later
				HealthCheckPort:     int64(port.NodePort),
				HealthCheckProtocol: elbv2.ProtocolEnumTcp,
			}

			certificateARN := annotations[ServiceAnnotationLoadBalancerCertificate]
			if certificateARN != "" && (sslPorts == nil || sslPorts.numbers.Has(int64(port.Port)) || sslPorts.names.Has(port.Name)) {
				portMapping.FrontendProtocol = elbv2.ProtocolEnumTls
				portMapping.SSLCertificateARN = certificateARN
				portMapping.SSLPolicy = annotations[ServiceAnnotationLoadBalancerSSLNegotiationPolicy]

				if backendProtocol := annotations[ServiceAnnotationLoadBalancerBEProtocol]; backendProtocol == "ssl" {
					portMapping.TrafficProtocol = elbv2.ProtocolEnumTls
				}
			}

			v2Mappings = append(v2Mappings, portMapping)
		}
		listener, err := buildListener(port, annotations, sslPorts)
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, listener)
	}

	if apiService.Spec.LoadBalancerIP != "" {
		return nil, fmt.Errorf("LoadBalancerIP cannot be specified for AWS ELB")
	}

	instances, err := c.findInstancesForELB(nodes)
	if err != nil {
		return nil, err
	}

	sourceRanges, err := servicehelpers.GetLoadBalancerSourceRanges(apiService)
	if err != nil {
		return nil, err
	}

	// Determine if this is tagged as an Internal ELB
	internalELB := false
	internalAnnotation := apiService.Annotations[ServiceAnnotationLoadBalancerInternal]
	if internalAnnotation == "false" {
		internalELB = false
	} else if internalAnnotation != "" {
		internalELB = true
	}

	if isNLB(annotations) {

		if path, healthCheckNodePort := servicehelpers.GetServiceHealthCheckPathPort(apiService); path != "" {
			for i := range v2Mappings {
				v2Mappings[i].HealthCheckPort = int64(healthCheckNodePort)
				v2Mappings[i].HealthCheckPath = path
				v2Mappings[i].HealthCheckProtocol = elbv2.ProtocolEnumHttp
			}
		}

		// Find the subnets that the ELB will live in
		subnetIDs, err := c.findELBSubnets(internalELB)
		if err != nil {
			klog.Errorf("Error listing subnets in VPC: %q", err)
			return nil, err
		}
		// Bail out early if there are no subnets
		if len(subnetIDs) == 0 {
			return nil, fmt.Errorf("could not find any suitable subnets for creating the ELB")
		}

		loadBalancerName := c.GetLoadBalancerName(ctx, clusterName, apiService)
		serviceName := types.NamespacedName{Namespace: apiService.Namespace, Name: apiService.Name}

		instanceIDs := []string{}
		for id := range instances {
			instanceIDs = append(instanceIDs, string(id))
		}

		v2LoadBalancer, err := c.ensureLoadBalancerv2(
			serviceName,
			loadBalancerName,
			v2Mappings,
			instanceIDs,
			subnetIDs,
			internalELB,
			annotations,
		)
		if err != nil {
			return nil, err
		}

		sourceRangeCidrs := []string{}
		for cidr := range sourceRanges {
			sourceRangeCidrs = append(sourceRangeCidrs, cidr)
		}
		if len(sourceRangeCidrs) == 0 {
			sourceRangeCidrs = append(sourceRangeCidrs, "0.0.0.0/0")
		}

		err = c.updateInstanceSecurityGroupsForNLB(v2Mappings, instances, loadBalancerName, sourceRangeCidrs)
		if err != nil {
			klog.Warningf("Error opening ingress rules for the load balancer to the instances: %q", err)
			return nil, err
		}

		// We don't have an `ensureLoadBalancerInstances()` function for elbv2
		// because `ensureLoadBalancerv2()` requires instance Ids

		// TODO: Wait for creation?
		return v2toStatus(v2LoadBalancer), nil
	}

	// Determine if we need to set the Proxy protocol policy
	proxyProtocol := false
	proxyProtocolAnnotation := apiService.Annotations[ServiceAnnotationLoadBalancerProxyProtocol]
	if proxyProtocolAnnotation != "" {
		if proxyProtocolAnnotation != "*" {
			return nil, fmt.Errorf("annotation %q=%q detected, but the only value supported currently is '*'", ServiceAnnotationLoadBalancerProxyProtocol, proxyProtocolAnnotation)
		}
		proxyProtocol = true
	}

	// Some load balancer attributes are required, so defaults are set. These can be overridden by annotations.
	loadBalancerAttributes := &elb.LoadBalancerAttributes{
		AccessLog:              &elb.AccessLog{Enabled: aws.Bool(false)},
		ConnectionDraining:     &elb.ConnectionDraining{Enabled: aws.Bool(false)},
		ConnectionSettings:     &elb.ConnectionSettings{IdleTimeout: aws.Int64(60)},
		CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{Enabled: aws.Bool(false)},
	}

	// Determine if an access log emit interval has been specified
	accessLogEmitIntervalAnnotation := annotations[ServiceAnnotationLoadBalancerAccessLogEmitInterval]
	if accessLogEmitIntervalAnnotation != "" {
		accessLogEmitInterval, err := strconv.ParseInt(accessLogEmitIntervalAnnotation, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerAccessLogEmitInterval,
				accessLogEmitIntervalAnnotation,
			)
		}
		loadBalancerAttributes.AccessLog.EmitInterval = &accessLogEmitInterval
	}

	// Determine if access log enabled/disabled has been specified
	accessLogEnabledAnnotation := annotations[ServiceAnnotationLoadBalancerAccessLogEnabled]
	if accessLogEnabledAnnotation != "" {
		accessLogEnabled, err := strconv.ParseBool(accessLogEnabledAnnotation)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerAccessLogEnabled,
				accessLogEnabledAnnotation,
			)
		}
		loadBalancerAttributes.AccessLog.Enabled = &accessLogEnabled
	}

	// Determine if access log s3 bucket name has been specified
	accessLogS3BucketNameAnnotation := annotations[ServiceAnnotationLoadBalancerAccessLogS3BucketName]
	if accessLogS3BucketNameAnnotation != "" {
		loadBalancerAttributes.AccessLog.S3BucketName = &accessLogS3BucketNameAnnotation
	}

	// Determine if access log s3 bucket prefix has been specified
	accessLogS3BucketPrefixAnnotation := annotations[ServiceAnnotationLoadBalancerAccessLogS3BucketPrefix]
	if accessLogS3BucketPrefixAnnotation != "" {
		loadBalancerAttributes.AccessLog.S3BucketPrefix = &accessLogS3BucketPrefixAnnotation
	}

	// Determine if connection draining enabled/disabled has been specified
	connectionDrainingEnabledAnnotation := annotations[ServiceAnnotationLoadBalancerConnectionDrainingEnabled]
	if connectionDrainingEnabledAnnotation != "" {
		connectionDrainingEnabled, err := strconv.ParseBool(connectionDrainingEnabledAnnotation)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerConnectionDrainingEnabled,
				connectionDrainingEnabledAnnotation,
			)
		}
		loadBalancerAttributes.ConnectionDraining.Enabled = &connectionDrainingEnabled
	}

	// Determine if connection draining timeout has been specified
	connectionDrainingTimeoutAnnotation := annotations[ServiceAnnotationLoadBalancerConnectionDrainingTimeout]
	if connectionDrainingTimeoutAnnotation != "" {
		connectionDrainingTimeout, err := strconv.ParseInt(connectionDrainingTimeoutAnnotation, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerConnectionDrainingTimeout,
				connectionDrainingTimeoutAnnotation,
			)
		}
		loadBalancerAttributes.ConnectionDraining.Timeout = &connectionDrainingTimeout
	}

	// Determine if connection idle timeout has been specified
	connectionIdleTimeoutAnnotation := annotations[ServiceAnnotationLoadBalancerConnectionIdleTimeout]
	if connectionIdleTimeoutAnnotation != "" {
		connectionIdleTimeout, err := strconv.ParseInt(connectionIdleTimeoutAnnotation, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerConnectionIdleTimeout,
				connectionIdleTimeoutAnnotation,
			)
		}
		loadBalancerAttributes.ConnectionSettings.IdleTimeout = &connectionIdleTimeout
	}

	// Determine if cross zone load balancing enabled/disabled has been specified
	crossZoneLoadBalancingEnabledAnnotation := annotations[ServiceAnnotationLoadBalancerCrossZoneLoadBalancingEnabled]
	if crossZoneLoadBalancingEnabledAnnotation != "" {
		crossZoneLoadBalancingEnabled, err := strconv.ParseBool(crossZoneLoadBalancingEnabledAnnotation)
		if err != nil {
			return nil, fmt.Errorf("error parsing service annotation: %s=%s",
				ServiceAnnotationLoadBalancerCrossZoneLoadBalancingEnabled,
				crossZoneLoadBalancingEnabledAnnotation,
			)
		}
		loadBalancerAttributes.CrossZoneLoadBalancing.Enabled = &crossZoneLoadBalancingEnabled
	}

	// Find the subnets that the ELB will live in
	subnetIDs, err := c.findELBSubnets(internalELB)
	if err != nil {
		klog.Errorf("Error listing subnets in VPC: %q", err)
		return nil, err
	}

	// Bail out early if there are no subnets
	if len(subnetIDs) == 0 {
		return nil, fmt.Errorf("could not find any suitable subnets for creating the ELB")
	}

	loadBalancerName := c.GetLoadBalancerName(ctx, clusterName, apiService)
	serviceName := types.NamespacedName{Namespace: apiService.Namespace, Name: apiService.Name}
	securityGroupIDs, err := c.buildELBSecurityGroupList(serviceName, loadBalancerName, annotations)
	if err != nil {
		return nil, err
	}
	if len(securityGroupIDs) == 0 {
		return nil, fmt.Errorf("[BUG] ELB can't have empty list of Security Groups to be assigned, this is a Kubernetes bug, please report")
	}

	{
		ec2SourceRanges := []*ec2.IpRange{}
		for _, sourceRange := range sourceRanges.StringSlice() {
			ec2SourceRanges = append(ec2SourceRanges, &ec2.IpRange{CidrIp: aws.String(sourceRange)})
		}

		permissions := NewIPPermissionSet()
		for _, port := range apiService.Spec.Ports {
			portInt64 := int64(port.Port)
			protocol := strings.ToLower(string(port.Protocol))

			permission := &ec2.IpPermission{}
			permission.FromPort = &portInt64
			permission.ToPort = &portInt64
			permission.IpRanges = ec2SourceRanges
			permission.IpProtocol = &protocol

			permissions.Insert(permission)
		}

		// Allow ICMP fragmentation packets, important for MTU discovery
		{
			permission := &ec2.IpPermission{
				IpProtocol: aws.String("icmp"),
				FromPort:   aws.Int64(3),
				ToPort:     aws.Int64(4),
				IpRanges:   ec2SourceRanges,
			}

			permissions.Insert(permission)
		}
		_, err = c.setSecurityGroupIngress(securityGroupIDs[0], permissions)
		if err != nil {
			return nil, err
		}
	}

	// Build the load balancer itself
	loadBalancer, err := c.ensureLoadBalancer(
		serviceName,
		loadBalancerName,
		listeners,
		subnetIDs,
		securityGroupIDs,
		internalELB,
		proxyProtocol,
		loadBalancerAttributes,
		annotations,
	)
	if err != nil {
		return nil, err
	}

	if sslPolicyName, ok := annotations[ServiceAnnotationLoadBalancerSSLNegotiationPolicy]; ok {
		err := c.ensureSSLNegotiationPolicy(loadBalancer, sslPolicyName)
		if err != nil {
			return nil, err
		}

		for _, port := range c.getLoadBalancerTLSPorts(loadBalancer) {
			err := c.setSSLNegotiationPolicy(loadBalancerName, sslPolicyName, port)
			if err != nil {
				return nil, err
			}
		}
	}

	if path, healthCheckNodePort := servicehelpers.GetServiceHealthCheckPathPort(apiService); path != "" {
		klog.V(4).Infof("service %v (%v) needs health checks on :%d%s)", apiService.Name, loadBalancerName, healthCheckNodePort, path)
		err = c.ensureLoadBalancerHealthCheck(loadBalancer, "HTTP", healthCheckNodePort, path, annotations)
		if err != nil {
			return nil, fmt.Errorf("Failed to ensure health check for localized service %v on node port %v: %q", loadBalancerName, healthCheckNodePort, err)
		}
	} else {
		klog.V(4).Infof("service %v does not need custom health checks", apiService.Name)
		// We only configure a TCP health-check on the first port
		var tcpHealthCheckPort int32
		for _, listener := range listeners {
			if listener.InstancePort == nil {
				continue
			}
			tcpHealthCheckPort = int32(*listener.InstancePort)
			break
		}
		annotationProtocol := strings.ToLower(annotations[ServiceAnnotationLoadBalancerBEProtocol])
		var hcProtocol string
		if annotationProtocol == "https" || annotationProtocol == "ssl" {
			hcProtocol = "SSL"
		} else {
			hcProtocol = "TCP"
		}
		// there must be no path on TCP health check
		err = c.ensureLoadBalancerHealthCheck(loadBalancer, hcProtocol, tcpHealthCheckPort, "", annotations)
		if err != nil {
			return nil, err
		}
	}

	err = c.updateInstanceSecurityGroupsForLoadBalancer(loadBalancer, instances)
	if err != nil {
		klog.Warningf("Error opening ingress rules for the load balancer to the instances: %q", err)
		return nil, err
	}

	err = c.ensureLoadBalancerInstances(aws.StringValue(loadBalancer.LoadBalancerName), loadBalancer.Instances, instances)
	if err != nil {
		klog.Warningf("Error registering instances with the load balancer: %q", err)
		return nil, err
	}

	klog.V(1).Infof("Loadbalancer %s (%v) has DNS name %s", loadBalancerName, serviceName, aws.StringValue(loadBalancer.DNSName))

	// TODO: Wait for creation?

	status := toStatus(loadBalancer)
	return status, nil
}
