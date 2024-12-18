package types

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
)

type TokenKind = string

const (
	TK_NUM      = TokenKind("TK_NUM")
	TK_STR      = TokenKind("TK_STR")
	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_PAREN    = TokenKind("TK_PAREN")
	TK_LIBCALL  = TokenKind("TK_LIBCALL")

	TK_RESERVED = TokenKind("TK_RESERVED")
	TK_DECLARE  = TokenKind("TK_DECLARE")
	TK_LAMBDA   = TokenKind("TK_LAMBDA")
	TK_BIND     = TokenKind("TK_BIND")

	TK_IF = TokenKind("TK_IF")

	TK_IDENT = TokenKind("TK_IDENT")
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}

func (tok *Token) IsKindAndVal(kind string, val string) bool {
	return tok != nil && tok.Kind == kind && tok.Val == val
}

func (tok *Token) IsNum() bool {
	return tok.Kind == TK_NUM
}

func (tok *Token) IsLibrary() bool {
	return tok.Kind == TK_LIBCALL
}
func (tok *Token) IsReserved() bool {
	return tok.Kind == TK_RESERVED
}
func (tok *Token) IsDeclare() bool {
	return tok.Kind == TK_DECLARE
}
func (tok *Token) IsLambda() bool {
	return tok.Kind == TK_LAMBDA
}
func (tok *Token) IsVar() bool {
	return tok.Kind == TK_IDENT
}
func (tok *Token) IsBind() bool {
	return tok.Kind == TK_BIND
}
func (tok *Token) IsIf() bool {
	return tok.Kind == TK_IF
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}
