package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mTypes "github.com/wf001/modo/pkg/types"
)

func TestSplitProgram(t *testing.T) {
	assert.ElementsMatch(t, []string{}, splitString(""))
	assert.ElementsMatch(t, []string{"1"}, splitString("1"))
	assert.ElementsMatch(t, []string{"+", "123", "45"}, splitString("+ 123 45"))
	assert.ElementsMatch(t, []string{"-", "123", "45"}, splitString("- 123 45"))
	assert.ElementsMatch(t, []string{"*", "123", "45"}, splitString("* 123 45"))
	assert.ElementsMatch(t, []string{"/", "123", "45"}, splitString("/ 123 45"))
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "123", "45", ")"},
		splitString("(+ 123 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "a", "45", ")"},
		splitString("(+ a 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "ae", "45", ")"},
		splitString("(+ ae 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "+", "ifa", "45", ")"},
		splitString("(+ ifa 45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "->", "(", "+", "123", "45", ")", ")"},
		splitString("(->(+ 123 45))"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "=", "1", "1", ")",
		},
		splitString("(= 1 1)"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "def", "main", "(", "fn", "[", "]", "(", "if", "(", "=", "1", "1", ")", "(", "prn",
			"(", "+", "1", "2", ")", ")", "(", "prn", "(", "+", "3", "-3", ")", ")", ")", ")", ")",
		},
		splitString("(def main (fn [] (if (= 1 1) (prn (+ 1 2)) (prn (+ 3 -3)))))"),
	)
	t.Log(mTypes.ALL_REG_EXP)
}

func TestLexOneInteger(t *testing.T) {
	assert.Equal(t, &mTypes.Token{Kind: mTypes.TK_NUM, Val: "1"}, Lex("1"))
}

func TestLexOperationAdd(t *testing.T) {
	want := &mTypes.Token{
		Kind: mTypes.TK_PAREN,
		Val:  "(",
		Next: &mTypes.Token{
			Kind: mTypes.TK_OPERATOR,
			Val:  "+",
			Next: &mTypes.Token{
				Kind: mTypes.TK_NUM,
				Val:  "1",
				Next: &mTypes.Token{
					Kind: mTypes.TK_NUM,
					Val:  "2",
					Next: &mTypes.Token{
						Kind: mTypes.TK_PAREN,
						Val:  ")",
					},
				},
			},
		},
	}

	assert.EqualValues(t, want, Lex("(+ 1 2)"))
}

func TestLexOperationAddTakingAdd(t *testing.T) {
	want := &mTypes.Token{
		Kind: mTypes.TK_PAREN,
		Val:  "(",
		Next: &mTypes.Token{
			Kind: mTypes.TK_OPERATOR,
			Val:  "+",
			Next: &mTypes.Token{
				Kind: mTypes.TK_PAREN,
				Val:  "(",
				Next: &mTypes.Token{
					Kind: mTypes.TK_OPERATOR,
					Val:  "+",
					Next: &mTypes.Token{
						Kind: mTypes.TK_NUM,
						Val:  "1",
						Next: &mTypes.Token{
							Kind: mTypes.TK_NUM,
							Val:  "2",
							Next: &mTypes.Token{
								Kind: mTypes.TK_PAREN,
								Val:  ")",
								Next: &mTypes.Token{
									Kind: mTypes.TK_PAREN,
									Val:  "(",
									Next: &mTypes.Token{
										Kind: mTypes.TK_OPERATOR,
										Val:  "+",
										Next: &mTypes.Token{
											Kind: mTypes.TK_NUM,
											Val:  "3",
											Next: &mTypes.Token{
												Kind: mTypes.TK_NUM,
												Val:  "4",
												Next: &mTypes.Token{
													Kind: mTypes.TK_PAREN,
													Val:  ")",
													Next: &mTypes.Token{
														Kind: mTypes.TK_PAREN,
														Val:  ")",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.EqualValues(t, want, Lex("(+ (+ 1 2) (+ 3 4))"))
}
