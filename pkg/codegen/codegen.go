package codegen

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
	"github.com/wf001/modo/util"
)

type assembler struct {
	program *mTypes.Program
}

var operatorMap = map[mTypes.NodeKind]func(*ir.Block, value.Value, value.Value) value.Value{
	// nary
	mTypes.ND_ADD: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	},
	mTypes.ND_SUB: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewSub(x, y)
	},
	mTypes.ND_MUL: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewMul(x, y)
	},
	// binary
	mTypes.ND_EQ: func(block *ir.Block, x, y value.Value) value.Value {
		res := block.NewICmp(enum.IPredEQ, x, y)
		return res
	},
}

func typeOf(v interface{}, ty interface{}) bool {
	return reflect.TypeOf(v) == reflect.TypeOf(ty)
}

func isNumericIR(arg value.Value) bool {
	return typeOf(arg, (*constant.Int)(nil)) ||
		typeOf(arg, (*ir.InstAdd)(nil)) ||
		typeOf(arg, (*ir.InstLoad)(nil)) ||
		typeOf(arg, (*ir.InstICmp)(nil))
}

func isStringIR(arg value.Value) bool {
	return typeOf(arg, (*ir.InstGetElementPtr)(nil))
}

var libraryMap = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, value.Value, *mTypes.Node){
	"prn": func(block *ir.Block, libs *mTypes.BuiltinLibProp, arg value.Value, node *mTypes.Node) {
		var formatStr *ir.Global

		if node.IsType(mTypes.TY_INT32) || node.IsKindNary() || node.IsKindBinary() {
			formatStr = libs.GlobalVar.FormatDigit

		} else if node.IsType(mTypes.TY_STR) {
			formatStr = libs.GlobalVar.FormatStr

		} else {
			log.Panic("unresolved type: have %+v", node)
		}
		block.NewCall(libs.Printf.FuncPtr, formatStr, arg)
	},
}

func newI32(s string) *constant.Int {

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newI32: %s", err)
	}
	return constant.NewInt(types.I32, i)
}

func newArray(length uint64) *types.ArrayType {
	return types.NewArray(length, types.I8)
}

func newStr(block *ir.Block, n *mTypes.Node) *ir.InstGetElementPtr {
	strType := newArray(n.Len)
	strPtr := block.NewAlloca(strType)

	strGEP := block.NewGetElementPtr(strType, strPtr)

	block.NewStore(constant.NewCharArray([]byte(n.Val)), strPtr)

	return strGEP
}

func getFuncName(v *mTypes.Node) string {
	s := v.Val
	if s == "" {
		s = "unnamed"
	}
	return fmt.Sprintf("fn.%s.%p", s, v)
}

func getBlockName(s string, v *mTypes.Node) string {
	return fmt.Sprintf("%s.%p", s, v)
}

func Assemble(llFile string, asmFile string) {
	// TODO: work it?
	out, err, errMsg := util.RunCommand("llc", llFile, "-o", asmFile)
	if err != nil {
		log.Debug("llFile: %s, asmFile: %s", llFile, asmFile)
		log.Panic("fail to asemble: out %+v, err %+v, message %+v", out, err, errMsg)
	}
	log.Debug("written asm: %s", asmFile)
}

// HACK: too many argument
func gen(
	mod *ir.Module,
	block *ir.Block,
	function *ir.Func,
	node *mTypes.Node,
	prog *mTypes.Program,
	scope *mTypes.Node,
) value.Value {

	if node.IsKind(mTypes.ND_SCALAR) && node.IsType(mTypes.TY_INT32) {
		return newI32(node.Val)

	} else if node.IsType(mTypes.TY_STR) {
		return newStr(block, node)

	} else if node.IsKindNary() {
		// nary takes more than 2 arguments
		child := node.Child
		fst := gen(mod, block, function, child, prog, scope)

		child = child.Next
		snd := gen(mod, block, function, child, prog, scope)

		nary := operatorMap[node.Kind]
		res := nary(block, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mod, block, function, child, prog, scope)
			res = nary(block, fst, snd)
		}
		return res

	} else if node.IsKindBinary() {
		// binary takes exactly 2 arguments
		child := node.Child
		fst := gen(mod, block, function, child, prog, scope)

		child = child.Next
		snd := gen(mod, block, function, child, prog, scope)

		binary := operatorMap[node.Kind]
		res := binary(block, fst, snd)
		return res

	} else if node.IsKind(mTypes.ND_LIBCALL) {
		// means calling standard library
		arg := gen(mod, block, function, node.Child, prog, scope)
		libFunc := libraryMap[node.Val]
		libFunc(block, prog.BuiltinLibs, arg, node.Child)
		return newI32("0")

	} else if node.IsKind(mTypes.ND_LAMBDA) {
		// TODO: validate
		funcFn := mod.NewFunc(
			getFuncName(node),
			types.I32,
		)
		llBlock := funcFn.NewBlock(getBlockName("fn.entry", node))
		res := gen(mod, llBlock, funcFn, node.Child, prog, scope)
		if res != nil {
			llBlock.NewRet(res)
		}
		return funcFn

	} else if node.IsKind(mTypes.ND_BIND) {
		for varDeclare := node.Bind; varDeclare != nil; varDeclare = varDeclare.Next {
			v := gen(mod, block, function, varDeclare.Child, prog, scope)
			if scope.Child == nil {
				scope = varDeclare
			} else {
				scope.Next = varDeclare
			}
			if isNumericIR(v) {
				varDeclare.VarPtr = block.NewAlloca(types.I32)
				block.NewStore(v, varDeclare.VarPtr)
			} else if isStringIR(v) {
				varDeclare.VarPtr = v
			}

		}
		return gen(mod, block, function, node.Child, prog, scope)

	} else if node.IsKind(mTypes.ND_IF) {
		condBlock := function.NewBlock(getBlockName("if.cond", node))
		block.NewBr(condBlock)

		thenBlock := function.NewBlock(getBlockName("if.then", node))
		elseBlock := function.NewBlock(getBlockName("if.else", node))
		exitBlock := function.NewBlock(getBlockName("if.exit", node))

		cond := gen(mod, condBlock, function, node.Cond, prog, scope)

		exitBlock.NewRet(newI32("0"))

		thenBlock.NewBr(exitBlock)
		elseBlock.NewBr(exitBlock)

		gen(mod, thenBlock, function, node.Then, prog, scope)
		gen(mod, elseBlock, function, node.Else, prog, scope)

		condBlock.NewCondBr(cond, thenBlock, elseBlock)

	} else if node.IsKind(mTypes.ND_EXPR) {
		var res value.Value
		for child := node.Child; child != nil; child = child.Next {
			res = gen(mod, block, function, child, prog, scope)
		}
		return res

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) && node.Val == "main" {
		// means declaring main function regarded as entrypoint

		fnc := mod.NewFunc(
			"main",
			types.I32,
		)
		llBlock := fnc.NewBlock("")

		res := gen(mod, llBlock, fnc, node.Child, prog, scope)
		llBlock.NewCall(res)
		llBlock.NewRet(newI32("0"))

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) {
		// means declaring global variable or function named except main

		retType := types.I32 // TODO: to be changable
		funcName := getFuncName(node)

		fnc := mod.NewFunc(
			funcName,
			retType,
		)
		llBlock := fnc.NewBlock("")

		res := gen(mod, llBlock, fnc, node.Child, prog, scope)
		node.FuncPtr = fnc
		llBlock.NewRet(res)

	} else if node.IsKind(mTypes.ND_VAR_REFERENCE) {
		// PERFORMANCE: too redundant
		// find in local variable which is declared with let
		for s := scope; s != nil; s = s.Next {
			if s.Val == node.Val {
				if s.Child.IsType(mTypes.TY_INT32) {
					node.Type = mTypes.TY_INT32
					return block.NewLoad(types.I32, s.VarPtr)

				} else if s.Child.IsType(mTypes.TY_STR) {
					node.Type = mTypes.TY_STR
					return block.NewGetElementPtr(newArray(s.Child.Len), s.VarPtr)

				} else {
					log.Panic("unresolved variable: have %+v, %+s", s)
				}
			}
		}

		// find in global variable which is declared with def
		for declare := prog.Declares; declare != nil; declare = declare.Next {
			if declare.Child.Val == node.Val {
				return block.NewCall(declare.Child.FuncPtr)
			}
		}

		log.Panic("unresolved symbol: '%s'", node.Val)

	} else if node.IsKind(mTypes.ND_DECLARE) {
		return gen(mod, block, function, node.Child, prog, scope)

	} else {
		log.Panic("unresolved Nodekind: have %+v", node)
	}
	return nil
}

func constructModule(prog *mTypes.Program) *ir.Module {
	module := ir.NewModule()
	prog.BuiltinLibs = &mTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, prog.BuiltinLibs)

	for declare := prog.Declares; declare != nil; declare = declare.Next {
		gen(module, nil, nil, declare, prog, &mTypes.Node{})
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) GenFrontend(llName string, asmName string) {
	log.DebugMessage("ir module constructing")
	module := constructModule(a.program)
	log.DebugMessage("ir module constructed")
	log.Debug("[IR]\n%s\n", module.String())

	err := os.WriteFile(llName, []byte(module.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

}
