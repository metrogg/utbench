func (d *AuthenticatedGossiper) processNetworkAnnouncement(
	nMsg *networkMsg) []networkMsg {

	isPremature := func(chanID lnwire.ShortChannelID, delta uint32) bool {
		// TODO(roasbeef) make height delta 6
		//  * or configurable
		bestHeight := atomic.LoadUint32(&d.bestHeight)
		return chanID.BlockHeight+delta > bestHeight
	}

	var announcements []networkMsg

	switch msg := nMsg.msg.(type) {

	// A new node announcement has arrived which either presents new
	// information about a node in one of the channels we know about, or a
	// updating previously advertised information.
	case *lnwire.NodeAnnouncement:
		timestamp := time.Unix(int64(msg.Timestamp), 0)

		// We'll quickly ask the router if it already has a
		// newer update for this node so we can skip validating
		// signatures if not required.
		if d.cfg.Router.IsStaleNode(msg.NodeID, timestamp) {
			nMsg.err <- nil
			return nil
		}

		if err := routing.ValidateNodeAnn(msg); err != nil {
			err := fmt.Errorf("unable to validate "+
				"node announcement: %v", err)
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		features := lnwire.NewFeatureVector(
			msg.Features, lnwire.GlobalFeatures,
		)
		node := &channeldb.LightningNode{
			HaveNodeAnnouncement: true,
			LastUpdate:           timestamp,
			Addresses:            msg.Addresses,
			PubKeyBytes:          msg.NodeID,
			Alias:                msg.Alias.String(),
			AuthSigBytes:         msg.Signature.ToSignatureBytes(),
			Features:             features,
			Color:                msg.RGBColor,
			ExtraOpaqueData:      msg.ExtraOpaqueData,
		}

		if err := d.cfg.Router.AddNode(node); err != nil {
			if routing.IsError(err, routing.ErrOutdated,
				routing.ErrIgnored) {

				log.Debug(err)
			} else {
				log.Error(err)
			}

			nMsg.err <- err
			return nil
		}

		// In order to ensure we don't leak unadvertised nodes, we'll
		// make a quick check to ensure this node intends to publicly
		// advertise itself to the network.
		isPublic, err := d.cfg.Router.IsPublicNode(node.PubKeyBytes)
		if err != nil {
			log.Errorf("Unable to determine if node %x is "+
				"advertised: %v", node.PubKeyBytes, err)
			nMsg.err <- err
			return nil
		}

		// If it does, we'll add their announcement to our batch so that
		// it can be broadcast to the rest of our peers.
		if isPublic {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    msg,
			})
		} else {
			log.Tracef("Skipping broadcasting node announcement "+
				"for %x due to being unadvertised", msg.NodeID)
		}

		nMsg.err <- nil
		// TODO(roasbeef): get rid of the above
		return announcements

	// A new channel announcement has arrived, this indicates the
	// *creation* of a new channel within the network. This only advertises
	// the existence of a channel and not yet the routing policies in
	// either direction of the channel.
	case *lnwire.ChannelAnnouncement:
		// We'll ignore any channel announcements that target any chain
		// other than the set of chains we know of.
		if !bytes.Equal(msg.ChainHash[:], d.cfg.ChainHash[:]) {
			err := fmt.Errorf("Ignoring ChannelAnnouncement from "+
				"chain=%v, gossiper on chain=%v", msg.ChainHash,
				d.cfg.ChainHash)
			log.Errorf(err.Error())

			d.rejectMtx.Lock()
			d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
			d.rejectMtx.Unlock()

			nMsg.err <- err
			return nil
		}

		// If the advertised inclusionary block is beyond our knowledge
		// of the chain tip, then we'll put the announcement in limbo
		// to be fully verified once we advance forward in the chain.
		if nMsg.isRemote && isPremature(msg.ShortChannelID, 0) {
			blockHeight := msg.ShortChannelID.BlockHeight
			log.Infof("Announcement for chan_id=(%v), is "+
				"premature: advertises height %v, only "+
				"height %v is known",
				msg.ShortChannelID.ToUint64(),
				msg.ShortChannelID.BlockHeight,
				atomic.LoadUint32(&d.bestHeight))

			d.Lock()
			d.prematureAnnouncements[blockHeight] = append(
				d.prematureAnnouncements[blockHeight],
				nMsg,
			)
			d.Unlock()
			return nil
		}

		// At this point, we'll now ask the router if this is a
		// zombie/known edge. If so we can skip all the processing
		// below.
		if d.cfg.Router.IsKnownEdge(msg.ShortChannelID) {
			nMsg.err <- nil
			return nil
		}

		// If this is a remote channel announcement, then we'll validate
		// all the signatures within the proof as it should be well
		// formed.
		var proof *channeldb.ChannelAuthProof
		if nMsg.isRemote {
			if err := routing.ValidateChannelAnn(msg); err != nil {
				err := fmt.Errorf("unable to validate "+
					"announcement: %v", err)
				d.rejectMtx.Lock()
				d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
				d.rejectMtx.Unlock()

				log.Error(err)
				nMsg.err <- err
				return nil
			}

			// If the proof checks out, then we'll save the proof
			// itself to the database so we can fetch it later when
			// gossiping with other nodes.
			proof = &channeldb.ChannelAuthProof{
				NodeSig1Bytes:    msg.NodeSig1.ToSignatureBytes(),
				NodeSig2Bytes:    msg.NodeSig2.ToSignatureBytes(),
				BitcoinSig1Bytes: msg.BitcoinSig1.ToSignatureBytes(),
				BitcoinSig2Bytes: msg.BitcoinSig2.ToSignatureBytes(),
			}
		}

		// With the proof validate (if necessary), we can now store it
		// within the database for our path finding and syncing needs.
		var featureBuf bytes.Buffer
		if err := msg.Features.Encode(&featureBuf); err != nil {
			log.Errorf("unable to encode features: %v", err)
			nMsg.err <- err
			return nil
		}

		edge := &channeldb.ChannelEdgeInfo{
			ChannelID:        msg.ShortChannelID.ToUint64(),
			ChainHash:        msg.ChainHash,
			NodeKey1Bytes:    msg.NodeID1,
			NodeKey2Bytes:    msg.NodeID2,
			BitcoinKey1Bytes: msg.BitcoinKey1,
			BitcoinKey2Bytes: msg.BitcoinKey2,
			AuthProof:        proof,
			Features:         featureBuf.Bytes(),
			ExtraOpaqueData:  msg.ExtraOpaqueData,
		}

		// If there were any optional message fields provided, we'll
		// include them in its serialized disk representation now.
		if nMsg.optionalMsgFields != nil {
			if nMsg.optionalMsgFields.capacity != nil {
				edge.Capacity = *nMsg.optionalMsgFields.capacity
			}
			if nMsg.optionalMsgFields.channelPoint != nil {
				edge.ChannelPoint = *nMsg.optionalMsgFields.channelPoint
			}
		}

		// We will add the edge to the channel router. If the nodes
		// present in this channel are not present in the database, a
		// partial node will be added to represent each node while we
		// wait for a node announcement.
		//
		// Before we add the edge to the database, we obtain
		// the mutex for this channel ID. We do this to ensure
		// no other goroutine has read the database and is now
		// making decisions based on this DB state, before it
		// writes to the DB.
		d.channelMtx.Lock(msg.ShortChannelID.ToUint64())
		defer d.channelMtx.Unlock(msg.ShortChannelID.ToUint64())
		if err := d.cfg.Router.AddEdge(edge); err != nil {
			// If the edge was rejected due to already being known,
			// then it may be that case that this new message has a
			// fresh channel proof, so we'll check.
			if routing.IsError(err, routing.ErrOutdated,
				routing.ErrIgnored) {

				// Attempt to process the rejected message to
				// see if we get any new announcements.
				anns, rErr := d.processRejectedEdge(msg, proof)
				if rErr != nil {
					d.rejectMtx.Lock()
					d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
					d.rejectMtx.Unlock()
					nMsg.err <- rErr
					return nil
				}

				// If while processing this rejected edge, we
				// realized there's a set of announcements we
				// could extract, then we'll return those
				// directly.
				if len(anns) != 0 {
					nMsg.err <- nil
					return anns
				}

				// Otherwise, this is just a regular rejected
				// edge.
				log.Debugf("Router rejected channel "+
					"edge: %v", err)
			} else {
				log.Tracef("Router rejected channel "+
					"edge: %v", err)
			}

			nMsg.err <- err
			return nil
		}

		// If we earlier received any ChannelUpdates for this channel,
		// we can now process them, as the channel is added to the
		// graph.
		shortChanID := msg.ShortChannelID.ToUint64()
		var channelUpdates []*networkMsg

		d.pChanUpdMtx.Lock()
		for _, cu := range d.prematureChannelUpdates[shortChanID] {
			channelUpdates = append(channelUpdates, cu)
		}

		// Now delete the premature ChannelUpdates, since we added them
		// all to the queue of network messages.
		delete(d.prematureChannelUpdates, shortChanID)
		d.pChanUpdMtx.Unlock()

		// Launch a new goroutine to handle each ChannelUpdate, this to
		// ensure we don't block here, as we can handle only one
		// announcement at a time.
		for _, cu := range channelUpdates {
			d.wg.Add(1)
			go func(nMsg *networkMsg) {
				defer d.wg.Done()

				switch msg := nMsg.msg.(type) {

				// Reprocess the message, making sure we return
				// an error to the original caller in case the
				// gossiper shuts down.
				case *lnwire.ChannelUpdate:
					log.Debugf("Reprocessing"+
						" ChannelUpdate for "+
						"shortChanID=%v",
						msg.ShortChannelID.ToUint64())

					select {
					case d.networkMsgs <- nMsg:
					case <-d.quit:
						nMsg.err <- ErrGossiperShuttingDown
					}

				// We don't expect any other message type than
				// ChannelUpdate to be in this map.
				default:
					log.Errorf("Unsupported message type "+
						"found among ChannelUpdates: "+
						"%T", msg)
				}
			}(cu)
		}

		// Channel announcement was successfully proceeded and know it
		// might be broadcast to other connected nodes if it was
		// announcement with proof (remote).
		if proof != nil {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    msg,
			})
		}

		nMsg.err <- nil
		return announcements

	// A new authenticated channel edge update has arrived. This indicates
	// that the directional information for an already known channel has
	// been updated.
	case *lnwire.ChannelUpdate:
		// We'll ignore any channel announcements that target any chain
		// other than the set of chains we know of.
		if !bytes.Equal(msg.ChainHash[:], d.cfg.ChainHash[:]) {
			err := fmt.Errorf("Ignoring ChannelUpdate from "+
				"chain=%v, gossiper on chain=%v", msg.ChainHash,
				d.cfg.ChainHash)
			log.Errorf(err.Error())

			d.rejectMtx.Lock()
			d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
			d.rejectMtx.Unlock()

			nMsg.err <- err
			return nil
		}

		blockHeight := msg.ShortChannelID.BlockHeight
		shortChanID := msg.ShortChannelID.ToUint64()

		// If the advertised inclusionary block is beyond our knowledge
		// of the chain tip, then we'll put the announcement in limbo
		// to be fully verified once we advance forward in the chain.
		if nMsg.isRemote && isPremature(msg.ShortChannelID, 0) {
			log.Infof("Update announcement for "+
				"short_chan_id(%v), is premature: advertises "+
				"height %v, only height %v is known",
				shortChanID, blockHeight,
				atomic.LoadUint32(&d.bestHeight))

			d.Lock()
			d.prematureAnnouncements[blockHeight] = append(
				d.prematureAnnouncements[blockHeight],
				nMsg,
			)
			d.Unlock()
			return nil
		}

		// Before we perform any of the expensive checks below, we'll
		// check whether this update is stale or is for a zombie
		// channel in order to quickly reject it.
		timestamp := time.Unix(int64(msg.Timestamp), 0)
		if d.cfg.Router.IsStaleEdgePolicy(
			msg.ShortChannelID, timestamp, msg.ChannelFlags,
		) {
			nMsg.err <- nil
			return nil
		}

		// Get the node pub key as far as we don't have it in channel
		// update announcement message. We'll need this to properly
		// verify message signature.
		//
		// We make sure to obtain the mutex for this channel ID
		// before we access the database. This ensures the state
		// we read from the database has not changed between this
		// point and when we call UpdateEdge() later.
		d.channelMtx.Lock(msg.ShortChannelID.ToUint64())
		defer d.channelMtx.Unlock(msg.ShortChannelID.ToUint64())
		chanInfo, _, _, err := d.cfg.Router.GetChannelByID(msg.ShortChannelID)
		switch err {
		// No error, break.
		case nil:
			break

		case channeldb.ErrZombieEdge:
			// Since we've deemed the update as not stale above,
			// before marking it live, we'll make sure it has been
			// signed by the correct party. The least-significant
			// bit in the flag on the channel update tells us which
			// edge is being updated.
			var pubKey *btcec.PublicKey
			switch {
			case msg.ChannelFlags&lnwire.ChanUpdateDirection == 0:
				pubKey, _ = chanInfo.NodeKey1()
			case msg.ChannelFlags&lnwire.ChanUpdateDirection == 1:
				pubKey, _ = chanInfo.NodeKey2()
			}

			err := routing.VerifyChannelUpdateSignature(msg, pubKey)
			if err != nil {
				err := fmt.Errorf("unable to verify channel "+
					"update signature: %v", err)
				log.Error(err)
				nMsg.err <- err
				return nil
			}

			// With the signature valid, we'll proceed to mark the
			// edge as live and wait for the channel announcement to
			// come through again.
			err = d.cfg.Router.MarkEdgeLive(msg.ShortChannelID)
			if err != nil {
				err := fmt.Errorf("unable to remove edge with "+
					"chan_id=%v from zombie index: %v",
					msg.ShortChannelID, err)
				log.Error(err)
				nMsg.err <- err
				return nil
			}

			log.Debugf("Removed edge with chan_id=%v from zombie "+
				"index", msg.ShortChannelID)

			// We'll fallthrough to ensure we stash the update until
			// we receive its corresponding ChannelAnnouncement.
			// This is needed to ensure the edge exists in the graph
			// before applying the update.
			fallthrough
		case channeldb.ErrGraphNotFound:
			fallthrough
		case channeldb.ErrGraphNoEdgesFound:
			fallthrough
		case channeldb.ErrEdgeNotFound:
			// If the edge corresponding to this ChannelUpdate was
			// not found in the graph, this might be a channel in
			// the process of being opened, and we haven't processed
			// our own ChannelAnnouncement yet, hence it is not
			// found in the graph. This usually gets resolved after
			// the channel proofs are exchanged and the channel is
			// broadcasted to the rest of the network, but in case
			// this is a private channel this won't ever happen.
			// This can also happen in the case of a zombie channel
			// with a fresh update for which we don't have a
			// ChannelAnnouncement for since we reject them. Because
			// of this, we temporarily add it to a map, and
			// reprocess it after our own ChannelAnnouncement has
			// been processed.
			d.pChanUpdMtx.Lock()
			d.prematureChannelUpdates[shortChanID] = append(
				d.prematureChannelUpdates[shortChanID], nMsg,
			)
			d.pChanUpdMtx.Unlock()

			log.Debugf("Got ChannelUpdate for edge not found in "+
				"graph(shortChanID=%v), saving for "+
				"reprocessing later", shortChanID)

			// NOTE: We don't return anything on the error channel
			// for this message, as we expect that will be done when
			// this ChannelUpdate is later reprocessed.
			return nil

		default:
			err := fmt.Errorf("unable to validate channel update "+
				"short_chan_id=%v: %v", shortChanID, err)
			log.Error(err)
			nMsg.err <- err

			d.rejectMtx.Lock()
			d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
			d.rejectMtx.Unlock()
			return nil
		}

		// The least-significant bit in the flag on the channel update
		// announcement tells us "which" side of the channels directed
		// edge is being updated.
		var pubKey *btcec.PublicKey
		switch {
		case msg.ChannelFlags&lnwire.ChanUpdateDirection == 0:
			pubKey, _ = chanInfo.NodeKey1()
		case msg.ChannelFlags&lnwire.ChanUpdateDirection == 1:
			pubKey, _ = chanInfo.NodeKey2()
		}

		// Validate the channel announcement with the expected public key and
		// channel capacity. In the case of an invalid channel update, we'll
		// return an error to the caller and exit early.
		err = routing.ValidateChannelUpdateAnn(pubKey, chanInfo.Capacity, msg)
		if err != nil {
			rErr := fmt.Errorf("unable to validate channel "+
				"update announcement for short_chan_id=%v: %v",
				spew.Sdump(msg.ShortChannelID), err)

			log.Error(rErr)
			nMsg.err <- rErr
			return nil
		}

		update := &channeldb.ChannelEdgePolicy{
			SigBytes:                  msg.Signature.ToSignatureBytes(),
			ChannelID:                 shortChanID,
			LastUpdate:                timestamp,
			MessageFlags:              msg.MessageFlags,
			ChannelFlags:              msg.ChannelFlags,
			TimeLockDelta:             msg.TimeLockDelta,
			MinHTLC:                   msg.HtlcMinimumMsat,
			MaxHTLC:                   msg.HtlcMaximumMsat,
			FeeBaseMSat:               lnwire.MilliSatoshi(msg.BaseFee),
			FeeProportionalMillionths: lnwire.MilliSatoshi(msg.FeeRate),
			ExtraOpaqueData:           msg.ExtraOpaqueData,
		}

		if err := d.cfg.Router.UpdateEdge(update); err != nil {
			if routing.IsError(err, routing.ErrOutdated,
				routing.ErrIgnored) {
				log.Debug(err)
			} else {
				d.rejectMtx.Lock()
				d.recentRejects[msg.ShortChannelID.ToUint64()] = struct{}{}
				d.rejectMtx.Unlock()
				log.Error(err)
			}

			nMsg.err <- err
			return nil
		}

		// If this is a local ChannelUpdate without an AuthProof, it
		// means it is an update to a channel that is not (yet)
		// supposed to be announced to the greater network. However,
		// our channel counter party will need to be given the update,
		// so we'll try sending the update directly to the remote peer.
		if !nMsg.isRemote && chanInfo.AuthProof == nil {
			// Get our peer's public key.
			remotePubKey := remotePubFromChanInfo(
				chanInfo, msg.ChannelFlags,
			)

			// Now, we'll attempt to send the channel update message
			// reliably to the remote peer in the background, so
			// that we don't block if the peer happens to be offline
			// at the moment.
			err := d.reliableSender.sendMessage(msg, remotePubKey)
			if err != nil {
				err := fmt.Errorf("unable to reliably send %v "+
					"for channel=%v to peer=%x: %v",
					msg.MsgType(), msg.ShortChannelID,
					remotePubKey, err)
				nMsg.err <- err
				return nil
			}
		}

		// Channel update announcement was successfully processed and
		// now it can be broadcast to the rest of the network. However,
		// we'll only broadcast the channel update announcement if it
		// has an attached authentication proof.
		if chanInfo.AuthProof != nil {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    msg,
			})
		}

		nMsg.err <- nil
		return announcements

	// A new signature announcement has been received. This indicates
	// willingness of nodes involved in the funding of a channel to
	// announce this new channel to the rest of the world.
	case *lnwire.AnnounceSignatures:
		needBlockHeight := msg.ShortChannelID.BlockHeight +
			d.cfg.ProofMatureDelta
		shortChanID := msg.ShortChannelID.ToUint64()

		prefix := "local"
		if nMsg.isRemote {
			prefix = "remote"
		}

		log.Infof("Received new %v channel announcement: %v", prefix,
			spew.Sdump(msg))

		// By the specification, channel announcement proofs should be
		// sent after some number of confirmations after channel was
		// registered in bitcoin blockchain. Therefore, we check if the
		// proof is premature.  If so we'll halt processing until the
		// expected announcement height.  This allows us to be tolerant
		// to other clients if this constraint was changed.
		if isPremature(msg.ShortChannelID, d.cfg.ProofMatureDelta) {
			d.Lock()
			d.prematureAnnouncements[needBlockHeight] = append(
				d.prematureAnnouncements[needBlockHeight],
				nMsg,
			)
			d.Unlock()
			log.Infof("Premature proof announcement, "+
				"current block height lower than needed: %v <"+
				" %v, add announcement to reprocessing batch",
				atomic.LoadUint32(&d.bestHeight), needBlockHeight)
			return nil
		}

		// Ensure that we know of a channel with the target channel ID
		// before proceeding further.
		//
		// We must acquire the mutex for this channel ID before getting
		// the channel from the database, to ensure what we read does
		// not change before we call AddProof() later.
		d.channelMtx.Lock(msg.ShortChannelID.ToUint64())
		defer d.channelMtx.Unlock(msg.ShortChannelID.ToUint64())

		chanInfo, e1, e2, err := d.cfg.Router.GetChannelByID(
			msg.ShortChannelID)
		if err != nil {
			// TODO(andrew.shvv) this is dangerous because remote
			// node might rewrite the waiting proof.
			proof := channeldb.NewWaitingProof(nMsg.isRemote, msg)
			err := d.cfg.WaitingProofStore.Add(proof)
			if err != nil {
				err := fmt.Errorf("unable to store "+
					"the proof for short_chan_id=%v: %v",
					shortChanID, err)
				log.Error(err)
				nMsg.err <- err
				return nil
			}

			log.Infof("Orphan %v proof announcement with "+
				"short_chan_id=%v, adding"+
				"to waiting batch", prefix, shortChanID)
			nMsg.err <- nil
			return nil
		}

		nodeID := nMsg.source.SerializeCompressed()
		isFirstNode := bytes.Equal(nodeID, chanInfo.NodeKey1Bytes[:])
		isSecondNode := bytes.Equal(nodeID, chanInfo.NodeKey2Bytes[:])

		// Ensure that channel that was retrieved belongs to the peer
		// which sent the proof announcement.
		if !(isFirstNode || isSecondNode) {
			err := fmt.Errorf("channel that was received not "+
				"belongs to the peer which sent the proof, "+
				"short_chan_id=%v", shortChanID)
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		// If proof was sent by a local sub-system, then we'll
		// send the announcement signature to the remote node
		// so they can also reconstruct the full channel
		// announcement.
		if !nMsg.isRemote {
			var remotePubKey [33]byte
			if isFirstNode {
				remotePubKey = chanInfo.NodeKey2Bytes
			} else {
				remotePubKey = chanInfo.NodeKey1Bytes
			}
			// Since the remote peer might not be online
			// we'll call a method that will attempt to
			// deliver the proof when it comes online.
			err := d.reliableSender.sendMessage(msg, remotePubKey)
			if err != nil {
				err := fmt.Errorf("unable to reliably send %v "+
					"for channel=%v to peer=%x: %v",
					msg.MsgType(), msg.ShortChannelID,
					remotePubKey, err)
				nMsg.err <- err
				return nil
			}
		}

		// Check if we already have the full proof for this channel.
		if chanInfo.AuthProof != nil {
			// If we already have the fully assembled proof, then
			// the peer sending us their proof has probably not
			// received our local proof yet. So be kind and send
			// them the full proof.
			if nMsg.isRemote {
				peerID := nMsg.source.SerializeCompressed()
				log.Debugf("Got AnnounceSignatures for " +
					"channel with full proof.")

				d.wg.Add(1)
				go func() {
					defer d.wg.Done()
					log.Debugf("Received half proof for "+
						"channel %v with existing "+
						"full proof. Sending full "+
						"proof to peer=%x",
						msg.ChannelID,
						peerID)

					chanAnn, _, _, err := CreateChanAnnouncement(
						chanInfo.AuthProof, chanInfo,
						e1, e2,
					)
					if err != nil {
						log.Errorf("unable to gen "+
							"ann: %v", err)
						return
					}
					err = nMsg.peer.SendMessage(
						false, chanAnn,
					)
					if err != nil {
						log.Errorf("Failed sending "+
							"full proof to "+
							"peer=%x: %v",
							peerID, err)
						return
					}
					log.Debugf("Full proof sent to peer=%x"+
						" for chanID=%v", peerID,
						msg.ChannelID)
				}()
			}

			log.Debugf("Already have proof for channel "+
				"with chanID=%v", msg.ChannelID)
			nMsg.err <- nil
			return nil
		}

		// Check that we received the opposite proof. If so, then we're
		// now able to construct the full proof, and create the channel
		// announcement. If we didn't receive the opposite half of the
		// proof than we should store it this one, and wait for
		// opposite to be received.
		proof := channeldb.NewWaitingProof(nMsg.isRemote, msg)
		oppositeProof, err := d.cfg.WaitingProofStore.Get(
			proof.OppositeKey(),
		)
		if err != nil && err != channeldb.ErrWaitingProofNotFound {
			err := fmt.Errorf("unable to get "+
				"the opposite proof for short_chan_id=%v: %v",
				shortChanID, err)
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		if err == channeldb.ErrWaitingProofNotFound {
			err := d.cfg.WaitingProofStore.Add(proof)
			if err != nil {
				err := fmt.Errorf("unable to store "+
					"the proof for short_chan_id=%v: %v",
					shortChanID, err)
				log.Error(err)
				nMsg.err <- err
				return nil
			}

			log.Infof("1/2 of channel ann proof received for "+
				"short_chan_id=%v, waiting for other half",
				shortChanID)

			nMsg.err <- nil
			return nil
		}

		// We now have both halves of the channel announcement proof,
		// then we'll reconstruct the initial announcement so we can
		// validate it shortly below.
		var dbProof channeldb.ChannelAuthProof
		if isFirstNode {
			dbProof.NodeSig1Bytes = msg.NodeSignature.ToSignatureBytes()
			dbProof.NodeSig2Bytes = oppositeProof.NodeSignature.ToSignatureBytes()
			dbProof.BitcoinSig1Bytes = msg.BitcoinSignature.ToSignatureBytes()
			dbProof.BitcoinSig2Bytes = oppositeProof.BitcoinSignature.ToSignatureBytes()
		} else {
			dbProof.NodeSig1Bytes = oppositeProof.NodeSignature.ToSignatureBytes()
			dbProof.NodeSig2Bytes = msg.NodeSignature.ToSignatureBytes()
			dbProof.BitcoinSig1Bytes = oppositeProof.BitcoinSignature.ToSignatureBytes()
			dbProof.BitcoinSig2Bytes = msg.BitcoinSignature.ToSignatureBytes()
		}
		chanAnn, e1Ann, e2Ann, err := CreateChanAnnouncement(
			&dbProof, chanInfo, e1, e2,
		)
		if err != nil {
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		// With all the necessary components assembled validate the
		// full channel announcement proof.
		if err := routing.ValidateChannelAnn(chanAnn); err != nil {
			err := fmt.Errorf("channel  announcement proof "+
				"for short_chan_id=%v isn't valid: %v",
				shortChanID, err)

			log.Error(err)
			nMsg.err <- err
			return nil
		}

		// If the channel was returned by the router it means that
		// existence of funding point and inclusion of nodes bitcoin
		// keys in it already checked by the router. In this stage we
		// should check that node keys are attest to the bitcoin keys
		// by validating the signatures of announcement.  If proof is
		// valid then we'll populate the channel edge with it, so we
		// can announce it on peer connect.
		err = d.cfg.Router.AddProof(msg.ShortChannelID, &dbProof)
		if err != nil {
			err := fmt.Errorf("unable add proof to the "+
				"channel chanID=%v: %v", msg.ChannelID, err)
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		err = d.cfg.WaitingProofStore.Remove(proof.OppositeKey())
		if err != nil {
			err := fmt.Errorf("unable remove opposite proof "+
				"for the channel with chanID=%v: %v",
				msg.ChannelID, err)
			log.Error(err)
			nMsg.err <- err
			return nil
		}

		// Proof was successfully created and now can announce the
		// channel to the remain network.
		log.Infof("Fully valid channel proof for short_chan_id=%v "+
			"constructed, adding to next ann batch",
			shortChanID)

		// Assemble the necessary announcements to add to the next
		// broadcasting batch.
		announcements = append(announcements, networkMsg{
			peer:   nMsg.peer,
			source: nMsg.source,
			msg:    chanAnn,
		})
		if e1Ann != nil {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    e1Ann,
			})
		}
		if e2Ann != nil {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    e2Ann,
			})
		}

		// We'll also send along the node announcements for each channel
		// participant if we know of them.
		node1Ann, err := d.fetchNodeAnn(chanInfo.NodeKey1Bytes)
		if err != nil {
			log.Debugf("Unable to fetch node announcement for "+
				"%x: %v", chanInfo.NodeKey1Bytes, err)
		} else {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    node1Ann,
			})
		}
		node2Ann, err := d.fetchNodeAnn(chanInfo.NodeKey2Bytes)
		if err != nil {
			log.Debugf("Unable to fetch node announcement for "+
				"%x: %v", chanInfo.NodeKey2Bytes, err)
		} else {
			announcements = append(announcements, networkMsg{
				peer:   nMsg.peer,
				source: nMsg.source,
				msg:    node2Ann,
			})
		}

		nMsg.err <- nil
		return announcements

	default:
		nMsg.err <- errors.New("wrong type of the announcement")
		return nil
	}
}
