package protocols

import (
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/parsers/contract"
	"github.com/mavryk-network/bcdhub/internal/parsers/migrations"
	"github.com/mavryk-network/bcdhub/internal/parsers/storage"
	"github.com/pkg/errors"
)

// Specific -
type Specific struct {
	StorageParser         storage.Parser
	ContractParser        contract.Parser
	MigrationParser       migrations.Parser
	NeedReceiveRawStorage bool
}

// Get -
func Get(ctx *config.Context, protocol string) (*Specific, error) {
	switch protocol {
	case "ProtoGenesisGenesisGenesisGenesisGenesisGenesk612im",
		"ProtoDemoNoopsDemoNoopsDemoNoopsDemoNoopsDemo6XBoYp",
		"PsUCFkqUrQ614xKsFEAf4AamoUXTAG4ygjMpFzsgEdKr3PGYreP":
		return &Specific{
			StorageParser:         storage.NewAlpha(),
			ContractParser:        contract.NewAlpha(ctx),
			MigrationParser:       migrations.NewAlpha(),
			NeedReceiveRawStorage: false,
		}, nil
	case "ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
		"PtAtLasdzXg4XxeVNtWheo13nG4wHXP22qYMqFcT3fyBpWkFero",
		"PtBoreasK2KPuKbeYtXeEdudEHS7YcMFHE9amwheUc4kejTxgRi":
		return &Specific{
			StorageParser:         storage.NewLazyBabylon(ctx.BigMapDiffs, ctx.Operations, ctx.Accounts),
			ContractParser:        contract.NewJakarta(ctx),
			MigrationParser:       migrations.NewJakarta(),
			NeedReceiveRawStorage: true,
		}, nil

	default:
		return nil, errors.Errorf("unknown protocol in parser's creation: %s", protocol)

	}
}

// NeedImplicitParsing -
func NeedImplicitParsing(protocol string) bool {
	switch protocol {
	case "ProtoGenesisGenesisGenesisGenesisGenesisGenesk612im",
		"ProtoDemoNoopsDemoNoopsDemoNoopsDemoNoopsDemo6XBoYp",
		"PsUCFkqUrQ614xKsFEAf4AamoUXTAG4ygjMpFzsgEdKr3PGYreP":
		return false
	case "ProtoALphaALphaALphaALphaALphaALphaALphaALphaDdp3zK",
		"PtAtLasdzXg4XxeVNtWheo13nG4wHXP22qYMqFcT3fyBpWkFero",
		"PtBoreasK2KPuKbeYtXeEdudEHS7YcMFHE9amwheUc4kejTxgRi":
		return true
	}
	return false
}
