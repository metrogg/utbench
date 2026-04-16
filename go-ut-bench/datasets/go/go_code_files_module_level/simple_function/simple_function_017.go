func ReadElement(r io.Reader, element interface{}) error {
	var err error
	switch e := element.(type) {
	case *bool:
		var b [1]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}

		if b[0] == 1 {
			*e = true
		}

	case *NodeAlias:
		var a [32]byte
		if _, err := io.ReadFull(r, a[:]); err != nil {
			return err
		}

		alias, err := NewNodeAlias(string(a[:]))
		if err != nil {
			return err
		}

		*e = alias
	case *ShortChanIDEncoding:
		var b [1]uint8
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		*e = ShortChanIDEncoding(b[0])
	case *uint8:
		var b [1]uint8
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		*e = b[0]
	case *FundingFlag:
		var b [1]uint8
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		*e = FundingFlag(b[0])
	case *uint16:
		var b [2]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint16(b[:])
	case *ChanUpdateMsgFlags:
		var b [1]uint8
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		*e = ChanUpdateMsgFlags(b[0])
	case *ChanUpdateChanFlags:
		var b [1]uint8
		if _, err := r.Read(b[:]); err != nil {
			return err
		}
		*e = ChanUpdateChanFlags(b[0])
	case *ErrorCode:
		var b [2]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = ErrorCode(binary.BigEndian.Uint16(b[:]))
	case *uint32:
		var b [4]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint32(b[:])
	case *uint64:
		var b [8]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint64(b[:])
	case *MilliSatoshi:
		var b [8]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = MilliSatoshi(int64(binary.BigEndian.Uint64(b[:])))
	case *btcutil.Amount:
		var b [8]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = btcutil.Amount(int64(binary.BigEndian.Uint64(b[:])))
	case **btcec.PublicKey:
		var b [btcec.PubKeyBytesLenCompressed]byte
		if _, err = io.ReadFull(r, b[:]); err != nil {
			return err
		}

		pubKey, err := btcec.ParsePubKey(b[:], btcec.S256())
		if err != nil {
			return err
		}
		*e = pubKey
	case **RawFeatureVector:
		f := NewRawFeatureVector()
		err = f.Decode(r)
		if err != nil {
			return err
		}

		*e = f

	case *[]Sig:
		var l [2]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return err
		}
		numSigs := binary.BigEndian.Uint16(l[:])

		var sigs []Sig
		if numSigs > 0 {
			sigs = make([]Sig, numSigs)
			for i := 0; i < int(numSigs); i++ {
				if err := ReadElement(r, &sigs[i]); err != nil {
					return err
				}
			}
		}

		*e = sigs

	case *Sig:
		if _, err := io.ReadFull(r, e[:]); err != nil {
			return err
		}
	case *OpaqueReason:
		var l [2]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return err
		}
		reasonLen := binary.BigEndian.Uint16(l[:])

		*e = OpaqueReason(make([]byte, reasonLen))
		if _, err := io.ReadFull(r, *e); err != nil {
			return err
		}
	case *ErrorData:
		var l [2]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return err
		}
		errorLen := binary.BigEndian.Uint16(l[:])

		*e = ErrorData(make([]byte, errorLen))
		if _, err := io.ReadFull(r, *e); err != nil {
			return err
		}
	case *PingPayload:
		var l [2]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return err
		}
		pingLen := binary.BigEndian.Uint16(l[:])

		*e = PingPayload(make([]byte, pingLen))
		if _, err := io.ReadFull(r, *e); err != nil {
			return err
		}
	case *PongPayload:
		var l [2]byte
		if _, err := io.ReadFull(r, l[:]); err != nil {
			return err
		}
		pongLen := binary.BigEndian.Uint16(l[:])

		*e = PongPayload(make([]byte, pongLen))
		if _, err := io.ReadFull(r, *e); err != nil {
			return err
		}
	case *[33]byte:
		if _, err := io.ReadFull(r, e[:]); err != nil {
			return err
		}
	case []byte:
		if _, err := io.ReadFull(r, e); err != nil {
			return err
		}
	case *PkScript:
		pkScript, err := wire.ReadVarBytes(r, 0, 34, "pkscript")
		if err != nil {
			return err
		}
		*e = pkScript
	case *wire.OutPoint:
		var h [32]byte
		if _, err = io.ReadFull(r, h[:]); err != nil {
			return err
		}
		hash, err := chainhash.NewHash(h[:])
		if err != nil {
			return err
		}

		var idxBytes [2]byte
		_, err = io.ReadFull(r, idxBytes[:])
		if err != nil {
			return err
		}
		index := binary.BigEndian.Uint16(idxBytes[:])

		*e = wire.OutPoint{
			Hash:  *hash,
			Index: uint32(index),
		}
	case *FailCode:
		if err := ReadElement(r, (*uint16)(e)); err != nil {
			return err
		}
	case *ChannelID:
		if _, err := io.ReadFull(r, e[:]); err != nil {
			return err
		}

	case *ShortChannelID:
		var blockHeight [4]byte
		if _, err = io.ReadFull(r, blockHeight[1:]); err != nil {
			return err
		}

		var txIndex [4]byte
		if _, err = io.ReadFull(r, txIndex[1:]); err != nil {
			return err
		}

		var txPosition [2]byte
		if _, err = io.ReadFull(r, txPosition[:]); err != nil {
			return err
		}

		*e = ShortChannelID{
			BlockHeight: binary.BigEndian.Uint32(blockHeight[:]),
			TxIndex:     binary.BigEndian.Uint32(txIndex[:]),
			TxPosition:  binary.BigEndian.Uint16(txPosition[:]),
		}

	case *[]net.Addr:
		// First, we'll read the number of total bytes that have been
		// used to encode the set of addresses.
		var numAddrsBytes [2]byte
		if _, err = io.ReadFull(r, numAddrsBytes[:]); err != nil {
			return err
		}
		addrsLen := binary.BigEndian.Uint16(numAddrsBytes[:])

		// With the number of addresses, read, we'll now pull in the
		// buffer of the encoded addresses into memory.
		addrs := make([]byte, addrsLen)
		if _, err := io.ReadFull(r, addrs[:]); err != nil {
			return err
		}
		addrBuf := bytes.NewReader(addrs)

		// Finally, we'll parse the remaining address payload in
		// series, using the first byte to denote how to decode the
		// address itself.
		var (
			addresses     []net.Addr
			addrBytesRead uint16
		)

		for addrBytesRead < addrsLen {
			var descriptor [1]byte
			if _, err = io.ReadFull(addrBuf, descriptor[:]); err != nil {
				return err
			}

			addrBytesRead++

			var address net.Addr
			switch aType := addressType(descriptor[0]); aType {
			case noAddr:
				addrBytesRead += aType.AddrLen()
				continue

			case tcp4Addr:
				var ip [4]byte
				if _, err := io.ReadFull(addrBuf, ip[:]); err != nil {
					return err
				}

				var port [2]byte
				if _, err := io.ReadFull(addrBuf, port[:]); err != nil {
					return err
				}

				address = &net.TCPAddr{
					IP:   net.IP(ip[:]),
					Port: int(binary.BigEndian.Uint16(port[:])),
				}
				addrBytesRead += aType.AddrLen()

			case tcp6Addr:
				var ip [16]byte
				if _, err := io.ReadFull(addrBuf, ip[:]); err != nil {
					return err
				}

				var port [2]byte
				if _, err := io.ReadFull(addrBuf, port[:]); err != nil {
					return err
				}

				address = &net.TCPAddr{
					IP:   net.IP(ip[:]),
					Port: int(binary.BigEndian.Uint16(port[:])),
				}
				addrBytesRead += aType.AddrLen()

			case v2OnionAddr:
				var h [tor.V2DecodedLen]byte
				if _, err := io.ReadFull(addrBuf, h[:]); err != nil {
					return err
				}

				var p [2]byte
				if _, err := io.ReadFull(addrBuf, p[:]); err != nil {
					return err
				}

				onionService := tor.Base32Encoding.EncodeToString(h[:])
				onionService += tor.OnionSuffix
				port := int(binary.BigEndian.Uint16(p[:]))

				address = &tor.OnionAddr{
					OnionService: onionService,
					Port:         port,
				}
				addrBytesRead += aType.AddrLen()

			case v3OnionAddr:
				var h [tor.V3DecodedLen]byte
				if _, err := io.ReadFull(addrBuf, h[:]); err != nil {
					return err
				}

				var p [2]byte
				if _, err := io.ReadFull(addrBuf, p[:]); err != nil {
					return err
				}

				onionService := tor.Base32Encoding.EncodeToString(h[:])
				onionService += tor.OnionSuffix
				port := int(binary.BigEndian.Uint16(p[:]))

				address = &tor.OnionAddr{
					OnionService: onionService,
					Port:         port,
				}
				addrBytesRead += aType.AddrLen()

			default:
				return &ErrUnknownAddrType{aType}
			}

			addresses = append(addresses, address)
		}

		*e = addresses
	case *color.RGBA:
		err := ReadElements(r,
			&e.R,
			&e.G,
			&e.B,
		)
		if err != nil {
			return err
		}
	case *DeliveryAddress:
		var addrLen [2]byte
		if _, err = io.ReadFull(r, addrLen[:]); err != nil {
			return err
		}
		length := binary.BigEndian.Uint16(addrLen[:])

		var addrBytes [34]byte
		if length > 34 {
			return fmt.Errorf("Cannot read %d bytes into addrBytes", length)
		}
		if _, err = io.ReadFull(r, addrBytes[:length]); err != nil {
			return err
		}
		*e = addrBytes[:length]
	default:
		return fmt.Errorf("Unknown type in ReadElement: %T", e)
	}

	return nil
}
