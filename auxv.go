package main

import "fmt"

type auxk uint64

/* copied from /usr/include/x86_64-linux-gnu/bits/auxv.h */
const (
	AT_NULL   auxk = 0  /* End of vector */
	AT_IGNORE auxk = 1  /* Entry should be ignored */
	AT_EXECFD auxk = 2  /* File descriptor of program */
	AT_PHDR   auxk = 3  /* Program headers for program */
	AT_PHENT  auxk = 4  /* Size of program header entry */
	AT_PHNUM  auxk = 5  /* Number of program headers */
	AT_PAGESZ auxk = 6  /* System page size */
	AT_BASE   auxk = 7  /* Base address of interpreter */
	AT_FLAGS  auxk = 8  /* Flags */
	AT_ENTRY  auxk = 9  /* Entry point of program */
	AT_NOTELF auxk = 10 /* Program is not ELF */
	AT_UID    auxk = 11 /* Real uid */
	AT_EUID   auxk = 12 /* Effective uid */
	AT_GID    auxk = 13 /* Real gid */
	AT_EGID   auxk = 14 /* Effective gid */
	AT_CLKTCK auxk = 17 /* Frequency of times() */

	/* Some more special a_type values describing the hardware.  */
	AT_PLATFORM auxk = 15 /* String identifying platform.  */
	AT_HWCAP    auxk = 16 /* Machine-dependent hints about
	   processor capabilities.  */

	/* This entry gives some information about the FPU initialization
	   performed by the kernel.  */
	AT_FPUCW auxk = 18 /* Used FPU control word.  */

	/* Cache block sizes.  */
	AT_DCACHEBSIZE auxk = 19 /* Data cache block size.  */
	AT_ICACHEBSIZE auxk = 20 /* Instruction cache block size.  */
	AT_UCACHEBSIZE auxk = 21 /* Unified cache block size.  */

	/* A special ignored value for PPC, used by the kernel to control the
	   interpretation of the AUXV. Must be > 16.  */
	AT_IGNOREPPC auxk = 22 /* Entry should be ignored.  */

	AT_SECURE auxk = 23 /* Boolean, was exec setuid-like?  */

	AT_BASE_PLATFORM auxk = 24 /* String identifying real platforms.*/

	AT_RANDOM auxk = 25 /* Address of 16 random bytes.  */

	AT_HWCAP2 auxk = 26 /* More machine-dependent hints about
	   processor capabilities.  */

	AT_EXECFN auxk = 31 /* Filename of executable.  */

	/* Pointer to the global system page used for system calls and other
	   nice things.  */
	AT_SYSINFO      auxk = 32
	AT_SYSINFO_EHDR auxk = 33

	/* Shapes of the caches.  Bits 0-3 contains associativity; bits 4-7 contains
	   log2 of line size; mask those to get cache size.  */
	AT_L1I_CACHESHAPE auxk = 34
	AT_L1D_CACHESHAPE auxk = 35
	AT_L2_CACHESHAPE  auxk = 36
	AT_L3_CACHESHAPE  auxk = 37

	/* Shapes of the caches, with more room to describe them.
	 *GEOMETRY are comprised of cache line size in bytes in the bottom 16 bits
	 and the cache associativity in the next 16 bits.  */
	AT_L1I_CACHESIZE     auxk = 40
	AT_L1I_CACHEGEOMETRY auxk = 41
	AT_L1D_CACHESIZE     auxk = 42
	AT_L1D_CACHEGEOMETRY auxk = 43
	AT_L2_CACHESIZE      auxk = 44
	AT_L2_CACHEGEOMETRY  auxk = 45
	AT_L3_CACHESIZE      auxk = 46
	AT_L3_CACHEGEOMETRY  auxk = 47

	AT_MINSIGSTKSZ auxk = 51 /* Stack needed for signal delivery
	   (AArch64).  */
)

// from the same file using:
// ```javascript
// [...s.matchAll(/\#define\s*(\w+)\s*(\d+).*/g)].map(m=>`    ${m[2]}: "${m[1]}",`).join("\n")
// ```
var auxkNames = map[auxk]string{
	0: "AT_NULL",
	1: "AT_IGNORE",
	2: "AT_EXECFD",
	3: "AT_PHDR",
	4: "AT_PHENT",
	5: "AT_PHNUM",
	6: "AT_PAGESZ",
	7: "AT_BASE",
	8: "AT_FLAGS",
	9: "AT_ENTRY",
	10: "AT_NOTELF",
	11: "AT_UID",
	12: "AT_EUID",
	13: "AT_GID",
	14: "AT_EGID",
	17: "AT_CLKTCK",
	15: "AT_PLATFORM",
	16: "AT_HWCAP",
	18: "AT_FPUCW",
	19: "AT_DCACHEBSIZE",
	20: "AT_ICACHEBSIZE",
	21: "AT_UCACHEBSIZE",
	22: "AT_IGNOREPPC",
	23: "AT_SECURE",
	24: "AT_BASE_PLATFORM",
	25: "AT_RANDOM",
	26: "AT_HWCAP2",
	31: "AT_EXECFN",
	32: "AT_SYSINFO",
	33: "AT_SYSINFO_EHDR",
	34: "AT_L1I_CACHESHAPE",
	35: "AT_L1D_CACHESHAPE",
	36: "AT_L2_CACHESHAPE",
	37: "AT_L3_CACHESHAPE",
	40: "AT_L1I_CACHESIZE",
	41: "AT_L1I_CACHEGEOMETRY",
	42: "AT_L1D_CACHESIZE",
	43: "AT_L1D_CACHEGEOMETRY",
	44: "AT_L2_CACHESIZE",
	45: "AT_L2_CACHEGEOMETRY",
	46: "AT_L3_CACHESIZE",
	47: "AT_L3_CACHEGEOMETRY",
	51: "AT_MINSIGSTKSZ",
}

func(k auxk) String() string {
	name := auxkNames[k]
	if name != "" {
		return fmt.Sprintf("%s(%d)", name, k)
	}
	return fmt.Sprintf("%d", k)
}