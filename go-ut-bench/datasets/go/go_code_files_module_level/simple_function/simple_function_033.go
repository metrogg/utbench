func generateCreationBundle(b *backend, data *dataBundle) error {
	// Read in names -- CN, DNS and email addresses
	var cn string
	var ridSerialNumber string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if data.csr != nil && data.role.UseCSRCommonName {
			cn = data.csr.Subject.CommonName
		}
		if cn == "" {
			cn = data.apiData.Get("common_name").(string)
			if cn == "" && data.role.RequireCN {
				return errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true, unless "require_cn" is set to false`}
			}
		}

		ridSerialNumber = data.apiData.Get("serial_number").(string)

		// only take serial number from CSR if one was not supplied via API
		if ridSerialNumber == "" && data.csr != nil {
			ridSerialNumber = data.csr.Subject.SerialNumber
		}

		if data.csr != nil && data.role.UseCSRSANs {
			dnsNames = data.csr.DNSNames
			emailAddresses = data.csr.EmailAddresses
		}

		if cn != "" && !data.apiData.Get("exclude_cn_from_sans").(bool) {
			if strings.Contains(cn, "@") {
				// Note: emails are not disallowed if the role's email protection
				// flag is false, because they may well be included for
				// informational purposes; it is up to the verifying party to
				// ensure that email addresses in a subject alternate name can be
				// used for the purpose for which they are presented
				emailAddresses = append(emailAddresses, cn)
			} else {
				// Only add to dnsNames if it's actually a DNS name but convert
				// idn first
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToASCII(cn)
				if err != nil {
					return errutil.UserError{Err: err.Error()}
				}
				if hostnameRegex.MatchString(converted) {
					dnsNames = append(dnsNames, converted)
				}
			}
		}

		if data.csr == nil || !data.role.UseCSRSANs {
			cnAltRaw, ok := data.apiData.GetOk("alt_names")
			if ok {
				cnAlt := strutil.ParseDedupLowercaseAndSortStrings(cnAltRaw.(string), ",")
				for _, v := range cnAlt {
					if strings.Contains(v, "@") {
						emailAddresses = append(emailAddresses, v)
					} else {
						// Only add to dnsNames if it's actually a DNS name but
						// convert idn first
						p := idna.New(
							idna.StrictDomainName(true),
							idna.VerifyDNSLength(true),
						)
						converted, err := p.ToASCII(v)
						if err != nil {
							return errutil.UserError{Err: err.Error()}
						}
						if hostnameRegex.MatchString(converted) {
							dnsNames = append(dnsNames, converted)
						}
					}
				}
			}
		}

		// Check the CN. This ensures that the CN is checked even if it's
		// excluded from SANs.
		if cn != "" {
			badName := validateNames(data, []string{cn})
			if len(badName) != 0 {
				return errutil.UserError{Err: fmt.Sprintf(
					"common name %s not allowed by this role", badName)}
			}
		}

		if ridSerialNumber != "" {
			badName := validateSerialNumber(data, ridSerialNumber)
			if len(badName) != 0 {
				return errutil.UserError{Err: fmt.Sprintf(
					"serial_number %s not allowed by this role", badName)}
			}
		}

		// Check for bad email and/or DNS names
		badName := validateNames(data, dnsNames)
		if len(badName) != 0 {
			return errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = validateNames(data, emailAddresses)
		if len(badName) != 0 {
			return errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	var otherSANs map[string][]string
	if sans := data.apiData.Get("other_sans").([]string); len(sans) > 0 {
		requested, err := parseOtherSANs(sans)
		if err != nil {
			return errutil.UserError{Err: errwrap.Wrapf("could not parse requested other SAN: {{err}}", err).Error()}
		}
		badOID, badName, err := validateOtherSANs(data, requested)
		switch {
		case err != nil:
			return errutil.UserError{Err: err.Error()}
		case len(badName) > 0:
			return errutil.UserError{Err: fmt.Sprintf(
				"other SAN %s not allowed for OID %s by this role", badName, badOID)}
		case len(badOID) > 0:
			return errutil.UserError{Err: fmt.Sprintf(
				"other SAN OID %s not allowed by this role", badOID)}
		default:
			otherSANs = requested
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	{
		if data.csr != nil && data.role.UseCSRSANs {
			if len(data.csr.IPAddresses) > 0 {
				if !data.role.AllowIPSANs {
					return errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided some via CSR")}
				}
				ipAddresses = data.csr.IPAddresses
			}
		} else {
			ipAlt := data.apiData.Get("ip_sans").([]string)
			if len(ipAlt) > 0 {
				if !data.role.AllowIPSANs {
					return errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
				}
				for _, v := range ipAlt {
					parsedIP := net.ParseIP(v)
					if parsedIP == nil {
						return errutil.UserError{Err: fmt.Sprintf(
							"the value '%s' is not a valid IP address", v)}
					}
					ipAddresses = append(ipAddresses, parsedIP)
				}
			}
		}
	}

	URIs := []*url.URL{}
	{
		if data.csr != nil && data.role.UseCSRSANs {
			if len(data.csr.URIs) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return errutil.UserError{Err: fmt.Sprintf(
						"URI Subject Alternative Names are not allowed in this role, but were provided via CSR"),
					}
				}

				// validate uri sans
				for _, uri := range data.csr.URIs {
					valid := false
					for _, allowed := range data.role.AllowedURISANs {
						validURI := glob.Glob(allowed, uri.String())
						if validURI {
							valid = true
							break
						}
					}

					if !valid {
						return errutil.UserError{Err: fmt.Sprintf(
							"URI Subject Alternative Names were provided via CSR which are not valid for this role"),
						}
					}

					URIs = append(URIs, uri)
				}
			}
		} else {
			uriAlt := data.apiData.Get("uri_sans").([]string)
			if len(uriAlt) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return errutil.UserError{Err: fmt.Sprintf(
						"URI Subject Alternative Names are not allowed in this role, but were provided via the API"),
					}
				}

				for _, uri := range uriAlt {
					valid := false
					for _, allowed := range data.role.AllowedURISANs {
						validURI := glob.Glob(allowed, uri)
						if validURI {
							valid = true
							break
						}
					}

					if !valid {
						return errutil.UserError{Err: fmt.Sprintf(
							"URI Subject Alternative Names were provided via CSR which are not valid for this role"),
						}
					}

					parsedURI, err := url.Parse(uri)
					if parsedURI == nil || err != nil {
						return errutil.UserError{Err: fmt.Sprintf(
							"the provided URI Subject Alternative Name '%s' is not a valid URI", uri),
						}
					}

					URIs = append(URIs, parsedURI)
				}
			}
		}
	}

	subject := pkix.Name{
		CommonName:         cn,
		SerialNumber:       ridSerialNumber,
		Country:            strutil.RemoveDuplicates(data.role.Country, false),
		Organization:       strutil.RemoveDuplicates(data.role.Organization, false),
		OrganizationalUnit: strutil.RemoveDuplicates(data.role.OU, false),
		Locality:           strutil.RemoveDuplicates(data.role.Locality, false),
		Province:           strutil.RemoveDuplicates(data.role.Province, false),
		StreetAddress:      strutil.RemoveDuplicates(data.role.StreetAddress, false),
		PostalCode:         strutil.RemoveDuplicates(data.role.PostalCode, false),
	}

	// Get the TTL and verify it against the max allowed
	var ttl time.Duration
	var maxTTL time.Duration
	var notAfter time.Time
	{
		ttl = time.Duration(data.apiData.Get("ttl").(int)) * time.Second

		if ttl == 0 && data.role.TTL > 0 {
			ttl = data.role.TTL
		}

		if data.role.MaxTTL > 0 {
			maxTTL = data.role.MaxTTL
		}

		if ttl == 0 {
			ttl = b.System().DefaultLeaseTTL()
		}
		if maxTTL == 0 {
			maxTTL = b.System().MaxLeaseTTL()
		}
		if ttl > maxTTL {
			ttl = maxTTL
		}

		notAfter = time.Now().Add(ttl)

		// If it's not self-signed, verify that the issued certificate won't be
		// valid past the lifetime of the CA certificate
		if data.signingBundle != nil &&
			notAfter.After(data.signingBundle.Certificate.NotAfter) && !data.role.AllowExpirationPastCA {

			return errutil.UserError{Err: fmt.Sprintf(
				"cannot satisfy request, as TTL would result in notAfter %s that is beyond the expiration of the CA certificate at %s", notAfter.Format(time.RFC3339Nano), data.signingBundle.Certificate.NotAfter.Format(time.RFC3339Nano))}
		}
	}

	data.params = &creationParameters{
		Subject:                       subject,
		DNSNames:                      dnsNames,
		EmailAddresses:                emailAddresses,
		IPAddresses:                   ipAddresses,
		URIs:                          URIs,
		OtherSANs:                     otherSANs,
		KeyType:                       data.role.KeyType,
		KeyBits:                       data.role.KeyBits,
		NotAfter:                      notAfter,
		KeyUsage:                      x509.KeyUsage(parseKeyUsages(data.role.KeyUsage)),
		ExtKeyUsage:                   parseExtKeyUsages(data.role),
		ExtKeyUsageOIDs:               data.role.ExtKeyUsageOIDs,
		PolicyIdentifiers:             data.role.PolicyIdentifiers,
		BasicConstraintsValidForNonCA: data.role.BasicConstraintsValidForNonCA,
		NotBeforeDuration:             data.role.NotBeforeDuration,
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if data.signingBundle == nil {
		return nil
	}

	// This will have been read in from the getURLs function
	data.params.URLs = data.signingBundle.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if data.role.MaxPathLength != nil {
		data.params.MaxPathLength = *data.role.MaxPathLength
	} else {
		switch {
		case data.signingBundle.Certificate.MaxPathLen < 0:
			data.params.MaxPathLength = -1
		case data.signingBundle.Certificate.MaxPathLen == 0 &&
			data.signingBundle.Certificate.MaxPathLenZero:
			// The signing function will ensure that we do not issue a CA cert
			data.params.MaxPathLength = 0
		default:
			// If this takes it to zero, we handle this case later if
			// necessary
			data.params.MaxPathLength = data.signingBundle.Certificate.MaxPathLen - 1
		}
	}

	return nil
}
