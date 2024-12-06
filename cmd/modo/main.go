package main

import (
	"fmt"
	"os"
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
	_, err, errMsg := io.RunCommand("clang", asmFile, "-o", executableFile)
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
	out, err, errMsg := io.RunCommand(executableFile)
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

	switch cmd {

	case runCmd.FullCommand():
		if *runExec != "" {
			os.Exit(doRun(*appOutput, *runExec))
		} else {
			if inputFile != nil {
				arg := io.ReadFile(inputFile)
				os.Exit(doRun(*appOutput, arg))
			} else {
				log.Panic("fail to run, input must be specified")
			}
		}
	}
}
