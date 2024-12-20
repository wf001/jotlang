package types

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
)

type TokenKind = string

const (
	TK_DECLARE = TokenKind("TK_DECLARE")
	TK_LAMBDA  = TokenKind("TK_LAMBDA")
	TK_BIND    = TokenKind("TK_BIND")
	TK_IF      = TokenKind("TK_IF")

	TK_PAREN = TokenKind("TK_PAREN")

	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_LIBCALL  = TokenKind("TK_LIBCALL")
	TK_IDENT    = TokenKind("TK_IDENT")

	TK_INT   = TokenKind("TK_INT")
	TK_FLOAT = TokenKind("TK_FLOAT")
	TK_STR   = TokenKind("TK_STR")
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}

func (tok *Token) IsKindAndVal(kind string, val string) bool {
	return tok != nil && tok.IsKind(kind) && tok.Val == val
}

func (tok *Token) IsKind(kind TokenKind) bool {
	return tok.Kind == kind
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}
