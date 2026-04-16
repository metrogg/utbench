func (bt *blockTree) shiftBlocksToFillHole(
	ctx context.Context, parents []ParentBlockAndChildIndex) (
	newDirtyPtrs []BlockPointer, newUnrefs []BlockInfo,
	newlyDirtiedChildBytes int64, err error) {
	// `parents` should represent the right side of the tree down to
	// the new rightmost indirect pointer, the offset of which should
	// match `newHoleStartOff`.  Keep swapping it with its sibling on
	// the left until its offset would be lower than that child's
	// offset.  If there are no children to the left, continue on with
	// the children in the cousin block to the left.  If we swap a
	// child between cousin blocks, we must update the offset in the
	// right cousin's parent block.  If *that* updated pointer is the
	// leftmost pointer in its parent block, update that one as well,
	// up to the root.
	//
	// We are guaranteed at least one level of indirection because
	// `newRightBlock` should have been called before
	// `shiftBlocksToFillHole`.
	immedParent := parents[len(parents)-1]
	currIndex := immedParent.childIndex
	_, newBlockStartOff := immedParent.childIPtr()

	bt.vlog.CLogf(
		ctx, libkb.VLog1, "Shifting block with offset %s for entry %v into "+
			"position", newBlockStartOff, bt.rootBlockPointer())

	// Swap left as needed.
	for loopedOnce := false; ; loopedOnce = true {
		var leftOff Offset
		var newParents []ParentBlockAndChildIndex
		immedPblock := immedParent.pblock
		if currIndex > 0 {
			_, leftOff = immedPblock.IndirectPtr(currIndex - 1)
		} else {
			if loopedOnce {
				// Now update the left side if needed, before looking into
				// swapping across blocks.
				bt.vlog.CLogf(ctx, libkb.VLog1, "Updating on left side")
				_, newOff := immedPblock.IndirectPtr(currIndex)
				ndp, nu, err := bt.setParentOffsets(
					ctx, newOff, parents, currIndex)
				if err != nil {
					return nil, nil, 0, err
				}
				newDirtyPtrs = append(newDirtyPtrs, ndp...)
				newUnrefs = append(newUnrefs, nu...)
			}

			// Construct the new set of parents for the shifted block,
			// by looking for the next left cousin.
			newParents = make([]ParentBlockAndChildIndex, len(parents))
			copy(newParents, parents)
			var level int
			for level = len(newParents) - 2; level >= 0; level-- {
				// The parent at the level being evaluated has a left
				// sibling, so we use that sibling.
				if newParents[level].childIndex > 0 {
					break
				}
				// Keep going up until we find a way back down a left branch.
			}

			if level < 0 {
				// We are already all the way on the left, we're done!
				return newDirtyPtrs, newUnrefs, newlyDirtiedChildBytes, nil
			}
			newParents[level].childIndex--

			// Walk back down, shifting the new parents into position.
			for ; level < len(newParents)-1; level++ {
				nextPtr := newParents[level].childBlockPtr()
				childBlock, _, err := bt.getter(
					ctx, bt.kmd, nextPtr, bt.file, BlockWrite)
				if err != nil {
					return nil, nil, 0, err
				}

				newParents[level+1].pblock = childBlock
				newParents[level+1].childIndex =
					childBlock.NumIndirectPtrs() - 1
				_, leftOff = childBlock.IndirectPtr(
					childBlock.NumIndirectPtrs() - 1)
			}
		}

		// We're done!
		if leftOff.Less(newBlockStartOff) {
			return newDirtyPtrs, newUnrefs, newlyDirtiedChildBytes, nil
		}

		// Otherwise, we need to swap the indirect file pointers.
		if currIndex > 0 {
			immedPblock.SwapIndirectPtrs(currIndex-1, immedPblock, currIndex)
			currIndex--
			continue
		}

		// Swap block pointers across cousins at the lowest level of
		// indirection.
		newImmedParent := newParents[len(newParents)-1]
		newImmedPblock := newImmedParent.pblock
		newCurrIndex := newImmedPblock.NumIndirectPtrs() - 1
		newImmedPblock.SwapIndirectPtrs(newCurrIndex, immedPblock, currIndex)

		// Cache the new immediate parent as dirty.  Also cache the
		// old immediate parent's right-most leaf child as dirty, to
		// make sure this path is captured in
		// getNextDirtyBlockAtOffset calls.  TODO: this is inefficient
		// since it might end up re-encoding and re-uploading a leaf
		// block that wasn't actually dirty; we should find a better
		// way to make sure ready() sees these parent blocks.
		if len(newParents) > 1 {
			i := len(newParents) - 2
			childPtr := newParents[i].childBlockPtr()
			if err := bt.cacher(
				ctx, childPtr, newImmedPblock); err != nil {
				return nil, nil, 0, err
			}
			newDirtyPtrs = append(newDirtyPtrs, childPtr)

			// Fetch the old parent's right leaf for writing, and mark
			// it as dirty.
			rightLeafInfo, _ := immedPblock.IndirectPtr(
				immedPblock.NumIndirectPtrs() - 1)
			leafBlock, _, err := bt.getter(
				ctx, bt.kmd, rightLeafInfo.BlockPointer, bt.file, BlockWrite)
			if err != nil {
				return nil, nil, 0, err
			}
			if err := bt.cacher(
				ctx, rightLeafInfo.BlockPointer, leafBlock); err != nil {
				return nil, nil, 0, err
			}
			newDirtyPtrs = append(newDirtyPtrs, rightLeafInfo.BlockPointer)
			// Remember the size of the dirtied leaf.
			if rightLeafInfo.EncodedSize != 0 {
				newlyDirtiedChildBytes += leafBlock.BytesCanBeDirtied()
				newUnrefs = append(newUnrefs, rightLeafInfo)
				immedPblock.ClearIndirectPtrSize(
					immedPblock.NumIndirectPtrs() - 1)
			}
		}

		// Now we need to update the parent offsets on the right side,
		// all the way up to the common ancestor (which is the one
		// with the one that doesn't have a childIndex of 0).
		_, newRightOff := immedPblock.IndirectPtr(currIndex)
		ndp, nu, err := bt.setParentOffsets(
			ctx, newRightOff, parents, currIndex)
		if err != nil {
			return nil, nil, 0, err
		}
		newDirtyPtrs = append(newDirtyPtrs, ndp...)
		newUnrefs = append(newUnrefs, nu...)

		immedParent = newImmedParent
		currIndex = newCurrIndex
		parents = newParents
	}
	// The loop above must exit via one of the returns.
}
