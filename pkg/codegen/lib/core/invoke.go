package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	mTypes "github.com/wf001/modo/pkg/types"
)

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, *mTypes.Node) value.Value{
	"prn": InvokePrn,
}
