func (j *Layer) UnmarshalJSONFFLexer(fs *fflib.FFLexer, state fflib.FFParseState) error {
	var err error
	currentKey := ffjtLayerbase
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
				currentKey = ffjtLayernosuchkey
				state = fflib.FFParse_want_colon
				goto mainparse
			} else {
				switch kn[0] {

				case 'c':

					if bytes.Equal(ffjKeyLayerCreated, kn) {
						currentKey = ffjtLayerCreated
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyLayerCompressedDigest, kn) {
						currentKey = ffjtLayerCompressedDigest
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyLayerCompressedSize, kn) {
						currentKey = ffjtLayerCompressedSize
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyLayerCompressionType, kn) {
						currentKey = ffjtLayerCompressionType
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'd':

					if bytes.Equal(ffjKeyLayerUncompressedDigest, kn) {
						currentKey = ffjtLayerUncompressedDigest
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyLayerUncompressedSize, kn) {
						currentKey = ffjtLayerUncompressedSize
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'f':

					if bytes.Equal(ffjKeyLayerFlags, kn) {
						currentKey = ffjtLayerFlags
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'g':

					if bytes.Equal(ffjKeyLayerGIDMap, kn) {
						currentKey = ffjtLayerGIDMap
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'i':

					if bytes.Equal(ffjKeyLayerID, kn) {
						currentKey = ffjtLayerID
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'm':

					if bytes.Equal(ffjKeyLayerMetadata, kn) {
						currentKey = ffjtLayerMetadata
						state = fflib.FFParse_want_colon
						goto mainparse

					} else if bytes.Equal(ffjKeyLayerMountLabel, kn) {
						currentKey = ffjtLayerMountLabel
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'n':

					if bytes.Equal(ffjKeyLayerNames, kn) {
						currentKey = ffjtLayerNames
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'p':

					if bytes.Equal(ffjKeyLayerParent, kn) {
						currentKey = ffjtLayerParent
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				case 'u':

					if bytes.Equal(ffjKeyLayerUIDMap, kn) {
						currentKey = ffjtLayerUIDMap
						state = fflib.FFParse_want_colon
						goto mainparse
					}

				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerGIDMap, kn) {
					currentKey = ffjtLayerGIDMap
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerUIDMap, kn) {
					currentKey = ffjtLayerUIDMap
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerFlags, kn) {
					currentKey = ffjtLayerFlags
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerCompressionType, kn) {
					currentKey = ffjtLayerCompressionType
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerUncompressedSize, kn) {
					currentKey = ffjtLayerUncompressedSize
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerUncompressedDigest, kn) {
					currentKey = ffjtLayerUncompressedDigest
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerCompressedSize, kn) {
					currentKey = ffjtLayerCompressedSize
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerCompressedDigest, kn) {
					currentKey = ffjtLayerCompressedDigest
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerCreated, kn) {
					currentKey = ffjtLayerCreated
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerMountLabel, kn) {
					currentKey = ffjtLayerMountLabel
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerMetadata, kn) {
					currentKey = ffjtLayerMetadata
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerParent, kn) {
					currentKey = ffjtLayerParent
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.EqualFoldRight(ffjKeyLayerNames, kn) {
					currentKey = ffjtLayerNames
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				if fflib.SimpleLetterEqualFold(ffjKeyLayerID, kn) {
					currentKey = ffjtLayerID
					state = fflib.FFParse_want_colon
					goto mainparse
				}

				currentKey = ffjtLayernosuchkey
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

				case ffjtLayerID:
					goto handle_ID

				case ffjtLayerNames:
					goto handle_Names

				case ffjtLayerParent:
					goto handle_Parent

				case ffjtLayerMetadata:
					goto handle_Metadata

				case ffjtLayerMountLabel:
					goto handle_MountLabel

				case ffjtLayerCreated:
					goto handle_Created

				case ffjtLayerCompressedDigest:
					goto handle_CompressedDigest

				case ffjtLayerCompressedSize:
					goto handle_CompressedSize

				case ffjtLayerUncompressedDigest:
					goto handle_UncompressedDigest

				case ffjtLayerUncompressedSize:
					goto handle_UncompressedSize

				case ffjtLayerCompressionType:
					goto handle_CompressionType

				case ffjtLayerFlags:
					goto handle_Flags

				case ffjtLayerUIDMap:
					goto handle_UIDMap

				case ffjtLayerGIDMap:
					goto handle_GIDMap

				case ffjtLayernosuchkey:
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

handle_ID:

	/* handler: j.ID type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.ID = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Names:

	/* handler: j.Names type=[]string kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.Names = nil
		} else {

			j.Names = []string{}

			wantVal := true

			for {

				var tmpJNames string

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

				/* handler: tmpJNames type=string kind=string quoted=false*/

				{

					{
						if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
							return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
						}
					}

					if tok == fflib.FFTok_null {

					} else {

						outBuf := fs.Output.Bytes()

						tmpJNames = string(string(outBuf))

					}
				}

				j.Names = append(j.Names, tmpJNames)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Parent:

	/* handler: j.Parent type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Parent = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Metadata:

	/* handler: j.Metadata type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.Metadata = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_MountLabel:

	/* handler: j.MountLabel type=string kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for string", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.MountLabel = string(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Created:

	/* handler: j.Created type=time.Time kind=struct quoted=false*/

	{
		if tok == fflib.FFTok_null {

		} else {

			tbuf, err := fs.CaptureField(tok)
			if err != nil {
				return fs.WrapErr(err)
			}

			err = j.Created.UnmarshalJSON(tbuf)
			if err != nil {
				return fs.WrapErr(err)
			}
		}
		state = fflib.FFParse_after_value
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_CompressedDigest:

	/* handler: j.CompressedDigest type=digest.Digest kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for Digest", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.CompressedDigest = digest.Digest(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_CompressedSize:

	/* handler: j.CompressedSize type=int64 kind=int64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int64", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.CompressedSize = int64(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_UncompressedDigest:

	/* handler: j.UncompressedDigest type=digest.Digest kind=string quoted=false*/

	{

		{
			if tok != fflib.FFTok_string && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for Digest", tok))
			}
		}

		if tok == fflib.FFTok_null {

		} else {

			outBuf := fs.Output.Bytes()

			j.UncompressedDigest = digest.Digest(string(outBuf))

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_UncompressedSize:

	/* handler: j.UncompressedSize type=int64 kind=int64 quoted=false*/

	{
		if tok != fflib.FFTok_integer && tok != fflib.FFTok_null {
			return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for int64", tok))
		}
	}

	{

		if tok == fflib.FFTok_null {

		} else {

			tval, err := fflib.ParseInt(fs.Output.Bytes(), 10, 64)

			if err != nil {
				return fs.WrapErr(err)
			}

			j.UncompressedSize = int64(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_CompressionType:

	/* handler: j.CompressionType type=archive.Compression kind=int quoted=false*/

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

			j.CompressionType = archive.Compression(tval)

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_Flags:

	/* handler: j.Flags type=map[string]interface {} kind=map quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_bracket && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.Flags = nil
		} else {

			j.Flags = make(map[string]interface{}, 0)

			wantVal := true

			for {

				var k string

				var tmpJFlags interface{}

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
				/* handler: tmpJFlags type=interface {} kind=interface quoted=false*/

				{
					/* Falling back. type=interface {} kind=interface */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmpJFlags)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				j.Flags[k] = tmpJFlags

				wantVal = false
			}

		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_UIDMap:

	/* handler: j.UIDMap type=[]idtools.IDMap kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.UIDMap = nil
		} else {

			j.UIDMap = []idtools.IDMap{}

			wantVal := true

			for {

				var tmpJUIDMap idtools.IDMap

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

				/* handler: tmpJUIDMap type=idtools.IDMap kind=struct quoted=false*/

				{
					/* Falling back. type=idtools.IDMap kind=struct */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmpJUIDMap)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				j.UIDMap = append(j.UIDMap, tmpJUIDMap)

				wantVal = false
			}
		}
	}

	state = fflib.FFParse_after_value
	goto mainparse

handle_GIDMap:

	/* handler: j.GIDMap type=[]idtools.IDMap kind=slice quoted=false*/

	{

		{
			if tok != fflib.FFTok_left_brace && tok != fflib.FFTok_null {
				return fs.WrapErr(fmt.Errorf("cannot unmarshal %s into Go value for ", tok))
			}
		}

		if tok == fflib.FFTok_null {
			j.GIDMap = nil
		} else {

			j.GIDMap = []idtools.IDMap{}

			wantVal := true

			for {

				var tmpJGIDMap idtools.IDMap

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

				/* handler: tmpJGIDMap type=idtools.IDMap kind=struct quoted=false*/

				{
					/* Falling back. type=idtools.IDMap kind=struct */
					tbuf, err := fs.CaptureField(tok)
					if err != nil {
						return fs.WrapErr(err)
					}

					err = json.Unmarshal(tbuf, &tmpJGIDMap)
					if err != nil {
						return fs.WrapErr(err)
					}
				}

				j.GIDMap = append(j.GIDMap, tmpJGIDMap)

				wantVal = false
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
