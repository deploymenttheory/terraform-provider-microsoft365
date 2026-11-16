package graphBetaEndpointSecurityIntent_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaEndpointSecurityIntent "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/endpoint_security_intent"
	graphBetaGroup "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaEndpointSecurityIntent.ResourceName

	// testResource is the test resource implementation for endpoint security intents
	testResource = graphBetaEndpointSecurityIntent.EndpointSecurityIntentTestResource{}

	groupResourceType = graphBetaGroup.ResourceName
	groupTestResource = graphBetaGroup.GroupTestResource{}
)

const accBitLockerTemplateId = "d1174162-1dd2-4976-affc-6667049ab0ae"

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

func accTestCase(t *testing.T, steps []resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        map[string]resource.ExternalProvider{"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion}},
		CheckDestroy:             destroy.CheckDestroyedTypesFunc(30*time.Second, destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource}, destroy.ResourceTypeMapping{ResourceType: groupResourceType, TestResource: groupTestResource}),
		Steps:                    steps,
	}
}

// Creates a BitLocker intent, adds settings and a group assignment, changes the settings in place and removes them again.
// The intent ID is checked at every step so a settings change can never recreate the intent.
func TestAccResourceEndpointSecurityIntent_01_Lifecycle(t *testing.T) {
	address := resourceType + ".test"
	var intentId string

	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating minimal BitLocker intent")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(address).ExistsInGraph(testResource),
				resource.TestCheckResourceAttrWith(address, "id", func(value string) error {
					intentId = value
					return nil
				}),
				check.That(address).Key("template_id").HasValue(accBitLockerTemplateId),
				check.That(address).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-esi-[a-z0-9]{8}$`)),
				check.That(address).Key("description").HasValue(""),
				check.That(address).Key("role_scope_tag_ids.#").HasValue("1"),
				check.That(address).Key("settings.#").DoesNotExist(),
				check.That(address).Key("assignments.#").DoesNotExist(),
				check.That(address).Key("is_assigned").HasValue("false"),
				check.That(address).Key("last_modified_date_time").IsSet(),
			),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Importing minimal BitLocker intent")
			},
			ResourceName:            address,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts"},
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Adding settings and a group assignment")
			},
			Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(address).ExistsInGraph(testResource),
				resource.TestCheckResourceAttrPtr(address, "id", &intentId),
				check.That(address).Key("description").HasValue("Acceptance test BitLocker intent"),
				check.That(address).Key("settings.#").HasValue("2"),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "setting_type": "boolean", "value_json": "true"}),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "setting_type": "string", "value_json": `"enabledForAzureAd"`}),
				check.That(address).Key("assignments.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(address, "assignments.*.group_id", groupResourceType+".test", "id"),
				check.That(address).Key("is_assigned").HasValue("true"),
			),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Importing BitLocker intent with its settings")
			},
			ResourceName:            address,
			ImportState:             true,
			ImportStateIdFunc:       importIdWithSettings(address, encryptDeviceId, recoveryRotationId),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts"},
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Updating settings in place and switching to an exclusion assignment")
			},
			Config: loadAcceptanceTestTerraform("resource_settings_updated.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(address).ExistsInGraph(testResource),
				resource.TestCheckResourceAttrPtr(address, "id", &intentId),
				check.That(address).Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-esi-updated-[a-z0-9]{8}$`)),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "value_json": "false"}),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "value_json": `"enabledForAzureAdAndHybrid"`}),
				check.That(address).Key("assignments.#").HasValue("1"),
				resource.TestCheckTypeSetElemNestedAttrs(address, "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget"}),
			),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Confirming an empty plan after the settings update")
			},
			Config:   loadAcceptanceTestTerraform("resource_settings_updated.tf"),
			PlanOnly: true,
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Removing settings and assignments")
			},
			Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(address).ExistsInGraph(testResource),
				resource.TestCheckResourceAttrPtr(address, "id", &intentId),
				check.That(address).Key("settings.#").DoesNotExist(),
				check.That(address).Key("assignments.#").DoesNotExist(),
				check.That(address).Key("is_assigned").HasValue("false"),
			),
		},
	}))
}

// Creates an intent with settings and an assignment in one apply, then imports it.
func TestAccResourceEndpointSecurityIntent_02_CreateWithSettings(t *testing.T) {
	address := resourceType + ".test"

	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Creating BitLocker intent with settings and an exclusion assignment")
			},
			Config: loadAcceptanceTestTerraform("resource_settings_updated.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(address).ExistsInGraph(testResource),
				check.That(address).Key("settings.#").HasValue("2"),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": encryptDeviceId, "setting_type": "boolean", "value_json": "false"}),
				resource.TestCheckTypeSetElemNestedAttrs(address, "settings.*", map[string]string{"definition_id": recoveryRotationId, "setting_type": "string", "value_json": `"enabledForAzureAdAndHybrid"`}),
				check.That(address).Key("assignments.#").HasValue("1"),
				resource.TestCheckTypeSetElemAttrPair(address, "assignments.*.group_id", groupResourceType+".test", "id"),
			),
		},
		{
			PreConfig: func() {
				testlog.StepAction(resourceType, "Importing BitLocker intent with its settings")
			},
			ResourceName:            address,
			ImportState:             true,
			ImportStateIdFunc:       importIdWithSettings(address, encryptDeviceId, recoveryRotationId),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts"},
		},
	}))
}
