package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToIntWithDefault(t *testing.T) {
	testCases := []struct {
		name         string
		input        *int
		defaultValue int
		expected     int
	}{
		{
			name:         "Non-nil pointer",
			input:        intPtr(5),
			defaultValue: 10,
			expected:     5,
		},
		{
			name:         "Nil pointer",
			input:        nil,
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toIntWithDefault(tc.input, tc.defaultValue)
			assert.Equal(t, tc.expected, result, "Failed test: "+tc.name)
		})
	}
}

// Helper function to create an int pointer
func intPtr(val int) *int {
	return &val
}
