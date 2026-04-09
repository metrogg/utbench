func (km *KeyManagerStandard) Rekey(ctx context.Context, md *RootMetadata, promptPaper bool) (
	mdChanged bool, cryptKey *kbfscrypto.TLFCryptKey, err error) {
	km.log.CDebugf(ctx, "Rekey %s (prompt for paper key: %t)",
		md.TlfID(), promptPaper)
	defer func() { km.deferLog.CDebugf(ctx, "Rekey %s done: %+v", md.TlfID(), err) }()

	if md.TypeForKeying() == tlf.TeamKeying {
		return false, nil, errors.New(
			"Rekeying is not done by KBFS for team TLFs")
	}

	currKeyGen := md.LatestKeyGeneration()
	if (md.TypeForKeying() == tlf.PublicKeying) != (currKeyGen == kbfsmd.PublicKeyGen) {
		return false, nil, errors.Errorf(
			"ID %v has keying type=%s but currKeyGen is %d (isPublic=%t)",
			md.TlfID(), md.TypeForKeying(), currKeyGen,
			currKeyGen == kbfsmd.PublicKeyGen)
	}

	if promptPaper && md.TypeForKeying() != tlf.PrivateKeying {
		return false, nil, errors.Errorf(
			"promptPaper set for non-private TLF %v", md.TlfID())
	}

	handle := md.GetTlfHandle()

	session, err := km.config.KBPKI().GetCurrentSession(ctx)
	if err != nil {
		return false, nil, err
	}

	idGetter := tlfhandle.ConstIDGetter{ID: md.TlfID()}
	resolvedHandle, err := handle.ResolveAgain(
		ctx, km.config.KBPKI(), idGetter, km.config)
	if err != nil {
		return false, nil, err
	}

	isWriter := resolvedHandle.IsWriter(session.UID)
	if md.TypeForKeying() == tlf.PrivateKeying && !isWriter {
		// If I was already a reader, there's nothing more to do
		if handle.IsReader(session.UID) {
			resolvedHandle = handle
			km.log.CDebugf(ctx, "Local user is not a writer, and was "+
				"already a reader; reverting back to the original handle")
		} else {
			// Only allow yourself to change
			resolvedHandle, err = handle.ResolveAgainForUser(
				ctx, km.config.KBPKI(), idGetter, km.config, session.UID)
			if err != nil {
				return false, nil, err
			}
		}
	}

	eq, err := handle.Equals(km.config.Codec(), *resolvedHandle)
	if err != nil {
		return false, nil, err
	}
	handleChanged := !eq
	if handleChanged {
		km.log.CDebugf(ctx, "handle for %s resolved to %s",
			handle.GetCanonicalPath(),
			resolvedHandle.GetCanonicalPath())

		// Check with the server to see if the handle became a conflict.
		latestHandle, err := km.config.MDOps().GetLatestHandleForTLF(ctx, md.TlfID())
		if err != nil {
			return false, nil, err
		}
		if latestHandle.ConflictInfo != nil {
			km.log.CDebugf(ctx, "handle for %s is conflicted",
				handle.GetCanonicalPath())
		}
		resolvedHandle, err = resolvedHandle.WithUpdatedConflictInfo(
			km.config.Codec(), latestHandle.ConflictInfo)
		if err != nil {
			return false, nil, err
		}
	}

	// For a public or team TLF there's no rekeying to be done, but we
	// should still update the writer list.
	if md.TypeForKeying() != tlf.PrivateKeying {
		if !handleChanged {
			km.log.CDebugf(ctx,
				"Skipping rekeying %s (%s): handle hasn't changed",
				md.TlfID(), md.TypeForKeying())
			return false, nil, nil
		}
		return true, nil, md.updateFromTlfHandle(resolvedHandle)
	}

	// Decide whether we have a new device and/or a revoked device, or neither.
	// Look up all the device public keys for all writers and readers first.

	incKeyGen := currKeyGen < kbfsmd.FirstValidKeyGen

	if !isWriter && incKeyGen {
		// Readers cannot create the first key generation
		return false, nil, tlfhandle.NewReadAccessError(
			resolvedHandle, session.Name, resolvedHandle.GetCanonicalPath())
	}

	offline := km.config.OfflineAvailabilityForID(md.TlfID())

	// All writer keys in the desired keyset
	updatedWriterKeys, err := km.generateKeyMapForUsers(
		ctx, resolvedHandle.ResolvedWriters(), offline)
	if err != nil {
		return false, nil, err
	}
	// All reader keys in the desired keyset
	updatedReaderKeys, err := km.generateKeyMapForUsers(
		ctx, resolvedHandle.ResolvedReaders(), offline)
	if err != nil {
		return false, nil, err
	}

	addNewReaderDevice := false
	addNewWriterDevice := false
	var newReaderUsers map[keybase1.UID]bool
	var newWriterUsers map[keybase1.UID]bool
	var readersToPromote map[keybase1.UID]bool

	// Figure out if we need to add or remove any keys.
	// If we're already incrementing the key generation then we don't need to
	// figure out the key delta.
	addNewReaderDeviceForSelf := false
	if !incKeyGen {
		// See if there is at least one new device in relation to the
		// current key bundle
		writers, readers, err := md.getUserDevicePublicKeys()
		if err != nil {
			return false, nil, err
		}

		newWriterUsers = km.usersWithNewDevices(
			ctx, md.TlfID(), writers, updatedWriterKeys)
		newReaderUsers = km.usersWithNewDevices(
			ctx, md.TlfID(), readers, updatedReaderKeys)
		addNewWriterDevice = len(newWriterUsers) > 0
		addNewReaderDevice = len(newReaderUsers) > 0

		wRemoved := km.usersWithRemovedDevices(
			ctx, md.TlfID(), writers, updatedWriterKeys)
		rRemoved := km.usersWithRemovedDevices(
			ctx, md.TlfID(), readers, updatedReaderKeys)

		readersToPromote = make(map[keybase1.UID]bool, len(rRemoved))

		// Before we add the removed devices, check if we are adding a
		// new reader device for ourselves.
		_, addNewReaderDeviceForSelf = newReaderUsers[session.UID]

		for u := range rRemoved {
			// FIXME (potential): this could cause a reader to attempt to rekey
			// in the case of a revocation for the currently logged-in user. I
			// _think_ incKeyGen above protects against this, but I'm not
			// confident.
			newReaderUsers[u] = true
			// Track which readers have been promoted. This must happen before
			// the following line adds all the removed writers to the writer
			// set
			if newWriterUsers[u] {
				readersToPromote[u] = true
			}
		}
		for u := range wRemoved {
			newWriterUsers[u] = true
		}

		incKeyGen = len(wRemoved) > 0 || (len(rRemoved) > len(readersToPromote))

		if err := km.identifyUIDSets(ctx, md.TlfID(), newWriterUsers, newReaderUsers); err != nil {
			return false, nil, err
		}
	}

	if !addNewReaderDevice && !addNewWriterDevice && !incKeyGen &&
		!handleChanged {
		km.log.CDebugf(ctx,
			"Skipping rekeying %s (private): no new or removed devices, no new keygen, and handle hasn't changed",
			md.TlfID())
		return false, nil, nil
	}

	if !isWriter {
		// This shouldn't happen; see the code above where we
		// either use the original handle or resolve only the
		// current UID.
		if len(readersToPromote) != 0 {
			return false, nil, errors.New(
				"promoted readers unexpectedly non-empty")
		}

		if _, userHasNewKeys := newReaderUsers[session.UID]; userHasNewKeys {
			// Only rekey the logged-in reader.
			updatedWriterKeys = nil
			updatedReaderKeys = kbfsmd.UserDevicePublicKeys{
				session.UID: updatedReaderKeys[session.UID],
			}
			delete(newReaderUsers, session.UID)
		} else {
			// No new reader device for our user, so the reader can't do
			// anything
			return false, nil, RekeyIncompleteError{}
		}
	}

	// If promotedReader is non-empty, then isWriter is true (see
	// check above).
	err = md.promoteReaders(readersToPromote)
	if err != nil {
		return false, nil, err
	}

	// Generate ephemeral keys to be used by addNewDevice,
	// incKeygen, or both. ePrivKey will be discarded at the end
	// of the function.
	ePubKey, ePrivKey, err :=
		km.config.Crypto().MakeRandomTLFEphemeralKeys()

	// Note: For MDv3, if incKeyGen is true, then all the
	// manipulations below aren't needed, since they'll just be
	// replaced by the new key generation. However, do them
	// anyway, as they may have some side effects, e.g. removing
	// server key halves.

	// If there's at least one new device, add that device to every key bundle.
	if addNewReaderDevice || addNewWriterDevice {
		start, end := md.KeyGenerationsToUpdate()
		if start >= end {
			return false, nil, errors.New(
				"Unexpected empty range for key generations to update")
		}
		cryptKeys := make([]kbfscrypto.TLFCryptKey, end-start)
		flags := getTLFCryptKeyAnyDevice
		if promptPaper {
			flags |= getTLFCryptKeyPromptPaper
		}
		for keyGen := start; keyGen < end; keyGen++ {
			currTlfCryptKey, err := km.getTLFCryptKey(
				ctx, md.ReadOnly(), keyGen, flags)
			if err != nil {
				return false, nil, err
			}
			cryptKeys[keyGen-start] = currTlfCryptKey
		}
		err = km.updateKeyBundles(ctx, md, updatedWriterKeys,
			updatedReaderKeys, ePubKey, ePrivKey, cryptKeys)
		if err != nil {
			return false, nil, err
		}
	}

	// Make sure the private MD is decrypted if it wasn't already.  We
	// have to do this here, before adding a new key generation, since
	// decryptMDPrivateData assumes that the MD is always encrypted
	// using the latest key gen.
	if !md.IsReadable() && len(md.GetSerializedPrivateMetadata()) > 0 {
		pmd, err := decryptMDPrivateData(
			ctx, km.config.Codec(), km.config.Crypto(),
			km.config.BlockCache(), km.config.BlockOps(), km, km.config.KBPKI(),
			km.config, km.config.Mode(), session.UID,
			md.GetSerializedPrivateMetadata(), md, md, km.log)
		if err != nil {
			return false, nil, err
		}
		md.data = pmd
	}

	defer func() {
		// On our way back out, update the md with the
		// resolved handle if at least part of a rekey was
		// performed.  Also, if we return true for mdChanged
		// with a nil error or RekeyIncompleteError{}, we must
		// call md.finalizeRekey() first.

		_, isRekeyIncomplete := err.(RekeyIncompleteError)
		if err == nil || isRekeyIncomplete {
			updateErr := md.updateFromTlfHandle(resolvedHandle)
			if updateErr != nil {
				mdChanged = false
				cryptKey = nil
				err = updateErr
			}

			if mdChanged {
				finalizeErr := md.finalizeRekey(
					km.config.Codec())
				if finalizeErr != nil {
					mdChanged = false
					cryptKey = nil
					err = finalizeErr
				}
			}
		}
	}()

	if !isWriter {
		if len(newReaderUsers) > 0 || addNewWriterDevice || incKeyGen {
			// If we're a reader but we haven't completed all the work, return
			// RekeyIncompleteError.
			return addNewReaderDeviceForSelf, nil, RekeyIncompleteError{}
		}
		// Otherwise, there's nothing left to do!
		return true, nil, nil
	} else if !incKeyGen {
		// we're done!
		return true, nil, nil
	}

	// Send rekey start notification once we're sure that this device
	// can perform the rekey.  Don't send this for the first key
	// generation though, since that's not really a "re"key.
	//
	// TODO: Shouldn't this happen earlier?
	if currKeyGen >= kbfsmd.FirstValidKeyGen {
		km.config.Reporter().Notify(ctx, rekeyNotification(
			ctx, km.config, resolvedHandle, false))
	}

	// Delete server-side key halves for any revoked devices, if
	// there are any previous key generations. Do this before
	// adding a new key generation, as MDv3 only keeps track of
	// the latest key generation.
	//
	// TODO: Add test coverage for this.
	if currKeyGen >= kbfsmd.FirstValidKeyGen {
		allRemovalInfo, err := md.revokeRemovedDevices(
			updatedWriterKeys, updatedReaderKeys)
		if err != nil {
			return false, nil, err
		}

		if len(allRemovalInfo) == 0 {
			return false, nil, errors.New(
				"Didn't revoke any devices, but indicated incrementing the key generation")
		}

		kops := km.config.KeyOps()
		for uid, userRemovalInfo := range allRemovalInfo {
			if userRemovalInfo.UserRemoved {
				km.log.CInfof(ctx, "Rekey %s: removed user %s entirely",
					md.TlfID(), uid)
			}
			for key, serverHalfIDs := range userRemovalInfo.DeviceServerHalfIDs {
				km.log.CInfof(ctx, "Rekey %s: removing %d server key halves "+
					"for device %s of user %s", md.TlfID(),
					len(serverHalfIDs), key, uid)
				for _, serverHalfID := range serverHalfIDs {
					err := kops.DeleteTLFCryptKeyServerHalf(
						ctx, uid, key, serverHalfID)
					if err != nil {
						return false, nil, err
					}
				}
			}
		}
	}

	pubKey, privKey, tlfCryptKey, err :=
		km.config.Crypto().MakeRandomTLFKeys()
	if err != nil {
		return false, nil, err
	}

	// Get the current TLF crypt key if needed. It's
	// symmetrically encrypted and appended to a list for MDv3
	// metadata.
	var currTLFCryptKey kbfscrypto.TLFCryptKey
	if md.StoresHistoricTLFCryptKeys() && currKeyGen >= kbfsmd.FirstValidKeyGen {
		flags := getTLFCryptKeyAnyDevice
		if promptPaper {
			flags |= getTLFCryptKeyPromptPaper
		}
		var err error
		currTLFCryptKey, err = km.getTLFCryptKey(
			ctx, md.ReadOnly(), currKeyGen, flags)
		if err != nil {
			return false, nil, err
		}
	}
	serverHalves, err := md.AddKeyGeneration(km.config.Codec(),
		updatedWriterKeys, updatedReaderKeys,
		ePubKey, ePrivKey, pubKey, privKey,
		currTLFCryptKey, tlfCryptKey)
	if err != nil {
		return false, nil, err
	}

	err = km.config.KeyOps().PutTLFCryptKeyServerHalves(ctx, serverHalves)
	if err != nil {
		return false, nil, err
	}

	return true, &tlfCryptKey, nil
}
