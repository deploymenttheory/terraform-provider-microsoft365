package graphBetaAppControlForBusinessPolicy_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// loadAcceptanceTestTerraform loads acceptance test terraform configuration files
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Test 01: Minimal Configuration - No Assignments
func TestAccResourceAppControlForBusinessPolicy_01_MinimalNoAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "01: Creating minimal app control policy - no assignments")
				},
				Config: loadAcceptanceTestTerraform("01_minimal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("name").HasValue("acc-test-app-control-policy-minimal"),
					check.That(resourceType+".minimal").Key("description").HasValue("Minimal app control policy for testing - no assignments"),
					check.That(resourceType+".minimal").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("assignments.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "01: Importing minimal policy")
				},
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 02: Maximal Configuration - No Assignments
func TestAccResourceAppControlForBusinessPolicy_02_MaximalNoAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "02: Creating maximal app control policy - no assignments")
				},
				Config: loadAcceptanceTestTerraform("02_maximal-no-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("name").HasValue("acc-test-app-control-policy-maximal"),
					check.That(resourceType+".maximal").Key("description").HasValue("Maximal app control policy for testing with enhanced description - no assignments"),
					check.That(resourceType+".maximal").Key("policy_xml").IsSet(),
					check.That(resourceType+".maximal").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".maximal").Key("assignments.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "02: Importing maximal policy")
				},
				ResourceName:            resourceType + ".maximal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 03: Minimal to Maximal Configuration Update
func TestAccResourceAppControlForBusinessPolicy_03_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "03: Step 1 - Creating minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("03_minimal-to-maximal-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".step_test").ExistsInGraph(testResource),
					check.That(resourceType+".step_test").Key("name").HasValue("acc-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - starts minimal"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "03: Step 2 - Updating to maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("03_minimal-to-maximal-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".step_test").ExistsInGraph(testResource),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - updated to maximal with enhanced description"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "03: Importing step test policy")
				},
				ResourceName:            resourceType + ".step_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 05: Minimal Configuration - Minimal Assignments (1 group)
func TestAccResourceAppControlForBusinessPolicy_05_MinimalMinimalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "05: Creating minimal policy with minimal assignments (1 group)")
				},
				Config: loadAcceptanceTestTerraform("05_minimal-minimal-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal_minimal_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".minimal_minimal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal_minimal_assignments").Key("name").HasValue("acc-test-app-control-policy-minimal-minimal-assignments"),
					check.That(resourceType+".minimal_minimal_assignments").Key("description").HasValue("Minimal app control policy with minimal assignments (1 group)"),
					check.That(resourceType+".minimal_minimal_assignments").Key("policy_xml").IsSet(),
					check.That(resourceType+".minimal_minimal_assignments").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".minimal_minimal_assignments").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_minimal_assignments", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "05: Importing minimal minimal assignments policy")
				},
				ResourceName:            resourceType + ".minimal_minimal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 06: Minimal Configuration - Maximal Assignments (allLicensedUsers, group, allDevices)
func TestAccResourceAppControlForBusinessPolicy_06_MinimalMaximalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "06: Creating minimal policy with maximal assignments (3 targets)")
				},
				Config: loadAcceptanceTestTerraform("06_minimal-maximal-assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal_maximal_assignments").ExistsInGraph(testResource),
					check.That(resourceType+".minimal_maximal_assignments").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal_maximal_assignments").Key("name").HasValue("acc-test-app-control-policy-minimal-maximal-assignments"),
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
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".minimal_maximal_assignments", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "06: Importing minimal maximal assignments policy")
				},
				ResourceName:            resourceType + ".minimal_maximal_assignments",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 04: Maximal to Minimal Configuration Update
func TestAccResourceAppControlForBusinessPolicy_04_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "04: Step 1 - Creating maximal configuration")
				},
				Config: loadAcceptanceTestTerraform("04_maximal-to-minimal-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".step_test").ExistsInGraph(testResource),
					check.That(resourceType+".step_test").Key("name").HasValue("acc-test-app-control-policy-step-test"),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - starts maximal with enhanced description"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "04: Step 2 - Downgrading to minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("04_maximal-to-minimal-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".step_test").ExistsInGraph(testResource),
					check.That(resourceType+".step_test").Key("description").HasValue("Step test policy - downgraded to minimal"),
					check.That(resourceType+".step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".step_test").Key("role_scope_tag_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "04: Importing step test policy")
				},
				ResourceName:            resourceType + ".step_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 07: Minimal Assignments to Maximal Assignments Update
func TestAccResourceAppControlForBusinessPolicy_07_MinimalAssignmentsToMaximalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "07: Step 1 - Creating policy with minimal assignments (1 group)")
				},
				Config: loadAcceptanceTestTerraform("07_minimal-assignments-to-maximal-assignments-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_step_test").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("acc-test-app-control-policy-assignment-step-test"),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - starts with minimal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "07: Step 2 - Expanding to maximal assignments (3 targets)")
				},
				Config: loadAcceptanceTestTerraform("07_minimal-assignments-to-maximal-assignments-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_step_test").ExistsInGraph(testResource),
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
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "07: Importing assignment step test policy")
				},
				ResourceName:            resourceType + ".assignment_step_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 08: Maximal Assignments to Minimal Assignments Update
func TestAccResourceAppControlForBusinessPolicy_08_MaximalAssignmentsToMinimalAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "08: Step 1 - Creating policy with maximal assignments (3 targets)")
				},
				Config: loadAcceptanceTestTerraform("08_maximal-assignments-to-minimal-assignments-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_step_test").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_step_test").Key("name").HasValue("acc-test-app-control-policy-assignment-step-test"),
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
						"filter_type": "none",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "allDevicesAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "08: Step 2 - Reducing to minimal assignments (1 group)")
				},
				Config: loadAcceptanceTestTerraform("08_maximal-assignments-to-minimal-assignments-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".assignment_step_test").ExistsInGraph(testResource),
					check.That(resourceType+".assignment_step_test").Key("description").HasValue("Assignment step test policy - reduced to minimal assignments"),
					check.That(resourceType+".assignment_step_test").Key("policy_xml").IsSet(),
					check.That(resourceType+".assignment_step_test").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".assignment_step_test").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".assignment_step_test", "assignments.*", map[string]string{
						"type":        "groupAssignmentTarget",
						"filter_type": "none",
					}),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "08: Importing assignment step test policy")
				},
				ResourceName:            resourceType + ".assignment_step_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test 09: Error Handling - Invalid Policy XML
func TestAccResourceAppControlForBusinessPolicy_09_ErrorInvalidPolicyXML(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "09: Testing error handling with invalid XML")
				},
				Config:      loadAcceptanceTestTerraform("09_error-invalid-policy-xml.tf"),
				ExpectError: regexp.MustCompile("(?i)(invalid|error|failed|xml)"),
			},
		},
	})
}

// Test 10: Removed Policy Configuration
func TestAccResourceAppControlForBusinessPolicy_10_RemovedPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "10: Testing removed block for controlled destruction")
				},
				Config: loadAcceptanceTestTerraform("10_removed-policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("app control policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
				),
			},
		},
	})
}
