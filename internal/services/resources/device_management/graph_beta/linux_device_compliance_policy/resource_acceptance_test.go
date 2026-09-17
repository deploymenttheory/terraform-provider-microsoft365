package graphBetaLinuxDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaLinuxDeviceCompliancePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_device_compliance_policy"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaLinuxDeviceCompliancePolicy.ResourceName

var testResource = graphBetaLinuxDeviceCompliancePolicy.LinuxDeviceCompliancePolicyTestResource{}

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceLinuxDeviceCompliancePolicy_01_Scenario_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaLinuxDeviceCompliancePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("name").Exists(),
					check.That(resourceType+".test_001").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_001").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_001").Key("created_date_time").Exists(),
					check.That(resourceType+".test_001").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_001").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_001", "assignments.#"),
				),
			},
			{
				ResourceName:            resourceType + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_02_Scenario_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaLinuxDeviceCompliancePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("name").Exists(),
					check.That(resourceType+".test_002").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_002").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_002").Key("created_date_time").Exists(),
					check.That(resourceType+".test_002").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_002").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(resourceType+".test_002").Key("settings_count").HasValue("8"),
					check.That(resourceType+".test_002").Key("device_encryption_required").HasValue("true"),
					check.That(resourceType+".test_002").Key("custom_compliance_required").HasValue("false"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(resourceType+".test_002").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(resourceType+".test_002").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(resourceType+".test_002").Key("password_policy_minimum_length").HasValue("12"),
					check.That(resourceType+".test_002").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(resourceType+".test_002").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(resourceType+".test_002").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(resourceType+".test_002").Key("role_scope_tag_ids.#").HasValue("1"),
					resource.TestCheckNoResourceAttr(resourceType+".test_002", "assignments.#"),
				),
			},
			{
				ResourceName:            resourceType + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaLinuxDeviceCompliancePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("name").Exists(),
					check.That(resourceType+".test_003").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_003").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_003").Key("created_date_time").Exists(),
					check.That(resourceType+".test_003").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_003").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_003").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "assignments.#"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_003").Key("name").Exists(),
					check.That(resourceType+".test_003").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_003").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_003").Key("created_date_time").Exists(),
					check.That(resourceType+".test_003").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_003").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(resourceType+".test_003").Key("settings_count").HasValue("8"),
					check.That(resourceType+".test_003").Key("device_encryption_required").HasValue("true"),
					check.That(resourceType+".test_003").Key("custom_compliance_required").HasValue("false"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(resourceType+".test_003").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(resourceType+".test_003").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(resourceType+".test_003").Key("password_policy_minimum_length").HasValue("12"),
					check.That(resourceType+".test_003").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(resourceType+".test_003").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(resourceType+".test_003").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(resourceType+".test_003").Key("role_scope_tag_ids.#").HasValue("1"),
					resource.TestCheckNoResourceAttr(resourceType+".test_003", "assignments.#"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			graphBetaLinuxDeviceCompliancePolicy.ResourceName,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("name").Exists(),
					check.That(resourceType+".test_004").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_004").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_004").Key("created_date_time").Exists(),
					check.That(resourceType+".test_004").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_004").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(resourceType+".test_004").Key("settings_count").HasValue("8"),
					check.That(resourceType+".test_004").Key("device_encryption_required").HasValue("true"),
					check.That(resourceType+".test_004").Key("custom_compliance_required").HasValue("false"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(resourceType+".test_004").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(resourceType+".test_004").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(resourceType+".test_004").Key("password_policy_minimum_length").HasValue("12"),
					check.That(resourceType+".test_004").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(resourceType+".test_004").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(resourceType+".test_004").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(resourceType+".test_004").Key("role_scope_tag_ids.#").HasValue("1"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "assignments.#"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_004").Key("name").Exists(),
					check.That(resourceType+".test_004").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_004").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_004").Key("created_date_time").Exists(),
					check.That(resourceType+".test_004").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_004").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_004").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_004", "assignments.#"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_004",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_05_AssignmentsMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaLinuxDeviceCompliancePolicy.ResourceName,
				TestResource: graphBetaLinuxDeviceCompliancePolicy.LinuxDeviceCompliancePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_005").Key("name").Exists(),
					check.That(resourceType+".test_005").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_005").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_005").Key("created_date_time").Exists(),
					check.That(resourceType+".test_005").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_005").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_005").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_005", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_005").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_005", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_005_1", "id"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_005",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_06_AssignmentsMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaLinuxDeviceCompliancePolicy.ResourceName,
				TestResource: graphBetaLinuxDeviceCompliancePolicy.LinuxDeviceCompliancePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_006").Key("name").Exists(),
					check.That(resourceType+".test_006").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_006").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_006").Key("created_date_time").Exists(),
					check.That(resourceType+".test_006").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_006").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_006").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_006", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_006").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_006", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_006_1", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_006", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_006_2", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_006", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_006_3", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".test_006", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_006",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaLinuxDeviceCompliancePolicy.ResourceName,
				TestResource: graphBetaLinuxDeviceCompliancePolicy.LinuxDeviceCompliancePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("name").Exists(),
					check.That(resourceType+".test_007").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_007").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_007").Key("created_date_time").Exists(),
					check.That(resourceType+".test_007").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_007").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_007", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_007_1", "id"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_007").Key("name").Exists(),
					check.That(resourceType+".test_007").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_007").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_007").Key("created_date_time").Exists(),
					check.That(resourceType+".test_007").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_007").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_007", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_007").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_007", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_007_1", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_007", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_007_2", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_007", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_007_3", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".test_007", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_007",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceLinuxDeviceCompliancePolicy_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaLinuxDeviceCompliancePolicy.ResourceName,
				TestResource: graphBetaLinuxDeviceCompliancePolicy.LinuxDeviceCompliancePolicyTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("name").Exists(),
					check.That(resourceType+".test_008").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_008").Key("created_date_time").Exists(),
					check.That(resourceType+".test_008").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("3"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_008", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_008_1", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_008", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_008_2", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_008", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_008_3", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".test_008", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("name").Exists(),
					check.That(resourceType+".test_008").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_008").Key("created_date_time").Exists(),
					check.That(resourceType+".test_008").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_uppercase"),
					check.That(resourceType+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test_008", "assignments.*.group_id", graphBetaGroup.ResourceName+".acc_test_group_008_1", "id"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_008").Key("name").Exists(),
					check.That(resourceType+".test_008").Key("platforms").HasValue("linux"),
					check.That(resourceType+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(resourceType+".test_008").Key("created_date_time").Exists(),
					check.That(resourceType+".test_008").Key("last_modified_date_time").Exists(),
					check.That(resourceType+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(resourceType+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(resourceType+".test_008", "assignments.#"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("linux device compliance policy", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test_008",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
