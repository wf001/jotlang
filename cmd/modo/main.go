package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/codegen"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app     = kingpin.New("modo", "Compiler for the modo programming language.").Version(VERSION).Author(AUTHOR)
	verbose = app.Flag("verbose", "Which log tags to show").Short('v').Default("false").Bool()

	buildCmd    = app.Command("build", "Build an executable.")
	buildOutput = buildCmd.Flag("output", "Output binary name.").Short('o').String()

	runCmd    = app.Command("run", "Build and run an executable.")
	buildExec = runCmd.Flag("exec", "evaluate passing string").String()
)

func asemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: Error: %v, %s\n", err, out)
	}
	log.Info(asmFile, "written asm: %s")
}
func compile(asmFile string, executableFile string) {
	out, err := exec.Command("clang", asmFile, "-o", executableFile).CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: Error: %v, %s\n", err, out)
	}
	log.Info(executableFile, "written executable: %s")
}

func prepareWorkingFile(artifactFilePrefix string) (string, string, string) {
	if artifactFilePrefix == "" {
		currentTime := time.Now().Unix()
		generated := "generated"
		artifactDir := fmt.Sprintf("%s/%d", generated, currentTime)
		_, err := exec.Command("mkdir", "-p", artifactDir).CombinedOutput()
		if err != nil {
			fmt.Printf("rouph: Error: %v", err)
		}
		log.Info(artifactDir, "make dir: %s")

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)
	return llName, asmName, executableName

}

func doRun(output string, evaluatee string) int {

	log.Info("generating LL")
	m := codegen.Codegen()
	log.Info("generated LL")
	llName, asmName, executableName := prepareWorkingFile(output)

	os.WriteFile(llName, []byte(m.String()), 0600)

	asemble(llName, asmName)
	compile(asmName, executableName)

	log.Info(llName, "written LL: %s")

	return 0

}

func main() {
	log.Info("main start")
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	log.Info(cmd, "%#+v")
	log.Info(*buildOutput, "output opt = %#+v")
	log.Info(*buildExec, "exec opt = %#+v")

	switch cmd {
	// Register user
	case runCmd.FullCommand():
		log.Info("call build")
		os.Exit(doRun(*buildOutput, *buildExec))
	}
	log.Info("main end")
	os.Exit(1)
}
