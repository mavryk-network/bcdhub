package migrations

import (
	"fmt"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/schollz/progressbar/v3"
)

// SetTotalWithdrawn - migration that set total_withdrawn to contracts in all networks
type SetTotalWithdrawn struct {
	Network string
}

// Description -
func (m *SetTotalWithdrawn) Description() string {
	return "set total_withdrawn to contracts in all networks"
}

// Do - migrate function
func (m *SetTotalWithdrawn) Do(ctx *Context) error {
	for _, network := range []string{consts.Mainnet, consts.Zeronet, consts.Carthage, consts.Babylon} {
		filter := make(map[string]interface{})
		filter["network"] = network

		contracts, err := ctx.ES.GetContracts(filter)
		if err != nil {
			return err
		}

		logger.Info("Found %d contracts in %s", len(contracts), network)

		bar := progressbar.NewOptions(len(contracts), progressbar.OptionSetPredictTime(false))

		for i, c := range contracts {
			bar.Add(1)

			totalWithdrawn, err := ctx.ES.GetContractWithdrawn(c.Address, c.Network)
			if err != nil {
				fmt.Print("\033[2K\r")
				return err
			}

			fmt.Println("total withdrawn:", totalWithdrawn)

			c.TotalWithdrawn = totalWithdrawn

			if _, err := ctx.ES.UpdateDoc(elastic.DocContracts, contracts[i].ID, contracts[i]); err != nil {
				fmt.Print("\033[2K\r")
				return err
			}
		}

		fmt.Print("\033[2K\r")
		logger.Info("[%s] done. Total contracts: %d", network, len(contracts))
	}
	return nil
}