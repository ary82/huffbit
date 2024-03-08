package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
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

func (n Node) getValue() int {
	return n.freq
}
func (l Leaf) getValue() int {
	return l.freq
}

type HuffmanHeap []HuffmanEle

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

func makeHuffTree(freqmap map[rune]int) HuffmanEle {
	var huff HuffmanHeap
	for i, v := range freqmap {
		huff = append(huff, Leaf{v, i})
	}
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
		currentCode = append(currentCode, 0)
		getCodes(i.left, currentCode, codeMap)
		currentCode = currentCode[:len(currentCode)-1]

		currentCode = append(currentCode, 1)
		getCodes(i.right, currentCode, codeMap)
		currentCode = currentCode[:len(currentCode)-1]

	case Leaf:
    // Assign a copy of the currentCode as the value of current char in map
    temp := make([]byte, len(currentCode))
    copy(temp, currentCode)
    codeMap[i.char] = temp
	}

}

func compress() {

	if len(os.Args) < 3 {
		log.Fatal("No file given")
	}

	filedata, err := os.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	freqMap := make(map[rune]int)
	for _, v := range filedata {
		freqMap[rune(v)] += 1
	}
	fmt.Println(freqMap)

	a := makeHuffTree(freqMap)
	codeMap := make(map[rune][]byte)

	getCodes(a, []byte{}, codeMap)
	fmt.Println(codeMap)
}

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
