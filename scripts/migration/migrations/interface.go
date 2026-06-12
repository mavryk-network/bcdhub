package migrations

import "github.com/mavryk-network/nexushub/internal/config"

// Migration - intreface need to realize for migrate
type Migration interface {
	Do(ctx *config.Context) error
	Key() string
	Description() string
}
