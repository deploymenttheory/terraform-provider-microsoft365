package graphBetaM365AdminSharepointSettings_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaM365AdminSharepointSettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/m365_admin/graph_beta/sharepoint_settings"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaM365AdminSharepointSettings.ResourceName

	// testResource is the test resource implementation for the SharePoint settings
	testResource = graphBetaM365AdminSharepointSettings.SharepointSettingsTestResource{}
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
			testlog.StepAction(resourceType, "Importing the SharePoint settings")
		},
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "sharepointSettings",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// The SharePoint settings are a tenant-wide singleton. Every test records the original settings in its PreCheck
// and a t.Cleanup PATCHes them back and verifies the restore with GET (R9). CheckDestroy asserts that destroy left
// the settings in the tenant (R7) instead of expecting them to be deleted. The tests only change low-impact settings:
// the deleted user OneDrive retention period, the sync app excluded file extensions, the Mac sync app, the OneDrive
// sync button and idle session sign-out. Sharing, Loop, legacy authentication and unmanaged device settings are never changed.

// Minimal: only the retention period is configured.
func TestAccResourceSharepointSettings_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sharepointSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the SharePoint settings with a minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("sharepointSettings"),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("45"),
					check.That(resourceType+".test").Key("sharing_capability").MatchesRegex(regexp.MustCompile(`^\S+$`)),
					check.That(resourceType+".test").Key("available_managed_paths_for_site_creation.#").MatchesRegex(regexp.MustCompile(`^[1-9]\d*$`)),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").MatchesRegex(regexp.MustCompile(`^(true|false)$`)),
				),
			},
			acceptanceImportStep(),
		},
	})
}

// Maximal: every low-impact value is set, then the configuration stops managing them and they keep their values.
func TestAccResourceSharepointSettings_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sharepointSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Applying the maximal low-impact SharePoint settings configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("60"),
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.#").HasValue("1"),
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.*").ContainsTypeSetElement("*.acc-test-tmp"),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("is_sync_button_hidden_on_personal_site").HasValue("true"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
				),
			},
			acceptanceImportStep(),
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing the managed settings except the retention period from the configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("45"),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("is_sync_button_hidden_on_personal_site").HasValue("true"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
				),
			},
		},
	})
}

// Update: the low-impact values change in place and idle session sign-out is turned back on.
func TestAccResourceSharepointSettings_03_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { sharepointSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Disabling the Mac sync app and idle session sign-out")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("60"),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Enabling the Mac sync app and idle session sign-out with long timeouts")
				},
				Config: loadAcceptanceTestTerraform("resource_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("90"),
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.#").HasValue("2"),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("is_sync_button_hidden_on_personal_site").HasValue("false"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("idle_session_sign_out.warn_after_in_seconds").HasValue("3300"),
					check.That(resourceType+".test").Key("idle_session_sign_out.sign_out_after_in_seconds").HasValue("3600"),
				),
			},
			acceptanceImportStep(),
		},
	})
}
