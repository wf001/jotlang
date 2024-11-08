package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/wf001/modo/internal/log"
)

func newInt32(s string) *constant.Int {

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newInt32: %s", err)
	}
	res := constant.NewInt(types.I32, i)
	log.Debug("%#v", res)
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
		log.Debug("make dir: %s", artifactDir)

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
	log.Debug("artifactFilePrefix = %s", artifactFilePrefix)
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
func Codegen(s string) *ir.Module {
	m := ir.NewModule()
	funcMain := m.NewFunc(
		"main",
		types.I32,
	)
	mb := funcMain.NewBlock("")

	mb.NewRet(newInt32(s))
	return m
}

func Assemble(node string, workingDirPrefix string, currentTime int64) (string, string) {
	m := Codegen(node)
	log.DebugMessage("code generated")
	log.Debug("IR = \n %s\n", m.String())

	llName, asmName, executableName := prepareWorkingFile(workingDirPrefix, currentTime)

	err := os.WriteFile(llName, []byte(m.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %s", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)
	doAsemble(llName, asmName)

	return asmName, executableName

}

// func Codegen(s string) *ir.Module {
// 	m := ir.NewModule()
//
// 	globalG := m.NewGlobalDef("g", constant.NewInt(types.I32, 58))
//
// 	funcAdd := m.NewFunc("add", types.I32,
// 		ir.NewParam("x", types.I32),
// 		ir.NewParam("y", types.I32),
// 	)
// 	ab := funcAdd.NewBlock("")
// 	ab.NewRet(ab.NewAdd(funcAdd.Params[0], funcAdd.Params[1]))
//
// 	funcMain := m.NewFunc(
// 		"main",
// 		types.I32,
// 	) // omit parameters
// 	mb := funcMain.NewBlock("") // llir/llvm would give correct default name for block without name
// 	mb.NewRet(mb.NewCall(funcAdd, constant.NewInt(types.I32, 59), mb.NewLoad(types.I32, globalG)))
// 	return m
//
// }
