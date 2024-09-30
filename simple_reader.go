package main

import (
	"bufio"
	"os"
)

func findUniqueIp(fileName string) (uint64, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	bitmap := createIpsBitMap()
	var uniqueCount uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipStr := scanner.Text()
		wordIndex, bitMask, err := ipStringToBitMap([]byte(ipStr))
		if err != nil {
			return 0, err
		}
		if bitmap[wordIndex]&bitMask == 0 {
			uniqueCount++
			bitmap[wordIndex] |= bitMask
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return uniqueCount, nil
}
