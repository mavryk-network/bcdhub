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
	"PsUCFkqUrQ614xKsFEAf4AamoUXTAG4ygjMpFzsgEdKr3PGYreP": SymLinkAlpha,
	"PtAtLasdzXg4XxeVNtWheo13nG4wHXP22qYMqFcT3fyBpWkFero": SymLinkAtlas, // Atlas
	"PtBzwViMCC1gfm98y5TDKqz2e3vjBXPAUoWu7jfEcN6yj2ZhCyT": SymLinkAtlas, // Boreas
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
	return "PtAtLasdzXg4XxeVNtWheo13nG4wHXP22qYMqFcT3fyBpWkFero"
}

// SymLink - returns last sym link
func SymLink() (string, error) {
	return GetProtoSymLink(GetCurrentProtocol())
}

// Symbolic links
const (
	SymLinkAlpha   = "alpha"
	SymLinkBabylon = "babylon"
	SymLinkAtlas   = "atlas"
)

var ChainID = map[string]string{
	"NetXXAAR1wWQhhe": "mainnet",
	"NetXUrNc8uioxP8": "atlasnet",
	"NetXR64bNAYkP4S": "boreasnet",
}
