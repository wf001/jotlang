package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	modoTypes "github.com/wf001/modo/pkg/types"
)

func genPrintf(module *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	printfFunc := module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true

	libs.Core["prn"] = &modoTypes.CoreProp{
		FuncPtr: printfFunc,
	}

	return module, libs

}
func genGlobal(module *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	formatStr := module.NewGlobalDef("formatDigit", constant.NewCharArray([]byte("%d\n\x00")))
	libs.GlobalVar["formatDigit"] = &modoTypes.CoreGlobalVars{
		Vars: formatStr,
	}
	return module, libs
}

func GenCore(ir *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	module, libs := genGlobal(ir, libs)
	module, libs = genPrintf(ir, libs)
	return module, libs
}
