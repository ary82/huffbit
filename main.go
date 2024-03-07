package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("needs more args")
	}

	mode := os.Args[1]

	switch mode {
	case "-c":
		fmt.Println("compression mode")
		compress()
	case "-d":
		fmt.Println("decompression mode")
	default:
		fmt.Println("not recognized")
	}

}
