func (bls *Streamer) parseEvents(ctx context.Context, events <-chan mysql.BinlogEvent) (mysql.Position, error) {
	var statements []FullBinlogStatement
	var format mysql.BinlogFormat
	var gtid mysql.GTID
	var pos = bls.startPos
	var autocommit = true
	var err error

	// Remember the RBR state.
	// tableMaps is indexed by tableID.
	tableMaps := make(map[uint64]*tableCacheEntry)

	// A begin can be triggered either by a BEGIN query, or by a GTID_EVENT.
	begin := func() {
		if statements != nil {
			// If this happened, it would be a legitimate error.
			log.Errorf("BEGIN in binlog stream while still in another transaction; dropping %d statements: %v", len(statements), statements)
			binlogStreamerErrors.Add("ParseEvents", 1)
		}
		statements = make([]FullBinlogStatement, 0, 10)
		autocommit = false
	}
	// A commit can be triggered either by a COMMIT query, or by an XID_EVENT.
	// Statements that aren't wrapped in BEGIN/COMMIT are committed immediately.
	commit := func(timestamp uint32) error {
		if int64(timestamp) >= bls.timestamp {
			eventToken := &querypb.EventToken{
				Timestamp: int64(timestamp),
				Position:  mysql.EncodePosition(pos),
			}
			if err = bls.sendTransaction(eventToken, statements); err != nil {
				if err == io.EOF {
					return ErrClientEOF
				}
				return fmt.Errorf("send reply error: %v", err)
			}
		}
		statements = nil
		autocommit = true
		return nil
	}

	// Parse events.
	for {
		var ev mysql.BinlogEvent
		var ok bool

		select {
		case ev, ok = <-events:
			if !ok {
				// events channel has been closed, which means the connection died.
				log.Infof("reached end of binlog event stream")
				return pos, ErrServerEOF
			}
		case <-ctx.Done():
			log.Infof("stopping early due to binlog Streamer service shutdown or client disconnect")
			return pos, ctx.Err()
		}

		// Validate the buffer before reading fields from it.
		if !ev.IsValid() {
			return pos, fmt.Errorf("can't parse binlog event, invalid data: %#v", ev)
		}

		// We need to keep checking for FORMAT_DESCRIPTION_EVENT even after we've
		// seen one, because another one might come along (e.g. on log rotate due to
		// binlog settings change) that changes the format.
		if ev.IsFormatDescription() {
			format, err = ev.Format()
			if err != nil {
				return pos, fmt.Errorf("can't parse FORMAT_DESCRIPTION_EVENT: %v, event data: %#v", err, ev)
			}
			continue
		}

		// We can't parse anything until we get a FORMAT_DESCRIPTION_EVENT that
		// tells us the size of the event header.
		if format.IsZero() {
			// The only thing that should come before the FORMAT_DESCRIPTION_EVENT
			// is a fake ROTATE_EVENT, which the master sends to tell us the name
			// of the current log file.
			if ev.IsRotate() {
				continue
			}
			return pos, fmt.Errorf("got a real event before FORMAT_DESCRIPTION_EVENT: %#v", ev)
		}

		// Strip the checksum, if any. We don't actually verify the checksum, so discard it.
		ev, _, err = ev.StripChecksum(format)
		if err != nil {
			return pos, fmt.Errorf("can't strip checksum from binlog event: %v, event data: %#v", err, ev)
		}

		switch {
		case ev.IsPseudo():
			gtid, _, err = ev.GTID(format)
			if err != nil {
				return pos, fmt.Errorf("can't get GTID from binlog event: %v, event data: %#v", err, ev)
			}
			oldpos := pos
			pos = mysql.AppendGTID(pos, gtid)
			// If the event is received outside of a transaction, it must
			// be sent. Otherwise, it will get lost and the targets will go out
			// of sync.
			if autocommit && !pos.Equal(oldpos) {
				if err = commit(ev.Timestamp()); err != nil {
					return pos, err
				}
			}
		case ev.IsGTID(): // GTID_EVENT: update current GTID, maybe BEGIN.
			var hasBegin bool
			gtid, hasBegin, err = ev.GTID(format)
			if err != nil {
				return pos, fmt.Errorf("can't get GTID from binlog event: %v, event data: %#v", err, ev)
			}
			pos = mysql.AppendGTID(pos, gtid)
			if hasBegin {
				begin()
			}
		case ev.IsXID(): // XID_EVENT (equivalent to COMMIT)
			if err = commit(ev.Timestamp()); err != nil {
				return pos, err
			}
		case ev.IsIntVar(): // INTVAR_EVENT
			typ, value, err := ev.IntVar(format)
			if err != nil {
				return pos, fmt.Errorf("can't parse INTVAR_EVENT: %v, event data: %#v", err, ev)
			}
			statements = append(statements, FullBinlogStatement{
				Statement: &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET %s=%d", mysql.IntVarNames[typ], value)),
				},
			})
		case ev.IsRand(): // RAND_EVENT
			seed1, seed2, err := ev.Rand(format)
			if err != nil {
				return pos, fmt.Errorf("can't parse RAND_EVENT: %v, event data: %#v", err, ev)
			}
			statements = append(statements, FullBinlogStatement{
				Statement: &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET @@RAND_SEED1=%d, @@RAND_SEED2=%d", seed1, seed2)),
				},
			})
		case ev.IsQuery(): // QUERY_EVENT
			// Extract the query string and group into transactions.
			q, err := ev.Query(format)
			if err != nil {
				return pos, fmt.Errorf("can't get query from binlog event: %v, event data: %#v", err, ev)
			}
			switch cat := getStatementCategory(q.SQL); cat {
			case binlogdatapb.BinlogTransaction_Statement_BL_BEGIN:
				begin()
			case binlogdatapb.BinlogTransaction_Statement_BL_ROLLBACK:
				// Rollbacks are possible under some circumstances. Since the stream
				// client keeps track of its replication position by updating the set
				// of GTIDs it's seen, we must commit an empty transaction so the client
				// can update its position.
				statements = nil
				fallthrough
			case binlogdatapb.BinlogTransaction_Statement_BL_COMMIT:
				if err = commit(ev.Timestamp()); err != nil {
					return pos, err
				}
			default: // BL_DDL, BL_SET, BL_INSERT, BL_UPDATE, BL_DELETE, BL_UNRECOGNIZED
				if q.Database != "" && q.Database != bls.cp.DbName {
					// Skip cross-db statements.
					continue
				}
				setTimestamp := &binlogdatapb.BinlogTransaction_Statement{
					Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
					Sql:      []byte(fmt.Sprintf("SET TIMESTAMP=%d", ev.Timestamp())),
				}
				statement := &binlogdatapb.BinlogTransaction_Statement{
					Category: cat,
					Sql:      []byte(q.SQL),
				}
				// If the statement has a charset and it's different than our client's
				// default charset, send it along with the statement.
				// If our client hasn't told us its charset, always send it.
				if bls.clientCharset == nil || (q.Charset != nil && !proto.Equal(q.Charset, bls.clientCharset)) {
					setTimestamp.Charset = q.Charset
					statement.Charset = q.Charset
				}
				statements = append(statements, FullBinlogStatement{
					Statement: setTimestamp,
				}, FullBinlogStatement{
					Statement: statement,
				})
				if autocommit {
					if err = commit(ev.Timestamp()); err != nil {
						return pos, err
					}
				}
			}
		case ev.IsPreviousGTIDs(): // PREVIOUS_GTIDS_EVENT
			// MySQL 5.6 only: The Binlogs contain an
			// event that gives us all the previously
			// applied commits. It is *not* an
			// authoritative value, unless we started from
			// the beginning of a binlog file.
			if !bls.usePreviousGTIDs {
				continue
			}
			newPos, err := ev.PreviousGTIDs(format)
			if err != nil {
				return pos, err
			}
			pos = newPos
			if err = commit(ev.Timestamp()); err != nil {
				return pos, err
			}
		case ev.IsTableMap():
			// Save all tables, even not in the same DB.
			tableID := ev.TableID(format)
			tm, err := ev.TableMap(format)
			if err != nil {
				return pos, err
			}
			// TODO(alainjobart) if table is already in map,
			// just use it.

			tce := &tableCacheEntry{
				tm:              tm,
				keyspaceIDIndex: -1,
			}
			tableMaps[tableID] = tce

			// Check we're in the right database, and if so, fill
			// in more data.
			if tm.Database != "" && tm.Database != bls.cp.DbName {
				continue
			}

			// Find and fill in the table schema.
			tce.ti = bls.se.GetTable(sqlparser.NewTableIdent(tm.Name))
			if tce.ti == nil {
				return pos, fmt.Errorf("unknown table %v in schema", tm.Name)
			}

			// Fill in the resolver if needed.
			if bls.resolverFactory != nil {
				tce.keyspaceIDIndex, tce.resolver, err = bls.resolverFactory(tce.ti)
				if err != nil {
					return pos, fmt.Errorf("cannot find column to use to find keyspace_id for table %v", tm.Name)
				}
			}

			// Fill in PK indexes if necessary.
			if bls.extractPK {
				tce.pkNames = make([]*querypb.Field, len(tce.ti.PKColumns))
				tce.pkIndexes = make([]int, len(tce.ti.Columns))
				for i := range tce.pkIndexes {
					// Put -1 as default in here.
					tce.pkIndexes[i] = -1
				}
				for i, c := range tce.ti.PKColumns {
					// Patch in every PK column index.
					tce.pkIndexes[c] = i
					// Fill in pknames
					tce.pkNames[i] = &querypb.Field{
						Name: tce.ti.Columns[c].Name.String(),
						Type: tce.ti.Columns[c].Type,
					}
				}
			}
		case ev.IsWriteRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				return pos, fmt.Errorf("unknown tableID %v in WriteRows event", tableID)
			}
			if tce.ti == nil {
				// Skip cross-db statements.
				continue
			}
			setTimestamp := &binlogdatapb.BinlogTransaction_Statement{
				Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
				Sql:      []byte(fmt.Sprintf("SET TIMESTAMP=%d", ev.Timestamp())),
			}
			statements = append(statements, FullBinlogStatement{
				Statement: setTimestamp,
			})

			rows, err := ev.Rows(format, tce.tm)
			if err != nil {
				return pos, err
			}

			statements = bls.appendInserts(statements, tce, &rows)

			if autocommit {
				if err = commit(ev.Timestamp()); err != nil {
					return pos, err
				}
			}
		case ev.IsUpdateRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				return pos, fmt.Errorf("unknown tableID %v in UpdateRows event", tableID)
			}
			if tce.ti == nil {
				// Skip cross-db statements.
				continue
			}
			setTimestamp := &binlogdatapb.BinlogTransaction_Statement{
				Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
				Sql:      []byte(fmt.Sprintf("SET TIMESTAMP=%d", ev.Timestamp())),
			}
			statements = append(statements, FullBinlogStatement{
				Statement: setTimestamp,
			})

			rows, err := ev.Rows(format, tce.tm)
			if err != nil {
				return pos, err
			}

			statements = bls.appendUpdates(statements, tce, &rows)

			if autocommit {
				if err = commit(ev.Timestamp()); err != nil {
					return pos, err
				}
			}
		case ev.IsDeleteRows():
			tableID := ev.TableID(format)
			tce, ok := tableMaps[tableID]
			if !ok {
				return pos, fmt.Errorf("unknown tableID %v in DeleteRows event", tableID)
			}
			if tce.ti == nil {
				// Skip cross-db statements.
				continue
			}
			setTimestamp := &binlogdatapb.BinlogTransaction_Statement{
				Category: binlogdatapb.BinlogTransaction_Statement_BL_SET,
				Sql:      []byte(fmt.Sprintf("SET TIMESTAMP=%d", ev.Timestamp())),
			}
			statements = append(statements, FullBinlogStatement{
				Statement: setTimestamp,
			})

			rows, err := ev.Rows(format, tce.tm)
			if err != nil {
				return pos, err
			}

			statements = bls.appendDeletes(statements, tce, &rows)

			if autocommit {
				if err = commit(ev.Timestamp()); err != nil {
					return pos, err
				}
			}
		}
	}
}
