package graphBetaIdentityAndAccessRoleDefinition_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleDefinitionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/role_definition/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *roleDefinitionMocks.RoleDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleDefinitionMock := &roleDefinitionMocks.RoleDefinitionMock{}
	roleDefinitionMock.RegisterMocks()
	return mockClient, roleDefinitionMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *roleDefinitionMocks.RoleDefinitionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleDefinitionMock := &roleDefinitionMocks.RoleDefinitionMock{}
	roleDefinitionMock.RegisterErrorMocks()
	return mockClient, roleDefinitionMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Minimal: required attributes and one permission block
func TestUnitResourceRoleDefinition_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-definition-minimal"),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("description").DoesNotExist(),
					check.That(resourceType+".minimal").Key("version").DoesNotExist(),
					check.That(resourceType+".minimal").Key("template_id").MatchesOtherKey(check.That(resourceType+".minimal").Key("id")),
					check.That(resourceType+".minimal").Key("is_built_in").HasValue("false"),
					check.That(resourceType+".minimal").Key("is_privileged").HasValue("false"),
					check.That(resourceType+".minimal").Key("resource_scopes.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("resource_scopes.*").ContainsTypeSetElement("/"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/applications/standard/read"),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Maximal: every optional attribute, a custom template ID and two permission blocks
func TestUnitResourceRoleDefinition_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".maximal").Key("template_id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-role-definition-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("Unit test custom role with every optional field set"),
					check.That(resourceType+".maximal").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".maximal").Key("version").HasValue("1"),
					check.That(resourceType+".maximal").Key("role_permissions.#").HasValue("2"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/applications/standard/read"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/groups/standard/read"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/devices/standard/read"),
					check.That(resourceType+".maximal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/bitlockerKeys/metadata/read"),
				),
			},
			{
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Minimal to maximal and back: in-place updates, including clearing description and version
func TestUnitResourceRoleDefinition_03_MinimalToMaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-definition-minimal"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-definition-minimal-updated"),
					check.That(resourceType+".minimal").Key("description").HasValue("Unit test custom role updated in place"),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".minimal").Key("version").HasValue("2"),
					check.That(resourceType+".minimal").Key("template_id").MatchesOtherKey(check.That(resourceType+".minimal").Key("id")),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("2"),
					check.That(resourceType+".minimal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/groups/standard/read"),
					check.That(resourceType+".minimal").Key("role_permissions.*.allowed_resource_actions.*").ContainsTypeSetElement("microsoft.directory/devices/standard/read"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-role-definition-minimal"),
					check.That(resourceType+".minimal").Key("description").DoesNotExist(),
					check.That(resourceType+".minimal").Key("version").DoesNotExist(),
					check.That(resourceType+".minimal").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".minimal").Key("role_permissions.#").HasValue("1"),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Unknown action: rejected against the microsoft.directory action catalogue before the write
func TestUnitResourceRoleDefinition_04_ValidationErrorUnknownAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_action.tf"),
				ExpectError: regexp.MustCompile(`resource actions not found in the microsoft\.directory namespace:\s+microsoft\.directory/applications/basic/read`),
			},
		},
	})
}

// Wrong namespace: rejected by schema validation at plan time
func TestUnitResourceRoleDefinition_05_ValidationErrorWrongNamespace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_namespace.tf"),
				ExpectError: regexp.MustCompile(`must be a microsoft\.directory resource action`),
			},
		},
	})
}

// Graph error: a 400 on create surfaces as a Terraform error
func TestUnitResourceRoleDefinition_06_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Invalid request|Bad Request|400`),
			},
		},
	})
}

// Built-in role import: rejected because only custom roles can be managed
func TestUnitResourceRoleDefinition_07_ImportBuiltInRoleRejected(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleDefinitionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleDefinitionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("resource_minimal.tf"),
				ResourceName:  resourceType + ".minimal",
				ImportState:   true,
				ImportStateId: "62e90394-69f5-4237-9190-012177145e10",
				ExpectError:   regexp.MustCompile(`Built-in role definitions are not supported`),
			},
		},
	})
}
