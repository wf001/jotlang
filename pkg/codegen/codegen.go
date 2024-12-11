package codegen

import (
	"fmt"
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
	"github.com/wf001/modo/pkg/io"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
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
		return block.NewZExt(res, types.I32)
	},
}
var libraryMap = map[string]func(*ir.Block, *mTypes.BuiltinLibProp, value.Value){
	"prn": func(block *ir.Block, libs *mTypes.BuiltinLibProp, arg value.Value) {
		block.NewCall(libs.Printf.FuncPtr, libs.GlobalVar.FormatDigit, arg)
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

func getFuncName(v string) string {
	return fmt.Sprintf("fn-%s", v)
}

func doAsemble(llFile string, asmFile string) {
	// TODO: work it?
	out, err, errMsg := io.RunCommand("llc", llFile, "-o", asmFile)
	if err != nil {
		log.Debug("llFile: %s, asmFile: %s", llFile, asmFile)
		// TODO: move to utility
		log.Panic(
			"fail to asemble: %+v",
			map[string]interface{}{
				"out":    out,
				"err":    err,
				"errMsg": errMsg,
			},
		)
	}
	log.Debug("written asm: %s", asmFile)
}

func gen(
	mod *ir.Module,
	block *ir.Block,
	node *mTypes.Node,
	prog *mTypes.Program,
	scope *mTypes.Node,
) value.Value {

	if node.IsInteger() {
		return newI32(node.Val)

	} else if node.IsNary() {
		// nary takes more than 2 arguments
		child := node.Child
		fst := gen(mod, block, child, prog, scope)

		child = child.Next
		snd := gen(mod, block, child, prog, scope)

		nary := operatorMap[node.Kind]
		res := nary(block, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mod, block, child, prog, scope)
			res = nary(block, fst, snd)
		}
		return res

	} else if node.IsBinary() {
		// binary takes exactly 2 arguments
		child := node.Child
		fst := gen(mod, block, child, prog, scope)

		child = child.Next
		snd := gen(mod, block, child, prog, scope)

		binary := operatorMap[node.Kind]
		res := binary(block, fst, snd)
		return res

	} else if node.IsLibCall() {
		// means calling standard library
		arg := gen(mod, block, node.Child, prog, scope)
		libFunc := libraryMap[node.Val]
		libFunc(block, prog.BuiltinLibs, arg)
		return newI32("0")

	} else if node.IsLambda() {
		// TODO: validate
		funcFn := mod.NewFunc(
			fmt.Sprintf("fn-%p", node),
			types.I32,
		)
		llBlock := funcFn.NewBlock("")
		res := gen(mod, llBlock, node.Child, prog, scope)
		llBlock.NewRet(res)
		return funcFn

	} else if node.IsBind() {
		for varDeclare := node.Bind; varDeclare != nil; varDeclare = varDeclare.Next {
			v := gen(mod, block, varDeclare.Child, prog, scope)
			if scope.Child == nil {
				scope = varDeclare
			} else {
				scope.Next = varDeclare
			}
			varDeclare.VarPtr = block.NewAlloca(types.I32)
			block.NewStore(v, varDeclare.VarPtr)
		}
		return gen(mod, block, node.Child, prog, scope)

	} else if node.IsExpr() {
		var res value.Value
		for child := node.Child; child != nil; child = child.Next {
			res = gen(mod, block, child, prog, scope)
		}
		return res

	} else if node.IsVarDeclare() && node.Val == "main" {
		// means declaring main function regarded as entrypoint

		function := mod.NewFunc(
			"main",
			types.I32,
		)
		llBlock := function.NewBlock("")

		res := gen(mod, llBlock, node.Child, prog, scope)
		llBlock.NewCall(res)
		llBlock.NewRet(newI32("0"))

	} else if node.IsVarDeclare() {
		// means declaring global variable or function named except main

		retType := types.I32 // TODO: to be changable
		funcName := getFuncName(node.Val)

		function := mod.NewFunc(
			funcName,
			retType,
		)
		llBlock := function.NewBlock("")

		res := gen(mod, llBlock, node.Child, prog, scope)
		node.FuncPtr = function
		llBlock.NewRet(res)

	} else if node.IsVarReference() {
		// PERFORMANCE: too redundant
		for s := scope; s != nil; s = s.Next {
			if s.Val == node.Val {
				return block.NewLoad(types.I32, s.VarPtr)
			}
		}

		for declare := prog.Declares; declare != nil; declare = declare.Next {
			if declare.Child.Val == node.Val {
				return block.NewCall(declare.Child.FuncPtr)
			}
		}

		log.Panic("unresolved symbol: '%s'", node.Val)

	} else if node.IsDeclare() {
		return gen(mod, block, node.Child, prog, scope)

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
		gen(module, nil, declare, prog, &mTypes.Node{})
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	log.DebugMessage("ir module constructing")
	module := constructModule(a.program)
	log.DebugMessage("ir module constructed")
	log.Debug("[IR]\n%s\n", module.String())

	err := os.WriteFile(llName, []byte(module.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
