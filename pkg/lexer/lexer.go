package lexer

import (
	"regexp"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func isMatched(s string, typ string) bool {
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

var tokenMap = map[string]string{
	mTypes.INTEGER_REG_EXP:   mTypes.TK_NUM,
	mTypes.OPERATORS_REG_EXP: mTypes.TK_OPERATOR,
	mTypes.BRACKETS_REG_EXP:  mTypes.TK_PAREN,
	mTypes.LIB_CORE_REG_EXP:  mTypes.TK_LIBCALL,
	mTypes.SYMBOL_FN:         mTypes.TK_LAMBDA,
	mTypes.SYMBOL_DEF:        mTypes.TK_DECLARE,
	mTypes.SYMBOL_LET:        mTypes.TK_BIND,
	mTypes.SYMBOL_IF:         mTypes.TK_IF,
}

func doLexicalAnalyse(splittedString []string) *mTypes.Token {
	prev := &mTypes.Token{}
	head := prev

	for _, p := range splittedString {
		matched := false

		for pattern, tokenType := range tokenMap {
			if isMatched(p, pattern) {
				prev = newToken(tokenType, prev, p)
				matched = true
				break
			}
		}
		if !matched {
			log.Debug("regard '%+v' as variable declaration or reference symbol", p)
			prev = newToken(mTypes.TK_IDENT, prev, p)
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
