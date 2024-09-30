package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	filePath := flag.String("file-path", "", "File path")
	mode := flag.String("mode", "multi", "Execution mode: 'simple' for single-threaded or 'multi' for multi-threaded")
	chunkSizeMb := flag.Int("chunk-size-mb", 20, "Size of chunks to process")
	threadsNumber := flag.Int("threads-number", 0, "Number of threads to use")

	flag.Parse()
	if *filePath == "" {
		fmt.Println("Error: File path is required")
		return
	}

	var uniqueCount uint64
	var err error

	if *mode == "simple" {
		uniqueCount, err = findUniqueIp(*filePath)
	} else {
		uniqueCount, err = findUniqueIpMultiThread(*filePath, *chunkSizeMb, *threadsNumber)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Time elapsed took:", time.Since(start))
	fmt.Println("Number of unique IP addresses:", uniqueCount)
}
