package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeString(t *testing.T) {
	assert.Equal(t, NormalizeString("abC dEF"), "abc def")
	assert.Equal(t, NormalizeString("žůžo"), "zuzo")
	assert.Equal(t, NormalizeString("Loïc"), "loic")
}
