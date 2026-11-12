package graphBetaMacOSDepEnrollmentProfile_test

import (
	"log"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMacOSDepEnrollmentProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_dep_enrollment_profile"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceType = graphBetaMacOSDepEnrollmentProfile.ResourceName

var testResource = graphBetaMacOSDepEnrollmentProfile.MacOSDepEnrollmentProfileTestResource{}

func TestAccResourceMacOSDepEnrollmentProfile_01_Lifecycle(t *testing.T) {
	var profileID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
		Steps: lifecycleSteps(&profileID),
	})
}

// lifecycleSteps exercises all wizard configurations on the same enrollment profile.
func lifecycleSteps(profileID *string) []resource.TestStep {
	address := resourceType + ".test"
	scenarios := []struct {
		name     string
		affinity string
		modern   string
		skips    string
		rotation string
	}{
		{"minimal", "false", "false", "0", ""},
		{"modern_auth", "true", "true", "27", "90"},
		{"legacy_auth", "true", "false", "27", "180"},
		{"without_user_affinity", "false", "false", "27", ""},
		{"minimal", "false", "false", "0", ""},
	}
	var steps []resource.TestStep
	for i, scenario := range scenarios {
		config := testAccMacOSDepEnrollmentProfileConfig(scenario.name)
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(address, "requires_user_authentication", scenario.affinity),
			resource.TestCheckResourceAttr(address, "configuration_web_url", scenario.modern),
			resource.TestCheckResourceAttr(address, "enabled_skip_keys.#", scenario.skips),
			resource.TestCheckResourceAttrSet(address, "configuration_endpoint_url"),
			resource.TestCheckResourceAttrSet(address, "dep_onboarding_settings_id"),
		}
		if i == 0 {
			checks = append(checks, resource.TestCheckResourceAttrWith(address, "id", func(value string) error { *profileID = value; return nil }))
		} else {
			checks = append(checks, resource.TestCheckResourceAttrPtr(address, "id", profileID))
		}
		if scenario.rotation == "" {
			checks = append(checks, resource.TestCheckNoResourceAttr(address, "admin_account_password_rotation.auto_rotation_period_in_days"))
		} else {
			checks = append(checks,
				resource.TestCheckResourceAttr(address, "admin_account_password_rotation.auto_rotation_period_in_days", scenario.rotation),
				resource.TestCheckNoResourceAttr(address, "admin_account_password_rotation.on_retrieval_auto_rotate_password_enabled"),
				resource.TestCheckNoResourceAttr(address, "admin_account_password_rotation.on_retrieval_delay_auto_rotate_password_in_hours"),
			)
		}
		steps = append(steps,
			resource.TestStep{Config: config, Check: resource.ComposeTestCheckFunc(checks...)},
			resource.TestStep{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"admin_account_password"}},
			resource.TestStep{Config: config, PlanOnly: true},
		)
	}
	return steps
}

// testAccMacOSDepEnrollmentProfileConfig loads a wizard scenario with the acceptance provider.
func testAccMacOSDepEnrollmentProfileConfig(scenario string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_" + scenario + ".tf")
	if err != nil {
		log.Fatalf("Failed to load %s test config: %v", scenario, err)
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}
