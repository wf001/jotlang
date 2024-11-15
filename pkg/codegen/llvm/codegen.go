package codegen

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	llirTypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/types"
)

type assembler struct {
	node *types.Node
}

func newInt32(s string) *constant.Int {

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newInt32: %s", err)
	}
	return constant.NewInt(llirTypes.I32, i)
}

func doAsemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic("fail to asemble: %+v", map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile})
	}
	log.Debug("written asm: %s", asmFile)
}

func codegen(mb *ir.Block, node *types.Node) value.Value {
	switch node.Kind {
	case types.ND_INT:
		return newInt32(node.Val)
	case types.ND_ADD:
		fst := codegen(mb, node.Child)
		snd := codegen(mb, node.Child.Next)
		res := mb.NewAdd(fst, snd)
		child := node.Child.Next.Next
		for ; child != nil; child = child.Next {
			fst = res
			snd = codegen(mb, child)
			res = mb.NewAdd(fst, snd)
		}
		return res
	}
	return nil
}

func Construct(node *types.Node) *assembler {
	return &assembler{
		node: node,
	}
}

func (a assembler) Assemble(llName string, asmName string) {
	ir := ir.NewModule()
	funcMain := ir.NewFunc(
		"main",
		llirTypes.I32,
	)
	llBlock := funcMain.NewBlock("")

	res := codegen(llBlock, a.node)
	llBlock.NewRet(res)

	log.DebugMessage("code generated")
	log.Debug("IR = \n %s\n", ir.String())

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)
}
