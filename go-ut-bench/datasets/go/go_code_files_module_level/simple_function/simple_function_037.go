func (db *ChainDB) SyncChainDB(ctx context.Context, client rpcutils.MasterBlockGetter,
	updateAllAddresses, updateAllVotes, newIndexes bool,
	updateExplorer chan *chainhash.Hash, barLoad chan *dbtypes.ProgressBarLoad) (int64, error) {
	// Note that we are doing a batch blockchain sync
	db.InBatchSync = true
	defer func() { db.InBatchSync = false }()

	// Get chain servers's best block
	nodeHeight, err := client.NodeHeight()
	if err != nil {
		return -1, fmt.Errorf("GetBestBlock failed: %v", err)
	}

	lastBlock, err := db.HeightDB()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info("blocks table is empty, starting fresh.")
		} else {
			return -1, fmt.Errorf("RetrieveBestBlockHeight: %v", err)
		}
	}

	// Remove indexes/constraints before an initial sync or when explicitly
	// requested to reindex and update spending information in the addresses
	// table.
	reindexing := newIndexes || lastBlock == -1
	if reindexing {
		// Remove any existing indexes.
		log.Info("Large bulk load: Removing indexes and disabling duplicate checks.")
		err = db.DeindexAll()
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			return lastBlock, err
		}

		// Disable duplicate checks on insert queries since the unique indexes
		// that enforce the constraints will not exist.
		db.EnableDuplicateCheckOnInsert(false)

		// Syncing blocks without indexes requires a UTXO cache to avoid
		// extremely expensive queries. Warm the UTXO cache if resuming an
		// interrupted initial sync.
		blocksToSync := nodeHeight - lastBlock
		if lastBlock > 0 && blocksToSync > 50 {
			log.Infof("Collecting all UTXO data prior to height %d...", lastBlock+1)
			utxos, err := RetrieveUTXOs(ctx, db.db)
			if err != nil {
				return -1, fmt.Errorf("RetrieveUTXOs: %v", err)
			}
			log.Infof("Pre-warming UTXO cache with %d UTXOs...", len(utxos))
			db.InitUtxoCache(utxos)
			log.Infof("UTXO cache is ready.")
		}
	} else {
		// When the unique indexes exist, inserts should check for conflicts
		// with the tables' constraints.
		db.EnableDuplicateCheckOnInsert(true)
	}

	// When reindexing or adding a large amount of data, ANALYZE tables.
	requireAnalyze := reindexing || nodeHeight-lastBlock > 10000

	if reindexing || updateAllAddresses || updateAllVotes {
		// Set meta.ibd_complete = FALSE.
		if err = SetIBDComplete(db.db, false); err != nil {
			return nodeHeight, fmt.Errorf("failed to set meta.ibd_complete: %v", err)
		}
	}

	// Safely send sync status updates on barLoad channel, and set the channel
	// to nil if the buffer is full.
	sendProgressUpdate := func(p *dbtypes.ProgressBarLoad) {
		if barLoad == nil {
			return
		}
		select {
		case barLoad <- p:
		default:
			log.Debugf("(*ChainDB).SyncChainDB: barLoad chan closed or full. Halting sync progress updates.")
			barLoad = nil
		}
	}

	// Safely send new block hash on updateExplorer channel, and set the channel
	// to nil if the buffer is full.
	sendPageData := func(hash *chainhash.Hash) {
		if updateExplorer == nil {
			return
		}
		select {
		case updateExplorer <- hash:
		default:
			log.Debugf("(*ChainDB).SyncChainDB: updateExplorer chan closed or full. Halting explorer updates.")
			updateExplorer = nil
		}
	}

	// Add the various updates that should run on successful sync.
	sendProgressUpdate(&dbtypes.ProgressBarLoad{
		Msg:   initialLoadSyncStatusMsg,
		BarID: dbtypes.InitialDBLoad,
	})
	// Addresses table sync should only run if bulk update is enabled.
	if updateAllAddresses {
		sendProgressUpdate(&dbtypes.ProgressBarLoad{
			Msg:   addressesSyncStatusMsg,
			BarID: dbtypes.AddressesTableSync,
		})
	}

	// Total and rate statistics
	var totalTxs, totalVins, totalVouts, totalAddresses int64
	var lastTxs, lastVins, lastVouts int64
	tickTime := 20 * time.Second
	ticker := time.NewTicker(tickTime)
	startTime := time.Now()
	o := sync.Once{}
	speedReporter := func() {
		ticker.Stop()
		totalElapsed := time.Since(startTime).Seconds()
		if int64(totalElapsed) == 0 {
			return
		}
		totalVoutPerSec := totalVouts / int64(totalElapsed)
		totalTxPerSec := totalTxs / int64(totalElapsed)
		if totalTxs == 0 {
			return
		}
		log.Infof("Avg. speed: %d tx/s, %d vout/s", totalTxPerSec, totalVoutPerSec)
	}
	speedReport := func() { o.Do(speedReporter) }
	defer speedReport()

	lastProgressUpdateTime := startTime

	// Start syncing blocks.
	startHeight := lastBlock + 1
	for ib := startHeight; ib <= nodeHeight; ib++ {
		// Check for quit signal.
		select {
		case <-ctx.Done():
			log.Infof("Rescan cancelled at height %d.", ib)
			return ib - 1, nil
		default:
		}

		// Progress logging
		if (ib-1)%rescanLogBlockChunk == 0 || ib == startHeight {
			if ib == 0 {
				log.Infof("Scanning genesis block into auxiliary chain db.")
			} else {
				endRangeBlock := rescanLogBlockChunk * (1 + (ib-1)/rescanLogBlockChunk)
				if endRangeBlock > nodeHeight {
					endRangeBlock = nodeHeight
				}
				log.Infof("Processing blocks %d to %d...", ib, endRangeBlock)

				if barLoad != nil {
					timeTakenPerBlock := (time.Since(lastProgressUpdateTime).Seconds() /
						float64(endRangeBlock-ib))
					sendProgressUpdate(&dbtypes.ProgressBarLoad{
						From:      ib,
						To:        nodeHeight,
						Timestamp: int64(timeTakenPerBlock * float64(nodeHeight-endRangeBlock)),
						Msg:       initialLoadSyncStatusMsg,
						BarID:     dbtypes.InitialDBLoad,
					})
					lastProgressUpdateTime = time.Now()
				}
			}
		}

		// Speed report
		select {
		case <-ticker.C:
			blocksPerSec := float64(ib-lastBlock) / tickTime.Seconds()
			txPerSec := float64(totalTxs-lastTxs) / tickTime.Seconds()
			vinsPerSec := float64(totalVins-lastVins) / tickTime.Seconds()
			voutPerSec := float64(totalVouts-lastVouts) / tickTime.Seconds()
			log.Infof("(%3d blk/s,%5d tx/s,%5d vin/sec,%5d vout/s)", int64(blocksPerSec),
				int64(txPerSec), int64(vinsPerSec), int64(voutPerSec))
			lastBlock, lastTxs = ib, totalTxs
			lastVins, lastVouts = totalVins, totalVouts
		default:
		}

		// Register for notification from stakedb when it connects this block.
		waitChan := db.stakeDB.WaitForHeight(ib)

		// Get the block, making it available to stakedb, which will signal on
		// the above channel when it is done connecting it.
		block, err := client.UpdateToBlock(ib)
		if err != nil {
			log.Errorf("UpdateToBlock (%d) failed: %v", ib, err)
			return ib - 1, fmt.Errorf("UpdateToBlock (%d) failed: %v", ib, err)
		}

		// Wait for our StakeDatabase to connect the block.
		var blockHash *chainhash.Hash
		select {
		case blockHash = <-waitChan:
		case <-ctx.Done():
			log.Infof("Rescan cancelled at height %d.", ib)
			return ib - 1, nil
		}
		if blockHash == nil {
			log.Errorf("stakedb says that block %d has come and gone", ib)
			return ib - 1, fmt.Errorf("stakedb says that block %d has come and gone", ib)
		}
		// If not master:
		//blockHash := <-client.WaitForHeight(ib)
		//block, err := client.Block(blockHash)
		// direct:
		//block, blockHash, err := rpcutils.GetBlock(ib, client)

		// Winning tickets from StakeDatabase, which just connected the block,
		// as signaled via the waitChan.
		tpi, ok := db.stakeDB.PoolInfo(*blockHash)
		if !ok {
			return ib - 1, fmt.Errorf("stakeDB.PoolInfo could not locate block %s", blockHash.String())
		}
		winners := tpi.Winners

		// Get the chainwork
		chainWork, err := client.GetChainWork(blockHash)
		if err != nil {
			return ib - 1, fmt.Errorf("GetChainWork failed (%s): %v", blockHash, err)
		}

		// Store data from this block in the database.
		isValid, isMainchain := true, true
		// updateExisting is ignored if dupCheck=false, but set it to true since
		// SyncChainDB is processing main chain blocks.
		updateExisting := true
		numVins, numVouts, numAddresses, err := db.StoreBlock(block.MsgBlock(), winners, isValid,
			isMainchain, updateExisting, !updateAllAddresses, !updateAllVotes, chainWork)
		if err != nil {
			return ib - 1, fmt.Errorf("StoreBlock failed: %v", err)
		}
		totalVins += numVins
		totalVouts += numVouts
		totalAddresses += numAddresses

		// Total transactions is the sum of regular and stake transactions.
		totalTxs += int64(len(block.STransactions()) + len(block.Transactions()))

		// Update explorer pages at intervals of 20 blocks if the update channel
		// is active (non-nil and not closed).
		if ib%20 == 0 && !updateAllAddresses {
			if updateExplorer != nil {
				log.Infof("Updating the explorer with information for block %v", ib)
				sendPageData(blockHash)
			}
		}

		// Update node height, the end condition for the loop.
		if nodeHeight, err = client.NodeHeight(); err != nil {
			return ib, fmt.Errorf("GetBestBlock failed: %v", err)
		}
	}

	// Final speed report
	speedReport()

	// After the last call to StoreBlock, synchronously update the project fund
	// and clear the general address balance cache.
	if err = db.FreshenAddressCaches(false, nil); err != nil {
		log.Warnf("FreshenAddressCaches: %v", err)
		err = nil // not an error with sync
	}

	// Signal the end of the initial load sync.
	sendProgressUpdate(&dbtypes.ProgressBarLoad{
		From:  nodeHeight,
		To:    nodeHeight,
		Msg:   initialLoadSyncStatusMsg,
		BarID: dbtypes.InitialDBLoad,
	})

	// Index and analyze tables.
	var analyzed bool
	if reindexing {
		// To build indexes, there must NOT be duplicate rows in terms of the
		// constraints defined by the unique indexes. Duplicate transactions,
		// vins, and vouts can end up in the tables when identical transactions
		// are included in multiple blocks. This happens when a block is
		// invalidated and the transactions are subsequently re-mined in another
		// block. Remove these before indexing.
		log.Infof("Finding and removing duplicate table rows before indexing...")
		if err = db.DeleteDuplicates(barLoad); err != nil {
			return 0, err
		}

		// Create all indexes.
		if err = db.IndexAll(barLoad); err != nil {
			return nodeHeight, fmt.Errorf("IndexAll failed: %v", err)
		}

		// Only reindex addresses and tickets tables here if not doing it below.
		if !updateAllAddresses {
			if err = db.IndexAddressTable(barLoad); err != nil {
				return nodeHeight, fmt.Errorf("IndexAddressTable failed: %v", err)
			}
		}
		if !updateAllVotes {
			if err = db.IndexTicketsTable(barLoad); err != nil {
				return nodeHeight, fmt.Errorf("IndexTicketsTable failed: %v", err)
			}
		}

		// Deep ANALYZE all tables.
		log.Infof("Performing an ANALYZE(%d) on all tables...", deepStatsTarget)
		if err = AnalyzeAllTables(db.db, deepStatsTarget); err != nil {
			return nodeHeight, fmt.Errorf("failed to ANALYZE tables: %v", err)
		}
		analyzed = true
	}

	// Batch update addresses table with spending info.
	if updateAllAddresses {
		// Analyze vins and table first.
		if !analyzed {
			log.Infof("Performing an ANALYZE(%d) on vins table...", deepStatsTarget)
			if err = AnalyzeTable(db.db, "vins", deepStatsTarget); err != nil {
				return nodeHeight, fmt.Errorf("failed to ANALYZE vins table: %v", err)
			}
		}

		// Remove existing indexes not on funding txns
		_ = db.DeindexAddressTable() // ignore errors for non-existent indexes
		log.Infof("Populating spending tx info in address table...")
		numAddresses, err := db.UpdateSpendingInfoInAllAddresses(barLoad)
		if err != nil {
			log.Errorf("UpdateSpendingInfoInAllAddresses FAILED: %v", err)
		}
		// Index addresses table
		log.Infof("Updated %d rows of address table", numAddresses)
		if err = db.IndexAddressTable(barLoad); err != nil {
			log.Errorf("IndexAddressTable FAILED: %v", err)
		}

		// Deep ANALYZE the newly indexed addresses table.
		log.Infof("Performing an ANALYZE(%d) on addresses table...", deepStatsTarget)
		if err = AnalyzeTable(db.db, "addresses", deepStatsTarget); err != nil {
			return nodeHeight, fmt.Errorf("failed to ANALYZE addresses table: %v", err)
		}
	}

	// Batch update tickets table with spending info.
	if updateAllVotes {
		// Analyze vins table first.
		if !analyzed {
			log.Infof("Performing an ANALYZE(%d) on vins table...", deepStatsTarget)
			if err = AnalyzeTable(db.db, "vins", deepStatsTarget); err != nil {
				return nodeHeight, fmt.Errorf("failed to ANALYZE vins table: %v", err)
			}
		}

		// Remove indexes not on funding txns (remove on tickets table indexes)
		_ = db.DeindexTicketsTable() // ignore errors for non-existent indexes
		db.EnableDuplicateCheckOnInsert(false)
		log.Infof("Populating spending tx info in tickets table...")
		numTicketsUpdated, err := db.UpdateSpendingInfoInAllTickets()
		if err != nil {
			log.Errorf("UpdateSpendingInfoInAllTickets FAILED: %v", err)
		}
		// Index tickets table
		log.Infof("Updated %d rows of address table", numTicketsUpdated)
		if err = db.IndexTicketsTable(barLoad); err != nil {
			log.Errorf("IndexTicketsTable FAILED: %v", err)
		}

		// Deep ANALYZE the newly indexed tickets table.
		log.Infof("Performing an ANALYZE(%d) on tickets table...", deepStatsTarget)
		if err = AnalyzeTable(db.db, "tickets", deepStatsTarget); err != nil {
			return nodeHeight, fmt.Errorf("failed to ANALYZE tickets table: %v", err)
		}
	}

	// Quickly ANALYZE all tables if not already done after indexing.
	if !analyzed && requireAnalyze {
		// Analyze all tables.
		log.Infof("Performing an ANALYZE(%d) on all tables...", quickStatsTarget)
		if err = AnalyzeAllTables(db.db, quickStatsTarget); err != nil {
			return nodeHeight, fmt.Errorf("failed to ANALYZE tables: %v", err)
		}
	}

	// After sync and indexing, must use upsert statement, which checks for
	// duplicate entries and updates instead of throwing and error and panicing.
	db.EnableDuplicateCheckOnInsert(true)

	// Set meta.ibd_complete = TRUE.
	if err = SetIBDComplete(db.db, true); err != nil {
		return nodeHeight, fmt.Errorf("failed to set meta.ibd_complete: %v", err)
	}

	if barLoad != nil {
		barID := dbtypes.InitialDBLoad
		if updateAllAddresses {
			barID = dbtypes.AddressesTableSync
		}
		sendProgressUpdate(&dbtypes.ProgressBarLoad{
			BarID:    barID,
			Subtitle: "sync complete",
		})
	}

	log.Infof("Sync finished at height %d. Delta: %d blocks, %d transactions, %d ins, %d outs, %d addresses",
		nodeHeight, nodeHeight-startHeight+1, totalTxs, totalVins, totalVouts, totalAddresses)

	return nodeHeight, err
}
