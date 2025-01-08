package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/codegen/lib/core"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func DeclareBuiltin(ir *ir.Module, libs *mTypes.BuiltinLibProp) {

	declareVariable(ir, libs)
	core.Declare(ir, libs)

	log.DebugMessage("built-in library declared")
}
