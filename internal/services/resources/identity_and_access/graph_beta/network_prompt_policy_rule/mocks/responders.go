package mocks

import (
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_prompt_policy/mocks"
)

// PromptPolicyRuleMock shares the policy mock's nested collection and cascade behavior.
type PromptPolicyRuleMock struct {
	policyMocks.PromptPolicyMock
}

var _ mocks.MockRegistrar = (*PromptPolicyRuleMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_prompt_policy_rule", &PromptPolicyRuleMock{})
}

// RegisterErrorMocks injects a synthetic rule-creation validation error while leaving policy setup intact.
func (m *PromptPolicyRuleMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		`=~^https://graph\.microsoft\.com/beta/networkaccess/promptPolicies/[^/]+/policyRules$`,
		httpmock.NewStringResponder(
			400,
			`{"error":{"code":"BadRequest","message":"Synthetic prompt rule validation error"}}`,
		),
	)
}
