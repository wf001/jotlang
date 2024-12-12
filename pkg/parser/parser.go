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

func parseExprs(
	rootToken *mTypes.Token,
	exprKind mTypes.NodeKind,
) (*mTypes.Token, *mTypes.Node) {

	nextToken, argHead := parseDeclare(rootToken, exprKind)
	prevNode := argHead
	// NOTE: what means?
	for nextToken.IsNum() ||
		nextToken.IsVar() ||
		nextToken.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {

		nextToken, prevNode.Next = parseDeclare(nextToken, exprKind)
		prevNode = prevNode.Next
	}
	rootToken = nextToken
	return rootToken, argHead
}

func parseBody(
	rootToken *mTypes.Token,
	rootNode *mTypes.Node, // TODO: unused?
	exprKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {

	nextToken, argHead := parseExprs(rootToken.Next, exprKind)
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
			tok, head = parseBody(tok.Next.Next, head, mTypes.ND_EXPR, "")
			head = newNode(
				mTypes.ND_LAMBDA,
				head,
				"",
			)

		} else if tok.IsBind() {
			tok = tok.Next

			if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
				log.Panic("let bindings require even number of forms")
			}

			prev := &mTypes.Node{}
			varHead := prev

			tok = tok.Next
			for {

				nextToken, nodeVar := parseDeclare(tok, mTypes.ND_DECLARE)
				tok = nextToken

				prev.Next = nodeVar
				prev = prev.Next

				if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
					break
				}
				if tok == nil {
					log.Panic("must be closed with ]")
				}
			}

			nextToken, body := parseBody(tok, head, mTypes.ND_EXPR, "")
			bind := newNode(mTypes.ND_BIND, body, "")
			bind.Bind = varHead.Next

			return nextToken, bind

		} else if kind, isOperatorCall := matchedOperator(tok); isOperatorCall {
			log.Debug("is Operator :have %+v", tok)
			tok, head = parseBody(tok, head, kind, tok.Val)

		} else if tok.IsLibrary() {
			log.Debug("is Library :have %+v", tok)
			tok, head = parseBody(tok, head, mTypes.ND_LIBCALL, tok.Val)

		} else if tok.IsIf() {
			log.Debug("is If :have %+v", tok)
			head.Kind = mTypes.ND_IF

			nextToken, cond := parseDeclare(tok.Next, mTypes.ND_IF)
			head.Cond = cond

			// FIXME: fail to parse (if (= 1 1) ((prn 2) (prn 3)) (prn 4))
			nextToken, then := parseDeclare(nextToken, mTypes.ND_IF)
			head.Then = newNode(mTypes.ND_EXPR, then, "")

			nextToken, els := parseDeclare(nextToken, mTypes.ND_IF)
			head.Else = newNode(mTypes.ND_EXPR, els, "")
			tok = nextToken
		}

		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_CLOSE) {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.IsVar() {
		if parentKind == mTypes.ND_DECLARE {
			log.Debug("is Variable declaration :have %+v", tok)
			v := tok.Val
			tok, head = parseDeclare(tok.Next, mTypes.ND_VAR_DECLARE)
			return tok, newNode(mTypes.ND_VAR_DECLARE, head, v)

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
	prog := parseProgram(token)
	log.DebugMessage("code parsed")

	prog.Debug(0)

	return prog
}
