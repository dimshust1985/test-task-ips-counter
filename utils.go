package main

import (
	"encoding/binary"
	"errors"
	"net"
)

func ipStringToBitMap(line []byte) (uint32, uint32, error) {
	ipStr := string(line)
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, 0, errors.New("Error parsing IP address: " + ipStr)
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, 0, errors.New("Error parsing IP address: " + ipStr)
	}
	ipUint32 := binary.BigEndian.Uint32(ip4)
	wordIndex := ipUint32 / 32
	bitIndex := ipUint32 % 32
	bitMask := uint32(1) << bitIndex
	return wordIndex, bitMask, nil
}

func createIpsBitMap() []uint32 {
	const totalBits = 1 << 32
	const bitsPerWord = 32
	numWords := totalBits / bitsPerWord
	return make([]uint32, numWords)
}
