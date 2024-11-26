package parser

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/types"
)

type parser struct {
	token *types.Token
}

func newNode(kind types.NodeKind, child *types.Node, val string) *types.Node {
	return &types.Node{
		Kind:  kind,
		Child: child,
		Val:   val,
	}
}

func newNodeInt(tok *types.Token) *types.Node {
	return newNode(types.ND_INT, nil, tok.Val)
}

func matchedNodeKind(tok *types.Token) (types.NodeKind, bool) {
	if tok.Kind != types.TK_OPERATOR {
		return types.ND_NIL, false
	}
	switch tok.Val {
	case types.NARY_OPERATOR_ADD:
		return types.ND_ADD, true
	case types.BINARY_OPERATOR_EQ:
		return types.ND_EQ, true
	}
	return types.ND_NIL, false
}

func expr(tok *types.Token, head *types.Node, kind types.NodeKind) (*types.Token, *types.Node) {
	nextToken, argHead := program(tok.Next)
	prevNode := argHead
	for nextToken.IsNum() || nextToken.IsParenOpen() {
		nextToken, prevNode.Next = program(nextToken)
		prevNode = prevNode.Next
	}
	head = newNode(kind, argHead, "")
	tok = nextToken
	return tok, head
}

func funcCall(
	tok *types.Token,
	head *types.Node,
	kind types.NodeKind,
	val string,
) (*types.Token, *types.Node) {
	nextToken, argHead := program(tok.Next)
	prevNode := argHead
	for nextToken.IsNum() || nextToken.IsParenOpen() {
		nextToken, prevNode.Next = program(nextToken)
		prevNode = prevNode.Next
	}
	head = newNode(kind, argHead, val)
	tok = nextToken
	return tok, head
}

func program(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &types.Node{}

	if tok.Kind == types.TK_EOL {
		return tok, nil
	}

	if tok.IsParenOpen() {
		tok = tok.Next

		if kind, isExprCall := matchedNodeKind(tok); isExprCall {
			tok, head = expr(tok, head, kind)
		} else if tok.IsLibrary() {
			nextToken, nextNode := funcCall(tok, head, types.ND_LIB, tok.Val)
			tok, head = nextToken, nextNode
		}

		if !tok.IsParenClose() {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.Kind == types.TK_NUM {
		return tok.Next, newNodeInt(tok)

	}
	return tok, head
}

func Construct(token *types.Token) *parser {
	return &parser{
		token: token,
	}
}

// take Token object, return Node object
func (p parser) Parse() *types.Node {
	_, node := program(p.token)
	node.DebugNode(0)

	return node
}
