package graphBetaRoleDefinition_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleDefinitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_role_definition/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *roleDefinitionMocks.RoleDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleDefinitionMock := &roleDefinitionMocks.RoleDefinitionMock{}
	roleDefinitionMock.RegisterMocks()
	return mockClient, roleDefinitionMock
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

const resourceType = "microsoft365_graph_beta_windows_365_cloud_pc_role_definition"

// TestUnitResourceCloudPcRoleDefinition_01_Schema validates the resource schema
func TestUnitResourceCloudPcRoleDefinition_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-minimal"),
					check.That(resourceType+".minimal").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".minimal").Key("is_built_in").Exists(),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("description").HasValue(""),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcRoleDefinition_02_Minimal tests basic CRUD operations
func TestUnitResourceCloudPcRoleDefinition_02_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-minimal"),
					check.That(resourceType+".minimal").Key("is_built_in_role_definition").Exists(),
					check.That(resourceType+".minimal").Key("is_built_in").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal_custom").Key("id").Exists(),
					check.That(resourceType+".maximal_custom").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-maximal"),
					check.That(resourceType+".maximal_custom").Key("description").HasValue("Comprehensive custom role definition for testing with all features"),
					check.That(resourceType+".maximal_custom").Key("role_permissions.#").HasValue("1"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcRoleDefinition_03_UpdateInPlace tests in-place updates
func TestUnitResourceCloudPcRoleDefinition_03_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-minimal"),
				),
			},
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal_custom").Key("id").Exists(),
					check.That(resourceType+".maximal_custom").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-maximal"),
					check.That(resourceType+".maximal_custom").Key("description").HasValue("Comprehensive custom role definition for testing with all features"),
					check.That(resourceType+".maximal_custom").Key("role_permissions.#").HasValue("1"),
				),
			},
		},
	})
}

// TestUnitResourceCloudPcRoleDefinition_04_RequiredFields tests required field validation
func TestUnitResourceCloudPcRoleDefinition_04_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	testCases := []struct {
		name          string
		configFile    string
		expectedError string
	}{
		{
			name:          "invalid_prefix_validation",
			configFile:    "resource_invalid_prefix.tf",
			expectedError: `must start with 'Microsoft.CloudPC/'`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testConfigHelper(tc.configFile),
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

// TestUnitResourceCloudPcRoleDefinition_05_ErrorHandling tests error scenarios
func TestUnitResourceCloudPcRoleDefinition_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("resource_error_handling.tf"),
				ExpectError: regexp.MustCompile(`invalid Cloud PC resource operation`),
			},
		},
	})
}

// TestUnitResourceCloudPcRoleDefinition_06_RolePermissions tests role permissions handling
func TestUnitResourceCloudPcRoleDefinition_06_RolePermissions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_role_permissions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("unit-test-cloud-pc-role-definition-role-permissions"),
					check.That(resourceType+".test").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".test").Key("role_permissions.0.allowed_resource_actions.#").HasValue("3"),
				),
			},
		},
	})
}
