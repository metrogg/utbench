func NewMountStatsCollector() (Collector, error) {
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}

	proc, err := fs.Self()
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/self: %v", err)
	}

	const (
		// For the time being, only NFS statistics are available via this mechanism.
		subsystem = "mountstats_nfs"
	)

	var (
		labels   = []string{"export", "protocol"}
		opLabels = []string{"export", "protocol", "operation"}
	)

	return &mountStatsCollector{
		NFSAgeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "age_seconds_total"),
			"The age of the NFS mount in seconds.",
			labels,
			nil,
		),

		NFSReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "read_bytes_total"),
			"Number of bytes read using the read() syscall.",
			labels,
			nil,
		),

		NFSWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "write_bytes_total"),
			"Number of bytes written using the write() syscall.",
			labels,
			nil,
		),

		NFSDirectReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "direct_read_bytes_total"),
			"Number of bytes read using the read() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSDirectWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "direct_write_bytes_total"),
			"Number of bytes written using the write() syscall in O_DIRECT mode.",
			labels,
			nil,
		),

		NFSTotalReadBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_read_bytes_total"),
			"Number of bytes read from the NFS server, in total.",
			labels,
			nil,
		),

		NFSTotalWriteBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_write_bytes_total"),
			"Number of bytes written to the NFS server, in total.",
			labels,
			nil,
		),

		NFSReadPagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "read_pages_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSWritePagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "write_pages_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSTransportBindTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_bind_total"),
			"Number of times the client has had to establish a connection from scratch to the NFS server.",
			labels,
			nil,
		),

		NFSTransportConnectTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_connect_total"),
			"Number of times the client has made a TCP connection to the NFS server.",
			labels,
			nil,
		),

		NFSTransportIdleTimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_idle_time_seconds"),
			"Duration since the NFS mount last saw any RPC traffic, in seconds.",
			labels,
			nil,
		),

		NFSTransportSendsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_sends_total"),
			"Number of RPC requests for this mount sent to the NFS server.",
			labels,
			nil,
		),

		NFSTransportReceivesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_receives_total"),
			"Number of RPC responses for this mount received from the NFS server.",
			labels,
			nil,
		),

		NFSTransportBadTransactionIDsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_bad_transaction_ids_total"),
			"Number of times the NFS server sent a response with a transaction ID unknown to this client.",
			labels,
			nil,
		),

		NFSTransportBacklogQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_backlog_queue_total"),
			"Total number of items added to the RPC backlog queue.",
			labels,
			nil,
		),

		NFSTransportMaximumRPCSlots: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_maximum_rpc_slots"),
			"Maximum number of simultaneously active RPC requests ever used.",
			labels,
			nil,
		),

		NFSTransportSendingQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_sending_queue_total"),
			"Total number of items added to the RPC transmission sending queue.",
			labels,
			nil,
		),

		NFSTransportPendingQueueTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transport_pending_queue_total"),
			"Total number of items added to the RPC transmission pending queue.",
			labels,
			nil,
		),

		NFSOperationsRequestsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_requests_total"),
			"Number of requests performed for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsTransmissionsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_transmissions_total"),
			"Number of times an actual RPC request has been transmitted for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsMajorTimeoutsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_major_timeouts_total"),
			"Number of times a request has had a major timeout for a given operation.",
			opLabels,
			nil,
		),

		NFSOperationsSentBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_sent_bytes_total"),
			"Number of bytes sent for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsReceivedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_received_bytes_total"),
			"Number of bytes received for a given operation, including RPC headers and payload.",
			opLabels,
			nil,
		),

		NFSOperationsQueueTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_queue_time_seconds_total"),
			"Duration all requests spent queued for transmission for a given operation before they were sent, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsResponseTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_response_time_seconds_total"),
			"Duration all requests took to get a reply back after a request for a given operation was transmitted, in seconds.",
			opLabels,
			nil,
		),

		NFSOperationsRequestTimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "operations_request_time_seconds_total"),
			"Duration all requests took from when a request was enqueued to when it was completely handled for a given operation, in seconds.",
			opLabels,
			nil,
		),

		NFSEventInodeRevalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_inode_revalidate_total"),
			"Number of times cached inode attributes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDnodeRevalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_dnode_revalidate_total"),
			"Number of times cached dentry nodes are re-validated from the server.",
			labels,
			nil,
		),

		NFSEventDataInvalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_data_invalidate_total"),
			"Number of times an inode cache is cleared.",
			labels,
			nil,
		),

		NFSEventAttributeInvalidateTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_attribute_invalidate_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSOpenTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_open_total"),
			"Number of times cached inode attributes are invalidated.",
			labels,
			nil,
		),

		NFSEventVFSLookupTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_lookup_total"),
			"Number of times a directory lookup has occurred.",
			labels,
			nil,
		),

		NFSEventVFSAccessTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_access_total"),
			"Number of times permissions have been checked.",
			labels,
			nil,
		),

		NFSEventVFSUpdatePageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_update_page_total"),
			"Number of updates (and potential writes) to pages.",
			labels,
			nil,
		),

		NFSEventVFSReadPageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_read_page_total"),
			"Number of pages read directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSReadPagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_read_pages_total"),
			"Number of times a group of pages have been read.",
			labels,
			nil,
		),

		NFSEventVFSWritePageTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_write_page_total"),
			"Number of pages written directly via mmap()'d files.",
			labels,
			nil,
		),

		NFSEventVFSWritePagesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_write_pages_total"),
			"Number of times a group of pages have been written.",
			labels,
			nil,
		),

		NFSEventVFSGetdentsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_getdents_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSSetattrTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_setattr_total"),
			"Number of times directory entries have been read with getdents().",
			labels,
			nil,
		),

		NFSEventVFSFlushTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_flush_total"),
			"Number of pending writes that have been forcefully flushed to the server.",
			labels,
			nil,
		),

		NFSEventVFSFsyncTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_fsync_total"),
			"Number of times fsync() has been called on directories and files.",
			labels,
			nil,
		),

		NFSEventVFSLockTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_lock_total"),
			"Number of times locking has been attempted on a file.",
			labels,
			nil,
		),

		NFSEventVFSFileReleaseTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_vfs_file_release_total"),
			"Number of times files have been closed and released.",
			labels,
			nil,
		),

		NFSEventTruncationTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_truncation_total"),
			"Number of times files have been truncated.",
			labels,
			nil,
		),

		NFSEventWriteExtensionTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_write_extension_total"),
			"Number of times a file has been grown due to writes beyond its existing end.",
			labels,
			nil,
		),

		NFSEventSillyRenameTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_silly_rename_total"),
			"Number of times a file was removed while still open by another process.",
			labels,
			nil,
		),

		NFSEventShortReadTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_short_read_total"),
			"Number of times the NFS server gave less data than expected while reading.",
			labels,
			nil,
		),

		NFSEventShortWriteTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_short_write_total"),
			"Number of times the NFS server wrote less data than expected while writing.",
			labels,
			nil,
		),

		NFSEventJukeboxDelayTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_jukebox_delay_total"),
			"Number of times the NFS server indicated EJUKEBOX; retrieving data from offline storage.",
			labels,
			nil,
		),

		NFSEventPNFSReadTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_pnfs_read_total"),
			"Number of NFS v4.1+ pNFS reads.",
			labels,
			nil,
		),

		NFSEventPNFSWriteTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "event_pnfs_write_total"),
			"Number of NFS v4.1+ pNFS writes.",
			labels,
			nil,
		),

		proc: proc,
	}, nil
}
