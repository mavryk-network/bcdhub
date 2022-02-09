package contract_metadata

import "time"

// ParserConfig -
type ParserConfig struct {
	IPFSGateways []string
	SharePath    string
	HTTPTimeout  time.Duration
}
