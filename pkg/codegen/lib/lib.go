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
		"format.digit",
		constant.NewCharArrayFromString("%d\x00"),
	)
	globalVars.FormatStr = ir.NewGlobalDef(
		"format.string",
		constant.NewCharArrayFromString("%s\x00"),
	)
	globalVars.FormatSpace = ir.NewGlobalDef(
		"format.space",
		constant.NewCharArrayFromString(" \x00"),
	)
	globalVars.FormatCR = ir.NewGlobalDef(
		"format.cr",
		constant.NewCharArrayFromString("\n\x00"),
	)

	libs.GlobalVar = globalVars

	declareCore(ir, libs)

	log.DebugMessage("built-in library declared")
}
