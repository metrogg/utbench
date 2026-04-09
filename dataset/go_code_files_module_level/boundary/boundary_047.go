func (b *blockDec) decodeCompressed(hist *history) error {
	in := b.data
	delayedHistory := hist == nil

	if delayedHistory {
		// We must always grab history.
		defer func() {
			if hist == nil {
				<-b.history
			}
		}()
	}
	// There must be at least one byte for Literals_Block_Type and one for Sequences_Section_Header
	if len(in) < 2 {
		return ErrBlockTooSmall
	}
	litType := literalsBlockType(in[0] & 3)
	var litRegenSize int
	var litCompSize int
	sizeFormat := (in[0] >> 2) & 3
	var fourStreams bool
	switch litType {
	case literalsBlockRaw, literalsBlockRLE:
		switch sizeFormat {
		case 0, 2:
			// Regenerated_Size uses 5 bits (0-31). Literals_Section_Header uses 1 byte.
			litRegenSize = int(in[0] >> 3)
			in = in[1:]
		case 1:
			// Regenerated_Size uses 12 bits (0-4095). Literals_Section_Header uses 2 bytes.
			litRegenSize = int(in[0]>>4) + (int(in[1]) << 4)
			in = in[2:]
		case 3:
			//  Regenerated_Size uses 20 bits (0-1048575). Literals_Section_Header uses 3 bytes.
			if len(in) < 3 {
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			}
			litRegenSize = int(in[0]>>4) + (int(in[1]) << 4) + (int(in[2]) << 12)
			in = in[3:]
		}
	case literalsBlockCompressed, literalsBlockTreeless:
		switch sizeFormat {
		case 0, 1:
			// Both Regenerated_Size and Compressed_Size use 10 bits (0-1023).
			if len(in) < 3 {
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			}
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12)
			litRegenSize = int(n & 1023)
			litCompSize = int(n >> 10)
			fourStreams = sizeFormat == 1
			in = in[3:]
		case 2:
			fourStreams = true
			if len(in) < 4 {
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			}
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20)
			litRegenSize = int(n & 16383)
			litCompSize = int(n >> 14)
			in = in[4:]
		case 3:
			fourStreams = true
			if len(in) < 5 {
				println("too small: litType:", litType, " sizeFormat", sizeFormat, len(in))
				return ErrBlockTooSmall
			}
			n := uint64(in[0]>>4) + (uint64(in[1]) << 4) + (uint64(in[2]) << 12) + (uint64(in[3]) << 20) + (uint64(in[4]) << 28)
			litRegenSize = int(n & 262143)
			litCompSize = int(n >> 18)
			in = in[5:]
		}
	}
	var literals []byte
	var huff *huff0.Scratch
	switch litType {
	case literalsBlockRaw:
		if len(in) < litRegenSize {
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litRegenSize)
			return ErrBlockTooSmall
		}
		literals = in[:litRegenSize]
		in = in[litRegenSize:]
		//printf("Found %d uncompressed literals\n", litRegenSize)
	case literalsBlockRLE:
		if len(in) < 1 {
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", 1)
			return ErrBlockTooSmall
		}
		if cap(b.literalBuf) < litRegenSize {
			if b.lowMem {
				b.literalBuf = make([]byte, litRegenSize)
			} else {
				if litRegenSize > maxCompressedLiteralSize {
					// Exceptional
					b.literalBuf = make([]byte, litRegenSize)
				} else {
					b.literalBuf = make([]byte, litRegenSize, maxCompressedLiteralSize)

				}
			}
		}
		literals = b.literalBuf[:litRegenSize]
		v := in[0]
		for i := range literals {
			literals[i] = v
		}
		in = in[1:]
		if debug {
			printf("Found %d RLE compressed literals\n", litRegenSize)
		}
	case literalsBlockTreeless:
		if len(in) < litCompSize {
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return ErrBlockTooSmall
		}
		// Store compressed literals, so we defer decoding until we get history.
		literals = in[:litCompSize]
		in = in[litCompSize:]
		if debug {
			printf("Found %d compressed literals\n", litCompSize)
		}
	case literalsBlockCompressed:
		if len(in) < litCompSize {
			println("too small: litType:", litType, " sizeFormat", sizeFormat, "remain:", len(in), "want:", litCompSize)
			return ErrBlockTooSmall
		}
		literals = in[:litCompSize]
		in = in[litCompSize:]

		huff = huffDecoderPool.Get().(*huff0.Scratch)
		var err error
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize {
			if b.lowMem {
				b.literalBuf = make([]byte, 0, litRegenSize)
			} else {
				b.literalBuf = make([]byte, 0, maxCompressedLiteralSize)
			}
		}
		if huff == nil {
			huff = &huff0.Scratch{}
		}
		huff.Out = b.literalBuf[:0]
		huff, literals, err = huff0.ReadTable(literals, huff)
		if err != nil {
			println("reading huffman table:", err)
			return err
		}
		// Use our out buffer.
		huff.Out = b.literalBuf[:0]
		if fourStreams {
			literals, err = huff.Decompress4X(literals, litRegenSize)
		} else {
			literals, err = huff.Decompress1X(literals)
		}
		if err != nil {
			return err
		}
		// Make sure we don't leak our literals buffer
		huff.Out = nil
		if len(literals) != litRegenSize {
			return fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		}
		if debug {
			printf("Decompressed %d literals into %d bytes\n", litCompSize, litRegenSize)
		}
	}

	// Decode Sequences
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#sequences-section
	if len(in) < 1 {
		return ErrBlockTooSmall
	}
	seqHeader := in[0]
	nSeqs := 0
	switch {
	case seqHeader == 0:
		in = in[1:]
	case seqHeader < 128:
		nSeqs = int(seqHeader)
		in = in[1:]
	case seqHeader < 255:
		if len(in) < 2 {
			return ErrBlockTooSmall
		}
		nSeqs = int(seqHeader-128)<<8 | int(in[1])
		in = in[2:]
	case seqHeader == 255:
		if len(in) < 3 {
			return ErrBlockTooSmall
		}
		nSeqs = 0x7f00 + int(in[1]) + (int(in[2]) << 8)
		in = in[3:]
	}

	// Allocate sequences
	if cap(b.sequenceBuf) < nSeqs {
		if b.lowMem {
			b.sequenceBuf = make([]seq, nSeqs)
		} else {
			// Allocate max
			b.sequenceBuf = make([]seq, nSeqs, 0x7f00+0xffff)
		}
	} else {
		// Reuse buffer
		b.sequenceBuf = b.sequenceBuf[:nSeqs]
	}
	var seqs = &sequenceDecs{}
	if nSeqs > 0 {
		if len(in) < 1 {
			return ErrBlockTooSmall
		}
		br := byteReader{b: in, off: 0}
		compMode := br.Uint8()
		br.advance(1)
		for i := uint(0); i < 3; i++ {
			mode := seqCompMode((compMode >> (6 - i*2)) & 3)
			//println("Table", tableIndex(i), "is", mode)
			var seq *sequenceDec
			switch tableIndex(i) {
			case tableLiteralLengths:
				seq = &seqs.litLengths
			case tableOffsets:
				seq = &seqs.offsets
			case tableMatchLengths:
				seq = &seqs.matchLengths
			}
			switch mode {
			case compModePredefined:
				seq.fse = &fsePredef[i]
			case compModeRLE:
				if br.remain() < 1 {
					return ErrBlockTooSmall
				}
				v := br.Uint8()
				br.advance(1)
				dec := fseDecoderPool.Get().(*fseDecoder)
				symb, err := decSymbolValue(v, symbolTableX[i])
				if err != nil {
					println("RLE Transform table error:", err)
					return err
				}
				dec.setRLE(symb)
				seq.fse = dec
				if debug {
					println("RLE set to ", symb)
				}
			case compModeFSE:
				println("Reading table for", tableIndex(i))
				dec := fseDecoderPool.Get().(*fseDecoder)
				err := dec.readNCount(&br, uint16(maxTableSymbol[i]))
				if err != nil {
					println("Read table error:", err)
					return err
				}
				println("Read table ok")
				err = dec.transform(symbolTableX[i])
				if err != nil {
					println("Transform table error:", err)
					return err
				}
				seq.fse = dec
			case compModeRepeat:
				seq.repeat = true
			}
			if br.overread() {
				return io.ErrUnexpectedEOF
			}
		}
		in = br.unread()
	}

	// Wait for history.
	// All time spent after this is critical since it is strictly sequential.
	if hist == nil {
		hist = <-b.history
		if hist.error {
			return ErrDecoderClosed
		}
	}

	// Decode treeless literal block.
	if litType == literalsBlockTreeless {
		// TODO: We could send the history early WITHOUT the stream history.
		//   This would allow decoding treeless literials before the byte history is available.
		//   Silencia stats: Treeless 4393, with: 32775, total: 37168, 11% treeless.
		//   So not much obvious gain here.

		if hist.huffTree == nil {
			return errors.New("literal block was treeless, but no history was defined")
		}
		// Ensure we have space to store it.
		if cap(b.literalBuf) < litRegenSize {
			if b.lowMem {
				b.literalBuf = make([]byte, 0, litRegenSize)
			} else {
				b.literalBuf = make([]byte, 0, maxCompressedLiteralSize)
			}
		}
		var err error
		// Use our out buffer.
		huff = hist.huffTree
		huff.Out = b.literalBuf[:0]
		if fourStreams {
			literals, err = huff.Decompress4X(literals, litRegenSize)
		} else {
			literals, err = huff.Decompress1X(literals)
		}
		// Make sure we don't leak our literals buffer
		huff.Out = nil
		if err != nil {
			println("decompressing literals:", err)
			return err
		}
		if len(literals) != litRegenSize {
			return fmt.Errorf("literal output size mismatch want %d, got %d", litRegenSize, len(literals))
		}
	} else {
		if hist.huffTree != nil && huff != nil {
			huffDecoderPool.Put(hist.huffTree)
			hist.huffTree = nil
		}
	}
	if huff != nil {
		huff.Out = nil
		hist.huffTree = huff
	}
	if debug {
		println("Final literals:", len(literals), "and", nSeqs, "sequences.")
	}

	if nSeqs == 0 {
		// Decompressed content is defined entirely as Literals Section content.
		b.dst = append(b.dst, literals...)
		if delayedHistory {
			hist.append(literals)
		}
		return nil
	}

	seqs, err := seqs.mergeHistory(&hist.decoders)
	if err != nil {
		return err
	}
	if debug {
		println("History merged ok")
	}
	br := &bitReader{}
	if err := br.init(in); err != nil {
		return err
	}

	// TODO: Investigate if sending history without decoders are faster.
	//   This would allow the sequences to be decoded async and only have to construct stream history.
	//   If only recent offsets were not transferred, this would be an obvious win.

	if err := seqs.initialize(br, hist, literals, b.dst); err != nil {
		println("initializing sequences:", err)
		return err
	}

	err = seqs.decode(nSeqs, br, hist.b)
	if err != nil {
		return err
	}
	if !br.finished() {
		// Disabled until this is resolved:
		// https://github.com/facebook/zstd/issues/1597
		// return fmt.Errorf("%d extra bits on block, should be 0", br.remain())
	}

	err = br.close()
	if err != nil {
		printf("Closing sequences: %v, %+v\n", err, *br)
	}

	// Set output and release references.
	b.dst = seqs.out
	seqs.out, seqs.literals, seqs.hist = nil, nil, nil

	if !delayedHistory {
		// If we don't have delayed history, no need to update.
		hist.recentOffsets = seqs.prevOffset
		return nil
	}
	if b.Last {
		// if last block we don't care about history.
		println("Last block, no history returned")
		hist.b = hist.b[:0]
		return nil
	}
	hist.append(b.dst)
	hist.recentOffsets = seqs.prevOffset
	return nil
}
