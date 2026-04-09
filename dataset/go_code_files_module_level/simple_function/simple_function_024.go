func (exp *explorerUI) TxPage(w http.ResponseWriter, r *http.Request) {
	// attempt to get tx hash string from URL path
	hash, ok := r.Context().Value(ctxTxHash).(string)
	if !ok {
		log.Trace("txid not set")
		exp.StatusPage(w, defaultErrorCode, "there was no transaction requested",
			"", ExpStatusNotFound)
		return
	}

	inout, _ := r.Context().Value(ctxTxInOut).(string)
	if inout != "in" && inout != "out" && inout != "" {
		exp.StatusPage(w, defaultErrorCode, "there was no transaction requested",
			"", ExpStatusNotFound)
		return
	}
	ioid, _ := r.Context().Value(ctxTxInOutId).(string)
	inoutid, _ := strconv.ParseInt(ioid, 10, 0)

	tx := exp.blockData.GetExplorerTx(hash)
	// If dcrd has no information about the transaction, pull the transaction
	// details from the auxiliary DB database.
	if tx == nil {
		// Search for occurrences of the transaction in the database.
		dbTxs, err := exp.explorerSource.Transaction(hash)
		if exp.timeoutErrorPage(w, err, "Transaction") {
			return
		}
		if err != nil {
			log.Errorf("Unable to retrieve transaction details for %s.", hash)
			exp.StatusPage(w, defaultErrorCode, "could not find that transaction",
				"", ExpStatusNotFound)
			return
		}
		if dbTxs == nil {
			exp.StatusPage(w, defaultErrorCode, "that transaction has not been recorded",
				"", ExpStatusNotFound)
			return
		}

		// Take the first one. The query order should put valid at the top of
		// the list. Regardless of order, the transaction web page will link to
		// all occurrences of the transaction.
		dbTx0 := dbTxs[0]
		fees := dcrutil.Amount(dbTx0.Fees)
		tx = &types.TxInfo{
			TxBasic: &types.TxBasic{
				TxID:          hash,
				FormattedSize: humanize.Bytes(uint64(dbTx0.Size)),
				Total:         dcrutil.Amount(dbTx0.Sent).ToCoin(),
				Fee:           fees,
				FeeRate:       dcrutil.Amount((1000 * int64(fees)) / int64(dbTx0.Size)),
				// VoteInfo TODO - check votes table
				Coinbase: dbTx0.BlockIndex == 0,
			},
			SpendingTxns: make([]types.TxInID, len(dbTx0.VoutDbIds)), // SpendingTxns filled below
			Type:         txhelpers.TxTypeToString(int(dbTx0.TxType)),
			// Vins - looked-up in vins table
			// Vouts - looked-up in vouts table
			BlockHeight:   dbTx0.BlockHeight,
			BlockIndex:    dbTx0.BlockIndex,
			BlockHash:     dbTx0.BlockHash,
			Confirmations: exp.Height() - dbTx0.BlockHeight + 1,
			Time:          types.TimeDef(dbTx0.Time),
		}

		// Coinbase transactions are regular, but call them coinbase for the page.
		if tx.Coinbase {
			tx.Type = "Coinbase"
		}

		// Retrieve vouts from DB.
		vouts, err := exp.explorerSource.VoutsForTx(dbTx0)
		if exp.timeoutErrorPage(w, err, "VoutsForTx") {
			return
		}
		if err != nil {
			log.Errorf("Failed to retrieve all vout details for transaction %s: %v",
				dbTx0.TxID, err)
			exp.StatusPage(w, defaultErrorCode, "VoutsForTx failed", "", ExpStatusError)
			return
		}

		// Convert to explorer.Vout, getting spending information from DB.
		for iv := range vouts {
			// Check pkScript for OP_RETURN
			var opReturn string
			asm, _ := txscript.DisasmString(vouts[iv].ScriptPubKey)
			if strings.Contains(asm, "OP_RETURN") {
				opReturn = asm
			}
			// Determine if the outpoint is spent
			spendingTx, _, _, err := exp.explorerSource.SpendingTransaction(hash, vouts[iv].TxIndex)
			if exp.timeoutErrorPage(w, err, "SpendingTransaction") {
				return
			}
			if err != nil && err != sql.ErrNoRows {
				log.Warnf("SpendingTransaction failed for outpoint %s:%d: %v",
					hash, vouts[iv].TxIndex, err)
			}
			amount := dcrutil.Amount(int64(vouts[iv].Value)).ToCoin()
			tx.Vout = append(tx.Vout, types.Vout{
				Addresses:       vouts[iv].ScriptPubKeyData.Addresses,
				Amount:          amount,
				FormattedAmount: humanize.Commaf(amount),
				Type:            txhelpers.TxTypeToString(int(vouts[iv].TxType)),
				Spent:           spendingTx != "",
				OP_RETURN:       opReturn,
				Index:           vouts[iv].TxIndex,
			})
		}

		// Retrieve vins from DB.
		vins, prevPkScripts, scriptVersions, err := exp.explorerSource.VinsForTx(dbTx0)
		if exp.timeoutErrorPage(w, err, "VinsForTx") {
			return
		}
		if err != nil {
			log.Errorf("Failed to retrieve all vin details for transaction %s: %v",
				dbTx0.TxID, err)
			exp.StatusPage(w, defaultErrorCode, "VinsForTx failed", "", ExpStatusError)
			return
		}

		// Convert to explorer.Vin from dbtypes.VinTxProperty.
		for iv := range vins {
			// Decode all addresses from previous outpoint's pkScript.
			var addresses []string
			pkScriptsStr, err := hex.DecodeString(prevPkScripts[iv])
			if err != nil {
				log.Errorf("Failed to decode pkgScript: %v", err)
			}
			_, scrAddrs, _, err := txscript.ExtractPkScriptAddrs(scriptVersions[iv],
				pkScriptsStr, exp.ChainParams)
			if err != nil {
				log.Errorf("Failed to decode pkScript: %v", err)
			} else {
				for ia := range scrAddrs {
					addresses = append(addresses, scrAddrs[ia].EncodeAddress())
				}
			}

			// If the scriptsig does not decode or disassemble, oh well.
			asm, _ := txscript.DisasmString(vins[iv].ScriptHex)

			txIndex := vins[iv].TxIndex
			amount := dcrutil.Amount(vins[iv].ValueIn).ToCoin()
			var coinbase, stakebase string
			if txIndex == 0 {
				if tx.Coinbase {
					coinbase = hex.EncodeToString(txhelpers.CoinbaseScript)
				} else if tx.IsVote() {
					stakebase = hex.EncodeToString(txhelpers.CoinbaseScript)
				}
			}
			tx.Vin = append(tx.Vin, types.Vin{
				Vin: &dcrjson.Vin{
					Coinbase:    coinbase,
					Stakebase:   stakebase,
					Txid:        hash,
					Vout:        vins[iv].PrevTxIndex,
					Tree:        dbTx0.Tree,
					Sequence:    vins[iv].Sequence,
					AmountIn:    amount,
					BlockHeight: uint32(tx.BlockHeight),
					BlockIndex:  tx.BlockIndex,
					ScriptSig: &dcrjson.ScriptSig{
						Asm: asm,
						Hex: hex.EncodeToString(vins[iv].ScriptHex),
					},
				},
				Addresses:       addresses,
				FormattedAmount: humanize.Commaf(amount),
				Index:           txIndex,
			})
		}

		// For coinbase and stakebase, get maturity status.
		if tx.Coinbase || tx.IsVote() {
			tx.Maturity = int64(exp.ChainParams.CoinbaseMaturity)
			if tx.IsVote() {
				tx.Maturity++ // TODO why as elsewhere for votes?
			}
			if tx.Confirmations >= int64(exp.ChainParams.CoinbaseMaturity) {
				tx.Mature = "True"
			} else if tx.IsVote() {
				tx.VoteFundsLocked = "True"
			}
			coinbaseMaturityInHours :=
				exp.ChainParams.TargetTimePerBlock.Hours() * float64(tx.Maturity)
			tx.MaturityTimeTill = coinbaseMaturityInHours *
				(1 - float64(tx.Confirmations)/float64(tx.Maturity))
		}

		// For ticket purchase, get status and maturity blocks, but compute
		// details in normal code branch below.
		if tx.IsTicket() {
			tx.TicketInfo.TicketMaturity = int64(exp.ChainParams.TicketMaturity)
			if tx.Confirmations >= tx.TicketInfo.TicketMaturity {
				tx.Mature = "True"
			}
		}
	} // tx == nil (not found by dcrd)

	// Check for any transaction outputs that appear unspent.
	unspents := types.UnspentOutputIndices(tx.Vout)
	if len(unspents) > 0 {
		// Grab the mempool transaction inputs that match this transaction.
		mempoolVins := exp.GetTxMempoolInputs(hash, tx.Type)
		if len(mempoolVins) > 0 {
			// A quick matching function.
			matchingVin := func(vout *types.Vout) (string, uint32) {
				for vindex := range mempoolVins {
					vin := mempoolVins[vindex]
					for inIdx := range vin.Inputs {
						input := vin.Inputs[inIdx]
						if input.Outdex == vout.Index {
							return vin.TxId, input.Index
						}
					}
				}
				return "", 0
			}
			for _, outdex := range unspents {
				vout := &tx.Vout[outdex]
				txid, vindex := matchingVin(vout)
				if txid == "" {
					continue
				}
				vout.Spent = true
				tx.SpendingTxns[vout.Index] = types.TxInID{
					Hash:  txid,
					Index: vindex,
				}
			}
		}
	}

	// Set ticket-related parameters.
	if tx.IsTicket() {
		blocksLive := tx.Confirmations - int64(exp.ChainParams.TicketMaturity)
		tx.TicketInfo.TicketPoolSize = int64(exp.ChainParams.TicketPoolSize) *
			int64(exp.ChainParams.TicketsPerBlock)
		tx.TicketInfo.TicketExpiry = int64(exp.ChainParams.TicketExpiry)
		expirationInDays := (exp.ChainParams.TargetTimePerBlock.Hours() *
			float64(exp.ChainParams.TicketExpiry)) / 24
		maturityInHours := (exp.ChainParams.TargetTimePerBlock.Hours() *
			float64(tx.TicketInfo.TicketMaturity))
		tx.TicketInfo.TimeTillMaturity = ((float64(exp.ChainParams.TicketMaturity) -
			float64(tx.Confirmations)) / float64(exp.ChainParams.TicketMaturity)) *
			maturityInHours
		ticketExpiryBlocksLeft := int64(exp.ChainParams.TicketExpiry) - blocksLive
		tx.TicketInfo.TicketExpiryDaysLeft = (float64(ticketExpiryBlocksLeft) /
			float64(exp.ChainParams.TicketExpiry)) * expirationInDays
	}

	// For any coinbase transactions look up the total block fees to include
	// as part of the inputs.
	if tx.Type == "Coinbase" {
		data := exp.blockData.GetExplorerBlock(tx.BlockHash)
		if data == nil {
			log.Errorf("Unable to get block %s", tx.BlockHash)
		} else {
			// BlockInfo.MiningFee is coin (float64), while
			// TxInfo.BlockMiningFee is int64 (atoms), so convert. If the
			// float64 is somehow invalid, use the default zero value.
			feeAmt, _ := dcrutil.NewAmount(data.MiningFee)
			tx.BlockMiningFee = int64(feeAmt)
		}
	}

	// Details on all the blocks containing this transaction
	blocks, blockInds, err := exp.explorerSource.TransactionBlocks(tx.TxID)
	if exp.timeoutErrorPage(w, err, "TransactionBlocks") {
		return
	}
	if err != nil {
		log.Errorf("Unable to retrieve blocks for transaction %s: %v",
			hash, err)
		exp.StatusPage(w, defaultErrorCode, defaultErrorMessage, tx.TxID, ExpStatusError)
		return
	}

	// See if any of these blocks are mainchain and stakeholder-approved
	// (a.k.a. valid).
	var isConfirmedMainchain bool
	for ib := range blocks {
		if blocks[ib].IsValid && blocks[ib].IsMainchain {
			isConfirmedMainchain = true
			break
		}
	}

	// For each output of this transaction, look up any spending transactions,
	// and the index of the spending transaction input.
	spendingTxHashes, spendingTxVinInds, voutInds, err :=
		exp.explorerSource.SpendingTransactions(hash)
	if exp.timeoutErrorPage(w, err, "SpendingTransactions") {
		return
	}
	if err != nil {
		log.Errorf("Unable to retrieve spending transactions for %s: %v", hash, err)
		exp.StatusPage(w, defaultErrorCode, defaultErrorMessage, hash, ExpStatusError)
		return
	}
	for i, vout := range voutInds {
		if int(vout) >= len(tx.SpendingTxns) {
			log.Errorf("Invalid spending transaction data (%s:%d)", hash, vout)
			continue
		}
		tx.SpendingTxns[vout] = types.TxInID{
			Hash:  spendingTxHashes[i],
			Index: spendingTxVinInds[i],
		}
	}

	if tx.IsTicket() {
		spendStatus, poolStatus, err := exp.explorerSource.PoolStatusForTicket(hash)
		if exp.timeoutErrorPage(w, err, "PoolStatusForTicket") {
			return
		}
		if err != nil && err != sql.ErrNoRows {
			log.Errorf("Unable to retrieve ticket spend and pool status for %s: %v",
				hash, err)
			exp.StatusPage(w, defaultErrorCode, defaultErrorMessage, "", ExpStatusError)
			return
		} else if err == sql.ErrNoRows {
			if tx.Confirmations != 0 {
				log.Warnf("Spend and pool status not found for ticket %s: %v", hash, err)
			}
		} else {
			if tx.Mature == "False" {
				tx.TicketInfo.PoolStatus = "immature"
			} else {
				tx.TicketInfo.PoolStatus = poolStatus.String()
			}
			tx.TicketInfo.SpendStatus = spendStatus.String()

			// For missed tickets, get the block in which it should have voted.
			if poolStatus == dbtypes.PoolStatusMissed {
				tx.TicketInfo.LotteryBlock, _, err = exp.explorerSource.TicketMiss(hash)
				if err != nil && err != sql.ErrNoRows {
					log.Errorf("Unable to retrieve miss information for ticket %s: %v",
						hash, err)
					exp.StatusPage(w, defaultErrorCode, defaultErrorMessage, "", ExpStatusError)
					return
				} else if err == sql.ErrNoRows {
					log.Warnf("No mainchain miss data for ticket %s: %v",
						hash, err)
				}
			}

			// Ticket luck and probability of voting.
			// blockLive < 0 for immature tickets
			blocksLive := tx.Confirmations - int64(exp.ChainParams.TicketMaturity)
			if tx.TicketInfo.SpendStatus == "Voted" {
				// Blocks from eligible until voted (actual luck)
				txhash, err := chainhash.NewHashFromStr(tx.SpendingTxns[0].Hash)
				if err != nil {
					exp.StatusPage(w, defaultErrorCode, err.Error(), "", ExpStatusError)
					return
				}
				tx.TicketInfo.TicketLiveBlocks = exp.blockData.TxHeight(txhash) -
					tx.BlockHeight - int64(exp.ChainParams.TicketMaturity) - 1
			} else if tx.Confirmations >= int64(exp.ChainParams.TicketExpiry+
				uint32(exp.ChainParams.TicketMaturity)) { // Expired
				// Blocks ticket was active before expiring (actual no luck)
				tx.TicketInfo.TicketLiveBlocks = int64(exp.ChainParams.TicketExpiry)
			} else { // Active
				// Blocks ticket has been active and eligible to vote
				tx.TicketInfo.TicketLiveBlocks = blocksLive
			}
			tx.TicketInfo.BestLuck = tx.TicketInfo.TicketExpiry / int64(exp.ChainParams.TicketPoolSize)
			tx.TicketInfo.AvgLuck = tx.TicketInfo.BestLuck - 1
			if tx.TicketInfo.TicketLiveBlocks == int64(exp.ChainParams.TicketExpiry) {
				tx.TicketInfo.VoteLuck = 0
			} else {
				tx.TicketInfo.VoteLuck = float64(tx.TicketInfo.BestLuck) -
					(float64(tx.TicketInfo.TicketLiveBlocks) / float64(exp.ChainParams.TicketPoolSize))
			}
			if tx.TicketInfo.VoteLuck >= float64(tx.TicketInfo.BestLuck-
				(1/int64(exp.ChainParams.TicketPoolSize))) {
				tx.TicketInfo.LuckStatus = "Perfection"
			} else if tx.TicketInfo.VoteLuck > (float64(tx.TicketInfo.BestLuck) - 0.25) {
				tx.TicketInfo.LuckStatus = "Very Lucky!"
			} else if tx.TicketInfo.VoteLuck > (float64(tx.TicketInfo.BestLuck) - 0.75) {
				tx.TicketInfo.LuckStatus = "Good Luck"
			} else if tx.TicketInfo.VoteLuck > (float64(tx.TicketInfo.BestLuck) - 1.25) {
				tx.TicketInfo.LuckStatus = "Normal"
			} else if tx.TicketInfo.VoteLuck > (float64(tx.TicketInfo.BestLuck) * 0.50) {
				tx.TicketInfo.LuckStatus = "Bad Luck"
			} else if tx.TicketInfo.VoteLuck > 0 {
				tx.TicketInfo.LuckStatus = "Horrible Luck!"
			} else if tx.TicketInfo.VoteLuck == 0 {
				tx.TicketInfo.LuckStatus = "No Luck"
			}

			// Chance for a ticket to NOT be voted in a given time frame:
			// C = (1 - P)^N
			// Where: P is the probability of a vote in one block. (votes
			// per block / current ticket pool size)
			// N is the number of blocks before ticket expiry. (ticket
			// expiry in blocks - (number of blocks since ticket purchase -
			// ticket maturity))
			// C is the probability (chance)
			exp.pageData.RLock()
			pVote := float64(exp.ChainParams.TicketsPerBlock) /
				float64(exp.pageData.HomeInfo.PoolInfo.Size)
			exp.pageData.RUnlock()

			remainingBlocksLive := float64(exp.ChainParams.TicketExpiry) -
				float64(blocksLive)
			tx.TicketInfo.Probability = 100 * math.Pow(1-pVote, remainingBlocksLive)
		}
	} // tx.IsTicket()

	// Prepare the string to display for previous outpoint.
	for idx := range tx.Vin {
		vin := &tx.Vin[idx]
		if vin.Coinbase != "" {
			vin.DisplayText = "Coinbase"
		} else if vin.Stakebase != "" {
			vin.DisplayText = "Stakebase"
		} else {
			voutStr := strconv.Itoa(int(vin.Vout))
			vin.DisplayText = vin.Txid + ":" + voutStr
			vin.TextIsHash = true
			vin.Link = "/tx/" + vin.Txid + "/out/" + voutStr
		}
	}

	// For an unconfirmed tx, get the time it was received in explorer's mempool.
	if tx.BlockHeight == 0 {
		tx.Time = exp.mempoolTime(tx.TxID)
	}

	pageData := struct {
		*CommonPageData
		Data                 *types.TxInfo
		Blocks               []*dbtypes.BlockStatus
		BlockInds            []uint32
		IsConfirmedMainchain bool
		HighlightInOut       string
		HighlightInOutID     int64
		Conversions          struct {
			Total *exchanges.Conversion
			Fees  *exchanges.Conversion
		}
	}{
		CommonPageData:       exp.commonData(r),
		Data:                 tx,
		Blocks:               blocks,
		BlockInds:            blockInds,
		IsConfirmedMainchain: isConfirmedMainchain,
		HighlightInOut:       inout,
		HighlightInOutID:     inoutid,
	}

	// Get a fiat-converted value for the total and the fees.
	if exp.xcBot != nil {
		pageData.Conversions.Total = exp.xcBot.Conversion(tx.Total)
		pageData.Conversions.Fees = exp.xcBot.Conversion(tx.Fee.ToCoin())
	}

	str, err := exp.templates.execTemplateToString("tx", pageData)
	if err != nil {
		log.Errorf("Template execute failure: %v", err)
		exp.StatusPage(w, defaultErrorCode, defaultErrorMessage, "", ExpStatusError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Turbolinks-Location", r.URL.RequestURI())
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, str)
}
