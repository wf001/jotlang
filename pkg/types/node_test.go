package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFuncName(t *testing.T) {
	node := &Node{Val: "testFunc"}
	want := "fn.testFunc"
	assert.Equal(t, want, node.GetFuncName())
}

func TestGetFuncNameUnnamed(t *testing.T) {
	node := &Node{}
	want := "fn.unnamed." + fmt.Sprintf("%p", node)
	assert.Equal(t, want, node.GetUnnamedFuncName())
}

func TestGetBlockName(t *testing.T) {
	node := &Node{}
	blockName := "block"
	want := blockName + "." + fmt.Sprintf("%p", node)
	assert.Equal(t, want, node.GetBlockName(blockName))
}
