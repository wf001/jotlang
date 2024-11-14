package parser

import (
	"fmt"

	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/types"
)

type Parser struct {
	Token *types.Token
}

func ConstructParser(token *types.Token) *Parser {
	return &Parser{
		Token: token,
	}
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

func expr(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	head := &types.Node{}

	if tok.Kind == types.TK_EOL {
		return tok, nil
	}
	if tok.IsParenOpen() {
		tok = tok.Next
		if tok.IsOperationAdd() {
			nextToken, childHead := expr(tok.Next)
			prevNode := childHead
			for nextToken.IsNum() || nextToken.IsParenOpen() {
				nextToken, prevNode.Next = expr(nextToken)
				prevNode = prevNode.Next
			}
			tok = nextToken
			head = newNode(types.ND_ADD, childHead)
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

// take Token object, return Node object
func (p Parser) Parse() *types.Node {
	_, node := expr(p.Token)
	node.DebugNode(0)

	return node
}
