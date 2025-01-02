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

	TK_TYPE_SIG   = TokenKind("TK_TYPE_SIG")
	TK_TYPE_ARROW = TokenKind("TK_TYPE_ARROW")
	TK_TYPE_INT   = TokenKind("TK_TYPE_INT")
	TK_TYPE_FLOAT = TokenKind("TK_TYPE_FLOAT")
	TK_TYPE_STR   = TokenKind("TK_TYPE_STR")
	TK_TYPE_NIL   = TokenKind("TK_TYPE_NIL")

	TK_INT   = TokenKind("TK_INT")
	TK_FLOAT = TokenKind("TK_FLOAT")
	TK_STR   = TokenKind("TK_STR")
	TK_NIL   = TokenKind("TK_NIL")
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

func (tok *Token) IsKindType() bool {
	return tok.IsKind(TK_TYPE_ARROW) ||
		tok.IsKind(TK_TYPE_INT) ||
		tok.IsKind(TK_TYPE_STR) ||
		tok.IsKind(TK_TYPE_NIL)
}

func (tok Token) MatchedType() (ModoType, bool) {
	var typeMap = map[string]ModoType{
		TK_TYPE_INT: TY_INT32,
		TK_TYPE_STR: TY_STR,
		TK_TYPE_NIL: TY_NIL,
	}

	if kind, exists := typeMap[tok.Kind]; exists {
		return kind, true
	}

	return "", false
}

func (tok *Token) MatchedNary() (NodeKind, bool) {
	var tokenMap = map[string]NodeKind{
		NARY_OPERATOR_ADD: ND_ADD,
	}

	return tok.MatchedOperator(tokenMap)
}

func (tok *Token) MatchedBinary() (NodeKind, bool) {
	var tokenMap = map[string]NodeKind{
		BINARY_OPERATOR_EQ: ND_EQ,
	}

	return tok.MatchedOperator(tokenMap)
}

func (tok *Token) MatchedOperator(

	tokenMap map[string]NodeKind,
) (NodeKind, bool) {

	if !tok.IsKind(TK_OPERATOR) {
		return "", false
	}

	if kind, exists := tokenMap[tok.Val]; exists {
		return kind, true
	}

	return "", false
}

func (tok *Token) DebugTokens() {
	log.Debug(log.BLUE("[token]"))
	for ; tok != nil; tok = tok.Next {
		log.Debug(log.BLUE(fmt.Sprintf("\t %p %#+v", tok, tok)))
	}
}
