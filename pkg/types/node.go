package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
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
	ND_SCALAR     = NodeKind("ND_SCALAR")     // int32, string, bool
	ND_COLLECTION = NodeKind("ND_COLLECTION") // vector, map, list
)

type ModoType string

const (
	TY_INT32 = ModoType("TY_INT32")
	TY_STR   = ModoType("TY_STR")
	TY_NIL   = ModoType("TY_NIL")
)

type Program struct {
	Declares    *Node
	BuiltinLibs *BuiltinLibProp
	GlobalStr   []*ir.InstLoad
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
	FormatSpace *ir.Global
	FormatCR    *ir.Global
}

type Node struct {
	Kind    NodeKind
	Next    *Node
	Type    ModoType
	Child   *Node
	Cond    *Node
	CondRet *ir.InstAlloca
	Then    *Node
	Else    *Node
	Val     string
	Len     uint64 // the number of bytes, only used with string type
	Bind    *Node
	Args    *Node
	VarPtr  value.Value // binded local variable
	FuncPtr *ir.Func    // declared function, library function
	IRValue value.Value //
}

// pred kind
func (node *Node) IsKind(kind NodeKind) bool {
	return node.Kind == kind
}

// Check if Kind is nary operator that takes more than 2 arguments
func (node *Node) IsKindNary() bool {
	return node.IsKind(ND_ADD)
}

// Check if Kind is binary operator that takes more than 2 arguments
func (node *Node) IsKindBinary() bool {
	return node.IsKind(ND_EQ)
}

// pred type
func (node *Node) IsType(ty ModoType) bool {
	return node.Type == ty
}

// naming
func (node *Node) GetUnnamedFuncName() string {
	return fmt.Sprintf("fn.%s.%p", "unnamed", node)
}

func (node *Node) GetFuncName() string {
	return fmt.Sprintf("fn.%s", node.Val)
}

func (node *Node) GetBlockName(s string) string {
	return fmt.Sprintf("%s.%p", s, node)
}

// Returns the last node of the linked list.
func (node *Node) GetLastNode() *Node {
	lastNode := node
	for lastNode.Next != nil {
		lastNode = lastNode.Next
	}
	return lastNode
}

// Get LLVM type from corresponding custom type
func (node *Node) GetLLVMType() types.Type {
	var retType types.Type

	if node.IsType(TY_INT32) {
		return types.I32
	} else if node.IsType(TY_STR) {
		return types.I8Ptr
	} else if node.IsType(TY_NIL) {
		return types.Void
	}
	return retType
}

// debug
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

	switch node.Kind {
	case ND_BIND:
		indicate(".Bind", depth+1)
		node.Bind.Debug(depth + 1)

		indicate(".Child", depth+1)
		node.Child.Debug(depth + 1)
	case ND_LAMBDA:
		indicate(".Args", depth+1)
		node.Args.Debug(depth + 1)

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
