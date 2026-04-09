package main

import (
	"errors"
	"fmt"
)

func VerifyConfig() error {
	if err := Config.verifyTransactionLimitConfig(); err != nil {
		return err
	}
	if actual, dryRun := Config.EnableHotRowProtection, Config.EnableHotRowProtectionDryRun; actual && dryRun {
		return errors.New("only one of two flags allowed: -enable_hot_row_protection or -enable_hot_row_protection_dry_run")
	}
	if v := Config.HotRowProtectionMaxQueueSize; v <= 0 {
		return fmt.Errorf("-hot_row_protection_max_queue_size must be > 0 (specified value: %v)", v)
	}
	if v := Config.HotRowProtectionMaxGlobalQueueSize; v <= 0 {
		return fmt.Errorf("-hot_row_protection_max_global_queue_size must be > 0 (specified value: %v)", v)
	}
	if globalSize, size := Config.HotRowProtectionMaxGlobalQueueSize, Config.HotRowProtectionMaxQueueSize; globalSize < size {
		return fmt.Errorf("global queue size must be >= per row (range) queue size: -hot_row_protection_max_global_queue_size < hot_row_protection_max_queue_size (%v < %v)", globalSize, size)
	}
	if v := Config.HotRowProtectionConcurrentTransactions; v <= 0 {
		return fmt.Errorf("-hot_row_protection_concurrent_transactions must be > 0 (specified value: %v)", v)
	}
	return nil
}
