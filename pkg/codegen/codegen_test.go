package codegen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mTypes "github.com/wf001/modo/pkg/types"
)

func TestGetFuncName(t *testing.T) {
	node := &mTypes.Node{Val: "testFunc"}
	want := "fn.testFunc"
	assert.Equal(t, want, node.GetFuncName())
}

func TestGetFuncNameUnnamed(t *testing.T) {
	node := &mTypes.Node{}
	want := "fn.unnamed." + fmt.Sprintf("%p", node)
	assert.Equal(t, want, node.GetUnnamedFuncName())
}

func TestGetBlockName(t *testing.T) {
	node := &mTypes.Node{}
	blockName := "block"
	want := blockName + "." + fmt.Sprintf("%p", node)
	assert.Equal(t, want, node.GetBlockName(blockName))
}
