package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRHostSplit(t *testing.T) {
	s, err := rHostSplit("[::]:6379:1:abc")
	assert.NoError(t, err)
	assert.Equal(t, "::", s.host)
	assert.Equal(t, 6379, s.port)
	assert.Equal(t, 1, s.db)
	assert.Equal(t, "abc", s.pass)

	s, err = rHostSplit("127.0.0.1:6379:1:abc")
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", s.host)
	assert.Equal(t, 6379, s.port)
	assert.Equal(t, 1, s.db)
	assert.Equal(t, "abc", s.pass)
}
