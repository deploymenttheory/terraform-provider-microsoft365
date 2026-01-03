package graphBetaRoleScopeTag_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleScopeTagMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_scope_tag/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *roleScopeTagMocks.RoleScopeTagMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleScopeTagMock := &roleScopeTagMocks.RoleScopeTagMock{}
	roleScopeTagMock.RegisterMocks()
	return mockClient, roleScopeTagMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *roleScopeTagMocks.RoleScopeTagMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleScopeTagMock := &roleScopeTagMocks.RoleScopeTagMock{}
	roleScopeTagMock.RegisterErrorMocks()
	return mockClient, roleScopeTagMock
}

func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestRoleScopeTagResource_Schema validates the resource schema
func TestRoleScopeTagResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-minimal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("description").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-minimal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("is_built_in").HasValue("false"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_Minimal tests basic CRUD operations
func TestRoleScopeTagResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-minimal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("description").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-minimal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("is_built_in").HasValue("false"),
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
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-maximal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("description").MatchesRegex(regexp.MustCompile(`^unit-test-role-scope-tag-maximal-[A-Za-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("is_built_in").HasValue("false"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_UpdateInPlace tests in-place updates
func TestRoleScopeTagResource_UpdateInPlace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").Exists(),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.minimal", "display_name", regexp.MustCompile(`^unit-test-role-scope-tag-minimal-[A-Za-z0-9]{8}$`)),
				),
			},
			{
				Config: testConfigHelper("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").Exists(),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "display_name", regexp.MustCompile(`^unit-test-role-scope-tag-maximal-[A-Za-z0-9]{8}$`)),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "description", regexp.MustCompile(`^unit-test-role-scope-tag-maximal-[A-Za-z0-9]{8}$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.maximal", "assignments.#", "2"),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_RequiredFields tests required field validation
func TestRoleScopeTagResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  description = "Test Role Scope Tag"
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestRoleScopeTagResource_DisplayNameUniqueness tests display name uniqueness validation
func TestRoleScopeTagResource_DisplayNameUniqueness(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create first role scope tag
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "first" {
  display_name = "Unique Test Role Scope Tag"
}
`,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".first").Key("id").Exists(),
					check.That(resourceType+".first").Key("display_name").HasValue("Unique Test Role Scope Tag"),
				),
			},
			// Try to create second role scope tag with same display name
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "first" {
  display_name = "Unique Test Role Scope Tag"
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "second" {
  display_name = "Unique Test Role Scope Tag"
}
`,
				ExpectError: regexp.MustCompile(`already exists|Display names must be unique`),
			},
		},
	})
}

// TestRoleScopeTagResource_ErrorHandling tests error scenarios
func TestRoleScopeTagResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag"
  description  = "Test role scope tag for error testing"
}
`,
				ExpectError: regexp.MustCompile(`failed to retrieve existing role scope tags for validation|Invalid role scope tag data|BadRequest|Internal server error`),
			},
		},
	})
}

// TestRoleScopeTagResource_Assignments tests assignments handling
func TestRoleScopeTagResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag with Assignments"
  description  = "Role scope tag with group assignments"
  
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    }
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "11111111-1111-1111-1111-111111111111",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_role_scope_tag.test", "assignments.*", map[string]string{
						"type":     "groupAssignmentTarget",
						"group_id": "22222222-2222-2222-2222-222222222222",
					}),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_EmptyDescription tests handling of empty description
func TestRoleScopeTagResource_EmptyDescription(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag Empty Description"
  description  = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Role Scope Tag Empty Description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", ""),
				),
			},
		},
	})
}

// TestRoleScopeTagResource_NoDescription tests handling when description is not provided
func TestRoleScopeTagResource_NoDescription(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleScopeTagMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleScopeTagMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Role Scope Tag No Description"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "display_name", "Test Role Scope Tag No Description"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "description", ""),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_role_scope_tag.test", "is_built_in", "false"),
				),
			},
		},
	})
}
