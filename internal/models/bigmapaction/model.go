package bigmapaction

import (
	"time"

	"github.com/mavryk-network/bcdhub/internal/models/types"
	"github.com/uptrace/bun"
)

// BigMapAction -
type BigMapAction struct {
	bun.BaseModel `bun:"big_map_actions"`

	ID             int64              `bun:"id,pk,notnull,autoincrement"`
	Action         types.BigMapAction `bun:"action,type:SMALLINT"`
	SourcePtr      *int64
	DestinationPtr *int64
	OperationID    int64
	Level          int64
	Address        string    `bun:"address,type:text"`
	Timestamp      time.Time `bun:"timestamp,pk,notnull"`
}

// GetID -
func (b *BigMapAction) GetID() int64 {
	return b.ID
}

func (b *BigMapAction) TableName() string {
	return "big_map_actions"
}
