package miguel

import (
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/formatter"
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/meta"
	"github.com/tidwall/gjson"
)

type lambdaDecoder struct{}

// Decode -
func (l *lambdaDecoder) Decode(node gjson.Result, path string, nm *meta.NodeMetadata, metadata meta.Metadata) (interface{}, error) {
	val, err := formatter.MichelineToMichelson(node, true)
	return map[string]interface{}{
		"value": val,
		"type":  nm.Type,
	}, err
}
