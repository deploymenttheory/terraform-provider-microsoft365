package graphBetaNetworkPromptPolicyRule

import "errors"

var (
	errInvalidResponse = errors.New("invalid prompt policy rule response")
	errEmptyResponse   = errors.New("prompt policy rule response is empty")
	errInvalidSchemes  = errors.New("invalid conversation schemes")
)
