package main

func initContextLUTs() {
	for i := 0; i < 256; i++ {
		for m := 0; m < numContextModes; m++ {
			base := m << 8

			// Operations performed here are specified in RFC section 7.1.
			switch m {
			case contextLSB6:
				contextP1LUT[base+i] = byte(i) & 0x3f
				contextP2LUT[base+i] = 0
			case contextMSB6:
				contextP1LUT[base+i] = byte(i) >> 2
				contextP2LUT[base+i] = 0
			case contextUTF8:
				contextP1LUT[base+i] = contextLUT0[byte(i)]
				contextP2LUT[base+i] = contextLUT1[byte(i)]
			case contextSigned:
				contextP1LUT[base+i] = contextLUT2[byte(i)] << 3
				contextP2LUT[base+i] = contextLUT2[byte(i)]
			default:
				panic("unknown context mode")
			}
		}
	}
}
