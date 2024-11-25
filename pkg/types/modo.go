package types

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/wf001/modo/pkg/log"
	llirTypes "github.com/wf001/modo/pkg/types/llir"
)

// ########
// Token
// ########

type TokenKind = string

const (
	TK_NUM      = TokenKind("TK_NUM")
	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_PAREN    = TokenKind("TK_PAREN")
	TK_LIB      = TokenKind("TK_LIB")
	TK_EOL      = TokenKind("TK_EOL")
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}

var (
	FRACTIONAL_REG_EXP = `\d+\.\d+`
	INTEGER_REG_EXP    = `\d+`

	PARREN_OPEN      = "("
	PARREN_CLOSE     = ")"
	BRACKET_OPEN     = "["
	BRACKET_CLOSE    = "]"
	BRACE_OPEN       = "{"
	BRACE_CLOSE      = "}"
	BRACKETS_REG_EXP = fmt.Sprintf(
		"[%s%s%s\\%s%s%s]",
		PARREN_OPEN,
		PARREN_CLOSE,
		BRACKET_OPEN,
		BRACKET_CLOSE,
		BRACE_OPEN,
		BRACE_CLOSE,
	)

	NARY_OPERATOR_ADD = "+"
	NARY_OPERATOR_SUB = "-"
	NARY_OPERATOR_MUL = "*"
	NARY_OPERATOR_DIV = "/"

	BINARY_OPERATOR_EQ = "="
	BINARY_OPERATOR_LT = "<"
	BINARY_OPERATOR_GT = ">"
	OPERATORS_REG_EXP  = fmt.Sprintf(
		"[%s\\%s%s%s%s%s%s]",
		NARY_OPERATOR_ADD,
		NARY_OPERATOR_SUB,
		NARY_OPERATOR_MUL,
		NARY_OPERATOR_DIV,
		BINARY_OPERATOR_EQ,
		BINARY_OPERATOR_LT,
		BINARY_OPERATOR_GT,
	)

	THREADING_FIRST   = "->"
	THREADING_LAST    = "->>"
	THREADING_REG_EXP = fmt.Sprintf("%s|%s", THREADING_FIRST, THREADING_LAST)

	EXPR_IF        = "if"
	EXPR_COND      = "cond"
	BRANCH_REG_EXP = fmt.Sprintf("\\b(%s|%s)\\b", EXPR_IF, EXPR_COND)

	EXPR_DEF           = "def"
	EXPR_DEFN          = "defn"
	EXPR_LET           = "let"
	DEFINITION_REG_EXP = fmt.Sprintf("\\b(%s|%s|%s)\\b", EXPR_DEF, EXPR_DEFN, EXPR_LET)

	RESERVED_REG_EXP = strings.Join(
		[]string{
			THREADING_REG_EXP,
			BRANCH_REG_EXP,
			DEFINITION_REG_EXP,
		},
		"|",
	)
	LIB_CORE_PRN     = "prn"
	LIB_CORE_REG_EXP = fmt.Sprintf("\\b(%s)\\b", LIB_CORE_PRN)

	USER_DEFINED_REG_EXP = `\w+`

	ALL_REG_EXP = fmt.Sprintf(
		"%s",
		strings.Join(
			[]string{
				FRACTIONAL_REG_EXP,
				INTEGER_REG_EXP,
				RESERVED_REG_EXP,
				OPERATORS_REG_EXP,
				BRACKETS_REG_EXP,
				LIB_CORE_PRN,
				USER_DEFINED_REG_EXP,
			},
			"|",
		),
	)
)

func (tok *Token) IsParenOpen() bool {
	return tok.Kind == TK_PAREN && tok.Val == PARREN_OPEN
}

func (tok *Token) IsParenClose() bool {
	return tok.Kind == TK_PAREN && tok.Val == PARREN_CLOSE
}

func (tok *Token) IsOperationAdd() bool {
	return tok.Kind == TK_OPERATOR && tok.Val == NARY_OPERATOR_ADD
}

func (tok *Token) IsNum() bool {
	return tok.Kind == TK_NUM
}

func (tok *Token) IsLibrary() bool {
	return tok.Kind == TK_LIB
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}

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
)

type Program struct {
	Declares   *Node
	FuncCalls  *Node
	GlobalVar  map[string]*llirTypes.LLIRAlloca
	BuiltinLib *BuiltinLibProp
}

type Node struct {
	Kind     NodeKind
	Next     *Node
	Child    *Node
	Cond     *Node
	Then     *Node
	Else     *Node
	Val      string
	LocalVar *llirTypes.LLIRAlloca
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

func (node *Node) DebugNode(depth int) {
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
		node.Next.DebugNode(depth)
	case ND_ADD:
		node.Child.DebugNode(depth + 1)
	case ND_EQ:
		node.Child.DebugNode(depth + 1)
	case ND_LIB:
		node.Child.DebugNode(depth + 1)
	}
	if node.Next != nil && node.Kind != ND_INT {
		node.Next.DebugNode(depth)
	}
}

// Built-in Library
type BuiltinLibProp struct {
	GlobalVar *BuiltinGlobalVarsProp
	Printf    *BuiltinProp
}

type BuiltinProp struct {
	FuncPtr *ir.Func
}

type BuiltinGlobalVarsProp struct {
	FormatDigit *ir.Global
}
