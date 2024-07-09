package bcd

import (
	"github.com/pkg/errors"
)

// This is the list of protocols BCD supports
// Every time new protocol is proposed we determine if everything works fine or implement a custom handler otherwise
// After that we append protocol to this list with a corresponding handler id (aka symlink)
var symLinks = map[string]string{
	"ProtoGenesisGenesisGenesisGenesisGenesisGenesk612im": SymLinkAlpha,
	"ProtoDemoNoopsDemoNoopsDemoNoopsDemoNoopsDemo6XBoYp": SymLinkAlpha,
	"Ps9mPmXaRzmzk35gbAYNCAw6UXdE2qoABTHbN2oEEc1qM7CwT9P": SymLinkAlpha,
	"PtAtLasomUEW99aVhVTrqjCHjJSpFUa8uHNEAEamx9v2SNeTaNp": SymLinkJakarta, // Atlas
	"PsParisCZo7KAh1Z1smVd9ZMZ1HHn5gkzbM94V3PLCpknFWhUAi": SymLinkJakarta, // Boreas
}

// GetProtoSymLink -
func GetProtoSymLink(protocol string) (string, error) {
	if protoSymLink, ok := symLinks[protocol]; ok {
		return protoSymLink, nil
	}
	return "", errors.Errorf("Unknown protocol: %s", protocol)
}

// GetCurrentProtocol - returns last supported protocol
func GetCurrentProtocol() string {
	return "PsParisCZo7KAh1Z1smVd9ZMZ1HHn5gkzbM94V3PLCpknFWhUAi"
}

// SymLink - returns last sym link
func SymLink() (string, error) {
	return GetProtoSymLink(GetCurrentProtocol())
}

// Symbolic links
const (
	SymLinkAlpha   = "alpha"
	SymLinkBabylon = "babylon"
	SymLinkJakarta = "jakarta"
)

var ChainID = map[string]string{
	"NetXdQprcVkpaWU": "mainnet",
	"NetXvyTAafh8goH": "atlasnet",
	"NetXR64bNAYkP4S": "boreasnet",
}
