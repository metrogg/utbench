func (cr *ConflictResolver) fixRenameConflicts(ctx context.Context,
	unmergedChains, mergedChains *crChains,
	mergedPaths map[data.BlockPointer]data.Path) ([]data.Path, error) {
	// For every renamed block pointer in the unmerged chains:
	//   * Check if any BlockPointer in its merged path contains a relative of
	//     itself
	//   * If so, replace the corresponding unmerged create operation with a
	//     symlink creation to the new merged path instead.
	// So, if in the merged branch someone did `mv b/ a/` and in the unmerged
	// branch someone did `mv a/ b/`, the conflict resolution would end up with
	// `a/b/a` where the second a is a symlink to "../".
	//
	// To calculate what the symlink should be, consider the following:
	//   * The unmerged path for the new parent of ptr P is u_1/u_2/.../u_n
	//   * u_i is the largest i <= n such that the corresponding block
	//     can be mapped to a node in merged branch (pointer m_j).
	//   * The full path to m_j in the merged branch is m_1/m_2/m_3/.../m_j
	//   * For a rename cycle to occur, some m_x where x <= j must be a
	//     descendant of P's original pointer.
	//   * The full merged path to the parent of the second copy of P will
	//     then be: m_1/m_2/.../m_x/.../m_j/u_i+1/.../u_n.
	//   * Then, the symlink to put under P's name in u_n is "../"*((n-i)+(j-x))
	// In the case that u_n is a directory that was newly-created in the
	// unmerged branch, we also need to construct a complete corresponding
	// merged path, for use in later stages (like executing actions).  This
	// merged path is just m_1/.../m_j/u_i+1/.../u_n, using the most recent
	// unmerged pointers.
	var newUnmergedPaths []data.Path
	var removeRenames []data.BlockPointer
	var doubleRenames []data.BlockPointer // merged most recent ptrs
	for ptr, info := range unmergedChains.renamedOriginals {
		if unmergedChains.isDeleted(ptr) {
			continue
		}

		// Also, we need to get the merged paths for anything that was
		// renamed in both branches, if they are different.
		if mergedInfo, ok := mergedChains.renamedOriginals[ptr]; ok &&
			(info.originalNewParent != mergedInfo.originalNewParent ||
				info.newName != mergedInfo.newName) {
			mergedMostRecent, err :=
				mergedChains.mostRecentFromOriginalOrSame(ptr)
			if err != nil {
				return nil, err
			}

			doubleRenames = append(doubleRenames, mergedMostRecent)
			continue
		}

		// If this node was modified in both branches, we need to fork
		// the node, so we can get rid of the unmerged remove op and
		// force a copy on the create op.
		unmergedChain := unmergedChains.byOriginal[ptr]
		mergedChain := mergedChains.byOriginal[ptr]
		if crConflictCheckQuick(unmergedChain, mergedChain) {
			cr.log.CDebugf(ctx, "File that was renamed on the unmerged "+
				"branch from %s -> %s has conflicting edits, forking "+
				"(original ptr %v)", info.oldName, info.newName, ptr)
			oldParent := unmergedChains.byOriginal[info.originalOldParent]
			for _, op := range oldParent.ops {
				ro, ok := op.(*rmOp)
				if !ok {
					continue
				}
				if ro.OldName == info.oldName {
					ro.dropThis = true
					break
				}
			}
			newParent := unmergedChains.byOriginal[info.originalNewParent]
			for _, npOp := range newParent.ops {
				co, ok := npOp.(*createOp)
				if !ok {
					continue
				}
				if co.NewName == info.newName && co.renamed {
					co.forceCopy = true
					co.renamed = false
					co.AddRefBlock(unmergedChain.mostRecent)
					co.DelRefBlock(ptr)
					// Clear out the ops on the file itself, as we
					// will be doing a fresh create instead.
					unmergedChain.ops = nil
					break
				}
			}
			// Reset the chain of the forked file to the most recent
			// pointer, since we want to avoid any local notifications
			// linking the old version of the file to the new one.
			if ptr != unmergedChain.mostRecent {
				err := unmergedChains.changeOriginal(
					ptr, unmergedChain.mostRecent)
				if err != nil {
					return nil, err
				}
				unmergedChains.createdOriginals[unmergedChain.mostRecent] = true
			}
			continue
		}

		// The merged path is keyed by the most recent unmerged tail
		// pointer.
		parent, err :=
			unmergedChains.mostRecentFromOriginal(info.originalNewParent)
		if err != nil {
			return nil, err
		}

		mergedPath, ok := mergedPaths[parent]
		unmergedWalkBack := 0 // (n-i) in the equation above
		var unmergedPath data.Path
		if !ok {
			// If this parent was newly created in the unmerged
			// branch, we need to look up its earliest parent that
			// existed in both branches.
			if !unmergedChains.isCreated(info.originalNewParent) {
				// There should definitely be a merged path for this
				// parent, since it doesn't have a create operation.
				return nil, fmt.Errorf("fixRenameConflicts: couldn't find "+
					"merged path for %v", parent)
			}

			chain := unmergedChains.byOriginal[info.originalNewParent]
			unmergedPath, err = cr.getSingleUnmergedPath(
				ctx, unmergedChains, chain)
			if err != nil {
				return nil, err
			}
			// Look backwards to find the first parent with a merged path.
			n := len(unmergedPath.Path) - 1
			for i := n; i >= 0; i-- {
				mergedPath, ok = mergedPaths[unmergedPath.Path[i].BlockPointer]
				if ok {
					unmergedWalkBack = n - i
					break
				}
			}
			if !ok {
				return nil, fmt.Errorf("fixRenameConflicts: couldn't find any "+
					"merged path for any parents of %v", parent)
			}
		}

		for x, pn := range mergedPath.Path {
			original, err :=
				mergedChains.originalFromMostRecent(pn.BlockPointer)
			if err != nil {
				// This node wasn't changed in the merged branch
				original = pn.BlockPointer
			}

			if original != ptr {
				continue
			}

			// If any node on this path matches the renamed pointer,
			// we have a cycle.
			chain, ok := unmergedChains.byMostRecent[parent]
			if !ok {
				return nil, fmt.Errorf("fixRenameConflicts: no chain for "+
					"parent %v", parent)
			}

			j := len(mergedPath.Path) - 1
			// (j-x) in the above equation
			mergedWalkBack := j - x
			walkBack := unmergedWalkBack + mergedWalkBack

			// Mark this as a symlink, and the resolver
			// will take care of making it a symlink in
			// the merged branch later. No need to copy
			// since this createOp must have been created
			// as part of conflict resolution.
			symPath := "./" + strings.Repeat("../", walkBack)
			cr.log.CDebugf(ctx, "Creating symlink %s at "+
				"merged path %s", symPath, mergedPath)

			err = cr.convertCreateIntoSymlinkOrCopy(ctx, ptr, info, chain,
				unmergedChains, mergedChains, symPath)
			if err != nil {
				return nil, err
			}

			if unmergedWalkBack > 0 {
				cr.log.CDebugf(ctx, "Adding new unmerged path %s",
					unmergedPath)
				newUnmergedPaths = append(newUnmergedPaths,
					unmergedPath)
				// Fake a merged path to make sure these
				// actions will be taken.
				mergedLen := len(mergedPath.Path)
				pLen := mergedLen + unmergedWalkBack
				p := data.Path{
					FolderBranch: mergedPath.FolderBranch,
					Path:         make([]data.PathNode, pLen),
				}
				unmergedStart := len(unmergedPath.Path) -
					unmergedWalkBack
				copy(p.Path[:mergedLen], mergedPath.Path)
				copy(p.Path[mergedLen:],
					unmergedPath.Path[unmergedStart:])
				mergedPaths[unmergedPath.TailPointer()] = p
				if !p.IsValid() {
					// Temporary debugging for KBFS-2507.
					cr.log.CDebugf(ctx, "Added invalid unmerged path for %v",
						unmergedPath.TailPointer())
				}
			}

			removeRenames = append(removeRenames, ptr)
		}
	}

	// A map from merged most recent pointers of the parent
	// directories of files that have been forked, to a list of child
	// pointers within those directories that need their merged paths
	// fixed up.
	forkedFromMergedRenames := make(map[data.BlockPointer][]data.PathNode)

	// Check the merged renames to see if any of them affect a
	// modified file that the unmerged branch did not rename.  If we
	// find one, fork the file and leave the unmerged version under
	// its unmerged name.
	for ptr, info := range mergedChains.renamedOriginals {
		if mergedChains.isDeleted(ptr) {
			continue
		}

		// Skip double renames, already dealt with them above.
		if unmergedInfo, ok := unmergedChains.renamedOriginals[ptr]; ok &&
			(info.originalNewParent != unmergedInfo.originalNewParent ||
				info.newName != unmergedInfo.newName) {
			continue
		}

		// If this is a file that was modified in both branches, we
		// need to fork the file and tell the unmerged copy to keep
		// its current name.
		unmergedChain := unmergedChains.byOriginal[ptr]
		mergedChain := mergedChains.byOriginal[ptr]
		if crConflictCheckQuick(unmergedChain, mergedChain) {
			cr.log.CDebugf(ctx, "File that was renamed on the merged "+
				"branch from %s -> %s has conflicting edits, forking "+
				"(original ptr %v)", info.oldName, info.newName, ptr)
			var unmergedParentPath data.Path
			for _, op := range unmergedChain.ops {
				switch realOp := op.(type) {
				case *syncOp:
					realOp.keepUnmergedTailName = true
					unmergedParentPath = *op.getFinalPath().ParentPath()
				case *setAttrOp:
					realOp.keepUnmergedTailName = true
					unmergedParentPath = *op.getFinalPath().ParentPath()
				}
			}
			if unmergedParentPath.IsValid() {
				// Reset the merged path for this file back to the
				// merged path corresponding to the unmerged parent.
				// Put the merged parent path on the list of paths to
				// search for.
				unmergedParent := unmergedParentPath.TailPointer()
				if _, ok := mergedPaths[unmergedParent]; !ok {
					upOriginal := unmergedChains.originals[unmergedParent]
					mergedParent, err :=
						mergedChains.mostRecentFromOriginalOrSame(upOriginal)
					if err != nil {
						return nil, err
					}
					forkedFromMergedRenames[mergedParent] =
						append(forkedFromMergedRenames[mergedParent],
							data.PathNode{
								BlockPointer: unmergedChain.mostRecent,
								Name:         info.oldName,
							})
					newUnmergedPaths =
						append(newUnmergedPaths, unmergedParentPath)
				}
			}
		}
	}

	for _, ptr := range removeRenames {
		delete(unmergedChains.renamedOriginals, ptr)
	}

	numRenamesToCheck := len(doubleRenames) + len(forkedFromMergedRenames)
	if numRenamesToCheck == 0 {
		return newUnmergedPaths, nil
	}

	// Make chains for the new merged parents of all the double renames.
	newPtrs := make(map[data.BlockPointer]bool)
	ptrs := make([]data.BlockPointer, len(doubleRenames), numRenamesToCheck)
	copy(ptrs, doubleRenames)
	for ptr := range forkedFromMergedRenames {
		ptrs = append(ptrs, ptr)
	}
	// Fake out the rest of the chains to populate newPtrs
	for ptr := range mergedChains.byMostRecent {
		newPtrs[ptr] = true
	}

	mergedNodeCache := newNodeCacheStandard(cr.fbo.folderBranch)
	nodeMap, _, err := cr.fbo.blocks.SearchForNodes(
		ctx, mergedNodeCache, ptrs, newPtrs,
		mergedChains.mostRecentChainMDInfo,
		mergedChains.mostRecentChainMDInfo.GetRootDirEntry().BlockPointer)
	if err != nil {
		return nil, err
	}

	for _, ptr := range doubleRenames {
		// Find the merged paths
		node, ok := nodeMap[ptr]
		if !ok || node == nil {
			return nil, fmt.Errorf("Couldn't find merged path for "+
				"doubly-renamed pointer %v", ptr)
		}

		original, err :=
			mergedChains.originalFromMostRecentOrSame(ptr)
		if err != nil {
			return nil, err
		}
		unmergedInfo, ok := unmergedChains.renamedOriginals[original]
		if !ok {
			return nil, fmt.Errorf("fixRenameConflicts: can't find the "+
				"unmerged rename info for %v during double-rename resolution",
				original)
		}
		mergedInfo, ok := mergedChains.renamedOriginals[original]
		if !ok {
			return nil, fmt.Errorf("fixRenameConflicts: can't find the "+
				"merged rename info for %v during double-rename resolution",
				original)
		}

		// If any node on this path matches the renamed pointer,
		// we have a cycle.
		chain, ok := unmergedChains.byOriginal[unmergedInfo.originalNewParent]
		if !ok {
			return nil, fmt.Errorf("fixRenameConflicts: no chain for "+
				"parent %v", unmergedInfo.originalNewParent)
		}

		// For directories, the symlinks traverse down the merged path
		// to the first common node, and then back up to the new
		// parent/name.  TODO: what happens when some element along
		// the merged path also got renamed by the unmerged branch?
		// The symlink would likely be wrong in that case.
		mergedPathOldParent, ok := mergedPaths[chain.mostRecent]
		if !ok {
			return nil, fmt.Errorf("fixRenameConflicts: couldn't find "+
				"merged path for old parent %v", chain.mostRecent)
		}
		mergedPathNewParent := mergedNodeCache.PathFromNode(node)
		symPath := "./"
		newParentStart := 0
	outer:
		for i := len(mergedPathOldParent.Path) - 1; i >= 0; i-- {
			mostRecent := mergedPathOldParent.Path[i].BlockPointer
			for j, pnode := range mergedPathNewParent.Path {
				original, err :=
					unmergedChains.originalFromMostRecentOrSame(mostRecent)
				if err != nil {
					return nil, err
				}
				mergedMostRecent, err :=
					mergedChains.mostRecentFromOriginalOrSame(original)
				if err != nil {
					return nil, err
				}
				if pnode.BlockPointer == mergedMostRecent {
					newParentStart = j
					break outer
				}
			}
			symPath += "../"
		}
		// Move up directories starting from beyond the common parent,
		// to right before the actual node.
		for i := newParentStart + 1; i < len(mergedPathNewParent.Path)-1; i++ {
			symPath += mergedPathNewParent.Path[i].Name + "/"
		}
		symPath += mergedInfo.newName

		err = cr.convertCreateIntoSymlinkOrCopy(ctx, original, unmergedInfo,
			chain, unmergedChains, mergedChains, symPath)
		if err != nil {
			return nil, err
		}
	}

	for ptr, pathNodes := range forkedFromMergedRenames {
		// Find the merged paths
		node, ok := nodeMap[ptr]
		if !ok || node == nil {
			return nil, fmt.Errorf("Couldn't find merged path for "+
				"forked parent pointer %v", ptr)
		}

		mergedPathNewParent := mergedNodeCache.PathFromNode(node)
		for _, pNode := range pathNodes {
			mergedPath := mergedPathNewParent.ChildPath(
				pNode.Name, pNode.BlockPointer)
			mergedPaths[pNode.BlockPointer] = mergedPath
		}
	}

	return newUnmergedPaths, nil
}
