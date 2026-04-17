func (a *APIServer) waitJob(pachClient *client.APIClient, jobInfo *pps.JobInfo, logger *taggedLogger) (retErr error) {
	logger.Logf("waitJob: %s", jobInfo.Job.ID)
	ctx, cancel := context.WithCancel(pachClient.Ctx())
	pachClient = pachClient.WithCtx(ctx)

	// Watch the output commit to see if it's terminated (KILLED, FAILED, or
	// SUCCESS) and if so, cancel the current context
	go func() {
		backoff.RetryNotify(func() error {
			commitInfo, err := pachClient.PfsAPIClient.InspectCommit(ctx,
				&pfs.InspectCommitRequest{
					Commit:     jobInfo.OutputCommit,
					BlockState: pfs.CommitState_FINISHED,
				})
			if err != nil {
				if pfsserver.IsCommitNotFoundErr(err) || pfsserver.IsCommitDeletedErr(err) {
					defer cancel() // whether we return error or nil, job is done
					// Output commit was deleted. Delete job as well
					if _, err := col.NewSTM(ctx, a.etcdClient, func(stm col.STM) error {
						// Delete the job if no other worker has deleted it yet
						jobPtr := &pps.EtcdJobInfo{}
						if err := a.jobs.ReadWrite(stm).Get(jobInfo.Job.ID, jobPtr); err != nil {
							return err
						}
						return a.deleteJob(stm, jobPtr)
					}); err != nil && !col.IsErrNotFound(err) {
						return err
					}
					return nil
				}
				return err
			}
			if commitInfo.Trees == nil {
				defer cancel() // whether job state update succeeds or not, job is done
				if _, err := col.NewSTM(ctx, a.etcdClient, func(stm col.STM) error {
					// Read an up to date version of the jobInfo so that we
					// don't overwrite changes that have happened since this
					// function started.
					jobPtr := &pps.EtcdJobInfo{}
					if err := a.jobs.ReadWrite(stm).Get(jobInfo.Job.ID, jobPtr); err != nil {
						return err
					}
					if !ppsutil.IsTerminal(jobPtr.State) {
						return ppsutil.UpdateJobState(a.pipelines.ReadWrite(stm), a.jobs.ReadWrite(stm), jobPtr, pps.JobState_JOB_KILLED, "")
					}
					return nil
				}); err != nil {
					return err
				}
			}
			return nil
		}, backoff.NewInfiniteBackOff(), func(err error, d time.Duration) error {
			if isDone(ctx) {
				return err // exit retry loop
			}
			return nil // retry again
		})
	}()
	if jobInfo.JobTimeout != nil {
		startTime, err := types.TimestampFromProto(jobInfo.Started)
		if err != nil {
			return err
		}
		timeout, err := types.DurationFromProto(jobInfo.JobTimeout)
		if err != nil {
			return err
		}
		afterTime := startTime.Add(timeout).Sub(time.Now())
		logger.Logf("cancelling job at: %+v", afterTime)
		timer := time.AfterFunc(afterTime, func() {
			if _, err := pachClient.PfsAPIClient.FinishCommit(ctx,
				&pfs.FinishCommitRequest{
					Commit: jobInfo.OutputCommit,
					Empty:  true,
				}); err != nil {
				logger.Logf("error from FinishCommit while timing out job: %+v", err)
			}
		})
		defer timer.Stop()
	}

	backoff.RetryNotify(func() (retErr error) {
		// block until job inputs are ready
		// TODO(bryce) This should be removed because it is no longer applicable.
		failedInputs, err := a.failedInputs(ctx, jobInfo)
		if err != nil {
			return err
		}
		if len(failedInputs) > 0 {
			reason := fmt.Sprintf("inputs %s failed", strings.Join(failedInputs, ", "))
			if err := a.updateJobState(ctx, jobInfo, nil, pps.JobState_JOB_FAILURE, reason); err != nil {
				return err
			}
			_, err := pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
				Commit: jobInfo.OutputCommit,
				Empty:  true,
			})
			return err
		}

		// Create a datum factory pointing at the job's inputs and split up the
		// input data into chunks
		df, err := NewDatumFactory(pachClient, jobInfo.Input)
		if err != nil {
			return err
		}
		parallelism, err := ppsutil.GetExpectedNumWorkers(a.kubeClient, a.pipelineInfo.ParallelismSpec)
		if err != nil {
			return fmt.Errorf("error from GetExpectedNumWorkers: %v", err)
		}
		numHashtrees, err := ppsutil.GetExpectedNumHashtrees(a.pipelineInfo.HashtreeSpec)
		if err != nil {
			return fmt.Errorf("error from GetExpectedNumHashtrees: %v", err)
		}
		plan := &Plan{}
		// Get stats commit
		var statsCommit *pfs.Commit
		var statsTrees []*pfs.Object
		var statsSize uint64
		if jobInfo.EnableStats {
			ci, err := pachClient.InspectCommit(jobInfo.OutputCommit.Repo.Name, jobInfo.OutputCommit.ID)
			if err != nil {
				return err
			}
			for _, commitRange := range ci.Subvenance {
				if commitRange.Lower.Repo.Name == jobInfo.OutputRepo.Name && commitRange.Upper.Repo.Name == jobInfo.OutputRepo.Name {
					statsCommit = commitRange.Lower
				}
			}
		}
		// Read the job document, and either resume (if we're recovering from a
		// crash) or mark it running. Also write the input chunks calculated above
		// into plansCol
		jobID := jobInfo.Job.ID
		if _, err := col.NewSTM(ctx, a.etcdClient, func(stm col.STM) error {
			jobs := a.jobs.ReadWrite(stm)
			jobPtr := &pps.EtcdJobInfo{}
			if err := jobs.Get(jobID, jobPtr); err != nil {
				return err
			}
			if jobPtr.State == pps.JobState_JOB_KILLED {
				return nil
			}
			jobPtr.DataTotal = int64(df.Len())
			jobPtr.StatsCommit = statsCommit
			if err := ppsutil.UpdateJobState(a.pipelines.ReadWrite(stm), a.jobs.ReadWrite(stm), jobPtr, pps.JobState_JOB_RUNNING, ""); err != nil {
				return err
			}
			plansCol := a.plans.ReadWrite(stm)
			if err := plansCol.Get(jobID, plan); err == nil {
				return nil
			}
			plan = newPlan(df, jobInfo.ChunkSpec, parallelism, numHashtrees)
			return plansCol.Put(jobID, plan)
		}); err != nil {
			return err
		}
		defer func() {
			if retErr == nil {
				if _, err := col.NewSTM(ctx, a.etcdClient, func(stm col.STM) error {
					chunksCol := a.chunks(jobID).ReadWrite(stm)
					chunksCol.DeleteAll()
					plansCol := a.plans.ReadWrite(stm)
					return plansCol.Delete(jobID)
				}); err != nil {
					retErr = err
				}
			}
		}()
		// Handle the case when there are no datums
		if df.Len() == 0 {
			if err := a.updateJobState(ctx, jobInfo, nil, pps.JobState_JOB_SUCCESS, ""); err != nil {
				return err
			}
			if jobInfo.EnableStats {
				if _, err = pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
					Commit:    statsCommit,
					Trees:     statsTrees,
					SizeBytes: statsSize,
				}); err != nil {
					return err
				}
			}
			_, err := pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
				Commit: jobInfo.OutputCommit,
				Empty:  true,
			})
			return err
		}
		// Watch the chunks in order
		chunks := a.chunks(jobInfo.Job.ID).ReadOnly(ctx)
		var failedDatumID string
		for _, high := range plan.Chunks {
			chunkState := &ChunkState{}
			if err := chunks.WatchOneF(fmt.Sprint(high), func(e *watch.Event) error {
				var key string
				if err := e.Unmarshal(&key, chunkState); err != nil {
					return err
				}
				if key != fmt.Sprint(high) {
					return nil
				}
				if chunkState.State != State_RUNNING {
					if chunkState.State == State_FAILED {
						failedDatumID = chunkState.DatumID
					}
					return errutil.ErrBreak
				}
				return nil
			}); err != nil {
				return err
			}
		}
		if err := a.updateJobState(ctx, jobInfo, nil, pps.JobState_JOB_MERGING, ""); err != nil {
			return err
		}
		var trees []*pfs.Object
		var size uint64
		if failedDatumID == "" || jobInfo.EnableStats {
			// Wait for all merges to happen.
			merges := a.merges(jobInfo.Job.ID).ReadOnly(ctx)
			for merge := int64(0); merge < plan.Merges; merge++ {
				mergeState := &MergeState{}
				if err := merges.WatchOneF(fmt.Sprint(merge), func(e *watch.Event) error {
					var key string
					if err := e.Unmarshal(&key, mergeState); err != nil {
						return err
					}
					if key != fmt.Sprint(merge) {
						return nil
					}
					if mergeState.State != State_RUNNING {
						trees = append(trees, mergeState.Tree)
						size += mergeState.SizeBytes
						statsTrees = append(statsTrees, mergeState.StatsTree)
						statsSize += mergeState.StatsSizeBytes
						return errutil.ErrBreak
					}
					return nil
				}); err != nil {
					return err
				}
			}
		}
		if jobInfo.EnableStats {
			if _, err = pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
				Commit:    statsCommit,
				Trees:     statsTrees,
				SizeBytes: statsSize,
			}); err != nil {
				return err
			}
		}
		// If the job failed we finish the commit with an empty tree but only
		// after we've set the state, otherwise the job will be considered
		// killed.
		if failedDatumID != "" {
			reason := fmt.Sprintf("failed to process datum: %v", failedDatumID)
			if err := a.updateJobState(ctx, jobInfo, statsCommit, pps.JobState_JOB_FAILURE, reason); err != nil {
				return err
			}
			_, err = pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
				Commit: jobInfo.OutputCommit,
				Empty:  true,
			})
			return err
		}
		// Write out the datums processed/skipped and merged for this job
		buf := &bytes.Buffer{}
		pbw := pbutil.NewWriter(buf)
		for i := 0; i < df.Len(); i++ {
			files := df.Datum(i)
			datumHash := HashDatum(a.pipelineInfo.Pipeline.Name, a.pipelineInfo.Salt, files)
			if _, err := pbw.WriteBytes([]byte(datumHash)); err != nil {
				return err
			}
		}
		datums, _, err := pachClient.PutObject(buf)
		if err != nil {
			return err
		}
		// Finish the job's output commit
		_, err = pachClient.PfsAPIClient.FinishCommit(ctx, &pfs.FinishCommitRequest{
			Commit:    jobInfo.OutputCommit,
			Trees:     trees,
			SizeBytes: size,
			Datums:    datums,
		})
		if err != nil && !pfsserver.IsCommitFinishedErr(err) {
			if pfsserver.IsCommitNotFoundErr(err) || pfsserver.IsCommitDeletedErr(err) {
				// output commit was deleted during e.g. FinishCommit, which means this job
				// should be deleted. Goro from top of waitJob() will observe the deletion,
				// delete the jobPtr and call cancel()--wait for that.
				<-ctx.Done() // wait for cancel()
				return nil
			}
			return err
		}
		// Handle egress
		if err := a.egress(pachClient, logger, jobInfo); err != nil {
			reason := fmt.Sprintf("egress error: %v", err)
			return a.updateJobState(ctx, jobInfo, statsCommit, pps.JobState_JOB_FAILURE, reason)
		}
		return a.updateJobState(ctx, jobInfo, statsCommit, pps.JobState_JOB_SUCCESS, "")
	}, backoff.NewInfiniteBackOff(), func(err error, d time.Duration) error {
		logger.Logf("error in waitJob %v, retrying in %v", err, d)
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				if err == context.DeadlineExceeded {
					reason := fmt.Sprintf("job exceeded timeout (%v)", jobInfo.JobTimeout)
					// Mark the job as failed.
					// Workers subscribe to etcd for this state change to cancel their work
					_, err := col.NewSTM(context.Background(), a.etcdClient, func(stm col.STM) error {
						jobs := a.jobs.ReadWrite(stm)
						jobID := jobInfo.Job.ID
						jobPtr := &pps.EtcdJobInfo{}
						if err := jobs.Get(jobID, jobPtr); err != nil {
							return err
						}
						err = ppsutil.UpdateJobState(a.pipelines.ReadWrite(stm), a.jobs.ReadWrite(stm), jobPtr, pps.JobState_JOB_FAILURE, reason)
						if err != nil {
							return nil
						}
						return nil
					})
					if err != nil {
						return err
					}
				}
				return err
			}
			return err
		default:
		}
		// Increment the job's restart count
		_, err = col.NewSTM(ctx, a.etcdClient, func(stm col.STM) error {
			jobs := a.jobs.ReadWrite(stm)
			jobID := jobInfo.Job.ID
			jobPtr := &pps.EtcdJobInfo{}
			if err := jobs.Get(jobID, jobPtr); err != nil {
				return err
			}
			jobPtr.Restart++
			return jobs.Put(jobID, jobPtr)
		})
		if err != nil {
			logger.Logf("error incrementing job %s's restart count", jobInfo.Job.ID)
		}
		return nil
	})
	return nil
}
