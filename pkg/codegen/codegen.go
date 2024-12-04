package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
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

func doAsemble(llFile string, asmFile string) {
	// TODO: show error message send to stderr
	out, err := exec.Command("llc", llFile, "-o", asmFile).Output()
	if err != nil {
		log.Panic(
			"fail to asemble: %+v",
			map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile},
		)
	}
	log.Debug("written asm: %s", asmFile)
}

func gen(
	mod *ir.Module,
	bl *ir.Block,
	funcCallNode *mTypes.Node,
	libs *mTypes.BuiltinLibProp,
) value.Value {

	if funcCallNode.IsInteger() {
		return newI32(funcCallNode.Val)

	} else if funcCallNode.IsNary() {
		// nary takes more than 2 arguments
		child := funcCallNode.Child
		fst := gen(mod, bl, child, libs)

		child = child.Next
		snd := gen(mod, bl, child, libs)

		nary := operatorMap[funcCallNode.Kind]
		res := nary(bl, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mod, bl, child, libs)
			res = nary(bl, fst, snd)
		}
		return res

	} else if funcCallNode.IsBinary() {
		// binary takes exactly 2 arguments
		child := funcCallNode.Child
		fst := gen(mod, bl, child, libs)

		child = child.Next
		snd := gen(mod, bl, child, libs)

		binary := operatorMap[funcCallNode.Kind]
		res := binary(bl, fst, snd)

		return res
	} else if funcCallNode.IsLibrary() {
		// means calling standard library
		arg := gen(mod, bl, funcCallNode.Child, libs)
		libFunc := libraryMap[funcCallNode.Val]
		libFunc(bl, libs, arg)

		return newI32("0")
	} else if funcCallNode.IsVar() {
		// means declaring global variable or function
		child := gen(mod, bl, funcCallNode.Child, libs)
		variable := bl.NewAlloca(types.I32)

		if funcCallNode.Child.Kind == mTypes.ND_DECLARE {
			bl.NewCall(child)
		} else {
			bl.NewStore(child, variable)
			funcCallNode.Child.VarPtr = variable
		}
		return child

	} else if funcCallNode.Val == "fn" {
		// TODO: validate
		funcFn := mod.NewFunc(
			fmt.Sprintf("fn-%p", funcCallNode),
			types.I32,
		)
		llBlock := funcFn.NewBlock("")
		res := gen(mod, llBlock, funcCallNode.Child, libs)
		llBlock.NewRet(res)
		return funcFn

	} else if funcCallNode.IsDeclare() {
		// means declaring global variable or function

		// TODO: validate
		funcName := funcCallNode.Child.Val
		retType := types.I32

		if funcName != "main" {
			retType = types.I32
		}

		function := mod.NewFunc(
			funcName,
			retType,
		)
		llBlock := function.NewBlock("")

		if funcName != "main" {
			res := gen(mod, llBlock, funcCallNode.Child, libs)
			llBlock.NewRet(res)
		} else {
			gen(mod, llBlock, funcCallNode.Child, libs)
			llBlock.NewRet(newI32("0"))
		}

	}
	return nil
}

func codegen(prog *mTypes.Program) *ir.Module {
	module := ir.NewModule()
	prog.BuiltinLibs = &mTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, prog.BuiltinLibs)

	for calls := prog.Declares; calls != nil; calls = calls.Next {
		gen(module, nil, calls, prog.BuiltinLibs)
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	ir := codegen(a.program)

	log.DebugMessage("code generated")
	log.Debug("[IR]\n%s\n", ir.String())

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
