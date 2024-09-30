package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	StartOffset int64
	EndOffset   int64
	current     int
	total       int
}

const maxChunkSizeMB = 200

func findUniqueIpMultiThread(fileName string, chunkSizeMB int, threadNumbers int) (uint64, error) {

	if chunkSizeMB > maxChunkSizeMB {
		fmt.Printf("chunkSizeMB > maxChunkSizeMB (%d), chunkSizeMB is set to maxChunkSizeMB\n", chunkSizeMB)
		chunkSizeMB = maxChunkSizeMB
	}

	chunkSizeBytes := int64(1024 * 1024 * chunkSizeMB)

	offsets, err := getChunkOffsets(fileName, chunkSizeBytes)
	if err != nil {
		return 0, err
	}

	bitmap := createIpsBitMap()

	var uniqueCount uint64
	jobChan := make(chan Job, len(offsets)-1)

	var numWorkers int
	if threadNumbers > 0 {
		numWorkers = threadNumbers
	} else if threadNumbers > runtime.NumCPU() {
		fmt.Printf("number of threads is more than number of CPU: %d\n", runtime.NumCPU())
		numWorkers = runtime.NumCPU()
	} else {
		fmt.Printf("default number of threads is used: %d\n", runtime.NumCPU())
		numWorkers = runtime.NumCPU()
	}
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobChan {
				start := time.Now()
				err := processChunk(fileName, job.StartOffset, job.EndOffset, bitmap, &uniqueCount)
				if err != nil {
					fmt.Printf("error proccessing chunk N: %d, %s", job.current, err.Error())
				} else {
					fmt.Printf("Chunk N: %d proccessd, time: %s, total chunks: %d \n", job.current, time.Since(start), job.total)
				}
			}
		}()
	}

	total := len(offsets) - 1
	for i := 0; i < len(offsets)-1; i++ {
		jobChan <- Job{
			StartOffset: offsets[i],
			EndOffset:   offsets[i+1],
			current:     i + 1,
			total:       total,
		}
	}
	close(jobChan)
	wg.Wait()
	return uniqueCount, nil
}

// precalculate all offsets to feed threads
func getChunkOffsets(filePath string, chunkSize int64) ([]int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	var offsets []int64
	startOffset := int64(0)
	offsets = append(offsets, startOffset)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	for startOffset < fileSize {
		tempOffset := startOffset + chunkSize
		if tempOffset >= fileSize {
			tempOffset = fileSize
		} else {
			tempOffset, err = adjustOffsetToNextNewline(f, tempOffset)
			if tempOffset-startOffset > chunkSize*2 {
				return nil, errors.New("Cant calculate offset at start offset: " + strconv.FormatInt(startOffset, 10))
			}
			if err != nil {
				return nil, err
			}
		}
		offsets = append(offsets, tempOffset)
		startOffset = tempOffset
	}

	return offsets, nil
}

func adjustOffsetToNextNewline(f *os.File, offset int64) (int64, error) {
	_, err := f.Seek(offset, io.SeekStart)
	if err != nil {
		return offset, err
	}
	reader := bufio.NewReader(f)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				newOffset, _ := f.Seek(0, io.SeekCurrent)
				return newOffset, nil
			}
			return offset, err
		}
		if b == '\n' {
			break
		}
	}
	currentFilePosition, _ := f.Seek(0, io.SeekCurrent)
	// current position moved to the buffer size, and buffered contains (buffer size - bytes red)
	newOffset := currentFilePosition - int64(reader.Buffered())
	return newOffset, nil
}

func processChunk(filePath string, startOffset, endOffset int64, bitmap []uint32, uniqueCount *uint64) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := io.NewSectionReader(f, startOffset, endOffset-startOffset)
	bufReader := bufio.NewReaderSize(reader, int(endOffset-startOffset))

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			processLine([]byte(line), bitmap, uniqueCount)
		}
		if err == io.EOF {
			break
		}
	}
	return nil

}

func processLine(line []byte, bitmap []uint32, uniqueCount *uint64) {
	wordIndex, bitMask, err := ipStringToBitMap(line)
	if err != nil {
		fmt.Printf("Error parsing line: %v, the line will be skipped\n", err)
		return
	}
	addr := &bitmap[wordIndex]
	oldValue := atomic.LoadUint32(addr)
	if oldValue&bitMask != 0 {
		return
	}
	for {
		oldValue = atomic.LoadUint32(addr)
		if oldValue&bitMask != 0 {
			break
		}
		newValue := oldValue | bitMask
		if atomic.CompareAndSwapUint32(addr, oldValue, newValue) {
			atomic.AddUint64(uniqueCount, 1)
			break
		}
	}
}
