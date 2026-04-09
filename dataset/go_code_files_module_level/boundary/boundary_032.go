func rescan(chain ChainSource, options ...RescanOption) error {
	// First, we'll apply the set of default options, then serially apply
	// all the options that've been passed in.
	ro := defaultRescanOptions()
	ro.endBlock = &waddrmgr.BlockStamp{
		Hash:   chainhash.Hash{},
		Height: 0,
	}
	for _, option := range options {
		option(ro)
	}

	// If we have something to watch, create a watch list. The watch list
	// can be composed of a set of scripts, outpoints, and txids.
	for _, addr := range ro.watchAddrs {
		script, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return err
		}

		ro.watchList = append(ro.watchList, script)
	}
	for _, input := range ro.watchInputs {
		ro.watchList = append(ro.watchList, input.PkScript)
	}

	// Check that we have either an end block or a quit channel.
	if ro.endBlock != nil {
		// If the end block hash is non-nil, then we'll query the
		// database to find out the stop height.
		if (ro.endBlock.Hash != chainhash.Hash{}) {
			_, height, err := chain.GetBlockHeader(
				&ro.endBlock.Hash,
			)
			if err != nil {
				ro.endBlock.Hash = chainhash.Hash{}
			} else {
				ro.endBlock.Height = int32(height)
			}
		}

		// If the ending hash it nil, then check to see if the target
		// height is non-nil. If not, then we'll use this to find the
		// stopping hash.
		if (ro.endBlock.Hash == chainhash.Hash{}) {
			if ro.endBlock.Height != 0 {
				header, err := chain.GetBlockHeaderByHeight(
					uint32(ro.endBlock.Height),
				)
				if err == nil {
					ro.endBlock.Hash = header.BlockHash()
				} else {
					ro.endBlock = &waddrmgr.BlockStamp{}
				}
			}
		}
	} else {
		ro.endBlock = &waddrmgr.BlockStamp{}
	}

	// If we don't have a quit channel, and the end height is still
	// unspecified, then we'll exit out here.
	if ro.quit == nil && ro.endBlock.Height == 0 {
		return fmt.Errorf("Rescan request must specify a quit channel" +
			" or valid end block")
	}

	// Track our position in the chain.
	var (
		curHeader wire.BlockHeader
		curStamp  waddrmgr.BlockStamp
	)

	// If no start block is specified, start the scan from our current best
	// block.
	if ro.startBlock == nil {
		bs, err := chain.BestBlock()
		if err != nil {
			return err
		}
		ro.startBlock = bs
	}
	curStamp = *ro.startBlock

	// To find our starting block, either the start hash should be set, or
	// the start height should be set. If neither is, then we'll be
	// starting from the genesis block.
	if (curStamp.Hash != chainhash.Hash{}) {
		header, height, err := chain.GetBlockHeader(&curStamp.Hash)
		if err == nil {
			curHeader = *header
			curStamp.Height = int32(height)
		} else {
			curStamp.Hash = chainhash.Hash{}
		}
	}
	if (curStamp.Hash == chainhash.Hash{}) {
		if curStamp.Height == 0 {
			curStamp.Hash = *chain.ChainParams().GenesisHash
		} else {
			header, err := chain.GetBlockHeaderByHeight(
				uint32(curStamp.Height),
			)
			if err == nil {
				curHeader = *header
				curStamp.Hash = curHeader.BlockHash()
			} else {
				curHeader = chain.ChainParams().GenesisBlock.Header
				curStamp.Hash = *chain.ChainParams().GenesisHash
				curStamp.Height = 0
			}
		}
	}

	// Now that we've determined the starting point of our rescan, we can
	// begin processing updates from the client.
	var updates []*updateOptions

	// We'll need to ensure that the backing chain has actually caught up to
	// the rescan's starting height.
	bestBlock, err := chain.BestBlock()
	if err != nil {
		return err
	}

	// If it hasn't, we'll subscribe for block notifications at tip and wait
	// until we receive a notification for a block with the rescan's
	// starting height.
	if bestBlock.Height < curStamp.Height {
		log.Debugf("Waiting to catch up to the rescan start height=%d "+
			"from height=%d", curStamp.Height, bestBlock.Height)

		blockSubscription, err := chain.Subscribe(
			uint32(bestBlock.Height),
		)
		if err != nil {
			return err
		}

	waitUntilSynced:
		for {
			select {
			// We'll make sure to process any updates while we're
			// syncing to prevent blocking the client.
			case update := <-ro.update:
				updates = append(updates, update)

			// A new block notification for the tip of the chain has
			// arrived. We'll determine we've caught up to the
			// rescan's starting height by receiving a block
			// connected notification for the same height.
			case ntfn, ok := <-blockSubscription.Notifications:
				if !ok {
					return errors.New("rescan block " +
						"subscription was canceled " +
						"while waiting to catch up")
				}

				if _, ok := ntfn.(*blockntfns.Connected); !ok {
					continue
				}
				if ntfn.Height() < uint32(curStamp.Height) {
					continue
				}

				break waitUntilSynced

			case <-ro.quit:
				blockSubscription.Cancel()
				return ErrRescanExit
			}
		}

		blockSubscription.Cancel()

		// If any updates were queued while waiting to catch up to the
		// start height of the rescan, apply them now.
		for _, upd := range updates {
			_, err := ro.updateFilter(
				chain, upd, &curStamp, &curHeader,
			)
			if err != nil {
				return err
			}
		}
	}

	log.Debugf("Starting rescan from known block %d (%s)", curStamp.Height,
		curStamp.Hash)

	// Compare the start time to the start block. If the start time is
	// later, cycle through blocks until we find a block timestamp later
	// than the start time, and begin filter download at that block. Since
	// time is non-monotonic between blocks, we look for the first block to
	// trip the switch, and download filters from there, rather than
	// checking timestamps at each block.
	scanning := ro.startTime.Before(curHeader.Timestamp)

	var blockSubscription *blockntfns.Subscription

	// blockRetryInterval is the interval in which we'll continually re-try
	// to fetch the latest filter from our peers.
	//
	// TODO(roasbeef): add exponential back-off
	blockRetryInterval := time.Millisecond * 100

	// blockReFetchTimer is a stoppable timer that we'll use to reminder
	// ourselves to refetch a block in the case that we're unable to fetch
	// the filter for a block the first time around.
	var (
		blockReFetchTimer *time.Timer
		reFetchMtx        sync.Mutex
	)

	resetBlockReFetchTimer := func(headerTip wire.BlockHeader, height uint32) {
		// If so, then we'll avoid notifying the block, and will
		// instead add this to our retry queue, as we should be getting
		// block disconnected notifications in short order.
		if blockReFetchTimer != nil {
			blockReFetchTimer.Stop()
		}

		log.Infof("Setting timer to attempt to re-fetch filter for "+
			"hash=%v, height=%v", headerTip.BlockHash(), height)

		blockReFetch := func() {
			// If we're unable to process notifications at the
			// moment (due to not being current), we'll reset our
			// timer.
			reFetchMtx.Lock()
			if blockReFetchTimer != nil && blockSubscription == nil {
				if !blockReFetchTimer.Stop() {
					<-blockReFetchTimer.C
				}
				blockReFetchTimer.Reset(blockRetryInterval)

				reFetchMtx.Unlock()
				return
			}
			reFetchMtx.Unlock()

			log.Infof("Resending rescan header for block hash=%v, "+
				"height=%v", headerTip.BlockHash(), height)

			ntfn := blockntfns.NewBlockConnected(headerTip, height)
			select {
			case blockSubscription.Notifications <- ntfn:
			case <-ro.quit:
			}
		}

		// We'll start a timer to re-send this header so we re-process
		// if in the case that we don't get a re-org soon afterwards.
		reFetchMtx.Lock()
		blockReFetchTimer = time.AfterFunc(
			blockRetryInterval, blockReFetch,
		)
		reFetchMtx.Unlock()
	}

	// We'll need to keep track of whether we are current with the chain in
	// order to properly recover from a re-org. We'll start by assuming that
	// we are not current in order to catch up from the starting point to
	// the tip of the chain.
	current := false

	// handleBlockConnected is a closure that handles a new block connected
	// notification.
	//
	// TODO(wilmer): refactor this and handleBlockDisconnected into their
	// own methods.
	handleBlockConnected := func(ntfn *blockntfns.Connected) error {
		// If we've somehow missed a header in the range, then we'll
		// mark ourselves as not current so we can walk down the chain
		// and notify the callers of blocks we may have missed.
		//
		// It's possible due to the nature of the current subscription
		// system that we get a duplicate block. We'll catch this and
		// continue forwards to avoid an unnecessary state transition
		// back to the !current state.
		header := ntfn.Header()
		if header.PrevBlock != curStamp.Hash &&
			header.BlockHash() != curStamp.Hash {

			current = false
			return fmt.Errorf("out of order block %v: "+
				"expected PrevBlock %v, got %v",
				header.BlockHash(), curStamp.Hash,
				header.PrevBlock)
		}

		// Do not process block until we have all filter headers. Don't
		// worry, the block will get re-queued every time there is a new
		// filter available. However, if it's a duplicate block
		// notification, then we can re-process it without any issues.
		nextBlockHeight := uint32(curStamp.Height + 1)
		_, err := chain.GetFilterHeaderByHeight(nextBlockHeight)
		if header.BlockHash() != curStamp.Hash && err != nil {
			log.Warnf("Missing filter header for height=%v, "+
				"skipping", curStamp.Height+1)

			return nil
		}

		// As this could be a re-try, we'll ensure that we don't
		// incorrectly increment our current time stamp.
		if curStamp.Hash != header.BlockHash() {
			curHeader = header
			curStamp.Hash = header.BlockHash()
			curStamp.Height++
		}

		log.Tracef("Rescan got block %d (%s)", curStamp.Height,
			curStamp.Hash)

		// We're only scanning if the header is beyond the horizon of
		// our start time.
		if !scanning {
			scanning = ro.startTime.Before(
				curHeader.Timestamp,
			)
		}

		// If we're not scanning or our watch list is empty, then we can
		// just notify the block without fetching any filters/blocks.
		if !scanning || len(ro.watchList) == 0 {
			if ro.ntfn.OnFilteredBlockConnected != nil {
				ro.ntfn.OnFilteredBlockConnected(
					curStamp.Height, &curHeader, nil,
				)
			}
			if ro.ntfn.OnBlockConnected != nil {
				ro.ntfn.OnBlockConnected(
					&curStamp.Hash, curStamp.Height,
					curHeader.Timestamp,
				)
			}

			return nil
		}

		// Otherwise, we'll attempt to fetch the filter to retrieve the
		// relevant transactions and notify them.
		queryOptions := NumRetries(0)
		blockFilter, err := chain.GetCFilter(
			curStamp.Hash, wire.GCSFilterRegular, queryOptions,
		)

		switch {
		// If the block index doesn't know about this block, then it's
		// likely we're mid re-org so we'll accept this as we account
		// for it below.
		case err == headerfs.ErrHashNotFound:

		case err != nil:
			return fmt.Errorf("unable to get filter for hash=%v: %v",
				curStamp.Hash, err)
		}

		// If the filter is nil, then this either means that we don't
		// have any peers to fetch this filter from, or the peer(s) that
		// we're trying to fetch from are in the progress of a re-org.
		if blockFilter == nil {
			// TODO(halseth): this is racy, as blocks can come in
			// before we refetch.
			resetBlockReFetchTimer(header, uint32(curStamp.Height))
			return nil
		}

		err = notifyBlockWithFilter(
			chain, ro, &curHeader, &curStamp, blockFilter,
		)
		if err != nil {
			return err
		}

		// We've successfully fetched this current block, so we'll reset
		// the retry timer back to nil.
		if blockReFetchTimer != nil {
			blockReFetchTimer.Stop()
			blockReFetchTimer = nil
		}

		return nil
	}

	// handleBlockDisconnected is a helper closure that handles a new block
	// disconnected notification.
	handleBlockDisconnected := func(ntfn *blockntfns.Disconnected) error {
		log.Debugf("Rescan disconnect block %d (%s)\n", curStamp.Height,
			curStamp.Hash)

		// Only deal with it if it's the current block we know about.
		// Otherwise, it's in the future.
		blockDisconnected := ntfn.Header()
		if blockDisconnected.BlockHash() != curStamp.Hash {
			return nil
		}

		// Run through notifications. This is all single-threaded. We
		// include deprecated calls as they're still used, for now.
		if ro.ntfn.OnFilteredBlockDisconnected != nil {
			ro.ntfn.OnFilteredBlockDisconnected(
				curStamp.Height, &curHeader,
			)
		}
		if ro.ntfn.OnBlockDisconnected != nil {
			ro.ntfn.OnBlockDisconnected(
				&curStamp.Hash, curStamp.Height,
				curHeader.Timestamp,
			)
		}

		curHeader = ntfn.ChainTip()
		curStamp.Hash = curHeader.BlockHash()
		curStamp.Height--

		// Now that we got a re-org, if we had a re-fetch timer going,
		// we'll reset it be at the new header tip.
		if blockReFetchTimer != nil {
			resetBlockReFetchTimer(
				curHeader, uint32(curStamp.Height),
			)
		}

		return nil
	}

	// Loop through blocks, one at a time. This relies on the underlying
	// chain source to deliver notifications in the correct order.
rescanLoop:
	for {
		// If we've reached the ending height or hash for this rescan,
		// then we'll exit.
		if curStamp.Hash == ro.endBlock.Hash ||
			(ro.endBlock.Height > 0 &&
				curStamp.Height == ro.endBlock.Height) {
			return nil
		}

		// If we're current, we wait for notifications that will be
		// delivered each time a block is connecting, disconnecting, or
		// we can an update to the filter we should be looking for.
		switch current {
		case true:
			// Wait for a signal that we have a newly connected
			// header and cfheader, or a newly disconnected header;
			// alternatively, forward ourselves to the next block
			// if possible.
			select {

			case <-ro.quit:
				return ErrRescanExit

			// An update mesage has just come across, if it points
			// to a prior point in the chain, then we may need to
			// rewind a bit in order to provide the client all its
			// requested client.
			case update := <-ro.update:
				rewound, err := ro.updateFilter(
					chain, update, &curStamp, &curHeader,
				)
				if err != nil {
					return err
				}

				// If we have to rewind our state, then we'll
				// mark ourselves as not current so we can walk
				// forward in the chain again until we we are
				// current. This is our way of doing a manual
				// rescan.
				if rewound {
					log.Tracef("Rewound to block %d (%s), "+
						"no longer current",
						curStamp.Height, curStamp.Hash)

					current = false
					blockSubscription.Cancel()
					blockSubscription = nil
				}

			case ntfn, ok := <-blockSubscription.Notifications:
				if !ok {
					return errors.New("rescan block " +
						"subscription was canceled")
				}

				switch ntfn := ntfn.(type) {
				case *blockntfns.Connected:
					err := handleBlockConnected(ntfn)
					if err != nil {
						log.Errorf("Unable to process "+
							"%v: %v", ntfn, err)
					}

				case *blockntfns.Disconnected:
					err := handleBlockDisconnected(ntfn)
					if err != nil {
						log.Errorf("Unable to process "+
							"%v: %v", ntfn, err)
					}

				default:
					log.Warnf("Received unhandled block "+
						"notification: %T", ntfn)
				}
			}

		// If we're not yet current, then we'll walk down the chain
		// until we reach the tip of the chain as we know it. At this
		// point, we'll be "current" again.
		case false:

			// Apply all queued filter updates.
		updateFilterLoop:
			for {
				select {
				case update := <-ro.update:
					_, err := ro.updateFilter(
						chain, update, &curStamp,
						&curHeader,
					)
					if err != nil {
						return err
					}

				default:
					break updateFilterLoop
				}
			}

			bestBlock, err := chain.BestBlock()
			if err != nil {
				return err
			}

			// Since we're not current, we try to manually advance
			// the block. If the next height is above the best
			// height known to the chain service, then we mark
			// ourselves as current and follow notifications.
			nextHeight := curStamp.Height + 1
			if nextHeight > bestBlock.Height {
				log.Debugf("Rescan became current at %d (%s), "+
					"subscribing to block notifications",
					curStamp.Height, curStamp.Hash)

				current = true

				// Ensure we cancel the old subscription if
				// we're going back to scan for missed blocks.
				if blockSubscription != nil {
					blockSubscription.Cancel()
				}

				// Subscribe to block notifications.
				blockSubscription, err = chain.Subscribe(
					uint32(curStamp.Height),
				)
				if err != nil {
					return fmt.Errorf("unable to register "+
						"block subscription: %v", err)
				}
				defer func() {
					if blockSubscription != nil {
						blockSubscription.Cancel()
						blockSubscription = nil
					}
				}()

				continue rescanLoop
			}

			// If the next height is known to the chain service,
			// then we'll fetch the next block and send a
			// notification, maybe also scanning the filters for
			// the block.
			header, err := chain.GetBlockHeaderByHeight(
				uint32(nextHeight),
			)
			if err != nil {
				return err
			}

			curHeader = *header
			curStamp.Height++
			curStamp.Hash = header.BlockHash()

			if !scanning {
				scanning = ro.startTime.Before(curHeader.Timestamp)
			}
			err = notifyBlock(chain, ro, curHeader, curStamp, scanning)
			if err != nil {
				return err
			}
		}
	}
}
