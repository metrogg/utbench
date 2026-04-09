package main

func CloseNtfnChans() {
	if NtfnChans.ConnectChan != nil {
		close(NtfnChans.ConnectChan)
	}
	if NtfnChans.ConnectChanWiredDB != nil {
		close(NtfnChans.ConnectChanWiredDB)
	}
	if NtfnChans.ConnectChanStakeDB != nil {
		close(NtfnChans.ConnectChanStakeDB)
	}
	if NtfnChans.ConnectChanDcrpgDB != nil {
		close(NtfnChans.ConnectChanDcrpgDB)
	}

	if NtfnChans.ReorgChanBlockData != nil {
		close(NtfnChans.ReorgChanBlockData)
	}
	if NtfnChans.ReorgChanWiredDB != nil {
		close(NtfnChans.ReorgChanWiredDB)
	}
	if NtfnChans.ReorgChanStakeDB != nil {
		close(NtfnChans.ReorgChanStakeDB)
	}
	if NtfnChans.ReorgChanDcrpgDB != nil {
		close(NtfnChans.ReorgChanDcrpgDB)
	}

	if NtfnChans.UpdateStatusNodeHeight != nil {
		close(NtfnChans.UpdateStatusNodeHeight)
	}
	if NtfnChans.UpdateStatusDBHeight != nil {
		close(NtfnChans.UpdateStatusDBHeight)
	}

	if NtfnChans.NewTxChan != nil {
		close(NtfnChans.NewTxChan)
	}
	if NtfnChans.RelevantTxMempoolChan != nil {
		close(NtfnChans.RelevantTxMempoolChan)
	}

	if NtfnChans.SpendTxBlockChan != nil {
		close(NtfnChans.SpendTxBlockChan)
	}
	if NtfnChans.RecvTxBlockChan != nil {
		close(NtfnChans.RecvTxBlockChan)
	}

	if NtfnChans.InsightNewTxChan != nil {
		close(NtfnChans.InsightNewTxChan)
	}

	if NtfnChans.ReorgChartsCache != nil {
		close(NtfnChans.ReorgChartsCache)
	}
}
