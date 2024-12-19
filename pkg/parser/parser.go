package parser

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

// TODO: unused
type parser struct {
	token *mTypes.Token
}

func newNodeParent(kind mTypes.NodeKind, child *mTypes.Node, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind:  kind,
		Child: child,
		Val:   val,
	}
}

func newNodeScalar(ty mTypes.ScalarType, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind: mTypes.ND_SCALAR,
		Type: ty,
		Val:  val,
	}
}

func matchedOperator(tok *mTypes.Token) (mTypes.NodeKind, bool) {
	var operatorTokenMap = map[string]mTypes.NodeKind{
		mTypes.NARY_OPERATOR_ADD:  mTypes.ND_ADD,
		mTypes.BINARY_OPERATOR_EQ: mTypes.ND_EQ,
	}

	if tok.Kind != mTypes.TK_OPERATOR {
		return "", false
	}

	if kind, exists := operatorTokenMap[tok.Val]; exists {
		return kind, true
	}

	return "", false
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
	return nextToken, argHead
}

func parseBody(
	rootToken *mTypes.Token,
	parentKind mTypes.NodeKind,
	exprName string,
) (*mTypes.Token, *mTypes.Node) {

	nextToken, argHead := parseExprs(rootToken.Next, parentKind)
	rootNode := newNodeParent(parentKind, argHead, exprName)
	return nextToken, rootNode
}

func parseDeclare(tok *mTypes.Token, parentKind mTypes.NodeKind) (*mTypes.Token, *mTypes.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &mTypes.Node{}

	// NOTE: too huge
	if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {
		tok = tok.Next

		if tok.IsDeclare() {
			tok, head = parseDeclare(tok.Next, mTypes.ND_DECLARE)
			head = newNodeParent(mTypes.ND_DECLARE, head, "")

		} else if tok.IsLambda() {
			tok, head = parseBody(tok.Next.Next, mTypes.ND_EXPR, "")
			head = newNodeParent(
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

			nextToken, body := parseBody(tok, mTypes.ND_EXPR, "")
			bind := newNodeParent(mTypes.ND_BIND, body, "")
			bind.Bind = varHead.Next

			return nextToken, bind

		} else if kind, isOperatorCall := matchedOperator(tok); isOperatorCall {
			log.Debug("is Operator :have %+v", tok)
			tok, head = parseBody(tok, kind, tok.Val)

		} else if tok.IsLibrary() {
			log.Debug("is Library :have %+v", tok)
			tok, head = parseBody(tok, mTypes.ND_LIBCALL, tok.Val)

		} else if tok.IsIf() {
			log.Debug("is If :have %+v", tok)
			head.Kind = mTypes.ND_IF

			nextToken, cond := parseDeclare(tok.Next, mTypes.ND_IF)
			head.Cond = cond
			head.Cond.Val = "cond"

			nextToken, then := parseDeclare(nextToken, mTypes.ND_IF)
			head.Then = newNodeParent(mTypes.ND_EXPR, then, "then")

			nextToken, els := parseDeclare(nextToken, mTypes.ND_IF)
			head.Else = newNodeParent(mTypes.ND_EXPR, els, "els")
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
			return tok, newNodeParent(mTypes.ND_VAR_DECLARE, head, v)

		} else {
			log.Debug("is Variable reference :have %+v", tok)
			return tok.Next, newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)

		}

	} else if tok.IsNum() {
		return tok.Next, newNodeScalar(mTypes.TY_INT32, tok.Val)

	} else if tok.IsStr() {
		return tok.Next, newNodeScalar(mTypes.TY_STR, tok.Val)

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
