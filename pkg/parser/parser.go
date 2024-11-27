package parser

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

type parser struct {
	token *mTypes.Token
}

func newNode(kind mTypes.NodeKind, child *mTypes.Node, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind:  kind,
		Child: child,
		Val:   val,
	}
}

func newNodeInt(tok *mTypes.Token) *mTypes.Node {
	return newNode(mTypes.ND_INT, nil, tok.Val)
}

func matchedNodeKind(tok *mTypes.Token) (mTypes.NodeKind, bool) {
	if tok.Kind != mTypes.TK_OPERATOR {
		return mTypes.ND_NIL, false
	}
	switch tok.Val {
	case mTypes.NARY_OPERATOR_ADD:
		return mTypes.ND_ADD, true
	case mTypes.BINARY_OPERATOR_EQ:
		return mTypes.ND_EQ, true
	}
	return mTypes.ND_NIL, false
}

func expr(
	rootToken *mTypes.Token,
	rootNode *mTypes.Node,
	exprKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {
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

func program(tok *mTypes.Token) (*mTypes.Token, *mTypes.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &mTypes.Node{}

	if tok.Kind == mTypes.TK_EOL {
		return tok, nil
	}

	if tok.IsParenOpen() {
		tok = tok.Next

		if kind, isOperatorCall := matchedNodeKind(tok); isOperatorCall {
			tok, head = expr(tok, head, kind, tok.Val)
		} else if tok.IsLibrary() {
			// TODO: put together?
			tok, head = expr(tok, head, mTypes.ND_LIB, tok.Val)
		}

		if !tok.IsParenClose() {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.Kind == mTypes.TK_NUM {
		return tok.Next, newNodeInt(tok)

	}
	return tok, head
}

func Construct(token *mTypes.Token) *parser {
	return &parser{
		token: token,
	}
}

// take Token object, return Node object
func (p parser) Parse() *mTypes.Program {
	prog := &mTypes.Program{}
	_, prog.FuncCalls = program(p.token)
	prog.Debug(0)

	return prog
}
