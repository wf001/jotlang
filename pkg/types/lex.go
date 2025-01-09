package types

import (
	"fmt"
	"strings"
)

var (
	// Scholar
	FRACTIONAL_REG_EXP = `-?\d+\.\d+`
	INTEGER_REG_EXP    = `-?\d+`
	STRING_REG_EXP     = `"([^"]*)"`

	// Brachets
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

	// operators
	OPERATOR_ADD = "+"
	OPERATOR_SUB = "-"
	OPERATOR_MUL = "*"
	OPERATOR_DIV = "/"

	OPERATOR_EQ = "="
	OPERATOR_LT = "<"
	OPERATOR_GT = ">"

	OPERATORS_REG_EXP = fmt.Sprintf(
		"[%s\\%s%s%s%s%s%s]",
		OPERATOR_ADD,
		OPERATOR_SUB,
		OPERATOR_MUL,
		OPERATOR_DIV,
		OPERATOR_EQ,
		OPERATOR_LT,
		OPERATOR_GT,
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

	// Core library
	LIB_CORE_PRN     = "prn"
	LIB_CORE_REG_EXP = fmt.Sprintf("\\b(%s)\\b", LIB_CORE_PRN)

	// Type signature
	SYMBOL_TYPE_SIG   = "::"
	SYMBOL_TYPE_ARROW = "=>"
	SYMBOL_TYPE_INT   = "int"
	SYMBOL_TYPE_STR   = "string"
	SYMBOL_TYPE_NIL   = "nil"

	TYPE_REG_EXP = fmt.Sprintf(
		"(%s|%s|\\b(%s|%s|%s)\\b)",
		SYMBOL_TYPE_SIG,
		SYMBOL_TYPE_ARROW,
		SYMBOL_TYPE_INT,
		SYMBOL_TYPE_STR,
		SYMBOL_TYPE_NIL,
	)

	SYMBOL_UNDEFINED_REG_EXP = `\w+`

	ALL_REG_EXP = fmt.Sprintf(
		"%s",
		strings.Join(
			[]string{
				FRACTIONAL_REG_EXP,
				INTEGER_REG_EXP,
				STRING_REG_EXP,
				TYPE_REG_EXP,
				THREADING_REG_EXP,
				BRANCH_REG_EXP,
				OPERATORS_REG_EXP,
				BRACKETS_REG_EXP,
				LIB_CORE_PRN,
				SYMBOL_UNDEFINED_REG_EXP,
			},
			"|",
		),
	)
)
