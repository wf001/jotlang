package core

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func InvokePrn(block *ir.Block, libs *mTypes.BuiltinLibProp, node *mTypes.Node) value.Value {
	var formatStr *ir.Global

	for n := node; n != nil; n = n.Next {
		v := n.IRValue
		if n.IsType(mTypes.TY_INT32) || n.IRValue.Type().Equal(types.I32) {
			formatStr = libs.GlobalVar.FormatDigit

		} else if n.IsType(mTypes.TY_BOOL) || n.IRValue.Type().Equal(types.I1) {
			formatStr = libs.GlobalVar.FormatStr
			v = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)

		} else if n.IsType(mTypes.TY_STR) || n.IRValue.Type().Equal(types.I8Ptr) {
			formatStr = libs.GlobalVar.FormatStr

		} else if n.IsKind(mTypes.ND_FUNCCALL) {
			t := node.IRValue.Type()
			if t.Equal(types.I32) {
				formatStr = libs.GlobalVar.FormatDigit
			} else if t.Equal(types.I1) {
				formatStr = libs.GlobalVar.FormatStr
				v = block.NewSelect(n.IRValue, libs.GlobalVar.TrueValue, libs.GlobalVar.FalseValue)
			} else if t.Equal(types.Void) {
				formatStr = libs.GlobalVar.FormatStr
			} else if _, ok := t.(*types.PointerType); ok {
				formatStr = libs.GlobalVar.FormatStr
			} else {
				log.Panic("unresolved type: have %+v", n)
			}
		} else {
			log.Panic("unresolved type: have %+v", n)
		}
		block.NewCall(libs.Printf.FuncPtr, formatStr, v)

		if n.Next == nil {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatCR)
		} else {
			block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatSpace)
		}

	}
	return nil
}
