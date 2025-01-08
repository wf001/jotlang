package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	mTypes "github.com/wf001/modo/pkg/types"
)

func declarePrintf(
	module *ir.Module,
	libs *mTypes.BuiltinLibProp,
) {
	printfFunc := module.NewFunc(
		"printf",
		types.Void,
		ir.NewParam("format", types.I8Ptr),
	)
	printfFunc.Sig.Variadic = true

	libs.Printf = &mTypes.BuiltinProp{
		FuncPtr: printfFunc,
	}

}

func Declare(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
}
