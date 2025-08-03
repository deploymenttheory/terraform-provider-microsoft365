package graphBetaAssignmentFilter_test

import (
	"os"
	"path/filepath"
	"regexp"

	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	assignmentFilterMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/assignment_filter/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *assignmentFilterMocks.AssignmentFilterMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	assignmentFilterMock := &assignmentFilterMocks.AssignmentFilterMock{}
	assignmentFilterMock.RegisterMocks()

	return mockClient, assignmentFilterMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *assignmentFilterMocks.AssignmentFilterMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	assignmentFilterMock := &assignmentFilterMocks.AssignmentFilterMock{}
	assignmentFilterMock.RegisterErrorMocks()

	return mockClient, assignmentFilterMock
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

// TestAssignmentFilterResource_Schema validates the resource schema
func TestAssignmentFilterResource_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "display_name", "Test Minimal Assignment Filter - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "rule", "(device.osVersion -startsWith \"10.0\")"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "assignment_filter_management_type", "devices"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.minimal", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.minimal", "last_modified_date_time"),
				),
			},
		},
	})
}

// TestAssignmentFilterResource_Minimal tests basic CRUD operations
func TestAssignmentFilterResource_Minimal(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_assignment_filter.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "display_name", "Test Minimal Assignment Filter - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "rule", "(device.osVersion -startsWith \"10.0\")"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "assignment_filter_management_type", "devices"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_assignment_filter.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_assignment_filter.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "display_name", "Test Maximal Assignment Filter - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "description", "Maximal assignment filter for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "assignment_filter_management_type", "devices"),
				),
			},
		},
	})
}

// TestAssignmentFilterResource_UpdateInPlace tests in-place updates
func TestAssignmentFilterResource_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_assignment_filter.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.minimal", "display_name", "Test Minimal Assignment Filter - Unique"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_assignment_filter.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "display_name", "Test Maximal Assignment Filter - Unique"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "description", "Maximal assignment filter for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "assignment_filter_management_type", "devices"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.maximal", "role_scope_tags.#", "2"),
				),
			},
		},
	})
}

// TestAssignmentFilterResource_PlatformValidation tests platform validation
func TestAssignmentFilterResource_PlatformValidation(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name = "Test Assignment Filter"
  platform     = "invalid_platform"
  rule         = "(device.osVersion -startsWith \"10.0\")"
}
`,
				ExpectError: regexp.MustCompile(`Attribute platform value must be one of:`),
			},
		},
	})
}

// TestAssignmentFilterResource_RequiredFields tests required field validation
func TestAssignmentFilterResource_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing display_name",
			config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  platform = "windows10AndLater"
  rule     = "(device.osVersion -startsWith \"10.0\")"
}
`,
			expectedError: `The argument "display_name" is required`,
		},
		{
			name: "missing platform",
			config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name = "Test Assignment Filter"
  rule         = "(device.osVersion -startsWith \"10.0\")"
}
`,
			expectedError: `The argument "platform" is required`,
		},
		{
			name: "missing rule",
			config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name = "Test Assignment Filter"
  platform     = "windows10AndLater"
}
`,
			expectedError: `The argument "rule" is required`,
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

// TestAssignmentFilterResource_ErrorHandling tests error scenarios
func TestAssignmentFilterResource_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name = "Test Assignment Filter"
  platform     = "windows10AndLater" 
  rule         = "(device.osVersion -startsWith \"10.0\")"
}
`,
				ExpectError: regexp.MustCompile(`Invalid assignment filter data|BadRequest`),
			},
		},
	})
}

// TestAssignmentFilterResource_RoleScopeTags tests role scope tags handling
func TestAssignmentFilterResource_RoleScopeTags(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name    = "Test Assignment Filter"
  platform        = "windows10AndLater"
  rule            = "(device.osVersion -startsWith \"10.0\")"
  role_scope_tags = ["0", "1", "2"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "role_scope_tags.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "role_scope_tags.*", "0"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "role_scope_tags.*", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "role_scope_tags.*", "2"),
				),
			},
		},
	})
}

// TestAssignmentFilterResource_ComplexRule tests complex assignment filter rules
func TestAssignmentFilterResource_ComplexRule(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, assignmentFilterMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer assignmentFilterMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_assignment_filter" "test" {
  display_name = "Test Complex Rule Assignment Filter"
  platform     = "windows10AndLater"
  rule         = "(device.osVersion -startsWith \"10.0\") -and (device.manufacturer -eq \"Microsoft Corporation\") -and (device.model -notIn [\"Virtual\"])"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_assignment_filter.test", "display_name", "Test Complex Rule Assignment Filter"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_assignment_filter.test", "rule"),
				),
			},
		},
	})
}
