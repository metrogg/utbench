func convertInMessage(
	inMsg *buffer.InMessage,
	outMsg *buffer.OutMessage,
	protocol fusekernel.Protocol) (o interface{}, err error) {
	switch inMsg.Header().Opcode {
	case fusekernel.OpLookup:
		buf := inMsg.ConsumeBytes(inMsg.Len())
		n := len(buf)
		if n == 0 || buf[n-1] != '\x00' {
			err = errors.New("Corrupt OpLookup")
			return
		}

		o = &fuseops.LookUpInodeOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(buf[:n-1]),
		}

	case fusekernel.OpGetattr:
		o = &fuseops.GetInodeAttributesOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}

	case fusekernel.OpSetattr:
		type input fusekernel.SetattrIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpSetattr")
			return
		}

		to := &fuseops.SetInodeAttributesOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}
		o = to

		valid := fusekernel.SetattrValid(in.Valid)
		if valid&fusekernel.SetattrSize != 0 {
			to.Size = &in.Size
		}

		if valid&fusekernel.SetattrMode != 0 {
			mode := convertFileMode(in.Mode)
			to.Mode = &mode
		}

		if valid&fusekernel.SetattrAtime != 0 {
			t := time.Unix(int64(in.Atime), int64(in.AtimeNsec))
			to.Atime = &t
		}

		if valid&fusekernel.SetattrMtime != 0 {
			t := time.Unix(int64(in.Mtime), int64(in.MtimeNsec))
			to.Mtime = &t
		}

	case fusekernel.OpForget:
		type input fusekernel.ForgetIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpForget")
			return
		}

		o = &fuseops.ForgetInodeOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
			N:     in.Nlookup,
		}

	case fusekernel.OpMkdir:
		in := (*fusekernel.MkdirIn)(inMsg.Consume(fusekernel.MkdirInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpMkdir")
			return
		}

		name := inMsg.ConsumeBytes(inMsg.Len())
		i := bytes.IndexByte(name, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpMkdir")
			return
		}
		name = name[:i]

		o = &fuseops.MkDirOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(name),

			// On Linux, vfs_mkdir calls through to the inode with at most
			// permissions and sticky bits set (cf. https://goo.gl/WxgQXk), and fuse
			// passes that on directly (cf. https://goo.gl/f31aMo). In other words,
			// the fact that this is a directory is implicit in the fact that the
			// opcode is mkdir. But we want the correct mode to go through, so ensure
			// that os.ModeDir is set.
			Mode: convertFileMode(in.Mode) | os.ModeDir,
		}

	case fusekernel.OpMknod:
		in := (*fusekernel.MknodIn)(inMsg.Consume(fusekernel.MknodInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpMknod")
			return
		}

		name := inMsg.ConsumeBytes(inMsg.Len())
		i := bytes.IndexByte(name, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpMknod")
			return
		}
		name = name[:i]

		o = &fuseops.MkNodeOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(name),
			Mode:   convertFileMode(in.Mode),
		}

	case fusekernel.OpCreate:
		in := (*fusekernel.CreateIn)(inMsg.Consume(fusekernel.CreateInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpCreate")
			return
		}

		name := inMsg.ConsumeBytes(inMsg.Len())
		i := bytes.IndexByte(name, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpCreate")
			return
		}
		name = name[:i]

		o = &fuseops.CreateFileOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(name),
			Mode:   convertFileMode(in.Mode),
		}

	case fusekernel.OpSymlink:
		// The message is "newName\0target\0".
		names := inMsg.ConsumeBytes(inMsg.Len())
		if len(names) == 0 || names[len(names)-1] != 0 {
			err = errors.New("Corrupt OpSymlink")
			return
		}
		i := bytes.IndexByte(names, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpSymlink")
			return
		}
		newName, target := names[0:i], names[i+1:len(names)-1]

		o = &fuseops.CreateSymlinkOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(newName),
			Target: string(target),
		}

	case fusekernel.OpRename:
		type input fusekernel.RenameIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpRename")
			return
		}

		names := inMsg.ConsumeBytes(inMsg.Len())
		// names should be "old\x00new\x00"
		if len(names) < 4 {
			err = errors.New("Corrupt OpRename")
			return
		}
		if names[len(names)-1] != '\x00' {
			err = errors.New("Corrupt OpRename")
			return
		}
		i := bytes.IndexByte(names, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpRename")
			return
		}
		oldName, newName := names[:i], names[i+1:len(names)-1]

		o = &fuseops.RenameOp{
			OldParent: fuseops.InodeID(inMsg.Header().Nodeid),
			OldName:   string(oldName),
			NewParent: fuseops.InodeID(in.Newdir),
			NewName:   string(newName),
		}

	case fusekernel.OpUnlink:
		buf := inMsg.ConsumeBytes(inMsg.Len())
		n := len(buf)
		if n == 0 || buf[n-1] != '\x00' {
			err = errors.New("Corrupt OpUnlink")
			return
		}

		o = &fuseops.UnlinkOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(buf[:n-1]),
		}

	case fusekernel.OpRmdir:
		buf := inMsg.ConsumeBytes(inMsg.Len())
		n := len(buf)
		if n == 0 || buf[n-1] != '\x00' {
			err = errors.New("Corrupt OpRmdir")
			return
		}

		o = &fuseops.RmDirOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(buf[:n-1]),
		}

	case fusekernel.OpOpen:
		o = &fuseops.OpenFileOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}

	case fusekernel.OpOpendir:
		o = &fuseops.OpenDirOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}

	case fusekernel.OpRead:
		in := (*fusekernel.ReadIn)(inMsg.Consume(fusekernel.ReadInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpRead")
			return
		}

		to := &fuseops.ReadFileOp{
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
			Handle: fuseops.HandleID(in.Fh),
			Offset: int64(in.Offset),
		}
		o = to

		readSize := int(in.Size)
		p := outMsg.GrowNoZero(readSize)
		if p == nil {
			err = fmt.Errorf("Can't grow for %d-byte read", readSize)
			return
		}

		sh := (*reflect.SliceHeader)(unsafe.Pointer(&to.Dst))
		sh.Data = uintptr(p)
		sh.Len = readSize
		sh.Cap = readSize

	case fusekernel.OpReaddir:
		in := (*fusekernel.ReadIn)(inMsg.Consume(fusekernel.ReadInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpReaddir")
			return
		}

		to := &fuseops.ReadDirOp{
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
			Handle: fuseops.HandleID(in.Fh),
			Offset: fuseops.DirOffset(in.Offset),
		}
		o = to

		readSize := int(in.Size)
		p := outMsg.GrowNoZero(readSize)
		if p == nil {
			err = fmt.Errorf("Can't grow for %d-byte read", readSize)
			return
		}

		sh := (*reflect.SliceHeader)(unsafe.Pointer(&to.Dst))
		sh.Data = uintptr(p)
		sh.Len = readSize
		sh.Cap = readSize

	case fusekernel.OpRelease:
		type input fusekernel.ReleaseIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpRelease")
			return
		}

		o = &fuseops.ReleaseFileHandleOp{
			Handle: fuseops.HandleID(in.Fh),
		}

	case fusekernel.OpReleasedir:
		type input fusekernel.ReleaseIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpReleasedir")
			return
		}

		o = &fuseops.ReleaseDirHandleOp{
			Handle: fuseops.HandleID(in.Fh),
		}

	case fusekernel.OpWrite:
		in := (*fusekernel.WriteIn)(inMsg.Consume(fusekernel.WriteInSize(protocol)))
		if in == nil {
			err = errors.New("Corrupt OpWrite")
			return
		}

		buf := inMsg.ConsumeBytes(inMsg.Len())
		if len(buf) < int(in.Size) {
			err = errors.New("Corrupt OpWrite")
			return
		}

		o = &fuseops.WriteFileOp{
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
			Handle: fuseops.HandleID(in.Fh),
			Data:   buf,
			Offset: int64(in.Offset),
		}

	case fusekernel.OpFsync:
		type input fusekernel.FsyncIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpFsync")
			return
		}

		o = &fuseops.SyncFileOp{
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
			Handle: fuseops.HandleID(in.Fh),
		}

	case fusekernel.OpFlush:
		type input fusekernel.FlushIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpFlush")
			return
		}

		o = &fuseops.FlushFileOp{
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
			Handle: fuseops.HandleID(in.Fh),
		}

	case fusekernel.OpReadlink:
		o = &fuseops.ReadSymlinkOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}

	case fusekernel.OpStatfs:
		o = &fuseops.StatFSOp{}

	case fusekernel.OpInterrupt:
		type input fusekernel.InterruptIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpInterrupt")
			return
		}

		o = &interruptOp{
			FuseID: in.Unique,
		}

	case fusekernel.OpInit:
		type input fusekernel.InitIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpInit")
			return
		}

		o = &initOp{
			Kernel:       fusekernel.Protocol{in.Major, in.Minor},
			MaxReadahead: in.MaxReadahead,
			Flags:        fusekernel.InitFlags(in.Flags),
		}

	case fusekernel.OpLink:
		type input fusekernel.LinkIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpLink")
			return
		}

		name := inMsg.ConsumeBytes(inMsg.Len())
		i := bytes.IndexByte(name, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpLink")
			return
		}
		name = name[:i]
		if len(name) == 0 {
			err = errors.New("Corrupt OpLink (Name not read)")
			return
		}

		o = &fuseops.CreateLinkOp{
			Parent: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:   string(name),
			Target: fuseops.InodeID(in.Oldnodeid),
		}

	case fusekernel.OpRemovexattr:
		buf := inMsg.ConsumeBytes(inMsg.Len())
		n := len(buf)
		if n == 0 || buf[n-1] != '\x00' {
			err = errors.New("Corrupt OpRemovexattr")
			return
		}

		o = &fuseops.RemoveXattrOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:  string(buf[:n-1]),
		}

	case fusekernel.OpGetxattr:
		type input fusekernel.GetxattrIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpGetxattr")
			return
		}

		name := inMsg.ConsumeBytes(inMsg.Len())
		i := bytes.IndexByte(name, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpGetxattr")
			return
		}
		name = name[:i]

		to := &fuseops.GetXattrOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:  string(name),
		}
		o = to

		readSize := int(in.Size)
		p := outMsg.GrowNoZero(readSize)
		if p == nil {
			err = fmt.Errorf("Can't grow for %d-byte read", readSize)
			return
		}

		sh := (*reflect.SliceHeader)(unsafe.Pointer(&to.Dst))
		sh.Data = uintptr(p)
		sh.Len = readSize
		sh.Cap = readSize

	case fusekernel.OpListxattr:
		type input fusekernel.ListxattrIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpListxattr")
			return
		}

		to := &fuseops.ListXattrOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
		}
		o = to

		readSize := int(in.Size)
		if readSize != 0 {
			p := outMsg.GrowNoZero(readSize)
			if p == nil {
				err = fmt.Errorf("Can't grow for %d-byte read", readSize)
				return
			}
			sh := (*reflect.SliceHeader)(unsafe.Pointer(&to.Dst))
			sh.Data = uintptr(p)
			sh.Len = readSize
			sh.Cap = readSize
		}
	case fusekernel.OpSetxattr:
		type input fusekernel.SetxattrIn
		in := (*input)(inMsg.Consume(unsafe.Sizeof(input{})))
		if in == nil {
			err = errors.New("Corrupt OpSetxattr")
			return
		}

		payload := inMsg.ConsumeBytes(inMsg.Len())
		// payload should be "name\x00value"
		if len(payload) < 3 {
			err = errors.New("Corrupt OpSetxattr")
			return
		}
		i := bytes.IndexByte(payload, '\x00')
		if i < 0 {
			err = errors.New("Corrupt OpSetxattr")
			return
		}

		name, value := payload[:i], payload[i+1:len(payload)]

		o = &fuseops.SetXattrOp{
			Inode: fuseops.InodeID(inMsg.Header().Nodeid),
			Name:  string(name),
			Value: value,
			Flags: in.Flags,
		}

	default:
		o = &unknownOp{
			OpCode: inMsg.Header().Opcode,
			Inode:  fuseops.InodeID(inMsg.Header().Nodeid),
		}
	}

	return
}
