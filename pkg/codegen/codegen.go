package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func Codegen(s string) *ir.Module {
	m := ir.NewModule()

	globalG := m.NewGlobalDef("g", constant.NewInt(types.I32, 58))

	funcAdd := m.NewFunc("add", types.I32,
		ir.NewParam("x", types.I32),
		ir.NewParam("y", types.I32),
	)
	ab := funcAdd.NewBlock("")
	ab.NewRet(ab.NewAdd(funcAdd.Params[0], funcAdd.Params[1]))

	funcMain := m.NewFunc(
		"main",
		types.I32,
	) // omit parameters
	mb := funcMain.NewBlock("") // llir/llvm would give correct default name for block without name
	mb.NewRet(mb.NewCall(funcAdd, constant.NewInt(types.I32, 59), mb.NewLoad(types.I32, globalG)))
	return m

}
