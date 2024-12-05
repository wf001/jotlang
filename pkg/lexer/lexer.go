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

func newToken(kind mTypes.TokenKind, prev *mTypes.Token, val string) *mTypes.Token {
	tok := &mTypes.Token{
		Kind: kind,
		Val:  val,
	}
	prev.Next = tok
	return tok
}

func splitString(expr string) []string {
	re := regexp.MustCompile(mTypes.ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted program: %#+v"), res)
	return res
}

func doLexicalAnalyse(splittedString []string) *mTypes.Token {
	prev := &mTypes.Token{}
	head := prev
	next := &mTypes.Token{}

	for _, p := range splittedString {
		// NOTE: should implement as map refer?
		if matched(p, mTypes.INTEGER_REG_EXP) {
			next = newToken(mTypes.TK_NUM, prev, p)
			prev = next
		} else if matched(p, mTypes.OPERATORS_REG_EXP) {
			next = newToken(mTypes.TK_OPERATOR, prev, p)
			prev = next
		} else if matched(p, mTypes.BRACKETS_REG_EXP) {
			next = newToken(mTypes.TK_PAREN, prev, p)
			prev = next
		} else if matched(p, mTypes.LIB_CORE_REG_EXP) {
			next = newToken(mTypes.TK_LIB, prev, p)
			prev = next
		} else if matched(p, mTypes.SYMBOL_FN) {
			next = newToken(mTypes.TK_LAMBDA, prev, p)
			prev = next
		} else if matched(p, mTypes.SYMBOL_DEF) {
			next = newToken(mTypes.TK_DECLARE, prev, p)
			prev = next
		} else {
			if prev.Kind == mTypes.TK_DECLARE {
				log.Debug("regard '%+v' as variable declaration symbol", p)
				next = newToken(mTypes.TK_VAR_DECLARE, prev, p)
				prev = next
			} else {
				log.Debug("regard '%+v' as variable reference symbol", p)
				next = newToken(mTypes.TK_VAR_REFERENCE, prev, p)
				prev = next
			}
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
