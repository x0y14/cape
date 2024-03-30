package from_inter

import (
	"cape/c"
	"cape/interlang"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestConvertNodeFromInterLang(t *testing.T) {
	tests := []struct {
		name   string
		in     []*interlang.Node
		expect []*c.Node
	}{
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
			[]*c.Node{
				c.NewNode(c.FunctionDefine, &c.FunctionDefineField{
					TType:  c.Integer,
					Ident:  c.NewNode(c.Ident, &c.IdentField{S: "main"}),
					Params: nil,
					Block: c.NewNode(c.Block, &c.BlockField{Stmts: []*c.Node{
						c.NewNode(c.Return, &c.ReturnField{Value: c.NewNode(c.Literal, &c.LiteralField{TType: c.Integer, I: 32})}),
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
