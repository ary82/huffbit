package main

import (
	"container/heap"
	"sort"
)

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

// Sort helper Function
func (h HuffmanHeap) LessLeaf(i int, j int) bool {
	return int(h[i].(Leaf).char) < int(h[j].(Leaf).char)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
