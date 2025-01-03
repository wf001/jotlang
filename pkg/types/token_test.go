package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchedNary(t *testing.T) {
	// is nary
	tok := &Token{Kind: TK_OPERATOR, Val: NARY_OPERATOR_ADD}
	res, ok := tok.MatchedNary()
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, ND_ADD, res)

	// is not nary
	tok = &Token{Kind: TK_BIND, Val: SYMBOL_LET}
	res, ok = tok.MatchedNary()
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, "", res)
}

func TestMatchedBinary(t *testing.T) {
	// is binary
	tok := &Token{Kind: TK_OPERATOR, Val: BINARY_OPERATOR_EQ}
	res, ok := tok.MatchedBinary()
	assert.EqualValues(t, true, ok)
	assert.EqualValues(t, ND_EQ, res)

	// is not binary
	tok = &Token{Kind: TK_BIND, Val: SYMBOL_LET}
	res, ok = tok.MatchedNary()
	assert.EqualValues(t, false, ok)
	assert.EqualValues(t, "", res)
}
