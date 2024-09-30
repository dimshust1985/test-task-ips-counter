package main

import (
	"os"
	"testing"
)

func TestFindUniqueIpMultiThread(t *testing.T) {
	// Prepare test cases
	tests := []struct {
		name          string
		fileContent   string
		chunkSizeMb   int
		threadsNumber int
		expectedCount uint64
		expectError   bool
	}{
		{
			name:          "Multiple unique IPs, multiple threads",
			fileContent:   "192.168.0.1\n192.168.0.2\n192.168.0.3\n192.168.0.4\n", // 4 unique IPs
			chunkSizeMb:   1,                                                      // Small chunk size to test multi-threading
			threadsNumber: 2,                                                      // Use 2 threads
			expectedCount: 4,                                                      // Expect 4 unique IPs
			expectError:   false,                                                  // No errors expected
		},
		{
			name:          "Duplicate IPs, multiple threads",
			fileContent:   "192.168.0.1\n192.168.0.1\n192.168.0.2\n192.168.0.2\n", // 2 unique IPs with duplicates
			chunkSizeMb:   1,                                                      // Small chunk size
			threadsNumber: 2,                                                      // Use 2 threads
			expectedCount: 2,                                                      // Expect 2 unique IPs
			expectError:   false,                                                  // No errors expected
		},
		{
			name:          "File with no newlines means invalid",
			fileContent:   "192.168.0.1192.168.0.2", // IPs without newline between them (invalid)
			chunkSizeMb:   1,                        // Small chunk size
			threadsNumber: 2,                        // Use 2 threads
			expectedCount: 0,                        // Only the first valid IP should be counted
			expectError:   false,                    // No errors expected
		},
		{
			name:          "Empty file",
			fileContent:   "",    // No content
			chunkSizeMb:   1,     // Chunk size is irrelevant here
			threadsNumber: 2,     // Any number of threads
			expectedCount: 0,     // No IPs
			expectError:   false, // No errors expected
		},
		{
			name:          "Single thread processing",
			fileContent:   "192.168.0.1\n192.168.0.2\n192.168.0.3\n", // Multiple unique IPs
			chunkSizeMb:   1,                                         // Small chunk size
			threadsNumber: 1,                                         // Use single thread
			expectedCount: 3,                                         // Expect 3 unique IPs
			expectError:   false,                                     // No errors expected
		},
		{
			name:          "Invalid IP format",
			fileContent:   "invalidIP\n192.168.0.1\n", // Invalid IP followed by valid one
			chunkSizeMb:   1,                          // Small chunk size
			threadsNumber: 2,                          // Use 2 threads
			expectedCount: 1,                          // Only 1 valid IP
			expectError:   false,                      // No errors expected
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a temporary file with the given content
			tmpFile, err := os.CreateTemp("", "test_ip_multi_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up after test

			// Write the test content to the temporary file
			if _, err := tmpFile.WriteString(test.fileContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Call the findUniqueIpMultiThread function
			uniqueCount, err := findUniqueIpMultiThread(tmpFile.Name(), test.chunkSizeMb, test.threadsNumber)

			// Check for errors if we expect one
			if test.expectError {
				if err == nil {
					t.Fatalf("Expected an error but got none")
				}
			} else {
				// Check for unexpected errors
				if err != nil {
					t.Fatalf("Did not expect an error, but got: %v", err)
				}

				// Verify that the count of unique IPs matches the expected result
				if uniqueCount != test.expectedCount {
					t.Errorf("Expected uniqueCount %d, got %d", test.expectedCount, uniqueCount)
				}
			}
		})
	}
}
