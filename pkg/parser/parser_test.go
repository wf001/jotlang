package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mTypes "github.com/wf001/modo/pkg/types"
)

func TestMatchedNary(t *testing.T) {
	// is nary
	tok := &mTypes.Token{Kind: mTypes.TK_OPERATOR, Val: mTypes.NARY_OPERATOR_ADD}
	res, ok := tok.MatchedNary()
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, mTypes.ND_ADD, res)

	// is not nary
	tok = &mTypes.Token{Kind: mTypes.TK_BIND, Val: mTypes.SYMBOL_LET}
	res, ok = tok.MatchedNary()
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, "", res)
}

func TestMatchedBinary(t *testing.T) {
	// is binary
	tok := &mTypes.Token{Kind: mTypes.TK_OPERATOR, Val: mTypes.BINARY_OPERATOR_EQ}
	res, ok := tok.MatchedBinary()
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, mTypes.ND_EQ, res)

	// is not binary
	tok = &mTypes.Token{Kind: mTypes.TK_BIND, Val: mTypes.SYMBOL_LET}
	res, ok = tok.MatchedNary()
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, "", res)
}
