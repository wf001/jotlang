package parser

import (
	"fmt"

	"github.com/wf001/modo/pkg/log"
	"github.com/wf001/modo/pkg/types"
)

type parser struct {
	token *types.Token
}

func newNode(kind types.NodeKind, child *types.Node) *types.Node {
	return &types.Node{
		Kind:  kind,
		Child: child,
	}
}

func newNodeNum(tok *types.Token) *types.Node {
	return &types.Node{
		Kind: types.ND_INT,
		Val:  tok.Val,
	}
}

func matchedFuncCall(tok *types.Token) (bool, types.NodeKind) {
	if tok.Kind != types.TK_OPERATOR {
		return false, types.ND_NIL
	}
	switch tok.Val {
	case types.OPERATOR_ADD:
		return true, types.ND_ADD
	}
	return false, types.ND_NIL
}

func program(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &types.Node{}

	if tok.Kind == types.TK_EOL {
		return tok, nil
	}
	if tok.IsParenOpen() {
		tok = tok.Next

		if isFuncCall, kind := matchedFuncCall(tok); isFuncCall {
			nextToken, argHead := program(tok.Next)
			prevNode := argHead
			for nextToken.IsNum() || nextToken.IsParenOpen() {
				nextToken, prevNode.Next = program(nextToken)
				prevNode = prevNode.Next
			}
			head = newNode(kind, argHead)
			tok = nextToken
		}

		if !tok.IsParenClose() {
			log.Panic("must be ) :have %+v", tok)
		}
		return tok.Next, head
	} else if tok.Kind == types.TK_NUM {
		return tok.Next, newNodeNum(tok)
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
