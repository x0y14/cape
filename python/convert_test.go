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
			"int",
			[]*interlang.Node{
				interlang.NewNode(interlang.FunctionDefine, &interlang.FunctionDefineField{
					TType:  nil,
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
