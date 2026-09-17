package graphBetaNetworkProxyAutoConfiguration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	target "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_proxy_auto_configuration"
	testmocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_proxy_auto_configuration/mocks"
)

const resourceType = target.ResourceName

func TestUnitResourceNetworkProxyAutoConfiguration_01_Lifecycle(t *testing.T) {
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
					Check:  check.That(resourceType + ".test").Key("is_enabled").HasValue("false"),
				},
				{
					ResourceName:      resourceType + ".test",
					ImportState:       true,
					ImportStateId:     "11111111-2222-3333-4444-555555555555",
					ImportStateVerify: true,
				},
				{Config: initial, PlanOnly: true},
				{Config: updated},
				{Config: updated, PlanOnly: true},
			},
		},
	)
	require.Equal(t, 1, m.Deletes)
	require.Nil(t, m.State)
	require.Len(t, m.Patches, 1)
	require.Len(t, m.Patches[0], 1)
	require.Contains(t, m.Patches[0], "content")
}
