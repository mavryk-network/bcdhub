package types

import (
	"strconv"

	"github.com/pkg/errors"
)

// Network -
type Network int64

// Network names
const (
	Empty Network = iota
	Mainnet
	Sandboxnet
	Basenet
	Mondaynet
	Dailynet
	Rollupnet
	Atlasnet
	Weeklynet
)

var networkNames = map[Network]string{
	Mainnet:    "mainnet",
	Sandboxnet: "sandboxnet",
	Basenet:    "basenet",
	Mondaynet:  "mondaynet",
	Weeklynet:  "weeklynet",
	Dailynet:   "dailynet",
	Rollupnet:  "rollupnet",
	Atlasnet:   "atlasnet",
}

var namesToNetwork = map[string]Network{
	"mainnet":    Mainnet,
	"sandboxnet": Sandboxnet,
	"basenet":    Basenet,
	"mondaynet":  Mondaynet,
	"dailynet":   Dailynet,
	"rollupnet":  Rollupnet,
	"atlasnet":   Atlasnet,
	"weeklynet":  Weeklynet,
}

// String - convert enum to string for printing
func (network Network) String() string {
	return networkNames[network]
}

// UnmarshalJSON -
func (network *Network) UnmarshalJSON(data []byte) error {
	name, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	newValue, ok := namesToNetwork[name]
	if !ok {
		return errors.Errorf("Unknown network: %d", network)
	}

	*network = newValue
	return nil
}

// MarshalJSON -
func (network Network) MarshalJSON() ([]byte, error) {
	name, ok := networkNames[network]
	if !ok {
		return nil, errors.Errorf("Unknown network: %d", network)
	}

	return []byte(strconv.Quote(name)), nil
}

// NewNetwork -
func NewNetwork(name string) Network {
	return namesToNetwork[name]
}

// Networks -
type Networks []Network

func (n Networks) Len() int           { return len(n) }
func (n Networks) Less(i, j int) bool { return n[i] < n[j] }
func (n Networks) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
