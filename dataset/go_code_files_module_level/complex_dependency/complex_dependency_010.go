func NewStatLine(oldMongo, newMongo MongoStatus, key string, all bool, sampleSecs int64) *StatLine {
	oldStat := *oldMongo.ServerStatus
	newStat := *newMongo.ServerStatus

	returnVal := &StatLine{
		Key:       key,
		Host:      newStat.Host,
		Mapped:    -1,
		Virtual:   -1,
		Resident:  -1,
		NonMapped: -1,
		Faults:    -1,
	}

	// set connection info
	returnVal.CurrentC = newStat.Connections.Current
	returnVal.AvailableC = newStat.Connections.Available
	returnVal.TotalCreatedC = newStat.Connections.TotalCreated

	// set the storage engine appropriately
	if newStat.StorageEngine != nil && newStat.StorageEngine["name"] != "" {
		returnVal.StorageEngine = newStat.StorageEngine["name"]
	} else {
		returnVal.StorageEngine = "mmapv1"
	}

	if newStat.Opcounters != nil && oldStat.Opcounters != nil {
		returnVal.Insert, returnVal.InsertCnt = diff(newStat.Opcounters.Insert, oldStat.Opcounters.Insert, sampleSecs)
		returnVal.Query, returnVal.QueryCnt = diff(newStat.Opcounters.Query, oldStat.Opcounters.Query, sampleSecs)
		returnVal.Update, returnVal.UpdateCnt = diff(newStat.Opcounters.Update, oldStat.Opcounters.Update, sampleSecs)
		returnVal.Delete, returnVal.DeleteCnt = diff(newStat.Opcounters.Delete, oldStat.Opcounters.Delete, sampleSecs)
		returnVal.GetMore, returnVal.GetMoreCnt = diff(newStat.Opcounters.GetMore, oldStat.Opcounters.GetMore, sampleSecs)
		returnVal.Command, returnVal.CommandCnt = diff(newStat.Opcounters.Command, oldStat.Opcounters.Command, sampleSecs)
	}

	if newStat.Metrics != nil && oldStat.Metrics != nil {
		if newStat.Metrics.TTL != nil && oldStat.Metrics.TTL != nil {
			returnVal.Passes, returnVal.PassesCnt = diff(newStat.Metrics.TTL.Passes, oldStat.Metrics.TTL.Passes, sampleSecs)
			returnVal.DeletedDocuments, returnVal.DeletedDocumentsCnt = diff(newStat.Metrics.TTL.DeletedDocuments, oldStat.Metrics.TTL.DeletedDocuments, sampleSecs)
		}
		if newStat.Metrics.Cursor != nil && oldStat.Metrics.Cursor != nil {
			returnVal.TimedOutC, returnVal.TimedOutCCnt = diff(newStat.Metrics.Cursor.TimedOut, oldStat.Metrics.Cursor.TimedOut, sampleSecs)
			if newStat.Metrics.Cursor.Open != nil && oldStat.Metrics.Cursor.Open != nil {
				returnVal.NoTimeoutC, returnVal.NoTimeoutCCnt = diff(newStat.Metrics.Cursor.Open.NoTimeout, oldStat.Metrics.Cursor.Open.NoTimeout, sampleSecs)
				returnVal.PinnedC, returnVal.PinnedCCnt = diff(newStat.Metrics.Cursor.Open.Pinned, oldStat.Metrics.Cursor.Open.Pinned, sampleSecs)
				returnVal.TotalC, returnVal.TotalCCnt = diff(newStat.Metrics.Cursor.Open.Total, oldStat.Metrics.Cursor.Open.Total, sampleSecs)
			}
		}
		if newStat.Metrics.Document != nil {
			returnVal.DeletedD = newStat.Metrics.Document.Deleted
			returnVal.InsertedD = newStat.Metrics.Document.Inserted
			returnVal.ReturnedD = newStat.Metrics.Document.Returned
			returnVal.UpdatedD = newStat.Metrics.Document.Updated
		}
	}

	if newStat.OpcountersRepl != nil && oldStat.OpcountersRepl != nil {
		returnVal.InsertR, returnVal.InsertRCnt = diff(newStat.OpcountersRepl.Insert, oldStat.OpcountersRepl.Insert, sampleSecs)
		returnVal.QueryR, returnVal.QueryRCnt = diff(newStat.OpcountersRepl.Query, oldStat.OpcountersRepl.Query, sampleSecs)
		returnVal.UpdateR, returnVal.UpdateRCnt = diff(newStat.OpcountersRepl.Update, oldStat.OpcountersRepl.Update, sampleSecs)
		returnVal.DeleteR, returnVal.DeleteRCnt = diff(newStat.OpcountersRepl.Delete, oldStat.OpcountersRepl.Delete, sampleSecs)
		returnVal.GetMoreR, returnVal.GetMoreRCnt = diff(newStat.OpcountersRepl.GetMore, oldStat.OpcountersRepl.GetMore, sampleSecs)
		returnVal.CommandR, returnVal.CommandRCnt = diff(newStat.OpcountersRepl.Command, oldStat.OpcountersRepl.Command, sampleSecs)
	}

	returnVal.CacheDirtyPercent = -1
	returnVal.CacheUsedPercent = -1
	if newStat.WiredTiger != nil {
		returnVal.CacheDirtyPercent = float64(newStat.WiredTiger.Cache.TrackedDirtyBytes) / float64(newStat.WiredTiger.Cache.MaxBytesConfigured)
		returnVal.CacheUsedPercent = float64(newStat.WiredTiger.Cache.CurrentCachedBytes) / float64(newStat.WiredTiger.Cache.MaxBytesConfigured)

		returnVal.TrackedDirtyBytes = newStat.WiredTiger.Cache.TrackedDirtyBytes
		returnVal.CurrentCachedBytes = newStat.WiredTiger.Cache.CurrentCachedBytes
		returnVal.MaxBytesConfigured = newStat.WiredTiger.Cache.MaxBytesConfigured
		returnVal.AppThreadsPageReadCount = newStat.WiredTiger.Cache.AppThreadsPageReadCount
		returnVal.AppThreadsPageReadTime = newStat.WiredTiger.Cache.AppThreadsPageReadTime
		returnVal.AppThreadsPageWriteCount = newStat.WiredTiger.Cache.AppThreadsPageWriteCount
		returnVal.BytesWrittenFrom = newStat.WiredTiger.Cache.BytesWrittenFrom
		returnVal.BytesReadInto = newStat.WiredTiger.Cache.BytesReadInto
		returnVal.PagesEvictedByAppThread = newStat.WiredTiger.Cache.PagesEvictedByAppThread
		returnVal.PagesQueuedForEviction = newStat.WiredTiger.Cache.PagesQueuedForEviction
		returnVal.PagesReadIntoCache = newStat.WiredTiger.Cache.PagesReadIntoCache
		returnVal.PagesRequestedFromCache = newStat.WiredTiger.Cache.PagesRequestedFromCache
		returnVal.ServerEvictingPages = newStat.WiredTiger.Cache.ServerEvictingPages
		returnVal.WorkerThreadEvictingPages = newStat.WiredTiger.Cache.WorkerThreadEvictingPages

		returnVal.InternalPagesEvicted = newStat.WiredTiger.Cache.InternalPagesEvicted
		returnVal.ModifiedPagesEvicted = newStat.WiredTiger.Cache.ModifiedPagesEvicted
		returnVal.UnmodifiedPagesEvicted = newStat.WiredTiger.Cache.UnmodifiedPagesEvicted

		returnVal.FlushesTotalTime = newStat.WiredTiger.Transaction.TransCheckpointsTotalTimeMsecs * int64(time.Millisecond)
	}
	if newStat.WiredTiger != nil && oldStat.WiredTiger != nil {
		returnVal.Flushes, returnVal.FlushesCnt = diff(newStat.WiredTiger.Transaction.TransCheckpoints, oldStat.WiredTiger.Transaction.TransCheckpoints, sampleSecs)
	} else if newStat.BackgroundFlushing != nil && oldStat.BackgroundFlushing != nil {
		returnVal.Flushes, returnVal.FlushesCnt = diff(newStat.BackgroundFlushing.Flushes, oldStat.BackgroundFlushing.Flushes, sampleSecs)
	}

	returnVal.Time = newMongo.SampleTime
	returnVal.IsMongos =
		(newStat.ShardCursorType != nil || strings.HasPrefix(newStat.Process, MongosProcess))

	// BEGIN code modification
	if oldStat.Mem.Supported.(bool) {
		// END code modification
		if !returnVal.IsMongos {
			returnVal.Mapped = newStat.Mem.Mapped
		}
		returnVal.Virtual = newStat.Mem.Virtual
		returnVal.Resident = newStat.Mem.Resident

		if !returnVal.IsMongos && all {
			returnVal.NonMapped = newStat.Mem.Virtual - newStat.Mem.Mapped
		}
	}

	if newStat.Repl != nil {
		setName, isReplSet := newStat.Repl.SetName.(string)
		if isReplSet {
			returnVal.ReplSetName = setName
		}
		// BEGIN code modification
		if newStat.Repl.IsMaster.(bool) {
			returnVal.NodeType = "PRI"
		} else if newStat.Repl.Secondary != nil && newStat.Repl.Secondary.(bool) {
			returnVal.NodeType = "SEC"
		} else if newStat.Repl.ArbiterOnly != nil && newStat.Repl.ArbiterOnly.(bool) {
			returnVal.NodeType = "ARB"
		} else {
			returnVal.NodeType = "UNK"
		}
		// END code modification
	} else if returnVal.IsMongos {
		returnVal.NodeType = "RTR"
	}

	if oldStat.ExtraInfo != nil && newStat.ExtraInfo != nil &&
		oldStat.ExtraInfo.PageFaults != nil && newStat.ExtraInfo.PageFaults != nil {
		returnVal.Faults, returnVal.FaultsCnt = diff(*(newStat.ExtraInfo.PageFaults), *(oldStat.ExtraInfo.PageFaults), sampleSecs)
	}
	if !returnVal.IsMongos && oldStat.Locks != nil {
		globalCheck, hasGlobal := oldStat.Locks["Global"]
		if hasGlobal && globalCheck.AcquireCount != nil {
			// This appears to be a 3.0+ server so the data in these fields do *not* refer to
			// actual namespaces and thus we can't compute lock %.
			returnVal.HighestLocked = nil

			// Check if it's a 3.0+ MMAP server so we can still compute collection locks
			collectionCheck, hasCollection := oldStat.Locks["Collection"]
			if hasCollection && collectionCheck.AcquireWaitCount != nil {
				readWaitCountDiff := newStat.Locks["Collection"].AcquireWaitCount.Read - oldStat.Locks["Collection"].AcquireWaitCount.Read
				readTotalCountDiff := newStat.Locks["Collection"].AcquireCount.Read - oldStat.Locks["Collection"].AcquireCount.Read
				writeWaitCountDiff := newStat.Locks["Collection"].AcquireWaitCount.Write - oldStat.Locks["Collection"].AcquireWaitCount.Write
				writeTotalCountDiff := newStat.Locks["Collection"].AcquireCount.Write - oldStat.Locks["Collection"].AcquireCount.Write
				readAcquireTimeDiff := newStat.Locks["Collection"].TimeAcquiringMicros.Read - oldStat.Locks["Collection"].TimeAcquiringMicros.Read
				writeAcquireTimeDiff := newStat.Locks["Collection"].TimeAcquiringMicros.Write - oldStat.Locks["Collection"].TimeAcquiringMicros.Write
				returnVal.CollectionLocks = &CollectionLockStatus{
					ReadAcquireWaitsPercentage:  percentageInt64(readWaitCountDiff, readTotalCountDiff),
					WriteAcquireWaitsPercentage: percentageInt64(writeWaitCountDiff, writeTotalCountDiff),
					ReadAcquireTimeMicros:       averageInt64(readAcquireTimeDiff, readWaitCountDiff),
					WriteAcquireTimeMicros:      averageInt64(writeAcquireTimeDiff, writeWaitCountDiff),
				}
			}
		} else {
			prevLocks := parseLocks(oldStat)
			curLocks := parseLocks(newStat)
			lockdiffs := computeLockDiffs(prevLocks, curLocks)
			if len(lockdiffs) == 0 {
				if newStat.GlobalLock != nil {
					returnVal.HighestLocked = &LockStatus{
						DBName:     "",
						Percentage: percentageInt64(newStat.GlobalLock.LockTime, newStat.GlobalLock.TotalTime),
						Global:     true,
					}
				}
			} else {
				// Get the entry with the highest lock
				highestLocked := lockdiffs[len(lockdiffs)-1]

				var timeDiffMillis int64
				timeDiffMillis = newStat.UptimeMillis - oldStat.UptimeMillis

				lockToReport := highestLocked.Writes

				// if the highest locked namespace is not '.'
				if highestLocked.Namespace != "." {
					for _, namespaceLockInfo := range lockdiffs {
						if namespaceLockInfo.Namespace == "." {
							lockToReport += namespaceLockInfo.Writes
						}
					}
				}

				// lock data is in microseconds and uptime is in milliseconds - so
				// divide by 1000 so that they units match
				lockToReport /= 1000

				returnVal.HighestLocked = &LockStatus{
					DBName:     highestLocked.Namespace,
					Percentage: percentageInt64(lockToReport, timeDiffMillis),
					Global:     false,
				}
			}
		}
	} else {
		returnVal.HighestLocked = nil
	}

	if newStat.GlobalLock != nil {
		hasWT := (newStat.WiredTiger != nil && oldStat.WiredTiger != nil)
		//If we have wiredtiger stats, use those instead
		if newStat.GlobalLock.CurrentQueue != nil {
			if hasWT {
				returnVal.QueuedReaders = newStat.GlobalLock.CurrentQueue.Readers + newStat.GlobalLock.ActiveClients.Readers - newStat.WiredTiger.Concurrent.Read.Out
				returnVal.QueuedWriters = newStat.GlobalLock.CurrentQueue.Writers + newStat.GlobalLock.ActiveClients.Writers - newStat.WiredTiger.Concurrent.Write.Out
				if returnVal.QueuedReaders < 0 {
					returnVal.QueuedReaders = 0
				}
				if returnVal.QueuedWriters < 0 {
					returnVal.QueuedWriters = 0
				}
			} else {
				returnVal.QueuedReaders = newStat.GlobalLock.CurrentQueue.Readers
				returnVal.QueuedWriters = newStat.GlobalLock.CurrentQueue.Writers
			}
		}

		if hasWT {
			returnVal.ActiveReaders = newStat.WiredTiger.Concurrent.Read.Out
			returnVal.ActiveWriters = newStat.WiredTiger.Concurrent.Write.Out
		} else if newStat.GlobalLock.ActiveClients != nil {
			returnVal.ActiveReaders = newStat.GlobalLock.ActiveClients.Readers
			returnVal.ActiveWriters = newStat.GlobalLock.ActiveClients.Writers
		}
	}

	if oldStat.Network != nil && newStat.Network != nil {
		returnVal.NetIn, returnVal.NetInCnt = diff(newStat.Network.BytesIn, oldStat.Network.BytesIn, sampleSecs)
		returnVal.NetOut, returnVal.NetOutCnt = diff(newStat.Network.BytesOut, oldStat.Network.BytesOut, sampleSecs)
	}

	if newStat.Connections != nil {
		returnVal.NumConnections = newStat.Connections.Current
	}

	newReplStat := *newMongo.ReplSetStatus

	if newReplStat.Members != nil {
		myName := newStat.Repl.Me
		// Find the master and myself
		master := ReplSetMember{}
		me := ReplSetMember{}
		for _, member := range newReplStat.Members {
			if member.Name == myName {
				// Store my state string
				returnVal.NodeState = member.StateStr
				if member.State == 1 {
					// I'm the master
					returnVal.ReplLag = 0
					break
				} else {
					// I'm secondary
					me = member
				}
			} else if member.State == 1 {
				// Master found
				master = member
			}
		}

		if me.State == 2 {
			// OptimeDate.Unix() type is int64
			lag := master.OptimeDate.Unix() - me.OptimeDate.Unix()
			if lag < 0 {
				returnVal.ReplLag = 0
			} else {
				returnVal.ReplLag = lag
			}
		}
	}

	newClusterStat := *newMongo.ClusterStatus
	returnVal.JumboChunksCount = newClusterStat.JumboChunksCount
	returnVal.OplogTimeDiff = newMongo.OplogStats.TimeDiff

	newDbStats := *newMongo.DbStats
	for _, db := range newDbStats.Dbs {
		dbStatsData := db.DbStatsData
		// mongos doesn't have the db key, so setting the db name
		if dbStatsData.Db == "" {
			dbStatsData.Db = db.Name
		}
		dbStatLine := &DbStatLine{
			Name:        dbStatsData.Db,
			Collections: dbStatsData.Collections,
			Objects:     dbStatsData.Objects,
			AvgObjSize:  dbStatsData.AvgObjSize,
			DataSize:    dbStatsData.DataSize,
			StorageSize: dbStatsData.StorageSize,
			NumExtents:  dbStatsData.NumExtents,
			Indexes:     dbStatsData.Indexes,
			IndexSize:   dbStatsData.IndexSize,
			Ok:          dbStatsData.Ok,
		}
		returnVal.DbStatsLines = append(returnVal.DbStatsLines, *dbStatLine)
	}

	// Set shard stats
	newShardStats := *newMongo.ShardStats
	returnVal.TotalInUse = newShardStats.TotalInUse
	returnVal.TotalAvailable = newShardStats.TotalAvailable
	returnVal.TotalCreated = newShardStats.TotalCreated
	returnVal.TotalRefreshing = newShardStats.TotalRefreshing
	returnVal.ShardHostStatsLines = map[string]ShardHostStatLine{}
	for host, stats := range newShardStats.Hosts {
		shardStatLine := &ShardHostStatLine{
			InUse:      stats.InUse,
			Available:  stats.Available,
			Created:    stats.Created,
			Refreshing: stats.Refreshing,
		}

		returnVal.ShardHostStatsLines[host] = *shardStatLine
	}

	return returnVal
}
