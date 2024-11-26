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
	node *modoTypes.Node
}

var coreLibs *modoTypes.BuiltinLibProp

var naryMap = map[modoTypes.NodeKind]func(*ir.Block, value.Value, value.Value) value.Value{
	modoTypes.ND_ADD: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewAdd(x, y)
	},
	modoTypes.ND_SUB: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewSub(x, y)
	},
	modoTypes.ND_MUL: func(block *ir.Block, x, y value.Value) value.Value {
		return block.NewMul(x, y)
	},
}
var binaryMap = map[modoTypes.NodeKind]func(*ir.Block, value.Value, value.Value) value.Value{
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

func gen(mb *ir.Block, node *modoTypes.Node) value.Value {

	if node.IsInteger() {
		return newI32(node.Val)

	} else if node.IsNary() {
		// nary takes more than 2 arguments
		child := node.Child
		fst := gen(mb, child)

		child = child.Next
		snd := gen(mb, child)

		nary := naryMap[node.Kind]
		res := nary(mb, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = gen(mb, child)
			res = nary(mb, fst, snd)
		}
		return res

	} else if node.IsBinary() {
		// binary takes exactly 2 arguments
		child := node.Child
		fst := gen(mb, child)

		child = child.Next
		snd := gen(mb, child)

		binary := binaryMap[node.Kind]
		res := binary(mb, fst, snd)

		return res
	} else if node.IsLibrary() {
		// means calling standard library
		arg := gen(mb, node.Child)
		libFunc := libraryMap[node.Val]
		libFunc(mb, coreLibs, arg)

		return newI32("0")
	}
	return nil
}

func codegen(node *modoTypes.Node) *ir.Module {
	module := ir.NewModule()
	coreLibs = &modoTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, coreLibs)

	funcMain := module.NewFunc(
		"main",
		types.I32,
	)
	llBlock := funcMain.NewBlock("")

	res := gen(llBlock, node)
	llBlock.NewRet(res)
	return module
}

func Construct(node *modoTypes.Node) *assembler {
	return &assembler{
		node: node,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	ir := codegen(a.node)

	log.DebugMessage("code generated")
	log.Debug("IR = \n %s\n", ir.String())

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
