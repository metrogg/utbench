func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	roleEntry, err := b.nonLockedAWSRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		roleEntry = &awsRoleEntry{
			RoleID:  roleID,
			Version: currentRoleStorageVersion,
		}
	} else {
		needUpdate, err := b.upgradeRoleEntry(ctx, req.Storage, roleEntry)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to update roleEntry: %v", err)), nil
		}
		if needUpdate {
			err = b.nonLockedSetAWSRole(ctx, req.Storage, roleName, roleEntry)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("failed to save upgraded roleEntry: %v", err)), nil
			}
		}
	}

	// Fetch and set the bound parameters. There can't be default values
	// for these.
	if boundAmiIDRaw, ok := data.GetOk("bound_ami_id"); ok {
		roleEntry.BoundAmiIDs = boundAmiIDRaw.([]string)
	}

	if boundAccountIDRaw, ok := data.GetOk("bound_account_id"); ok {
		roleEntry.BoundAccountIDs = boundAccountIDRaw.([]string)
	}

	if boundRegionRaw, ok := data.GetOk("bound_region"); ok {
		roleEntry.BoundRegions = boundRegionRaw.([]string)
	}

	if boundVpcIDRaw, ok := data.GetOk("bound_vpc_id"); ok {
		roleEntry.BoundVpcIDs = boundVpcIDRaw.([]string)
	}

	if boundSubnetIDRaw, ok := data.GetOk("bound_subnet_id"); ok {
		roleEntry.BoundSubnetIDs = boundSubnetIDRaw.([]string)
	}

	if resolveAWSUniqueIDsRaw, ok := data.GetOk("resolve_aws_unique_ids"); ok {
		switch {
		case req.Operation == logical.CreateOperation:
			roleEntry.ResolveAWSUniqueIDs = resolveAWSUniqueIDsRaw.(bool)
		case roleEntry.ResolveAWSUniqueIDs && !resolveAWSUniqueIDsRaw.(bool):
			return logical.ErrorResponse("changing resolve_aws_unique_ids from true to false is not allowed"), nil
		default:
			roleEntry.ResolveAWSUniqueIDs = resolveAWSUniqueIDsRaw.(bool)
		}
	} else if req.Operation == logical.CreateOperation {
		roleEntry.ResolveAWSUniqueIDs = data.Get("resolve_aws_unique_ids").(bool)
	}

	if boundIamRoleARNRaw, ok := data.GetOk("bound_iam_role_arn"); ok {
		roleEntry.BoundIamRoleARNs = boundIamRoleARNRaw.([]string)
	}

	if boundIamInstanceProfileARNRaw, ok := data.GetOk("bound_iam_instance_profile_arn"); ok {
		roleEntry.BoundIamInstanceProfileARNs = boundIamInstanceProfileARNRaw.([]string)
	}

	if boundEc2InstanceIDRaw, ok := data.GetOk("bound_ec2_instance_id"); ok {
		roleEntry.BoundEc2InstanceIDs = boundEc2InstanceIDRaw.([]string)
	}

	if boundIamPrincipalARNRaw, ok := data.GetOk("bound_iam_principal_arn"); ok {
		principalARNs := boundIamPrincipalARNRaw.([]string)
		roleEntry.BoundIamPrincipalARNs = principalARNs
		roleEntry.BoundIamPrincipalIDs = []string{}
	}
	if roleEntry.ResolveAWSUniqueIDs && len(roleEntry.BoundIamPrincipalIDs) == 0 {
		// we might be turning on resolution on this role, so ensure we update the IDs
		for _, principalARN := range roleEntry.BoundIamPrincipalARNs {
			if !strings.HasSuffix(principalARN, "*") {
				principalID, err := b.resolveArnToUniqueIDFunc(ctx, req.Storage, principalARN)
				if err != nil {
					return logical.ErrorResponse(fmt.Sprintf("unable to resolve ARN %#v to internal ID: %s", principalARN, err.Error())), nil
				}
				roleEntry.BoundIamPrincipalIDs = append(roleEntry.BoundIamPrincipalIDs, principalID)
			}
		}
	}

	if inferRoleTypeRaw, ok := data.GetOk("inferred_entity_type"); ok {
		roleEntry.InferredEntityType = inferRoleTypeRaw.(string)
	}

	if inferredAWSRegionRaw, ok := data.GetOk("inferred_aws_region"); ok {
		roleEntry.InferredAWSRegion = inferredAWSRegionRaw.(string)
	}

	// auth_type is a special case as it's immutable and can't be changed once a role is created
	if authTypeRaw, ok := data.GetOk("auth_type"); ok {
		// roleEntry.AuthType should only be "" when it's a new role; existing roles without an
		// auth_type should have already been upgraded to have one before we get here
		if roleEntry.AuthType == "" {
			switch authTypeRaw.(string) {
			case ec2AuthType, iamAuthType:
				roleEntry.AuthType = authTypeRaw.(string)
			default:
				return logical.ErrorResponse(fmt.Sprintf("unrecognized auth_type: %v", authTypeRaw.(string))), nil
			}
		} else if authTypeRaw.(string) != roleEntry.AuthType {
			return logical.ErrorResponse("changing auth_type on a role is not allowed"), nil
		}
	} else if req.Operation == logical.CreateOperation {
		switch req.MountType {
		// maintain backwards compatibility for old aws-ec2 auth types
		case "aws-ec2":
			roleEntry.AuthType = ec2AuthType
		// but default to iamAuth for new mounts going forward
		case "aws":
			roleEntry.AuthType = iamAuthType
		default:
			roleEntry.AuthType = iamAuthType
		}
	}

	allowEc2Binds := roleEntry.AuthType == ec2AuthType

	if roleEntry.InferredEntityType != "" {
		switch {
		case roleEntry.AuthType != iamAuthType:
			return logical.ErrorResponse("specified inferred_entity_type but didn't allow iam auth_type"), nil
		case roleEntry.InferredEntityType != ec2EntityType:
			return logical.ErrorResponse(fmt.Sprintf("specified invalid inferred_entity_type: %s", roleEntry.InferredEntityType)), nil
		case roleEntry.InferredAWSRegion == "":
			return logical.ErrorResponse("specified inferred_entity_type but not inferred_aws_region"), nil
		}
		allowEc2Binds = true
	} else if roleEntry.InferredAWSRegion != "" {
		return logical.ErrorResponse("specified inferred_aws_region but not inferred_entity_type"), nil
	}

	numBinds := 0

	if len(roleEntry.BoundAccountIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_account_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundRegions) > 0 {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified bound_region but not specifying ec2 auth_type"), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundAmiIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_ami_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamInstanceProfileARNs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_instance_profile_arn but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundEc2InstanceIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_ec2_instance_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamRoleARNs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_iam_role_arn but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundIamPrincipalARNs) > 0 {
		if roleEntry.AuthType != iamAuthType {
			return logical.ErrorResponse("specified bound_iam_principal_arn but not specifying iam auth_type"), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundVpcIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_vpc_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if len(roleEntry.BoundSubnetIDs) > 0 {
		if !allowEc2Binds {
			return logical.ErrorResponse(fmt.Sprintf("specified bound_subnet_id but not specifying ec2 auth_type or inferring %s", ec2EntityType)), nil
		}
		numBinds++
	}

	if numBinds == 0 {
		return logical.ErrorResponse("at least one bound parameter should be specified on the role"), nil
	}

	policiesRaw, ok := data.GetOk("policies")
	if ok {
		roleEntry.Policies = policyutil.ParsePolicies(policiesRaw)
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Policies = []string{}
	}

	disallowReauthenticationBool, ok := data.GetOk("disallow_reauthentication")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified disallow_reauthentication when not using ec2 auth type"), nil
		}
		roleEntry.DisallowReauthentication = disallowReauthenticationBool.(bool)
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
		roleEntry.DisallowReauthentication = data.Get("disallow_reauthentication").(bool)
	}

	allowInstanceMigrationBool, ok := data.GetOk("allow_instance_migration")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("specified allow_instance_migration when not using ec2 auth type"), nil
		}
		roleEntry.AllowInstanceMigration = allowInstanceMigrationBool.(bool)
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
		roleEntry.AllowInstanceMigration = data.Get("allow_instance_migration").(bool)
	}

	if roleEntry.AllowInstanceMigration && roleEntry.DisallowReauthentication {
		return logical.ErrorResponse("cannot specify both disallow_reauthentication=true and allow_instance_migration=true"), nil
	}

	var resp logical.Response

	ttlRaw, ok := data.GetOk("ttl")
	if ok {
		ttl := time.Duration(ttlRaw.(int)) * time.Second
		defaultLeaseTTL := b.System().DefaultLeaseTTL()
		if ttl > defaultLeaseTTL {
			resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time", ttl/time.Second, defaultLeaseTTL/time.Second))
		}
		roleEntry.TTL = ttl
	} else if req.Operation == logical.CreateOperation {
		roleEntry.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	maxTTLInt, ok := data.GetOk("max_ttl")
	if ok {
		maxTTL := time.Duration(maxTTLInt.(int)) * time.Second
		systemMaxTTL := b.System().MaxLeaseTTL()
		if maxTTL > systemMaxTTL {
			resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time", maxTTL/time.Second, systemMaxTTL/time.Second))
		}

		if maxTTL < time.Duration(0) {
			return logical.ErrorResponse("max_ttl cannot be negative"), nil
		}

		roleEntry.MaxTTL = maxTTL
	} else if req.Operation == logical.CreateOperation {
		roleEntry.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	if roleEntry.MaxTTL != 0 && roleEntry.MaxTTL < roleEntry.TTL {
		return logical.ErrorResponse("ttl should be shorter than max_ttl"), nil
	}

	periodRaw, ok := data.GetOk("period")
	if ok {
		roleEntry.Period = time.Second * time.Duration(periodRaw.(int))
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Period = time.Second * time.Duration(data.Get("period").(int))
	}

	if roleEntry.Period > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("'period' of '%s' is greater than the backend's maximum lease TTL of '%s'", roleEntry.Period.String(), b.System().MaxLeaseTTL().String())), nil
	}

	roleTagStr, ok := data.GetOk("role_tag")
	if ok {
		if roleEntry.AuthType != ec2AuthType {
			return logical.ErrorResponse("tried to enable role_tag when not using ec2 auth method"), nil
		}
		roleEntry.RoleTag = roleTagStr.(string)
		// There is a limit of 127 characters on the tag key for AWS EC2 instances.
		// Complying to that requirement, do not allow the value of 'key' to be more than that.
		if len(roleEntry.RoleTag) > 127 {
			return logical.ErrorResponse("length of role tag exceeds the EC2 key limit of 127 characters"), nil
		}
	} else if req.Operation == logical.CreateOperation && roleEntry.AuthType == ec2AuthType {
		roleEntry.RoleTag = data.Get("role_tag").(string)
	}

	if roleEntry.HMACKey == "" {
		roleEntry.HMACKey, err = uuid.GenerateUUID()
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate role HMAC key: {{err}}", err)
		}
	}

	if err := b.nonLockedSetAWSRole(ctx, req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}
