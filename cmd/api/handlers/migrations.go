package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/bcd/encoding"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/models/migration"
)

// GetContractMigrations godoc
// @Summary Get contract migrations
// @Description Get contract migrations
// @Tags contract
// @ID get-contract-migrations
// @Param network path string true "Network"
// @Param address path string true "KT address" minlength(36) maxlength(36)
// @Accept json
// @Produce json
// @Success 200 {array} Migration
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /v1/contract/{network}/{address}/migrations [get]
func GetContractMigrations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)

		var req getContractRequest
		if err := c.ShouldBindUri(&req); handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		contract, err := ctx.Contracts.Get(c.Request.Context(), req.Address)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		migrations, err := ctx.Migrations.Get(c.Request.Context(), contract.ID)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		result, err := prepareMigrations(c.Request.Context(), ctx, migrations)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		c.SecureJSON(http.StatusOK, result)
	}
}

func prepareMigrations(c context.Context, ctx *config.Context, data []migration.Migration) ([]Migration, error) {
	result := make([]Migration, len(data))
	for i := range data {
		proto, err := ctx.Cache.ProtocolByID(c, data[i].ProtocolID)
		if err != nil && !ctx.Storage.IsRecordNotFound(err) {
			return nil, err
		}
		prevProto, err := ctx.Cache.ProtocolByID(c, data[i].PrevProtocolID)
		if err != nil && !ctx.Storage.IsRecordNotFound(err) {
			return nil, err
		}
		var hash string
		if len(data[i].Hash) > 0 {
			hash = encoding.MustEncodeOperationHash(data[i].Hash)
		}
		result[i] = Migration{
			Level:        data[i].Level,
			Timestamp:    data[i].Timestamp,
			Hash:         hash,
			Protocol:     proto.Hash,
			PrevProtocol: prevProto.Hash,
			Kind:         data[i].Kind.String(),
		}
	}
	return result, nil
}
