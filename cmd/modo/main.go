package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/pkg/codegen"
	"github.com/wf001/modo/pkg/io"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/parser"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app = kingpin.New("modo", "Compiler for the modo programming language.").
		Version(VERSION).
		Author(AUTHOR)
	appVerbose      = app.Flag("verbose", "Use verbose log").Bool()
	appDebug        = app.Flag("debug", "Use debug log").Bool()
	appOutput       = app.Flag("output", "Write output to <OUTPUT>").Short('o').String()
	appPersistTemps = app.Flag("persist-temps", "Persist all of build artifacts").Bool()

	buildCmd = app.Command("build", "Build an executable.")

	runCmd    = app.Command("run", "Build and run an executable.")
	runExec   = runCmd.Flag("exec", "evaluate <EXEC>").String()
	inputFile = runCmd.Arg("file", "source file").String()
)

type IAssebler interface {
	Assemble()
}

func Assemble(a IAssebler) {
	a.Assemble()
}

// HACK: It might be better to move it to different package?
func showOpts(cmd string) {
	m := map[string]interface{}{}
	m["cmd"] = cmd
	m["appOutput"] = *appOutput
	m["appPersistTemps"] = *appPersistTemps
	m["exec"] = *runExec
	m["debug"] = *appDebug
	m["verbose"] = *appVerbose
	m["inputFile"] = *inputFile
	log.Debug("options = %#+v", m)
}

func setLogLevel() {
	if *appVerbose {
		log.SetLevelInfo()
	}
	if *appDebug {
		log.SetLevelDebug()
	}
}

func compile(asmFile string, executableFile string) {
	out, err := exec.Command("clang", asmFile, "-o", executableFile).CombinedOutput()
	if err != nil {
		log.Panic(
			"fail to compile: %+v",
			map[string]interface{}{"err": err, "out": out, "artifactDir": executableFile},
		)
	}
	log.Debug("written executable: %s", executableFile)
}

func run(executableFile string) int {
	cmd := exec.Command(executableFile)
	// TODO: it works, but correctly?
	out, err := cmd.Output()
	log.Debug("executed: %s", executableFile)
	if err != nil {
		log.Error("fail to run: %s", err)
	}
	fmt.Println(string(out))

	return cmd.ProcessState.ExitCode()
}

func doBuild(workingDirPrefix string, evaluatee string) (error, string) {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)

	// Token -> Node
	node := parser.Parse(token)

	llName, asmName, executableName := io.PrepareWorkingFile(workingDirPrefix, currentTime)

	// Node -> AST -> write assembly
	codegen.Construct(node).Assemble(llName, asmName)

	// assembly file -> write executable
	compile(asmName, executableName)

	return nil, executableName
}

// HACK: It might be better if the return type matches that of doBuild
func doRun(workingDirPrefix string, evaluatee string) int {
	err, executableName := doBuild(workingDirPrefix, evaluatee)
	if err != nil {
		log.Panic("fail to run: %s", map[string]interface{}{"err": err, "llName": executableName})
	}
	return run(executableName)
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel()
	showOpts(cmd)

	arg := io.ReadFile(inputFile)

	switch cmd {

	case runCmd.FullCommand():
		if *runExec != "" {
			os.Exit(doRun(*appOutput, *runExec))
		} else {
			os.Exit(doRun(*appOutput, arg))
		}
	}
}
