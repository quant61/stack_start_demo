package main

import "fmt"

func main() {
	//run()
	// TODO: add more platforms
	// TODO: make it configurable
	// TODO: add foreign platform support
	//
	fmt.Println("RUN FOR 32-BIT")
	elfRun(buildElfBinary32, 4)
	fmt.Println("RUN FOR 64-BIT")
	elfRun(buildElfBinary64, 8)
}
