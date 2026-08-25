package graphBetaDeviceAndAppManagementAppAssignment_test

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
	graphBetaIosStoreApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_store_app"
	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_assignment"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ============================================================================
// Test Strategy & Timing Considerations
// ============================================================================
//
// These tests exercise the intent/settings compatibility rules the Intune
// service enforces on Apple app assignments. The service rejects an
// unsupported combination with an HTTP 400 rather than ignoring the field:
//
//	isRemovable + any intent other than required
//	    -> "IsRemovable setting is only supported for Required intent."
//	uninstallOnDeviceRemoval + uninstall intent
//	    -> "UninstallOnDeviceRemoval setting is not supported Uninstall intent."
//	preventManagedAppBackup + uninstall intent
//	    -> "PreventManagedAppBackup setting is not supported for Uninstall intent."
//
// So the assertions below are as much about which attributes are *absent* from
// state as which are present: an attribute that is absent is one the provider
// did not send.
//
// ## Dependency propagation
//
// Every scenario creates its own iOS store app and target group, then waits on
// a time_sleep before creating the assignment. Intune validates the group and
// app referenced by an assignment against its own copy of directory data, which
// lags group creation by a few seconds. Without the pause the create
// intermittently fails with a 400 naming an unknown group.
//
// ## Replacement rather than update (test 005)
//
// Both intent and settings carry RequiresReplace, and Graph rejects a PATCH
// that carries intent, source or target ("Cannot patch read-only properties"),
// so changing an assignment always destroys and recreates it. Terraform
// destroys before creating, so the replacement does not collide with the
// inclusion intent the original assignment holds on the same group. A
// consistency wait between the two steps keeps the refresh in step 2 from
// seeing the old assignment as drift.
//
// ## Group hard delete
//
// Groups are created with hard_delete = true, which is a two phase operation in
// Graph. CheckDestroy therefore waits 60 seconds before asserting, since the
// permanent delete can take 60-90 seconds to propagate.
//
// ## No import steps
//
// ImportState is a passthrough on id alone, but mobile_app_id is required to
// address an assignment and is not part of the imported state, so Read cannot
// resolve an imported resource. Import is not exercised here; it needs fixing
// separately before it can be tested.
// ============================================================================

// loadAcceptanceTestTerraform loads an acceptance test config and prepends the provider block
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

var (
	testResource = graphBetaMobileAppAssignment.MobileAppAssignmentTestResource{}

	// externalProviders are needed by every scenario that builds its own dependencies:
	// random for the name suffix, time for the propagation pause.
	externalProviders = map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: constants.ExternalProviderTimeVersion,
		},
	}
)

// checkDestroyAllTypes asserts the assignment and everything it was built on are gone.
func checkDestroyAllTypes() resource.TestCheckFunc {
	return destroy.CheckDestroyedTypesFunc(
		60*time.Second, // groups with hard_delete can take 60-90s to fully propagate
		destroy.ResourceTypeMapping{
			ResourceType: graphBetaMobileAppAssignment.ResourceName,
			TestResource: graphBetaMobileAppAssignment.MobileAppAssignmentTestResource{},
		},
		destroy.ResourceTypeMapping{
			ResourceType: graphBetaIosStoreApp.ResourceName,
			TestResource: graphBetaIosStoreApp.IosStoreAppTestResource{},
		},
		destroy.ResourceTypeMapping{
			ResourceType: graphBetaGroup.ResourceName,
			TestResource: graphBetaGroup.GroupTestResource{},
		},
	)
}

// Test 001: available intent with is_removable omitted — the configuration reported in #3692
func TestAccResourceMobileAppAssignment_01_Ios_Available_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders,
		CheckDestroy:             checkDestroyAllTypes(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_ios_available_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").ExistsInGraph(testResource),
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_\d+_\d+$`)),
					check.That(resourceType+".test_001").Key("intent").HasValue("available"),
					check.That(resourceType+".test_001").Key("source").HasValue("direct"),
					check.That(resourceType+".test_001").Key("target.target_type").HasValue("groupAssignment"),
					check.That(resourceType+".test_001").Key("target.group_id").IsUUID(),
					// Absent from state means absent from the request. With the old schema
					// defaults these were populated and the create failed with a 400.
					check.That(resourceType+".test_001").Key("settings.ios_store.is_removable").DoesNotExist(),
					check.That(resourceType+".test_001").Key("settings.ios_store.uninstall_on_device_removal").DoesNotExist(),
					check.That(resourceType+".test_001").Key("settings.ios_store.prevent_managed_app_backup").DoesNotExist(),
				),
			},
		},
	})
}

// Test 002: required intent with every intent-restricted setting configured explicitly
func TestAccResourceMobileAppAssignment_02_Ios_Required_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders,
		CheckDestroy:             checkDestroyAllTypes(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_ios_required_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").ExistsInGraph(testResource),
					check.That(resourceType+".test_002").Key("intent").HasValue("required"),
					check.That(resourceType+".test_002").Key("settings.ios_store.is_removable").HasValue("true"),
					check.That(resourceType+".test_002").Key("settings.ios_store.prevent_managed_app_backup").HasValue("true"),
					check.That(resourceType+".test_002").Key("settings.ios_store.uninstall_on_device_removal").HasValue("true"),
					check.That(resourceType+".test_002").Key("target.device_and_app_management_assignment_filter_type").HasValue("none"),
				),
			},
		},
	})
}

// Test 003: uninstall intent, which accepts no Apple app settings at all
func TestAccResourceMobileAppAssignment_03_Ios_Uninstall_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders,
		CheckDestroy:             checkDestroyAllTypes(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("003_ios_uninstall_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_003").ExistsInGraph(testResource),
					check.That(resourceType+".test_003").Key("intent").HasValue("uninstall"),
					check.That(resourceType+".test_003").Key("settings.ios_store.is_removable").DoesNotExist(),
					check.That(resourceType+".test_003").Key("settings.ios_store.uninstall_on_device_removal").DoesNotExist(),
					check.That(resourceType+".test_003").Key("settings.ios_store.prevent_managed_app_backup").DoesNotExist(),
				),
			},
		},
	})
}

// Test 004: availableWithoutEnrollment, the BYOD case from the issue
func TestAccResourceMobileAppAssignment_04_Ios_AvailableWithoutEnrollment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders,
		CheckDestroy:             checkDestroyAllTypes(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("004_ios_available_without_enrollment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_004").ExistsInGraph(testResource),
					check.That(resourceType+".test_004").Key("intent").HasValue("availableWithoutEnrollment"),
					// Accepted for this intent, so it must still be sent.
					check.That(resourceType+".test_004").Key("settings.ios_store.prevent_managed_app_backup").HasValue("true"),
					check.That(resourceType+".test_004").Key("settings.ios_store.is_removable").DoesNotExist(),
				),
			},
		},
	})
}

// Test 005: lifecycle from a required assignment with settings to an available one without.
// This transition was impossible before the fix: the schema defaults put isRemovable back
// into the request for the available intent and the service returned a 400.
func TestAccResourceMobileAppAssignment_05_Ios_Lifecycle_RequiredToAvailable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders,
		CheckDestroy:             checkDestroyAllTypes(),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("005_ios_lifecycle_required_to_available_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").ExistsInGraph(testResource),
					check.That(resourceType+".test_005").Key("intent").HasValue("required"),
					check.That(resourceType+".test_005").Key("settings.ios_store.is_removable").HasValue("true"),
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("mobile app assignment intent change", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
				),
			},
			{
				Config: loadAcceptanceTestTerraform("005_ios_lifecycle_required_to_available_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_005").ExistsInGraph(testResource),
					check.That(resourceType+".test_005").Key("intent").HasValue("available"),
					check.That(resourceType+".test_005").Key("settings.ios_store.prevent_managed_app_backup").HasValue("true"),
					check.That(resourceType+".test_005").Key("settings.ios_store.is_removable").DoesNotExist(),
				),
			},
		},
	})
}

// Test 006: is_removable alongside an intent that does not support it is rejected at plan
// time, rather than as an opaque 400 part-way through the apply
func TestAccResourceMobileAppAssignment_06_Ios_Validation_IsRemovableUnsupportedIntent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadAcceptanceTestTerraform("006_ios_validation_is_removable_unsupported_intent.tf"),
				ExpectError: regexp.MustCompile(`(?s)is_removable.{0,4}\s+can\s+only\s+be\s+set\s+when`),
			},
		},
	})
}

// Test 007: the two settings the service refuses specifically for an uninstall intent are
// rejected at plan time as well
func TestAccResourceMobileAppAssignment_07_Ios_Validation_UninstallUnsupportedSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadAcceptanceTestTerraform("007_ios_validation_uninstall_unsupported_settings.tf"),
				ExpectError: regexp.MustCompile(`(?s)uninstall_on_device_removal.{0,4}\s+cannot\s+be\s+set\s+when`),
			},
		},
	})
}
