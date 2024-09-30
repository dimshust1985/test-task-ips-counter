package main

import (
	"os"
	"testing"
)

func TestProcessChunk(t *testing.T) {
	// Prepare test cases
	tests := []struct {
		name           string
		chunkContent   string
		expectedCount  uint64
		expectedBitmap []uint32
	}{
		{
			name:          "Valid IPs",
			chunkContent:  "192.168.0.1\n192.168.0.2\n192.168.0.1\n", // Repeated IPs
			expectedCount: 2,                                         // Only 2 unique IPs
			expectedBitmap: func() []uint32 {
				bitmap := createIpsBitMap()
				// Manually set bits for 192.168.0.1 and 192.168.0.2
				wordIndex1, bitMask1, _ := ipStringToBitMap([]byte("192.168.0.1"))
				wordIndex2, bitMask2, _ := ipStringToBitMap([]byte("192.168.0.2"))
				bitmap[wordIndex1] |= bitMask1
				bitmap[wordIndex2] |= bitMask2
				return bitmap
			}(),
		},
		{
			name:           "Empty Chunk",
			chunkContent:   "",
			expectedCount:  0,
			expectedBitmap: createIpsBitMap(),
		},
		{
			name:          "Invalid IPs",
			chunkContent:  "invalidIP\n192.168.0.1\n", // One valid, one invalid
			expectedCount: 1,
			expectedBitmap: func() []uint32 {
				bitmap := createIpsBitMap()
				// Manually set bit for 192.168.0.1
				wordIndex1, bitMask1, _ := ipStringToBitMap([]byte("192.168.0.1"))
				bitmap[wordIndex1] |= bitMask1
				return bitmap
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a temp file for chunk content using os.TempFile
			tmpFile, err := os.CreateTemp("", "test_chunk_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up after test

			// Write chunk content to file
			if _, err := tmpFile.WriteString(test.chunkContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Prepare bitmap and unique count
			bitmap := createIpsBitMap()
			var uniqueCount uint64

			// Run processChunk on the temp file (simulating a chunk)
			err = processChunk(tmpFile.Name(), 0, int64(len(test.chunkContent)), bitmap, &uniqueCount)
			if err != nil {
				t.Fatalf("Error processing chunk: %v", err)
			}

			// Verify the unique count
			if uniqueCount != test.expectedCount {
				t.Errorf("Expected uniqueCount %d, got %d", test.expectedCount, uniqueCount)
			}

			// Verify the bitmap
			for i, val := range bitmap {
				if val != test.expectedBitmap[i] {
					t.Errorf("Bitmap mismatch at index %d: expected %d, got %d", i, test.expectedBitmap[i], val)
				}
			}
		})
	}
}
