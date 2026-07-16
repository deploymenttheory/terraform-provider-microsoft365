package graphBetaWindowsCustomConfiguration_test

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
	graphBetaWindowsCustomConfiguration "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_custom_configuration"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaWindowsCustomConfiguration.ResourceName
	testResource = graphBetaWindowsCustomConfiguration.WindowsCustomConfigurationTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsCustomConfiguration_01_Lifecycle(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating windows custom configuration")
				},
				Config: loadAcceptanceTestTerraform("resource_windows_custom_configuration_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows custom configuration", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".custom_configuration_example").ExistsInGraph(testResource),
					check.That(resourceType+".custom_configuration_example").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_configuration_example").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-custom-config-[a-z0-9]{8}$`)),
					check.That(resourceType+".custom_configuration_example").Key("description").HasValue("Example Windows custom configuration profile using OMA-URI settings"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.#").HasValue("3"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.odata_type").HasValue("#microsoft.graph.omaSettingString"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.0.value").Exists(),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.2.odata_type").HasValue("#microsoft.graph.omaSettingInteger"),
					check.That(resourceType+".custom_configuration_example").Key("oma_settings.2.value").HasValue("30"),
					check.That(resourceType+".custom_configuration_example").Key("assignments.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing windows custom configuration")
				},
				ResourceName:      resourceType + ".custom_configuration_example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
