package types

import (
	"fmt"
	"strings"

	"github.com/wf001/modo/internal/log"
)

// ########
// Token
// ########

type TokenKind = string

const (
	TK_NUM      = TokenKind("TK_NUM")
	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_PAREN    = TokenKind("TK_PAREN")
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
	BRACKETS_REG_EXP = fmt.Sprintf("[%s%s%s\\%s%s%s]", PARREN_OPEN, PARREN_CLOSE, BRACKET_OPEN, BRACKET_CLOSE, BRACE_OPEN, BRACE_CLOSE)

	OPERATOR_ADD      = "+"
	OPERATOR_SUB      = "-"
	OPERATOR_MUL      = "*"
	OPERATOR_DIV      = "/"
	OPERATOR_LT       = "<"
	OPERATOR_GT       = ">"
	OPERATORS_REG_EXP = fmt.Sprintf("[%s\\%s%s%s%s%s]", OPERATOR_ADD, OPERATOR_SUB, OPERATOR_MUL, OPERATOR_DIV, OPERATOR_LT, OPERATOR_GT)

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
	return tok.Kind == TK_OPERATOR && tok.Val == OPERATOR_ADD
}

func (tok *Token) IsNum() bool {
	return tok.Kind == TK_NUM
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
	ND_ADD = NodeKind("ND_ADD") // +
	ND_SUB = NodeKind("ND_SUB") // -
	ND_MUL = NodeKind("ND_MUL") // *
	ND_DIV = NodeKind("ND_DIV") // /
	ND_INT = NodeKind("ND_INT") // 0-9
)

type Node struct {
	Kind  NodeKind
	Next  *Node
	Child *Node
	Cond  *Node
	Then  *Node
	Else  *Node
	Init  *Node
	Inc   *Node
	Body  *Node
	Func  string
	Args  *Node
	Val   string
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
	}
	if node.Next != nil && node.Kind != ND_INT {
		node.Next.DebugNode(depth)
	}
}
