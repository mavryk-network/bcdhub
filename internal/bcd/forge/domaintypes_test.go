package forge

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Contract(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    string
		wantErr bool
	}{
		{
			name: "tz1eLWfccL46VAUjtyz9kEKgzuKnwyZH4rTA",
			val:  "tz1eLWfccL46VAUjtyz9kEKgzuKnwyZH4rTA",
			want: "0000cd1a410ffd5315ded34337f5f76edff48a13999a",
		}, {
			name: "Case 1",
			val:  "tz1d75oB6T4zUMexzkr5WscGktZ1Nss1JrT7",
			want: "0000bf97f5f1dbfd6ada0cf986d0a812f1bf0a572abc",
		},
		{
			name: "tz1 address",
			val:  `tz1a5fMLLY5WCarCzH7RKTJHX9mJFN8eaaWG`,
			want: "00009e6ac2e529a49aedbcdd0ac9542d5c0f4ce76f77",
		},
		{
			name: "tz1 address",
			val:  "tz1MBqYpcoGU93c1bePp5A6dmwKYjmHXRopS",
			want: "000010fc2282886d9cf8a1eebdc2733e302c7b110f38",
		},
		{
			name: "tz1 address",
			val:  "tz1LFEVYR7YRCxT6Nm3Zfjdnfj77xZqhbR5U",
			want: "000006a868bd80219eb1f6a25108d1bdaa98ae27b2d9",
		},
		{
			name: "tz1 address",
			val:  "tz1RABAzdLWVvxAFf1wpeUALAkp32mVhSGXX",
			want: "00003c8c2fe0f75ce212558df94c7a7306c2eeadd979",
		},
		{
			name: "tz1 address",
			val:  "tz1SZZgtvMVXaBKPcez4gfjKUsDz1gs6vg6X",
			want: "00004bf0acca4cc9e034b1d5f0f783c78e5ed44d866e",
		},
		{
			name: "tz1 address",
			val:  "tz1WkaeRycRr999GrVFepJd9Nqi1FWqGyGqq",
			want: "000079e68d8f0a8d64ec856e193efc0a347ef4adf8ee",
		},
		{
			name: "tz1 address",
			val:  "tz1eq3gqb2iZHjHVHoPJqV84gZdBF2TMQiH4",
			want: "0000d27fcbd31910d2226ba4c8f646d3d4c7b2f3a756",
		},
		{
			name: "tz1 address",
			val:  "tz1M9CMEtsXm3QxA7FmMU2Qh7xzsuGXVbcDr",
			want: "0000107c4009f2bcfcc248d6952998af5b7203b8ff59",
		},
		{
			name: "tz2 address",
			val:  "tz28YZoayJjVz2bRgGeVjxE8NonMiJ3r2Wdu",
			want: "0001028562fb176188114cf437a757cdc75bc4aa8cae",
		},
		{
			name: "tz3 address",
			val:  "tz3agP9LGe2cXmKQyYn6T68BHKjjktDbbSWX",
			want: "00029d6a61cd3510193e257128da8f09a0b173bff695",
		},
		{
			name: "KT address",
			val:  "KT1J8T7U6J1BAo9fJAxvedHsNErnejwvPyUH",
			want: "0168b709e887ddc34c3c9e468b5819b2f012b60ef700",
		},
		{
			name: "KT address",
			val:  "KT1BUKeJTemAaVBfRz6cqxeUBQGQqMxfG19A",
			want: "011fb03e3ff9fedaf3a2200ffc64d27812da734bba00",
		},
		{
			name: "KT address",
			val:  "KT1U1JZaXoG4u1EPnhHL4R4otzkWc1L34q3c",
			want: "01d50e3f6f059dc86f5591455549313ce42d0c50f100",
		},
		{
			name: "KT address",
			val:  "KT1XHAmdRKugP1Q38CxDmpcRSxq143KpEiYx",
			want: "01f8f6c6a0af7c20251bc7df108f2a6e2879a06c9a00",
		}, {
			name: "KT address with entrypoint",
			val:  `KT1Nh9wK8W3j3CXeTVm5DTTaiU5RE8CxLWZ4%receive_bunny_balance`,
			want: "019ac6ee79c4e87a21d094bb8bf00f37fe51717e8700726563656976655f62756e6e795f62616c616e6365",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contract(tt.val)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnforgeBakerHash(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			str:  "94697e9229c88fac7d19d62e139ca6735f9569dd",
			want: "SG1d1wsgMKvSstzZQ8L4WoskCesdWGzVt5k4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnforgeBakerHash(tt.str)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnforgeAddress(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			str:  "016e4943f7a23ab9cbe56f48ff72f6c27e8956762400",
			want: "KT1JdufSdfg3WyxWJcCRNsBFV9V3x9TQBkJ2",
		}, {
			name: "test 2",
			str:  "00003a96709901319a5da2968782279dae581b9ba4",
			want: "mv182iBTCasWb9JM4wuv5avQfvJxTJMdCLQG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnforgeAddress(tt.str)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnforgePublicKey(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			str:  "0103682c3aaa998fd9adfe8111cd42cc0daedb5d97647e6020eb629fbc91b613f721",
			want: "sppk7c3Fz7QqhZqY2FZUWWAnDuqTwx4KwDjgFA4VeLPiV8n4tnbsVzG",
		}, {
			name: "test 2",
			str:  "0028fc6875ca69a6f5bde4f377bfcde72fd618bcfa52e7272c7b788d1165449eb4",
			want: "edpktxGsKjnk43ZZ7v6gJe6PFV85peHvoWqVUzDQjTfN8idYwVkBwN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnforgePublicKey(tt.str)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestUnforgeContract(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			str:  "019ac6ee79c4e87a21d094bb8bf00f37fe51717e8700726563656976655f62756e6e795f62616c616e6365",
			want: "KT1Nh9wK8W3j3CXeTVm5DTTaiU5RE8CxLWZ4%receive_bunny_balance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnforgeContract(tt.str)
			require.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
