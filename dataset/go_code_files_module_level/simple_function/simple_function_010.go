func WriteElement(w io.Writer, element interface{}) error {
	switch e := element.(type) {
	case NodeAlias:
		if _, err := w.Write(e[:]); err != nil {
			return err
		}

	case ShortChanIDEncoding:
		var b [1]byte
		b[0] = uint8(e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case uint8:
		var b [1]byte
		b[0] = e
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case FundingFlag:
		var b [1]byte
		b[0] = uint8(e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case uint16:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case ChanUpdateMsgFlags:
		var b [1]byte
		b[0] = uint8(e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case ChanUpdateChanFlags:
		var b [1]byte
		b[0] = uint8(e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case ErrorCode:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(e))
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case MilliSatoshi:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(e))
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case btcutil.Amount:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(e))
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case uint32:
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case uint64:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], e)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case *btcec.PublicKey:
		if e == nil {
			return fmt.Errorf("cannot write nil pubkey")
		}

		var b [33]byte
		serializedPubkey := e.SerializeCompressed()
		copy(b[:], serializedPubkey)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	case []Sig:
		var b [2]byte
		numSigs := uint16(len(e))
		binary.BigEndian.PutUint16(b[:], numSigs)
		if _, err := w.Write(b[:]); err != nil {
			return err
		}

		for _, sig := range e {
			if err := WriteElement(w, sig); err != nil {
				return err
			}
		}
	case Sig:
		// Write buffer
		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case PingPayload:
		var l [2]byte
		binary.BigEndian.PutUint16(l[:], uint16(len(e)))
		if _, err := w.Write(l[:]); err != nil {
			return err
		}

		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case PongPayload:
		var l [2]byte
		binary.BigEndian.PutUint16(l[:], uint16(len(e)))
		if _, err := w.Write(l[:]); err != nil {
			return err
		}

		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case ErrorData:
		var l [2]byte
		binary.BigEndian.PutUint16(l[:], uint16(len(e)))
		if _, err := w.Write(l[:]); err != nil {
			return err
		}

		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case OpaqueReason:
		var l [2]byte
		binary.BigEndian.PutUint16(l[:], uint16(len(e)))
		if _, err := w.Write(l[:]); err != nil {
			return err
		}

		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case [33]byte:
		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case []byte:
		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case PkScript:
		// The largest script we'll accept is a p2wsh which is exactly
		// 34 bytes long.
		scriptLength := len(e)
		if scriptLength > 34 {
			return fmt.Errorf("'PkScript' too long")
		}

		if err := wire.WriteVarBytes(w, 0, e); err != nil {
			return err
		}
	case *RawFeatureVector:
		if e == nil {
			return fmt.Errorf("cannot write nil feature vector")
		}

		if err := e.Encode(w); err != nil {
			return err
		}

	case wire.OutPoint:
		var h [32]byte
		copy(h[:], e.Hash[:])
		if _, err := w.Write(h[:]); err != nil {
			return err
		}

		if e.Index > math.MaxUint16 {
			return fmt.Errorf("index for outpoint (%v) is "+
				"greater than max index of %v", e.Index,
				math.MaxUint16)
		}

		var idx [2]byte
		binary.BigEndian.PutUint16(idx[:], uint16(e.Index))
		if _, err := w.Write(idx[:]); err != nil {
			return err
		}

	case ChannelID:
		if _, err := w.Write(e[:]); err != nil {
			return err
		}
	case FailCode:
		if err := WriteElement(w, uint16(e)); err != nil {
			return err
		}
	case ShortChannelID:
		// Check that field fit in 3 bytes and write the blockHeight
		if e.BlockHeight > ((1 << 24) - 1) {
			return errors.New("block height should fit in 3 bytes")
		}

		var blockHeight [4]byte
		binary.BigEndian.PutUint32(blockHeight[:], e.BlockHeight)

		if _, err := w.Write(blockHeight[1:]); err != nil {
			return err
		}

		// Check that field fit in 3 bytes and write the txIndex
		if e.TxIndex > ((1 << 24) - 1) {
			return errors.New("tx index should fit in 3 bytes")
		}

		var txIndex [4]byte
		binary.BigEndian.PutUint32(txIndex[:], e.TxIndex)
		if _, err := w.Write(txIndex[1:]); err != nil {
			return err
		}

		// Write the txPosition
		var txPosition [2]byte
		binary.BigEndian.PutUint16(txPosition[:], e.TxPosition)
		if _, err := w.Write(txPosition[:]); err != nil {
			return err
		}

	case *net.TCPAddr:
		if e == nil {
			return fmt.Errorf("cannot write nil TCPAddr")
		}

		if e.IP.To4() != nil {
			var descriptor [1]byte
			descriptor[0] = uint8(tcp4Addr)
			if _, err := w.Write(descriptor[:]); err != nil {
				return err
			}

			var ip [4]byte
			copy(ip[:], e.IP.To4())
			if _, err := w.Write(ip[:]); err != nil {
				return err
			}
		} else {
			var descriptor [1]byte
			descriptor[0] = uint8(tcp6Addr)
			if _, err := w.Write(descriptor[:]); err != nil {
				return err
			}
			var ip [16]byte
			copy(ip[:], e.IP.To16())
			if _, err := w.Write(ip[:]); err != nil {
				return err
			}
		}
		var port [2]byte
		binary.BigEndian.PutUint16(port[:], uint16(e.Port))
		if _, err := w.Write(port[:]); err != nil {
			return err
		}

	case *tor.OnionAddr:
		if e == nil {
			return errors.New("cannot write nil onion address")
		}

		var suffixIndex int
		switch len(e.OnionService) {
		case tor.V2Len:
			descriptor := []byte{byte(v2OnionAddr)}
			if _, err := w.Write(descriptor); err != nil {
				return err
			}
			suffixIndex = tor.V2Len - tor.OnionSuffixLen
		case tor.V3Len:
			descriptor := []byte{byte(v3OnionAddr)}
			if _, err := w.Write(descriptor); err != nil {
				return err
			}
			suffixIndex = tor.V3Len - tor.OnionSuffixLen
		default:
			return errors.New("unknown onion service length")
		}

		host, err := tor.Base32Encoding.DecodeString(
			e.OnionService[:suffixIndex],
		)
		if err != nil {
			return err
		}
		if _, err := w.Write(host); err != nil {
			return err
		}

		var port [2]byte
		binary.BigEndian.PutUint16(port[:], uint16(e.Port))
		if _, err := w.Write(port[:]); err != nil {
			return err
		}

	case []net.Addr:
		// First, we'll encode all the addresses into an intermediate
		// buffer. We need to do this in order to compute the total
		// length of the addresses.
		var addrBuf bytes.Buffer
		for _, address := range e {
			if err := WriteElement(&addrBuf, address); err != nil {
				return err
			}
		}

		// With the addresses fully encoded, we can now write out the
		// number of bytes needed to encode them.
		addrLen := addrBuf.Len()
		if err := WriteElement(w, uint16(addrLen)); err != nil {
			return err
		}

		// Finally, we'll write out the raw addresses themselves, but
		// only if we have any bytes to write.
		if addrLen > 0 {
			if _, err := w.Write(addrBuf.Bytes()); err != nil {
				return err
			}
		}
	case color.RGBA:
		if err := WriteElements(w, e.R, e.G, e.B); err != nil {
			return err
		}

	case DeliveryAddress:
		var length [2]byte
		binary.BigEndian.PutUint16(length[:], uint16(len(e)))
		if _, err := w.Write(length[:]); err != nil {
			return err
		}
		if _, err := w.Write(e[:]); err != nil {
			return err
		}

	case bool:
		var b [1]byte
		if e {
			b[0] = 1
		}
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown type in WriteElement: %T", e)
	}

	return nil
}
