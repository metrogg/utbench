func validateCAConfig(ctx context.Context, securityConfig *ca.SecurityConfig, cluster *api.Cluster) (*api.RootCA, error) {
	newConfig := cluster.Spec.CAConfig.Copy()
	newConfig.SigningCACert = ca.NormalizePEMs(newConfig.SigningCACert) // ensure this is normalized before we use it

	if len(newConfig.SigningCAKey) > 0 && len(newConfig.SigningCACert) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "if a signing CA key is provided, the signing CA cert must also be provided")
	}

	normalizedRootCA := ca.NormalizePEMs(cluster.RootCA.CACert)
	extCAs, err := getNormalizedExtCAs(newConfig, normalizedRootCA) // validate that the list of external CAs is not malformed
	if err != nil {
		return nil, err
	}

	var oldCertExtCAs []*api.ExternalCA
	if !hasSigningKey(&cluster.RootCA) {

		// If we are going from external -> internal, but providing the external CA's signing key,
		// then we don't need to validate any external CAs.  We can in fact abort any outstanding root
		// rotations if we are just adding a key.  Because we have a key, we don't care if there are
		// no external CAs matching the certificate.
		if bytes.Equal(normalizedRootCA, newConfig.SigningCACert) && hasSigningKey(newConfig) {
			// validate that the key and cert indeed match - if they don't then just fail now rather
			// than go through all the external CA URLs, which is a more expensive operation
			if _, err := ca.NewRootCA(newConfig.SigningCACert, newConfig.SigningCACert, newConfig.SigningCAKey, ca.DefaultNodeCertExpiration, nil); err != nil {
				return nil, err
			}
			copied := cluster.RootCA.Copy()
			copied.CAKey = newConfig.SigningCAKey
			copied.RootRotation = nil
			copied.LastForcedRotation = newConfig.ForceRotate
			return copied, nil
		}

		oldCertExtCAs, err = validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, normalizedRootCA, "current")
		if err != nil {
			return nil, err
		}
	}

	// if the desired CA cert and key are not set, then we are happy with the current root CA configuration, unless
	// the ForceRotate version has changed
	if len(newConfig.SigningCACert) == 0 {
		if cluster.RootCA.LastForcedRotation != newConfig.ForceRotate {
			newRootCA, err := ca.CreateRootCA(ca.DefaultRootCN)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return newRootRotationObject(ctx, securityConfig, &cluster.RootCA, newRootCA, oldCertExtCAs, newConfig.ForceRotate)
		}

		// we also need to make sure that if the current root rotation requires an external CA, those external CAs are
		// still valid
		if cluster.RootCA.RootRotation != nil && !hasSigningKey(cluster.RootCA.RootRotation) {
			_, err := validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, ca.NormalizePEMs(cluster.RootCA.RootRotation.CACert), "next")
			if err != nil {
				return nil, err
			}
		}

		return &cluster.RootCA, nil // no change, return as is
	}

	// A desired cert and maybe key were provided - we need to make sure the cert and key (if provided) match.
	var signingCert []byte
	if hasSigningKey(newConfig) {
		signingCert = newConfig.SigningCACert
	}
	newRootCA, err := ca.NewRootCA(newConfig.SigningCACert, signingCert, newConfig.SigningCAKey, ca.DefaultNodeCertExpiration, nil)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if len(newRootCA.Pool.Subjects()) != 1 {
		return nil, status.Errorf(codes.InvalidArgument, "the desired CA certificate cannot contain multiple certificates")
	}

	parsedCert, err := helpers.ParseCertificatePEM(newConfig.SigningCACert)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not parse the desired CA certificate")
	}

	// The new certificate's expiry must be at least one year away
	if parsedCert.NotAfter.Before(time.Now().Add(minRootExpiration)) {
		return nil, status.Errorf(codes.InvalidArgument, "CA certificate expires too soon")
	}

	if !hasSigningKey(newConfig) {
		if _, err := validateHasAtLeastOneExternalCA(ctx, extCAs, securityConfig, newConfig.SigningCACert, "desired"); err != nil {
			return nil, err
		}
	}

	// check if we can abort any existing root rotations
	if bytes.Equal(normalizedRootCA, newConfig.SigningCACert) {
		copied := cluster.RootCA.Copy()
		copied.CAKey = newConfig.SigningCAKey
		copied.RootRotation = nil
		copied.LastForcedRotation = newConfig.ForceRotate
		return copied, nil
	}

	// check if this is the same desired cert as an existing root rotation
	if r := cluster.RootCA.RootRotation; r != nil && bytes.Equal(ca.NormalizePEMs(r.CACert), newConfig.SigningCACert) {
		copied := cluster.RootCA.Copy()
		copied.RootRotation.CAKey = newConfig.SigningCAKey
		copied.LastForcedRotation = newConfig.ForceRotate
		return copied, nil
	}

	// ok, everything's different; we have to begin a new root rotation which means generating a new cross-signed cert
	return newRootRotationObject(ctx, securityConfig, &cluster.RootCA, newRootCA, oldCertExtCAs, newConfig.ForceRotate)
}
