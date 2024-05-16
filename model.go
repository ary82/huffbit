package main

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
