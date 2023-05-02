package main

import (
	"flag"
	"fmt"
	"github.com/quant61/stack_start_demo/internal/binfmt/elf"
)

func main() {
	//run()
	// TODO: add more platforms
	// TODO: make it configurable
	// TODO: add foreign platform support
	mode := flag.String("mode", "both", "32 or 64-bit mode")
	flag.Parse()

	switch *mode {
	case "32":
		fmt.Println("RUN FOR 32-BIT")
		elfRun(elf.BuildElfBinary32, 4)
	case "64":
		fmt.Println("RUN FOR 64-BIT")
		elfRun(elf.BuildElfBinary64, 8)
	case "both":
		fmt.Println("RUN FOR 32-BIT")
		elfRun(elf.BuildElfBinary32, 4)
		fmt.Println("RUN FOR 64-BIT")
		elfRun(elf.BuildElfBinary64, 8)
	}
}
