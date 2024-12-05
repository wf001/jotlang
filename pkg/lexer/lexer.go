package lexer

import (
	"regexp"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func matched(s string, typ string) bool {
	re := regexp.MustCompile(typ)
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

func splitString(expr string) []string {
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
		if matched(p, mTypes.INTEGER_REG_EXP) {
			next = newToken(mTypes.TK_NUM, cur, p)
			cur = next
		} else if matched(p, mTypes.OPERATORS_REG_EXP) {
			next = newToken(mTypes.TK_OPERATOR, cur, p)
			cur = next
		} else if matched(p, mTypes.BRACKETS_REG_EXP) {
			next = newToken(mTypes.TK_PAREN, cur, p)
			cur = next
		} else if matched(p, mTypes.LIB_CORE_REG_EXP) {
			next = newToken(mTypes.TK_LIB, cur, p)
			cur = next
		} else if matched(p, mTypes.SYMBOL_FN) {
			next = newToken(mTypes.TK_LAMBDA, cur, p)
			cur = next
		} else if matched(p, mTypes.SYMBOL_DEF) {
			next = newToken(mTypes.TK_DECLARE, cur, p)
			cur = next
		} else {
			log.Debug("regard '%+v' as user defined symbol", p)
			next = newToken(mTypes.TK_VARIABLE, cur, p)
			cur = next
		}
	}
	head = head.Next
	head.DebugTokens()
	return head
}

// take string, return Token object
func Lex(s string) *mTypes.Token {
	log.Debug(log.YELLOW("original source: '%s'"), s)
	log.DebugMessage("code lexing")
	arr := splitString(s)
	res := doLexicalAnalyse(arr)
	log.DebugMessage("code lexed")
	return res
}
