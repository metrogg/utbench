func (p *Decoder) parseMetaMsg(e *Event) (nextChunkType, bool, error) {
	var err error
	// channels aren't really channels
	switch e.MsgChan {
	case 0xF:
		var b byte
		if b, err = p.ReadByte(); err != nil {
			return eventChunk, false, err
		}
		e.Cmd = uint8(b)

		switch e.Cmd {

		// Sequence Number
		//This optional event, which must occur at the beginning of a
		//track, before any nonzero delta-times, and before any
		//transmittable MIDI events, specifies the number of a sequence. In a
		//format 2 MIDI File, it is used to identify each "pattern" so that a
		//"song" sequence using the Cue message to refer to the patterns. If
		//the ID numbers are omitted, the sequences' locations in order in the
		//file are used as defaults. In a format 0 or 1 MIDI File, which only
		//contain one sequence, this number should be contained in the first
		//(or only) track. If transfer of several multitrack sequences is
		//required, this must be done as a group of format 1 files, each with
		//a different sequence number.
		case 0x0:
			if b, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}
			if uint8(b) != 0 {
				var seqB []byte
				if err = p.Read(seqB); err != nil {
					return eventChunk, false, err
				}
				e.SeqNum = binary.BigEndian.Uint16(seqB)
			}

		// Text Event
		//Any amount of text describing anything. It is a good idea to put
		//a text event right at the beginning of a track, with the name of the
		//track, a description of its intended orchestration, and any other
		//information which the user wants to put there. Text events may also
		//occur at other times in a track, to be used as lyrics, or descriptions
		//of cue points. The text in this event should be printable ASCII
		//characters for maximum interchange. However, other characters codes
		//using the high-order bit may be used for interchange of files between
		//different programs on the same computer which supports an extended
		//character set. Programs on a computer which does not support
		//non-ASCII characters should ignore those characters.
		case 0x01:
			if e.Text, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Copyright Notice
		case 0x02:
			if e.Copyright, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Sequence/Track Name
		//If in a format 0 track, or the first track in a format 1 file, the
		//name of the sequence. Otherwise, the name of the track.
		case 0x03:
			if e.SeqTrackName, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Instrument name
		//A description of the type of instrumentation to be used in that track.
		//May be used with the MIDI Prefix meta-event to specify which MIDI
		//channel the description applies to, or the channel may be specified
		//as text in the event itself.
		//
		case 0x04:
			if e.InstrumentName, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Lyric
		case 0x05:
			if e.Lyric, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Marker
		//Normally in a format 0 track, or the first track in a format 1
		//file. The name of that point in the sequence, such as a rehersal letter
		case 0x06:
			if e.Marker, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// Cue point
		case 0x07:
			if e.CuePoint, _, err = p.VarLenTxt(); err != nil {
				return eventChunk, false, err
			}

		// MIDI Channel Prefix
		//The MIDI channel (0-15) containted in this event may be used
		//to associate a MIDI channel with all events which follow, including
		//System exclusive and meta-events. This channel is "effective" until
		//the next normal MIDI event (which contains a channel) or the next MIDI
		//Channel Prefix meta-event. If MIDI channels refer to "tracks", this
		//message may into a format 0 file, keeping their non-MIDI data
		//associated with a track. This capability is also present in Yamaha's
		//ESEQ file format.
		case 0x20:
			var b byte
			if b, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}
			if uint8(b) != 1 {
				return eventChunk, false, fmt.Errorf("MIDI Channel Prefix event error - %s", ErrUnexpectedData)
			}

			if e.Channel, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}

		// End of track
		//This event is not optional. It is included so that an exact
		//ending point may be specified for the track, so that an exect length,
		//which is necessary for tracks which are looped or concatenated.
		case 0x2F:
			_, err = p.ReadByte()
			return trackChunk, true, err

			/*
			   Set Tempo(in microseconds per MIDI quarter-note)
			      This meta event sets the sequence tempo in terms of microseconds per quarter-note which is encoded in three bytes.
			      It usually is found in the first track chunk, time-aligned to occur at the same time as a MIDI clock message
			      to promote more accurate synchronization. If no set tempo event is present, 120 beats per minute is assumed.
			      The following formula's can be used to translate the tempo from microseconds per quarter-note to beats per minute and back.
			*/
		case 0x51:
			var l uint32
			if l, _, err = p.VarLen(); err != nil {
				return eventChunk, false, err
			}
			if l != 3 {
				return eventChunk, false, fmt.Errorf("Set Tempo event error - %s", ErrUnexpectedData)
			}

			if e.MsPerQuartNote, err = p.Uint24(); err != nil {
				return eventChunk, false, err
			}

			e.Bpm = 60000000 / e.MsPerQuartNote

			/*
			     SMPTE Offset
			   This meta event is used to specify the SMPTE starting point offset from the beginning of the track.
			   It is defined in terms of hours, minutes, seconds, frames and sub-frames (always 100 sub-frames per frame,
			   no matter what sub-division is specified in the MIDI header chunk).
			   The byte used to specify the hour offset also specifies the frame rate in the following format:
			   0rrhhhhhh where rr is two bits for the frame rate where 00=24 fps, 01=25 fps, 10=30 fps (drop frame),
			   11=30 fps and hhhhhh is six bits for the hour (0-23). The hour byte's top bit is always 0.
			   The frame byte's possible range depends on the encoded frame rate in the hour byte.
			   A 25 fps frame rate means that a maximum value of 24 may be set for the frame byte.
			*/
		case 0x54:
			var l uint32
			if l, _, err = p.VarLen(); err != nil {
				return eventChunk, false, err
			}
			if l != 5 {
				return eventChunk, false, fmt.Errorf("error parsing SMPTE Offset - %s (%d)", ErrUnexpectedData, l)
			}
			so := &SmpteOffset{}
			if err = p.Read(&so.Hour); err != nil {
				return eventChunk, false, err
			}
			if err = p.Read(&so.Min); err != nil {
				return eventChunk, false, err
			}
			if err = p.Read(&so.Sec); err != nil {
				return eventChunk, false, err
			}
			if err = p.Read(&so.Fr); err != nil {
				return eventChunk, false, err
			}
			if err = p.Read(&so.SubFr); err != nil {
				return eventChunk, false, err
			}
			e.SmpteOffset = so

			// Time signature
			// FF 58 04 nn dd cc bb Time Signature
			// The time signature is expressed as four numbers. nn and dd
			// represent the numerator and denominator of the time signature as it
			// would be notated. The denominator is a negative power of two: 2
			// represents a quarter-note, 3 represents an eighth-note, etc.
			// The cc parameter expresses the number of MIDI clocks in a
			// metronome click. The bb parameter expresses the number of
			// notated 32nd-notes in a MIDI quarter-note (24 MIDI clocks). This
			// was added because there are already multiple programs which allow a
			// user to specify that what MIDI thinks of as a quarter-note (24 clocks)
			// is to be notated as, or related to in terms of, something else.
			//
			// This meta event is used to set a sequences time signature.
			// The time signature defined with 4 bytes, a numerator, a denominator, a metronome
			// pulse and number of 32nd notes per MIDI quarter-note. The numerator is specified as
			// a literal value, but the denominator is specified as (get ready) the value to which
			// the power of 2 must be raised to equal the number of subdivisions per whole note.
			// For example, a value of 0 means a whole note because 2 to the power of 0 is 1
			// (whole note), a value of 1 means a half-note because 2 to the power of 1 is 2
			// (half-note), and so on. The metronome pulse specifies how often the metronome should
			// click in terms of the number of clock signals per click, which come at a rate of 24
			// per quarter-note. For example, a value of 24 would mean to click once every quarter-note
			// (beat) and a value of 48 would mean to click once every half-note (2 beats).
			// And finally, the fourth byte specifies the number of 32nd notes per 24 MIDI clock signals.
			// This value is usually 8 because there are usually 8 32nd notes in a quarter-note.
			// At least one Time Signature Event should appear in the first track chunk
			// (or all track chunks in a Type 2 file) before any non-zero delta time events.
			// If one is not specified 4/4, 24, 8 should be assumed.
		case 0x58:
			var l uint32
			if l, _, err = p.VarLen(); err != nil {
				return eventChunk, false, err
			}

			if l != 4 {
				return eventChunk, false, fmt.Errorf("error parsing TimeSignature - %s (%d)", ErrUnexpectedData, l)
			}

			var num, denom, clocksPerClick, thirtySecondNotesPerQuarter byte
			if num, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}
			//The denominator is a neqative power of two: 2
			//represents a quarter-note, 3 represents an eighth-note, etc
			if denom, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}
			if clocksPerClick, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}
			// The number of notated 32nd-notes in a MIDI quarter-note (24 MIDI clocks).
			if thirtySecondNotesPerQuarter, err = p.ReadByte(); err != nil {
				return eventChunk, false, err
			}

			e.TimeSignature = &TimeSignature{uint8(num), uint8(denom), uint8(clocksPerClick), uint8(thirtySecondNotesPerQuarter)}

		// Key signature
		// This meta event is used to specify the key (number of sharps or flats) and scale (major or minor) of a sequence.
		// A positive value for the key specifies the number of sharps and a negative value specifies the number of flats.
		// A value of 0 for the scale specifies a major key and a value of 1 specifies a minor key.
		case 0x59:
			var l uint32
			if l, _, err = p.VarLen(); err != nil {
				return eventChunk, false, err
			}
			if l != 2 {
				return eventChunk, false, fmt.Errorf("Time Signature length not 2 as expected but %d", l)
			}

			// key
			b, err := p.ReadByte()
			if err != nil {
				return eventChunk, false, err
			}
			e.Key = int32(b)

			// scale
			b, err = p.ReadByte()
			if err != nil {
				return eventChunk, false, err
			}
			e.Scale = uint32(b)

		// Sequencer Specific
		//This meta event is used to specify information specific to a hardware or software sequencer.
		// The first Data byte (or three bytes if the first byte is 0) specifies the manufacturer's ID
		// and the following bytes contain information specified by the manufacturer.
		// The individual manufacturers may document this information in their respective manuals.
		case 0x7F:
			var l uint32
			if l, _, err = p.VarLen(); err != nil {
				return eventChunk, false, err
			}
			// Information not currently stored
			tmp := make([]byte, l)
			err = p.Read(tmp)
		default:
			fmt.Printf("Skipped meta cmd %#X\n", e.Cmd)
		}
	}

	return eventChunk, true, nil
}
