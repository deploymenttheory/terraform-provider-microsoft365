package graphBetaNetworkTLSInspectionPolicyRule

import "errors"

var (
	errInvalidResponse     = errors.New("invalid TLS inspection policy rule response")
	errEmptyResponse       = errors.New("TLS inspection policy rule response is empty")
	errInvalidDestinations = errors.New("invalid TLS inspection destinations")
	errSystemRule          = errors.New(
		"system TLS inspection rules cannot be managed by this resource",
	)
)
