package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mTypes "github.com/wf001/modo/pkg/types"
)

func TestNewTokenMap(t *testing.T) {
	res, ok := matchedOperator(
		&mTypes.Token{Kind: mTypes.TK_OPERATOR, Val: mTypes.NARY_OPERATOR_ADD},
	)
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, mTypes.ND_ADD, res)

	res, ok = matchedOperator(
		&mTypes.Token{Kind: mTypes.TK_BIND, Val: mTypes.SYMBOL_LET},
	)
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, mTypes.ND_NIL, res)
}
