package mocks

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_mcp_policy/mocks"
	"github.com/jarcoal/httpmock"
)

// MCPPolicyRuleMock shares the policy mock's nested collection and cascade behavior.
type MCPPolicyRuleMock struct {
	policyMocks.MCPPolicyMock
}

var _ mocks.MockRegistrar = (*MCPPolicyRuleMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_mcp_policy_rule", &MCPPolicyRuleMock{})
}

// RegisterErrorMocks injects a synthetic rule-creation validation error while leaving policy setup intact.
func (m *MCPPolicyRuleMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/networkaccess/mcpPolicies/[^/]+/policyRules$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Synthetic MCP rule validation error"}}`))
}
