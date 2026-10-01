package graphBetaExternalIdentitiesPolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/external_identities_policy"
	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/external_identities_policy/mocks"
)

func setupMockEnvironment(t *testing.T) *policymocks.ExternalIdentitiesPolicyMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	m := &policymocks.ExternalIdentitiesPolicyMock{}
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

// TestUnitResourceExternalIdentitiesPolicy_01_Lifecycle covers adoption, both Boolean values, import, and state-only deletion.
func TestUnitResourceExternalIdentitiesPolicy_01_Lifecycle(t *testing.T) {
	m := setupMockEnvironment(t)
	address := policy.ResourceName + ".test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform(t, "initial"), Check: resource.ComposeTestCheckFunc(
				check.That(address).Key("id").HasValue("externalIdentityPolicy"),
				check.That(address).Key("allow_external_identities_to_leave").HasValue("true"),
				check.That(address).Key("allow_deleted_identities_data_removal").HasValue("false"),
				check.That(address).Key("display_name").HasValue("External Identities Policy"),
			)},
			{
				Config: loadUnitTestTerraform(t, "updated"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).Key("allow_external_identities_to_leave").HasValue("false"),
					check.That(address).
						Key("allow_deleted_identities_data_removal").
						HasValue("true"),
				),
			},
			{
				ResourceName:      address,
				ImportState:       true,
				ImportStateId:     "externalIdentityPolicy",
				ImportStateVerify: true,
			},
			{Config: loadUnitTestTerraform(t, "updated"), PlanOnly: true},
			{Config: loadUnitTestTerraform(t, "updated"), Destroy: true},
		},
	})
	enabled, removal, patches := m.Snapshot()
	require.False(t, enabled)
	require.True(t, removal)
	require.Equal(t, 2, patches)
	require.Zero(t, httpmock.GetCallCountInfo()["DELETE "+policymocks.URL])
	require.Zero(t, httpmock.GetCallCountInfo()["POST "+policymocks.URL])
}
