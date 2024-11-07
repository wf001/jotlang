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
		log.Panic(err, "fail to newInt32: %s")
	}
	res := constant.NewInt(types.I32, i)
	log.Debug(res, "%#v")
	return res
}

func prepareWorkingFile(artifactFilePrefix string, currentTime int64) (string, string, string) {
	if artifactFilePrefix == "" {
		generated := "generated"
		artifactDir := fmt.Sprintf("%s/%d", generated, currentTime)
		out, err := exec.Command("mkdir", "-p", artifactDir).CombinedOutput()
		if err != nil {
			log.Panic(map[string]interface{}{"err": err, "out": out, "artifactDir": artifactDir}, "fail to make directory: %s")
		}
		log.Debug(artifactDir, "make dir: %s")

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
	log.Debug(artifactFilePrefix, "artifactFilePrefix = %s")
	log.Info(artifactFilePrefix, "persist all of build artifact in %s")

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)

	return llName, asmName, executableName
}

func doAsemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic(map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile}, "fail to asemble: %s")
	}
	log.Debug(asmFile, "written asm: %s")
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
	log.Debug("code generated")
	log.Debug(m.String(), "IR = \n %s\n")

	llName, asmName, executableName := prepareWorkingFile(workingDirPrefix, currentTime)

	err := os.WriteFile(llName, []byte(m.String()), 0600)
	if err != nil {
		log.Panic(map[string]interface{}{"err": err, "llName": llName}, "fail to write ll: %s")
	}
	log.Debug(llName, "written ll: %s")
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
