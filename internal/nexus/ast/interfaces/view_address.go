package interfaces

import "github.com/mavryk-network/nexushub/internal/nexus/consts"

// ViewAddress -
type ViewAddress struct{}

// GetName -
func (f *ViewAddress) GetName() string {
	return consts.ViewAddressTag
}

// GetContractInterface -
func (f *ViewAddress) GetContractInterface() string {
	return `{
		"entrypoints": {
			"default": {
				"prim": "address"
			}
		},
		"is_root": true
	}`
}
