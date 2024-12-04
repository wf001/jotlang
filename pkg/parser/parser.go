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
	nextToken, argHead := parseDeclare(rootToken.Next)
	prevNode := argHead
	for nextToken.IsNum() || nextToken.IsParenOpen() {
		nextToken, prevNode.Next = parseDeclare(nextToken)
		prevNode = prevNode.Next
	}
	rootNode = newNode(exprKind, argHead, exprName)
	rootToken = nextToken
	return rootToken, rootNode
}

func declare(
	rootToken *mTypes.Token,
	rootNode *mTypes.Node,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {
	nextToken, argHead := parseDeclare(rootToken.Next)

	rootNode = newNode(mTypes.ND_LAMBDA, argHead, exprName)
	rootToken = nextToken
	return rootToken, rootNode
}

// TODO: rename
func parseDeclare(tok *mTypes.Token) (*mTypes.Token, *mTypes.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &mTypes.Node{}

	if tok.Kind == mTypes.TK_EOL {
		return tok, nil
	}

	if tok.IsParenOpen() {
		tok = tok.Next

		// TODO: rename to matchedOperator
		if kind, isOperatorCall := matchedNodeKind(tok); isOperatorCall {
			log.Debug("is Operator :have %+v", tok)
			tok, head = expr(tok, head, kind, tok.Val)

		} else if tok.IsLibrary() {
			// TODO: put together?
			log.Debug("is Library :have %+v", tok)
			tok, head = expr(tok, head, mTypes.ND_LIBCALL, tok.Val)

		} else if tok.IsReserved() {
			tok, head = declare(tok, head, tok.Val)
		}

		if !tok.IsParenClose() {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.IsBracketOpen() {
		// skip
		return parseDeclare(tok.Next.Next)

	} else if tok.Kind == mTypes.TK_NUM {
		return tok.Next, newNodeInt(tok)

	} else {
		log.Debug("is Variable :have %+v", tok)
		tok, head = expr(tok, head, mTypes.ND_VAR, tok.Val)
	}

	return tok, head

}

// TODO: rename to parseProgram
func program(tok *mTypes.Token) *mTypes.Program {
	p := &mTypes.Program{}
	prevDeclare := p.Declares

	for tok != nil && tok.IsParenOpen() {

		if !tok.Next.IsReserved() {
			log.Panic("must be reserved token: have %+v", tok.Next)
		}

		if prevDeclare == nil {
			tok, p.Declares = parseDeclare(tok)
			prevDeclare = p.Declares
		} else {
			tok, prevDeclare.Next = parseDeclare(tok)
			prevDeclare = prevDeclare.Next
		}
	}
	return p
}

// take Token object, return Node object
func Parse(token *mTypes.Token) *mTypes.Program {
	prog := &mTypes.Program{}
	prog = program(token)
	prog.Debug(0)

	return prog
}
