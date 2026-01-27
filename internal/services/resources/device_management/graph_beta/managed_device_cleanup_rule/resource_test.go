package graphBetaManagedDeviceCleanupRule_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	cleanupMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/managed_device_cleanup_rule/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *cleanupMocks.ManagedDeviceCleanupRuleMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	cleanupMock := &cleanupMocks.ManagedDeviceCleanupRuleMock{}
	cleanupMock.RegisterMocks()

	return mockClient, cleanupMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *cleanupMocks.ManagedDeviceCleanupRuleMock) {
	httpmock.Activate()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	cleanupMock := &cleanupMocks.ManagedDeviceCleanupRuleMock{}
	cleanupMock.RegisterErrorMocks()

	return mockClient, cleanupMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// Schema validation
func TestUnitResourceManagedDeviceCleanupRule_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, cleanupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer cleanupMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigFromFile("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "display_name", "Test Minimal Managed Device Cleanup Rule - Unique"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_managed_device_cleanup_rule.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
				),
			},
		},
	})
}

// Platform-specific creation tests (each value once)
func TestUnitResourceManagedDeviceCleanupRule_02_Platforms_Create(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, cleanupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer cleanupMock.CleanupMockState()

	tests := []struct {
		name   string
		path   string
		resRef string
	}{
		{"all", "tests/terraform/unit/platform_all.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.all"},
		{"androidAOSP", "tests/terraform/unit/platform_androidAOSP.tf", "microsoft365_graph_beta_device_management_managed_device_cleanup_rule.androidAOSP"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config: testConfigFromFile(tc.path),
					Check: resource.ComposeTestCheckFunc(
						testCheckExists(tc.resRef),
					),
				}},
			})
		})
	}
}

// Duplicate platform should fail (server returns 500)
func TestUnitResourceManagedDeviceCleanupRule_03_DuplicatePlatform_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, cleanupMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer cleanupMock.CleanupMockState()

	configFirst := testConfigFromFile("tests/terraform/unit/platform_all.tf")
	configDuplicate := testConfigFromFile("tests/terraform/unit/platform_all.tf")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: configFirst},
			{Config: configDuplicate, ExpectError: regexp.MustCompile("A cleanup rule already exists for the specified platform type")},
		},
	})
}

func testConfigFromFile(rel string) string {
	unitTestConfig, err := helpers.ParseHCLFile(rel)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	return unitTestConfig
}
