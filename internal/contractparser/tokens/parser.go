package tokens

import (
	"fmt"
	"strings"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/contractparser/meta"
	"github.com/baking-bad/bcdhub/internal/contractparser/storage"
	"github.com/baking-bad/bcdhub/internal/contractparser/unpack"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/noderpc"
	"github.com/tidwall/gjson"
)

// TokenMetadataParser -
type TokenMetadataParser struct {
	es        elastic.IElastic
	rpc       noderpc.INode
	sharePath string
	network   string

	sources map[string]string
}

// NewTokenMetadataParser -
func NewTokenMetadataParser(es elastic.IElastic, rpc noderpc.INode, sharePath, network string) TokenMetadataParser {
	return TokenMetadataParser{
		es: es, rpc: rpc, sharePath: sharePath, network: network,
		sources: map[string]string{
			"carthagenet": "tz1grSQDByRpnVs7sPtaprNZRp531ZKz6Jmm",
			"mainnet":     "tz2FCNBrERXtaTtNX6iimR1UJ5JSDxvdHM93",
			"dalphanet":   "tz1eTHvnf1WrEXHPhrYFY3EWcSt5XpDA1u97",
			"delphinet":   "tz1ME9SBiGDCzLwgoShUMs2d9zRr23aJHf4w",
		},
	}
}

// Parse -
func (t TokenMetadataParser) Parse(address string, level int64) ([]Metadata, error) {
	state, err := t.getState(level)
	if err != nil {
		return nil, err
	}
	registryAddress, err := t.getTokenMetadataRegistry(address, state)
	if err != nil {
		return nil, err
	}
	return t.parse(address, registryAddress, state)
}

// ParseWithRegistry -
func (t TokenMetadataParser) ParseWithRegistry(address, registry string, level int64) ([]Metadata, error) {
	state, err := t.getState(level)
	if err != nil {
		return nil, err
	}
	return t.parse(address, registry, state)
}

func (t TokenMetadataParser) parse(address, registry string, state models.Block) ([]Metadata, error) {
	ptr, err := t.getBigMapPtr(registry, state)
	if err != nil {
		return nil, err
	}

	bmd, err := t.es.GetBigMapKeys(ptr, t.network, "", 1000, 0)
	if err != nil {
		return nil, err
	}

	metadata := make([]Metadata, len(bmd))
	for i := range bmd {
		value := gjson.Parse(bmd[i].Value)
		m, err := t.parseMetadata(value)
		if err != nil {
			continue
		}
		m.RegistryAddress = registry
		m.Timestamp = bmd[i].Timestamp
		m.Level = bmd[i].Level
		metadata[i] = m
	}

	return metadata, nil
}

func (t TokenMetadataParser) getState(level int64) (models.Block, error) {
	if level > 0 {
		return t.es.GetBlock(t.network, level)
	}
	return t.es.GetLastBlock(t.network)
}

func (t TokenMetadataParser) getTokenMetadataRegistry(address string, state models.Block) (string, error) {
	metadata, err := t.hasTokenMetadataRegistry(address, state.Protocol)
	if err != nil {
		return "", err
	} else if metadata == nil {
		return "", ErrNoTokenMetadataRegistryMethod
	}

	result, err := t.es.SearchByText("view_address", 0, nil, map[string]interface{}{
		"networks": []string{t.network},
		"indices":  []string{elastic.DocContracts},
	}, false)
	if err != nil {
		return "", err
	}
	if result.Count == 0 {
		return "", ErrNoViewAddressContract
	}

	source, ok := t.sources[t.network]
	if !ok {
		return "", ErrUnknownNetwork
	}

	counter, err := t.rpc.GetCounter(source)
	if err != nil {
		return "", err
	}

	protocol, err := t.es.GetProtocol(t.network, "", state.Level)
	if err != nil {
		return "", err
	}

	parameters := gjson.Parse(fmt.Sprintf(`{"entrypoint": "%s", "value": {"string": "%s"}}`, TokenMetadataRegistry, result.Items[0].Value))
	response, err := t.rpc.RunOperation(
		state.ChainID,
		state.Hash,
		source,
		address,
		0,
		protocol.Constants.HardGasLimitPerOperation,
		protocol.Constants.HardStorageLimitPerOperation,
		counter+1,
		0,
		parameters,
	)
	if err != nil {
		return "", err
	}

	registryAddress, err := t.parseRegistryAddress(response)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(address, "KT") {
		return "", ErrInvalidRegistryAddress
	}
	if registryAddress == selfAddress {
		registryAddress = address
	}
	return registryAddress, nil
}

func (t TokenMetadataParser) parseRegistryAddress(response gjson.Result) (string, error) {
	value := response.Get("contents.0.metadata.internal_operation_results.0.parameters.value")
	if value.Exists() {
		if value.Get("bytes").Exists() {
			return unpack.Address(value.Get("bytes").String())
		} else if value.Get("string").Exists() {
			return value.Get("string").String(), nil
		}
	}
	return "", ErrInvalidContractParameter
}

func (t TokenMetadataParser) hasTokenMetadataRegistry(address, protocol string) (meta.Metadata, error) {
	metadata, err := meta.GetMetadata(t.es, address, consts.PARAMETER, protocol)
	if err != nil {
		return nil, err
	}

	for _, node := range metadata {
		if node.Name == TokenMetadataRegistry {
			return metadata, nil
		}
	}

	return nil, nil
}

func (t TokenMetadataParser) getBigMapPtr(address string, state models.Block) (int64, error) {
	resistryStorageMetadata, err := meta.GetMetadata(t.es, address, consts.STORAGE, state.Protocol)
	if err != nil {
		return 0, err
	}

	var bmPath string
	for binPath, node := range resistryStorageMetadata {
		if node.Name == TokenMetadataRegistryStorageKey {
			bmPath = binPath
			break
		}
	}
	if bmPath == "" {
		return 0, ErrNoMetadataKeyInStorage
	}

	registryStorage, err := t.rpc.GetScriptStorageJSON(address, state.Level)
	if err != nil {
		return 0, err
	}

	ptrs, err := storage.FindBigMapPointers(resistryStorageMetadata, registryStorage)
	if err != nil {
		return 0, err
	}
	for ptr, path := range ptrs {
		if path == bmPath {
			return ptr, nil
		}
	}

	return 0, ErrNoMetadataKeyInStorage
}

const (
	keyTokenID  = "args.0.int"
	keySymbol   = "args.1.args.0.string"
	keyName     = "args.1.args.1.args.0.string"
	keyDecimals = "args.1.args.1.args.1.args.0.int"
	keyExtras   = "args.1.args.1.args.1.args.1"
)

func (t TokenMetadataParser) parseMetadata(value gjson.Result) (Metadata, error) {
	extras := make(map[string]interface{})
	for _, item := range value.Get(keyExtras).Array() {
		k := item.Get("args.0.string").String()
		if item.Get("args.1.string").Exists() || item.Get("args.1.bytes").Exists() {
			extras[k] = item.Get("args.1.string").String()
		} else if item.Get("args.1.int").Exists() {
			extras[k] = item.Get("args.1.int").Int()
		}
	}

	if !value.Get(keyTokenID).Exists() {
		return Metadata{}, ErrInvalidStorageStructure
	}
	if !value.Get(keySymbol).Exists() {
		return Metadata{}, ErrInvalidStorageStructure
	}
	if !value.Get(keyName).Exists() {
		return Metadata{}, ErrInvalidStorageStructure
	}
	if !value.Get(keyDecimals).Exists() {
		return Metadata{}, ErrInvalidStorageStructure
	}

	return Metadata{
		TokenID:  value.Get(keyTokenID).Int(),
		Symbol:   value.Get(keySymbol).String(),
		Name:     value.Get(keyName).String(),
		Decimals: value.Get(keyDecimals).Int(),
		Extras:   extras,
	}, nil
}