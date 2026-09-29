package graphBetaMacosDeviceCompliancePolicy_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaMacosDeviceCompliancePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_compliance_policy"
)

var (
	resourceType = graphBetaMacosDeviceCompliancePolicy.ResourceName
	testResource = graphBetaMacosDeviceCompliancePolicy.MacosDeviceCompliancePolicyTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// TestAccResourceMacosDeviceCompliancePolicy_01_OmittedRuleName verifies the scheduled-action lifecycle.
func TestAccResourceMacosDeviceCompliancePolicy_01_OmittedRuleName(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating compliance policy")
				},
				Config: loadAcceptanceTestTerraform("resource_omitted_rule_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "72"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_omitted_rule_name.tf"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating scheduled actions")
				},
				Config: loadAcceptanceTestTerraform("resource_updated_actions.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_updated_actions.tf"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceMacosDeviceCompliancePolicy_02_ExplicitRuleName verifies the scheduled-action lifecycle.
func TestAccResourceMacosDeviceCompliancePolicy_02_ExplicitRuleName(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating compliance policy")
				},
				Config: loadAcceptanceTestTerraform("resource_explicit_rule_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "72"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_explicit_rule_name.tf"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating scheduled actions")
				},
				Config: loadAcceptanceTestTerraform("resource_updated_explicit_rule_name.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_updated_explicit_rule_name.tf"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceMacosDeviceCompliancePolicy_03_LegacyRuleName verifies the scheduled-action lifecycle.
func TestAccResourceMacosDeviceCompliancePolicy_03_LegacyRuleName(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating compliance policy")
				},
				Config: loadAcceptanceTestTerraform("resource_legacy_rule_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("unavailable", "72"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_legacy_rule_name.tf"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating scheduled actions")
				},
				Config: loadAcceptanceTestTerraform("resource_updated_explicit_rule_name.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_updated_explicit_rule_name.tf"),
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceMacosDeviceCompliancePolicy_04_Assignments verifies the scheduled-action lifecycle.
func TestAccResourceMacosDeviceCompliancePolicy_04_Assignments(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating compliance policy")
				},
				Config: loadAcceptanceTestTerraform("resource_assignment_initial.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "72"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_assignment_initial.tf"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Assigning policy to an empty test group")
				},
				Config: loadAcceptanceTestTerraform("resource_assignment_added.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "72"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating scheduled actions")
				},
				Config: loadAcceptanceTestTerraform("resource_assignment_updated.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_assignment_updated.tf"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Removing policy assignment")
				},
				Config: loadAcceptanceTestTerraform("resource_assignment_removed.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
					resource.TestCheckNoResourceAttr(resourceType+".test", "assignments.#"),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_assignment_removed.tf"),
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceMacosDeviceCompliancePolicy_05_EmptyNotificationTemplate verifies the scheduled-action lifecycle.
func TestAccResourceMacosDeviceCompliancePolicy_05_EmptyNotificationTemplate(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating compliance policy")
				},
				Config: loadAcceptanceTestTerraform("resource_empty_notification_template.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "72"),
					check.That(resourceType+".test").
						Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.notification_template_id").
						HasValue(""),
				),
			},
			{
				Config:   loadAcceptanceTestTerraform("resource_empty_notification_template.tf"),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating scheduled actions")
				},
				Config: loadAcceptanceTestTerraform(
					"resource_updated_empty_notification_template.tf",
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					scheduledActionChecks("PasswordRequired", "24"),
					check.That(resourceType+".test").
						Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.notification_template_id").
						HasValue(""),
				),
			},
			{
				Config: loadAcceptanceTestTerraform(
					"resource_updated_empty_notification_template.tf",
				),
				PlanOnly: true,
			},
		},
	})
}
