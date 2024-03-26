package python

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGen(t *testing.T) {
	tests := []struct {
		name   string
		in     []*Node
		expect string
	}{
		{
			"int",
			[]*Node{
				NewNode(FunctionDefine, &FunctionDefineField{
					TType:  nil,
					Ident:  NewNode(Ident, &IdentField{S: "main"}),
					Params: nil,
					Block: NewNode(Block, &BlockField{Stmts: []*Node{
						NewNode(Return, &ReturnField{Value: NewNode(Literal, &LiteralField{TType: Integer, I: 32})}),
					}}),
				}),
			},
			"def main():\n    return 32\nif __name__ == \"__main__\":\n    main()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gen(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.expect, got); diff != "" {
				t.Fatalf("%v", diff)
			}
		})
	}
}
