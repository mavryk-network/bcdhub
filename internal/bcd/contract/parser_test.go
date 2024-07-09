package contract

import (
	"testing"

	"github.com/mavryk-network/bcdhub/internal/bcd/base"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
	"github.com/stretchr/testify/require"
)

func Test_parseConstants(t *testing.T) {
	testConstant := "expru54tk2k4E81xQy63P6x3RijnTz51s2m7BV7pr3fDQH8YDqiYvR"
	tests := []struct {
		name string
		node *base.Node
		want string
	}{
		{
			name: "nil node",
			node: nil,
			want: "",
		}, {
			name: "not constant",
			node: &base.Node{
				Prim: consts.ADDRESS,
			},
			want: "",
		}, {
			name: "constant without args",
			node: &base.Node{
				Prim: consts.CONSTANT,
			},
			want: "",
		}, {
			name: "constant with arg but with nil value",
			node: &base.Node{
				Prim: consts.CONSTANT,
				Args: []*base.Node{
					{
						Prim: consts.STRING,
					},
				},
			},
			want: "",
		}, {
			name: "good",
			node: &base.Node{
				Prim: consts.CONSTANT,
				Args: []*base.Node{
					{
						Prim:        consts.STRING,
						StringValue: &testConstant,
					},
				},
			},
			want: "expru54tk2k4E81xQy63P6x3RijnTz51s2m7BV7pr3fDQH8YDqiYvR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseConstants(tt.node)
			require.Equal(t, tt.want, got)
		})
	}
}
