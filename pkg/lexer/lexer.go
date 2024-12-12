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

	for _, p := range splittedString {
		// NOTE: should implement as map refer?
		if matched(p, mTypes.INTEGER_REG_EXP) {
			prev = newToken(mTypes.TK_NUM, prev, p)

		} else if matched(p, mTypes.OPERATORS_REG_EXP) {
			prev = newToken(mTypes.TK_OPERATOR, prev, p)

		} else if matched(p, mTypes.BRACKETS_REG_EXP) {
			prev = newToken(mTypes.TK_PAREN, prev, p)

		} else if matched(p, mTypes.LIB_CORE_REG_EXP) {
			prev = newToken(mTypes.TK_LIBCALL, prev, p)

		} else if matched(p, mTypes.SYMBOL_FN) {
			prev = newToken(mTypes.TK_LAMBDA, prev, p)

		} else if matched(p, mTypes.SYMBOL_DEF) {
			prev = newToken(mTypes.TK_DECLARE, prev, p)

		} else if matched(p, mTypes.SYMBOL_LET) {
			prev = newToken(mTypes.TK_BIND, prev, p)

		} else {
			log.Debug("regard '%+v' as variable declaration or reference symbol", p)
			prev = newToken(mTypes.TK_VAR, prev, p)
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
