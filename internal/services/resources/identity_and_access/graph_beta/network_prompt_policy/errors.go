package graphBetaNetworkPromptPolicy

import "errors"

var (
	errInvalidResponse = errors.New("invalid prompt policy response")
	errEmptyResponse   = errors.New("prompt policy response is empty")
)
