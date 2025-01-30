package core

import (
	"github.com/llir/llvm/ir/value"
	mTypes "github.com/wf001/modo/pkg/types"
)

func invokeFold(
	node *mTypes.Node,
	operate func(value.Value, value.Value) value.Value,
) value.Value {

	fst := node.IRValue

	node = node.Next
	snd := node.IRValue

	res := operate(fst, snd)

	for node = node.Next; node != nil; node = node.Next {
		fst := res
		snd := node.IRValue
		res = operate(fst, snd)
	}
	return res
}
