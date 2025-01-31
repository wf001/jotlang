package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, *mTypes.Node) value.Value{
	mTypes.LIB_CORE_PRN: InvokePrn,
	// nary
	mTypes.OPERATOR_ADD: InvokeAdd,
	mTypes.OPERATOR_SUB: InvokeSub,
	mTypes.OPERATOR_MUL: InvokeMul,
	mTypes.OPERATOR_DIV: InvokeDiv,
	mTypes.OPERATOR_EQ:  InvokeEq,
	mTypes.OPERATOR_GT:  InvokeGt,
	mTypes.OPERATOR_LT:  InvokeLt,
	// binary
	mTypes.OPERATOR_MOD: InvokeMod,
}
