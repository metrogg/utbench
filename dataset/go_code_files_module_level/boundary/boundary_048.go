func decodeLineOfMediaPlaylist(p *MediaPlaylist, wv *WV, state *decodingState, line string, strict bool) error {
	var err error

	line = strings.TrimSpace(line)
	switch {
	case !state.tagInf && strings.HasPrefix(line, "#EXTINF:"):
		state.tagInf = true
		state.listType = MEDIA
		sepIndex := strings.Index(line, ",")
		if sepIndex == -1 {
			if strict {
				return fmt.Errorf("could not parse: %q", line)
			}
			sepIndex = len(line)
		}
		duration := line[8:sepIndex]
		if len(duration) > 0 {
			if state.duration, err = strconv.ParseFloat(duration, 64); strict && err != nil {
				return fmt.Errorf("Duration parsing error: %s", err)
			}
		}
		if len(line) > sepIndex {
			state.title = line[sepIndex+1:]
		}
	case !strings.HasPrefix(line, "#"):
		if state.tagInf {
			err := p.Append(line, state.duration, state.title)
			if err == ErrPlaylistFull {
				// Extend playlist by doubling size, reset internal state, try again.
				// If the second Append fails, the if err block will handle it.
				// Retrying instead of being recursive was chosen as the state maybe
				// modified non-idempotently.
				p.Segments = append(p.Segments, make([]*MediaSegment, p.Count())...)
				p.capacity = uint(len(p.Segments))
				p.tail = p.count
				err = p.Append(line, state.duration, state.title)
			}
			// Check err for first or subsequent Append()
			if err != nil {
				return err
			}
			state.tagInf = false
		}
		if state.tagRange {
			if err = p.SetRange(state.limit, state.offset); strict && err != nil {
				return err
			}
			state.tagRange = false
		}
		if state.tagSCTE35 {
			state.tagSCTE35 = false
			if err = p.SetSCTE35(state.scte); strict && err != nil {
				return err
			}
		}
		if state.tagDiscontinuity {
			state.tagDiscontinuity = false
			if err = p.SetDiscontinuity(); strict && err != nil {
				return err
			}
		}
		if state.tagProgramDateTime && p.Count() > 0 {
			state.tagProgramDateTime = false
			if err = p.SetProgramDateTime(state.programDateTime); strict && err != nil {
				return err
			}
		}
		// If EXT-X-KEY appeared before reference to segment (EXTINF) then it linked to this segment
		if state.tagKey {
			p.Segments[p.last()].Key = &Key{state.xkey.Method, state.xkey.URI, state.xkey.IV, state.xkey.Keyformat, state.xkey.Keyformatversions}
			// First EXT-X-KEY may appeared in the header of the playlist and linked to first segment
			// but for convenient playlist generation it also linked as default playlist key
			if p.Key == nil {
				p.Key = state.xkey
			}
			state.tagKey = false
		}
		// If EXT-X-MAP appeared before reference to segment (EXTINF) then it linked to this segment
		if state.tagMap {
			p.Segments[p.last()].Map = &Map{state.xmap.URI, state.xmap.Limit, state.xmap.Offset}
			// First EXT-X-MAP may appeared in the header of the playlist and linked to first segment
			// but for convenient playlist generation it also linked as default playlist map
			if p.Map == nil {
				p.Map = state.xmap
			}
			state.tagMap = false
		}
	// start tag first
	case line == "#EXTM3U":
		state.m3u = true
	case line == "#EXT-X-ENDLIST":
		state.listType = MEDIA
		p.Closed = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &p.TargetDuration); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &p.SeqNo); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
		state.listType = MEDIA
		var playlistType string
		_, err = fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &playlistType)
		if err != nil {
			if strict {
				return err
			}
		} else {
			switch playlistType {
			case "EVENT":
				p.MediaType = EVENT
			case "VOD":
				p.MediaType = VOD
			}
		}
	case strings.HasPrefix(line, "#EXT-X-DISCONTINUITY-SEQUENCE:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-DISCONTINUITY-SEQUENCE:%d", &p.DiscontinuitySeq); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-START:"):
		state.listType = MEDIA
		for k, v := range decodeParamsLine(line[13:]) {
			switch k {
			case "TIME-OFFSET":
				st, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("Invalid TIME-OFFSET: %s: %v", v, err)
				}
				p.StartTime = st
			case "PRECISE":
				p.StartTimePrecise = v == "YES"
			}
		}
	case strings.HasPrefix(line, "#EXT-X-KEY:"):
		state.listType = MEDIA
		state.xkey = new(Key)
		for k, v := range decodeParamsLine(line[11:]) {
			switch k {
			case "METHOD":
				state.xkey.Method = v
			case "URI":
				state.xkey.URI = v
			case "IV":
				state.xkey.IV = v
			case "KEYFORMAT":
				state.xkey.Keyformat = v
			case "KEYFORMATVERSIONS":
				state.xkey.Keyformatversions = v
			}
		}
		state.tagKey = true
	case strings.HasPrefix(line, "#EXT-X-MAP:"):
		state.listType = MEDIA
		state.xmap = new(Map)
		for k, v := range decodeParamsLine(line[11:]) {
			switch k {
			case "URI":
				state.xmap.URI = v
			case "BYTERANGE":
				if _, err = fmt.Sscanf(v, "%d@%d", &state.xmap.Limit, &state.xmap.Offset); strict && err != nil {
					return fmt.Errorf("Byterange sub-range length value parsing error: %s", err)
				}
			}
		}
		state.tagMap = true
	case !state.tagProgramDateTime && strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME:"):
		state.tagProgramDateTime = true
		state.listType = MEDIA
		if state.programDateTime, err = TimeParse(line[25:]); strict && err != nil {
			return err
		}
	case !state.tagRange && strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
		state.tagRange = true
		state.listType = MEDIA
		state.offset = 0
		params := strings.SplitN(line[17:], "@", 2)
		if state.limit, err = strconv.ParseInt(params[0], 10, 64); strict && err != nil {
			return fmt.Errorf("Byterange sub-range length value parsing error: %s", err)
		}
		if len(params) > 1 {
			if state.offset, err = strconv.ParseInt(params[1], 10, 64); strict && err != nil {
				return fmt.Errorf("Byterange sub-range offset value parsing error: %s", err)
			}
		}
	case !state.tagSCTE35 && strings.HasPrefix(line, "#EXT-SCTE35:"):
		state.tagSCTE35 = true
		state.listType = MEDIA
		state.scte = new(SCTE)
		state.scte.Syntax = SCTE35_67_2014
		for attribute, value := range decodeParamsLine(line[12:]) {
			switch attribute {
			case "CUE":
				state.scte.Cue = value
			case "ID":
				state.scte.ID = value
			case "TIME":
				state.scte.Time, _ = strconv.ParseFloat(value, 64)
			}
		}
	case !state.tagSCTE35 && strings.HasPrefix(line, "#EXT-OATCLS-SCTE35:"):
		// EXT-OATCLS-SCTE35 contains the SCTE35 tag, EXT-X-CUE-OUT contains duration
		state.tagSCTE35 = true
		state.scte = new(SCTE)
		state.scte.Syntax = SCTE35_OATCLS
		state.scte.Cue = line[19:]
	case state.tagSCTE35 && state.scte.Syntax == SCTE35_OATCLS && strings.HasPrefix(line, "#EXT-X-CUE-OUT:"):
		// EXT-OATCLS-SCTE35 contains the SCTE35 tag, EXT-X-CUE-OUT contains duration
		state.scte.Time, _ = strconv.ParseFloat(line[15:], 64)
		state.scte.CueType = SCTE35Cue_Start
	case !state.tagSCTE35 && strings.HasPrefix(line, "#EXT-X-CUE-OUT-CONT:"):
		state.tagSCTE35 = true
		state.scte = new(SCTE)
		state.scte.Syntax = SCTE35_OATCLS
		state.scte.CueType = SCTE35Cue_Mid
		for attribute, value := range decodeParamsLine(line[20:]) {
			switch attribute {
			case "SCTE35":
				state.scte.Cue = value
			case "Duration":
				state.scte.Time, _ = strconv.ParseFloat(value, 64)
			case "ElapsedTime":
				state.scte.Elapsed, _ = strconv.ParseFloat(value, 64)
			}
		}
	case !state.tagSCTE35 && line == "#EXT-X-CUE-IN":
		state.tagSCTE35 = true
		state.scte = new(SCTE)
		state.scte.Syntax = SCTE35_OATCLS
		state.scte.CueType = SCTE35Cue_End
	case !state.tagDiscontinuity && strings.HasPrefix(line, "#EXT-X-DISCONTINUITY"):
		state.tagDiscontinuity = true
		state.listType = MEDIA
	case strings.HasPrefix(line, "#EXT-X-I-FRAMES-ONLY"):
		state.listType = MEDIA
		p.Iframe = true
	case strings.HasPrefix(line, "#WV-AUDIO-CHANNELS"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-CHANNELS %d", &wv.AudioChannels); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-FORMAT"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-FORMAT %d", &wv.AudioFormat); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-PROFILE-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-PROFILE-IDC %d", &wv.AudioProfileIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLE-SIZE"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLE-SIZE %d", &wv.AudioSampleSize); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLING-FREQUENCY"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLING-FREQUENCY %d", &wv.AudioSamplingFrequency); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-CYPHER-VERSION"):
		state.listType = MEDIA
		wv.CypherVersion = line[19:]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-ECM"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-ECM %s", &wv.ECM); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-FORMAT"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-FORMAT %d", &wv.VideoFormat); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-FRAME-RATE"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-FRAME-RATE %d", &wv.VideoFrameRate); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-LEVEL-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-LEVEL-IDC %d", &wv.VideoLevelIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-PROFILE-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-PROFILE-IDC %d", &wv.VideoProfileIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-RESOLUTION"):
		state.listType = MEDIA
		wv.VideoResolution = line[21:]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-VIDEO-SAR"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-SAR %s", &wv.VideoSAR); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#"): // unknown tags treated as comments
		return err
	}
	return err
}
