func (d *driver) makeCommit(pachClient *client.APIClient, ID string, parent *pfs.Commit, branch string, provenance []*pfs.CommitProvenance, treeRef *pfs.Object, recordFiles []string, records []*pfs.PutFileRecords, description string) (*pfs.Commit, error) {
	// Validate arguments:
	if parent == nil {
		return nil, fmt.Errorf("parent cannot be nil")
	}

	// Check that caller is authorized
	if err := d.checkIsAuthorized(pachClient, parent.Repo, auth.Scope_WRITER); err != nil {
		return nil, err
	}

	// New commit and commitInfo
	newCommit := &pfs.Commit{
		Repo: parent.Repo,
		ID:   ID,
	}
	if newCommit.ID == "" {
		newCommit.ID = uuid.NewWithoutDashes()
	}
	newCommitInfo := &pfs.CommitInfo{
		Commit:      newCommit,
		Started:     now(),
		Description: description,
	}

	// FinishCommit case. We need to create AND finish 'newCommit' in two cases:
	// 1. PutFile has been called on a finished commit (records != nil), and
	//    we want to apply 'records' to the parentCommit's HashTree, use that new
	//    HashTree for this commit's filesystem, and then finish this commit
	// 2. BuildCommit has been called by migration (treeRef != nil) and we
	//    want a new, finished commit with the given treeRef
	// - In either case, store this commit's HashTree in 'tree', so we have its
	//   size, and store a pointer to the tree (in object store) in 'treeRef', to
	//   put in newCommitInfo.Tree.
	// - We also don't want to resolve 'branch' or 'parent.ID' (if it's a branch)
	//   outside the txn below, so the 'PutFile' case is handled (by computing
	//   'tree' and 'treeRef') below as well
	var tree hashtree.HashTree
	if treeRef != nil {
		var err error
		tree, err = hashtree.GetHashTreeObject(pachClient, d.storageRoot, treeRef)
		if err != nil {
			return nil, err
		}
	}

	// Txn: create the actual commit in etcd and update the branch + parent/child
	if _, err := col.NewSTM(pachClient.Ctx(), d.etcdClient, func(stm col.STM) error {
		// Clone the parent, as this stm modifies it and might wind up getting
		// run more than once (if there's a conflict.)
		parent := proto.Clone(parent).(*pfs.Commit)
		repos := d.repos.ReadWrite(stm)
		commits := d.commits(parent.Repo.Name).ReadWrite(stm)
		branches := d.branches(parent.Repo.Name).ReadWrite(stm)

		// Check if repo exists
		repoInfo := new(pfs.RepoInfo)
		if err := repos.Get(parent.Repo.Name, repoInfo); err != nil {
			return err
		}

		// create/update 'branch' (if it was set) and set parent.ID (if, in
		// addition, 'parent.ID' was not set)
		if branch != "" {
			branchInfo := &pfs.BranchInfo{}
			if err := branches.Upsert(branch, branchInfo, func() error {
				// validate branch
				if parent.ID == "" && branchInfo.Head != nil {
					parent.ID = branchInfo.Head.ID
				}
				// Don't count the __spec__ repo towards the provenance count
				// since spouts will have __spec__ as provenance, but need to accept commits
				provenanceCount := len(branchInfo.Provenance)
				for _, p := range branchInfo.Provenance {
					if p.Repo.Name == ppsconsts.SpecRepo {
						provenanceCount--
						break
					}
				}
				if provenanceCount > 0 && treeRef == nil {
					return fmt.Errorf("cannot start a commit on an output branch")
				}
				// Point 'branch' at the new commit
				branchInfo.Name = branch // set in case 'branch' is new
				branchInfo.Head = newCommit
				branchInfo.Branch = client.NewBranch(newCommit.Repo.Name, branch)
				return nil
			}); err != nil {
				return err
			}
			// Add branch to repo (see "Update repoInfo" below)
			add(&repoInfo.Branches, branchInfo.Branch)
			// and add the branch to the commit info
			newCommitInfo.Branch = branchInfo.Branch
		}

		// Set newCommit.ParentCommit (if 'parent' and/or 'branch' was set) and add
		// newCommit to parent's ChildCommits
		if parent.ID != "" {
			// Resolve parent.ID if it's a branch that isn't 'branch' (which can
			// happen if 'branch' is new and diverges from the existing branch in
			// 'parent.ID')
			parentCommitInfo, err := d.resolveCommit(stm, parent)
			if err != nil {
				return fmt.Errorf("parent commit not found: %v", err)
			}
			// fail if the parent commit has not been finished
			if parentCommitInfo.Finished == nil {
				return fmt.Errorf("parent commit %s has not been finished", parent.ID)
			}
			if err := commits.Update(parent.ID, parentCommitInfo, func() error {
				newCommitInfo.ParentCommit = parent
				// If we don't know the branch the commit belongs to at this point, assume it is the same as the parent branch
				if newCommitInfo.Branch == nil {
					newCommitInfo.Branch = parentCommitInfo.Branch
				}
				parentCommitInfo.ChildCommits = append(parentCommitInfo.ChildCommits, newCommit)
				return nil
			}); err != nil {
				// Note: error is emitted if parent.ID is a missing/invalid branch OR a
				// missing/invalid commit ID
				return fmt.Errorf("could not resolve parent commit \"%s\": %v", parent.ID, err)
			}
		}

		// 1. Write 'newCommit' to 'openCommits' collection OR
		// 2. Finish 'newCommit' (if treeRef != nil or records != nil); see
		//    "FinishCommit case" above)
		if treeRef != nil || records != nil {
			if records != nil {
				parentTree, err := d.getTreeForCommit(pachClient, parent)
				if err != nil {
					return err
				}
				tree, err = parentTree.Copy()
				if err != nil {
					return err
				}
				for i, record := range records {
					if err := d.applyWrite(recordFiles[i], record, tree); err != nil {
						return err
					}
				}
				if err := tree.Hash(); err != nil {
					return err
				}
				treeRef, err = hashtree.PutHashTree(pachClient, tree)
				if err != nil {
					return err
				}
			}

			// now 'treeRef' is guaranteed to be set
			newCommitInfo.Tree = treeRef
			newCommitInfo.SizeBytes = uint64(tree.FSSize())
			newCommitInfo.Finished = now()

			// If we're updating the master branch, also update the repo size (see
			// "Update repoInfo" below)
			if branch == "master" {
				repoInfo.SizeBytes = newCommitInfo.SizeBytes
			}
		} else {
			if err := d.openCommits.ReadWrite(stm).Put(newCommit.ID, newCommit); err != nil {
				return err
			}
		}

		// Update repoInfo (potentially with new branch and new size)
		if err := repos.Put(parent.Repo.Name, repoInfo); err != nil {
			return err
		}

		// Build newCommit's full provenance. B/c commitInfo.Provenance is a
		// transitive closure, there's no need to search the full provenance graph,
		// just take the union of the immediate parents' (in the 'provenance' arg)
		// commitInfo.Provenance
		newCommitProv := make(map[string]*pfs.CommitProvenance)
		for _, prov := range provenance {
			newCommitProv[prov.Commit.ID] = prov
			provCommitInfo := &pfs.CommitInfo{}
			if err := d.commits(prov.Commit.Repo.Name).ReadWrite(stm).Get(prov.Commit.ID, provCommitInfo); err != nil {
				return err
			}
			for _, c := range provCommitInfo.Provenance {
				newCommitProv[c.Commit.ID] = c
			}
		}

		// Copy newCommitProv into newCommitInfo.Provenance, and update upstream subv
		for _, prov := range newCommitProv {
			newCommitInfo.Provenance = append(newCommitInfo.Provenance, prov)
			provCommitInfo := &pfs.CommitInfo{}
			if err := d.commits(prov.Commit.Repo.Name).ReadWrite(stm).Update(prov.Commit.ID, provCommitInfo, func() error {
				appendSubvenance(provCommitInfo, newCommitInfo)
				return nil
			}); err != nil {
				return err
			}
		}

		// Finally, create the commit
		if err := commits.Create(newCommit.ID, newCommitInfo); err != nil {
			return err
		}
		// We propagate the branch last so propagateCommit can write to the
		// now-existing commit's subvenance
		if branch != "" {
			return d.propagateCommit(stm, client.NewBranch(newCommit.Repo.Name, branch))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return newCommit, nil
}
