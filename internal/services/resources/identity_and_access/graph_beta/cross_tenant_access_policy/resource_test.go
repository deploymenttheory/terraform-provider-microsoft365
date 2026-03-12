package graphBetaCrossTenantAccessPolicy_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	crossTenantMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *crossTenantMocks.CrossTenantAccessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	crossTenantMock := &crossTenantMocks.CrossTenantAccessPolicyMock{}
	crossTenantMock.RegisterMocks()
	return mockClient, crossTenantMock
}

func loadUnitConfig(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit config " + filename + ": " + err.Error())
	}
	return config
}

// TestUnitResourceCrossTenantAccessPolicy_01_WithNoB2B verifies creation of the singleton policy
// with no allowed_cloud_endpoints (no cross-cloud B2B collaboration).
func TestUnitResourceCrossTenantAccessPolicy_01_WithNoB2B(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_cross_tenant_access_policy_with_no_b2b.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_no_b2b").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".with_no_b2b").Key("display_name").HasValue("CrossTenantAccessPolicy"),
					check.That(resourceType+".with_no_b2b").Key("allowed_cloud_endpoints.#").HasValue("0"),
					check.That(resourceType+".with_no_b2b").Key("restore_defaults_on_destroy").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".with_no_b2b",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"restore_defaults_on_destroy",
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPolicy_02_WithAllowedCloudEndpoints verifies the policy
// can be configured with allowed_cloud_endpoints.
func TestUnitResourceCrossTenantAccessPolicy_02_WithAllowedCloudEndpoints(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_cross_tenant_access_policy_with_endpoints.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_endpoints").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".with_endpoints").Key("display_name").Exists(),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.#").HasValue("2"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("microsoftonline.us"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("partner.microsoftonline.cn"),
					check.That(resourceType+".with_endpoints").Key("restore_defaults_on_destroy").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".with_endpoints",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"restore_defaults_on_destroy",
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPolicy_03_RestoreDefaultsOnDestroy verifies the restore_defaults_on_destroy
// flag is correctly stored in state and accepted by the schema.
func TestUnitResourceCrossTenantAccessPolicy_03_RestoreDefaultsOnDestroy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_cross_tenant_access_policy_restore_defaults.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_endpoints").Key("id").HasValue("crossTenantAccessPolicy"),
					check.That(resourceType+".with_endpoints").Key("display_name").Exists(),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.#").HasValue("2"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("microsoftonline.us"),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.*").ContainsTypeSetElement("partner.microsoftonline.cn"),
					check.That(resourceType+".with_endpoints").Key("restore_defaults_on_destroy").HasValue("true"),
				),
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPolicy_04_UpdateEndpoints verifies that updating
// allowed_cloud_endpoints in a subsequent plan step is applied correctly via PATCH.
func TestUnitResourceCrossTenantAccessPolicy_04_UpdateEndpoints(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitConfig("resource_cross_tenant_access_policy_with_no_b2b.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_no_b2b").Key("allowed_cloud_endpoints.#").HasValue("0"),
				),
			},
			{
				Config: loadUnitConfig("resource_cross_tenant_access_policy_with_endpoints.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".with_endpoints").Key("display_name").Exists(),
					check.That(resourceType+".with_endpoints").Key("allowed_cloud_endpoints.#").HasValue("2"),
				),
			},
		},
	})
}
