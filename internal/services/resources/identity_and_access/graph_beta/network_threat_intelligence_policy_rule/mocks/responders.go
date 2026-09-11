package mocks

import (
	common "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_threat_intelligence_policy/mocks"
)

// ThreatIntelligencePolicyRuleMock shares parent and child state for cascade and replacement checks.
type ThreatIntelligencePolicyRuleMock = policy.ThreatIntelligencePolicyMock

func init() {
	common.GlobalRegistry.Register(
		"network_threat_intelligence_policy_rule",
		&ThreatIntelligencePolicyRuleMock{},
	)
}
