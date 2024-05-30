package tokenize

import (
	"github.com/google/go-cmp/cmp"
	"log"
	"testing"
)

func createTokenWithPair(open, child, close *Token) *Token {
	child.Next = close
	open.Next = child
	return open
}

func getLast(token *Token) *Token {
	cur := token
	for {
		if cur.Next != nil {
			cur = cur.Next
		} else {
			break
		}
	}
	return cur
}

func TestGetLast(t *testing.T) {
	t.Logf("%v", getLast(&Token{Kind: Int, I: 1, Next: &Token{Kind: Int, I: 2}}).I == 2)
	t.Logf("%v", getLast(createTokenWithPair(&Token{Kind: Int, I: 1}, &Token{Kind: Int, I: 2}, &Token{Kind: Int, I: 3})).I == 3)
}

func createChain(tokens []*Token) *Token {
	head := Token{}
	cur := &head
	cur.Next = tokens[0]
	tokens = tokens[1:]
	for _, tok := range tokens {
		cur = getLast(cur)
		cur.Next = tok
	}

	return head.Next
}

func TestCreateTokenWithPair(t *testing.T) {
	tests := []struct {
		name   string
		in     [3]*Token
		expect *Token
	}{
		{
			"Rb",
			[3]*Token{
				{Kind: Lrb},
				{Kind: Int, I: 100},
				{Kind: Rrb},
			},
			&Token{
				Kind: Lrb,
				Next: &Token{
					Kind: Int,
					I:    100,
					Next: &Token{
						Kind: Rrb,
						Next: nil,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tok := createTokenWithPair(tt.in[0], tt.in[1], tt.in[2])
		if diff := cmp.Diff(tt.expect, tok); diff != "" {
			t.Errorf("%v", diff)
		}
	}
}

func TestCreateTokenChain(t *testing.T) {
	tests := []struct {
		name   string
		in     []*Token
		expect *Token
	}{
		{
			"123",
			[]*Token{
				{Kind: Int, I: 1},
				{Kind: Int, I: 2},
				{Kind: Int, I: 3},
			},
			&Token{
				Kind: Int,
				I:    1,
				Next: &Token{
					Kind: Int,
					I:    2,
					Next: &Token{
						Kind: Int,
						I:    3,
						Next: nil,
					},
				},
			},
		},
		{"123123",
			[]*Token{
				createTokenWithPair(&Token{Kind: Int, I: 1}, &Token{Kind: Int, I: 2}, &Token{Kind: Int, I: 3}),
				createTokenWithPair(&Token{Kind: Int, I: 4}, &Token{Kind: Int, I: 5}, &Token{Kind: Int, I: 6}),
			},
			&Token{
				Kind: Int,
				I:    1,
				Next: &Token{
					Kind: Int,
					I:    2,
					Next: &Token{
						Kind: Int,
						I:    3,
						Next: &Token{
							Kind: Int,
							I:    4,
							Next: &Token{
								Kind: Int,
								I:    5,
								Next: &Token{
									Kind: Int,
									I:    6,
									Next: nil,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tok := createChain(tt.in)
		if diff := cmp.Diff(tt.expect, tok); diff != "" {
			t.Errorf("%v", diff)
		}
	}
}

func TestTokenize(t *testing.T) {

	tests := []struct {
		name   string
		code   string
		expect *Token
	}{
		{
			"test",
			`

int main(void) {
	printf("hello\n");
	return 0;
}
`,
			&Token{
				Kind: Ident,
				S:    "int",
				Next: &Token{
					Kind: Ident,
					S:    "main",
					Next: &Token{
						Kind: Lrb,
						Next: &Token{
							Kind: Ident,
							S:    "void",
							Next: &Token{
								Kind: Rrb,
								Next: &Token{
									Kind: Lcb,
									Next: &Token{
										Kind: Ident,
										S:    "printf",
										Next: &Token{
											Kind: Lrb,
											Next: &Token{
												Kind: String,
												S:    "hello\n",
												Next: &Token{
													Kind: Rrb,
													S:    "",
													Next: &Token{
														Kind: Semi,
														Next: &Token{
															Kind: Ident,
															S:    "return",
															Next: &Token{
																Kind: Int,
																I:    0,
																Next: &Token{
																	Kind: Semi,
																	Next: &Token{
																		Kind: Rcb,
																		Next: &Token{
																			Kind: Eof,
																			Next: nil,
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
						},
					},
				},
			},
		},
		//		{
		//			name: "fizzbuzz",
		//			code: `
		//int main(void) {
		//    int i;
		//    for (i = 1; i <= 100; i=i+1) {
		//        if (i % 3 == 0 && i % 5 == 0) {
		//            printf("FizzBuzz\n");
		//        } else if (i % 3 == 0) {
		//            printf("Fizz\n");
		//        } else if (i % 5 == 0) {
		//            printf("Buzz\n");
		//        } else {
		//            printf("%d\n", i);
		//        }
		//    }
		//    return 0;
		//}`,
		//			expect: createChain([]*Token{
		//				{Kind: Ident, S: "int"}, {Kind: Ident, S: "main"}, createTokenWithPair(&Token{Kind: Lrb}, &Token{Kind: Ident, S: "void"}, &Token{Kind: Rrb}),
		//				createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // main func block-
		//					{Kind: Ident, S: "int"}, {Kind: Ident, S: "i"}, {Kind: Semi},
		//					{Kind: Ident, S: "for"}, createTokenWithPair(&Token{Kind: Lrb}, createChain([]*Token{ // for cond-
		//						{Kind: Ident, S: "i"}, {Kind: Assign}, {Kind: Int, I: 1}, {Kind: Semi},
		//						{Kind: Ident, S: "i"}, {Kind: Le}, {Kind: Int, I: 100}, {Kind: Semi},
		//						{Kind: Ident, S: "i"}, {Kind: Assign}, {Kind: Ident, S: "i"}, {Kind: Add}, {Kind: Int, I: 1},
		//					}), &Token{Kind: Rrb}), // -for cond
		//					createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // for block-
		//						{Kind: Ident, S: "if"}, createTokenWithPair(&Token{Kind: Lrb}, createChain([]*Token{ // if cond-
		//							{Kind: Ident, S: "i"}, {Kind: Mod}, {Kind: Int, I: 3}, {Kind: Eq}, {Kind: Int, I: 0},
		//							{Kind: And},
		//							{Kind: Ident, S: "i"}, {Kind: Mod}, {Kind: Int, I: 5}, {Kind: Eq}, {Kind: Int, I: 0},
		//						}), &Token{Kind: Rrb}), // -if cond
		//						createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // if block-
		//							{Kind: Ident, S: "printf"}, createTokenWithPair(&Token{Kind: Lrb}, &Token{Kind: String, S: "FizzBuzz\n"}, &Token{Kind: Rrb}), {Kind: Semi},
		//						}), &Token{Kind: Rcb}), // -if block
		//
		//						{Kind: Ident, S: "else"}, {Kind: Ident, S: "if"}, createTokenWithPair(&Token{Kind: Lrb}, createChain([]*Token{ // else cond I % 3-
		//							{Kind: Ident, S: "i"}, {Kind: Mod}, {Kind: Int, I: 3}, {Kind: Eq}, {Kind: Int, I: 0},
		//						}), &Token{Kind: Rrb}), // -else cond I % 3
		//						createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // else block I % 3-
		//							{Kind: Ident, S: "printf"}, createTokenWithPair(&Token{Kind: Lrb}, &Token{Kind: String, S: "Fizz\n"}, &Token{Kind: Rrb}), {Kind: Semi},
		//						}), &Token{Kind: Rcb}), // -else block I % 3
		//
		//						{Kind: Ident, S: "else"}, {Kind: Ident, S: "if"}, createTokenWithPair(&Token{Kind: Lrb}, createChain([]*Token{ // else cond I % 5-
		//							{Kind: Ident, S: "i"}, {Kind: Mod}, {Kind: Int, I: 5}, {Kind: Eq}, {Kind: Int, I: 0},
		//						}), &Token{Kind: Rrb}), // -else cond I % 5
		//						createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // else block I % 5-
		//							{Kind: Ident, S: "printf"}, createTokenWithPair(&Token{Kind: Lrb}, &Token{Kind: String, S: "Buzz\n"}, &Token{Kind: Rrb}), {Kind: Semi},
		//						}), &Token{Kind: Rcb}), // -else block I % 5
		//
		//						{Kind: Ident, S: "else"}, createTokenWithPair(&Token{Kind: Lcb}, createChain([]*Token{ // else block-
		//							{Kind: Ident, S: "printf"}, createTokenWithPair(&Token{Kind: Lrb}, createChain([]*Token{ // printf("%d\n", I)-
		//								{Kind: String, S: "%d\n"}, {Kind: Comma}, {Kind: Ident, S: "i"},
		//							}), &Token{Kind: Rrb}), {Kind: Semi}, // -printf("%d\n", I);
		//						}), &Token{Kind: Rcb}), // -else block
		//					}), &Token{Kind: Rcb}), // -for block
		//					{Kind: Ident, S: "return"}, {Kind: Int, I: 0}, {Kind: Semi},
		//				}), &Token{Kind: Rcb}), // -main func block
		//				{Kind: Eof},
		//			}),
		//		},
	}

	for _, tt := range tests {
		tok, err := tokenize(tt.code)
		if err != nil {
			log.Fatalf("failed: %v", err)
		}
		log.Printf("%v", tok)
		if diff := cmp.Diff(tt.expect, tok); diff != "" {
			t.Errorf("%v", diff)
		}
	}
}
