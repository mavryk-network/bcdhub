package core

import (
	"context"

	"github.com/mavryk-network/bcdhub/internal/models"
)

// GetByID -
func (p *Postgres) GetByID(ctx context.Context, output models.Model) error {
	err := p.DB.NewSelect().Model(output).Where("id = ?", output.GetID()).Scan(ctx)
	return err
}
