package graphBetaOperationApprovalPolicy_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	operationApprovalPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/operation_approval_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupUnitTestEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TF_ACC", "0")
	os.Setenv("MS365_TEST_MODE", "true")
}

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *operationApprovalPolicyMocks.OperationApprovalPolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register local mocks directly
	operationApprovalPolicyMock := &operationApprovalPolicyMocks.OperationApprovalPolicyMock{}
	operationApprovalPolicyMock.RegisterMocks()

	return mockClient, operationApprovalPolicyMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *operationApprovalPolicyMocks.OperationApprovalPolicyMock) {
	// Activate httpmock
	httpmock.Activate()

	// Create a new Mocks instance and register authentication mocks
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Register error mocks
	operationApprovalPolicyMock := &operationApprovalPolicyMocks.OperationApprovalPolicyMock{}
	operationApprovalPolicyMock.RegisterErrorMocks()

	return mockClient, operationApprovalPolicyMock
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

// TestUnitResourceOperationApprovalPolicy_01_Schema validates the resource schema
func TestUnitResourceOperationApprovalPolicy_01_Schema(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					// Check required attributes
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "display_name", "Test Minimal Operation Approval Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "policy_set.policy_type", "app"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "policy_set.policy_platform", "notApplicable"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "approver_group_ids.#", "1"),

					// Check computed attributes are set
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "last_modified_date_time"),
				),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_02_Minimal tests basic CRUD operations
func TestUnitResourceOperationApprovalPolicy_02_Minimal(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_operation_approval_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "display_name", "Test Minimal Operation Approval Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "policy_set.policy_type", "app"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "approver_group_ids.#", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_operation_approval_policy.minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_operation_approval_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "display_name", "Test Maximal Operation Approval Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "description", "Maximal operation approval policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "policy_set.policy_type", "script"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "approver_group_ids.#", "2"),
				),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_03_UpdateInPlace tests in-place updates
func TestUnitResourceOperationApprovalPolicy_03_UpdateInPlace(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_operation_approval_policy.minimal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.minimal", "display_name", "Test Minimal Operation Approval Policy"),
				),
			},
			{
				Config: testConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_operation_approval_policy.maximal"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "display_name", "Test Maximal Operation Approval Policy"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "description", "Maximal operation approval policy for testing with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "policy_platform", "windows10AndLater"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.maximal", "approver_group_ids.#", "2"),
				),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_04_PolicyTypeValidation tests policy_type enum validation
func TestUnitResourceOperationApprovalPolicy_04_PolicyTypeValidation(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Operation Approval Policy"

  policy_set = {
    policy_type     = "invalid_policy_type"
    policy_platform = "notApplicable"
  }

  approver_group_ids = ["11111111-1111-1111-1111-111111111111"]
}
`,
				ExpectError: regexp.MustCompile(`Attribute policy_set.policy_type value must be one of:`),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_05_RequiredFields tests required field validation
func TestUnitResourceOperationApprovalPolicy_05_RequiredFields(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	testCases := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "missing display_name",
			config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }
  approver_group_ids = ["11111111-1111-1111-1111-111111111111"]
}
`,
			expectedError: `The argument "display_name" is required`,
		},
		{
			name: "missing policy_set",
			config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name       = "Test Operation Approval Policy"
  approver_group_ids = ["11111111-1111-1111-1111-111111111111"]
}
`,
			expectedError: `The argument "policy_set" is required`,
		},
		{
			name: "missing approver_group_ids",
			config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Operation Approval Policy"
  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }
}
`,
			expectedError: `The argument "approver_group_ids" is required`,
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

// TestUnitResourceOperationApprovalPolicy_06_ErrorHandling tests error scenarios
func TestUnitResourceOperationApprovalPolicy_06_ErrorHandling(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Operation Approval Policy"

  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }

  approver_group_ids = ["11111111-1111-1111-1111-111111111111"]
}
`,
				ExpectError: regexp.MustCompile(`Invalid operation approval policy data|BadRequest`),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_07_ApproverGroupIds tests approver group IDs handling (typeset)
func TestUnitResourceOperationApprovalPolicy_07_ApproverGroupIds(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Operation Approval Policy"

  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }

  approver_group_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333",
  ]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_operation_approval_policy.test", "approver_group_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_operation_approval_policy.test", "approver_group_ids.*", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_operation_approval_policy.test", "approver_group_ids.*", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_operation_approval_policy.test", "approver_group_ids.*", "33333333-3333-3333-3333-333333333333"),
				),
			},
		},
	})
}

// TestUnitResourceOperationApprovalPolicy_08_InvalidApproverGroupId tests GUID format validation
func TestUnitResourceOperationApprovalPolicy_08_InvalidApproverGroupId(t *testing.T) {
	setupUnitTestEnvironment(t)
	_, operationApprovalPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer operationApprovalPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "test" {
  display_name = "Test Operation Approval Policy"

  policy_set = {
    policy_type     = "app"
    policy_platform = "notApplicable"
  }

  approver_group_ids = ["not-a-guid"]
}
`,
				ExpectError: regexp.MustCompile(`Must be a valid GUID format`),
			},
		},
	})
}
