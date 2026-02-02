package graphBetaTenantAppManagementPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTenantAppManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/tenant_app_management_policy"
	tenantAppPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/tenant_app_management_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaTenantAppManagementPolicy.ResourceName

	// testResource is the test resource implementation for tenant app management policy
	testResource = graphBetaTenantAppManagementPolicy.TenantAppManagementPolicyTestResource{}
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *tenantAppPolicyMocks.TenantAppManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tenantAppPolicyMock := &tenantAppPolicyMocks.TenantAppManagementPolicyMock{}
	tenantAppPolicyMock.RegisterMocks()
	return mockClient, tenantAppPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *tenantAppPolicyMocks.TenantAppManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	tenantAppPolicyMock := &tenantAppPolicyMocks.TenantAppManagementPolicyMock{}
	tenantAppPolicyMock.RegisterErrorMocks()
	return mockClient, tenantAppPolicyMock
}

// TestUnitResourceTenantAppManagementPolicy_01_Minimal tests creating a tenant app management policy with minimal configuration
func TestUnitResourceTenantAppManagementPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tenantAppPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tenantAppPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.state").HasValue("enabled"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.max_lifetime").HasValue("P90D"),
					check.That(resourceType+".minimal").Key("application_restrictions.password_credentials.0.restrict_for_apps_created_after_date_time").HasValue("2024-01-01T00:00:00Z"),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitResourceTenantAppManagementPolicy_02_Maximal tests creating a tenant app management policy with maximal configuration
func TestUnitResourceTenantAppManagementPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, tenantAppPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer tenantAppPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").HasValue("00000000-0000-0000-0000-000000000000"),
					check.That(resourceType+".maximal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".maximal").Key("display_name").HasValue("Custom Tenant App Management Policy"),
					check.That(resourceType+".maximal").Key("description").HasValue("Enforces comprehensive app management restrictions"),
					check.That(resourceType+".maximal").Key("restore_to_default_upon_delete").HasValue("false"),
					check.That(resourceType+".maximal").Key("application_restrictions.password_credentials.0.restriction_type").HasValue("passwordAddition"),
					check.That(resourceType+".maximal").Key("application_restrictions.password_credentials.1.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".maximal").Key("application_restrictions.password_credentials.2.restriction_type").HasValue("symmetricKeyLifetime"),
					check.That(resourceType+".maximal").Key("application_restrictions.key_credentials.0.restriction_type").HasValue("asymmetricKeyLifetime"),
					check.That(resourceType+".maximal").Key("application_restrictions.key_credentials.1.restriction_type").HasValue("trustedCertificateAuthority"),
					check.That(resourceType+".maximal").Key("service_principal_restrictions.password_credentials.0.restriction_type").HasValue("passwordLifetime"),
					check.That(resourceType+".maximal").Key("service_principal_restrictions.key_credentials.0.restriction_type").HasValue("asymmetricKeyLifetime"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitResourceTenantAppManagementPolicy_03_Error tests error handling
func TestUnitResourceTenantAppManagementPolicy_03_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, _ = setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigError(),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testConfigError() string {
	return loadUnitTestTerraform("resource_03_error.tf")
}
