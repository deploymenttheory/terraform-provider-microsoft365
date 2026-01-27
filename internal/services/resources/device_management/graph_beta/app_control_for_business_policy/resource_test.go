package graphBetaAppControlForBusinessPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAppControlForBusinessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_policy"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/app_control_for_business_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaAppControlForBusinessPolicy.ResourceName

	// testResource is the test resource implementation for app control policies
	testResource = graphBetaAppControlForBusinessPolicy.AppControlForBusinessPolicyTestResource{}
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.AppControlForBusinessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.AppControlForBusinessPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.AppControlForBusinessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.AppControlForBusinessPolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

// loadUnitTestTerraform loads unit test terraform configuration files
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Minimal Configuration Schema Validation
func TestUnitResourceAppControlForBusinessPolicy_01_MinimalSchema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_minimal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("name").HasValue("unit-test-app-control-policy-minimal"),
					check.That(resourceType+".minimal").Key("description").HasValue("Minimal app control policy for testing - no assignments"),
					check.That(resourceType+".minimal").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
					check.That(resourceType+".minimal").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 02: XML Policy Validation
func TestUnitResourceAppControlForBusinessPolicy_02_XMLPolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_minimal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").IsSet(),
					check.That(resourceType+".minimal").Key("name").HasValue("unit-test-app-control-policy-minimal"),
					check.That(resourceType+".minimal").Key("policy_xml").IsSet(),
				),
			},
			{
				Config: loadUnitTestTerraform("02_maximal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("id").IsSet(),
					check.That(resourceType+".maximal").Key("name").HasValue("unit-test-app-control-policy-maximal"),
					check.That(resourceType+".maximal").Key("policy_xml").IsSet(),
				),
			},
		},
	})
}

// Test 03: Assignments Validation
func TestUnitResourceAppControlForBusinessPolicy_03_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_minimal-minimal-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_minimal_assignments").Key("id").IsSet(),
					check.That(resourceType+".minimal_minimal_assignments").Key("name").HasValue("unit-test-app-control-policy-minimal-minimal-assignments"),
					check.That(resourceType+".minimal_minimal_assignments").Key("description").HasValue("Minimal app control policy with minimal assignments (1 group)"),
					check.That(resourceType+".minimal_minimal_assignments").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal_minimal_assignments").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal_minimal_assignments").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_minimal_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("06_minimal-maximal-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_maximal_assignments").Key("id").IsSet(),
					check.That(resourceType+".minimal_maximal_assignments").Key("name").HasValue("unit-test-app-control-policy-minimal-maximal-assignments"),
					check.That(resourceType+".minimal_maximal_assignments").Key("description").HasValue("Minimal app control policy with maximal assignments (allLicensedUsers, group, allDevices)"),
					check.That(resourceType+".minimal_maximal_assignments").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal_maximal_assignments").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal_maximal_assignments").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "allLicensedUsersAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
		},
	})
}

// Test 04: Import Validation
func TestUnitResourceAppControlForBusinessPolicy_04_Import(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_minimal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").IsSet(),
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

// Test 05: Error Handling
func TestUnitResourceAppControlForBusinessPolicy_05_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("09_error-invalid-policy-xml.tf"),
				ExpectError: regexp.MustCompile("(?i)(invalid|error|failed|xml)"),
			},
		},
	})
}

// Test 06: Minimal Configuration with Maximal Assignments
func TestUnitResourceAppControlForBusinessPolicy_06_MinimalMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_minimal-maximal-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal_maximal_assignments").Key("id").IsSet(),
					check.That(resourceType+".minimal_maximal_assignments").Key("name").HasValue("unit-test-app-control-policy-minimal-maximal-assignments"),
					check.That(resourceType+".minimal_maximal_assignments").Key("description").HasValue("Minimal app control policy with maximal assignments (allLicensedUsers, group, allDevices)"),
					check.That(resourceType+".minimal_maximal_assignments").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal_maximal_assignments").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal_maximal_assignments").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "allLicensedUsersAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
		},
	})
}

// Test 07: Step Test - Minimal Assignments to Maximal Assignments
func TestUnitResourceAppControlForBusinessPolicy_07_MinimalAssignmentsToMaximalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_minimal-assignments-to-maximal-assignments-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_step_test").Key("id").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("unit-test-app-control-policy-assignment-step-test"),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - starts with minimal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("07_minimal-assignments-to-maximal-assignments-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_step_test").Key("id").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("unit-test-app-control-policy-assignment-step-test"),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - expanded to maximal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allLicensedUsersAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
		},
	})
}

// Test 08: Step Test - Maximal Assignments to Minimal Assignments
func TestUnitResourceAppControlForBusinessPolicy_08_MaximalAssignmentsToMinimalAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_maximal-assignments-to-minimal-assignments-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_step_test").Key("id").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("unit-test-app-control-policy-assignment-step-test"),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - starts with maximal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allLicensedUsersAssignmentTarget",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("08_maximal-assignments-to-minimal-assignments-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".assignment_step_test").Key("id").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("unit-test-app-control-policy-assignment-step-test"),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - reduced to minimal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"group_id":    "33333333-3333-3333-3333-333333333333",
						"filter_type": "none",
					}),
				),
			},
		},
	})
}

// Test 09: Step Test - Minimal to Maximal Configuration
func TestUnitResourceAppControlForBusinessPolicy_09_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_minimal-to-maximal-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".step_test").Key("id").IsSet(),
					check.That(resourceType+".step_test").Key("name").HasValue("unit-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - starts minimal"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".step_test").Key("assignments.#").HasValue("0"),
				),
			},
			{
				Config: loadUnitTestTerraform("03_minimal-to-maximal-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".step_test").Key("id").IsSet(),
					check.That(resourceType+".step_test").Key("name").HasValue("unit-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - updated to maximal with enhanced description"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".step_test").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}

// Test 10: Step Test - Maximal to Minimal Configuration
func TestUnitResourceAppControlForBusinessPolicy_10_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_maximal-to-minimal-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".step_test").Key("id").IsSet(),
					check.That(resourceType+".step_test").Key("name").HasValue("unit-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - starts maximal with enhanced description"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".step_test").Key("assignments.#").HasValue("0"),
				),
			},
			{
				Config: loadUnitTestTerraform("04_maximal-to-minimal-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".step_test").Key("id").IsSet(),
					check.That(resourceType+".step_test").Key("name").HasValue("unit-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - downgraded to minimal"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".step_test").Key("assignments.#").HasValue("0"),
				),
			},
		},
	})
}
