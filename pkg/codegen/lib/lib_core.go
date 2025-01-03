package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
	"github.com/wf001/modo/util"
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

func declareCore(ir *ir.Module, libs *mTypes.BuiltinLibProp) {
	declarePrintf(ir, libs)
}

var LibInsts = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, *mTypes.Node) value.Value{
	"prn": func(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
		var formatStr *ir.Global

		for n := node; n != nil; n = n.Next {
			if n.IsType(mTypes.TY_INT32) || n.IsKindNary() || n.IsKindBinary() ||
				util.EqualType(n.IRValue.Type(), types.I32) {
				formatStr = libs.GlobalVar.FormatDigit

			} else if n.IsType(mTypes.TY_STR) || util.EqualType(n.IRValue.Type(), types.I8Ptr) {
				formatStr = libs.GlobalVar.FormatStr

			} else if n.IsKind(mTypes.ND_FUNCCALL) {
				t := node.IRValue.Type()
				if util.EqualType(t, types.I32) {
					formatStr = libs.GlobalVar.FormatDigit
				} else if _, ok := t.(*types.PointerType); ok {
					formatStr = libs.GlobalVar.FormatStr
				} else {
					log.Panic("unresolved type: have %+v", n)
				}
			} else {
				log.Panic("unresolved type: have %+v", n)
			}
			block.NewCall(libs.Printf.FuncPtr, formatStr, n.IRValue)

			if n.Next == nil {
				block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatCR)
			} else {
				block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
			}

		}
		// todo: return nil
		return nil
	},
}
