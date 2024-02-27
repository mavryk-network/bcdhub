package bcd

import (
	"regexp"
	"strings"

	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
)

// IsLiteral -
func IsLiteral(prim string) bool {
	for _, s := range []string{
		consts.CONTRACT, consts.BYTES, consts.ADDRESS, consts.KEYHASH,
		consts.KEY, consts.TIMESTAMP, consts.BOOL, consts.MUTEZ,
		consts.NAT, consts.STRING, consts.INT, consts.SIGNATURE,
	} {
		if prim == s {
			return true
		}
	}
	return false
}

// IsContractLazy -
func IsContractLazy(str string) bool {
	return len(str) == 36 && strings.HasPrefix(str, "KT1")
}

// IsAddressLazy -
func IsAddressLazy(address string) bool {
	return (len(address) == 36 && (strings.HasPrefix(address, "KT") ||
		strings.HasPrefix(address, "mv") ||
		strings.HasPrefix(address, "sr1"))) ||
		(len(address) == 37 && strings.HasPrefix(address, "txr"))
}

// IsRollupAddressLazy -
func IsRollupAddressLazy(address string) bool {
	return len(address) == 37 && strings.HasPrefix(address, "txr")
}

// IsRollupAddressLazy -
func IsSmartRollupAddressLazy(address string) bool {
	return len(address) == 36 && strings.HasPrefix(address, "sr1")
}

var (
	addressRegex     = regexp.MustCompile("(mv|KT|txr|sr)[0-9A-Za-z]{34}")
	contractRegex    = regexp.MustCompile("(KT1)[0-9A-Za-z]{33}")
	bakerHashRegex   = regexp.MustCompile("(SG1)[0-9A-Za-z]{33}")
	operationRegex   = regexp.MustCompile("^o[1-9A-HJ-NP-Za-km-z]{50}$")
	smartRollupRegex = regexp.MustCompile("(sr)[0-9A-Za-z]{34}")
)

// IsAddress -
func IsAddress(str string) bool {
	return addressRegex.MatchString(str)
}

// IsContract -
func IsContract(str string) bool {
	return contractRegex.MatchString(str)
}

// IsBakerHash -
func IsBakerHash(str string) bool {
	return bakerHashRegex.MatchString(str)
}

// IsOperationHash -
func IsOperationHash(str string) bool {
	return operationRegex.MatchString(str)
}

// IsSmartRollupHash -
func IsSmartRollupHash(str string) bool {
	return smartRollupRegex.MatchString(str)
}
