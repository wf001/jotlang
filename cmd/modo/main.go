package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/codegen/llvm"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/parser"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app             = kingpin.New("modo", "Compiler for the modo programming language.").Version(VERSION).Author(AUTHOR)
	appVerbose      = app.Flag("verbose", "Use verbose log").Bool()
	appDebug        = app.Flag("debug", "Use debug log").Bool()
	appOutput       = app.Flag("output", "Write output to <OUTPUT>").Short('o').String()
	appPersistTemps = app.Flag("persist-temps", "Persist all of build artifacts").Bool()

	buildCmd = app.Command("build", "Build an executable.")

	runCmd  = app.Command("run", "Build and run an executable.")
	runExec = runCmd.Flag("exec", "evaluate <EXEC>").String()

	currentTime = time.Now().Unix()
)

func showOpts(cmd string) {
	m := map[string]interface{}{}
	m["cmd"] = cmd
	m["appOutput"] = *appOutput
	m["appPersistTemps"] = *appPersistTemps
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

func prepareWorkingFile(artifactFilePrefix string) (string, string, string) {
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

func asemble(llFile string, asmFile string) {
	out, err := exec.Command("llc", llFile, "-o", asmFile).CombinedOutput()
	if err != nil {
		log.Panic(map[string]interface{}{"err": err, "out": out, "llFile": llFile, "asmFile": asmFile}, "fail to asemble: %s")
	}
	log.Debug(asmFile, "written asm: %s")
}

func compile(asmFile string, executableFile string) {
	out, err := exec.Command("clang", asmFile, "-o", executableFile).CombinedOutput()
	if err != nil {
		log.Panic(map[string]interface{}{"err": err, "out": out, "artifactDir": executableFile}, "fail to compile: %s")
	}
	log.Debug(executableFile, "written executable: %s")
}

func run(executableFile string) int {
	cmd := exec.Command(executableFile)
	err := cmd.Run()
	log.Debug(executableFile, "executed: %s")
	if err != nil {
		log.Error(err, "fail to run: %s")
	}

	return cmd.ProcessState.ExitCode()
}

func doRun(workingDirPrefix string, evaluatee string) int {
  token := lexer.Lex(evaluatee)
  log.Debug(token, "code lexed token: %#+v")

  parser.Parse(token)
	log.Debug("code parsed")

	m := codegen.Codegen(token)
	log.Debug("code generated")
  log.Debug(m.String(), "IR = \n %s\n")

	llName, asmName, executableName := prepareWorkingFile(workingDirPrefix)

	err := os.WriteFile(llName, []byte(m.String()), 0600)
	if err != nil {
		log.Panic(map[string]interface{}{"err": err, "llName": llName}, "fail to write ll: %s")
	}
	log.Debug(llName, "written ll: %s")
	asemble(llName, asmName)
	compile(asmName, executableName)

	return run(executableName)
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel()
	showOpts(cmd)

	switch cmd {

	case runCmd.FullCommand():
		status := doRun(*appOutput, *runExec)
		os.Exit(status)
	}
}
