package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToInt64(t *testing.T) {
	tests := []struct {
		input       string
		expected    uint64
		expectedErr error
	}{
		{"123", 123, nil},
		{"0", 0, nil},
		{"999999999999", 999999999999, nil}, // example of a large number
		{"abc", 0, ErrInvalidID},            // example of a non-integer string
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := StringToUint64(tc.input)
			if err != nil {
				// t.Errorf("unexpected error: %v", err.Error())
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}
