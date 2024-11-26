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

func expr(
	rootToken *types.Token,
	rootNode *types.Node,
	exprKind types.NodeKind,
	exprName string,
) (*types.Token, *types.Node) {
	nextToken, argHead := program(rootToken.Next)
	prevNode := argHead
	for nextToken.IsNum() || nextToken.IsParenOpen() {
		nextToken, prevNode.Next = program(nextToken)
		prevNode = prevNode.Next
	}
	rootNode = newNode(exprKind, argHead, exprName)
	rootToken = nextToken
	return rootToken, rootNode
}

func program(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &types.Node{}

	if tok.Kind == types.TK_EOL {
		return tok, nil
	}

	if tok.IsParenOpen() {
		tok = tok.Next

		if kind, isOperatorCall := matchedNodeKind(tok); isOperatorCall {
			tok, head = expr(tok, head, kind, tok.Val)
		} else if tok.IsLibrary() {
			// TODO: put together?
			tok, head = expr(tok, head, types.ND_LIB, tok.Val)
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
func (p parser) Parse() *types.Program {
	prog := &types.Program{}
	_, prog.FuncCalls = program(p.token)
	prog.Debug(0)

	return prog
}
