func newServer(listenAddrs []net.Addr, chanDB *channeldb.DB, cc *chainControl,
	privKey *btcec.PrivateKey,
	chansToRestore walletunlocker.ChannelsToRecover) (*server, error) {

	var err error

	listeners := make([]net.Listener, len(listenAddrs))
	for i, listenAddr := range listenAddrs {
		// Note: though brontide.NewListener uses ResolveTCPAddr, it
		// doesn't need to call the general lndResolveTCP function
		// since we are resolving a local address.
		listeners[i], err = brontide.NewListener(
			privKey, listenAddr.String(),
		)
		if err != nil {
			return nil, err
		}
	}

	globalFeatures := lnwire.NewRawFeatureVector()

	var serializedPubKey [33]byte
	copy(serializedPubKey[:], privKey.PubKey().SerializeCompressed())

	// Initialize the sphinx router, placing it's persistent replay log in
	// the same directory as the channel graph database.
	graphDir := chanDB.Path()
	sharedSecretPath := filepath.Join(graphDir, "sphinxreplay.db")
	replayLog := htlcswitch.NewDecayedLog(sharedSecretPath, cc.chainNotifier)
	sphinxRouter := sphinx.NewRouter(privKey, activeNetParams.Params, replayLog)

	writeBufferPool := pool.NewWriteBuffer(
		pool.DefaultWriteBufferGCInterval,
		pool.DefaultWriteBufferExpiryInterval,
	)

	writePool := pool.NewWrite(
		writeBufferPool, cfg.Workers.Write, pool.DefaultWorkerTimeout,
	)

	readBufferPool := pool.NewReadBuffer(
		pool.DefaultReadBufferGCInterval,
		pool.DefaultReadBufferExpiryInterval,
	)

	readPool := pool.NewRead(
		readBufferPool, cfg.Workers.Read, pool.DefaultWorkerTimeout,
	)

	decodeFinalCltvExpiry := func(payReq string) (uint32, error) {
		invoice, err := zpay32.Decode(payReq, activeNetParams.Params)
		if err != nil {
			return 0, err
		}
		return uint32(invoice.MinFinalCLTVExpiry()), nil
	}

	s := &server{
		chanDB:         chanDB,
		cc:             cc,
		sigPool:        lnwallet.NewSigPool(cfg.Workers.Sig, cc.signer),
		writePool:      writePool,
		readPool:       readPool,
		chansToRestore: chansToRestore,

		invoices: invoices.NewRegistry(chanDB, decodeFinalCltvExpiry),

		channelNotifier: channelnotifier.New(chanDB),

		identityPriv: privKey,
		nodeSigner:   netann.NewNodeSigner(privKey),

		listenAddrs: listenAddrs,

		// TODO(roasbeef): derive proper onion key based on rotation
		// schedule
		sphinx: htlcswitch.NewOnionProcessor(sphinxRouter),

		persistentPeers:         make(map[string]struct{}),
		persistentPeersBackoff:  make(map[string]time.Duration),
		persistentConnReqs:      make(map[string][]*connmgr.ConnReq),
		persistentRetryCancels:  make(map[string]chan struct{}),
		ignorePeerTermination:   make(map[*peer]struct{}),
		scheduledPeerConnection: make(map[string]func()),

		peersByPub:                make(map[string]*peer),
		inboundPeers:              make(map[string]*peer),
		outboundPeers:             make(map[string]*peer),
		peerConnectedListeners:    make(map[string][]chan<- lnpeer.Peer),
		peerDisconnectedListeners: make(map[string][]chan<- struct{}),

		globalFeatures: lnwire.NewFeatureVector(globalFeatures,
			lnwire.GlobalFeatures),
		quit: make(chan struct{}),
	}

	s.witnessBeacon = &preimageBeacon{
		invoices:    s.invoices,
		wCache:      chanDB.NewWitnessCache(),
		subscribers: make(map[uint64]*preimageSubscriber),
	}

	// If the debug HTLC flag is on, then we invoice a "master debug"
	// invoice which all outgoing payments will be sent and all incoming
	// HTLCs with the debug R-Hash immediately settled.
	if cfg.DebugHTLC {
		kiloCoin := btcutil.Amount(btcutil.SatoshiPerBitcoin * 1000)
		s.invoices.AddDebugInvoice(kiloCoin, invoices.DebugPre)
		srvrLog.Debugf("Debug HTLC invoice inserted, preimage=%x, hash=%x",
			invoices.DebugPre[:], invoices.DebugHash[:])
	}

	_, currentHeight, err := s.cc.chainIO.GetBestBlock()
	if err != nil {
		return nil, err
	}

	s.htlcSwitch, err = htlcswitch.New(htlcswitch.Config{
		DB:      chanDB,
		SelfKey: s.identityPriv.PubKey(),
		LocalChannelClose: func(pubKey []byte,
			request *htlcswitch.ChanClose) {

			peer, err := s.FindPeerByPubStr(string(pubKey))
			if err != nil {
				srvrLog.Errorf("unable to close channel, peer"+
					" with %v id can't be found: %v",
					pubKey, err,
				)
				return
			}

			select {
			case peer.localCloseChanReqs <- request:
				srvrLog.Infof("Local close channel request "+
					"delivered to peer: %x", pubKey[:])
			case <-peer.quit:
				srvrLog.Errorf("Unable to deliver local close "+
					"channel request to peer %x, err: %v",
					pubKey[:], err)
			}
		},
		FwdingLog:              chanDB.ForwardingLog(),
		SwitchPackager:         channeldb.NewSwitchPackager(),
		ExtractErrorEncrypter:  s.sphinx.ExtractErrorEncrypter,
		FetchLastChannelUpdate: s.fetchLastChanUpdate(),
		Notifier:               s.cc.chainNotifier,
		FwdEventTicker: ticker.New(
			htlcswitch.DefaultFwdEventInterval),
		LogEventTicker: ticker.New(
			htlcswitch.DefaultLogInterval),
		NotifyActiveChannel:   s.channelNotifier.NotifyActiveChannelEvent,
		NotifyInactiveChannel: s.channelNotifier.NotifyInactiveChannelEvent,
	}, uint32(currentHeight))
	if err != nil {
		return nil, err
	}

	chanStatusMgrCfg := &netann.ChanStatusConfig{
		ChanStatusSampleInterval: cfg.ChanStatusSampleInterval,
		ChanEnableTimeout:        cfg.ChanEnableTimeout,
		ChanDisableTimeout:       cfg.ChanDisableTimeout,
		OurPubKey:                privKey.PubKey(),
		MessageSigner:            s.nodeSigner,
		IsChannelActive:          s.htlcSwitch.HasActiveLink,
		ApplyChannelUpdate:       s.applyChannelUpdate,
		DB:                       chanDB,
		Graph:                    chanDB.ChannelGraph(),
	}

	chanStatusMgr, err := netann.NewChanStatusManager(chanStatusMgrCfg)
	if err != nil {
		return nil, err
	}
	s.chanStatusMgr = chanStatusMgr

	// If enabled, use either UPnP or NAT-PMP to automatically configure
	// port forwarding for users behind a NAT.
	if cfg.NAT {
		srvrLog.Info("Scanning local network for a UPnP enabled device")

		discoveryTimeout := time.Duration(10 * time.Second)

		ctx, cancel := context.WithTimeout(
			context.Background(), discoveryTimeout,
		)
		defer cancel()
		upnp, err := nat.DiscoverUPnP(ctx)
		if err == nil {
			s.natTraversal = upnp
		} else {
			// If we were not able to discover a UPnP enabled device
			// on the local network, we'll fall back to attempting
			// to discover a NAT-PMP enabled device.
			srvrLog.Errorf("Unable to discover a UPnP enabled "+
				"device on the local network: %v", err)

			srvrLog.Info("Scanning local network for a NAT-PMP " +
				"enabled device")

			pmp, err := nat.DiscoverPMP(discoveryTimeout)
			if err != nil {
				err := fmt.Errorf("Unable to discover a "+
					"NAT-PMP enabled device on the local "+
					"network: %v", err)
				srvrLog.Error(err)
				return nil, err
			}

			s.natTraversal = pmp
		}
	}

	// If we were requested to automatically configure port forwarding,
	// we'll use the ports that the server will be listening on.
	externalIPStrings := make([]string, len(cfg.ExternalIPs))
	for idx, ip := range cfg.ExternalIPs {
		externalIPStrings[idx] = ip.String()
	}
	if s.natTraversal != nil {
		listenPorts := make([]uint16, 0, len(listenAddrs))
		for _, listenAddr := range listenAddrs {
			// At this point, the listen addresses should have
			// already been normalized, so it's safe to ignore the
			// errors.
			_, portStr, _ := net.SplitHostPort(listenAddr.String())
			port, _ := strconv.Atoi(portStr)

			listenPorts = append(listenPorts, uint16(port))
		}

		ips, err := s.configurePortForwarding(listenPorts...)
		if err != nil {
			srvrLog.Errorf("Unable to automatically set up port "+
				"forwarding using %s: %v",
				s.natTraversal.Name(), err)
		} else {
			srvrLog.Infof("Automatically set up port forwarding "+
				"using %s to advertise external IP",
				s.natTraversal.Name())
			externalIPStrings = append(externalIPStrings, ips...)
		}
	}

	// If external IP addresses have been specified, add those to the list
	// of this server's addresses.
	externalIPs, err := lncfg.NormalizeAddresses(
		externalIPStrings, strconv.Itoa(defaultPeerPort),
		cfg.net.ResolveTCPAddr,
	)
	if err != nil {
		return nil, err
	}
	selfAddrs := make([]net.Addr, 0, len(externalIPs))
	for _, ip := range externalIPs {
		selfAddrs = append(selfAddrs, ip)
	}

	// If we were requested to route connections through Tor and to
	// automatically create an onion service, we'll initiate our Tor
	// controller and establish a connection to the Tor server.
	if cfg.Tor.Active && (cfg.Tor.V2 || cfg.Tor.V3) {
		s.torController = tor.NewController(cfg.Tor.Control)
	}

	chanGraph := chanDB.ChannelGraph()

	// We'll now reconstruct a node announcement based on our current
	// configuration so we can send it out as a sort of heart beat within
	// the network.
	//
	// We'll start by parsing the node color from configuration.
	color, err := parseHexColor(cfg.Color)
	if err != nil {
		srvrLog.Errorf("unable to parse color: %v\n", err)
		return nil, err
	}

	// If no alias is provided, default to first 10 characters of public
	// key.
	alias := cfg.Alias
	if alias == "" {
		alias = hex.EncodeToString(serializedPubKey[:10])
	}
	nodeAlias, err := lnwire.NewNodeAlias(alias)
	if err != nil {
		return nil, err
	}
	selfNode := &channeldb.LightningNode{
		HaveNodeAnnouncement: true,
		LastUpdate:           time.Now(),
		Addresses:            selfAddrs,
		Alias:                nodeAlias.String(),
		Features:             s.globalFeatures,
		Color:                color,
	}
	copy(selfNode.PubKeyBytes[:], privKey.PubKey().SerializeCompressed())

	// Based on the disk representation of the node announcement generated
	// above, we'll generate a node announcement that can go out on the
	// network so we can properly sign it.
	nodeAnn, err := selfNode.NodeAnnouncement(false)
	if err != nil {
		return nil, fmt.Errorf("unable to gen self node ann: %v", err)
	}

	// With the announcement generated, we'll sign it to properly
	// authenticate the message on the network.
	authSig, err := discovery.SignAnnouncement(
		s.nodeSigner, s.identityPriv.PubKey(), nodeAnn,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to generate signature for "+
			"self node announcement: %v", err)
	}
	selfNode.AuthSigBytes = authSig.Serialize()
	nodeAnn.Signature, err = lnwire.NewSigFromRawSignature(
		selfNode.AuthSigBytes,
	)
	if err != nil {
		return nil, err
	}

	// Finally, we'll update the representation on disk, and update our
	// cached in-memory version as well.
	if err := chanGraph.SetSourceNode(selfNode); err != nil {
		return nil, fmt.Errorf("can't set self node: %v", err)
	}
	s.currentNodeAnn = nodeAnn

	s.chanRouter, err = routing.New(routing.Config{
		Graph:     chanGraph,
		Chain:     cc.chainIO,
		ChainView: cc.chainView,
		SendToSwitch: func(firstHop lnwire.ShortChannelID,
			htlcAdd *lnwire.UpdateAddHTLC,
			circuit *sphinx.Circuit) ([32]byte, error) {

			// Using the created circuit, initialize the error
			// decrypter so we can parse+decode any failures
			// incurred by this payment within the switch.
			errorDecryptor := &htlcswitch.SphinxErrorDecrypter{
				OnionErrorDecrypter: sphinx.NewOnionErrorDecrypter(circuit),
			}

			return s.htlcSwitch.SendHTLC(
				firstHop, htlcAdd, errorDecryptor,
			)
		},
		ChannelPruneExpiry: routing.DefaultChannelPruneExpiry,
		GraphPruneInterval: time.Duration(time.Hour),
		QueryBandwidth: func(edge *channeldb.ChannelEdgeInfo) lnwire.MilliSatoshi {
			// If we aren't on either side of this edge, then we'll
			// just thread through the capacity of the edge as we
			// know it.
			if !bytes.Equal(edge.NodeKey1Bytes[:], selfNode.PubKeyBytes[:]) &&
				!bytes.Equal(edge.NodeKey2Bytes[:], selfNode.PubKeyBytes[:]) {

				return lnwire.NewMSatFromSatoshis(edge.Capacity)
			}

			cid := lnwire.NewChanIDFromOutPoint(&edge.ChannelPoint)
			link, err := s.htlcSwitch.GetLink(cid)
			if err != nil {
				// If the link isn't online, then we'll report
				// that it has zero bandwidth to the router.
				return 0
			}

			// If the link is found within the switch, but it isn't
			// yet eligible to forward any HTLCs, then we'll treat
			// it as if it isn't online in the first place.
			if !link.EligibleToForward() {
				return 0
			}

			// Otherwise, we'll return the current best estimate
			// for the available bandwidth for the link.
			return link.Bandwidth()
		},
		AssumeChannelValid: cfg.Routing.UseAssumeChannelValid(),
	})
	if err != nil {
		return nil, fmt.Errorf("can't create router: %v", err)
	}

	chanSeries := discovery.NewChanSeries(s.chanDB.ChannelGraph())
	gossipMessageStore, err := discovery.NewMessageStore(s.chanDB)
	if err != nil {
		return nil, err
	}
	waitingProofStore, err := channeldb.NewWaitingProofStore(s.chanDB)
	if err != nil {
		return nil, err
	}

	s.authGossiper = discovery.New(discovery.Config{
		Router:               s.chanRouter,
		Notifier:             s.cc.chainNotifier,
		ChainHash:            *activeNetParams.GenesisHash,
		Broadcast:            s.BroadcastMessage,
		ChanSeries:           chanSeries,
		NotifyWhenOnline:     s.NotifyWhenOnline,
		NotifyWhenOffline:    s.NotifyWhenOffline,
		ProofMatureDelta:     0,
		TrickleDelay:         time.Millisecond * time.Duration(cfg.TrickleDelay),
		RetransmitDelay:      time.Minute * 30,
		WaitingProofStore:    waitingProofStore,
		MessageStore:         gossipMessageStore,
		AnnSigner:            s.nodeSigner,
		RotateTicker:         ticker.New(discovery.DefaultSyncerRotationInterval),
		HistoricalSyncTicker: ticker.New(cfg.HistoricalSyncInterval),
		NumActiveSyncers:     cfg.NumGraphSyncPeers,
	},
		s.identityPriv.PubKey(),
	)

	utxnStore, err := newNurseryStore(activeNetParams.GenesisHash, chanDB)
	if err != nil {
		srvrLog.Errorf("unable to create nursery store: %v", err)
		return nil, err
	}

	srvrLog.Tracef("Sweeper batch window duration: %v",
		sweep.DefaultBatchWindowDuration)

	sweeperStore, err := sweep.NewSweeperStore(
		chanDB, activeNetParams.GenesisHash,
	)
	if err != nil {
		srvrLog.Errorf("unable to create sweeper store: %v", err)
		return nil, err
	}

	s.sweeper = sweep.New(&sweep.UtxoSweeperConfig{
		FeeEstimator: cc.feeEstimator,
		GenSweepScript: func() ([]byte, error) {
			return newSweepPkScript(cc.wallet)
		},
		Signer:             cc.wallet.Cfg.Signer,
		PublishTransaction: cc.wallet.PublishTransaction,
		NewBatchTimer: func() <-chan time.Time {
			return time.NewTimer(sweep.DefaultBatchWindowDuration).C
		},
		SweepTxConfTarget:    6,
		Notifier:             cc.chainNotifier,
		ChainIO:              cc.chainIO,
		Store:                sweeperStore,
		MaxInputsPerTx:       sweep.DefaultMaxInputsPerTx,
		MaxSweepAttempts:     sweep.DefaultMaxSweepAttempts,
		NextAttemptDeltaFunc: sweep.DefaultNextAttemptDeltaFunc,
	})

	s.utxoNursery = newUtxoNursery(&NurseryConfig{
		ChainIO:             cc.chainIO,
		ConfDepth:           1,
		FetchClosedChannels: chanDB.FetchClosedChannels,
		FetchClosedChannel:  chanDB.FetchClosedChannel,
		Notifier:            cc.chainNotifier,
		PublishTransaction:  cc.wallet.PublishTransaction,
		Store:               utxnStore,
		SweepInput:          s.sweeper.SweepInput,
	})

	// Construct a closure that wraps the htlcswitch's CloseLink method.
	closeLink := func(chanPoint *wire.OutPoint,
		closureType htlcswitch.ChannelCloseType) {
		// TODO(conner): Properly respect the update and error channels
		// returned by CloseLink.
		s.htlcSwitch.CloseLink(chanPoint, closureType, 0)
	}

	// We will use the following channel to reliably hand off contract
	// breach events from the ChannelArbitrator to the breachArbiter,
	contractBreaches := make(chan *ContractBreachEvent, 1)

	s.chainArb = contractcourt.NewChainArbitrator(contractcourt.ChainArbitratorConfig{
		ChainHash:              *activeNetParams.GenesisHash,
		IncomingBroadcastDelta: defaultIncomingBroadcastDelta,
		OutgoingBroadcastDelta: defaultOutgoingBroadcastDelta,
		NewSweepAddr: func() ([]byte, error) {
			return newSweepPkScript(cc.wallet)
		},
		PublishTx: cc.wallet.PublishTransaction,
		DeliverResolutionMsg: func(msgs ...contractcourt.ResolutionMsg) error {
			for _, msg := range msgs {
				err := s.htlcSwitch.ProcessContractResolution(msg)
				if err != nil {
					return err
				}
			}
			return nil
		},
		IncubateOutputs: func(chanPoint wire.OutPoint,
			commitRes *lnwallet.CommitOutputResolution,
			outHtlcRes *lnwallet.OutgoingHtlcResolution,
			inHtlcRes *lnwallet.IncomingHtlcResolution,
			broadcastHeight uint32) error {

			var (
				inRes  []lnwallet.IncomingHtlcResolution
				outRes []lnwallet.OutgoingHtlcResolution
			)
			if inHtlcRes != nil {
				inRes = append(inRes, *inHtlcRes)
			}
			if outHtlcRes != nil {
				outRes = append(outRes, *outHtlcRes)
			}

			return s.utxoNursery.IncubateOutputs(
				chanPoint, commitRes, outRes, inRes,
				broadcastHeight,
			)
		},
		PreimageDB:   s.witnessBeacon,
		Notifier:     cc.chainNotifier,
		Signer:       cc.wallet.Cfg.Signer,
		FeeEstimator: cc.feeEstimator,
		ChainIO:      cc.chainIO,
		MarkLinkInactive: func(chanPoint wire.OutPoint) error {
			chanID := lnwire.NewChanIDFromOutPoint(&chanPoint)
			s.htlcSwitch.RemoveLink(chanID)
			return nil
		},
		IsOurAddress: cc.wallet.IsOurAddress,
		ContractBreach: func(chanPoint wire.OutPoint,
			breachRet *lnwallet.BreachRetribution) error {
			event := &ContractBreachEvent{
				ChanPoint:         chanPoint,
				ProcessACK:        make(chan error, 1),
				BreachRetribution: breachRet,
			}

			// Send the contract breach event to the breachArbiter.
			select {
			case contractBreaches <- event:
			case <-s.quit:
				return ErrServerShuttingDown
			}

			// Wait for the breachArbiter to ACK the event.
			select {
			case err := <-event.ProcessACK:
				return err
			case <-s.quit:
				return ErrServerShuttingDown
			}
		},
		DisableChannel:      s.chanStatusMgr.RequestDisable,
		Sweeper:             s.sweeper,
		Registry:            s.invoices,
		NotifyClosedChannel: s.channelNotifier.NotifyClosedChannelEvent,
	}, chanDB)

	s.breachArbiter = newBreachArbiter(&BreachConfig{
		CloseLink: closeLink,
		DB:        chanDB,
		Estimator: s.cc.feeEstimator,
		GenSweepScript: func() ([]byte, error) {
			return newSweepPkScript(cc.wallet)
		},
		Notifier:           cc.chainNotifier,
		PublishTransaction: cc.wallet.PublishTransaction,
		ContractBreaches:   contractBreaches,
		Signer:             cc.wallet.Cfg.Signer,
		Store:              newRetributionStore(chanDB),
	})

	// Select the configuration and furnding parameters for Bitcoin or
	// Litecoin, depending on the primary registered chain.
	primaryChain := registeredChains.PrimaryChain()
	chainCfg := cfg.Bitcoin
	minRemoteDelay := minBtcRemoteDelay
	maxRemoteDelay := maxBtcRemoteDelay
	if primaryChain == litecoinChain {
		chainCfg = cfg.Litecoin
		minRemoteDelay = minLtcRemoteDelay
		maxRemoteDelay = maxLtcRemoteDelay
	}

	var chanIDSeed [32]byte
	if _, err := rand.Read(chanIDSeed[:]); err != nil {
		return nil, err
	}
	s.fundingMgr, err = newFundingManager(fundingConfig{
		IDKey:              privKey.PubKey(),
		Wallet:             cc.wallet,
		PublishTransaction: cc.wallet.PublishTransaction,
		Notifier:           cc.chainNotifier,
		FeeEstimator:       cc.feeEstimator,
		SignMessage: func(pubKey *btcec.PublicKey,
			msg []byte) (*btcec.Signature, error) {

			if pubKey.IsEqual(privKey.PubKey()) {
				return s.nodeSigner.SignMessage(pubKey, msg)
			}

			return cc.msgSigner.SignMessage(pubKey, msg)
		},
		CurrentNodeAnnouncement: func() (lnwire.NodeAnnouncement, error) {
			return s.genNodeAnnouncement(true)
		},
		SendAnnouncement: func(msg lnwire.Message,
			optionalFields ...discovery.OptionalMsgField) chan error {

			return s.authGossiper.ProcessLocalAnnouncement(
				msg, privKey.PubKey(), optionalFields...,
			)
		},
		NotifyWhenOnline: s.NotifyWhenOnline,
		TempChanIDSeed:   chanIDSeed,
		FindChannel: func(chanID lnwire.ChannelID) (
			*channeldb.OpenChannel, error) {

			dbChannels, err := chanDB.FetchAllChannels()
			if err != nil {
				return nil, err
			}

			for _, channel := range dbChannels {
				if chanID.IsChanPoint(&channel.FundingOutpoint) {
					return channel, nil
				}
			}

			return nil, fmt.Errorf("unable to find channel")
		},
		DefaultRoutingPolicy: cc.routingPolicy,
		NumRequiredConfs: func(chanAmt btcutil.Amount,
			pushAmt lnwire.MilliSatoshi) uint16 {
			// For large channels we increase the number
			// of confirmations we require for the
			// channel to be considered open. As it is
			// always the responder that gets to choose
			// value, the pushAmt is value being pushed
			// to us. This means we have more to lose
			// in the case this gets re-orged out, and
			// we will require more confirmations before
			// we consider it open.
			// TODO(halseth): Use Litecoin params in case
			// of LTC channels.

			// In case the user has explicitly specified
			// a default value for the number of
			// confirmations, we use it.
			defaultConf := uint16(chainCfg.DefaultNumChanConfs)
			if defaultConf != 0 {
				return defaultConf
			}

			// If not we return a value scaled linearly
			// between 3 and 6, depending on channel size.
			// TODO(halseth): Use 1 as minimum?
			minConf := uint64(3)
			maxConf := uint64(6)
			maxChannelSize := uint64(
				lnwire.NewMSatFromSatoshis(maxFundingAmount))
			stake := lnwire.NewMSatFromSatoshis(chanAmt) + pushAmt
			conf := maxConf * uint64(stake) / maxChannelSize
			if conf < minConf {
				conf = minConf
			}
			if conf > maxConf {
				conf = maxConf
			}
			return uint16(conf)
		},
		RequiredRemoteDelay: func(chanAmt btcutil.Amount) uint16 {
			// We scale the remote CSV delay (the time the
			// remote have to claim funds in case of a unilateral
			// close) linearly from minRemoteDelay blocks
			// for small channels, to maxRemoteDelay blocks
			// for channels of size maxFundingAmount.
			// TODO(halseth): Litecoin parameter for LTC.

			// In case the user has explicitly specified
			// a default value for the remote delay, we
			// use it.
			defaultDelay := uint16(chainCfg.DefaultRemoteDelay)
			if defaultDelay > 0 {
				return defaultDelay
			}

			// If not we scale according to channel size.
			delay := uint16(btcutil.Amount(maxRemoteDelay) *
				chanAmt / maxFundingAmount)
			if delay < minRemoteDelay {
				delay = minRemoteDelay
			}
			if delay > maxRemoteDelay {
				delay = maxRemoteDelay
			}
			return delay
		},
		WatchNewChannel: func(channel *channeldb.OpenChannel,
			peerKey *btcec.PublicKey) error {

			// First, we'll mark this new peer as a persistent peer
			// for re-connection purposes.
			s.mu.Lock()
			pubStr := string(peerKey.SerializeCompressed())
			s.persistentPeers[pubStr] = struct{}{}
			s.mu.Unlock()

			// With that taken care of, we'll send this channel to
			// the chain arb so it can react to on-chain events.
			return s.chainArb.WatchNewChannel(channel)
		},
		ReportShortChanID: func(chanPoint wire.OutPoint) error {
			cid := lnwire.NewChanIDFromOutPoint(&chanPoint)
			return s.htlcSwitch.UpdateShortChanID(cid)
		},
		RequiredRemoteChanReserve: func(chanAmt,
			dustLimit btcutil.Amount) btcutil.Amount {

			// By default, we'll require the remote peer to maintain
			// at least 1% of the total channel capacity at all
			// times. If this value ends up dipping below the dust
			// limit, then we'll use the dust limit itself as the
			// reserve as required by BOLT #2.
			reserve := chanAmt / 100
			if reserve < dustLimit {
				reserve = dustLimit
			}

			return reserve
		},
		RequiredRemoteMaxValue: func(chanAmt btcutil.Amount) lnwire.MilliSatoshi {
			// By default, we'll allow the remote peer to fully
			// utilize the full bandwidth of the channel, minus our
			// required reserve.
			reserve := lnwire.NewMSatFromSatoshis(chanAmt / 100)
			return lnwire.NewMSatFromSatoshis(chanAmt) - reserve
		},
		RequiredRemoteMaxHTLCs: func(chanAmt btcutil.Amount) uint16 {
			// By default, we'll permit them to utilize the full
			// channel bandwidth.
			return uint16(input.MaxHTLCNumber / 2)
		},
		ZombieSweeperInterval:  1 * time.Minute,
		ReservationTimeout:     10 * time.Minute,
		MinChanSize:            btcutil.Amount(cfg.MinChanSize),
		NotifyOpenChannelEvent: s.channelNotifier.NotifyOpenChannelEvent,
	})
	if err != nil {
		return nil, err
	}

	// Next, we'll assemble the sub-system that will maintain an on-disk
	// static backup of the latest channel state.
	chanNotifier := &channelNotifier{
		chanNotifier: s.channelNotifier,
		addrs:        s.chanDB,
	}
	backupFile := chanbackup.NewMultiFile(cfg.BackupFilePath)
	startingChans, err := chanbackup.FetchStaticChanBackups(s.chanDB)
	if err != nil {
		return nil, err
	}
	s.chanSubSwapper, err = chanbackup.NewSubSwapper(
		startingChans, chanNotifier, s.cc.keyRing, backupFile,
	)
	if err != nil {
		return nil, err
	}

	// Create the connection manager which will be responsible for
	// maintaining persistent outbound connections and also accepting new
	// incoming connections
	cmgr, err := connmgr.New(&connmgr.Config{
		Listeners:      listeners,
		OnAccept:       s.InboundPeerConnected,
		RetryDuration:  time.Second * 5,
		TargetOutbound: 100,
		Dial:           noiseDial(s.identityPriv),
		OnConnection:   s.OutboundPeerConnected,
	})
	if err != nil {
		return nil, err
	}
	s.connMgr = cmgr

	return s, nil
}
