package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockName(t *testing.T) {
	assert.Equal(t, true, EqualType("str", "str"))
	assert.Equal(t, false, EqualType(1, "str"))
}
