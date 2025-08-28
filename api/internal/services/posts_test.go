package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "single hashtag",
			text:     "This is a #test",
			expected: []string{"test"},
		},
		{
			name:     "multiple hashtags",
			text:     "Check out #golang and #testing",
			expected: []string{"golang", "testing"},
		},
		{
			name:     "no hashtags",
			text:     "This is just plain text",
			expected: nil,
		},
		{
			name:     "hashtags with special chars",
			text:     "Learning #machine-learning and #deep_learning",
			expected: []string{"machine", "deep_learning"},
		},
		{
			name:     "hashtags with numbers",
			text:     "Version #v2.0 and #test123",
			expected: []string{"v2", "test123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHashtags(tt.text)
			assert.Equal(t, tt.expected, result)
		})
	}
}
