package main

import (
	"testing"
)

func TestIpStringToBitMap(t *testing.T) {
	tests := []struct {
		input      string
		wantWord   uint32
		wantMask   uint32
		shouldFail bool
	}{
		// Valid IPs
		{"192.168.0.1", 101007360, 2, false},
		{"0.0.0.0", 0, 1 << 0, false},
		{"255.255.255.255", 134217727, 1 << 31, false},

		// Invalid IPs
		{"invalidIP", 0, 0, true},
		{"300.300.300.300", 0, 0, true},
		{"", 0, 0, true}, // Empty string
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			word, mask, err := ipStringToBitMap([]byte(test.input))
			if test.shouldFail {
				if err == nil {
					t.Errorf("expected error for input %s, got none", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for input %s, got: %v", test.input, err)
				}
				if word != test.wantWord {
					t.Errorf("expected word %d, got %d", test.wantWord, word)
				}
				if mask != test.wantMask {
					t.Errorf("expected mask %d, got %d", test.wantMask, mask)
				}
			}
		})
	}
}
