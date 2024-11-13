package types

type TokenKind = string

const (
	TK_NUM      = TokenKind("TK_NUM")
	TK_OPERATOR = TokenKind("TK_OPERATOR")
	TK_PAREN    = TokenKind("TK_PAREN")
	TK_EOL      = TokenKind("TK_EOL")
)

type Token struct {
	Kind TokenKind
	Next *Token
	Val  string
}

type NodeKind string

const (
	ND_ADD = NodeKind("ND_KIND_ADD") // +
	ND_SUB = NodeKind("ND_SUB")      // -
	ND_MUL = NodeKind("ND_MUL")      // *
	ND_DIV = NodeKind("ND_DIV")      // /
	ND_INT = NodeKind("ND_INT")      // /
)

type Node struct {
	Kind  NodeKind
	Next  *Node
	Child *Node
	Cond  *Node
	Then  *Node
	Else  *Node
	Init  *Node
	Inc   *Node
	Body  *Node
	Func  string
	Args  *Node
	Val   string
}
