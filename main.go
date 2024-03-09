package main

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// We need two types of structs implementing HuffmanEle: Node and Leaf
type HuffmanEle interface {
	getValue() int
}

type Node struct {
	freq  int
	left  HuffmanEle
	right HuffmanEle
}

type Leaf struct {
	freq int
	char rune
}

type HuffmanHeap []HuffmanEle

// Functions to implement Interfaces
func (n Node) getValue() int {
	return n.freq
}
func (l Leaf) getValue() int {
	return l.freq
}

func (h HuffmanHeap) Len() int {
	return len(h)
}
func (h HuffmanHeap) Less(i int, j int) bool {
	return h[i].getValue() < h[j].getValue()
}
func (h HuffmanHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *HuffmanHeap) Push(i any) {
	*h = append(*h, i.(HuffmanEle))
}
func (h *HuffmanHeap) Pop() any {
	poppedEle := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return poppedEle
}

// Extra Function for Line:68
func (h HuffmanHeap) LessLeaf(i int, j int) bool {
	return int(h[i].(Leaf).char) < int(h[j].(Leaf).char)
}

// Makes the huffman tree, returns the main mode
func makeHuffTree(freqmap map[rune]int) HuffmanEle {
	var huff HuffmanHeap
	for i, v := range freqmap {
		huff = append(huff, Leaf{v, i})
	}
	// Sort before initiating heap for uniformity
	sort.Slice(huff, huff.LessLeaf)
	heap.Init(&huff)

	for huff.Len() > 1 {
		lowest := heap.Pop(&huff).(HuffmanEle)
		secondLowest := heap.Pop(&huff).(HuffmanEle)
		heap.Push(&huff, Node{lowest.getValue() + secondLowest.getValue(), lowest, secondLowest})
	}

	return heap.Pop(&huff).(HuffmanEle)
}

func getCodes(h HuffmanEle, currentCode []byte, codeMap map[rune][]byte) {
	switch i := h.(type) {
	case Node:
		currentCode = append(currentCode, '0')
		getCodes(i.left, currentCode, codeMap)
		currentCode = currentCode[:len(currentCode)-1]

		currentCode = append(currentCode, '1')
		getCodes(i.right, currentCode, codeMap)
		currentCode = currentCode[:len(currentCode)-1]

	case Leaf:
		// Assign a copy of the currentCode as the value of current char in map
		temp := make([]byte, len(currentCode))
		copy(temp, currentCode)
		codeMap[i.char] = temp
	}
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
	// Prevent Data Loss at the cost of appending some garbage characters at the end of the file
	if written != 0 {
		currentByte = currentByte << (8 - written)
		compressedData = append(compressedData, currentByte)
	}
	return compressedData
}

func compress() {
	fmt.Println("compression mode")

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

func decompress() {
	fmt.Println("decompression mode")

	if len(os.Args) < 3 {
		panic("No file given")
	}
	filedata, err := os.ReadFile(os.Args[2])
	panicErr(err)

	freqMap, endOfHeader := parseHeader(filedata)

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

	toBeWritten := getUncompressedData(filedata[endOfHeader:], newcodeMap)
	_, err = outputFile.Write(toBeWritten)
	panicErr(err)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

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
