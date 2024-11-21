package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/lib/llir"
	"github.com/wf001/modo/pkg/types"
)

func Gen(ir *ir.Module) (*ir.Module, *types.Libs) {
	return core.GenCore(ir, &types.Libs{})
}
