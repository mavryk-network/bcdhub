package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress_Compare(t *testing.T) {
	tests := []struct {
		name       string
		first      string
		firstType  int
		second     string
		secondType int
		want       int
		wantErr    bool
	}{
		{
			name:       "equal",
			first:      "KT1Hbwyp8D39d3681bG4FtZ1rE1uopVmU4wK",
			firstType:  valueKindString,
			second:     "KT1Hbwyp8D39d3681bG4FtZ1rE1uopVmU4wK",
			secondType: valueKindString,
			want:       0,
		}, {
			name:       "unequal",
			first:      "KT1Hbwyp8D39d3681bG4FtZ1rE1uopVmU4wK",
			firstType:  valueKindString,
			second:     "KT1MjT5jseoujXvy1w2PjdaFXYo8jeh8k5S2",
			secondType: valueKindString,
			want:       -1,
		}, {
			name:       "equal",
			first:      "mv1ShyyCvhMT4SFy3JYzz41i9vBmWN4sfob7",
			firstType:  valueKindString,
			second:     "0000cd1a410ffd5315ded34337f5f76edff48a13999a",
			secondType: valueKindBytes,
			want:       0,
		}, {
			name:       "equal",
			first:      "0000cd1a410ffd5315ded34337f5f76edff48a13999a",
			firstType:  valueKindBytes,
			second:     "mv1ShyyCvhMT4SFy3JYzz41i9vBmWN4sfob7",
			secondType: valueKindString,
			want:       0,
		}, {
			name:       "equal",
			first:      "0000cd1a410ffd5315ded34337f5f76edff48a13999a",
			firstType:  valueKindBytes,
			second:     "mv1ShyyCvhMT4SFy3JYzz41i9vBmWN4sfob7",
			secondType: valueKindString,
			want:       0,
		}, {
			name:       "unequal",
			first:      "0000cd1a410ffd5315ded34337f5f76edff48a13999a",
			firstType:  valueKindBytes,
			second:     "KT1DEkR3cErDAn6oH4jK8Z7n9a4oCXRZZwYa",
			secondType: valueKindString,
			want:       -1,
		}, {
			name:       "unequal",
			first:      "mv1ShyyCvhMT4SFy3JYzz41i9vBmWN4sfob7",
			firstType:  valueKindString,
			second:     "KT1DEkR3cErDAn6oH4jK8Z7n9a4oCXRZZwYa",
			secondType: valueKindString,
			want:       -1,
		}, {
			name:       "unequal",
			first:      "KT1DEkR3cErDAn6oH4jK8Z7n9a4oCXRZZwYa",
			firstType:  valueKindString,
			second:     "mv1ShyyCvhMT4SFy3JYzz41i9vBmWN4sfob7",
			secondType: valueKindString,
			want:       1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first := NewAddress(0)
			first.Value = tt.first
			first.ValueKind = tt.firstType

			second := NewAddress(0)
			second.Value = tt.second
			second.ValueKind = tt.secondType

			got, err := first.Compare(second)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
