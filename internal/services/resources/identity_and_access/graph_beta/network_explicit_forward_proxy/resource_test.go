package graphBetaNetworkExplicitForwardProxy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_explicit_forward_proxy"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_explicit_forward_proxy/mocks"
)

const resourceType = target.ResourceName

func TestUnitResourceNetworkExplicitForwardProxy_01_Lifecycle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &testmocks.Mock{}
	m.RegisterMocks()
	t.Cleanup(func() { httpmock.DeactivateAndReset(); m.CleanupMockState() })
	initial, err := helpers.ParseHCLFile("tests/terraform/unit/resource.tf")
	require.NoError(t, err)
	updated, err := helpers.ParseHCLFile("tests/terraform/unit/resource_updated.tf")
	require.NoError(t, err)
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: initial,
					Check: check.That(resourceType + ".test").
						Key("internet_access.is_source_ip_session_affinity_enabled").
						HasValue("true"),
				},
				{
					ResourceName:      resourceType + ".test",
					ImportState:       true,
					ImportStateId:     "explicitForwardProxyConfig",
					ImportStateVerify: true,
				},
				{Config: initial, PlanOnly: true},
				{Config: updated},
				{Config: updated, PlanOnly: true},
			},
		},
	)
	require.Equal(t, 0, m.Deletes)
	require.Len(t, m.Patches, 2)
	require.Len(t, m.Patches[0], 1)
	require.NotContains(t, m.Patches[0]["internetAccess"], "isMtlsRequired")
	require.Equal(
		t,
		map[string]any{"internetAccess": map[string]any{"isSourceIpSessionAffinityEnabled": false}},
		m.Patches[1],
	)
	access := m.State["internetAccess"].(map[string]any)
	require.Equal(t, true, access["isMtlsRequired"])
	require.Equal(t, "preserve", access["futureSetting"])
	require.Equal(t, "none", access["sourceIpSessionAffinityOptions"])
}
