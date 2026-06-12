package migrations

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/mavryk-network/nexushub/internal/models"
	modelsContract "github.com/mavryk-network/nexushub/internal/models/contract"
	"github.com/mavryk-network/nexushub/internal/models/migration"
	"github.com/mavryk-network/nexushub/internal/models/protocol"
	"github.com/mavryk-network/nexushub/internal/models/types"
	"github.com/mavryk-network/nexushub/internal/nexus"
	"github.com/mavryk-network/nexushub/internal/nexus/contract"
	"github.com/mavryk-network/nexushub/internal/noderpc"
)

// Alpha -
type Alpha struct{}

// NewAlpha -
func NewAlpha() *Alpha {
	return &Alpha{}
}

// Parse -
func (p *Alpha) Parse(ctx context.Context, script noderpc.Script, old *modelsContract.Contract, previous, next protocol.Protocol, timestamp time.Time, tx models.Transaction) error {
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

	var s nexus.RawScript
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

	old.AlphaID = contractScript.ID

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
func (p *Alpha) IsMigratable(address string) bool {
	return true
}
