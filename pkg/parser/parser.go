package parser

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
)

func newNodeParent(kind mTypes.NodeKind, child *mTypes.Node, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind:  kind,
		Child: child,
		Val:   val,
	}
}

func newNodeScalar(ty mTypes.ModoType, val string) *mTypes.Node {
	return &mTypes.Node{
		Kind: mTypes.ND_SCALAR,
		Type: ty,
		Val:  val,
	}
}

func matchedNary(tok *mTypes.Token) (mTypes.NodeKind, bool) {
	var tokenMap = map[string]mTypes.NodeKind{
		mTypes.NARY_OPERATOR_ADD: mTypes.ND_ADD,
	}

	return matchedOperator(tok, tokenMap)
}

func matchedBinary(tok *mTypes.Token) (mTypes.NodeKind, bool) {
	var tokenMap = map[string]mTypes.NodeKind{
		mTypes.BINARY_OPERATOR_EQ: mTypes.ND_EQ,
	}

	return matchedOperator(tok, tokenMap)
}

func matchedOperator(
	tok *mTypes.Token,
	tokenMap map[string]mTypes.NodeKind,
) (mTypes.NodeKind, bool) {

	if !tok.IsKind(mTypes.TK_OPERATOR) {
		return "", false
	}

	if kind, exists := tokenMap[tok.Val]; exists {
		return kind, true
	}

	return "", false
}

func matchedType(
	tok *mTypes.Token,
) (mTypes.ModoType, bool) {
	var typeMap = map[string]mTypes.ModoType{
		mTypes.TK_TYPE_INT: mTypes.TY_INT32,
		mTypes.TK_TYPE_STR: mTypes.TY_STR,
		mTypes.TK_TYPE_NIL: mTypes.TY_NIL,
	}

	if kind, exists := typeMap[tok.Kind]; exists {
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
	for nextToken.IsKind(mTypes.TK_INT) ||
		nextToken.IsKind(mTypes.TK_STR) ||
		nextToken.IsKind(mTypes.TK_IDENT) ||
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

		if tok.IsKind(mTypes.TK_DECLARE) {
			tok, head = parseDeclare(tok.Next, mTypes.ND_DECLARE)
			head = newNodeParent(mTypes.ND_DECLARE, head, "")

		} else if tok.IsKind(mTypes.TK_LAMBDA) {
			// arguments
			argHead := &mTypes.Node{}
			argCur := argHead

			tok = tok.Next

			for tok = tok.Next; !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE); tok = tok.Next {
				if !tok.IsKind(mTypes.TK_IDENT) {
					log.Panic("unresolved token kind, must be TK_IDENT: have %+v", tok)
				}

				argCur.Next = newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)
				argCur = argCur.Next
			}

			// expressions body
			// NOTE: why passing ]?
			tok, head = parseBody(tok, mTypes.ND_EXPR, "")
			head = newNodeParent(
				mTypes.ND_LAMBDA,
				head,
				"",
			)
			head.Args = argHead.Next

		} else if tok.IsKind(mTypes.TK_BIND) {
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
			// NOTE: why must proceed additionally?
			nextToken = nextToken.Next

			return nextToken, bind

		} else if kind, isNary := matchedNary(tok); isNary {
			log.Debug("is nary operator :have %+v", tok)
			tok, head = parseBody(tok, kind, tok.Val)
			head.Type = mTypes.TY_INT32

		} else if kind, isBinary := matchedBinary(tok); isBinary {
			log.Debug("is binary operator :have %+v", tok)
			tok, head = parseBody(tok, kind, tok.Val)
			// FIXME: TY_INT32 => TY_BOOL
			head.Type = mTypes.TY_INT32

		} else if tok.IsKind(mTypes.TK_LIBCALL) {
			log.Debug("is Library :have %+v", tok)
			tok, head = parseBody(tok, mTypes.ND_LIBCALL, tok.Val)

		} else if tok.IsKind(mTypes.TK_IF) {
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

		} else if tok.IsKind(mTypes.TK_IDENT) {
			log.Debug("is calling function :have %+v", tok)
			return parseBody(tok, mTypes.ND_FUNCCALL, tok.Val)
		}

		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_CLOSE) {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.IsKind(mTypes.TK_IDENT) {
		if parentKind == mTypes.ND_DECLARE {
			log.Debug("is Variable declaration :have %+v", tok)

			v := tok.Val
			t := &mTypes.Node{}
			h := t

			if tok.Next.IsKind(mTypes.TK_TYPE_SIG) {
				tok = tok.Next.Next
				for {
					if !tok.IsKindType() {
						break
					}
					ty, _ := matchedType(tok)
					t.Type = ty
					t.Next = &mTypes.Node{}
					t = t.Next

					tok = tok.Next
					if tok.IsKind(mTypes.TK_TYPE_ARROW) {
						tok = tok.Next
					}
				}
			}

			tok, head = parseDeclare(tok, mTypes.ND_VAR_DECLARE)
			if head.Args != nil {
				for a := head.Args; a != nil; a = a.Next {
					if h.Type == "" {
						log.Panic("type required :have %+v, %+v", t, a)
					}
					a.Type = h.Type
					h = h.Next
				}

			}

			res := newNodeParent(mTypes.ND_VAR_DECLARE, head, v)
			res.Type = h.Type
			return tok, res

		} else {
			log.Debug("is Variable reference :have %+v", tok)
			return tok.Next, newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)

		}

	} else if tok.IsKind(mTypes.TK_INT) {
		return tok.Next, newNodeScalar(mTypes.TY_INT32, tok.Val)

	} else if tok.IsKind(mTypes.TK_STR) {
		strNode := newNodeScalar(mTypes.TY_STR, tok.Val)
		strNode.Len = uint64(len(tok.Val))
		return tok.Next, strNode

	} else {
		log.Panic("unresolved token :have %+v", tok)
	}

	return tok, head

}

func parseProgram(tok *mTypes.Token) *mTypes.Program {
	p := &mTypes.Program{}
	prevDeclare := p.Declares

	for tok != nil && tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_OPEN) {

		if !tok.Next.IsKind(mTypes.TK_DECLARE) {
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
	// TODO: confirm all token read
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
