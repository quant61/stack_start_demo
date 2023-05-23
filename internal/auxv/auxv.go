package auxv

import "fmt"

type Auxk uint64

/* copied from /usr/include/x86_64-linux-gnu/bits/auxv.h */
const (
	AT_NULL   Auxk = 0  /* End of vector */
	AT_IGNORE Auxk = 1  /* Entry should be ignored */
	AT_EXECFD Auxk = 2  /* File descriptor of program */
	AT_PHDR   Auxk = 3  /* Program headers for program */
	AT_PHENT  Auxk = 4  /* Size of program header entry */
	AT_PHNUM  Auxk = 5  /* Number of program headers */
	AT_PAGESZ Auxk = 6  /* System page size */
	AT_BASE   Auxk = 7  /* Base address of interpreter */
	AT_FLAGS  Auxk = 8  /* Flags */
	AT_ENTRY  Auxk = 9  /* Entry point of program */
	AT_NOTELF Auxk = 10 /* Program is not ELF */
	AT_UID    Auxk = 11 /* Real uid */
	AT_EUID   Auxk = 12 /* Effective uid */
	AT_GID    Auxk = 13 /* Real gid */
	AT_EGID   Auxk = 14 /* Effective gid */
	AT_CLKTCK Auxk = 17 /* Frequency of times() */

	/* Some more special a_type values describing the hardware.  */
	AT_PLATFORM Auxk = 15 /* String identifying platform.  */
	AT_HWCAP    Auxk = 16 /* Machine-dependent hints about
	   processor capabilities.  */

	/* This entry gives some information about the FPU initialization
	   performed by the kernel.  */
	AT_FPUCW Auxk = 18 /* Used FPU control word.  */

	/* Cache block sizes.  */
	AT_DCACHEBSIZE Auxk = 19 /* Data cache block size.  */
	AT_ICACHEBSIZE Auxk = 20 /* Instruction cache block size.  */
	AT_UCACHEBSIZE Auxk = 21 /* Unified cache block size.  */

	/* A special ignored value for PPC, used by the kernel to control the
	   interpretation of the AUXV. Must be > 16.  */
	AT_IGNOREPPC Auxk = 22 /* Entry should be ignored.  */

	AT_SECURE Auxk = 23 /* Boolean, was exec setuid-like?  */

	AT_BASE_PLATFORM Auxk = 24 /* String identifying real platforms.*/

	AT_RANDOM Auxk = 25 /* Address of 16 random bytes.  */

	AT_HWCAP2 Auxk = 26 /* More machine-dependent hints about
	   processor capabilities.  */

	AT_EXECFN Auxk = 31 /* Filename of executable.  */

	/* Pointer to the global system page used for system calls and other
	   nice things.  */
	AT_SYSINFO      Auxk = 32
	AT_SYSINFO_EHDR Auxk = 33

	/* Shapes of the caches.  Bits 0-3 contains associativity; bits 4-7 contains
	   log2 of line size; mask those to get cache size.  */
	AT_L1I_CACHESHAPE Auxk = 34
	AT_L1D_CACHESHAPE Auxk = 35
	AT_L2_CACHESHAPE  Auxk = 36
	AT_L3_CACHESHAPE  Auxk = 37

	/* Shapes of the caches, with more room to describe them.
	 *GEOMETRY are comprised of cache line size in bytes in the bottom 16 bits
	 and the cache associativity in the next 16 bits.  */
	AT_L1I_CACHESIZE     Auxk = 40
	AT_L1I_CACHEGEOMETRY Auxk = 41
	AT_L1D_CACHESIZE     Auxk = 42
	AT_L1D_CACHEGEOMETRY Auxk = 43
	AT_L2_CACHESIZE      Auxk = 44
	AT_L2_CACHEGEOMETRY  Auxk = 45
	AT_L3_CACHESIZE      Auxk = 46
	AT_L3_CACHEGEOMETRY  Auxk = 47

	AT_MINSIGSTKSZ Auxk = 51 /* Stack needed for signal delivery
	   (AArch64).  */
)

// from the same file using:
// ```javascript
// [...s.matchAll(/\#define\s*(\w+)\s*(\d+).*/g)].map(m=>`    ${m[2]}: "${m[1]}",`).join("\n")
// ```
var auxkNames = map[Auxk]string{
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

func IsCStringPtr(k Auxk) bool {
	switch k {
	// note: not all are probably listed; false negatives are possible
	case AT_PLATFORM, AT_BASE_PLATFORM, AT_EXECFN:
		return true
	default:
		return false
	}
}

func AuxkNames() map[Auxk]string {
	m := make(map[Auxk]string, len(auxkNames))
	for k, v := range auxkNames {
		m[k] = v
	}
	return m
}

func(k Auxk) String() string {
	name := auxkNames[k]
	if name != "" {
		return fmt.Sprintf("%s(%d)", name, k)
	}
	return fmt.Sprintf("%d", k)
}