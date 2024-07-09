package types

import (
	stdJSON "encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/mavryk-network/bcdhub/internal/bcd/consts"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Parameters -
type Parameters struct {
	Entrypoint string             `json:"entrypoint"`
	Value      stdJSON.RawMessage `json:"value"`
}

type params struct {
	Entrypoint *string            `json:"entrypoint,omitempty"`
	Value      stdJSON.RawMessage `json:"value,omitempty"`
}

// NewParameters -
func NewParameters(data []byte) *Parameters {
	var p params
	if err := json.Unmarshal(data, &p); err != nil || p.Entrypoint == nil {
		return &Parameters{
			Entrypoint: consts.DefaultEntrypoint,
			Value:      data,
		}
	}
	return &Parameters{
		Entrypoint: *p.Entrypoint,
		Value:      p.Value,
	}
}
