package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavryk-network/bcdhub/internal/config"
)

// GetConfig -
func GetConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxs := c.MustGet("contexts").(config.Contexts)
		ctx := ctxs.Any()

		cfg := ConfigResponse{
			Networks:       ctx.Config.API.Networks,
			RPCEndpoints:   ctx.Config.API.Frontend.RPC,
			GaEnabled:      ctx.Config.API.Frontend.GaEnabled,
			MempoolEnabled: ctx.Config.API.Frontend.MempoolEnabled,
			SandboxMode:    ctx.Config.API.Frontend.SandboxMode,
		}

		if ctx.Config.API.SentryEnabled {
			cfg.SentryDSN = ctx.Config.Sentry.FrontURI
		}

		c.SecureJSON(http.StatusOK, cfg)
	}
}
