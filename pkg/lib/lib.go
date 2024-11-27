package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/wf001/modo/pkg/lib/llir"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declareGlobalFormatDigit(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {
	formatStr := module.NewGlobalDef("formatDigit", constant.NewCharArray([]byte("%d\n\x00")))

	libs.GlobalVar = &mTypes.BuiltinGlobalVarsProp{}
	libs.GlobalVar.FormatDigit = formatStr
}

func DeclareBuiltin(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declareGlobalFormatDigit(ir, libs)
	core.Declare(ir, libs)
	log.DebugMessage("built-in library declared")
}
