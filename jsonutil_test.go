package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	out := toJSON(map[string]int{"a": 1})
	assert.Contains(t, out, `"a"`)
	assert.Contains(t, out, `1`)
}
