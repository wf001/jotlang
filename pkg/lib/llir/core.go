package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	modoTypes "github.com/wf001/modo/pkg/types"
)

func genPrintf(module *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	formatStr := module.NewGlobalDef("format", constant.NewCharArray([]byte("%d\n\x00")))
	printfFunc := module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true

	libs.Core = map[string]*modoTypes.CoreProp{}
	libs.Core["prn"] = &modoTypes.CoreProp{
		FuncPtr: printfFunc,
		Args:    []*ir.Global{formatStr},
	}

	return module, libs

}

func GenCore(ir *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	module, libs := genPrintf(ir, &modoTypes.Libs{})
	return module, libs
}
