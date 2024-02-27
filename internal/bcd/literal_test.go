package bcd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsContract(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "KT1HBy1L43tiLe5MVJZ5RoxGy53Kx8kMgyoU",
			address: "KT1HBy1L43tiLe5MVJZ5RoxGy53Kx8kMgyoU",
			want:    true,
		}, {
			name:    "mv1CjxUBZGvKdUuh7nhxyUrV6Q1A8e7A1WnX",
			address: "mv1CjxUBZGvKdUuh7nhxyUrV6Q1A8e7A1WnX",
			want:    false,
		}, {
			name:    "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			address: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want:    true,
		}, {
			name:    "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
			address: "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsContract(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			address: "KT1Ap287P1NzsnToSJdA4aqSNjPomRaHBZSr",
			want:    true,
		}, {
			name:    "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
			address: "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UAi",
			want:    true,
		}, {
			name:    "mv1CjxUBZGvKdUuh7nhxyUrV6Q1A8e7A1WnX",
			address: "mv1CjxUBZGvKdUuh7nhxyUrV6Q1A8e7A1WnX",
			want:    true,
		}, {
			name:    "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UA",
			address: "txr1YNMEtkj5Vkqsbdmt7xaxBTMRZjzS96UA",
			want:    false,
		}, {
			name:    "sr1J1ECygUgzE7urU3Ayr5HZaty83hpjbs28",
			address: "sr1J1ECygUgzE7urU3Ayr5HZaty83hpjbs28",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAddress(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsBakerHash(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "SG1d1wsgMKvSstzZQ8L4WoskCesdWGzVt5k4",
			str:  "SG1d1wsgMKvSstzZQ8L4WoskCesdWGzVt5k4",
			want: true,
		}, {
			name: "SG1d1wsgMKvSstzZQ8L4WoskCesdWGzVt5k",
			str:  "SG1d1wsgMKvSstzZQ8L4WoskCesdWGzVt5k",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBakerHash(tt.str)
			require.Equal(t, tt.want, got)
		})
	}
}
