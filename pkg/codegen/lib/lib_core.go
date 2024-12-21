package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func declarePrintf(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {
	printfFunc := module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true

	libs.Printf = &mTypes.BuiltinProp{
		FuncPtr: printfFunc,
	}

}

func declareCore(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
}

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, value.Value, *mTypes.Node) value.Value{
	"prn": func(block *ir.Block, libs *mTypes.BuiltinLibProp, arg value.Value, node *mTypes.Node) value.Value {
		var formatStr *ir.Global

		if node.IsType(mTypes.TY_INT32) || node.IsKindNary() || node.IsKindBinary() {
			formatStr = libs.GlobalVar.FormatDigit

		} else if node.IsType(mTypes.TY_STR) {
			formatStr = libs.GlobalVar.FormatStr

		} else {
			log.Panic("unresolved type: have %+v", node)
		}
		block.NewCall(libs.Printf.FuncPtr, formatStr, arg)
		// todo: return nil
		return constant.NewInt(types.I32, 0)
	},
}
