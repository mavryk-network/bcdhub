package migrations

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/bcd/contract"
	"github.com/mavryk-network/bcdhub/internal/models"
	modelsContract "github.com/mavryk-network/bcdhub/internal/models/contract"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
	"github.com/mavryk-network/bcdhub/internal/models/protocol"
	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/mavryk-network/bcdhub/internal/noderpc"
)

// Carthage -
type Carthage struct {
}

// NewCarthage -
func NewCarthage() *Carthage {
	return &Carthage{}
}

// Parse -
func (p *Carthage) Parse(ctx context.Context, script noderpc.Script, old *modelsContract.Contract, previous, next protocol.Protocol, timestamp time.Time, tx models.Transaction) error {
	codeBytes, err := json.Marshal(script.Code)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, codeBytes); err != nil {
		return err
	}

	newHash, err := contract.ComputeHash(buf.Bytes())
	if err != nil {
		return err
	}

	var s bcd.RawScript
	if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
		return err
	}

	contractScript := modelsContract.Script{
		Level:     next.StartLevel,
		Hash:      newHash,
		Code:      s.Code,
		Storage:   s.Storage,
		Parameter: s.Parameter,
		Views:     s.Views,
	}

	if err := tx.Scripts(ctx, &contractScript); err != nil {
		return err
	}

	old.BabylonID = contractScript.ID

	m := &migration.Migration{
		ContractID:     old.ID,
		Contract:       *old,
		Level:          next.StartLevel,
		ProtocolID:     next.ID,
		PrevProtocolID: previous.ID,
		Timestamp:      timestamp,
		Kind:           types.MigrationKindUpdate,
	}

	return tx.Migrations(ctx, m)
}

// IsMigratable -
func (p *Carthage) IsMigratable(address string) bool {
	return true
}
