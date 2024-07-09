package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/bcd"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/rs/zerolog/log"
)

// GetInfo godoc
// @Summary Get account info
// @Description Get account info
// @Tags account
// @ID get-account-info
// @Param network path string true "Network"
// @Param address path string true "Address" minlength(36) maxlength(36)
// @Accept  json
// @Produce  json
// @Success 200 {object} AccountInfo
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /v1/account/{network}/{address} [get]
func GetInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)

		var req getAccountRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, Error{Message: err.Error()})
			return
		}

		acc, err := ctx.Accounts.Get(c.Request.Context(), req.Address)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		var balance int64
		if !(bcd.IsRollupAddressLazy(acc.Address) || bcd.IsSmartRollupAddressLazy(acc.Address)) {
			block, err := ctx.Blocks.Last(c.Request.Context())
			if handleError(c, ctx.Storage, err, 0) {
				return
			}
			balance, err = ctx.Cache.TezosBalance(c, acc.Address, block.Level)
			if err != nil {
				log.Err(err).Msg("receiving tezos balance")
			}
		}
		c.SecureJSON(http.StatusOK, AccountInfo{
			Address:            acc.Address,
			OperationsCount:    acc.OperationsCount,
			EventsCount:        acc.EventsCount,
			MigrationsCount:    acc.MigrationsCount,
			TicketUpdatesCount: acc.TicketUpdatesCount,
			Balance:            balance,
			LastAction:         acc.LastAction.UTC(),
			AccountType:        acc.Type.String(),
		})
	}

}
