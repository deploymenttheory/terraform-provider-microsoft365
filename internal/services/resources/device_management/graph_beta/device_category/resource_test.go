package graphBetaDeviceCategory_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	deviceCategoryMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_category/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *deviceCategoryMocks.DeviceCategoryMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	deviceCategoryMock := &deviceCategoryMocks.DeviceCategoryMock{}
	deviceCategoryMock.RegisterMocks()

	return mockClient, deviceCategoryMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *deviceCategoryMocks.DeviceCategoryMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	deviceCategoryMock := &deviceCategoryMocks.DeviceCategoryMock{}
	deviceCategoryMock.RegisterErrorMocks()

	return mockClient, deviceCategoryMock
}

// testCheckExists is a basic check to ensure the resource exists in the state
func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

// testConfigMinimal returns the minimal configuration for testing
func testConfigMinimal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_minimal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// testConfigMaximal returns the maximal configuration for testing
func testConfigMaximal() string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "unit", "resource_maximal.tf"))
	if err != nil {
		return ""
	}
	return string(content)
}

// TestDeviceCategoryResource_Schema validates the resource schema
func TestDeviceCategoryResource_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.minimal", "display_name", "Test Minimal Device Category - Unique"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_device_category.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.minimal", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.minimal", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

// TestDeviceCategoryResource_Minimal tests basic CRUD operations
func TestDeviceCategoryResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_category.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.minimal", "display_name", "Test Minimal Device Category - Unique"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_device_category.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_category.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "display_name", "Test Maximal Device Category - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "description", "Maximal device category for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestDeviceCategoryResource_UpdateInPlace tests in-place updates
func TestDeviceCategoryResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_category.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.minimal", "display_name", "Test Minimal Device Category - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_device_category.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "display_name", "Test Maximal Device Category - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "description", "Maximal device category for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.maximal", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

// TestDeviceCategoryResource_RequiredFields tests required field validation
func TestDeviceCategoryResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_device_category" "test" {
  # Missing display_name
}
`,
				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
			},
		},
	})
}

// TestDeviceCategoryResource_ErrorHandling tests error scenarios
func TestDeviceCategoryResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_device_category" "test" {
  display_name = "Test Device Category"
}
`,
				ExpectError: regexp.MustCompile(`Invalid device category data|BadRequest`),
			},
		},
	})
}

// TestDeviceCategoryResource_RoleScopeTags tests role scope tags handling
func TestDeviceCategoryResource_RoleScopeTags(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_device_category" "test" {
  display_name       = "Test Device Category"
  role_scope_tag_ids = ["0", "1", "2"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_device_category.test", "role_scope_tag_ids.*", "2"),
				),
			},
		},
	})
}

// TestDeviceCategoryResource_Description tests optional description field
func TestDeviceCategoryResource_Description(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, deviceCategoryMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceCategoryMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_device_category" "test" {
  display_name = "Test Device Category"
  description  = "Test description for device category"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "display_name", "Test Device Category"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_device_category.test", "description", "Test description for device category"),
				),
			},
		},
	})
}
