package mocks

import policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy/mocks"

// CloudFirewallPolicyRuleMock shares parent and child state to exercise dependency ordering and cascade deletion.
type CloudFirewallPolicyRuleMock = policyMocks.CloudFirewallMock
