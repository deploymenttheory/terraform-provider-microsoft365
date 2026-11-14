package graphBetaIdentityAndAccessDeviceRegistrationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaIdentityAndAccessDeviceRegistrationPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/device_registration_policy"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessDeviceRegistrationPolicy.ResourceName

	// testResource is the test resource implementation for the device registration policy
	testResource = graphBetaIdentityAndAccessDeviceRegistrationPolicy.DeviceRegistrationPolicyTestResource{}

	// externalProviders are the providers the acceptance fixtures use besides microsoft365
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
	}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// acceptanceImportStep verifies import by the fixed singleton ID.
func acceptanceImportStep() resource.TestStep {
	return resource.TestStep{
		PreConfig: func() {
			testlog.StepAction(resourceType, "Importing the device registration policy")
		},
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "deviceRegistrationPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// The device registration policy is a tenant-wide singleton. Every test records the original policy in
// its PreCheck and a t.Cleanup PUTs it back and verifies the restore with GET (R9). CheckDestroy asserts
// that destroy left the policy in the tenant (R7) instead of expecting it to be deleted.

// Minimal: only user_device_quota is configured.
func TestAccResourceDeviceRegistrationPolicy_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { deviceRegistrationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the device registration policy with a minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("deviceRegistrationPolicy"),
					check.That(resourceType+".test").Key("user_device_quota").HasValue("49"),
					check.That(resourceType+".test").Key("multi_factor_auth_configuration").MatchesRegex(regexp.MustCompile(`^(notRequired|required)$`)),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").Exists(),
					check.That(resourceType+".test").Key("azure_ad_registration.is_admin_configurable").Exists(),
					check.That(resourceType+".test").Key("azure_ad_registration.allowed_to_register.type").Exists(),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").Exists(),
				),
			},
			acceptanceImportStep(),
		},
	})
}

// Maximal: every configurable value is set, with an acc-test group allowed to join devices.
func TestAccResourceDeviceRegistrationPolicy_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { deviceRegistrationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Applying the maximal device registration policy configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("user_device_quota").HasValue("45"),
					check.That(resourceType+".test").Key("multi_factor_auth_configuration").HasValue("notRequired"),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("azure_ad_registration.allowed_to_register.type").HasValue("all"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test", "azure_ad_join.allowed_to_join.groups.*", "microsoft365_graph_beta_groups_group.join", "id"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.enable_global_admins").HasValue("true"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.registering_users.type").HasValue("none"),
				),
			},
			acceptanceImportStep(),
		},
	})
}

// QuotaUpdate: user_device_quota changes in place.
func TestAccResourceDeviceRegistrationPolicy_03_QuotaUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { deviceRegistrationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Setting the device quota")
				},
				Config: loadAcceptanceTestTerraform("resource_quota_initial.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("user_device_quota").HasValue("48"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating the device quota")
				},
				Config: loadAcceptanceTestTerraform("resource_quota_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("user_device_quota").HasValue("47"),
				),
			},
			acceptanceImportStep(),
		},
	})
}

// MembershipTypeChanges: allowed_to_join moves from all to enumerated (an acc-test group) to none.
func TestAccResourceDeviceRegistrationPolicy_04_MembershipTypeChanges(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { deviceRegistrationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Allowing all users to join devices")
				},
				Config: loadAcceptanceTestTerraform("resource_join_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("all"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").DoesNotExist(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Allowing an acc-test group to join devices")
				},
				Config: loadAcceptanceTestTerraform("resource_join_enumerated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("1"),
					resource.TestCheckTypeSetElemAttrPair(resourceType+".test", "azure_ad_join.allowed_to_join.groups.*", "microsoft365_graph_beta_groups_group.join", "id"),
				),
			},
			acceptanceImportStep(),
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Allowing nobody to join devices")
				},
				Config: loadAcceptanceTestTerraform("resource_join_none.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("none"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").DoesNotExist(),
				),
			},
		},
	})
}

// LAPSToggle: local_admin_password.is_enabled is switched off and on again.
func TestAccResourceDeviceRegistrationPolicy_05_LAPSToggle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { deviceRegistrationPolicyPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyStillExists,
		ExternalProviders:        externalProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Disabling Windows LAPS")
				},
				Config: loadAcceptanceTestTerraform("resource_laps_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("local_admin_password.is_enabled").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Enabling Windows LAPS")
				},
				Config: loadAcceptanceTestTerraform("resource_laps_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("local_admin_password.is_enabled").HasValue("true"),
				),
			},
			acceptanceImportStep(),
		},
	})
}
