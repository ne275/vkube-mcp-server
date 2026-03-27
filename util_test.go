package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetenv(t *testing.T) {
	t.Setenv("VKUBE_MCP_TEST_X", "  trimmed  ")
	assert.Equal(t, "trimmed", getenv("VKUBE_MCP_TEST_X", "def"))
	assert.Equal(t, "def", getenv("VKUBE_MCP_UNSET_XYZ", "def"))
}

func TestDefaultPublicBaseURL(t *testing.T) {
	tests := []struct {
		listen string
		want   string
	}{
		{":3100", "http://127.0.0.1:3100"},
		{"0.0.0.0:9999", "http://127.0.0.1:9999"},
		{"127.0.0.1:8080", "http://127.0.0.1:8080"},
		{"[::]:5000", "http://127.0.0.1:5000"},
		{"bad-address-no-colon", "http://127.0.0.1:3100"},
	}
	for _, tt := range tests {
		t.Run(tt.listen, func(t *testing.T) {
			assert.Equal(t, tt.want, defaultPublicBaseURL(tt.listen))
		})
	}
}
