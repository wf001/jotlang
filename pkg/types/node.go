package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/log"
)

// ########
// Token
// ########

// ########
// Node
// ########
type NodeKind string

const (
	// Arithmetic ND_CORE
	ND_ADD = NodeKind("ND_ADD") // +
	ND_SUB = NodeKind("ND_SUB") // -
	ND_MUL = NodeKind("ND_MUL") // *
	ND_DIV = NodeKind("ND_DIV") // /
	ND_MOD = NodeKind("ND_MOD") // mod
	ND_EQ  = NodeKind("ND_EQ")  // =
	// Logical
	ND_AND = NodeKind("ND_AND") // and
	ND_OR  = NodeKind("ND_OR")  // or
	// type
	ND_NIL = NodeKind("ND_NIL") // nil
	ND_INT = NodeKind("ND_INT") // 0-9
	// lib
	ND_LIB = NodeKind("ND_LIB") // standard library
	// func
	ND_FUNC    = NodeKind("ND_FUNC")
	ND_VAR     = NodeKind("ND_VAR")
	ND_DECLARE = NodeKind("ND_DECLARE")
)

type Program struct {
	Declares    *Node
	FuncCalls   *Node
	BuiltinLibs *BuiltinLibProp
}

type BuiltinLibProp struct {
	GlobalVar *BuiltinGlobalVarsProp
	Printf    *BuiltinProp
}

// HACK: should remove llir/ir reference from this namespace
type BuiltinProp struct {
	FuncPtr *ir.Func
}

// HACK: should remove llir/ir reference from this namespace
type BuiltinGlobalVarsProp struct {
	FormatDigit *ir.Global
}

type Node struct {
	Kind       NodeKind
	Next       *Node
	Child      *Node
	Cond       *Node
	Then       *Node
	Else       *Node
	Val        string
	GlobalVars *ir.Global
	LocalVars  *ir.InstAlloca
}

func (node *Node) IsInteger() bool {
	return node.Kind == ND_INT
}
func (node *Node) IsNary() bool {
	return node.Kind == ND_ADD
}
func (node *Node) IsBinary() bool {
	return node.Kind == ND_EQ
}
func (node *Node) IsLibrary() bool {
	return node.Kind == ND_LIB
}

func (node *Node) Debug(depth int) {
	if node == nil {
		return
	}
	log.Debug(
		log.BLUE(
			fmt.Sprintf(
				"%s %p %#+v %#+v",
				strings.Repeat("\t", depth),
				node,
				node.Kind,
				node.Val),
		),
	)

	switch node.Kind {
	case ND_INT:
		node.Next.Debug(depth)
	case ND_ADD:
		node.Child.Debug(depth + 1)
	case ND_EQ:
		node.Child.Debug(depth + 1)
	case ND_LIB:
		node.Child.Debug(depth + 1)
	case ND_DECLARE:
		node.Child.Debug(depth + 1)
	case ND_VAR:
		node.Child.Debug(depth + 1)
	}
	if node.Next != nil && node.Kind != ND_INT {
		node.Next.Debug(depth)
	}
}
func (prog *Program) Debug(depth int) {
	if prog.Declares != nil {
		log.DebugMessage("[Declares]")
		prog.Declares.Debug(0)
	}

	if prog.FuncCalls != nil {
		log.DebugMessage("[FuncCalls]")
		prog.FuncCalls.Debug(0)
	}

	if prog.BuiltinLibs != nil {
		log.DebugMessage("[BuiltinLib]")
	}
}
