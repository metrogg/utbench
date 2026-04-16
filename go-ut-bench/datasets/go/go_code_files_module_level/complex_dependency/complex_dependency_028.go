func (c *vaultClient) renew(req *vaultClientRenewalRequest) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if req == nil {
		return fmt.Errorf("nil renewal request")
	}
	if req.errCh == nil {
		return fmt.Errorf("renewal request error channel nil")
	}

	if !c.config.IsEnabled() {
		close(req.errCh)
		return fmt.Errorf("vault client not enabled")
	}
	if !c.running {
		close(req.errCh)
		return fmt.Errorf("vault client is not running")
	}
	if req.id == "" {
		close(req.errCh)
		return fmt.Errorf("missing id in renewal request")
	}
	if req.increment < 1 {
		close(req.errCh)
		return fmt.Errorf("increment cannot be less than 1")
	}

	var renewalErr error
	leaseDuration := req.increment
	if req.isToken {
		// Set the token in the API client to the one that needs
		// renewal
		c.client.SetToken(req.id)

		// Renew the token
		renewResp, err := c.client.Auth().Token().RenewSelf(req.increment)
		if err != nil {
			renewalErr = fmt.Errorf("failed to renew the vault token: %v", err)
		} else if renewResp == nil || renewResp.Auth == nil {
			renewalErr = fmt.Errorf("failed to renew the vault token")
		} else {
			// Don't set this if renewal fails
			leaseDuration = renewResp.Auth.LeaseDuration
		}

		// Reset the token in the API client before returning
		c.client.SetToken("")
	} else {
		// Renew the secret
		renewResp, err := c.client.Sys().Renew(req.id, req.increment)
		if err != nil {
			renewalErr = fmt.Errorf("failed to renew vault secret: %v", err)
		} else if renewResp == nil {
			renewalErr = fmt.Errorf("failed to renew vault secret")
		} else {
			// Don't set this if renewal fails
			leaseDuration = renewResp.LeaseDuration
		}
	}

	// Determine the next renewal time
	renewalDuration := renewalTime(rand.Intn, leaseDuration)
	next := time.Now().Add(renewalDuration)

	fatal := false
	if renewalErr != nil &&
		(strings.Contains(renewalErr.Error(), "lease not found or lease is not renewable") ||
			strings.Contains(renewalErr.Error(), "lease is not renewable") ||
			strings.Contains(renewalErr.Error(), "token not found") ||
			strings.Contains(renewalErr.Error(), "permission denied")) {
		fatal = true
	} else if renewalErr != nil {
		c.logger.Debug("renewal error details", "req.increment", req.increment, "lease_duration", leaseDuration, "renewal_duration", renewalDuration)
		c.logger.Error("error during renewal of lease or token failed due to a non-fatal error; retrying",
			"error", renewalErr, "period", next)
	}

	if c.isTracked(req.id) {
		if fatal {
			// If encountered with an error where in a lease or a
			// token is not valid at all with vault, and if that
			// item is tracked by the renewal loop, stop renewing
			// it by removing the corresponding heap entry.
			if err := c.heap.Remove(req.id); err != nil {
				return fmt.Errorf("failed to remove heap entry: %v", err)
			}

			// Report the fatal error to the client
			req.errCh <- renewalErr
			close(req.errCh)

			return renewalErr
		}

		// If the identifier is already tracked, this indicates a
		// subsequest renewal. In this case, update the existing
		// element in the heap with the new renewal time.
		if err := c.heap.Update(req, next); err != nil {
			return fmt.Errorf("failed to update heap entry. err: %v", err)
		}

		// There is no need to signal an update to the renewal loop
		// here because this case is hit from the renewal loop itself.
	} else {
		if fatal {
			// If encountered with an error where in a lease or a
			// token is not valid at all with vault, and if that
			// item is not tracked by renewal loop, don't add it.

			// Report the fatal error to the client
			req.errCh <- renewalErr
			close(req.errCh)

			return renewalErr
		}

		// If the identifier is not already tracked, this is a first
		// renewal request. In this case, add an entry into the heap
		// with the next renewal time.
		if err := c.heap.Push(req, next); err != nil {
			return fmt.Errorf("failed to push an entry to heap.  err: %v", err)
		}

		// Signal an update for the renewal loop to trigger a fresh
		// computation for the next best candidate for renewal.
		if c.running {
			select {
			case c.updateCh <- struct{}{}:
			default:
			}
		}
	}

	return nil
}
