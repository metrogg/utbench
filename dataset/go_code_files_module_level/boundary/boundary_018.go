func (pm *ProtocolManager) handleMsg(p *peer) error {
	select {
	case err := <-p.errCh:
		return err
	default:
	}
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	p.Log().Trace("Light Ethereum message arrived", "code", msg.Code, "bytes", msg.Size)

	p.responseCount++
	responseCount := p.responseCount
	var (
		maxCost uint64
		task    *servingTask
	)

	accept := func(reqID, reqCnt, maxCnt uint64) bool {
		if reqCnt == 0 {
			return false
		}
		if p.fcClient == nil || reqCnt > maxCnt {
			return false
		}
		maxCost = p.fcCosts.getCost(msg.Code, reqCnt)

		if accepted, bufShort, servingPriority := p.fcClient.AcceptRequest(reqID, responseCount, maxCost); !accepted {
			if bufShort > 0 {
				p.Log().Error("Request came too early", "remaining", common.PrettyDuration(time.Duration(bufShort*1000000/p.fcParams.MinRecharge)))
			}
			return false
		} else {
			task = pm.servingQueue.newTask(servingPriority)
		}
		return task.start()
	}

	if msg.Size > ProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	defer msg.Discard()

	var deliverMsg *Msg

	sendResponse := func(reqID, amount uint64, reply *reply, servingTime uint64) {
		p.responseLock.Lock()
		defer p.responseLock.Unlock()

		var replySize uint32
		if reply != nil {
			replySize = reply.size()
		}
		var realCost uint64
		if pm.server.costTracker != nil {
			realCost = pm.server.costTracker.realCost(servingTime, msg.Size, replySize)
			pm.server.costTracker.updateStats(msg.Code, amount, servingTime, realCost)
		} else {
			realCost = maxCost
		}
		bv := p.fcClient.RequestProcessed(reqID, responseCount, maxCost, realCost)
		if reply != nil {
			p.queueSend(func() {
				if err := reply.send(bv); err != nil {
					select {
					case p.errCh <- err:
					default:
					}
				}
			})
		}
	}

	// Handle the message depending on its contents
	switch msg.Code {
	case StatusMsg:
		p.Log().Trace("Received status message")
		// Status messages should never arrive after the handshake
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")

	// Block header query, collect the requested headers and reply
	case AnnounceMsg:
		p.Log().Trace("Received announce message")
		var req announceData
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		update, size := req.Update.decode()
		if p.rejectUpdate(size) {
			return errResp(ErrRequestRejected, "")
		}
		p.updateFlowControl(update)

		if req.Hash != (common.Hash{}) {
			if p.announceType == announceTypeNone {
				return errResp(ErrUnexpectedResponse, "")
			}
			if p.announceType == announceTypeSigned {
				if err := req.checkSignature(p.ID(), update); err != nil {
					p.Log().Trace("Invalid announcement signature", "err", err)
					return err
				}
				p.Log().Trace("Valid announcement signature")
			}

			p.Log().Trace("Announce message content", "number", req.Number, "hash", req.Hash, "td", req.Td, "reorg", req.ReorgDepth)
			if pm.fetcher != nil {
				pm.fetcher.announce(p, &req)
			}
		}

	case GetBlockHeadersMsg:
		p.Log().Trace("Received block header request")
		// Decode the complex header query
		var req struct {
			ReqID uint64
			Query getBlockHeadersData
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		query := req.Query
		if !accept(req.ReqID, query.Amount, MaxHeaderFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			hashMode := query.Origin.Hash != (common.Hash{})
			first := true
			maxNonCanonical := uint64(100)

			// Gather headers until the fetch or network limits is reached
			var (
				bytes   common.StorageSize
				headers []*types.Header
				unknown bool
			)
			for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit {
				if !first && !task.waitOrStop() {
					return
				}
				// Retrieve the next header satisfying the query
				var origin *types.Header
				if hashMode {
					if first {
						origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
						if origin != nil {
							query.Origin.Number = origin.Number.Uint64()
						}
					} else {
						origin = pm.blockchain.GetHeader(query.Origin.Hash, query.Origin.Number)
					}
				} else {
					origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
				}
				if origin == nil {
					break
				}
				headers = append(headers, origin)
				bytes += estHeaderRlpSize

				// Advance to the next header of the query
				switch {
				case hashMode && query.Reverse:
					// Hash based traversal towards the genesis block
					ancestor := query.Skip + 1
					if ancestor == 0 {
						unknown = true
					} else {
						query.Origin.Hash, query.Origin.Number = pm.blockchain.GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
						unknown = (query.Origin.Hash == common.Hash{})
					}
				case hashMode && !query.Reverse:
					// Hash based traversal towards the leaf block
					var (
						current = origin.Number.Uint64()
						next    = current + query.Skip + 1
					)
					if next <= current {
						infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
						p.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
						unknown = true
					} else {
						if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
							nextHash := header.Hash()
							expOldHash, _ := pm.blockchain.GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
							if expOldHash == query.Origin.Hash {
								query.Origin.Hash, query.Origin.Number = nextHash, next
							} else {
								unknown = true
							}
						} else {
							unknown = true
						}
					}
				case query.Reverse:
					// Number based traversal towards the genesis block
					if query.Origin.Number >= query.Skip+1 {
						query.Origin.Number -= query.Skip + 1
					} else {
						unknown = true
					}

				case !query.Reverse:
					// Number based traversal towards the leaf block
					query.Origin.Number += query.Skip + 1
				}
				first = false
			}
			sendResponse(req.ReqID, query.Amount, p.ReplyBlockHeaders(req.ReqID, headers), task.done())
		}()

	case BlockHeadersMsg:
		if pm.downloader == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received block header response message")
		// A batch of headers arrived to one of our previous requests
		var resp struct {
			ReqID, BV uint64
			Headers   []*types.Header
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		if pm.fetcher != nil && pm.fetcher.requestedID(resp.ReqID) {
			pm.fetcher.deliverHeaders(p, resp.ReqID, resp.Headers)
		} else {
			err := pm.downloader.DeliverHeaders(p.id, resp.Headers)
			if err != nil {
				log.Debug(fmt.Sprint(err))
			}
		}

	case GetBlockBodiesMsg:
		p.Log().Trace("Received block bodies request")
		// Decode the retrieval message
		var req struct {
			ReqID  uint64
			Hashes []common.Hash
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			bytes  int
			bodies []rlp.RawValue
		)
		reqCnt := len(req.Hashes)
		if !accept(req.ReqID, uint64(reqCnt), MaxBodyFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			for i, hash := range req.Hashes {
				if i != 0 && !task.waitOrStop() {
					return
				}
				if bytes >= softResponseLimit {
					break
				}
				// Retrieve the requested block body, stopping if enough was found
				if number := rawdb.ReadHeaderNumber(pm.chainDb, hash); number != nil {
					if data := rawdb.ReadBodyRLP(pm.chainDb, hash, *number); len(data) != 0 {
						bodies = append(bodies, data)
						bytes += len(data)
					}
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyBlockBodiesRLP(req.ReqID, bodies), task.done())
		}()

	case BlockBodiesMsg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received block bodies response")
		// A batch of block bodies arrived to one of our previous requests
		var resp struct {
			ReqID, BV uint64
			Data      []*types.Body
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		deliverMsg = &Msg{
			MsgType: MsgBlockBodies,
			ReqID:   resp.ReqID,
			Obj:     resp.Data,
		}

	case GetCodeMsg:
		p.Log().Trace("Received code request")
		// Decode the retrieval message
		var req struct {
			ReqID uint64
			Reqs  []CodeReq
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			bytes int
			data  [][]byte
		)
		reqCnt := len(req.Reqs)
		if !accept(req.ReqID, uint64(reqCnt), MaxCodeFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			for i, req := range req.Reqs {
				if i != 0 && !task.waitOrStop() {
					return
				}
				// Look up the root hash belonging to the request
				number := rawdb.ReadHeaderNumber(pm.chainDb, req.BHash)
				if number == nil {
					p.Log().Warn("Failed to retrieve block num for code", "hash", req.BHash)
					continue
				}
				header := rawdb.ReadHeader(pm.chainDb, req.BHash, *number)
				if header == nil {
					p.Log().Warn("Failed to retrieve header for code", "block", *number, "hash", req.BHash)
					continue
				}
				triedb := pm.blockchain.StateCache().TrieDB()

				account, err := pm.getAccount(triedb, header.Root, common.BytesToHash(req.AccKey))
				if err != nil {
					p.Log().Warn("Failed to retrieve account for code", "block", header.Number, "hash", header.Hash(), "account", common.BytesToHash(req.AccKey), "err", err)
					continue
				}
				code, err := triedb.Node(common.BytesToHash(account.CodeHash))
				if err != nil {
					p.Log().Warn("Failed to retrieve account code", "block", header.Number, "hash", header.Hash(), "account", common.BytesToHash(req.AccKey), "codehash", common.BytesToHash(account.CodeHash), "err", err)
					continue
				}
				// Accumulate the code and abort if enough data was retrieved
				data = append(data, code)
				if bytes += len(code); bytes >= softResponseLimit {
					break
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyCode(req.ReqID, data), task.done())
		}()

	case CodeMsg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received code response")
		// A batch of node state data arrived to one of our previous requests
		var resp struct {
			ReqID, BV uint64
			Data      [][]byte
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		deliverMsg = &Msg{
			MsgType: MsgCode,
			ReqID:   resp.ReqID,
			Obj:     resp.Data,
		}

	case GetReceiptsMsg:
		p.Log().Trace("Received receipts request")
		// Decode the retrieval message
		var req struct {
			ReqID  uint64
			Hashes []common.Hash
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			bytes    int
			receipts []rlp.RawValue
		)
		reqCnt := len(req.Hashes)
		if !accept(req.ReqID, uint64(reqCnt), MaxReceiptFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			for i, hash := range req.Hashes {
				if i != 0 && !task.waitOrStop() {
					return
				}
				if bytes >= softResponseLimit {
					break
				}
				// Retrieve the requested block's receipts, skipping if unknown to us
				var results types.Receipts
				if number := rawdb.ReadHeaderNumber(pm.chainDb, hash); number != nil {
					results = rawdb.ReadRawReceipts(pm.chainDb, hash, *number)
				}
				if results == nil {
					if header := pm.blockchain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
						continue
					}
				}
				// If known, encode and queue for response packet
				if encoded, err := rlp.EncodeToBytes(results); err != nil {
					log.Error("Failed to encode receipt", "err", err)
				} else {
					receipts = append(receipts, encoded)
					bytes += len(encoded)
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyReceiptsRLP(req.ReqID, receipts), task.done())
		}()

	case ReceiptsMsg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received receipts response")
		// A batch of receipts arrived to one of our previous requests
		var resp struct {
			ReqID, BV uint64
			Receipts  []types.Receipts
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		deliverMsg = &Msg{
			MsgType: MsgReceipts,
			ReqID:   resp.ReqID,
			Obj:     resp.Receipts,
		}

	case GetProofsV2Msg:
		p.Log().Trace("Received les/2 proofs request")
		// Decode the retrieval message
		var req struct {
			ReqID uint64
			Reqs  []ProofReq
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			lastBHash common.Hash
			root      common.Hash
		)
		reqCnt := len(req.Reqs)
		if !accept(req.ReqID, uint64(reqCnt), MaxProofsFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			nodes := light.NewNodeSet()

			for i, req := range req.Reqs {
				if i != 0 && !task.waitOrStop() {
					return
				}
				// Look up the root hash belonging to the request
				var (
					number *uint64
					header *types.Header
					trie   state.Trie
				)
				if req.BHash != lastBHash {
					root, lastBHash = common.Hash{}, req.BHash

					if number = rawdb.ReadHeaderNumber(pm.chainDb, req.BHash); number == nil {
						p.Log().Warn("Failed to retrieve block num for proof", "hash", req.BHash)
						continue
					}
					if header = rawdb.ReadHeader(pm.chainDb, req.BHash, *number); header == nil {
						p.Log().Warn("Failed to retrieve header for proof", "block", *number, "hash", req.BHash)
						continue
					}
					root = header.Root
				}
				// Open the account or storage trie for the request
				statedb := pm.blockchain.StateCache()

				switch len(req.AccKey) {
				case 0:
					// No account key specified, open an account trie
					trie, err = statedb.OpenTrie(root)
					if trie == nil || err != nil {
						p.Log().Warn("Failed to open storage trie for proof", "block", header.Number, "hash", header.Hash(), "root", root, "err", err)
						continue
					}
				default:
					// Account key specified, open a storage trie
					account, err := pm.getAccount(statedb.TrieDB(), root, common.BytesToHash(req.AccKey))
					if err != nil {
						p.Log().Warn("Failed to retrieve account for proof", "block", header.Number, "hash", header.Hash(), "account", common.BytesToHash(req.AccKey), "err", err)
						continue
					}
					trie, err = statedb.OpenStorageTrie(common.BytesToHash(req.AccKey), account.Root)
					if trie == nil || err != nil {
						p.Log().Warn("Failed to open storage trie for proof", "block", header.Number, "hash", header.Hash(), "account", common.BytesToHash(req.AccKey), "root", account.Root, "err", err)
						continue
					}
				}
				// Prove the user's request from the account or stroage trie
				if err := trie.Prove(req.Key, req.FromLevel, nodes); err != nil {
					p.Log().Warn("Failed to prove state request", "block", header.Number, "hash", header.Hash(), "err", err)
					continue
				}
				if nodes.DataSize() >= softResponseLimit {
					break
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyProofsV2(req.ReqID, nodes.NodeList()), task.done())
		}()

	case ProofsV2Msg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received les/2 proofs response")
		// A batch of merkle proofs arrived to one of our previous requests
		var resp struct {
			ReqID, BV uint64
			Data      light.NodeList
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		deliverMsg = &Msg{
			MsgType: MsgProofsV2,
			ReqID:   resp.ReqID,
			Obj:     resp.Data,
		}

	case GetHelperTrieProofsMsg:
		p.Log().Trace("Received helper trie proof request")
		// Decode the retrieval message
		var req struct {
			ReqID uint64
			Reqs  []HelperTrieReq
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Gather state data until the fetch or network limits is reached
		var (
			auxBytes int
			auxData  [][]byte
		)
		reqCnt := len(req.Reqs)
		if !accept(req.ReqID, uint64(reqCnt), MaxHelperTrieProofsFetch) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {

			var (
				lastIdx  uint64
				lastType uint
				root     common.Hash
				auxTrie  *trie.Trie
			)
			nodes := light.NewNodeSet()
			for i, req := range req.Reqs {
				if i != 0 && !task.waitOrStop() {
					return
				}
				if auxTrie == nil || req.Type != lastType || req.TrieIdx != lastIdx {
					auxTrie, lastType, lastIdx = nil, req.Type, req.TrieIdx

					var prefix string
					if root, prefix = pm.getHelperTrie(req.Type, req.TrieIdx); root != (common.Hash{}) {
						auxTrie, _ = trie.New(root, trie.NewDatabase(rawdb.NewTable(pm.chainDb, prefix)))
					}
				}
				if req.AuxReq == auxRoot {
					var data []byte
					if root != (common.Hash{}) {
						data = root[:]
					}
					auxData = append(auxData, data)
					auxBytes += len(data)
				} else {
					if auxTrie != nil {
						auxTrie.Prove(req.Key, req.FromLevel, nodes)
					}
					if req.AuxReq != 0 {
						data := pm.getHelperTrieAuxData(req)
						auxData = append(auxData, data)
						auxBytes += len(data)
					}
				}
				if nodes.DataSize()+auxBytes >= softResponseLimit {
					break
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyHelperTrieProofs(req.ReqID, HelperTrieResps{Proofs: nodes.NodeList(), AuxData: auxData}), task.done())
		}()

	case HelperTrieProofsMsg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received helper trie proof response")
		var resp struct {
			ReqID, BV uint64
			Data      HelperTrieResps
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)
		deliverMsg = &Msg{
			MsgType: MsgHelperTrieProofs,
			ReqID:   resp.ReqID,
			Obj:     resp.Data,
		}

	case SendTxV2Msg:
		if pm.txpool == nil {
			return errResp(ErrRequestRejected, "")
		}
		// Transactions arrived, parse all of them and deliver to the pool
		var req struct {
			ReqID uint64
			Txs   []*types.Transaction
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		reqCnt := len(req.Txs)
		if !accept(req.ReqID, uint64(reqCnt), MaxTxSend) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			stats := make([]txStatus, len(req.Txs))
			for i, tx := range req.Txs {
				if i != 0 && !task.waitOrStop() {
					return
				}
				hash := tx.Hash()
				stats[i] = pm.txStatus(hash)
				if stats[i].Status == core.TxStatusUnknown {
					if errs := pm.txpool.AddRemotes([]*types.Transaction{tx}); errs[0] != nil {
						stats[i].Error = errs[0].Error()
						continue
					}
					stats[i] = pm.txStatus(hash)
				}
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyTxStatus(req.ReqID, stats), task.done())
		}()

	case GetTxStatusMsg:
		if pm.txpool == nil {
			return errResp(ErrUnexpectedResponse, "")
		}
		// Transactions arrived, parse all of them and deliver to the pool
		var req struct {
			ReqID  uint64
			Hashes []common.Hash
		}
		if err := msg.Decode(&req); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}
		reqCnt := len(req.Hashes)
		if !accept(req.ReqID, uint64(reqCnt), MaxTxStatus) {
			return errResp(ErrRequestRejected, "")
		}
		go func() {
			stats := make([]txStatus, len(req.Hashes))
			for i, hash := range req.Hashes {
				if i != 0 && !task.waitOrStop() {
					return
				}
				stats[i] = pm.txStatus(hash)
			}
			sendResponse(req.ReqID, uint64(reqCnt), p.ReplyTxStatus(req.ReqID, stats), task.done())
		}()

	case TxStatusMsg:
		if pm.odr == nil {
			return errResp(ErrUnexpectedResponse, "")
		}

		p.Log().Trace("Received tx status response")
		var resp struct {
			ReqID, BV uint64
			Status    []txStatus
		}
		if err := msg.Decode(&resp); err != nil {
			return errResp(ErrDecode, "msg %v: %v", msg, err)
		}

		p.fcServer.ReceivedReply(resp.ReqID, resp.BV)

	default:
		p.Log().Trace("Received unknown message", "code", msg.Code)
		return errResp(ErrInvalidMsgCode, "%v", msg.Code)
	}

	if deliverMsg != nil {
		err := pm.retriever.deliver(p, deliverMsg)
		if err != nil {
			p.responseErrors++
			if p.responseErrors > maxResponseErrors {
				return err
			}
		}
	}
	return nil
}
