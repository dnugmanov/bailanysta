package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "short text",
			text:     "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "long text",
			text:     "This is a very long text that should be truncated",
			maxLen:   20,
			expected: "This is a very long ...",
		},
		{
			name:     "exact length",
			text:     "Exactly 10 chars",
			maxLen:   16,
			expected: "Exactly 10 chars",
		},
		{
			name:     "empty text",
			text:     "",
			maxLen:   10,
			expected: "",
		},
		{
			name:     "unicode text",
			text:     "Привет мир 🌍",
			maxLen:   6,
			expected: "При...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLen)
			assert.Equal(t, tt.expected, result)
			assert.LessOrEqual(t, len(result), tt.maxLen+3) // +3 for "..."
		})
	}
}

func TestNotificationTypes(t *testing.T) {
	tests := []struct {
		name      string
		notifType NotificationType
		expected  string
	}{
		{
			name:      "like notification",
			notifType: NotificationTypeLike,
			expected:  "like",
		},
		{
			name:      "comment notification",
			notifType: NotificationTypeComment,
			expected:  "comment",
		},
		{
			name:      "follow notification",
			notifType: NotificationTypeFollow,
			expected:  "follow",
		},
		{
			name:      "mention notification",
			notifType: NotificationTypeMention,
			expected:  "mention",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.notifType))
		})
	}
}
