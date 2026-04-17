package main

func haveAVX() bool {
	_, _, c, _ := cpuid(1)

	// Check XGETBV, OXSAVE and AVX bits
	if c&(1<<26) != 0 && c&(1<<27) != 0 && c&(1<<28) != 0 {
		// Check for OS support
		eax, _ := xgetbv(0)
		return (eax & 0x6) == 0x6
	}
	return false
}
