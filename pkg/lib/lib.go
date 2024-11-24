package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/lib/llir"
	"github.com/wf001/modo/pkg/types"
)

func Gen(ir *ir.Module) (*ir.Module, *types.Libs) {
	libs := &types.Libs{}
	libs.Core = map[string]*types.CoreProp{}
	libs.GlobalVar = map[string]*types.CoreGlobalVars{}
	return core.GenCore(ir, libs)
}
