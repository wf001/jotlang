package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/internal/log"
	codegen "github.com/wf001/modo/pkg/codegen/llvm"
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

func Compile(asmFile string, executableFile string) {
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

func doBuild(workingDirPrefix string, evaluatee string) (int, string) {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)
	log.Debug(token, "code lexed token: %#+v")

	// Token -> Node
	parser.Parse(token)
	log.Debug("code parsed")

	// Node -> LLVM IR -> write assembly
	asmName, executableName := codegen.DoAssemble(token, workingDirPrefix, currentTime)

  // assembly file -> write executable
	Compile(asmName, executableName)

	return 0, executableName
}

func doRun(workingDirPrefix string, evaluatee string) int {
	err, executableName := doBuild(workingDirPrefix, evaluatee)
	if err != 0 {
		log.Panic(map[string]interface{}{"err": err, "llName": executableName}, "fail to run: %s")
	}
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
