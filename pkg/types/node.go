package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/log"
)

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
	ND_IF            = NodeKind("ND_IF")            // if
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

	ND_SCALAR     = NodeKind("ND_SCALAR") // int32, string, bool
	ND_COLLECTION = NodeKind("ND_COLLECTION")
)

type ScalarType string

const (
	TY_INT32 = ScalarType("TY_INT32")
	TY_STR   = ScalarType("TY_STR")
	TY_NIL   = ScalarType("TY_NIL")
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
	FormatStr   *ir.Global
}

type Node struct {
	Kind    NodeKind
	Next    *Node
	Type    ScalarType
	Child   *Node
	Cond    *Node
	Then    *Node
	Else    *Node
	Val     string
	Len     uint64
	Bind    *Node
	VarPtr  value.Value // local variable
	FuncPtr *ir.Func    // global variable
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
func (node *Node) IsIf() bool {
	return node.Kind == ND_IF
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
func (node *Node) IsNary() bool {
	return node.Kind == ND_ADD
}
func (node *Node) IsBinary() bool {
	return node.Kind == ND_EQ
}
func (node *Node) IsInt32() bool {
	return node.Type == TY_INT32
}
func (node *Node) IsStr() bool {
	return node.Type == TY_STR
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
	if node.Kind == ND_SCALAR {
		log.Debug(
			log.BLUE(
				fmt.Sprintf(
					"%s %p %#+v %#+v %#+v %d",
					strings.Repeat("  ", depth),
					node,
					node.Kind,
					node.Type,
					node.Val,
					node.Len,
				),
			),
		)

	} else {
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
	}

	switch node.Kind {
	case ND_BIND:
		indicate(".Bind", depth+1)
		node.Bind.Debug(depth + 1)

		indicate(".Child", depth+1)
		node.Child.Debug(depth + 1)
	case ND_IF:
		indicate(".Cond", depth+1)
		node.Cond.Debug(depth + 1)

		indicate(".Then", depth+1)
		node.Then.Debug(depth + 1)

		indicate(".Else", depth+1)
		node.Else.Debug(depth + 1)
	default:
		node.Child.Debug(depth + 1)
	}
	if node.Next != nil {
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
