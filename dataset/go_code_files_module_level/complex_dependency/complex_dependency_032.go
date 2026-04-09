func (fd *FileData) Write(ctx context.Context, data []byte, off Int64Offset,
	topBlock *FileBlock, oldDe DirEntry, df *DirtyFile) (
	newDe DirEntry, dirtyPtrs []BlockPointer, unrefs []BlockInfo,
	newlyDirtiedChildBytes int64, bytesExtended int64, err error) {
	n := int64(len(data))
	nCopied := int64(0)
	oldSizeWithoutHoles := oldDe.Size
	newDe = oldDe

	fd.tree.vlog.CLogf(ctx, libkb.VLog1, "Writing %d bytes at off %d", n, off)

	dirtyMap := make(map[BlockPointer]bool)
	for nCopied < n {
		ptr, parentBlocks, block, nextBlockOff, startOff, wasDirty, err :=
			fd.GetFileBlockAtOffset(
				ctx, topBlock, off+Int64Offset(nCopied), BlockWrite)
		if err != nil {
			return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
		}

		oldLen := len(block.Contents)

		// Take care not to write past the beginning of the next block
		// by using max.
		max := Int64Offset(len(data))
		if nextBlockOff > 0 {
			if room := nextBlockOff - off; room < max {
				max = room
			}
		}
		oldNCopied := nCopied
		nCopied += fd.tree.bsplit.CopyUntilSplit(
			block, nextBlockOff < 0, data[nCopied:max],
			int64(off+Int64Offset(nCopied)-startOff))

		// If we need another block but there are no more, then make one.
		switchToIndirect := false
		if nCopied < n {
			needExtendFile := nextBlockOff < 0
			needFillHole := off+Int64Offset(nCopied) < nextBlockOff
			newBlockOff := startOff + Int64Offset(len(block.Contents))
			if nCopied == 0 {
				if newBlockOff < off {
					// We are writing past the end of a file, or
					// somewhere inside a hole, not right at the start
					// of it; all we have done so far it reached the
					// end of an existing block (possibly zero-filling
					// it out to its capacity).  Make sure the next
					// block starts right at the offset we care about.
					newBlockOff = off
				}
			} else if newBlockOff != off+Int64Offset(nCopied) {
				return newDe, nil, unrefs, newlyDirtiedChildBytes, 0,
					fmt.Errorf("Copied %d bytes, but newBlockOff=%d does not "+
						"match off=%d plus new bytes",
						nCopied, newBlockOff, off)
			}
			var rightParents []ParentBlockAndChildIndex
			if needExtendFile || needFillHole {
				// Make a new right block and update the parent's
				// indirect block list, adding a level of indirection
				// if needed.  If we're just filling a hole, the block
				// will end up all the way to the right of the range,
				// and its offset will be smaller than the block to
				// its left -- we'll fix that up below.
				var newDirtyPtrs []BlockPointer
				fd.tree.vlog.CLogf(
					ctx, libkb.VLog1, "Making new right block at "+
						"nCopied=%d, newBlockOff=%d", nCopied, newBlockOff)
				wasIndirect := topBlock.IsInd
				rightParents, newDirtyPtrs, err = fd.tree.newRightBlock(
					ctx, parentBlocks, newBlockOff,
					DefaultNewBlockDataVersion(false), NewFileBlockWithPtrs,
					fd.fileTopBlocker(df),
				)
				if err != nil {
					return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
				}
				topBlock = rightParents[0].pblock.(*FileBlock)
				for _, p := range newDirtyPtrs {
					dirtyMap[p] = true
				}
				if topBlock.IsInd != wasIndirect {
					// The whole direct data block needs to be
					// re-uploaded as a child block with a new block
					// pointer, so below we'll need to track the dirty
					// bytes of the direct block and cache the block
					// as dirty.  (Note that currently we don't track
					// dirty bytes for indirect blocks.)
					switchToIndirect = true
					ptr = topBlock.IPtrs[0].BlockPointer
				}
			}
			// If we're filling a hole, swap the new right block into
			// the hole and shift everything else over.
			if needFillHole {
				newDirtyPtrs, newUnrefs, bytes, err :=
					fd.tree.shiftBlocksToFillHole(ctx, rightParents)
				if err != nil {
					return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
				}
				for _, p := range newDirtyPtrs {
					dirtyMap[p] = true
				}
				unrefs = append(unrefs, newUnrefs...)
				newlyDirtiedChildBytes += bytes
				if oldSizeWithoutHoles == oldDe.Size {
					// For the purposes of calculating the newly-dirtied
					// bytes for the deferral calculation, disregard the
					// existing "hole" in the file.
					oldSizeWithoutHoles = uint64(newBlockOff)
				}
			}
		}

		// Nothing was copied, no need to dirty anything.  This can
		// happen when trying to append to the contents of the file
		// (i.e., either to the end of the file or right before the
		// "hole"), and the last block is already full.
		if nCopied == oldNCopied && oldLen == len(block.Contents) &&
			!switchToIndirect {
			continue
		}

		// Only in the last block does the file size grow.
		if oldLen != len(block.Contents) && nextBlockOff < 0 {
			newDe.EncodedSize = 0
			// Since this is the last block, the end of this block
			// marks the file size.
			newDe.Size = uint64(startOff + Int64Offset(len(block.Contents)))
		}

		// Calculate the amount of bytes we've newly-dirtied as part
		// of this write.
		newlyDirtiedChildBytes += int64(len(block.Contents))
		if wasDirty {
			newlyDirtiedChildBytes -= int64(oldLen)
		}

		newDirtyPtrs, newUnrefs, err := fd.tree.markParentsDirty(
			ctx, parentBlocks)
		unrefs = append(unrefs, newUnrefs...)
		if err != nil {
			return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
		}
		for _, p := range newDirtyPtrs {
			dirtyMap[p] = true
		}

		// keep the old block ID while it's dirty
		if err = fd.tree.cacher(ctx, ptr, block); err != nil {
			return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
		}
		dirtyMap[ptr] = true
	}

	// Always make the top block dirty, so we will sync any indirect
	// blocks.  This has the added benefit of ensuring that any write
	// to a file while it's being sync'd will be deferred, even if
	// it's to a block that's not currently being sync'd, since this
	// top-most block will always be in the dirtyFiles map.  We do
	// this even for 0-byte writes, which indicate a forced sync.
	if err = fd.tree.cacher(ctx, fd.rootBlockPointer(), topBlock); err != nil {
		return newDe, nil, unrefs, newlyDirtiedChildBytes, 0, err
	}
	dirtyMap[fd.rootBlockPointer()] = true

	lastByteWritten := int64(off) + int64(len(data)) // not counting holes
	bytesExtended = 0
	if lastByteWritten > int64(oldSizeWithoutHoles) {
		bytesExtended = lastByteWritten - int64(oldSizeWithoutHoles)
	}

	dirtyPtrs = make([]BlockPointer, 0, len(dirtyMap))
	for p := range dirtyMap {
		dirtyPtrs = append(dirtyPtrs, p)
	}

	return newDe, dirtyPtrs, unrefs, newlyDirtiedChildBytes, bytesExtended, nil
}
