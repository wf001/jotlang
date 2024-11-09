package types

type TokenKind = int

const (
	TK_NUM TokenKind = iota + 1
	TK_OPERATOR
	TK_EOL
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}
