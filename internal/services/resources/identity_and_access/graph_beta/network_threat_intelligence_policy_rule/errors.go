package graphBetaNetworkThreatIntelligencePolicyRule

import "errors"

var (
	errInvalidResponse = errors.New("invalid threat intelligence response")
	errInvalidRequest  = errors.New("invalid threat intelligence request")
	// errUnknownRuleStatus reports a rule status other than enabled or disabled.
	errUnknownRuleStatus = errors.New("unsupported observed rule status")
)
