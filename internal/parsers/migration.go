package parsers

import (
	"encoding/json"
	"time"

	"github.com/baking-bad/bcdhub/internal/contractparser"
	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/contractparser/meta"
	"github.com/baking-bad/bcdhub/internal/contractparser/storage"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/helpers"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// MigrationParser -
type MigrationParser struct {
	es             elastic.IElastic
	filesDirectory string
}

// NewMigrationParser -
func NewMigrationParser(es elastic.IElastic, filesDirectory string) *MigrationParser {
	return &MigrationParser{
		es:             es,
		filesDirectory: filesDirectory,
	}
}

// Parse -
func (p *MigrationParser) Parse(script gjson.Result, old models.Contract, previous, next models.Protocol, timestamp time.Time) ([]elastic.Model, []elastic.Model, error) {
	metadata := models.Metadata{ID: old.Address}
	if err := p.es.GetByID(&metadata); err != nil {
		return nil, nil, err
	}

	if err := NewMetadataParser(next.SymLink).UpdateMetadata(script, old.Address, &metadata); err != nil {
		return nil, nil, err
	}

	var updates []elastic.Model
	if previous.SymLink == "alpha" {
		newUpdates, err := p.getUpdates(script, old, next, metadata)
		if err != nil {
			return nil, nil, err
		}
		updates = newUpdates
	}

	newHash, err := contractparser.ComputeContractHash(script.Get("code").Raw)
	if err != nil {
		return nil, nil, err
	}
	if newHash == old.Hash {
		return []elastic.Model{&metadata}, updates, nil
	}

	migration := models.Migration{
		ID:          helpers.GenerateID(),
		IndexedTime: time.Now().UnixNano() / 1000,

		Network:      old.Network,
		Level:        previous.EndLevel,
		Protocol:     next.Hash,
		PrevProtocol: previous.Hash,
		Address:      old.Address,
		Timestamp:    timestamp,
		Kind:         consts.MigrationUpdate,
	}

	return []elastic.Model{&metadata, &migration}, updates, nil
}

func (p *MigrationParser) getUpdates(script gjson.Result, contract models.Contract, protocol models.Protocol, metadata models.Metadata) ([]elastic.Model, error) {
	stringMetadata, ok := metadata.Storage[protocol.SymLink]
	if !ok {
		return nil, errors.Errorf("[MigrationParser.getUpdates] Unknown metadata sym link: %s", protocol.SymLink)
	}

	var m meta.Metadata
	if err := json.Unmarshal([]byte(stringMetadata), &m); err != nil {
		return nil, err
	}

	storageJSON := script.Get("storage")
	newMapPtr, err := storage.FindBigMapPointers(m, storageJSON)
	if err != nil {
		return nil, err
	}
	if len(newMapPtr) != 1 {
		return nil, nil
	}
	var newPath string
	var newPtr int64
	for p, b := range newMapPtr {
		newPath = b
		newPtr = p
	}

	bmd, err := p.es.GetBigMapsForAddress(contract.Network, contract.Address)
	if err != nil {
		return nil, err
	}
	if len(bmd) == 0 {
		return nil, nil
	}

	updates := make([]elastic.Model, len(bmd))
	for i := range bmd {
		bmd[i].BinPath = newPath
		bmd[i].Ptr = newPtr
		updates[i] = &bmd[i]
	}
	return updates, nil
}