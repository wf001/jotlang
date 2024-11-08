package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestParseIfNumber(t *testing.T) {
  assert.Exactly(t, "0", parseIfNumber("0"))
  assert.Exactly(t, "1", parseIfNumber("1"))
  assert.Exactly(t, "2", parseIfNumber("2"))
  assert.Exactly(t, "3", parseIfNumber("3"))
  assert.Exactly(t, "4", parseIfNumber("4"))
  assert.Exactly(t, "5", parseIfNumber("5"))
  assert.Exactly(t, "6", parseIfNumber("6"))
  assert.Exactly(t, "7", parseIfNumber("7"))
  assert.Exactly(t, "8", parseIfNumber("8"))
  assert.Exactly(t, "9", parseIfNumber("9"))
}

func TestParseIfOperator(t *testing.T) {
  assert.Exactly(t, "+", parseIfOperator("+"))
  assert.Exactly(t, "-", parseIfOperator("-"))
  assert.Exactly(t, "*", parseIfOperator("*"))
  assert.Exactly(t, "/", parseIfOperator("/"))
  assert.Exactly(t, "", parseIfOperator("|"))
}

