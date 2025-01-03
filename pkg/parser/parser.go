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

func parseIdent(
	tok *mTypes.Token,
	head *mTypes.Node,
	parentKind mTypes.NodeKind,
) (*mTypes.Token, *mTypes.Node) {
	if parentKind == mTypes.ND_DECLARE {
		log.Debug("is Variable declaration :have %+v", tok)

		identName := tok.Val
		typeList := &mTypes.Node{}
		h := typeList

		if tok.Next.IsKind(mTypes.TK_TYPE_SIG) {
			tok = tok.Next.Next
			for {
				if !tok.IsKindType() {
					break
				}
				ty, _ := tok.MatchedType()
				typeList.Type = ty
				typeList.Next = &mTypes.Node{}
				typeList = typeList.Next

				tok = tok.Next
				if tok.IsKind(mTypes.TK_TYPE_ARROW) {
					tok = tok.Next
				}
			}
		} else {
			log.Panic("must be :: :have %+v", tok)
		}

		tok, head = parseDeclare(tok, mTypes.ND_VAR_DECLARE)
		if head.Args != nil {
			for a := head.Args; a != nil; a = a.Next {
				if h.Type == "" {
					log.Panic("type required :have %+v, %+v", typeList, a)
				}
				a.Type = h.Type
				h = h.Next
			}

		}

		child := newNodeParent(mTypes.ND_VAR_DECLARE, head, identName)
		child.Type = h.Type
		return tok, child

	} else {
		log.Debug("is Variable reference :have %+v", tok)
		return tok.Next, newNodeParent(mTypes.ND_VAR_REFERENCE, nil, tok.Val)

	}

}

func parseLambda(tok *mTypes.Token, head *mTypes.Node) (*mTypes.Token, *mTypes.Node) {

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

	return tok, head
}

// NOTE: typed at here: ND_SCALAR, ND_VAR_DECLARE, ND_VAR_REFERENCE(Args), ND_EQ, ND_ADD
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
			tok, head = parseLambda(tok, head)

		} else if tok.IsKind(mTypes.TK_BIND) {
			tok = tok.Next

			if tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.BRACKET_CLOSE) {
				log.Panic("let bindings require even number of forms")
			}

			prev := &mTypes.Node{}
			varHead := prev

			tok = tok.Next
			for {

				tok, prev.Next = parseDeclare(tok, mTypes.ND_DECLARE)
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

		} else if kind, isNary := tok.MatchedNary(); isNary {
			log.Debug("is nary operator :have %+v", tok)
			tok, head = parseBody(tok, kind, tok.Val)
			head.Type = mTypes.TY_INT32

		} else if kind, isBinary := tok.MatchedBinary(); isBinary {
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
			tok, head = parseBody(tok, mTypes.ND_FUNCCALL, tok.Val)
		}

		if !tok.IsKindAndVal(mTypes.TK_PAREN, mTypes.PARREN_CLOSE) {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head

	} else if tok.IsKind(mTypes.TK_IDENT) {
		return parseIdent(tok, head, parentKind)

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
