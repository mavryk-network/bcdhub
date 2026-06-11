package types

import "github.com/mavryk-network/nexushub/internal/nexus"

// AccountType -
type AccountType int

// account types
const (
	AccountTypeUnknown = iota
	AccountTypeContract
	AccountTypeTz
	AccountTypeRollup
	AccountTypeSmartRollup
)

// NewAccountType -
func NewAccountType(address string) AccountType {
	switch {
	case nexus.IsContract(address):
		return AccountTypeContract
	case nexus.IsAddressLazy(address):
		return AccountTypeTz
	case nexus.IsRollupAddressLazy(address):
		return AccountTypeRollup
	case nexus.IsSmartRollupAddressLazy(address):
		return AccountTypeSmartRollup
	default:
		return AccountTypeUnknown
	}
}

// String -
func (typ AccountType) String() string {
	switch typ {
	case AccountTypeContract:
		return "contract"
	case AccountTypeRollup:
		return "rollup"
	case AccountTypeSmartRollup:
		return "smart_rollup"
	case AccountTypeTz:
		return "account"
	default:
		return "unknown"
	}
}
