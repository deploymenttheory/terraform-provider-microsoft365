package mocks

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_tls_inspection_policy/mocks"
	"github.com/jarcoal/httpmock"
)

// TLSInspectionPolicyRuleMock shares the policy mock's nested collection and cascade behavior.
type TLSInspectionPolicyRuleMock struct {
	policyMocks.TLSInspectionPolicyMock
}

var _ mocks.MockRegistrar = (*TLSInspectionPolicyRuleMock)(nil)

func init() {
	mocks.GlobalRegistry.Register("network_tls_inspection_policy_rule", &TLSInspectionPolicyRuleMock{})
}

// RegisterErrorMocks injects a synthetic rule-creation validation error while leaving policy setup intact.
func (m *TLSInspectionPolicyRuleMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/networkaccess/tlsInspectionPolicies/[^/]+/policyRules$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Synthetic TLS inspection rule validation error"}}`))
}
