package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokeAdd(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return InvokeFold(node, func(x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	})
}

func InvokeSub(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return InvokeFold(node, func(x, y value.Value) value.Value {
		return block.NewSub(x, y)
	})
}

func InvokeMul(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return InvokeFold(node, func(x, y value.Value) value.Value {
		return block.NewMul(x, y)
	})
}

func InvokeEq(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	return InvokeFold(node, func(x, y value.Value) value.Value {
		return block.NewICmp(enum.IPredEQ, x, y)
	})
}

func InvokeFold(
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
