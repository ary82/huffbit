package main

import (
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
	left  *HuffmanEle
	right *HuffmanEle
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

type PriorityQueue []HuffmanEle

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i int, j int) bool {
	return pq[i].getValue() < pq[j].getValue()
}
func (pq PriorityQueue) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(i any) {
	*pq = append(*pq, i.(HuffmanEle))
}
func (pq *PriorityQueue) Pop() any {
	poppedEle := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return poppedEle
}
func makeHuffTree() {

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

	makeHuffTree()

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
