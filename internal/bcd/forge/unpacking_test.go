package forge

import (
	"encoding/hex"
	"testing"

	"github.com/mavryk-network/bcdhub/internal/bcd/base"
	"github.com/mavryk-network/bcdhub/internal/testsuite"
	"github.com/stretchr/testify/require"
)

func TestCollectStrings(t *testing.T) {
	tests := []struct {
		want      []string
		name      string
		tree      string
		tryUnpack bool
		wantErr   bool
	}{
		{
			name:      "pair without unpack",
			tree:      `{"args":[{"bytes": "00000c9b9e93efaac92e71f2c1ec48bb35848efeba70"},{"bytes": "0000b240dadc291b4fd6f1328f60ed463264c0d17e97"}],"prim": "Pair"}`,
			tryUnpack: false,
			want:      []string{},
		}, {
			name:      "pair with unpack",
			tree:      `{"args":[{"bytes": "00000c9b9e93efaac92e71f2c1ec48bb35848efeba70"},{"bytes": "0000b240dadc291b4fd6f1328f60ed463264c0d17e97"}],"prim": "Pair"}`,
			tryUnpack: true,
			want: []string{
				"mv19AAXe9xMRRVw9KrundnJSxKMSNas4cisy",
				"mv1QG1y4BSr3sqDtfVafhUjMDT9eGRrZMHwd",
			},
		}, {
			name:      "bytes with unpack",
			tree:      `{"bytes": "0507070100000004636f6465010000000863616c6c4275726e"}`,
			tryUnpack: true,
			want: []string{
				"code",
				"callBurn",
			},
		}, {
			name:      "bytes failed unpack",
			tree:      `{"bytes": "056459e92a8506c310fb72e5af90d00dbc1b15dc9288efc9a2ff47925ef9625bed7f969e938c3de13cb2fd60b3ab148816c2643c5625795ec81a183fe956c838"}`,
			tryUnpack: true,
			want:      []string{},
		}, {
			name:      "bytes with unpack",
			tree:      `{"bytes": "05020000013f03210316051f02000000020317050d036e072f0200000029034f07430368010000001a55706172616d417267756d656e74556e7061636b4661696c6564034203270200000000034203210316051f02000000020317051f02000000af0321074303690a0000000b0501000000056f776e65720329072f02000000210743036801000000165553746f72653a206e6f206669656c64206f776e657203270200000000050d036e072f020000002907430368010000001e5553746f72653a206661696c656420746f20756e7061636b206f776e657203270200000000034803190325072c0200000000020000001f034f07430368010000001053656e64657249734e6f744f776e6572034203270346030c0346074303690a0000000e0501000000086e65774f776e65720350053d036d034203210316051f020000000203170342"}`,
			tryUnpack: true,
			want: []string{
				"UparamArgumentUnpackFailed",
				"owner",
				"UStore: no field owner",
				"UStore: failed to unpack owner",
				"SenderIsNotOwner",
				"newOwner",
			},
		}, {
			name:      "bytes with unpack",
			tree:      `{"bytes": "0507070100000004636f6465010000001563616c6c5472616e736665724f776e657273686970"}`,
			tryUnpack: true,
			want: []string{
				"code",
				"callTransferOwnership",
			},
		}, {
			name:      "bytes with unpack",
			tree:      `{"bytes": "05020000014203210316051f02000000020317050d036e072f0200000029034f07430368010000001a55706172616d417267756d656e74556e7061636b4661696c6564034203270200000000034203210316051f02000000020317051f02000000af0321074303690a0000000b0501000000056f776e65720329072f02000000210743036801000000165553746f72653a206e6f206669656c64206f776e657203270200000000050d036e072f020000002907430368010000001e5553746f72653a206661696c656420746f20756e7061636b206f776e657203270200000000034803190325072c0200000000020000001f034f07430368010000001053656e64657249734e6f744f776e657203420327030c0346074303690a0000001305010000000d72656465656d416464726573730350053d036d034203210316051f020000000203170342"}`,
			tryUnpack: true,
			want: []string{
				"UparamArgumentUnpackFailed",
				"owner",
				"UStore: no field owner",
				"UStore: failed to unpack owner",
				"SenderIsNotOwner",
				"redeemAddress",
			},
		}, {
			name:      "tz2",
			tree:      `{"bytes": "00012ffebbf1560632ca767bc960ccdb84669d284c2c"}`,
			tryUnpack: true,
			want: []string{
				"mv2QQ5sHsmFuksCRmRgkZpp2DUHBxrZkQzcZ",
			},
		}, {
			name:      "tz3",
			tree:      `{"bytes": "000247d8c0238fc2f5a3b6c2e16b19a2283323dfdbba"}`,
			tryUnpack: true,
			want: []string{
				"mv3FFQYxx8JR4vUjtMKBLAefZyYHGycEazFF",
			},
		}, {
			name:      "KT1",
			tree:      `{"bytes": "0127cdfb0a9737d1e97e9ac47b71406d0b6b8bd8a500"}`,
			tryUnpack: true,
			want: []string{
				"KT1CDEg2oY3VfMa1neB7hK5LoVMButvivKYv",
			},
		}, {
			name:      "simple string",
			tree:      `{"string":"BAL-USDT"}`,
			tryUnpack: true,
			want: []string{
				"BAL-USDT",
			},
		}, {
			name:      "simple bytes",
			tree:      `{"bytes":"74657a6f732d73746f726167653a636f6e74656e74"}`,
			tryUnpack: true,
			want:      []string{"mavryk-storage:content"},
		}, {
			name:      "ipfs test",
			tree:      `{"bytes":"050100000035697066733a2f2f516d585a4846695a5a35566747794c634b514c4d6b5032314e733855394e47316d6f707945777348446663575835"}`,
			tryUnpack: true,
			want:      []string{"ipfs://QmXZHFiZZ5VgGyLcKQLMkP21Ns8U9NG1mopyEwsHDfcWX5"},
		}, {
			name:      "tezos domain test",
			tree:      `{"bytes":"62616c6c732e74657a"}`,
			tryUnpack: true,
			want:      []string{"balls.tez"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node base.Node
			err := json.UnmarshalFromString(tt.tree, &node)
			require.NoError(t, err)

			got, err := CollectStrings(&node, tt.tryUnpack)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []*base.Node
		wantErr bool
	}{
		{
			name: "test 1",
			data: "050100000035697066733a2f2f516d585a4846695a5a35566747794c634b514c4d6b5032314e733855394e47316d6f707945777348446663575835",
			want: []*base.Node{
				{
					StringValue: testsuite.Ptr("ipfs://QmXZHFiZZ5VgGyLcKQLMkP21Ns8U9NG1mopyEwsHDfcWX5"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := hex.DecodeString(tt.data)
			require.NoError(t, err)

			got, err := Unpack(b)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
