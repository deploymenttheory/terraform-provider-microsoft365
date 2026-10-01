package graphBetaAuthenticationFlowsPolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_flows_policy"
	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_flows_policy/mocks"
)

func setupMockEnvironment(t *testing.T) *policymocks.AuthenticationFlowsPolicyMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &policymocks.AuthenticationFlowsPolicyMock{}
	m.RegisterMocks()
	t.Cleanup(httpmock.DeactivateAndReset)
	t.Cleanup(m.CleanupMockState)
	return m
}

func loadUnitTestTerraform(t *testing.T, name string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_" + name + ".tf")
	require.NoError(t, err)
	return config
}

// TestUnitResourceAuthenticationFlowsPolicy_01_Lifecycle covers adoption, both Boolean values, import, and state-only deletion.
func TestUnitResourceAuthenticationFlowsPolicy_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	address := policy.ResourceName + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform(t, "enabled"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("id").HasValue("authenticationFlowsPolicy"),
				check.That(address).Key("self_service_sign_up.is_enabled").HasValue("true"),
				check.That(address).Key("display_name").HasValue("Authentication flows policy"),
			)},
			{
				Config: loadUnitTestTerraform(t, "disabled"),
				Check: check.That(address).
					Key("self_service_sign_up.is_enabled").
					HasValue("false"),
			},
			{
				ResourceName:      address,
				ImportState:       true,
				ImportStateId:     "authenticationFlowsPolicy",
				ImportStateVerify: true,
			},
			{Config: loadUnitTestTerraform(t, "disabled"), PlanOnly: true},
			{Config: loadUnitTestTerraform(t, "disabled"), Destroy: true},
		},
	})
	enabled, patches := m.Snapshot()
	require.False(t, enabled)
	require.Equal(t, 2, patches)
	require.Zero(t, httpmock.GetCallCountInfo()["DELETE "+policymocks.URL])
	require.Zero(t, httpmock.GetCallCountInfo()["POST "+policymocks.URL])
}
