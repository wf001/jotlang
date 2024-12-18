package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/pkg/codegen"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/parser"
	"github.com/wf001/modo/util"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app = kingpin.
		New("modo", "Compiler for the modo programming language.").
		Version(VERSION).
		Author(AUTHOR)
	appVerbose = app.Flag("verbose", "Use verbose log").Bool()
	appDebug   = app.Flag("debug", "Use debug log").Bool()
	appOutput  = app.Flag("output", "Write output to <OUTPUT>").Short('o').String()
	// TODO: rename to keep-intermediates
	appPersistTemps = app.Flag("persist-temps", "Persist all of build artifacts").Bool()

	buildCmd = app.Command("build", "Build an executable.")

	runCmd    = app.Command("run", "Build and run an executable.")
	runExec   = runCmd.Flag("exec", "evaluate <EXEC>").String()
	inputFile = runCmd.Arg("file", "source file").String()
)

type IAssebler interface {
	GenFrontend()
}

func GenFrontend(a IAssebler) {
	a.GenFrontend()
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
	_, err, errMsg := util.RunCommand("clang", asmFile, "-o", executableFile)
	if err != nil {
		log.Debug("artifactDir: %s", executableFile)
		log.Panic(
			"fail to compile: %v",
			map[string]interface{}{
				"err":     err,
				"message": errMsg,
			},
		)
	}
	log.Debug("written executable: %s", executableFile)
}

func run(executableFile string) int {
	out, err, errMsg := util.RunCommand(executableFile)
	// TODO: it works, but correctly?
	log.Debug("executed: %s", executableFile)
	if err != nil {
		log.Debug("artifactDir: %s", executableFile)
		log.Error(
			"fail to run: %+v",
			map[string]interface{}{
				"err":     err,
				"message": errMsg,
			},
		)
		return 1
	}
	fmt.Println(out)

	return 0
}

func doBuild(workingDirPrefix string, evaluatee string) (error, string) {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)

	// Token -> Node
	node := parser.Parse(token)

	llName, asmName, executableName := util.PrepareWorkingFile(workingDirPrefix, currentTime)

	// Node -> AST -> write assembly
	codegen.Construct(node).GenFrontend(llName, asmName)
	codegen.Assemble(llName, asmName)

	// assembly file -> write executable
	compile(asmName, executableName)

	return nil, executableName
}

// TODO: refactoring
func doRunLLI(workingDirPrefix string, evaluatee string) int {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)

	// Token -> Node
	node := parser.Parse(token)

	llName, asmName, _ := util.PrepareWorkingFile(workingDirPrefix, currentTime)

	// Node -> AST -> write assembly
	codegen.Construct(node).GenFrontend(llName, asmName)

	out, err, errMsg := util.RunCommand("lli", llName)
	// TODO: it works, but correctly?
	log.Debug("executed: %s", llName)
	if err != nil {
		log.Debug("artifactDir: %s", llName)
		log.Error(
			"fail to run: %+v",
			map[string]interface{}{
				"err":     err,
				"message": errMsg,
			},
		)
		return 1
	}
	fmt.Println(out)

	return 0
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

	switch cmd {

	case runCmd.FullCommand():
		if *runExec != "" {
			os.Exit(doRun(*appOutput, *runExec))
		} else {
			if inputFile != nil {
				arg := util.ReadFile(inputFile)
				os.Exit(doRunLLI(*appOutput, arg))
			} else {
				log.Panic("fail to run, input must be specified")
			}
		}
	}
}
