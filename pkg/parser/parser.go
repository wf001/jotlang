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
	default:
		log.Error("unresolved token :have %+v", tok)
	}
	return mTypes.ND_NIL, false
}

func parseExpr(
	rootToken *mTypes.Token,
	rootNode *mTypes.Node,
	exprKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {
	nextToken, argHead := parseDeclare(rootToken.Next, exprKind)
	prevNode := argHead
	// NOTE: what means?
	for nextToken.IsNum() || nextToken.IsVar() || nextToken.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		nextToken, prevNode.Next = parseDeclare(nextToken, exprKind)
		prevNode = prevNode.Next
	}
	rootNode = newNode(exprKind, argHead, exprName)
	rootToken = nextToken
	return rootToken, rootNode
}

func parseDeclare(tok *mTypes.Token, parentKind mTypes.NodeKind) (*mTypes.Token, *mTypes.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &mTypes.Node{}

	if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		tok = tok.Next

		if tok.IsDeclare() {
			tok, head = parseDeclare(tok.Next, mTypes.ND_DECLARE)
			head = newNode(mTypes.ND_DECLARE, head, "")

		} else if tok.IsLambda() {
			tok, head = parseDeclare(tok.Next.Next.Next, mTypes.ND_LAMBDA)
			head = newNode(
				mTypes.ND_LAMBDA,
				head,
				"",
			)
		} else if tok.IsBind() {
			tok = tok.Next

			if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_OPEN) {
				log.Panic("must be [ :have %+v", tok)
			}

			prev := &mTypes.Node{}
			varHead := &mTypes.Node{}

			for {
				if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) ||
					tok.Next.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
					log.Panic("bindings require even number of forms")
				}
				tok = tok.Next

				tok, prev = parseDeclare(tok, mTypes.ND_DECLARE)

				if varHead.Child == nil {
					varHead.Child = prev
				} else {
					prev.Next = prev
					prev = prev.Next
				}

				if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
					tok = tok.Next
					nextToken, expr := parseDeclare(tok, mTypes.ND_EXPR)
					res := newNode(mTypes.ND_EXPR, expr, "")
					res.Bind = newNode(mTypes.ND_BIND, varHead.Child, "")
					tok = nextToken
					return nextToken, res
				}
				if tok == nil {
					log.Panic("must be closed with ]")
				}
			}
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

	} else if tok.IsVar() {
		if parentKind == mTypes.ND_DECLARE {
			log.Debug("is Variable declaration :have %+v", tok)
			tok, head = parseExpr(tok, head, mTypes.ND_VAR_DECLARE, tok.Val)

		} else {

			log.Debug("is Variable reference :have %+v", tok)
			return tok.Next, newNode(mTypes.ND_VAR_REFERENCE, nil, tok.Val)
		}

	} else if tok.IsNum() {
		return tok.Next, newNodeInt(tok)

	} else {
		log.Panic("unresolved token :have %+v", tok)
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
			tok, p.Declares = parseDeclare(tok, mTypes.ND_PROGRAM_ROOT)
			prevDeclare = p.Declares
		} else {
			tok, prevDeclare.Next = parseDeclare(tok, mTypes.ND_PROGRAM_ROOT)
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
