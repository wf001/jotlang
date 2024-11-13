package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/types"
)

var FRACTIONAL_REG_EXP = `\d+\.\d+`
var INTEGER_REG_EXP = `\d+`
var BINARY_OPERATORS_REG_EXP = `[+\-*/<>]`
var PARREN_REG_EXP = `[()[\]{}"]` // paren, bracket, brace
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
	"|",
)

var ALL_REG_EXP = fmt.Sprintf(
	"%s",
	strings.Join(
		[]string{
			FRACTIONAL_REG_EXP,
			INTEGER_REG_EXP,
			RESERVED_REG_EXP,
			BINARY_OPERATORS_REG_EXP,
			PARREN_REG_EXP,
			USER_DEFINED_REG_EXP,
		},
		"|",
	),
)

func isInteger(s string) bool {
	re := regexp.MustCompile(INTEGER_REG_EXP)
	return re.MatchString(s)
}
func isParen(s string) bool {
	re := regexp.MustCompile(PARREN_REG_EXP)
	return re.MatchString(s)
}
func isBinaryOperator(s string) bool {
	re := regexp.MustCompile(BINARY_OPERATORS_REG_EXP)
	return re.MatchString(s)
}

func newToken(kind types.TokenKind, cur *types.Token, val string) *types.Token {
	tok := new(types.Token)
	tok.Kind = kind
	tok.Val = val
	cur.Next = tok
	return tok
}
func splitExpression(expr string) []string {
	re := regexp.MustCompile(ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted expr: %#+v"), res)
	return res
}

func doLexicalAnalyse(expr []string) *types.Token {
	cur := new(types.Token)
	head := cur
	next := new(types.Token)

	for _, e := range expr {
		if isInteger(e) {
			next = newToken(types.TK_NUM, cur, e)
			cur = next
		} else if isBinaryOperator(e) {
			next = newToken(types.TK_OPERATOR, cur, e)
			cur = next
		} else if isParen(e) {
			next = newToken(types.TK_PAREN, cur, e)
			cur = next
		} else {
			log.Panic("include invalid signature: have '%+v'", e)
		}
	}
	head = head.Next
	log.DebugTokens(head)
	return head
}

// return Token array from string
func Lex(s string) *types.Token {
	log.Debug(log.YELLOW("original source: '%s'"), s)
	arr := splitExpression(s)
	return doLexicalAnalyse(arr)
}
