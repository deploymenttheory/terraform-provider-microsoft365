package graphBetaNetworkMCPPolicyRule

import "errors"

var (
	errInvalidResponse   = errors.New("invalid MCP policy rule response")
	errEmptyResponse     = errors.New("MCP policy rule response is empty")
	errInvalidConditions = errors.New("invalid MCP matching conditions")
)
