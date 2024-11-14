package main

import (
	"bufio"
	"os"
	"os/exec"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/wf001/modo/internal/log"
	codegen "github.com/wf001/modo/pkg/codegen/llvm"
	"github.com/wf001/modo/pkg/lexer"
	"github.com/wf001/modo/pkg/parser"
	"github.com/wf001/modo/pkg/types"
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

	runCmd    = app.Command("run", "Build and run an executable.")
	runExec   = runCmd.Flag("exec", "evaluate <EXEC>").String()
	inputFile = runCmd.Arg("file", "source file").String()
)

type IParser interface {
	Parse() *types.Node
}

func Parse(p IParser) {
	p.Parse()
}

type IAssebler interface {
	Assemble() (string, string)
}

func Assemble(a IAssebler) {
	a.Assemble()
}

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
		log.Panic("fail to compile: %+v", map[string]interface{}{"err": err, "out": out, "artifactDir": executableFile})
	}
	log.Debug("written executable: %s", executableFile)
}

func run(executableFile string) int {
	cmd := exec.Command(executableFile)
	err := cmd.Run()
	log.Debug("executed: %s", executableFile)
	if err != nil {
		log.Error("fail to run: %s", err)
	}

	return cmd.ProcessState.ExitCode()
}

func doBuild(workingDirPrefix string, evaluatee string) (int, string) {
	currentTime := time.Now().Unix()
	// string -> Token
	token := lexer.Lex(evaluatee)
	log.DebugMessage("code lexed ")

	// Token -> Node
	node := parser.ConstructParser(token).Parse()
	log.DebugMessage("code parsed")

	// Node -> AST -> write assembly
	asmName, executableName := codegen.ConstructAssembler(node).Assemble(workingDirPrefix, currentTime)

	// assembly file -> write executable
	compile(asmName, executableName)

	return 0, executableName
}

func doRun(workingDirPrefix string, evaluatee string) int {
	err, executableName := doBuild(workingDirPrefix, evaluatee)
	if err != 0 {
		log.Panic("fail to run: %s", map[string]interface{}{"err": err, "llName": executableName})
	}
	return run(executableName)
}

func readFile(inputFile *string) string {
	data, _ := os.Open(*inputFile)
	defer data.Close()

	scanner := bufio.NewScanner(data)
	var input_arr = ""
	for scanner.Scan() {
		input_arr += scanner.Text()
	}
	return input_arr
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel()
	showOpts(cmd)

	arg := readFile(inputFile)

	switch cmd {

	case runCmd.FullCommand():
		if *runExec != "" {
			os.Exit(doRun(*appOutput, *runExec))
		} else {
			os.Exit(doRun(*appOutput, arg))
		}
	}
}
