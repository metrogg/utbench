package main

func init() {
	RegisterCharset(&Charset{
		Name:    "GBK",
		Aliases: []string{"GB2312"}, // GBK is a superset of GB2312.
		NewDecoder: func() Decoder {
			return decodeGBKRune
		},
		NewEncoder: func() Encoder {
			return encodeGBKRune
		},
	})
}
