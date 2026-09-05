package graphBetaNetworkTLSInspectionPolicy

import "errors"

var (
	errInvalidResponse = errors.New("invalid TLS inspection policy response")
	errEmptyResponse   = errors.New("TLS inspection policy response is empty")
)
