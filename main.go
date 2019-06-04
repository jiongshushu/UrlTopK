package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"math"
	"os"
)

func urlTopK(bigFilePath string, k, maxCap int) {
	if k <= 0 || maxCap <= 0 {
		return
	}
	smallFileNames := splitFile(bigFilePath, maxCap)

	urlMinHeaps := make([]*MinHeap, 0, len(smallFileNames))
	for i, smallFileName := range smallFileNames {
		urlMinHeaps = append(urlMinHeaps, calcSmallFileTopK(smallFileName, i, k))
	}
	topKUrls := mergeSmallFileTopK(urlMinHeaps, k)
	for i, url := range topKUrls {
		fmt.Printf("No.%d %s --- %d\n", i+1, url.Url, url.Num)
	}
}

func calcSmallFileTopK(smallFileName string, pos, k int) *MinHeap {
	file, err := os.Open(smallFileName)
	if err != nil {
		fmt.Errorf("[urlTopK] fail to open fail: %v", err)
	}
	defer file.Close()

	urlNumHashMap := make(map[string]int, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		urlNumHashMap[line]++
	}

	data := make([]*UrlNumNode, 0, len(urlNumHashMap))
	for url, num := range urlNumHashMap {
		d := &UrlNumNode{
			Url: url,
			Num: num,
		}
		data = append(data, d)
	}

	return sortTopK(data, pos, k)
}

func sortTopK(data []*UrlNumNode, pos, k int) *MinHeap {
	h := &MinHeap{}
	h.BuildMinHeap(data, k, pos)
	return h
}

func mergeSmallFileTopK(urlMinHeaps []*MinHeap, k int) []*UrlNumNode {
	urlMaxHeaps := make([]*MaxHeap, 0, len(urlMinHeaps))
	for _, minHeap := range urlMinHeaps {
		maxHeap := &MaxHeap{}
		maxHeap.BuildMaxHeap(minHeap.Nodes, minHeap.K, minHeap.Pos)
		urlMaxHeaps = append(urlMaxHeaps, maxHeap)
	}

	data := make([]*UrlNumNode, 0, len(urlMaxHeaps))
	for _, maxHeap := range urlMaxHeaps {
		data = append(data, maxHeap.Nodes[0])
	}

	totalMaxHeap := &MaxHeap{}
	totalMaxHeap.BuildMaxHeap(data, len(data), -1)

	topKUrls := make([]*UrlNumNode, 0, k)
	nilUrlNumNode := &UrlNumNode{
		Url: "",
		Num: 0,
	}
	var (
		tmpHeap  *MaxHeap
		nextNode *UrlNumNode
	)
	for i := 0; i < k; i++ {
		node := totalMaxHeap.Nodes[0]
		if node == nilUrlNumNode {
			break
		}
		topKUrls = append(topKUrls, node)

		tmpHeap = urlMaxHeaps[node.Pos]

		tmpHeap.Remove(0)
		if len(tmpHeap.Nodes) > 0 {
			nextNode = tmpHeap.Nodes[0]
		} else {
			nextNode = nil
		}
		if nextNode == nil {
			totalMaxHeap.Add(nilUrlNumNode)
		} else {
			totalMaxHeap.Add(nextNode)
		}
	}

	return topKUrls
}

func splitFile(filePath string, maxCap int) []string {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("[urlTopK] fail to open fail: %v", err)
	}

	fileScanner := bufio.NewScanner(file)

	var smallFileName string
	smallFileNameMap := make(map[string]int, 0)

	fdMap := make(map[string]*os.File, 0)

	fileWriterMap := make(map[string]*bufio.Writer, 0)

	defer func() {
		file.Close()
		for _, fd := range fdMap {
			fd.Close()
		}
	}()

	mod := uint32(math.Ceil(float64(fileSize) / float64(maxCap) * 2)) // half

	var smallFileWriter *bufio.Writer
	for fileScanner.Scan() {
		line := fileScanner.Text()

		h := fnv.New32a()
		h.Write([]byte(line))
		num := h.Sum32()
		smallFileName = fmt.Sprintf("small_file%d.txt", num%mod)

		line += "\n"

		if _, existed := fdMap[smallFileName]; !existed {
			smallFile, err := os.OpenFile(smallFileName, os.O_APPEND|os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("[urlTopK] fail to open write file: %v\n", err)
				continue
			} else {
				fdMap[smallFileName] = smallFile
				fileWriterMap[smallFileName] = bufio.NewWriter(smallFile)
			}
		}
		smallFileWriter = fileWriterMap[smallFileName]

		if len(line) > smallFileWriter.Available() {
			smallFileWriter.Flush()
		}

		n, err := smallFileWriter.WriteString(line)
		if err != nil {
			fmt.Printf("[urlTopK] fail to write file: %v\n", err)
		}
		smallFileNameMap[smallFileName] += n
	}

	for _, w := range fileWriterMap {
		w.Flush()
	}

	smallFileNames := make([]string, 0)
	for smallFileName, _ := range smallFileNameMap {
		smallFileNames = append(smallFileNames, smallFileName)
	}
	return smallFileNames
}

func main() {
	file := flag.String("file", "./test.txt", "specify source file")
	k := flag.Int("k", 100, "k value")
	maxCap := flag.Int("cap", 100, "max memory size")

	flag.Parse()

	urlTopK(*file, *k, *maxCap)
}
