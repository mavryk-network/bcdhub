package handlers

import (
	"net/http"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/database"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/gin-gonic/gin"
)

// GetDAppList -
func (ctx *Context) GetDAppList(c *gin.Context) {
	dapps, err := ctx.DB.GetDApps()
	if handleError(c, err, 0) {
		return
	}

	results := make([]DApp, len(dapps))
	for i := range dapps {
		result, err := ctx.appendDAppInfo(dapps[i], false)
		if handleError(c, err, 0) {
			return
		}
		results[i] = result
	}

	c.JSON(http.StatusOK, results)
}

// GetDApp -
func (ctx *Context) GetDApp(c *gin.Context) {
	var req getDappRequest
	if err := c.BindUri(&req); handleError(c, err, http.StatusBadRequest) {
		return
	}

	dapp, err := ctx.DB.GetDAppBySlug(req.Slug)
	if handleError(c, err, 0) {
		return
	}

	response, err := ctx.appendDAppInfo(dapp, true)
	if handleError(c, err, 0) {
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctx *Context) appendDAppInfo(dapp database.DApp, withDetails bool) (DApp, error) {
	result := DApp{
		Name:              dapp.Name,
		ShortDescription:  dapp.ShortDescription,
		FullDescription:   dapp.FullDescription,
		Version:           dapp.Version,
		License:           dapp.License,
		WebSite:           dapp.WebSite,
		Slug:              dapp.Slug,
		AgoraReviewPostID: dapp.AgoraReviewPostID,
		AgoraQAPostID:     dapp.AgoraQAPostID,
		Authors:           dapp.Authors,
		SocialLinks:       dapp.SocialLinks,
		Interfaces:        dapp.Interfaces,
		Categories:        dapp.Categories,
		Soon:              dapp.Soon,
	}

	if len(dapp.Pictures) > 0 {
		screenshots := make([]Screenshot, 0)
		for _, pic := range dapp.Pictures {
			switch pic.Type {
			case "logo":
				result.Logo = pic.Link
			case "cover":
				result.Cover = pic.Link
			default:
				screenshots = append(screenshots, Screenshot{
					Type: pic.Type,
					Link: pic.Link,
				})
			}
		}

		result.Screenshots = screenshots
	}

	if withDetails {
		if len(dapp.DexTokens) > 0 {
			result.DexTokens = make([]models.TokenMetadata, 0)
			for _, token := range dapp.DexTokens {
				tokenMetadata, err := ctx.ES.GetTokenMetadata(elastic.GetTokenMetadataContext{
					Contract: token.Contract,
					Network:  consts.Mainnet,
					TokenID:  int64(token.TokenID),
				})
				if err != nil {
					if elastic.IsRecordNotFound(err) {
						continue
					}
					return result, err
				}
				result.DexTokens = append(result.DexTokens, tokenMetadata...)
			}
		}

		if len(dapp.Contracts) > 0 {
			result.Contracts = make([]DAppContract, 0)

			for _, address := range dapp.Contracts {
				contract, err := ctx.ES.GetContract(map[string]interface{}{
					"network": consts.Mainnet,
					"address": address,
				})
				if err != nil {
					return result, err
				}
				result.Contracts = append(result.Contracts, DAppContract{
					Network:     contract.Network,
					Address:     contract.Address,
					Alias:       contract.Alias,
					ReleaseDate: contract.Timestamp.UTC(),
				})

				tokens, err := ctx.getTokens(consts.Mainnet, address)
				if err != nil {
					if elastic.IsRecordNotFound(err) {
						continue
					}
					return result, err
				}
				result.Tokens = append(result.Tokens, tokens...)
			}
		}
	}

	return result, nil
}