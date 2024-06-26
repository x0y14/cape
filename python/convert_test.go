package python

import (
	"cape/interlang"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestConvertNodeFromInterLang(t *testing.T) {
	tests := []struct {
		name   string
		in     []*interlang.Node
		expect []*Node
	}{
		{
			"def",
			[]*interlang.Node{
				interlang.NewNode(interlang.FunctionDefine, &interlang.FunctionDefineField{
					TType:  interlang.Integer,
					Ident:  interlang.NewNode(interlang.Ident, &interlang.IdentField{S: "main"}),
					Params: nil,
					Block:  interlang.NewNode(interlang.Block, &interlang.BlockField{Stmts: nil}),
				}),
			},
			[]*Node{
				NewNode(FunctionDefine, &FunctionDefineField{
					TType:  nil,
					Ident:  NewNode(Ident, &IdentField{S: "main"}),
					Params: nil,
					Block:  NewNode(Block, &BlockField{Stmts: nil}),
				}),
			},
		},
		{
			"int",
			[]*interlang.Node{
				interlang.NewNode(interlang.FunctionDefine, &interlang.FunctionDefineField{
					TType:  interlang.Integer,
					Ident:  interlang.NewNode(interlang.Ident, &interlang.IdentField{S: "main"}),
					Params: nil,
					Block: interlang.NewNode(interlang.Block, &interlang.BlockField{Stmts: []*interlang.Node{
						interlang.NewNode(interlang.Return, &interlang.ReturnField{Value: interlang.NewNode(interlang.Literal, &interlang.LiteralField{TType: interlang.Integer, I: 32})}),
					}}),
				}),
			},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertNodeFromInterLang(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.expect, got); diff != "" {
				t.Fatalf("%v", diff)
			}
		})
	}
}
