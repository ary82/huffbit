package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Primary function for compression
func compress() {
	fmt.Println("compressing file")

	if len(os.Args) < 3 {
		panic("No file given")
	}

	filedata, err := os.ReadFile(os.Args[2])
	panicErr(err)

	freqMap := make(map[rune]int)
	for _, v := range filedata {
		freqMap[rune(v)] += 1
	}

	huffmanTree := makeHuffTree(freqMap)
	codeMap := make(map[rune][]byte)

	getCodes(huffmanTree, []byte{}, codeMap)

	// Create new file
	outputFile, err := os.Create(fmt.Sprint(filepath.Base(os.Args[2]), ".huffbit"))
	panicErr(err)

	// Write Header(Huffman tree). Slicing to avoid writing 'map'
	header := (fmt.Sprint(freqMap))[3:]
	_, err = outputFile.Write([]byte(header))
	panicErr(err)

	// Write Compressed data
	toBeWritten := getCompressedData(filedata, codeMap)
	_, err = outputFile.Write(toBeWritten)
	panicErr(err)
}

func getCompressedData(file []byte, codeMap map[rune][]byte) []byte {
	var currentByte byte
	var compressedData []byte
	written := 0

	for _, v := range file {
		for _, b := range codeMap[rune(v)] {
			if written == 8 {
				written = 0
				compressedData = append(compressedData, currentByte)
			}
			if b == '1' {
				currentByte = currentByte << 1
				currentByte += 1
				written += 1
			} else {
				currentByte = currentByte << 1
				written += 1
			}
		}
	}
	// Prevent character loss at the end at the cost of appending some garbage characters at the end of the file.
	// These are then handled in the decompression phase.
	if written != 0 {
		currentByte = currentByte << (8 - written)
		compressedData = append(compressedData, currentByte)
	}
	return compressedData
}
