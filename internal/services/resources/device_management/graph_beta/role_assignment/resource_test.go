package graphBetaRoleDefinitionAssignment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleAssignmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_assignment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *roleAssignmentMocks.RoleAssignmentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleAssignmentMock := &roleAssignmentMocks.RoleAssignmentMock{}
	roleAssignmentMock.RegisterMocks()
	return mockClient, roleAssignmentMock
}

// Helper function to load test configs from unit directory
func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestRoleAssignmentResource_Schema validates the resource schema
func TestUnitResourceRoleAssignment_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-assignment-minimal"),
					check.That(resourceType+".minimal").Key("description").HasValue("Minimal role assignment for unit testing"),
					check.That(resourceType+".minimal").Key("role_definition_id").HasValue("0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					check.That(resourceType+".minimal").Key("members.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("scope_configuration.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("scope_configuration.0.type").HasValue("AllLicensedUsers"),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_Minimal tests basic CRUD operations
func TestUnitResourceRoleAssignment_02_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-assignment-minimal"),
					check.That(resourceType+".minimal").Key("description").HasValue("Minimal role assignment for unit testing"),
					check.That(resourceType+".minimal").Key("role_definition_id").HasValue("0bd113fe-6be5-400c-a28f-ae5553f9c0be"),
					check.That(resourceType+".minimal").Key("scope_configuration.0.type").HasValue("AllLicensedUsers"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".minimal"]
					if !ok {
						return "", fmt.Errorf("not found: %s.minimal", resourceType)
					}
					id := rs.Primary.ID
					roleDefId := rs.Primary.Attributes["role_definition_id"]
					compositeId := fmt.Sprintf("%s/%s", id, roleDefId)
					return compositeId, nil
				},
			},
		},
	})
}

// TestRoleAssignmentResource_Maximal tests maximal configuration
func TestUnitResourceRoleAssignment_03_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-role-assignment-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("Comprehensive role assignment for unit testing with all features"),
					check.That(resourceType+".maximal").Key("role_definition_id").HasValue("9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"),
					check.That(resourceType+".maximal").Key("members.#").HasValue("3"),
					check.That(resourceType+".maximal").Key("scope_configuration.0.type").HasValue("ResourceScopes"),
					check.That(resourceType+".maximal").Key("scope_configuration.0.resource_scopes.#").HasValue("3"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_UpdateInPlace tests in-place updates
func TestUnitResourceRoleAssignment_04_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-assignment-minimal"),
					check.That(resourceType+".minimal").Key("scope_configuration.0.type").HasValue("AllLicensedUsers"),
				),
			},
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-role-assignment-maximal"),
					check.That(resourceType+".maximal").Key("scope_configuration.0.type").HasValue("ResourceScopes"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_ScopeValidation tests scope configuration validation
func TestUnitResourceRoleAssignment_05_ScopeValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing_scope_configuration",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
}
`,
			expectedError: `Missing Scope Configuration`,
		},
		{
			name: "resource_scopes_without_scopes",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  
  scope_configuration {
    type = "ResourceScopes"
  }
}
`,
			expectedError: `Missing Resource Scopes`,
		},
		{
			name: "all_users_with_scopes",
			config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "Test Role Assignment"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  members = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  
  scope_configuration {
    type = "AllLicensedUsers"
    resource_scopes = ["ea8e2fb8-e909-44e6-bae7-56757cf6f347"]
  }
}
`,
			expectedError: `Invalid Resource Scopes`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      tc.config,
						ExpectError: regexp.MustCompile(tc.expectedError),
					},
				},
			})
		})
	}
}

// TestRoleAssignmentResource_AllDevicesScope tests AllDevices scope type
func TestUnitResourceRoleAssignment_06_AllDevicesScope(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "unit-test-role-assignment-all-devices"
  description        = "Role assignment for all devices scope"
  role_definition_id = "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"
  
  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
    "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
  ]
  
  scope_configuration {
    type = "AllDevices"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("unit-test-role-assignment-all-devices"),
					check.That(resourceType+".test").Key("scope_configuration.0.type").HasValue("AllDevices"),
					check.That(resourceType+".test").Key("members.#").HasValue("2"),
				),
			},
		},
	})
}

// TestRoleAssignmentResource_Members tests members handling
func TestUnitResourceRoleAssignment_07_Members(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleAssignmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleAssignmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_assignment" "test" {
  display_name       = "unit-test-role-assignment-members"
  description        = "Minimal role assignment for unit testing"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
  
  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
    "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2",
    "35d09841-af73-43e6-a59f-024fef1b6b95"
  ]
  
  scope_configuration {
    type = "AllLicensedUsers"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("members.#").HasValue("3"),
					check.That(resourceType+".test").Key("members.*").ContainsTypeSetElement("ea8e2fb8-e909-44e6-bae7-56757cf6f347"),
					check.That(resourceType+".test").Key("members.*").ContainsTypeSetElement("b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"),
					check.That(resourceType+".test").Key("members.*").ContainsTypeSetElement("35d09841-af73-43e6-a59f-024fef1b6b95"),
				),
			},
		},
	})
}
