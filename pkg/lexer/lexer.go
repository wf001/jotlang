package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wf001/modo/internal/log"
)

var FRACTIONAL_REG_EXP = `\d+\.\d+`
var INTEGER_REG_EXP = `\d+`
var BINARY_OPERATORS_REG_EXP = `[+\-*/<>]`
var PARREN_REG_EXP = `[()[\]{}"]`
var USER_DEFINED_REG_EXP = `\w+`

var THREADING_REG_EXP = `->|->>`
var BRANCH_REG_EXP = `\b(if|cond)\b`
var DEFINITION_REG_EXP = `\b(def|defn|let)\b`
var RESERVED_REG_EXP = strings.Join(
	[]string{
		THREADING_REG_EXP,
		BRANCH_REG_EXP,
    DEFINITION_REG_EXP,
	},
	"|")

var REG_EXP = fmt.Sprintf(
	"%s",
	strings.Join([]string{
		FRACTIONAL_REG_EXP,
		INTEGER_REG_EXP,
		RESERVED_REG_EXP,
		BINARY_OPERATORS_REG_EXP,
		PARREN_REG_EXP,
		USER_DEFINED_REG_EXP,
	}, "|"),
)

type TokenKind = int

const (
	TK_NUM TokenKind = iota + 1
	TK_OPERATOR
)

type ModoToken struct {
}

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}

func newToken(kind TokenKind, cur *Token, val string) *Token {
	tok := new(Token)
	tok.Kind = kind
	tok.Val = val
	cur.Next = tok
	return tok
}
func splitExpression(expr string) []string {
	re := regexp.MustCompile(REG_EXP)
	return re.FindAllString(expr, -1)
}

// return Token array from string
func Lex(s string) string {
	log.Debug(s)
	return s
}
