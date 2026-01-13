package graphBetaGroupPolicyConfiguration_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaGroupPolicyConfiguration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_configuration"
	groupPolicyConfigurationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/group_policy_configuration/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = graphBetaGroupPolicyConfiguration.ResourceName

// Helper function to load test configs
func testConfigHelper(filename string) string {
	config, err := helpers.ParseHCLFile(filename)
	if err != nil {
		panic("failed to load config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *groupPolicyConfigurationMocks.GroupPolicyConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupPolicyConfigurationMock := &groupPolicyConfigurationMocks.GroupPolicyConfigurationMock{}
	groupPolicyConfigurationMock.RegisterMocks()
	return mockClient, groupPolicyConfigurationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *groupPolicyConfigurationMocks.GroupPolicyConfigurationMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	groupPolicyConfigurationMock := &groupPolicyConfigurationMocks.GroupPolicyConfigurationMock{}
	groupPolicyConfigurationMock.RegisterErrorMocks()
	return mockClient, groupPolicyConfigurationMock
}

// TestUnitGroupPolicyConfigurationResource_Minimal tests creating a minimal group policy configuration
func TestUnitGroupPolicyConfigurationResource_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-001-minimal"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("policy_configuration_ingestion_type").Exists(),
					check.That(resourceType+".minimal").Key("created_date_time").Exists(),
					check.That(resourceType+".minimal").Key("last_modified_date_time").Exists(),
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

// TestUnitGroupPolicyConfigurationResource_Maximal tests creating a maximal group policy configuration
func TestUnitGroupPolicyConfigurationResource_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("display_name").HasValue("unit-test-002-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("unit-test-002-maximal"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".maximal").Key("policy_configuration_ingestion_type").Exists(),
					check.That(resourceType+".maximal").Key("created_date_time").Exists(),
					check.That(resourceType+".maximal").Key("last_modified_date_time").Exists(),
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

// TestUnitGroupPolicyConfigurationResource_MinimalAssignment tests creating a configuration with minimal assignment
func TestUnitGroupPolicyConfigurationResource_MinimalAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_minimal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal_assignment").Key("display_name").HasValue("unit-test-003-minimal-assignment"),
					check.That(resourceType+".minimal_assignment").Key("assignments.#").HasValue("1"),
					check.That(resourceType+".minimal_assignment").Key("assignments.0.type").HasValue("allDevicesAssignmentTarget"),
				),
			},
			{
				ResourceName:      resourceType + ".minimal_assignment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitGroupPolicyConfigurationResource_MaximalAssignment tests creating a configuration with maximal assignments
func TestUnitGroupPolicyConfigurationResource_MaximalAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_maximal_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal_assignment").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal_assignment").Key("display_name").HasValue("unit-test-004-maximal-assignment"),
					check.That(resourceType+".maximal_assignment").Key("description").HasValue("unit-test-004-maximal-assignment"),
					check.That(resourceType+".maximal_assignment").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".maximal_assignment").Key("assignments.#").HasValue("4"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal_assignment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitGroupPolicyConfigurationResource_MinimalToMaximal tests transitioning from minimal to maximal configuration
func TestUnitGroupPolicyConfigurationResource_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("display_name").HasValue("unit-test-001-minimal"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				Config: testConfigHelper("tests/terraform/unit/resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("unit-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("description").HasValue("unit-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("3"),
				),
			},
			{
				ResourceName:      resourceType + ".transition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitGroupPolicyConfigurationResource_MaximalToMinimal tests transitioning from maximal to minimal configuration
func TestUnitGroupPolicyConfigurationResource_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigHelper("tests/terraform/unit/resource_minimal_to_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("unit-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("description").HasValue("unit-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".transition").Key("assignments.#").HasValue("3"),
				),
			},
			{
				Config: testConfigHelper("tests/terraform/unit/resource_maximal_to_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".transition").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".transition").Key("display_name").HasValue("unit-test-006-lifecycle-minimal"),
					// Description persists from previous state when not explicitly cleared
					check.That(resourceType+".transition").Key("description").HasValue("unit-test-005-lifecycle-maximal"),
					check.That(resourceType+".transition").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".transition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestUnitGroupPolicyConfigurationResource_ErrorHandling tests error scenarios
func TestUnitGroupPolicyConfigurationResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, groupPolicyConfigurationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer groupPolicyConfigurationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigHelper("tests/terraform/unit/resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(BadRequest|Invalid request body|Error)`),
			},
		},
	})
}
