func (j *TarOptions) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtTarOptionsbase
	_ = currentKey
	tok := fflib.FFTok_init
	wantedTok := fflib.FFTok_init

mainparse:
	for {
		tok = fs.Scan()
		//	println(fmt.Sprintf("debug: tok: %v  state: %v", tok, state))
		if tok == fflib.FFTok_error {
			goto tokerror
		}

		switch state {

		case fflib.FFParse_map_start:
			if tok != fflib.FFTok_left_bracket {
				wantedTok = fflib.FFTok_left_bracket
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_key
			continue

		case fflib.FFParse_after_value:
			if tok == fflib.FFTok_comma {
				state = fflib.FFParse_want_key
			} else if tok == fflib.FFTok_right_bracket {
				goto done
			} else {
				wantedTok = fflib.FFTok_comma
				goto wrongtokenerror
			}

		case fflib.FFParse_want_key:
			// json {} ended. goto exit. woo.
			if tok == fflib.FFTok_right_bracket {
				goto done
			}
			if tok != fflib.FFTok_string {
				wantedTok = fflib.FFTok_string
				goto wrongtokenerror
			}

			kn := fs.Output.Bytes()
			if len(kn) <= 0 {
				// "" case. hrm.
				currentKey = ffjtTarOptionsnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'C':

					if bytes.Equal(ffjKeyTarOptionsCompression, kn) {
						currentKey = ffjtTarOptionsCompression
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsChownOpts, kn) {
						currentKey = ffjtTarOptionsChownOpts
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsCopyPass, kn) {
						currentKey = ffjtTarOptionsCopyPass
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'E':

					if bytes.Equal(ffjKeyTarOptionsExcludePatterns, kn) {
						currentKey = ffjtTarOptionsExcludePatterns
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'G':

					if bytes.Equal(ffjKeyTarOptionsGIDMaps, kn) {
						currentKey = ffjtTarOptionsGIDMaps
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'I':

					if bytes.Equal(ffjKeyTarOptionsIncludeFiles, kn) {
						currentKey = ffjtTarOptionsIncludeFiles
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsIncludeSourceDir, kn) {
						currentKey = ffjtTarOptionsIncludeSourceDir
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsInUserNS, kn) {
						currentKey = ffjtTarOptionsInUserNS
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'N':

					if bytes.Equal(ffjKeyTarOptionsNoLchown, kn) {
						currentKey = ffjtTarOptionsNoLchown
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsNoOverwriteDirNonDir, kn) {
						currentKey = ffjtTarOptionsNoOverwriteDirNonDir
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'R':

					if bytes.Equal(ffjKeyTarOptionsRebaseNames, kn) {
						currentKey = ffjtTarOptionsRebaseNames
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'U':

					if bytes.Equal(ffjKeyTarOptionsUIDMaps, kn) {
						currentKey = ffjtTarOptionsUIDMaps
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'W':

					if bytes.Equal(ffjKeyTarOptionsWhiteoutFormat, kn) {
						currentKey = ffjtTarOptionsWhiteoutFormat
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyTarOptionsWhiteoutData, kn) {
						currentKey = ffjtTarOptionsWhiteoutData
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsCopyPass, kn) {
					currentKey = ffjtTarOptionsCopyPass
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsInUserNS, kn) {
					currentKey = ffjtTarOptionsInUserNS
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsRebaseNames, kn) {
					currentKey = ffjtTarOptionsRebaseNames
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyTarOptionsNoOverwriteDirNonDir, kn) {
					currentKey = ffjtTarOptionsNoOverwriteDirNonDir
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyTarOptionsWhiteoutData, kn) {
					currentKey = ffjtTarOptionsWhiteoutData
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyTarOptionsWhiteoutFormat, kn) {
					currentKey = ffjtTarOptionsWhiteoutFormat
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsIncludeSourceDir, kn) {
					currentKey = ffjtTarOptionsIncludeSourceDir
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsChownOpts, kn) {
					currentKey = ffjtTarOptionsChownOpts
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsGIDMaps, kn) {
					currentKey = ffjtTarOptionsGIDMaps
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsUIDMaps, kn) {
					currentKey = ffjtTarOptionsUIDMaps
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyTarOptionsNoLchown, kn) {
					currentKey = ffjtTarOptionsNoLchown
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsCompression, kn) {
					currentKey = ffjtTarOptionsCompression
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsExcludePatterns, kn) {
					currentKey = ffjtTarOptionsExcludePatterns
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyTarOptionsIncludeFiles, kn) {
					currentKey = ffjtTarOptionsIncludeFiles
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtTarOptionsnosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			}

		case fflib.FFParse_want_colon:
			if tok != fflib.FFTok_colon {
				wantedTok = fflib.FFTok_colon
				goto wrongtokenerror
			}
			state = fflib.FFParse_want_value
			continue
		case fflib.FFParse_want_value:

			if tok == fflib.FFTok_left_brace || tok == fflib.FFTok_left_bracket || tok == fflib.FFTok_integer || tok == fflib.FFTok_double || tok == fflib.FFTok_string || tok == fflib.FFTok_bool || tok == fflib.FFTok_null {
				switch currentKey {

				case ffjtTarOptionsIncludeFiles:
					goto handle_IncludeFiles

				case ffjtTarOptionsExcludePatterns:
					goto handle_ExcludePatterns

				case ffjtTarOptionsCompression:
					goto handle_Compression

				case ffjtTarOptionsNoLchown:
					goto handle_NoLchown

				case ffjtTarOptionsUIDMaps:
					goto handle_UIDMaps

				case ffjtTarOptionsGIDMaps:
					goto handle_GIDMaps

				case ffjtTarOptionsChownOpts:
					goto handle_ChownOpts

				case ffjtTarOptionsIncludeSourceDir:
					goto handle_IncludeSourceDir

				case ffjtTarOptionsWhiteoutFormat:
					goto handle_WhiteoutFormat

				case ffjtTarOptionsWhiteoutData:
					goto handle_WhiteoutData

				case ffjtTarOptionsNoOverwriteDirNonDir:
					goto handle_NoOverwriteDirNonDir

				case ffjtTarOptionsRebaseNames:
					goto handle_RebaseNames

				case ffjtTarOptionsInUserNS:
					goto handle_InUserNS

				case ffjtTarOptionsCopyPass:
					goto handle_CopyPass

				case ffjtTarOptionsnosuchkey:
					err = fs.SkipField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}
					state = fflib.FFParse_after_value
					goto mainparse
				}
			} else {
				goto wantedvalue
			}
		}
	}

handle_IncludeFiles:

	/* handler: j.IncludeFiles type=[]string kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.IncludeFiles = nil
		} else {

			j.IncludeFiles = []string{}

			wantVal := true

			for {

				var tmpJIncludeFiles string

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmpJIncludeFiles type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						tmpJIncludeFiles = string(string(outBuf))

					}
				}

				j.IncludeFiles = append(j.IncludeFiles, tmpJIncludeFiles)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_ExcludePatterns:

	/* handler: j.ExcludePatterns type=[]string kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.ExcludePatterns = nil
		} else {

			j.ExcludePatterns = []string{}

			wantVal := true

			for {

				var tmpJExcludePatterns string

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmpJExcludePatterns type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						tmpJExcludePatterns = string(string(outBuf))

					}
				}

				j.ExcludePatterns = append(j.ExcludePatterns, tmpJExcludePatterns)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Compression:

	/* handler: j.Compression type=archive.Compression kind=int quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for Compression", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.Compression = Compression(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_NoLchown:

	/* handler: j.NoLchown type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				j.NoLchown = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				j.NoLchown = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_UIDMaps:

	/* handler: j.UIDMaps type=[]idtools.IDMap kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.UIDMaps = nil
		} else {

			j.UIDMaps = []idtools.IDMap{}

			wantVal := true

			for {

				var tmpJUIDMaps idtools.IDMap

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmpJUIDMaps type=idtools.IDMap kind=struct quoted=false*/

				{
					/* Falling back. type=idtools.IDMap kind=struct */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmpJUIDMaps)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				j.UIDMaps = append(j.UIDMaps, tmpJUIDMaps)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_GIDMaps:

	/* handler: j.GIDMaps type=[]idtools.IDMap kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.GIDMaps = nil
		} else {

			j.GIDMaps = []idtools.IDMap{}

			wantVal := true

			for {

				var tmpJGIDMaps idtools.IDMap

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_brace {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: tmpJGIDMaps type=idtools.IDMap kind=struct quoted=false*/

				{
					/* Falling back. type=idtools.IDMap kind=struct */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmpJGIDMaps)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				j.GIDMaps = append(j.GIDMaps, tmpJGIDMaps)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_ChownOpts:

	/* handler: j.ChownOpts type=idtools.IDPair kind=struct quoted=false*/

	{
		/* Falling back. type=idtools.IDPair kind=struct */
		tbuf, err := fs.CaptureField(tok)
		if err != nil {
			return fs.WrapErr(err)
		}

		err = json.Unmarshal(tbuf, &j.ChownOpts)
		if err != nil {
			return fs.WrapErr(err)
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_IncludeSourceDir:

	/* handler: j.IncludeSourceDir type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				j.IncludeSourceDir = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				j.IncludeSourceDir = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_WhiteoutFormat:

	/* handler: j.WhiteoutFormat type=archive.WhiteoutFormat kind=int quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for WhiteoutFormat", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.WhiteoutFormat = WhiteoutFormat(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_WhiteoutData:

	/* handler: j.WhiteoutData type=interface {} kind=interface quoted=false*/

	{
		/* Falling back. type=interface {} kind=interface */
		tbuf, err := fs.CaptureField(tok)
		if err != nil {
			return fs.WrapErr(err)
		}

		err = json.Unmarshal(tbuf, &j.WhiteoutData)
		if err != nil {
			return fs.WrapErr(err)
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_NoOverwriteDirNonDir:

	/* handler: j.NoOverwriteDirNonDir type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				j.NoOverwriteDirNonDir = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				j.NoOverwriteDirNonDir = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_RebaseNames:

	/* handler: j.RebaseNames type=map[string]string kind=map quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_bracket && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.RebaseNames = nil
		} else {

			j.RebaseNames = make(map[string]string, 0)

			wantVal := true

			for {

				var k string

				var tmpJRebaseNames string

				tok = fs.Scan()
				if tok == fflib.FFTok_error {
					goto tokerror
				}
				if tok == fflib.FFTok_right_bracket {
					break
				}

				if tok == fflib.FFTok_comma {
					if wantVal == true {
						// TODO(pquerna): this isn't an ideal error message, this handles
						// things like [,,,] as an array value.
						return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
					}
					continue
				} else {
					wantVal = true
				}

				/* handler: k type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						k = string(string(outBuf))

					}
				}

				// Expect ':' after key
				tok = fs.Scan()
				if tok != fflib.FFTok_colon {
					return fs.WrapErr(fmt.Errorf("wanted colon token, but got token: %v", tok))
				}

				tok = fs.Scan()
				/* handler: tmpJRebaseNames type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						tmpJRebaseNames = string(string(outBuf))

					}
				}

				j.RebaseNames[k] = tmpJRebaseNames

				wantVal = false
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_InUserNS:

	/* handler: j.InUserNS type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				j.InUserNS = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				j.InUserNS = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_CopyPass:

	/* handler: j.CopyPass type=bool kind=bool quoted=false*/

	{
		if tok != fflib.FFTok_bool && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for bool", tok))
		}
	}

	{
		if tok == fflib.FFTok_null {

		} else {
			tmpb := fs.Output.Bytes()

			if bytes.Compare([]byte{'t', 'r', 'u', 'e'}, tmpb) == 0 {

				j.CopyPass = true

			} else if bytes.Compare([]byte{'f', 'a', 'l', 's', 'e'}, tmpb) == 0 {

				j.CopyPass = false

			} else {
				err = errors.New("unexpected bytes for true/false value")
				return fs.WrapErr(err)
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

wantedvalue:
	return fs.WrapErr(fmt.Errorf("wanted value token, but got token: %v", tok))
wrongtokenerror:
	return fs.WrapErr(fmt.Errorf("ffjson: wanted token: %v, but got token: %v output=%s", wantedTok, tok, fs.Output.String()))
tokerror:
	if fs.BigError != nil {
		return fs.WrapErr(fs.BigError)
	}
	err = fs.Error.ToError()
	if err != nil {
		return fs.WrapErr(err)
	}
	panic("ffjson-generated: unreachable, please report bug.")
done:

	return nil
}
