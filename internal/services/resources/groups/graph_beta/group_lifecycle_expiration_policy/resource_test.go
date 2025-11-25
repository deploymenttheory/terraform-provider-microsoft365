package graphBetaGroupLifecycleExpirationPolicy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group_lifecycle_expiration_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.GroupLifecyclePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.GroupLifecyclePolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.GroupLifecyclePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.GroupLifecyclePolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func testConfigAll() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_all.tf")
	return config
}

func testConfigSelected() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_selected.tf")
	return config
}

func testConfigNone() string {
	config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_none.tf")
	return config
}

// TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_All tests managed_group_types = "All"
func TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_All(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".all").Key("id").Exists(),
					check.That(resourceType+".all").Key("group_lifetime_in_days").HasValue("180"),
					check.That(resourceType+".all").Key("managed_group_types").HasValue("All"),
					check.That(resourceType+".all").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com"),
				),
			},
			{
				ResourceName:            resourceType + ".all",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_Selected tests managed_group_types = "Selected"
func TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_Selected(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigSelected(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".selected").Key("id").Exists(),
					check.That(resourceType+".selected").Key("group_lifetime_in_days").HasValue("365"),
					check.That(resourceType+".selected").Key("managed_group_types").HasValue("Selected"),
					check.That(resourceType+".selected").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com;notifications@deploymenttheory.com"),
				),
			},
			{
				ResourceName:            resourceType + ".selected",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_None tests managed_group_types = "None"
func TestUnitGroupLifecycleExpirationPolicyResource_ManagedGroupTypes_None(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNone(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".none").Key("id").Exists(),
					check.That(resourceType+".none").Key("group_lifetime_in_days").HasValue("365"),
					check.That(resourceType+".none").Key("managed_group_types").HasValue("None"),
					check.That(resourceType+".none").Key("alternate_notification_emails").HasValue("admin@deploymenttheory.com"),
				),
			},
			{
				ResourceName:            resourceType + ".none",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite_existing_policy"},
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_Delete tests resource deletion
func TestUnitGroupLifecycleExpirationPolicyResource_Delete(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAll(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".all").Key("id").Exists(),
				),
			},
			{
				Config: `# Empty config for deletion test`,
				Check: func(s *terraform.State) error {
					_, exists := s.RootModule().Resources[resourceType+".all"]
					if exists {
						return fmt.Errorf("resource still exists after deletion")
					}
					return nil
				},
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_RequiredFields tests required field validation
func TestUnitGroupLifecycleExpirationPolicyResource_RequiredFields(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "test" {
  managed_group_types           = "All"
  alternate_notification_emails = "admin@deploymenttheory.com"
}
`,
				ExpectError: regexp.MustCompile(`The argument "group_lifetime_in_days" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "test" {
  group_lifetime_in_days        = 180
  alternate_notification_emails = "admin@deploymenttheory.com"
}
`,
				ExpectError: regexp.MustCompile(`The argument "managed_group_types" is required`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "test" {
  group_lifetime_in_days = 180
  managed_group_types    = "All"
}
`,
				ExpectError: regexp.MustCompile(`The argument "alternate_notification_emails" is required`),
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_InvalidValues tests invalid value validation
func TestUnitGroupLifecycleExpirationPolicyResource_InvalidValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "test" {
  group_lifetime_in_days        = 0
  managed_group_types           = "All"
  alternate_notification_emails = "admin@deploymenttheory.com"
}
`,
				ExpectError: regexp.MustCompile(`Attribute group_lifetime_in_days value must be at least 30`),
			},
			{
				Config: `
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "test" {
  group_lifetime_in_days        = 180
  managed_group_types           = "Invalid"
  alternate_notification_emails = "admin@deploymenttheory.com"
}
`,
				ExpectError: regexp.MustCompile(`Attribute managed_group_types value must be one of`),
			},
		},
	})
}

// TestUnitGroupLifecycleExpirationPolicyResource_ErrorHandling tests API error handling
func TestUnitGroupLifecycleExpirationPolicyResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigAll(),
				ExpectError: regexp.MustCompile(`Bad Request|400|ApiError`),
			},
		},
	})
}
