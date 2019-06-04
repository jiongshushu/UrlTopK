package main

import "fmt"

type UrlNumNode struct {
	Url        string
	Num        int
	Pos        int
	LeftChild  *UrlNumNode
	RightChild *UrlNumNode
}

type MinHeap struct {
	Nodes []*UrlNumNode
	Pos   int
	K     int
}

func (h *MinHeap) BuildMinHeap(data []*UrlNumNode, k, pos int) {
	h.Nodes = make([]*UrlNumNode, 0, k)

	for i, d := range data {
		if pos != -1 {
			d.Pos = pos
		}
		if i < k {
			h.Nodes = append([]*UrlNumNode{d}, h.Nodes...)
			h.MinAdjust(0)
		} else {
			if d.Num > h.Nodes[0].Num {
				h.Nodes[0] = d
				h.MinAdjust(0)
			}
		}
	}
	h.Pos = pos
	h.K = k
}

func (h *MinHeap) MinAdjust(i int) {
	root := i
	left, right := 2*root+1, len(h.Nodes)-1
	min := left
	for left <= right {
		if left+1 <= right && h.Nodes[left].Num > h.Nodes[left+1].Num {
			min = left + 1
		}
		if h.Nodes[root].Num < h.Nodes[min].Num {
			break
		} else {
			h.Nodes[root], h.Nodes[min] = h.Nodes[min], h.Nodes[root]
			root = min
			left = 2*root + 1
			min = left
		}
	}
}

func (h *MinHeap) Print() {
	fmt.Println("============")
	for i, d := range h.Nodes {
		fmt.Printf("%d, %s, %d, %d\n", i, d.Url, d.Num, d.Pos)
	}
	fmt.Println("============")
}

type MaxHeap struct {
	Nodes []*UrlNumNode
	K     int
	Pos   int
}

func (h *MaxHeap) BuildMaxHeap(data []*UrlNumNode, k, pos int) {
	h.Nodes = make([]*UrlNumNode, 0, k)

	for i, d := range data {
		if pos != -1 {
			d.Pos = pos
		}
		if i < k {
			h.Nodes = append([]*UrlNumNode{d}, h.Nodes...)
			h.MaxAdjust(0)
		} else {
			if d.Num > h.Nodes[0].Num {
				h.Nodes[0] = d
				h.MaxAdjust(0)
			}
		}
	}
	h.Pos = pos
	h.K = k
}

func (h *MaxHeap) MaxAdjust(i int) {
	root := i
	left, right := 2*root+1, len(h.Nodes)-1
	max := left
	for left <= right {
		if left+1 <= right && h.Nodes[left].Num <= h.Nodes[left+1].Num {
			max = left + 1
		}
		if h.Nodes[root].Num >= h.Nodes[max].Num {
			break
		} else {
			h.Nodes[root], h.Nodes[max] = h.Nodes[max], h.Nodes[root]
			root = max
			left = 2*root + 1
			left = max
		}
	}
}

func (h *MaxHeap) GetNextMaxNode(i int) *UrlNumNode {
	n := len(h.Nodes)
	if i >= n {
		return nil
	}

	left, right := 2*i+1, 2*i+2
	if left < n && right < n {
		if h.Nodes[left].Num > h.Nodes[right].Num {
			return h.Nodes[left]
		} else {
			return h.Nodes[right]
		}
	} else if left < n {
		return h.Nodes[left]
	} else if right < n {
		return h.Nodes[right]
	}

	return nil
}

func (h *MaxHeap) Remove(i int) {
	n := len(h.Nodes)
	if i >= n {
		return
	}
	h.Nodes[i] = h.Nodes[n-1]
	h.Nodes = h.Nodes[:n-1]
	h.MaxAdjust(i)
}

func (h *MaxHeap) Add(node *UrlNumNode) {
	h.Nodes[0] = node
	h.MaxAdjust(0)
}

func (h *MaxHeap) Print() {
	fmt.Println("============")
	for i, d := range h.Nodes {
		fmt.Printf("%d, %s, %d, %d\n", i, d.Url, d.Num, d.Pos)
	}
	fmt.Println("============")
}
