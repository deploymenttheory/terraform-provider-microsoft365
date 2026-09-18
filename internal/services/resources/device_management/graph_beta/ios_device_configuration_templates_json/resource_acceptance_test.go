package graphBetaIosDeviceConfigurationTemplatesJson_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIosDeviceConfigurationTemplatesJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/ios_device_configuration_templates_json"
)

var (
	resourceType = graphBetaIosDeviceConfigurationTemplatesJson.ResourceName
	testResource = graphBetaIosDeviceConfigurationTemplatesJson.IosDeviceConfigurationTemplatesJsonTestResource{}
)

func loadAcceptanceTestTerraform(t *testing.T, filename string) string {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		t.Skipf("skipping acceptance test: fixture file not found: %s", filename)
		return ""
	}
	return config
}

func TestAccResourceIosDeviceConfigurationTemplatesJson_01_GeneralConfiguration(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating general device configuration")
				},
				Config: loadAcceptanceTestTerraform(t, "resource_general_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("iOS/iPadOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".general_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".general_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".general_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-iOS-json-general-[a-z0-9]{8}$`)),
					check.That(resourceType+".general_minimal").Key("settings_json").HasValue(`{"cameraBlocked":true}`),
					// hard-delete test disabled
					// check.That(resourceType+".general_minimal").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying no drift against a live tenant")
				},
				Config:             loadAcceptanceTestTerraform(t, "resource_general_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating general device configuration (exercises PATCH)")
				},
				Config: loadAcceptanceTestTerraform(t, "resource_general_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".general_minimal").Key("settings_json").MatchesRegex(
						regexp.MustCompile(`"airDropBlocked":true`),
					),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing general device configuration")
				},
				ResourceName:      resourceType + ".general_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"settings_json"},
			},
		},
	})
}

func TestAccResourceIosDeviceConfigurationTemplatesJson_02_DeviceFeaturesConfiguration(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating device features configuration")
				},
				Config: loadAcceptanceTestTerraform(t, "resource_device_features_home_screen.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("iOS/iPadOS device configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".device_features").ExistsInGraph(testResource),
					check.That(resourceType+".device_features").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".device_features").Key("settings_json").MatchesRegex(
						regexp.MustCompile(`#microsoft\.graph\.iosHomeScreenApp`),
					),
					check.That(resourceType+".device_features").Key("settings_json").MatchesRegex(
						regexp.MustCompile(`#microsoft\.graph\.iosHomeScreenFolder`),
					),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Verifying the home screen tree does not drift")
				},
				Config:             loadAcceptanceTestTerraform(t, "resource_device_features_home_screen.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
