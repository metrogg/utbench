func ValidateSetSystemVar(vars *SessionVars, name string, value string) (string, error) {
	if strings.EqualFold(value, "DEFAULT") {
		if val := GetSysVar(name); val != nil {
			return val.Value, nil
		}
		return value, UnknownSystemVar.GenWithStackByArgs(name)
	}
	switch name {
	case ConnectTimeout:
		return checkUInt64SystemVar(name, value, 2, secondsPerYear, vars)
	case DefaultWeekFormat:
		return checkUInt64SystemVar(name, value, 0, 7, vars)
	case DelayKeyWrite:
		if strings.EqualFold(value, "ON") || value == "1" {
			return "ON", nil
		} else if strings.EqualFold(value, "OFF") || value == "0" {
			return "OFF", nil
		} else if strings.EqualFold(value, "ALL") || value == "2" {
			return "ALL", nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case FlushTime:
		return checkUInt64SystemVar(name, value, 0, secondsPerYear, vars)
	case ForeignKeyChecks:
		if strings.EqualFold(value, "ON") || value == "1" {
			// TiDB does not yet support foreign keys.
			// For now, resist the change and show a warning.
			vars.StmtCtx.AppendWarning(ErrUnsupportedValueForVar.GenWithStackByArgs(name, value))
			return "OFF", nil
		} else if strings.EqualFold(value, "OFF") || value == "0" {
			return "OFF", nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case GroupConcatMaxLen:
		// The reasonable range of 'group_concat_max_len' is 4~18446744073709551615(64-bit platforms)
		// See https://dev.mysql.com/doc/refman/8.0/en/server-system-variables.html#sysvar_group_concat_max_len for details
		return checkUInt64SystemVar(name, value, 4, math.MaxUint64, vars)
	case InteractiveTimeout:
		return checkUInt64SystemVar(name, value, 1, secondsPerYear, vars)
	case InnodbCommitConcurrency:
		return checkUInt64SystemVar(name, value, 0, 1000, vars)
	case InnodbFastShutdown:
		return checkUInt64SystemVar(name, value, 0, 2, vars)
	case InnodbLockWaitTimeout:
		return checkUInt64SystemVar(name, value, 1, 1073741824, vars)
	// See "https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_max_allowed_packet"
	case MaxAllowedPacket:
		return checkUInt64SystemVar(name, value, 1024, 1073741824, vars)
	case MaxConnections:
		return checkUInt64SystemVar(name, value, 1, 100000, vars)
	case MaxConnectErrors:
		return checkUInt64SystemVar(name, value, 1, math.MaxUint64, vars)
	case MaxSortLength:
		return checkUInt64SystemVar(name, value, 4, 8388608, vars)
	case MaxSpRecursionDepth:
		return checkUInt64SystemVar(name, value, 0, 255, vars)
	case MaxUserConnections:
		return checkUInt64SystemVar(name, value, 0, 4294967295, vars)
	case OldPasswords:
		return checkUInt64SystemVar(name, value, 0, 2, vars)
	case SessionTrackGtids:
		if strings.EqualFold(value, "OFF") || value == "0" {
			return "OFF", nil
		} else if strings.EqualFold(value, "OWN_GTID") || value == "1" {
			return "OWN_GTID", nil
		} else if strings.EqualFold(value, "ALL_GTIDS") || value == "2" {
			return "ALL_GTIDS", nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case SQLSelectLimit:
		return checkUInt64SystemVar(name, value, 0, math.MaxUint64, vars)
	case SyncBinlog:
		return checkUInt64SystemVar(name, value, 0, 4294967295, vars)
	case TableDefinitionCache:
		return checkUInt64SystemVar(name, value, 400, 524288, vars)
	case TmpTableSize:
		return checkUInt64SystemVar(name, value, 1024, math.MaxUint64, vars)
	case WaitTimeout:
		return checkUInt64SystemVar(name, value, 1, 31536000, vars)
	case MaxPreparedStmtCount:
		return checkInt64SystemVar(name, value, -1, 1048576, vars)
	case TimeZone:
		if strings.EqualFold(value, "SYSTEM") {
			return "SYSTEM", nil
		}
		_, err := parseTimeZone(value)
		return value, err
	case ValidatePasswordLength, ValidatePasswordNumberCount:
		return checkUInt64SystemVar(name, value, 0, math.MaxUint64, vars)
	case WarningCount, ErrorCount:
		return value, ErrReadOnly.GenWithStackByArgs(name)
	case GeneralLog, TiDBGeneralLog, AvoidTemporalUpgrade, BigTables, CheckProxyUsers, LogBin,
		CoreFile, EndMakersInJSON, SQLLogBin, OfflineMode, PseudoSlaveMode, LowPriorityUpdates,
		SkipNameResolve, SQLSafeUpdates, TiDBConstraintCheckInPlace, serverReadOnly:
		if strings.EqualFold(value, "ON") || value == "1" {
			return "1", nil
		} else if strings.EqualFold(value, "OFF") || value == "0" {
			return "0", nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case AutocommitVar, TiDBSkipUTF8Check, TiDBOptAggPushDown,
		TiDBOptInSubqToJoinAndAgg, TiDBEnableFastAnalyze,
		TiDBBatchInsert, TiDBDisableTxnAutoRetry, TiDBEnableStreaming,
		TiDBBatchDelete, TiDBBatchCommit, TiDBEnableCascadesPlanner, TiDBEnableWindowFunction,
		TiDBCheckMb4ValueInUTF8:
		if strings.EqualFold(value, "ON") || value == "1" || strings.EqualFold(value, "OFF") || value == "0" {
			return value, nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case TiDBEnableTablePartition:
		switch {
		case strings.EqualFold(value, "ON") || value == "1":
			return "on", nil
		case strings.EqualFold(value, "OFF") || value == "0":
			return "off", nil
		case strings.EqualFold(value, "AUTO"):
			return "auto", nil
		}
		return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
	case TiDBDDLReorgBatchSize:
		return checkUInt64SystemVar(name, value, uint64(MinDDLReorgBatchSize), uint64(MaxDDLReorgBatchSize), vars)
	case TiDBDDLErrorCountLimit:
		return checkUInt64SystemVar(name, value, uint64(0), math.MaxInt64, vars)
	case TiDBIndexLookupConcurrency, TiDBIndexLookupJoinConcurrency, TiDBIndexJoinBatchSize,
		TiDBIndexLookupSize,
		TiDBHashJoinConcurrency,
		TiDBHashAggPartialConcurrency,
		TiDBHashAggFinalConcurrency,
		TiDBDistSQLScanConcurrency,
		TiDBIndexSerialScanConcurrency, TiDBDDLReorgWorkerCount,
		TiDBBackoffLockFast,
		TiDBDMLBatchSize, TiDBOptimizerSelectivityLevel:
		v, err := strconv.Atoi(value)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v <= 0 {
			return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
		}
		return value, nil
	case TiDBOptCorrelationExpFactor:
		v, err := strconv.Atoi(value)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v < 0 {
			return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
		}
		return value, nil
	case TiDBOptCorrelationThreshold:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v < 0 || v > 1 {
			return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
		}
		return value, nil
	case TiDBProjectionConcurrency,
		TIDBMemQuotaQuery,
		TIDBMemQuotaHashJoin,
		TIDBMemQuotaMergeJoin,
		TIDBMemQuotaSort,
		TIDBMemQuotaTopn,
		TIDBMemQuotaIndexLookupReader,
		TIDBMemQuotaIndexLookupJoin,
		TIDBMemQuotaNestedLoopApply,
		TiDBRetryLimit,
		TiDBSlowLogThreshold,
		TiDBQueryLogMaxLen:
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return value, ErrWrongValueForVar.GenWithStackByArgs(name)
		}
		return value, nil
	case TiDBAutoAnalyzeStartTime, TiDBAutoAnalyzeEndTime:
		v, err := setAnalyzeTime(vars, value)
		if err != nil {
			return "", err
		}
		return v, nil
	case TxnIsolation, TransactionIsolation:
		upVal := strings.ToUpper(value)
		_, exists := TxIsolationNames[upVal]
		if !exists {
			return "", ErrWrongValueForVar.GenWithStackByArgs(name, value)
		}
		switch upVal {
		case "SERIALIZABLE", "READ-UNCOMMITTED":
			skipIsolationLevelCheck, err := GetSessionSystemVar(vars, TiDBSkipIsolationLevelCheck)
			returnErr := ErrUnsupportedIsolationLevel.GenWithStackByArgs(value)
			if err != nil {
				returnErr = err
			}
			if !TiDBOptOn(skipIsolationLevelCheck) || err != nil {
				return "", returnErr
			}
			//SET TRANSACTION ISOLATION LEVEL will affect two internal variables:
			// 1. tx_isolation
			// 2. transaction_isolation
			// The following if condition is used to deduplicate two same warnings.
			if name == "transaction_isolation" {
				vars.StmtCtx.AppendWarning(returnErr)
			}
		}
		return upVal, nil
	case TiDBInitChunkSize:
		v, err := strconv.Atoi(value)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v <= 0 {
			return value, ErrWrongValueForVar.GenWithStackByArgs(name, value)
		}
		if v > initChunkSizeUpperBound {
			return value, errors.Errorf("tidb_init_chunk_size(%d) cannot be bigger than %d", v, initChunkSizeUpperBound)
		}
		return value, nil
	case TiDBMaxChunkSize:
		v, err := strconv.Atoi(value)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v < maxChunkSizeLowerBound {
			return value, errors.Errorf("tidb_max_chunk_size(%d) cannot be smaller than %d", v, maxChunkSizeLowerBound)
		}
		return value, nil
	case TiDBOptJoinReorderThreshold:
		v, err := strconv.Atoi(value)
		if err != nil {
			return value, ErrWrongTypeForVar.GenWithStackByArgs(name)
		}
		if v < 0 || v >= 64 {
			return value, errors.Errorf("tidb_join_order_algo_threshold(%d) cannot be smaller than 0 or larger than 63", v)
		}
	}
	return value, nil
}
