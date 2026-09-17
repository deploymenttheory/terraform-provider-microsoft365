package graphBetaSettingsCatalogConfigurationPolicy_test

import (
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaRoleScopeTag "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/role_scope_tag"
	graphBetaSettingsCatalogConfigurationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceType = graphBetaSettingsCatalogConfigurationPolicy.ResourceName

var testResource = graphBetaSettingsCatalogConfigurationPolicy.SettingsCatalogTestResource{}

func TestAccResourceSettingsCatalogConfigurationPolicy_01_Lifecycle(t *testing.T) {
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
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			// Create with minimal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "name", "macos mdm filevault2 settings"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings", "role_scope_tag_ids.*", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.macos_mdm_filevault2_settings",
				ImportState:       true,
				ImportStateVerify: true,
				// Ignore secret password field as Microsoft Graph returns different UUIDs for security
				ImportStateVerifyIgnore: []string{
					"configuration_policy.settings.0.setting_instance.group_setting_collection_value.0.children.6.simple_setting_value.value",
				},
			},
		},
	})
}

func TestAccResourceSettingsCatalogConfigurationPolicy_02_Maximal(t *testing.T) {
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
				ResourceType: graphBetaSettingsCatalogConfigurationPolicy.ResourceName,
				TestResource: graphBetaSettingsCatalogConfigurationPolicy.SettingsCatalogTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaRoleScopeTag.ResourceName,
				TestResource: graphBetaRoleScopeTag.RoleScopeTagTestResource{},
			},
		),
		Steps: []resource.TestStep{
			// Create with maximal configuration
			{
				Config: testAccSettingsCatalogConfigurationPolicyConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "id"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "name", "Test Acceptance Settings Catalog Policy - Updated"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "description", "Updated description for acceptance testing"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "platforms", "macOS"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.test", "role_scope_tag_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceSettingsCatalogConfigurationPolicy_03_Assignments(t *testing.T) {
	t.Log("=== ASSIGNMENTS TEST START ===")
	t.Log("Starting assignments acceptance test with comprehensive logging")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			t.Log("=== PRE-CHECK START ===")
			mocks.TestAccPreCheck(t)
			t.Log("=== PRE-CHECK COMPLETE ===")
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaSettingsCatalogConfigurationPolicy.ResourceName,
				TestResource: graphBetaSettingsCatalogConfigurationPolicy.SettingsCatalogTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaRoleScopeTag.ResourceName,
				TestResource: graphBetaRoleScopeTag.RoleScopeTagTestResource{},
			},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaGroup.ResourceName,
				TestResource: graphBetaGroup.GroupTestResource{},
			},
		),
		Steps: []resource.TestStep{
			// Create with all assignment types
			{
				PreConfig: func() {
					t.Log("=== STEP PRE-CONFIG ===")
					t.Log("About to apply configuration with assignments, groups, and role scope tags")
					t.Log("Expected dependencies: 3 groups, 2 role scope tags")
				},
				Config: testAccSettingsCatalogConfigurationPolicyConfig_assignments(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						t.Log("=== STEP CHECK START ===")
						t.Log("Verifying resource creation and attributes")
						return nil
					},
					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "id"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Resource ID is set")
						return nil
					},
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "name", "Test All Assignment Types Settings Catalog Policy"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Resource name verified")
						return nil
					},
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.#", "5"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: 5 assignments verified")
						return nil
					},
					// Verify all assignment types are present
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "groupAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: groupAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: allLicensedUsersAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: allDevicesAssignmentTarget verified")
						return nil
					},
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
					func(s *terraform.State) error {
						t.Log("SUCCESS: exclusionGroupAssignmentTarget verified")
						return nil
					},
					// Verify role scope tags
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_settings_catalog_configuration_policy.assignments", "role_scope_tag_ids.#", "2"),
					func(s *terraform.State) error {
						t.Log("SUCCESS: Role scope tags verified")
						t.Log("=== STEP CHECK COMPLETE ===")
						return nil
					},
				),
			},
		},
	})
	t.Log("=== ASSIGNMENTS TEST COMPLETE ===")
}

func TestAccResourceSettingsCatalogConfigurationPolicy_03_RequiredFields(t *testing.T) {
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
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccResourceSettingsCatalogConfigurationPolicy_05_InvalidValues(t *testing.T) {
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
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform(),
				ExpectError: regexp.MustCompile("Attribute platforms value must be one of"),
			},
		},
	})
}

func testAccSettingsCatalogConfigurationPolicyConfig_minimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		log.Fatalf("Failed to load minimal test config: %v", err)
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_maximal() string {
	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}

	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		log.Fatalf("Failed to load maximal test config: %v", err)
	}

	return acceptance.ConfiguredM365ProviderBlock(roleScopeTags + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_assignments() string {
	groups, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/groups.tf")
	if err != nil {
		log.Fatalf("Failed to load groups config: %v", err)
	}

	roleScopeTags, err := helpers.ParseHCLFile("../../../../../acceptance/terraform_dependancies/device_management/role_scope_tags.tf")
	if err != nil {
		log.Fatalf("Failed to load role scope tags config: %v", err)
	}

	// Use local assignment filters with proper dependency management to preserve the correct destroy order
	//assignmentFilters, err := helpers.ParseHCLFile("assignment_filters_with_dependencies.tf")
	//if err != nil {
	//	log.Fatalf("Failed to load assignment filters config: %v", err)
	//}

	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_assignments.tf")
	if err != nil {
		log.Fatalf("Failed to load test config: %v", err)
	}

	return acceptance.ConfiguredM365ProviderBlock(groups + "\n" + roleScopeTags + "\n" + accTestConfig)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingName() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  platforms = "windows10"
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_missingPlatforms() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name = "Test Policy"
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccSettingsCatalogConfigurationPolicyConfig_invalidPlatform() string {
	config := `
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "test" {
  name      = "Test Policy"
  platforms = "invalid"
	template_reference = {
    template_id = ""
  }
  configuration_policy = {
    settings = []
  }
}
`
	return acceptance.ConfiguredM365ProviderBlock(config)
}
