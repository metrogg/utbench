func (b *backend) pathLoginUpdateEc2(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	identityDocB64 := data.Get("identity").(string)
	var identityDocBytes []byte
	var err error
	if identityDocB64 != "" {
		identityDocBytes, err = base64.StdEncoding.DecodeString(identityDocB64)
		if err != nil || len(identityDocBytes) == 0 {
			return logical.ErrorResponse("failed to base64 decode the instance identity document"), nil
		}
	}

	signatureB64 := data.Get("signature").(string)
	var signatureBytes []byte
	if signatureB64 != "" {
		signatureBytes, err = base64.StdEncoding.DecodeString(signatureB64)
		if err != nil {
			return logical.ErrorResponse("failed to base64 decode the SHA256 RSA signature of the instance identity document"), nil
		}
	}

	pkcs7B64 := data.Get("pkcs7").(string)

	// Either the pkcs7 signature of the instance identity document, or
	// the identity document itself along with its SHA256 RSA signature
	// needs to be provided.
	if pkcs7B64 == "" && (len(identityDocBytes) == 0 && len(signatureBytes) == 0) {
		return logical.ErrorResponse("either pkcs7 or a tuple containing the instance identity document and its SHA256 RSA signature needs to be provided"), nil
	} else if pkcs7B64 != "" && (len(identityDocBytes) != 0 && len(signatureBytes) != 0) {
		return logical.ErrorResponse("both pkcs7 and a tuple containing the instance identity document and its SHA256 RSA signature is supplied; provide only one"), nil
	}

	// Verify the signature of the identity document and unmarshal it
	var identityDocParsed *identityDocument
	if pkcs7B64 != "" {
		identityDocParsed, err = b.parseIdentityDocument(ctx, req.Storage, pkcs7B64)
		if err != nil {
			return nil, err
		}
		if identityDocParsed == nil {
			return logical.ErrorResponse("failed to verify the instance identity document using pkcs7"), nil
		}
	} else {
		identityDocParsed, err = b.verifyInstanceIdentitySignature(ctx, req.Storage, identityDocBytes, signatureBytes)
		if err != nil {
			return nil, err
		}
		if identityDocParsed == nil {
			return logical.ErrorResponse("failed to verify the instance identity document using the SHA256 RSA digest"), nil
		}
	}

	roleName := data.Get("role").(string)

	// If roleName is not supplied, a role in the name of the instance's AMI ID will be looked for
	if roleName == "" {
		roleName = identityDocParsed.AmiID
	}

	// Get the entry for the role used by the instance
	roleEntry, err := b.lockedAWSRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry for role %q not found", roleName)), nil
	}

	if roleEntry.AuthType != ec2AuthType {
		return logical.ErrorResponse(fmt.Sprintf("auth method ec2 not allowed for role %s", roleName)), nil
	}

	identityConfigEntry, err := identityConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	identityAlias := ""

	switch identityConfigEntry.EC2Alias {
	case identityAliasRoleID:
		identityAlias = roleEntry.RoleID
	case identityAliasEC2InstanceID:
		identityAlias = identityDocParsed.InstanceID
	case identityAliasEC2ImageID:
		identityAlias = identityDocParsed.AmiID
	}

	// If we're just looking up for MFA, return the Alias info
	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: identityAlias,
				},
			},
		}, nil
	}

	// Validate the instance ID by making a call to AWS EC2 DescribeInstances API
	// and fetching the instance description. Validation succeeds only if the
	// instance is in 'running' state.
	instance, err := b.validateInstance(ctx, req.Storage, identityDocParsed.InstanceID, identityDocParsed.Region, identityDocParsed.AccountID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to verify instance ID: %v", err)), nil
	}

	// Verify that the `Region` of the instance trying to login matches the
	// `Region` specified as a constraint on role
	if len(roleEntry.BoundRegions) > 0 && !strutil.StrListContains(roleEntry.BoundRegions, identityDocParsed.Region) {
		return logical.ErrorResponse(fmt.Sprintf("Region %q does not satisfy the constraint on role %q", identityDocParsed.Region, roleName)), nil
	}

	validationError, err := b.verifyInstanceMeetsRoleRequirements(ctx, req.Storage, instance, roleEntry, roleName, identityDocParsed)
	if err != nil {
		return nil, err
	}
	if validationError != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error validating instance: %v", validationError)), nil
	}

	// Get the entry from the identity whitelist, if there is one
	storedIdentity, err := whitelistIdentityEntry(ctx, req.Storage, identityDocParsed.InstanceID)
	if err != nil {
		return nil, err
	}

	// disallowReauthentication value that gets cached at the stored
	// identity-whitelist entry is determined not just by the role entry.
	// If client explicitly sets nonce to be empty, it implies intent to
	// disable reauthentication. Also, role tag can override the 'false'
	// value with 'true' (the other way around is not allowed).

	// Read the value from the role entry
	disallowReauthentication := roleEntry.DisallowReauthentication

	clientNonce := ""

	// Check if the nonce is supplied by the client
	clientNonceRaw, clientNonceSupplied := data.GetOk("nonce")
	if clientNonceSupplied {
		clientNonce = clientNonceRaw.(string)

		// Nonce explicitly set to empty implies intent to disable
		// reauthentication by the client. Set a predefined nonce which
		// indicates reauthentication being disabled.
		if clientNonce == "" {
			clientNonce = reauthenticationDisabledNonce

			// Ensure that the intent lands in the whitelist
			disallowReauthentication = true
		}
	}

	// This is NOT a first login attempt from the client
	if storedIdentity != nil {
		// Check if the client nonce match the cached nonce and if the pending time
		// of the identity document is not before the pending time of the document
		// with which previous login was made. If 'allow_instance_migration' is
		// enabled on the registered role, client nonce requirement is relaxed.
		if err = validateMetadata(clientNonce, identityDocParsed.PendingTime, storedIdentity, roleEntry); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		// Don't let subsequent login attempts to bypass the initial
		// intent of disabling reauthentication, despite the properties
		// of role getting updated. For example: Role has the value set
		// to 'false', a role-tag login sets the value to 'true', then
		// role gets updated to not use a role-tag, and a login attempt
		// is made with role's value set to 'false'. Removing the entry
		// from the identity-whitelist should be the only way to be
		// able to login from the instance again.
		disallowReauthentication = disallowReauthentication || storedIdentity.DisallowReauthentication
	}

	// If we reach this point without erroring and if the client nonce was
	// not supplied, a first time login is implied and that the client
	// intends that the nonce be generated by the backend. Create a random
	// nonce to be associated for the instance ID.
	if !clientNonceSupplied {
		if clientNonce, err = uuid.GenerateUUID(); err != nil {
			return nil, fmt.Errorf("failed to generate random nonce")
		}
	}

	// Load the current values for max TTL and policies from the role entry,
	// before checking for overriding max TTL in the role tag.  The shortest
	// max TTL is used to cap the token TTL; the longest max TTL is used to
	// make the whitelist entry as long as possible as it controls for replay
	// attacks.
	shortestMaxTTL := b.System().MaxLeaseTTL()
	longestMaxTTL := b.System().MaxLeaseTTL()
	if roleEntry.MaxTTL > time.Duration(0) && roleEntry.MaxTTL < shortestMaxTTL {
		shortestMaxTTL = roleEntry.MaxTTL
	}
	if roleEntry.MaxTTL > longestMaxTTL {
		longestMaxTTL = roleEntry.MaxTTL
	}

	policies := roleEntry.Policies
	rTagMaxTTL := time.Duration(0)
	var roleTagResp *roleTagLoginResponse
	if roleEntry.RoleTag != "" {
		roleTagResp, err = b.handleRoleTagLogin(ctx, req.Storage, roleName, roleEntry, instance)
		if err != nil {
			return nil, err
		}
		if roleTagResp == nil {
			return logical.ErrorResponse("failed to fetch and verify the role tag"), nil
		}
	}

	if roleTagResp != nil {
		// Role tag is enabled on the role.

		// Overwrite the policies with the ones returned from processing the role tag
		// If there are no policies on the role tag, policies on the role are inherited.
		// If policies on role tag are set, by this point, it is verified that it is a subset of the
		// policies on the role. So, apply only those.
		if len(roleTagResp.Policies) != 0 {
			policies = roleTagResp.Policies
		}

		// If roleEntry had disallowReauthentication set to 'true', do not reset it
		// to 'false' based on role tag having it not set. But, if role tag had it set,
		// be sure to override the value.
		if !disallowReauthentication {
			disallowReauthentication = roleTagResp.DisallowReauthentication
		}

		// Cache the value of role tag's max_ttl value
		rTagMaxTTL = roleTagResp.MaxTTL

		// Scope the shortestMaxTTL to the value set on the role tag
		if roleTagResp.MaxTTL > time.Duration(0) && roleTagResp.MaxTTL < shortestMaxTTL {
			shortestMaxTTL = roleTagResp.MaxTTL
		}
		if roleTagResp.MaxTTL > longestMaxTTL {
			longestMaxTTL = roleTagResp.MaxTTL
		}
	}

	// Save the login attempt in the identity whitelist
	currentTime := time.Now()
	if storedIdentity == nil {
		// Role, ClientNonce and CreationTime of the identity entry,
		// once set, should never change.
		storedIdentity = &whitelistIdentity{
			Role:         roleName,
			ClientNonce:  clientNonce,
			CreationTime: currentTime,
		}
	}

	// DisallowReauthentication, PendingTime, LastUpdatedTime and
	// ExpirationTime may change.
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(longestMaxTTL)
	storedIdentity.PendingTime = identityDocParsed.PendingTime
	storedIdentity.DisallowReauthentication = disallowReauthentication

	// Don't cache the nonce if DisallowReauthentication is set
	if storedIdentity.DisallowReauthentication {
		storedIdentity.ClientNonce = ""
	}

	// Sanitize the nonce to a reasonable length
	if len(clientNonce) > 128 && !storedIdentity.DisallowReauthentication {
		return logical.ErrorResponse("client nonce exceeding the limit of 128 characters"), nil
	}

	if err = setWhitelistIdentityEntry(ctx, req.Storage, identityDocParsed.InstanceID, storedIdentity); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period:   roleEntry.Period,
			Policies: policies,
			Metadata: map[string]string{
				"instance_id":      identityDocParsed.InstanceID,
				"region":           identityDocParsed.Region,
				"account_id":       identityDocParsed.AccountID,
				"role_tag_max_ttl": rTagMaxTTL.String(),
				"role":             roleName,
				"ami_id":           identityDocParsed.AmiID,
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       roleEntry.TTL,
				MaxTTL:    shortestMaxTTL,
			},
			Alias: &logical.Alias{
				Name: identityAlias,
			},
		},
	}

	// Return the nonce only if reauthentication is allowed and if the nonce
	// was not supplied by the user.
	if !disallowReauthentication && !clientNonceSupplied {
		// Echo the client nonce back. If nonce param was not supplied
		// to the endpoint at all (setting it to empty string does not
		// qualify here), callers should extract out the nonce from
		// this field for reauthentication requests.
		resp.Auth.Metadata["nonce"] = clientNonce
	}

	return resp, nil
}
