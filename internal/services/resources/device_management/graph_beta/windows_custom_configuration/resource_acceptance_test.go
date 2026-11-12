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
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
	name := resourceType + ".custom_configuration_example"
	var profileID string
	minimal := loadAcceptanceTestTerraform("resource_windows_custom_configuration_minimal.tf")
	maximal := loadAcceptanceTestTerraform("resource_windows_custom_configuration_maximal.tf")
	updated := loadAcceptanceTestTerraform("resource_windows_custom_configuration_updated.tf")
	var steps []resource.TestStep
	for i, phase := range []struct {
		name, config, description, settings, assignments, integerValue string
	}{
		{"Creating minimal configuration", minimal, "Minimal Windows custom configuration profile", "1", "", ""},
		{"Expanding to maximal configuration", maximal, "Example Windows custom configuration profile using OMA-URI settings", "3", "2", "30"},
		{"Updating maximal configuration", updated, "Updated Windows custom configuration profile using OMA-URI settings", "3", "2", "60"},
		{"Reducing to minimal configuration", minimal, "Minimal Windows custom configuration profile", "1", "", ""},
		{"Restoring maximal configuration", maximal, "Example Windows custom configuration profile using OMA-URI settings", "3", "2", "30"},
	} {
		checks := []resource.TestCheckFunc{
			check.That(name).ExistsInGraph(testResource),
			check.That(name).Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
			check.That(name).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-windows-custom-config-[a-z0-9]{8}$`)),
			check.That(name).Key("description").HasValue(phase.description),
			check.That(name).Key("oma_settings.#").HasValue(phase.settings),
			check.That(name).Key("oma_settings.0.odata_type").HasValue("#microsoft.graph.omaSettingString"),
			check.That(name).Key("oma_settings.0.value").Exists(),
		}
		action := plancheck.ResourceActionUpdate
		if i == 0 {
			action = plancheck.ResourceActionCreate
			checks = append(checks, resource.TestCheckResourceAttrWith(name, "id", func(value string) error { profileID = value; return nil }))
		} else {
			checks = append(checks, resource.TestCheckResourceAttrPtr(name, "id", &profileID))
		}
		if phase.assignments == "" {
			checks = append(checks, resource.TestCheckNoResourceAttr(name, "assignments.#"))
		} else {
			checks = append(checks,
				check.That(name).Key("assignments.#").HasValue(phase.assignments),
				check.That(name).Key("oma_settings.2.odata_type").HasValue("#microsoft.graph.omaSettingInteger"),
				check.That(name).Key("oma_settings.2.value").HasValue(phase.integerValue),
			)
			mode := "manual"
			if phase.integerValue == "60" {
				mode = "none"
			}
			checks = append(checks, check.That(name).Key("oma_settings.1.value").HasValue("<enabled/>\n<data id=\"UpdateMode\" value=\""+mode+"\"/>"))
		}
		steps = append(steps,
			resource.TestStep{
				PreConfig: func() { testlog.StepAction(resourceType, phase.name) },
				Config:    phase.config,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(name, action),
				}},
				Check: resource.ComposeTestCheckFunc(checks...),
			},
			resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true},
			resource.TestStep{Config: phase.config, PlanOnly: true},
		)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: graphBetaGroup.ResourceName, TestResource: graphBetaGroup.GroupTestResource{}},
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: steps,
	})
}
