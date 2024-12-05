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

func matchedOperator(tok *mTypes.Token) (mTypes.NodeKind, bool) {
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

func parseExpr(
	rootToken *mTypes.Token,
	rootNode *mTypes.Node,
	exprKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {
	nextToken, argHead := parseDeclare(rootToken.Next)
	prevNode := argHead
	for nextToken.IsNum() || nextToken.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		nextToken, prevNode.Next = parseDeclare(nextToken)
		prevNode = prevNode.Next
	}
	rootNode = newNode(exprKind, argHead, exprName)
	rootToken = nextToken
	return rootToken, rootNode
}

func parseDeclare(tok *mTypes.Token) (*mTypes.Token, *mTypes.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &mTypes.Node{}

	if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		tok = tok.Next

		if tok.IsDeclare() {
			tok, head = parseDeclare(tok.Next)
			head = newNode(mTypes.ND_DECLARE, head, "")
		} else if tok.IsLambda() {
			tok, head = parseDeclare(tok.Next)
			head = newNode(
				mTypes.ND_LAMBDA,
				newNode(
					mTypes.ND_EXPR,
					head,
					"",
				),
				"",
			)

		} else if kind, isOperatorCall := matchedOperator(tok); isOperatorCall {
			log.Debug("is Operator :have %+v", tok)
			tok, head = parseExpr(tok, head, kind, tok.Val)

		} else if tok.IsLibrary() {
			log.Debug("is Library :have %+v", tok)
			tok, head = parseExpr(tok, head, mTypes.ND_LIBCALL, tok.Val)
		}

		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_CLOSE) {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_OPEN) {
		// skip
		return parseDeclare(tok.Next.Next)

	} else if tok.Kind == mTypes.TK_NUM {
		return tok.Next, newNodeInt(tok)

	} else if tok.Kind == mTypes.TK_VAR_DECLARE {
		log.Debug("is Variable declaration :have %+v", tok)
		tok, head = parseExpr(tok, head, mTypes.ND_VAR_DECLARE, tok.Val)

	} else if tok.Kind == mTypes.TK_VAR_REFERENCE {
		log.Debug("is Variable reference :have %+v", tok)
		return tok.Next, newNode(mTypes.ND_VAR_REFERENCE, nil, tok.Val)

	} else {
		log.Panic("undefined token :have %+v", tok)
	}

	return tok, head

}

func parseProgram(tok *mTypes.Token) *mTypes.Program {
	p := &mTypes.Program{}
	prevDeclare := p.Declares

	for tok != nil && tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {

		if !tok.Next.IsDeclare() {
			log.Panic("must be declare token: have %+v", tok.Next)
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

// take Token object, return Program object
func Parse(token *mTypes.Token) *mTypes.Program {
	log.DebugMessage("code parsing")
	prog := &mTypes.Program{}
	prog = parseProgram(token)
	log.DebugMessage("code parsed")
	prog.Debug(0)

	return prog
}
