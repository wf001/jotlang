package lexer

import (
	"regexp"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func isInteger(s string) bool {
	re := regexp.MustCompile(mTypes.INTEGER_REG_EXP)
	return re.MatchString(s)
}

func isParen(s string) bool {
	re := regexp.MustCompile(mTypes.BRACKETS_REG_EXP)
	return re.MatchString(s)
}

func isBinaryOperator(s string) bool {
	re := regexp.MustCompile(mTypes.OPERATORS_REG_EXP)
	return re.MatchString(s)
}
func isLibCore(s string) bool {
	re := regexp.MustCompile(mTypes.LIB_CORE_REG_EXP)
	return re.MatchString(s)
}
func isReserved(s string) bool {
	re := regexp.MustCompile(mTypes.RESERVED_REG_EXP)
	return re.MatchString(s)
}

func newToken(kind mTypes.TokenKind, cur *mTypes.Token, val string) *mTypes.Token {
	tok := &mTypes.Token{
		Kind: kind,
		Val:  val,
	}
	cur.Next = tok
	return tok
}

func splitProgram(expr string) []string {
	re := regexp.MustCompile(mTypes.ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted program: %#+v"), res)
	return res
}

func doLexicalAnalyse(splittedProgram []string) *mTypes.Token {
	cur := new(mTypes.Token)
	head := cur
	next := new(mTypes.Token)

	for _, p := range splittedProgram {
		if isInteger(p) {
			next = newToken(mTypes.TK_NUM, cur, p)
			cur = next
		} else if isBinaryOperator(p) {
			next = newToken(mTypes.TK_OPERATOR, cur, p)
			cur = next
		} else if isParen(p) {
			next = newToken(mTypes.TK_PAREN, cur, p)
			cur = next
		} else if isLibCore(p) {
			next = newToken(mTypes.TK_LIB, cur, p)
			cur = next
		} else if isReserved(p) {
			next = newToken(mTypes.TK_RESERVED, cur, p)
			cur = next
		} else {
			log.Debug("use '%+v' as user defined variable", p)
		}
	}
	head = head.Next
	head.DebugTokens()
	return head
}

// take string, return Token object
func Lex(s string) *mTypes.Token {
	log.Debug(log.YELLOW("original source: '%s'"), s)
	arr := splitProgram(s)
	return doLexicalAnalyse(arr)
}
