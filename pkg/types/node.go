package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/log"
)

// ########
// Node
// ########
type NodeKind string

const (
	ND_PROGRAM_ROOT = NodeKind("ND_PROGRAM_ROOT")

	// reserved symbol
	ND_VAR_DECLARE   = NodeKind("ND_VAR_DECLARE")   // variables and functions
	ND_VAR_REFERENCE = NodeKind("ND_VAR_REFERENCE") // variables and functions
	ND_DECLARE       = NodeKind("ND_DECLARE")       // def
	ND_LAMBDA        = NodeKind("ND_LAMBDA")        // fn
	ND_BIND          = NodeKind("ND_BIND")          // let
	ND_EXPR          = NodeKind("ND_EXPR")          // set of functions
	ND_FUNCCALL      = NodeKind("ND_FUNCCALL")
	ND_LIBCALL       = NodeKind("ND_LIBCALL")

	// Arithmetic signature
	ND_ADD = NodeKind("ND_ADD") // +
	ND_SUB = NodeKind("ND_SUB") // -
	ND_MUL = NodeKind("ND_MUL") // *
	ND_DIV = NodeKind("ND_DIV") // /
	ND_MOD = NodeKind("ND_MOD") // mod
	ND_EQ  = NodeKind("ND_EQ")  // =

	// Logical signature
	ND_AND = NodeKind("ND_AND") // and
	ND_OR  = NodeKind("ND_OR")  // or

	// type
	ND_NIL = NodeKind("ND_NIL") // nil
	ND_INT = NodeKind("ND_INT") // 0-9
)

type Program struct {
	Declares    *Node
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
	Kind    NodeKind
	Next    *Node
	Child   *Node
	Cond    *Node
	Then    *Node
	Else    *Node
	Val     string
	Bind    *Node
	VarPtr  *ir.InstAlloca // local variable
	FuncPtr *ir.Func       // global variable
}

func (node *Node) IsDeclare() bool {
	return node.Kind == ND_DECLARE
}
func (node *Node) IsLambda() bool {
	return node.Kind == ND_LAMBDA
}
func (node *Node) IsBind() bool {
	return node.Kind == ND_BIND
}
func (node *Node) IsExpr() bool {
	return node.Kind == ND_EXPR
}
func (node *Node) IsVarDeclare() bool {
	return node.Kind == ND_VAR_DECLARE
}
func (node *Node) IsVarReference() bool {
	return node.Kind == ND_VAR_REFERENCE
}
func (node *Node) IsLibCall() bool {
	return node.Kind == ND_LIBCALL
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

func indicate(s string, depth int) {
	log.Debug(
		log.YELLOW(
			fmt.Sprintf(
				"%s [%s]",
				strings.Repeat("  ", depth),
				s,
			),
		))

}

func (node *Node) Debug(depth int) {
	if node == nil {
		return
	}
	log.Debug(
		log.BLUE(
			fmt.Sprintf(
				"%s %p %#+v %#+v",
				strings.Repeat("  ", depth),
				node,
				node.Kind,
				node.Val),
		),
	)

	switch node.Kind {
	case ND_BIND:
		indicate("bind", depth+1)
		node.Bind.Debug(depth + 1)

		indicate("child", depth+1)
		node.Child.Debug(depth + 1)
	default:
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

	// TODO: implement
	// NOTE: is type of BuiltinLibs Node?
	if prog.BuiltinLibs != nil {
		log.DebugMessage("[BuiltinLib]")
	}
}
