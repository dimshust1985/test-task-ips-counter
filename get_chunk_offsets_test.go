package main

import (
	"os"
	"testing"
)

func TestGetChunkOffsets(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string
		chunkSize     int64
		expected      []int64
		expectError   bool
		errorContains string
	}{
		{
			name:        "Small file, no newlines",
			fileContent: "192.168.0.1",
			chunkSize:   1024,
			expected:    []int64{0, 11},
			expectError: false,
		},
		{
			name:        "Large file, multiple lines",
			fileContent: "192.168.0.1\n192.168.0.2\n192.168.0.3\n",
			chunkSize:   10,
			expected:    []int64{0, 12, 24, 36},
			expectError: false,
		},
		{
			name:          "File with long line (trigger large chunk error)",
			fileContent:   "192.168.0.1 192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5",
			chunkSize:     10,
			expected:      nil,
			expectError:   true,
			errorContains: "Cant calculate offset at start offset:",
		},
		{
			name:        "File ends with newline",
			fileContent: "192.168.0.1\n192.168.0.2\n",
			chunkSize:   15,
			expected:    []int64{0, 24},
			expectError: false,
		},
		{
			name:        "Empty file",
			fileContent: "",
			chunkSize:   1024,
			expected:    []int64{0},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_offsets_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up after test

			if _, err := tmpFile.WriteString(test.fileContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			offsets, err := getChunkOffsets(tmpFile.Name(), test.chunkSize)

			if test.expectError {
				if err == nil {
					t.Fatalf("Expected an error but got none")
				}
				if !containsErrorSubstring(err.Error(), test.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", test.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect an error, but got: %v", err)
				}

				if len(offsets) != len(test.expected) {
					t.Fatalf("Expected %d offsets, got %d", len(test.expected), len(offsets))
				}

				for i, expectedOffset := range test.expected {
					if offsets[i] != expectedOffset {
						t.Errorf("Expected offset %d at index %d, got %d", expectedOffset, i, offsets[i])
					}
				}
			}
		})
	}
}

func containsErrorSubstring(errMsg, expectedSubstring string) bool {
	return len(errMsg) >= len(expectedSubstring) && errMsg[:len(expectedSubstring)] == expectedSubstring
}
