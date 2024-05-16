package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Needs more args")
	}

	mode := os.Args[1]

	switch mode {
	case "-c":
		compress()
	case "-d":
		decompress()
	default:
		fmt.Println("not a recognized option")
	}
}
