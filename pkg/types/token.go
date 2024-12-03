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
	TK_RESERVED = TokenKind("TK_RESERVED")
	TK_LIB      = TokenKind("TK_LIB")
	TK_EOL      = TokenKind("TK_EOL")
	TK_VARIABLE = TokenKind("TK_VARIABLE")
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
	EXPR_FN            = "fn"
	EXPR_DEFN          = "defn"
	EXPR_LET           = "let"
	DEFINITION_REG_EXP = fmt.Sprintf("\\b(%s|%s|%s|%s)\\b", EXPR_DEF, EXPR_DEFN, EXPR_LET, EXPR_FN)

	SIGNATURE_REG_EXP = strings.Join(
		[]string{
			THREADING_REG_EXP,
			BRANCH_REG_EXP,
		},
		"|",
	)
	RESERVED_REG_EXP = strings.Join(
		[]string{
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
				SIGNATURE_REG_EXP,
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

// TODO: define isKindAndVal
func (tok *Token) IsParenOpen() bool {
	return tok.Kind == TK_PAREN && tok.Val == PARREN_OPEN
}

func (tok *Token) IsParenClose() bool {
	return tok.Kind == TK_PAREN && tok.Val == PARREN_CLOSE
}
func (tok *Token) IsBracketOpen() bool {
	return tok.Kind == TK_PAREN && tok.Val == BRACKET_OPEN
}

func (tok *Token) IsBracketClose() bool {
	return tok.Kind == TK_PAREN && tok.Val == BRACKET_CLOSE
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
func (tok *Token) IsReserved() bool {
	return tok.Kind == TK_RESERVED
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}
