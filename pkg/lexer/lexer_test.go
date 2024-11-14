package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wf001/modo/pkg/types"
)

func TestSplitProgram(t *testing.T) {
	assert.ElementsMatch(t, []string{}, splitProgram(""))
	assert.ElementsMatch(t, []string{"1"}, splitProgram("1"))
	assert.ElementsMatch(t, []string{"123", "+", "45"}, splitProgram("123+45"))
	assert.ElementsMatch(t, []string{"123", "-", "45"}, splitProgram("123-45"))
	assert.ElementsMatch(t, []string{"123", "*", "45"}, splitProgram("123*45"))
	assert.ElementsMatch(t, []string{"123", "/", "45"}, splitProgram("123/45"))
	assert.ElementsMatch(
		t,
		[]string{"(", "123", "+", "45", ")"},
		splitProgram("(123+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "a", "+", "45", ")"},
		splitProgram("(a+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "ae", "+", "45", ")"},
		splitProgram("(ae+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "ifa", "+", "45", ")"},
		splitProgram("(ifa+45)"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "->", "(", "123", "+", "45", ")", ")"},
		splitProgram("(->(123+45))"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "if", "a", ">", "1", "(", "2", "+", "3", ")", "(", "4", "-", "5", ")", ")"},
		splitProgram("(if a>1 (2+3) (4-5))"),
	)
	assert.ElementsMatch(
		t,
		[]string{"(", "if", "a", ">", "1", "(", "2", "+", "3", ")", "(", "4", "-", "5", ")", ")"},
		splitProgram("(if a>1 (2+3) (4-5))"),
	)
	assert.ElementsMatch(
		t,
		[]string{
			"(", "defn", "f", "[", "arg", "]", "(", "let", "[", "x", "1", "y", "3", "]",
			"(", "if", "x", ">", "1", "(", "2", "+", "y", ")", "(", "4", "-", "5", ")",
			")", ")", ")",
		},
		splitProgram("(defn f [arg] (let [x 1 y 3] (if x>1 (2+y) (4-5))))"),
	)
	t.Log(types.ALL_REG_EXP)
}

func TestLexOneInteger(t *testing.T) {
	assert.Equal(t, &types.Token{Kind: types.TK_NUM, Val: "1"}, Lex("1"))
}

func TestLexOperationAdd(t *testing.T) {
	want := &types.Token{
		Kind: types.TK_PAREN,
		Val:  "(",
		Next: &types.Token{
			Kind: types.TK_OPERATOR,
			Val:  "+",
			Next: &types.Token{
				Kind: types.TK_NUM,
				Val:  "1",
				Next: &types.Token{
					Kind: types.TK_NUM,
					Val:  "2",
					Next: &types.Token{
						Kind: types.TK_PAREN,
						Val:  ")",
					},
				},
			},
		},
	}

	assert.EqualValues(t, want, Lex("(+ 1 2)"))
}

func TestLexOperationAddTakingAdd(t *testing.T) {
	want := &types.Token{
		Kind: types.TK_PAREN,
		Val:  "(",
		Next: &types.Token{
			Kind: types.TK_OPERATOR,
			Val:  "+",
			Next: &types.Token{
				Kind: types.TK_PAREN,
				Val:  "(",
				Next: &types.Token{
					Kind: types.TK_OPERATOR,
					Val:  "+",
					Next: &types.Token{
						Kind: types.TK_NUM,
						Val:  "1",
						Next: &types.Token{
							Kind: types.TK_NUM,
							Val:  "2",
							Next: &types.Token{
								Kind: types.TK_PAREN,
								Val:  ")",
								Next: &types.Token{
									Kind: types.TK_PAREN,
									Val:  "(",
									Next: &types.Token{
										Kind: types.TK_OPERATOR,
										Val:  "+",
										Next: &types.Token{
											Kind: types.TK_NUM,
											Val:  "3",
											Next: &types.Token{
												Kind: types.TK_NUM,
												Val:  "4",
												Next: &types.Token{
													Kind: types.TK_PAREN,
													Val:  ")",
													Next: &types.Token{
														Kind: types.TK_PAREN,
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
