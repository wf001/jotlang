package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

const (
	VERSION = "0.0.1"
	AUTHOR  = "wf001"
)

var (
	app = kingpin.New("jotc", "Compiler for the Jot programming language.").Version(VERSION).Author(AUTHOR)

	//logLevel = app.Flag("loglevel", "Set the level of logging to show").Default("info").Enum("debug", "verbose", "info", "warning", "error")
	verbose  = app.Flag("verbose", "Which log tags to show").Short('v').Default("false").Bool()

	buildCom    = app.Command("build", "Build an executable.")
	buildOutput = buildCom.Flag("output", "Output binary name.").Short('o').Default("out").String()

	runCom    = app.Command("run", "Build and run an executable.")
)

func doMain() {
	m := ir.NewModule()

	globalG := m.NewGlobalDef("g", constant.NewInt(types.I32, 4))

	funcAdd := m.NewFunc("add", types.I32,
		ir.NewParam("x", types.I32),
		ir.NewParam("y", types.I32),
	)
	ab := funcAdd.NewBlock("")
	ab.NewRet(ab.NewAdd(funcAdd.Params[0], funcAdd.Params[1]))

	funcMain := m.NewFunc(
		"main",
		types.I32,
	) // omit parameters
	mb := funcMain.NewBlock("") // llir/llvm would give correct default name for block without name
	mb.NewRet(mb.NewCall(funcAdd, constant.NewInt(types.I32, 1), mb.NewLoad(types.I32, globalG)))

	fmt.Println(m)

}

func main() {
  switch kingpin.MustParse(app.Parse(os.Args[1:])) {
  // Register user
  case buildCom.FullCommand():
    doMain()
  }
}
