package types

import (
	"fmt"
	"strings"

	"github.com/wf001/modo/pkg/log"
)

type TokenKind = string

const (
	TK_NUM      = TokenKind("TK_NUM")
	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_PAREN    = TokenKind("TK_PAREN")
	TK_LIB      = TokenKind("TK_LIB")
	TK_VARIABLE = TokenKind("TK_VARIABLE")

	TK_RESERVED = TokenKind("TK_RESERVED")
	TK_DECLARE  = TokenKind("TK_DECLARE")
	TK_LAMBDA   = TokenKind("TK_LAMBDA")
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

	SYMBOL_THREADING_FIRST = "->"
	SYMBOL_THREADING_LAST  = "->>"
	THREADING_REG_EXP      = fmt.Sprintf("%s|%s", SYMBOL_THREADING_FIRST, SYMBOL_THREADING_LAST)

	SYMBOL_IF      = "if"
	SYMBOL_COND    = "cond"
	BRANCH_REG_EXP = fmt.Sprintf("\\b(%s|%s)\\b", SYMBOL_IF, SYMBOL_COND)

	SYMBOL_DEF      = "def" // NOTE: is core library in clojure
	SYMBOL_FN       = "fn"
	SYMBOL_LET      = "let"
	DECLARE_REG_EXP = fmt.Sprintf("\\b(%s|%s|%s)\\b", SYMBOL_DEF, SYMBOL_LET, SYMBOL_FN)

	LIB_CORE_PRN     = "prn"
	LIB_CORE_REG_EXP = fmt.Sprintf("\\b(%s)\\b", LIB_CORE_PRN)

	SYMBOL_USER_DEFINED_REG_EXP = `\w+`

	ALL_REG_EXP = fmt.Sprintf(
		"%s",
		strings.Join(
			[]string{
				FRACTIONAL_REG_EXP,
				INTEGER_REG_EXP,
				THREADING_REG_EXP,
				BRANCH_REG_EXP,
				OPERATORS_REG_EXP,
				BRACKETS_REG_EXP,
				LIB_CORE_PRN,
				SYMBOL_USER_DEFINED_REG_EXP,
			},
			"|",
		),
	)
)

func (tok *Token) IsKindAndVal(kind string, val string) bool {
	return tok.Kind == kind && tok.Val == val
}

func (tok *Token) IsNum() bool {
	return tok.Kind == TK_NUM
}

func (tok *Token) IsLibrary() bool {
	return tok.Kind == TK_LIB
}
func (tok *Token) IsReserved() bool {
	return tok.Kind == TK_RESERVED
}
func (tok *Token) IsDeclare() bool {
	return tok.Kind == TK_DECLARE
}
func (tok *Token) IsLambda() bool {
	return tok.Kind == TK_LAMBDA
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}
