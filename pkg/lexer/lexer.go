package lexer

import (
	"regexp"

	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/types"
)

func isInteger(s string) bool {
	re := regexp.MustCompile(types.INTEGER_REG_EXP)
	return re.MatchString(s)
}

func isParen(s string) bool {
	re := regexp.MustCompile(types.BRACKETS_REG_EXP)
	return re.MatchString(s)
}

func isBinaryOperator(s string) bool {
	re := regexp.MustCompile(types.OPERATORS_REG_EXP)
	return re.MatchString(s)
}

func newToken(kind types.TokenKind, cur *types.Token, val string) *types.Token {
	tok := &types.Token{
		Kind: kind,
		Val:  val,
	}
	cur.Next = tok
	return tok
}

func splitProgram(expr string) []string {
	re := regexp.MustCompile(types.ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted program: %#+v"), res)
	return res
}

func doLexicalAnalyse(splittedProgram []string) *types.Token {
	cur := new(types.Token)
	head := cur
	next := new(types.Token)

	for _, p := range splittedProgram {
		if isInteger(p) {
			next = newToken(types.TK_NUM, cur, p)
			cur = next
		} else if isBinaryOperator(p) {
			next = newToken(types.TK_OPERATOR, cur, p)
			cur = next
		} else if isParen(p) {
			next = newToken(types.TK_PAREN, cur, p)
			cur = next
		} else {
			log.Panic("include invalid signature: have '%+v'", p)
		}
	}
	head = head.Next
	head.DebugTokens()
	return head
}

// take string, return Token object
func Lex(s string) *types.Token {
	log.Debug(log.YELLOW("original source: '%s'"), s)
	arr := splitProgram(s)
	return doLexicalAnalyse(arr)
}
