package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultVKubeFileSchema(t *testing.T) {
	s := defaultVKubeFileSchema()
	assert.Equal(t, "object", s.Type)
	assert.ElementsMatch(t, []string{"Kind", "vkubeToken", "containers"}, s.Required)

	kind, ok := s.Properties["Kind"]
	require.True(t, ok)
	assert.Equal(t, "vkube", kind.Fixed)

	_, ok = s.Properties["containers"]
	require.True(t, ok)
	assert.Equal(t, "array", s.Properties["containers"].Type)
}
