package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, *mTypes.Node) value.Value{
	"prn": InvokePrn,
	// nary
	mTypes.NARY_OPERATOR_ADD: InvokeAdd,
	mTypes.NARY_OPERATOR_SUB: InvokeSub,
	mTypes.NARY_OPERATOR_MUL: InvokeMul,

	// binary
	mTypes.BINARY_OPERATOR_EQ: InvokeEq,
}
