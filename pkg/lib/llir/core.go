package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	modoTypes "github.com/wf001/modo/pkg/types"
)

func declarePrintf(
	module *ir.Module,
	libs *modoTypes.BuiltinLibProp,
) {
	printfFunc := module.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("format", types.NewPointer(types.I8)),
	)
	printfFunc.Sig.Variadic = true

	libs.Printf = &modoTypes.BuiltinProp{}
	libs.Printf.FuncPtr = printfFunc

}

func Declare(ir *ir.Module, libs *modoTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
}
