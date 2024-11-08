package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIfNumber(t *testing.T) {
	assert.ElementsMatch(t, []string{}, splitExpression(""))
	assert.ElementsMatch(t, []string{"1"}, splitExpression("1"))
	assert.ElementsMatch(t, []string{"123", "+", "45"}, splitExpression("123+45"))
	assert.ElementsMatch(t, []string{"123", "-", "45"}, splitExpression("123-45"))
	assert.ElementsMatch(t, []string{"123", "*", "45"}, splitExpression("123*45"))
	assert.ElementsMatch(t, []string{"123", "/", "45"}, splitExpression("123/45"))
	assert.ElementsMatch(
		t,
		[]string{"(", "123", "+", "45", ")"},
		splitExpression("(123+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "a", "+", "45", ")"},
		splitExpression("(a+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "ae", "+", "45", ")"},
		splitExpression("(ae+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "ifa", "+", "45", ")"},
		splitExpression("(ifa+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "->", "(", "123", "+", "45", ")", ")"},
		splitExpression("(->(123+45))"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "if", "a", ">", "1", "(", "2", "+", "3", ")", "(", "4", "-", "5", ")", ")"},
		splitExpression("(if a>1 (2+3) (4-5))"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "if", "a", ">", "1", "(", "2", "+", "3", ")", "(", "4", "-", "5", ")", ")"},
		splitExpression("(if a>1 (2+3) (4-5))"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "defn", "f", "[", "arg", "]", "(", "let", "[", "x", "1", "y", "3", "]",
			"(", "if", "x", ">", "1", "(", "2", "+", "y", ")", "(", "4", "-", "5", ")",
			")",")",")",
		},
    splitExpression("(defn f [arg] (let [x 1 y 3] (if x>1 (2+y) (4-5))))"),
	)
	t.Log(REG_EXP)
}
