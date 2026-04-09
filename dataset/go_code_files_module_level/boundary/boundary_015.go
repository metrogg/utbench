func (ts *TokenStore) handleCreateCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, orphan bool, role *tsRoleEntry) (*logical.Response, error) {
	// Read the parent policy
	parent, err := ts.Lookup(ctx, req.ClientToken)
	if err != nil {
		return nil, errwrap.Wrapf("parent token lookup failed: {{err}}", err)
	}
	if parent == nil {
		return logical.ErrorResponse("parent token lookup failed: no parent found"), logical.ErrInvalidRequest
	}
	if parent.Type == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot create more tokens"), nil
	}

	// A token with a restricted number of uses cannot create a new token
	// otherwise it could escape the restriction count.
	if parent.NumUses > 0 {
		return logical.ErrorResponse("restricted use token cannot generate child tokens"),
			logical.ErrInvalidRequest
	}

	// Check if the client token has sudo/root privileges for the requested path
	isSudo := ts.System().SudoPrivilege(ctx, req.MountPoint+req.Path, req.ClientToken)

	// Read and parse the fields
	var data struct {
		ID              string
		Policies        []string
		Metadata        map[string]string `mapstructure:"meta"`
		NoParent        bool              `mapstructure:"no_parent"`
		NoDefaultPolicy bool              `mapstructure:"no_default_policy"`
		Lease           string
		TTL             string
		Renewable       *bool
		ExplicitMaxTTL  string `mapstructure:"explicit_max_ttl"`
		DisplayName     string `mapstructure:"display_name"`
		NumUses         int    `mapstructure:"num_uses"`
		Period          string
		Type            string `mapstructure:"type"`
	}
	if err := mapstructure.WeakDecode(req.Data, &data); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding request: %s", err)), logical.ErrInvalidRequest
	}

	// If the context's namespace is different from the parent and this is an
	// orphan token creation request, then this is an admin token generation for
	// the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != parent.NamespaceID {
		parentNS, err := NamespaceByID(ctx, parent.NamespaceID, ts.core)
		if err != nil {
			ts.logger.Error("error looking up parent namespace", "error", err, "parent_namespace", parent.NamespaceID)
			return nil, ErrInternalError
		}
		if parentNS == nil {
			ts.logger.Error("could not find information for parent namespace", "parent_namespace", parent.NamespaceID)
			return nil, ErrInternalError
		}

		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to directly generate a token in a child namespace"), logical.ErrInvalidRequest
		}

		if strutil.StrListContains(data.Policies, "root") {
			return logical.ErrorResponse("root tokens may not be created from a parent namespace"), logical.ErrInvalidRequest
		}
	}

	renewable := true
	if data.Renewable != nil {
		renewable = *data.Renewable
	}

	tokenType := logical.TokenTypeService
	tokenTypeStr := data.Type
	if role != nil {
		switch role.TokenType {
		case logical.TokenTypeDefault, logical.TokenTypeDefaultService:
			// Use the user-given value, but fall back to service
		case logical.TokenTypeDefaultBatch:
			// Use the user-given value, but fall back to batch
			if tokenTypeStr == "" {
				tokenTypeStr = logical.TokenTypeBatch.String()
			}
		case logical.TokenTypeService:
			tokenTypeStr = logical.TokenTypeService.String()
		case logical.TokenTypeBatch:
			tokenTypeStr = logical.TokenTypeBatch.String()
		default:
			return logical.ErrorResponse(fmt.Sprintf("role being used for token creation contains invalid token type %q", role.TokenType.String())), nil
		}
	}
	switch tokenTypeStr {
	case "", "service":
	case "batch":
		var badReason string
		switch {
		case data.ExplicitMaxTTL != "":
			dur, err := parseutil.ParseDurationSecond(data.ExplicitMaxTTL)
			if err != nil {
				return logical.ErrorResponse(`"explicit_max_ttl" value could not be parsed`), nil
			}
			if dur != 0 {
				badReason = "explicit_max_ttl"
			}
		case data.NumUses != 0:
			badReason = "num_uses"
		case data.Period != "":
			dur, err := parseutil.ParseDurationSecond(data.Period)
			if err != nil {
				return logical.ErrorResponse(`"period" value could not be parsed`), nil
			}
			if dur != 0 {
				badReason = "period"
			}
		}
		if badReason != "" {
			return logical.ErrorResponse(fmt.Sprintf("batch tokens cannot have %q set", badReason)), nil
		}
		tokenType = logical.TokenTypeBatch
		renewable = false
	default:
		return logical.ErrorResponse("invalid 'token_type' value"), logical.ErrInvalidRequest
	}

	// Verify the number of uses is positive
	if data.NumUses < 0 {
		return logical.ErrorResponse("number of uses cannot be negative"),
			logical.ErrInvalidRequest
	}

	// Setup the token entry
	te := logical.TokenEntry{
		Parent: req.ClientToken,

		// The mount point is always the same since we have only one token
		// store; using req.MountPoint causes trouble in tests since they don't
		// have an official mount
		Path: fmt.Sprintf("auth/token/%s", req.Path),

		Meta:         data.Metadata,
		DisplayName:  "token",
		NumUses:      data.NumUses,
		CreationTime: time.Now().Unix(),
		NamespaceID:  ns.ID,
		Type:         tokenType,
	}

	// If the role is not nil, we add the role name as part of the token's
	// path. This makes it much easier to later revoke tokens that were issued
	// by a role (using revoke-prefix). Users can further specify a PathSuffix
	// in the role; that way they can use something like "v1", "v2" to indicate
	// role revisions, and revoke only tokens issued with a previous revision.
	if role != nil {
		te.Role = role.Name

		// If renewable hasn't been disabled in the call and the role has
		// renewability disabled, set renewable false
		if renewable && !role.Renewable {
			renewable = false
		}

		if role.PathSuffix != "" {
			te.Path = fmt.Sprintf("%s/%s", te.Path, role.PathSuffix)
		}
	}

	// Attach the given display name if any
	if data.DisplayName != "" {
		full := "token-" + data.DisplayName
		full = displayNameSanitize.ReplaceAllString(full, "-")
		full = strings.TrimSuffix(full, "-")
		te.DisplayName = full
	}

	// Allow specifying the ID of the token if the client has root or sudo privileges
	if data.ID != "" {
		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to specify token id"),
				logical.ErrInvalidRequest
		}
		if ns.ID != namespace.RootNamespaceID {
			return logical.ErrorResponse("token IDs can only be manually specified in the root namespace"),
				logical.ErrInvalidRequest
		}
		te.ID = data.ID
	}

	resp := &logical.Response{}

	var addDefault bool

	// N.B.: The logic here uses various calculations as to whether default
	// should be added. In the end we decided that if NoDefaultPolicy is set it
	// should be stripped out regardless, *but*, the logic of when it should
	// and shouldn't be added is kept because we want to do subset comparisons
	// based on adding default when it's correct to do so.
	switch {
	case role != nil && (len(role.AllowedPolicies) > 0 || len(role.DisallowedPolicies) > 0):
		// Holds the final set of policies as they get munged
		var finalPolicies []string

		// We don't make use of the global one because roles with allowed or
		// disallowed set do their own policy rules
		var localAddDefault bool

		// If the request doesn't say not to add "default" and if "default"
		// isn't in the disallowed list, add it. This is in line with the idea
		// that roles, when allowed/disallowed ar set, allow a subset of
		// policies to be set disjoint from the parent token's policies.
		if !data.NoDefaultPolicy && !strutil.StrListContains(role.DisallowedPolicies, "default") {
			localAddDefault = true
		}

		// Start with passed-in policies as a baseline, if they exist
		if len(data.Policies) > 0 {
			finalPolicies = policyutil.SanitizePolicies(data.Policies, localAddDefault)
		}

		var sanitizedRolePolicies []string

		// First check allowed policies; if policies are specified they will be
		// checked, otherwise if an allowed set exists that will be the set
		// that is used
		if len(role.AllowedPolicies) > 0 {
			// Note that if "default" is already in allowed, and also in
			// disallowed, this will still result in an error later since this
			// doesn't strip out default
			sanitizedRolePolicies = policyutil.SanitizePolicies(role.AllowedPolicies, localAddDefault)

			if len(finalPolicies) == 0 {
				finalPolicies = sanitizedRolePolicies
			} else {
				if !strutil.StrListSubset(sanitizedRolePolicies, finalPolicies) {
					return logical.ErrorResponse(fmt.Sprintf("token policies (%q) must be subset of the role's allowed policies (%q)", finalPolicies, sanitizedRolePolicies)), logical.ErrInvalidRequest
				}
			}
		} else {
			// Assign parent policies if none have been requested. As this is a
			// role, add default unless explicitly disabled.
			if len(finalPolicies) == 0 {
				finalPolicies = policyutil.SanitizePolicies(parent.Policies, localAddDefault)
			}
		}

		if len(role.DisallowedPolicies) > 0 {
			// We don't add the default here because we only want to disallow it if it's explicitly set
			sanitizedRolePolicies = strutil.RemoveDuplicates(role.DisallowedPolicies, true)

			for _, finalPolicy := range finalPolicies {
				if strutil.StrListContains(sanitizedRolePolicies, finalPolicy) {
					return logical.ErrorResponse(fmt.Sprintf("token policy %q is disallowed by this role", finalPolicy)), logical.ErrInvalidRequest
				}
			}
		}

		data.Policies = finalPolicies

	// We are creating a token from a parent namespace. We should only use the input
	// policies.
	case ns.ID != parent.NamespaceID:
		addDefault = !data.NoDefaultPolicy

	// No policies specified, inherit parent
	case len(data.Policies) == 0:
		// Only inherit "default" if the parent already has it, so don't touch addDefault here
		data.Policies = policyutil.SanitizePolicies(parent.Policies, policyutil.DoNotAddDefaultPolicy)

	// When a role is not in use or does not specify allowed/disallowed, only
	// permit policies to be a subset unless the client has root or sudo
	// privileges. Default is added in this case if the parent has it, unless
	// the client specified for it not to be added.
	case !isSudo:
		// Sanitize passed-in and parent policies before comparison
		sanitizedInputPolicies := policyutil.SanitizePolicies(data.Policies, policyutil.DoNotAddDefaultPolicy)
		sanitizedParentPolicies := policyutil.SanitizePolicies(parent.Policies, policyutil.DoNotAddDefaultPolicy)

		if !strutil.StrListSubset(sanitizedParentPolicies, sanitizedInputPolicies) {
			return logical.ErrorResponse("child policies must be subset of parent"), logical.ErrInvalidRequest
		}

		// If the parent has default, and they haven't requested not to get it,
		// add it. Note that if they have explicitly put "default" in
		// data.Policies it will still be added because NoDefaultPolicy
		// controls *automatic* adding.
		if !data.NoDefaultPolicy && strutil.StrListContains(parent.Policies, "default") {
			addDefault = true
		}

	// Add default by default in this case unless requested not to
	case isSudo:
		addDefault = !data.NoDefaultPolicy
	}

	te.Policies = policyutil.SanitizePolicies(data.Policies, addDefault)

	// Yes, this is a little inefficient to do it like this, but meh
	if data.NoDefaultPolicy {
		te.Policies = strutil.StrListDelete(te.Policies, "default")
	}

	// Prevent internal policies from being assigned to tokens
	for _, policy := range te.Policies {
		if strutil.StrListContains(nonAssignablePolicies, policy) {
			return logical.ErrorResponse(fmt.Sprintf("cannot assign policy %q", policy)), nil
		}
	}

	if strutil.StrListContains(te.Policies, "root") {
		// Prevent attempts to create a root token without an actual root token as parent.
		// This is to thwart privilege escalation by tokens having 'sudo' privileges.
		if !strutil.StrListContains(parent.Policies, "root") {
			return logical.ErrorResponse("root tokens may not be created without parent token being root"), logical.ErrInvalidRequest
		}

		if te.Type == logical.TokenTypeBatch {
			// Batch tokens cannot be revoked so we should never have root batch tokens
			return logical.ErrorResponse("batch tokens cannot be root tokens"), nil
		}
	}

	//
	// NOTE: Do not modify policies below this line. We need the checks above
	// to be the last checks as they must look at the final policy set.
	//

	switch {
	case role != nil:
		if role.Orphan {
			te.Parent = ""
		}

		if len(role.BoundCIDRs) > 0 {
			te.BoundCIDRs = role.BoundCIDRs
		}

	case data.NoParent:
		// Only allow an orphan token if the client has sudo policy
		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to create orphan token"),
				logical.ErrInvalidRequest
		}

		te.Parent = ""

	default:
		// This comes from create-orphan, which can be properly ACLd
		if orphan {
			te.Parent = ""
		}
	}

	// At this point, it is clear whether the token is going to be an orphan or
	// not. If the token is not going to be an orphan, inherit the parent's
	// entity identifier into the child token.
	if te.Parent != "" {
		te.EntityID = parent.EntityID

		// If the parent has bound CIDRs, copy those into the child. We don't
		// do this if role is not nil because then we always use the role's
		// bound CIDRs; roles allow escalation of privilege in proper
		// circumstances.
		if role == nil {
			te.BoundCIDRs = parent.BoundCIDRs
		}
	}

	var explicitMaxTTLToUse time.Duration
	if data.ExplicitMaxTTL != "" {
		dur, err := parseutil.ParseDurationSecond(data.ExplicitMaxTTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("explicit_max_ttl must be positive"), logical.ErrInvalidRequest
		}
		te.ExplicitMaxTTL = dur
		explicitMaxTTLToUse = dur
	}

	var periodToUse time.Duration
	if data.Period != "" {
		dur, err := parseutil.ParseDurationSecond(data.Period)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}

		switch {
		case dur < 0:
			return logical.ErrorResponse("period must be positive"), logical.ErrInvalidRequest
		case dur == 0:
		default:
			if !isSudo {
				return logical.ErrorResponse("root or sudo privileges required to create periodic token"),
					logical.ErrInvalidRequest
			}
			te.Period = dur
			periodToUse = dur
		}
	}

	// Parse the TTL/lease if any
	if data.TTL != "" {
		dur, err := parseutil.ParseDurationSecond(data.TTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("ttl must be positive"), logical.ErrInvalidRequest
		}
		te.TTL = dur
	} else if data.Lease != "" {
		// This block is compatibility
		dur, err := time.ParseDuration(data.Lease)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("lease must be positive"), logical.ErrInvalidRequest
		}
		te.TTL = dur
	}

	// Set the lesser period/explicit max TTL if defined both in arguments and
	// in role. Batch tokens will error out if not set via role, but here we
	// need to explicitly check
	if role != nil && te.Type != logical.TokenTypeBatch {
		if role.ExplicitMaxTTL != 0 {
			switch {
			case explicitMaxTTLToUse == 0:
				explicitMaxTTLToUse = role.ExplicitMaxTTL
			default:
				if role.ExplicitMaxTTL < explicitMaxTTLToUse {
					explicitMaxTTLToUse = role.ExplicitMaxTTL
				}
				resp.AddWarning(fmt.Sprintf("Explicit max TTL specified both during creation call and in role; using the lesser value of %d seconds", int64(explicitMaxTTLToUse.Seconds())))
			}
		}
		if role.Period != 0 {
			switch {
			case periodToUse == 0:
				periodToUse = role.Period
			default:
				if role.Period < periodToUse {
					periodToUse = role.Period
				}
				resp.AddWarning(fmt.Sprintf("Period specified both during creation call and in role; using the lesser value of %d seconds", int64(periodToUse.Seconds())))
			}
		}
	}

	sysView := ts.System()

	// Only calculate a TTL if you are A) periodic, B) have a TTL, C) do not have a TTL and are not a root token
	if periodToUse > 0 || te.TTL > 0 || (te.TTL == 0 && !strutil.StrListContains(te.Policies, "root")) {
		ttl, warnings, err := framework.CalculateTTL(sysView, 0, te.TTL, periodToUse, 0, explicitMaxTTLToUse, time.Unix(te.CreationTime, 0))
		if err != nil {
			return nil, err
		}
		for _, warning := range warnings {
			resp.AddWarning(warning)
		}
		te.TTL = ttl
	}

	// Root tokens are still bound by explicit max TTL
	if te.TTL == 0 && explicitMaxTTLToUse > 0 {
		te.TTL = explicitMaxTTLToUse
	}

	// Don't advertise non-expiring root tokens as renewable, as attempts to
	// renew them are denied. Don't CIDR-restrict these either.
	if te.TTL == 0 {
		if parent.TTL != 0 {
			return logical.ErrorResponse("expiring root tokens cannot create non-expiring root tokens"), logical.ErrInvalidRequest
		}
		renewable = false
		te.BoundCIDRs = nil
	}

	if te.ID != "" {
		resp.AddWarning("Supplying a custom ID for the token uses the weaker SHA1 hashing instead of the more secure SHA2-256 HMAC for token obfuscation. SHA1 hashed tokens on the wire leads to less secure lookups.")
	}

	// Create the token
	if err := ts.create(ctx, &te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp.Auth = &logical.Auth{
		NumUses:     te.NumUses,
		DisplayName: te.DisplayName,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: renewable,
		},
		ClientToken:    te.ID,
		Accessor:       te.Accessor,
		EntityID:       te.EntityID,
		Period:         periodToUse,
		ExplicitMaxTTL: explicitMaxTTLToUse,
		CreationPath:   te.Path,
		TokenType:      te.Type,
		Orphan:         te.Parent == "",
	}

	for _, p := range te.Policies {
		policy, err := ts.core.policyStore.GetPolicy(ctx, p, PolicyTypeToken)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("could not look up policy %s", p)), nil
		}
		if policy == nil {
			resp.AddWarning(fmt.Sprintf("Policy %q does not exist", p))
		}
	}

	return resp, nil
}
