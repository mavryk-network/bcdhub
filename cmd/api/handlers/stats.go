package handlers

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/config"
	"github.com/mavryk-network/bcdhub/internal/models/types"
)

// GetStats godoc
// @Summary Show indexer stats
// @Description get indexer states for all networks
// @Tags statistics
// @ID get-stats
// @Accept  json
// @Produce  json
// @Success 200 {array} Block
// @Failure 500 {object} Error
// @Router /v1/stats [get]
func GetStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxs := c.MustGet("contexts").(config.Contexts)

		networks := make(types.Networks, 0)
		for n := range ctxs {
			networks = append(networks, n)
		}

		sort.Sort(networks)

		blocks := make([]Block, 0)
		for _, network := range networks {
			last, err := ctxs[network].Blocks.Last(c.Request.Context())
			if err != nil {
				if ctxs[network].Storage.IsRecordNotFound(err) {
					continue
				}
				handleError(c, ctxs[network].Storage, err, 0)
				return
			}
			var block Block
			block.FromModel(last)
			predecessor, err := ctxs[network].Blocks.Get(c.Request.Context(), last.Level)
			if err != nil && !ctxs[network].Storage.IsRecordNotFound(err) {
				handleError(c, ctxs[network].Storage, err, 0)
				return
			}
			block.Network = network.String()
			block.Predecessor = predecessor.Hash

			stats, err := ctxs[network].Stats.Get(c.Request.Context())
			if err != nil && !ctxs[network].Storage.IsRecordNotFound(err) {
				handleError(c, ctxs[network].Storage, err, 0)
				return
			}
			block.Stats = NewStats(stats)

			blocks = append(blocks, block)
		}

		c.SecureJSON(http.StatusOK, blocks)
	}
}

// RecentlyCalledContracts godoc
// @Summary Show recently called contracts
// @Description Show recently called contracts
// @Tags statistics
// @ID get-recenly-called-contracts
// @Param network path string true "Network"
// @Param size query integer false "Contracts count" mininum(1) maximum(10)
// @Param offset query integer false "Offset" mininum(1)
// @Accept  json
// @Produce  json
// @Success 200 {array} RecentlyCalledContract
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/stats/{network}/recently_called_contracts [get]
func RecentlyCalledContracts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)
		var req getByNetwork
		if err := c.ShouldBindUri(&req); handleError(c, ctx.Storage, err, http.StatusBadRequest) {
			return
		}

		var page pageableRequest
		if err := c.ShouldBindQuery(&page); handleError(c, ctx.Storage, err, http.StatusBadRequest) {
			return
		}

		if page.Size > 10 || page.Size == 0 {
			page.Size = 10
		}

		contracts, err := ctx.Accounts.RecentlyCalledContracts(c.Request.Context(), page.Offset, page.Size)
		if handleError(c, ctx.Storage, err, 0) {
			return
		}

		response := make([]RecentlyCalledContract, len(contracts))
		for i := range contracts {
			var res RecentlyCalledContract
			res.FromModel(contracts[i])
			response[i] = res
		}

		c.SecureJSON(http.StatusOK, response)
	}
}

// ContractsCount godoc
// @Summary Get contracts count
// @Description Get contracts count
// @Tags statistics
// @ID get-contracts-count
// @Param network path string true "Network"
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /v1/stats/{network}/contracts_count [get]
func ContractsCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.MustGet("context").(*config.Context)
		stats, err := ctx.Stats.Get(c.Request.Context())
		if handleError(c, ctx.Storage, err, 0) {
			return
		}
		c.SecureJSON(http.StatusOK, stats.ContractsCount)
	}
}
