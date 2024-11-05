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
	app        = kingpin.New("modo", "Compiler for the modo programming language.").Version(VERSION).Author(AUTHOR)
	appVerbose = app.Flag("verbose", "Which log tags to show").Bool()
	appDebug   = app.Flag("debug", "Which log tags to show").Bool()

	buildCmd    = app.Command("build", "Build an executable.")
	buildOutput = buildCmd.Flag("output", "Output binary name.").Short('o').String()

	runCmd  = app.Command("run", "Build and run an executable.")
	runExec = runCmd.Flag("exec", "evaluate passing string").String()
)

func asemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: Error: %v, %s\n", err, out)
	}
	log.Debug(asmFile, "written asm: %s")
}
func compile(asmFile string, executableFile string) {
	out, err := exec.Command("clang", asmFile, "-o", executableFile).CombinedOutput()
	if err != nil {
		fmt.Printf("rouph: Error: %v, %s\n", err, out)
	}
	log.Debug(executableFile, "written executable: %s")
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
		log.Debug(artifactDir, "make dir: %s")

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
  log.Debug(artifactFilePrefix, "artifactFilePrefix = %s")

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)
	return llName, asmName, executableName

}

func doRun(output string, evaluatee string) int {
	m := codegen.Codegen()
	log.Debug("code generated")

	llName, asmName, executableName := prepareWorkingFile(output)

	os.WriteFile(llName, []byte(m.String()), 0600)
  log.Debug(llName, "written ll: %s")
	asemble(llName, asmName)
	compile(asmName, executableName)

	return 0

}
func showOpts(cmd string) {
	m := map[string]interface{}{}
	m["cmd"] = cmd
	m["buildOutput"] = *buildOutput
	m["exec"] = *runExec
	m["debug"] = *appDebug
	m["verbose"] = *appVerbose
	log.Debug(m, "options = %#+v")
}
func setLogLevel() {
	if *appVerbose {
		log.SetLevelInfo()
	}
	if *appDebug {
		log.SetLevelDebug()
	}
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel()
	showOpts(cmd)

	switch cmd {
	case runCmd.FullCommand():
		status := doRun(*buildOutput, *runExec)
		os.Exit(status)
	}
}
