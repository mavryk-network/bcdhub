package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/models/types"
)

// ContractsHelpers -
func ContractsHelpers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)

		var req getByNetwork
		if err := c.ShouldBindUri(&req); handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		var args findContract
		if err := c.ShouldBindQuery(&args); handleError(c, ctx.Storage, err, http.StatusBadRequest) {
			return
		}

		splitted := strings.Split(args.Tags, ",")
		tags := types.NewTags(splitted)
		contract, err := ctx.Contracts.FindOne(c.Request.Context(), tags)
		if handleError(c, ctx.Storage, err, http.StatusNotFound) {
			return
		}

		response, err := contractPostprocessing(ctx, contract)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		c.SecureJSON(http.StatusOK, response)
	}
}

func getStringPointer(s string) *string {
	return &s
}
