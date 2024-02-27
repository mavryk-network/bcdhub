package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/config"
)

// ListEvents -
// @Summary List contract events
// @Description List contract events
// @Tags operations
// @ID list-events
// @Param network path string true "Network"
// @Param address path string true "KT address" minlength(36) maxlength(36)
// @Param offset query string false "Offset"
// @Param size query integer false "Expected events count" mininum(1) maximum(10)
// @Accept  json
// @Produce  json
// @Success 200 {array} Event
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /v1/contract/{network}/{address}/events [get]
func ListEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)

		var req getAccountRequest
		if err := c.ShouldBindUri(&req); handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		var page pageableRequest
		if err := c.ShouldBindQuery(&page); handleError(c, ctx.Storage, err, http.StatusBadRequest) {
			return
		}

		account, err := ctx.Accounts.Get(c.Request.Context(), req.Address)
		if handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		operations, err := ctx.Operations.ListEvents(c.Request.Context(), account.ID, page.Size, page.Offset)
		if handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		events := make([]Event, 0)
		for i := range operations {
			e, err := NewEvent(operations[i])
			if handleError(c, ctx.Storage, err, http.StatusNotFound) {
				return
			}
			if e == nil {
				continue
			}
			events = append(events, *e)
		}
		c.SecureJSON(http.StatusOK, events)
	}
}
