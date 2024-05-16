package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Primary function for decompression
func decompress() {
	fmt.Println("decompressing file")

	if len(os.Args) < 3 {
		panic("No file given")
	}
	filedata, err := os.ReadFile(os.Args[2])
	panicErr(err)

	freqMap, endOfHeader := parseHeader(filedata)

	// Total characters
	totalChar := 0
	for _, freq := range freqMap {
		totalChar += freq
	}

	huffmanTree := makeHuffTree(freqMap)

	// Make codeMap from HuffmanTree
	codeMap := make(map[rune][]byte)
	getCodes(huffmanTree, []byte{}, codeMap)

	newcodeMap := make(map[string]rune)
	for key, value := range codeMap {
		newcodeMap[string(value)] = key
	}

	// Create new file
	outputFile, err := os.Create(fmt.Sprint(filepath.Base(os.Args[2]), ".original"))
	panicErr(err)

	// Passing filedata except the header, then slicing till the total number of characters in the original file.
	toBeWritten := (getUncompressedData(filedata[endOfHeader:], newcodeMap))[:totalChar]
	_, err = outputFile.Write(toBeWritten)
	panicErr(err)
}

func getUncompressedData(compressedData []byte, codeMap map[string]rune) []byte {
	var uncompressedData []byte
	var buffer []byte

	for _, v := range compressedData {
		for i := 0; i < 8; i++ {

			if (v>>(7-i))&1 != 0 {
				buffer = append(buffer, 49)
			} else {
				buffer = append(buffer, 48)
			}

			// Check if value exists in map
			value, exists := codeMap[string(buffer)]
			if exists {
				uncompressedData = append(uncompressedData, byte(value))
				buffer = []byte{}
			}
		}
	}
	return uncompressedData
}

// Parses Header to a Frequency map, gives index of start of compressed data
func parseHeader(header []byte) (map[rune]int, int) {
	freqMap := make(map[rune]int)
	var end int
	var count int
	var char int

	for i, v := range header {
		switch v {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			count *= 10
			count = count + int(v) - 48
			if header[i+1] == ':' {
				char = count
				count = 0
			} else if header[i+1] == ' ' || header[i+1] == ']' {
				freqMap[rune(char)] = count
				count = 0
				char = 0
			}
		case ']':
			end = i + 1
			return freqMap, end
		}
	}
	return freqMap, end
}
