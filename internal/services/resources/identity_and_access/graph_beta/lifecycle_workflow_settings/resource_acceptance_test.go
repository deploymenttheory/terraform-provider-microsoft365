package graphBetaIdentityAndAccessLifecycleWorkflowSettings_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	graphBetaIdentityAndAccessLifecycleWorkflowSettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow_settings"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaIdentityAndAccessLifecycleWorkflowSettings.ResourceName

	// testResource is the test resource implementation for the lifecycle workflow settings
	testResource = graphBetaIdentityAndAccessLifecycleWorkflowSettings.LifecycleWorkflowSettingsTestResource{}
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
			testlog.StepAction(resourceType, "Importing the lifecycle workflow settings")
		},
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "lifecycleWorkflowSettings",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// The lifecycle workflow settings are a tenant-wide singleton. Every test records the original settings in
// its PreCheck and a t.Cleanup PATCHes them back and verifies the restore with GET (R9). CheckDestroy asserts
// that destroy left the settings in the tenant (R7) instead of expecting them to be deleted.

// Minimal: only the schedule interval is configured.
func TestAccResourceLifecycleWorkflowSettings_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { lifecycleWorkflowSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Adopting the lifecycle workflow settings with a minimal configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue("lifecycleWorkflowSettings"),
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("5"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").MatchesRegex(regexp.MustCompile(`^\S+$`)),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").MatchesRegex(regexp.MustCompile(`^(true|false)$`)),
				),
			},
			acceptanceImportStep(),
		},
	})
}

// Maximal: every configurable value is set, then email_settings stops being managed and keeps its value.
func TestAccResourceLifecycleWorkflowSettings_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { lifecycleWorkflowSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Applying the maximal lifecycle workflow settings configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("7"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
			acceptanceImportStep(),
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing email_settings from the configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("5"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
		},
	})
}

// Update: the interval and company branding change in place.
func TestAccResourceLifecycleWorkflowSettings_03_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { lifecycleWorkflowSettingsPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsStillExist,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Setting the schedule interval and enabling company branding")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("7"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating the schedule interval and disabling company branding")
				},
				Config: loadAcceptanceTestTerraform("resource_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("9"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("false"),
				),
			},
			acceptanceImportStep(),
		},
	})
}
