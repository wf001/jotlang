package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/wf001/modo/internal/log"
	modoTypes "github.com/wf001/modo/pkg/types"
)

func newInt32(s string) *constant.Int {

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newInt32: %s", err)
	}
	res := constant.NewInt(types.I32, i)
	return res
}

func prepareWorkingFile(artifactFilePrefix string, currentTime int64) (string, string, string) {
	if artifactFilePrefix == "" {
		generated := "generated"
		artifactDir := fmt.Sprintf("%s/%d", generated, currentTime)
		out, err := exec.Command("mkdir", "-p", artifactDir).CombinedOutput()
		if err != nil {
			log.Panic("fail to make directory: %s", map[string]interface{}{"err": err, "out": out, "artifactDir": artifactDir})
		}
		log.Debug(log.YELLOW("make dir: %s"), artifactDir)

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
	log.Debug(log.YELLOW("artifactFilePrefix = %s"), artifactFilePrefix)
	log.Info("persist all of build artifact in %s", artifactFilePrefix)

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)

	return llName, asmName, executableName
}

func doAsemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic("fail to asemble: %s", map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile})
	}
	log.Debug("written asm: %s", asmFile)
}
func Codegen(mb *ir.Block, node *modoTypes.Node) value.Value {
	switch node.Kind {
	case modoTypes.ND_INT:
		return newInt32(node.Val)
	case modoTypes.ND_ADD:
		fst := Codegen(mb, node.Child)
		snd := Codegen(mb, node.Child.Next)
		res := mb.NewAdd(fst, snd)
		child := node.Child.Next.Next
		for ; child != nil; child = child.Next {
			fst = res
			snd = Codegen(mb, child)
			res = mb.NewAdd(fst, snd)
		}
		return res
	}
	return nil
}

func Assemble(node *modoTypes.Node, workingDirPrefix string, currentTime int64) (string, string) {
	ir := ir.NewModule()
	funcMain := ir.NewFunc(
		"main",
		types.I32,
	)
	llBlock := funcMain.NewBlock("")

	res := Codegen(llBlock, node)
	llBlock.NewRet(res)

	log.DebugMessage("code generated")
	log.Debug("IR = \n %s\n", ir.String())

	llName, asmName, executableName := prepareWorkingFile(workingDirPrefix, currentTime)

	err := os.WriteFile(llName, []byte(ir.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %s", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

	doAsemble(llName, asmName)

	return asmName, executableName

}
