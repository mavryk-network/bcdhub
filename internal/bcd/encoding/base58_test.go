package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBase58String(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{
			name: "mv1",
			data: "mv18cho8jUqmnEEKX5cQuZKopjy6WxQ37w8s",
			want: "06a868bd80219eb1f6a25108d1bdaa98ae27b2d9",
		},
		{
			name: "mv1",
			data: "mv1EHAJYdXRjXBkzjbBpGmSTcwXPsD5cmfkP",
			want: "44c6f8bc6088cd3b64f0bca87f812634c3f0ed30",
		},
		{
			name: "mv1",
			data: "mv1NT8eveuNrmrdS8bgGZGzJgAdGokb6bUuh",
			want: "9e6ac2e529a49aedbcdd0ac9542d5c0f4ce76f77",
		},
		{
			name: "mv3",
			data: "mv3DafLtx9N5aM5LjVrQrqqFdfW8EsbUjxaK",
			want: "358cbffa97149631cfb999fa47f0035fb1ea8636",
		},
		{
			name: "KT",
			data: "KT1BUKeJTemAaVBfRz6cqxeUBQGQqMxfG19A",
			want: "1fb03e3ff9fedaf3a2200ffc64d27812da734bba",
		},
		{
			name: "secp256k1_public_key",
			data: "sppk7bMuoa8w2LSKz3XEuPsKx1WavsMLCWgbWG9CZNAsJg9eTmkXRPd",
			want: "030ed412d33412ab4b71df0aaba07df7ddd2a44eb55c87bf81868ba09a358bc0e0",
		},
		{
			name: "rollup address",
			data: "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
			want: "76a57f87ee7624b92ab1453c75ba5d29ed8fe0bf",
		},
		{
			name: "smart rollup commitment hash",
			data: "src13MtM1eBzxCH1FBhLAkAiWGW6JbjvycLeH6vuz5k9GSiTYTCTja",
			want: "751b92ce705ebc551917bb488310498e969d7a1261fda86b509e7da2c780ec8d",
		}, {
			name: "operation",
			data: "opMhUv6wNFN1k5DofyPonyts1mSo4cTjCsk19i4Li7reT1oAQ1K",
			want: "e1e56182ae0a8738f804ac7edfd199a6ce617a6ee6751a0c69a4ad972098ebc7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase58String(tt.data)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeBase58String(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		prefix  string
		want    string
		wantErr bool
	}{
		{
			name:   "mv1",
			data:   "06a868bd80219eb1f6a25108d1bdaa98ae27b2d9",
			prefix: "mv1",
			want:   "mv18cho8jUqmnEEKX5cQuZKopjy6WxQ37w8s",
		},
		{
			name:   "mv1",
			data:   "44c6f8bc6088cd3b64f0bca87f812634c3f0ed30",
			prefix: "mv1",
			want:   "mv1EHAJYdXRjXBkzjbBpGmSTcwXPsD5cmfkP",
		},
		{
			name:   "mv1",
			data:   "9e6ac2e529a49aedbcdd0ac9542d5c0f4ce76f77",
			prefix: "mv1",
			want:   "mv1NT8eveuNrmrdS8bgGZGzJgAdGokb6bUuh",
		},
		{
			name:   "mv3",
			data:   "358cbffa97149631cfb999fa47f0035fb1ea8636",
			prefix: "mv3",
			want:   "mv3DafLtx9N5aM5LjVrQrqqFdfW8EsbUjxaK",
		},
		{
			name:   "KT",
			data:   "1fb03e3ff9fedaf3a2200ffc64d27812da734bba",
			prefix: "KT1",
			want:   "KT1BUKeJTemAaVBfRz6cqxeUBQGQqMxfG19A",
		},
		{
			name:   "KT",
			data:   "6f516588d2ee560385e386708a13bd63da907cf3",
			prefix: "KT1",
			want:   "KT1JjN5bTE9yayzYHiBm6ruktwEWSHRF8aDm",
		},
		{
			name:   "KT",
			data:   "e5bae183211979a662665319a0900df3542e65ba",
			prefix: "KT1",
			want:   "KT1VXUBQbYMt58yoKhNo73Zf8HTMfAd8Fqge",
		},
		{
			name:   "sig",
			data:   "bdc36db614aaa6084549020d376bb2469b5ea888dca2f7afbe5a0095bcc45ca0d8b5f00a051969437fe092debbcfe19d66378fbb74104de7eb1ecd895a64a80a",
			prefix: "sig",
			want:   "signpEFVQ1rW3TnVhc3PXf6SHRj7PvxwfJhBukWfB5X9rDhzpEk3ms5gRh763e922n52uQcjeqhqPdYi7WbFs2ERrNAPmCZJ",
		},
		{
			name:   "sig",
			data:   "a04991b4e938cc42d6c01c42be3649a81a9f80d244d9b90e7ec4edf8e0a7b68b6c212da2fef076e48fed66802fa83442b960a36afdb3e60c3cf14d4010f41f03",
			prefix: "sig",
			want:   "sigixZejtj1GfDpyiWAQAmvbtnNmCXKyADqVvCaXJH9xHyhSnYYV8696Z3kkns5DNV7oMnMPfNzo3qm84DfEx1XG6saZmHiA",
		},
		{
			name:   "chainID/main",
			data:   "7a06a770",
			prefix: "Net",
			want:   "NetXdQprcVkpaWU",
		},
		{
			name:   "chainID/babylon",
			data:   "458aa837",
			prefix: "Net",
			want:   "NetXUdfLh6Gm88t",
		},
		{
			name:   "chainID/carthage",
			data:   "9caecab9",
			prefix: "Net",
			want:   "NetXjD3HPJJjmcd",
		},
		{
			name:   "chainID/zeronet",
			data:   "0f6f0310",
			prefix: "Net",
			want:   "NetXKakFj1A7ouL",
		},
		{
			name:   "ed25519_public_key",
			data:   "4e4ca2abb4baeed702a0ac5b0de9b5607dd1fedb399c0ce25e15b3868f67269e",
			prefix: "edpk",
			want:   "edpkuEhzJqdFBCWMw6TU3deADRK2fq3GuwWFUphwyH7ero1Na4oGFP",
		},
		{
			name:   "secp256k1_public_key",
			data:   "030ed412d33412ab4b71df0aaba07df7ddd2a44eb55c87bf81868ba09a358bc0e0",
			prefix: "sppk",
			want:   "sppk7bMuoa8w2LSKz3XEuPsKx1WavsMLCWgbWG9CZNAsJg9eTmkXRPd",
		},
		{
			name:   "p256_public_key",
			data:   "031a3ad5ea94de6912f9bc83fd31de49816e90602c5252d77b5b233bfe711b0dd2",
			prefix: "p2pk",
			want:   "p2pk66iTZwLmRPshQgUr2HE3RUzSFwAN5MNaBQ5rfduT1dGKXd25pNN",
		},
		{
			name:   "rollup address",
			data:   "76a57f87ee7624b92ab1453c75ba5d29ed8fe0bf",
			prefix: "txr1",
			want:   "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
		},
		{
			name:   "smart rollup commitment hash",
			prefix: "src1",
			data:   "751b92ce705ebc551917bb488310498e969d7a1261fda86b509e7da2c780ec8d",
			want:   "src13MtM1eBzxCH1FBhLAkAiWGW6JbjvycLeH6vuz5k9GSiTYTCTja",
		}, {
			name:   "operation",
			prefix: "o",
			data:   "53E4FC89A78AD0D8C5B02A03C679DE44A42ED2EF8283D7CB759DD57877646EBC",
			want:   "ooHAF4skSPJe5NwF7Xy9tJy4GSGZUDF5UCNfDFQ4GNKjg1uKajf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeBase58String(tt.data, []byte(tt.prefix))
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
