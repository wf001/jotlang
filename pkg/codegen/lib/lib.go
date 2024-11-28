package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func DeclareBuiltin(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	globalVars := &mTypes.BuiltinGlobalVarsProp{}
	globalVars.FormatDigit = ir.NewGlobalDef(
		"formatDigit",
		constant.NewCharArray([]byte("%d\n\x00")),
	)
	libs.GlobalVar = globalVars

	declareCore(ir, libs)

	log.DebugMessage("built-in library declared")
}
