package codegen

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/lib"
	"github.com/wf001/modo/pkg/log"
	modoTypes "github.com/wf001/modo/pkg/types"
)

type assembler struct {
	program *modoTypes.Program
}

var operatorMap = map[modoTypes.NodeKind]func(*ir.Block, value.Value, value.Value) value.Value{
	// nary
	modoTypes.ND_ADD: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	},
	modoTypes.ND_SUB: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewSub(x, y)
	},
	modoTypes.ND_MUL: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewMul(x, y)
	},
	// binary
	modoTypes.ND_EQ: func(block *ir.Block, x, y value.Value) value.Value {
		res := block.NewICmp(enum.IPredEQ, x, y)
		return block.NewZExt(res, types.I32)
	},
}
var libraryMap = map[string]func(*ir.Block, *modoTypes.BuiltinLibProp, value.Value){
	"prn": func(block *ir.Block, libs *modoTypes.BuiltinLibProp, arg value.Value) {
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
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic(
			"fail to asemble: %+v",
			map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile},
		)
	}
	log.Debug("written asm: %s", asmFile)
}

func gen(mb *ir.Block, funcCallNode *modoTypes.Node, libs *modoTypes.BuiltinLibProp) value.Value {

	if funcCallNode.IsInteger() {
		return newI32(funcCallNode.Val)

	} else if funcCallNode.IsNary() {
		// nary takes more than 2 arguments
		child := funcCallNode.Child
		fst := gen(mb, child, libs)

		child = child.Next
		snd := gen(mb, child, libs)

		nary := operatorMap[funcCallNode.Kind]
		res := nary(mb, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mb, child, libs)
			res = nary(mb, fst, snd)
		}
		return res

	} else if funcCallNode.IsBinary() {
		// binary takes exactly 2 arguments
		child := funcCallNode.Child
		fst := gen(mb, child, libs)

		child = child.Next
		snd := gen(mb, child, libs)

		binary := operatorMap[funcCallNode.Kind]
		res := binary(mb, fst, snd)

		return res
	} else if funcCallNode.IsLibrary() {
		// means calling standard library
		arg := gen(mb, funcCallNode.Child, libs)
		libFunc := libraryMap[funcCallNode.Val]
		libFunc(mb, libs, arg)

		return newI32("0")
	}
	return nil
}

func codegen(prog *modoTypes.Program) *ir.Module {
	module := ir.NewModule()
	prog.BuiltinLibs = &modoTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, prog.BuiltinLibs)

	funcMain := module.NewFunc(
		"main",
		types.I32,
	)
	llBlock := funcMain.NewBlock("")

	res := gen(llBlock, prog.FuncCalls, prog.BuiltinLibs)
	llBlock.NewRet(res)
	return module
}

func Construct(program *modoTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	ir := codegen(a.program)

	log.DebugMessage("code generated")
	log.Debug("[IR]  \n%s\n", ir.String())

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
