package graphBetaNetworkMCPPolicy

import "errors"

var (
	errInvalidResponse = errors.New("invalid MCP policy response")
	errEmptyResponse   = errors.New("MCP policy response is empty")
)
