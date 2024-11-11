package parser

import (
	"fmt"

	"github.com/wf001/modo/internal/log"
	"github.com/wf001/modo/pkg/types"
)

func newNode(kind types.NodeKind, lhs *types.Node, rhs *types.Node) *types.Node {
	node := new(types.Node)
	node.Kind = kind
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}

func newNodeNum(tok *types.Token) *types.Node {
	node := new(types.Node)
	node.Kind = types.ND_INT
	node.Val = tok.Val
	return node
}

func primary(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))

	if tok.Kind != types.TK_NUM {
		log.Panic("must be number: have %+v", tok)
	}

	//validate if number
	return tok.Next, newNodeNum(tok)
}

func add(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))

	if tok.Val == "+" {
		tok, lNode := primary(tok.Next)
		tok, rNode := primary(tok)
		return tok, newNode(types.ND_ADD, lNode, rNode)
	} else {
		log.Panic("must be + :have %+v", tok)
	}
	return tok, nil
}

func expr(tok *types.Token) (*types.Token, *types.Node) {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))

	if tok.Val == "(" {
		tok = tok.Next
		tok, node := expr(tok)
		if tok.Val != ")" {
			log.Panic("must be ) : have %+v", tok)
		}
		tok = tok.Next
		return tok, node
	}

	return add(tok)
}

// return Node object from Token array
func Parse(tok *types.Token) *types.Node {
	log.Debug(log.GREEN(fmt.Sprintf("%+v", tok)))
	_, node := expr(tok)
	log.DebugNode(node, 0)

	return node
}
