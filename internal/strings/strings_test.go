package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_Remove(t *testing.T) {
	tt := []struct {
		name     string
		s        String
		i        int
		expected String
		removed  bool
	}{
		{
			name:     "i ∈ [0:len(s)-1]",
			s:        New("tesit"),
			i:        3,
			expected: New("test"),
			removed:  true,
		},
		{
			name:     "i ∉ [0:len(s)-1]",
			s:        New("test"),
			i:        10,
			expected: New("test"),
			removed:  false,
		},
		{
			name:     "i < 0",
			s:        New("test"),
			i:        -1,
			expected: New("test"),
			removed:  false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, removed := tc.s.Remove(tc.i)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.removed, removed)
		})
	}
}
