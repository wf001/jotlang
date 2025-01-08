package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
	mTypes "github.com/wf001/modo/pkg/types"
)

var operatorInsts = map[string]func(*ir.Block, value.Value, value.Value) value.Value{
	// nary
	mTypes.NARY_OPERATOR_ADD: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	},
	mTypes.NARY_OPERATOR_SUB: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewSub(x, y)
	},
	mTypes.NARY_OPERATOR_MUL: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewMul(x, y)
	},
	// binary
	mTypes.BINARY_OPERATOR_EQ: func(block *ir.Block, x, y value.Value) value.Value {
		res := block.NewICmp(enum.IPredEQ, x, y)
		return res
	},
}
