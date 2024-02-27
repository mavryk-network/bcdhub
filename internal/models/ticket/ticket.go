package ticket

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/mavryk-network/bcdhub/internal/models/account"
	"github.com/uptrace/bun"
)

type Ticket struct {
	bun.BaseModel `bun:"tickets"`

	ID           int64  `bun:"id,pk,notnull,autoincrement"`
	Level        int64  `bun:"level"`
	TicketerID   int64  `bun:"ticketer_id"`
	ContentType  []byte `bun:"content_type,type:bytea"`
	Content      []byte `bun:"content,type:bytea"`
	UpdatesCount int    `bun:"updates_count"`
	Hash         string `bun:"hash,unique:ticket_key"`

	Ticketer account.Account `bun:"rel:belongs-to"`
}

func (t Ticket) GetID() int64 {
	return t.ID
}

func (Ticket) TableName() string {
	return "tickets"
}

func (t Ticket) GetHash() string {
	data := make([]byte, len(t.ContentType))
	copy(data, t.ContentType)
	data = append(data, t.Content...)
	data = append(data, []byte(t.Ticketer.Address)...)
	h := sha256.New()
	_, _ = h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// LogFields -
func (t Ticket) LogFields() map[string]interface{} {
	return map[string]interface{}{
		"ticketer_id":  t.TicketerID,
		"content":      string(t.Content),
		"content_type": string(t.ContentType),
	}
}
