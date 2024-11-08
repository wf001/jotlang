package lexer

import (
	"regexp"

	"github.com/wf001/modo/internal/log"
)

var OPERATORS_REG_EXP = `^[+\-*/]$`
var NUMBER_REG_EXP = `^[+-]?\d+$`

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
func parseIf(s string, expr string) string {
	re := regexp.MustCompile(expr)
	if re.MatchString(s) {
		return s
	}
	return ""
}

func parseIfOperator(c string) string {
	return parseIf(c, OPERATORS_REG_EXP)
}

func parseIfNumber(c string) string {
	return parseIf(c, NUMBER_REG_EXP)
}

// return Token array from string
func Lex(s string) string {
	log.Debug(s)
	return s
}
