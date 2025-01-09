package lexer

import (
	"regexp"
	"strings"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

// defines the structure for token patterns as a linked list
type tokenPattern struct {
	Pattern   string
	TokenType string
	Next      *tokenPattern
}

func isMatched(s string, typ string) bool {
	re := regexp.MustCompile(typ)
	return re.MatchString(s)
}

// traverses the linked list to find a matching pattern and returns the TokenType
func (tp *tokenPattern) matchTokenType(s string) (string, bool) {
	current := tp
	for current != nil {
		if isMatched(s, current.Pattern) {
			return current.TokenType, true
		}
		current = current.Next
	}
	return "", false
}

// initializes a TokenPattern with predefined patterns and token types
func newTokenPattern() *tokenPattern {
	head := &tokenPattern{
		Pattern:   mTypes.INTEGER_REG_EXP,
		TokenType: mTypes.TK_INT,
	}
	current := head

	add := func(pattern, tokenType string) {
		newNode := &tokenPattern{
			Pattern:   pattern,
			TokenType: tokenType,
		}
		current.Next = newNode
		current = newNode
	}

	add(mTypes.SYMBOL_TYPE_SIG, mTypes.TK_TYPE_SIG)
	add(mTypes.SYMBOL_TYPE_ARROW, mTypes.TK_TYPE_ARROW)
	add(mTypes.SYMBOL_TYPE_INT, mTypes.TK_TYPE_INT)
	add(mTypes.SYMBOL_TYPE_STR, mTypes.TK_TYPE_STR)
	add(mTypes.SYMBOL_TYPE_NIL, mTypes.TK_TYPE_NIL)

	add(mTypes.STRING_REG_EXP, mTypes.TK_STR)
	add(mTypes.OPERATORS_REG_EXP, mTypes.TK_LIBCALL)
	add(mTypes.BRACKETS_REG_EXP, mTypes.TK_PAREN)
	add(mTypes.LIB_CORE_REG_EXP, mTypes.TK_LIBCALL)
	add(mTypes.SYMBOL_FN, mTypes.TK_LAMBDA)
	add(mTypes.SYMBOL_DEF, mTypes.TK_DECLARE)
	add(mTypes.SYMBOL_LET, mTypes.TK_BIND)
	add(mTypes.SYMBOL_IF, mTypes.TK_IF)

	return head
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
	// Replace strings enclosed in double quotes with \"
	re := regexp.MustCompile(mTypes.STRING_REG_EXP)
	expr = re.ReplaceAllString(expr, `"$1"`)

	log.Debug(log.YELLOW("preprocessed %#+v"), expr)

	re = regexp.MustCompile(mTypes.ALL_REG_EXP)
	res := re.FindAllString(expr, -1)
	log.Debug(log.YELLOW("splitted program: %#+v"), res)
	return res
}

// performs lexical analysis on the input strings
func doLexicalAnalyse(splittedString []string) *mTypes.Token {
	prev := &mTypes.Token{}
	head := prev

	tokenMap := newTokenPattern()

	for _, p := range splittedString {
		if tokenType, matched := tokenMap.matchTokenType(p); matched {
			prev = newToken(tokenType, prev, p)
		} else {
			log.Debug("regard '%+v' as variable declaration or reference symbol", p)
			prev = newToken(mTypes.TK_IDENT, prev, p)
		}
	}
	trimQuote(head)

	head = head.Next
	head.DebugTokens()
	return head
}

// remove escaped \" from origin string
func trimQuote(head *mTypes.Token) {
	for t := head.Next; t != nil; t = t.Next {
		if t.IsKind(mTypes.TK_STR) {
			s := strings.Trim(t.Val, "\"")
			s = s + string([]byte{0})
			t.Val = s
		}
	}
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
