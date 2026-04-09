func segmentWords(data []byte, maxTokens int, atEOF bool, val [][]byte, types []int) ([][]byte, []int, int, error) {
  cs, p, pe := 0, 0, len(data)
  cap := maxTokens
  if cap < 0 {
    cap = 1000
  }
  if val == nil {
    val = make([][]byte, 0, cap)
  }
  if types == nil {
    types = make([]int, 0, cap)
  }

  // added for scanner
  ts := 0
  te := 0
  act := 0
  eof := pe
  _ = ts // compiler not happy
  _ = te
  _ = act

  // our state
  startPos := 0
  endPos := 0
  totalConsumed := 0
  
//line segment_words.go:18574
	{
	cs = s_start
	ts = 0
	te = 0
	act = 0
	}

//line segment_words.go:18582
	{
	var _klen int
	var _keys int
	var _trans int

	if p == pe {
		goto _test_eof
	}
_resume:
	switch _s_from_state_actions[cs] {
	case 23:
//line NONE:1
ts = p


//line segment_words.go:18598
	}

	_keys = int(_s_key_offsets[cs])
	_trans = int(_s_index_offsets[cs])

	_klen = int(_s_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _s_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _s_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_s_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _s_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _s_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_s_indicies[_trans])
_eof_trans:
	cs = int(_s_trans_targs[_trans])

	if _s_trans_actions[_trans] == 0 {
		goto _again
	}

	switch _s_trans_actions[_trans] {
	case 10:
//line segment_words.rl:72

    endPos = p
  

	case 34:
//line segment_words.rl:76
te = p
p--
{
    if !atEOF {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Number)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 37:
//line segment_words.rl:89
te = p
p--
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 39:
//line segment_words.rl:104
te = p
p--
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 35:
//line segment_words.rl:119
te = p
p--
{
    if !atEOF {
      return val, types, totalConsumed, nil
    }
    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 38:
//line segment_words.rl:131
te = p
p--
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 40:
//line segment_words.rl:146
te = p
p--
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 41:
//line segment_words.rl:161
te = p
p--
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 32:
//line segment_words.rl:161
te = p
p--
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 36:
//line segment_words.rl:161
te = p
p--
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 31:
//line segment_words.rl:161
te = p
p--
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 4:
//line segment_words.rl:76
p = (te) - 1
{
    if !atEOF {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Number)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 12:
//line segment_words.rl:89
p = (te) - 1
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 17:
//line segment_words.rl:104
p = (te) - 1
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 7:
//line segment_words.rl:119
p = (te) - 1
{
    if !atEOF {
      return val, types, totalConsumed, nil
    }
    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 15:
//line segment_words.rl:131
p = (te) - 1
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 19:
//line segment_words.rl:146
p = (te) - 1
{
    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 21:
//line segment_words.rl:161
p = (te) - 1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 8:
//line segment_words.rl:161
p = (te) - 1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 1:
//line segment_words.rl:161
p = (te) - 1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 3:
//line NONE:1
	switch act {
	case 1:
	{p = (te) - 1

    if !atEOF {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Number)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 2:
	{p = (te) - 1

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 3:
	{p = (te) - 1

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 4:
	{p = (te) - 1

    if !atEOF {
      return val, types, totalConsumed, nil
    }
    val = append(val, data[startPos:endPos+1])
    types = append(types, Letter)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 5:
	{p = (te) - 1

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }

    val = append(val, data[startPos:endPos+1])
    types = append(types, Ideo)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 7:
	{p = (te) - 1

    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 12:
	{p = (te) - 1

    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	case 13:
	{p = (te) - 1

    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }
	}
	

	case 28:
//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  

	case 33:
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
te = p+1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 13:
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
te = p+1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 18:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  

	case 26:
//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
te = p+1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 27:
//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
te = p+1
{
    lastPos := startPos
    for lastPos <= endPos {
      _, size := utf8.DecodeRune(data[lastPos:])
      lastPos += size
    }
    endPos = lastPos -1
    p = endPos

    if endPos+1 == pe && !atEOF {
      return val, types, totalConsumed, nil
    } else if dr, size := utf8.DecodeRune(data[endPos+1:]); dr == utf8.RuneError && size == 1 {
      return val, types, totalConsumed, nil
    }
    // otherwise, consume this as well
    val = append(val, data[startPos:endPos+1])
    types = append(types, None)
    totalConsumed = endPos+1
    if maxTokens > 0 && len(val) >= maxTokens {
      return val, types, totalConsumed, nil
    }
  }

	case 5:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:76
act = 1;

	case 11:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:89
act = 2;

	case 16:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:104
act = 3;

	case 6:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:119
act = 4;

	case 14:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:131
act = 5;

	case 20:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
act = 7;

	case 9:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
act = 12;

	case 2:
//line NONE:1
te = p+1

//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
act = 13;

	case 29:
//line NONE:1
te = p+1

//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:76
act = 1;

	case 30:
//line NONE:1
te = p+1

//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:119
act = 4;

	case 25:
//line NONE:1
te = p+1

//line segment_words.rl:68

    startPos = p
  
//line segment_words.rl:72

    endPos = p
  
//line segment_words.rl:161
act = 13;

//line segment_words.go:19500
	}

_again:
	switch _s_to_state_actions[cs] {
	case 22:
//line NONE:1
ts = 0


//line segment_words.go:19510
	}

	if p++; p != pe {
		goto _resume
	}
	_test_eof: {}
	if p == eof {
		if _s_eof_trans[cs] > 0 {
			_trans = int(_s_eof_trans[cs] - 1)
			goto _eof_trans
		}
		switch _s_eof_actions[cs] {
		case 24:
//line segment_words.rl:68

    startPos = p
  

//line segment_words.go:19529
		}
	}

	}

//line segment_words.rl:278


  if cs < s_first_final {
    return val, types, totalConsumed, ParseError
  }

  return val, types, totalConsumed, nil
}
