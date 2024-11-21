package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	modoTypes "github.com/wf001/modo/pkg/types"
)

func genPrintf(module *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	stdLibFunc := module.NewFunc(
		"core.printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	stdLibFunc.Sig.Variadic = true

	entry := stdLibFunc.NewBlock("entry")

	// フォーマット文字列 "%f\n" を定義
	formatStr := entry.NewAlloca(types.NewArray(4, types.I8)) // "%f\n" + null terminator
	entry.NewStore(constant.NewCharArray([]byte("%s\n\x00")), formatStr)

	entry.NewRet(constant.NewInt(types.I32, 0))

	libs.Core["printf"] = &modoTypes.CoreProp{
		FuncPtr: stdLibFunc,
		Args:    []*ir.InstAlloca{formatStr},
	}

	return module, libs

}

func GenCore(ir *ir.Module, libs *modoTypes.Libs) (*ir.Module, *modoTypes.Libs) {
	module, libs := genPrintf(ir, &modoTypes.Libs{})
	return module, libs
}
