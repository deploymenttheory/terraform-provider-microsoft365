package graphBetaMacosDeviceCompliancePolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/macos_device_compliance_policy/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.MacosDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &policyMocks.MacosDeviceCompliancePolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestUnitResourceMacosDeviceCompliancePolicy_01_OmittedRuleName verifies the scheduled-action lifecycle.
func TestUnitResourceMacosDeviceCompliancePolicy_01_OmittedRuleName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_omitted_rule_name.tf"),
				Check:  scheduledActionChecks("PasswordRequired", "72"),
			},
			{
				Config:   loadUnitTestTerraform("resource_omitted_rule_name.tf"),
				PlanOnly: true,
			},
			{
				Config: loadUnitTestTerraform("resource_updated_actions.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: scheduledActionChecks("PasswordRequired", "24"),
			},
			{
				Config:   loadUnitTestTerraform("resource_updated_actions.tf"),
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

// TestUnitResourceMacosDeviceCompliancePolicy_02_ExplicitRuleName verifies the scheduled-action lifecycle.
func TestUnitResourceMacosDeviceCompliancePolicy_02_ExplicitRuleName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_explicit_rule_name.tf"),
				Check:  scheduledActionChecks("PasswordRequired", "72"),
			},
			{
				Config:   loadUnitTestTerraform("resource_explicit_rule_name.tf"),
				PlanOnly: true,
			},
			{
				Config: loadUnitTestTerraform("resource_updated_explicit_rule_name.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: scheduledActionChecks("PasswordRequired", "24"),
			},
			{
				Config:   loadUnitTestTerraform("resource_updated_explicit_rule_name.tf"),
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

// TestUnitResourceMacosDeviceCompliancePolicy_03_LegacyRuleName verifies the scheduled-action lifecycle.
func TestUnitResourceMacosDeviceCompliancePolicy_03_LegacyRuleName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_legacy_rule_name.tf"),
				Check:  scheduledActionChecks("unavailable", "72"),
			},
			{
				Config:   loadUnitTestTerraform("resource_legacy_rule_name.tf"),
				PlanOnly: true,
			},
			{
				Config: loadUnitTestTerraform("resource_updated_explicit_rule_name.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: scheduledActionChecks("PasswordRequired", "24"),
			},
			{
				Config:   loadUnitTestTerraform("resource_updated_explicit_rule_name.tf"),
				PlanOnly: true,
			},
		},
	})
}

// TestUnitResourceMacosDeviceCompliancePolicy_04_Assignments verifies the scheduled-action lifecycle.
func TestUnitResourceMacosDeviceCompliancePolicy_04_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_assignment_initial.tf"),
				Check:  scheduledActionChecks("PasswordRequired", "72"),
			},
			{
				Config:   loadUnitTestTerraform("resource_assignment_initial.tf"),
				PlanOnly: true,
			},
			{
				Config: loadUnitTestTerraform("resource_assignment_added.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					scheduledActionChecks("PasswordRequired", "72"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_assignment_updated.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					scheduledActionChecks("PasswordRequired", "24"),
					check.That(resourceType+".test").Key("assignments.#").HasValue("1"),
				),
			},
			{
				Config:   loadUnitTestTerraform("resource_assignment_updated.tf"),
				PlanOnly: true,
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config: loadUnitTestTerraform("resource_assignment_removed.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					scheduledActionChecks("PasswordRequired", "24"),
					resource.TestCheckNoResourceAttr(resourceType+".test", "assignments.#"),
				),
			},
			{
				Config:   loadUnitTestTerraform("resource_assignment_removed.tf"),
				PlanOnly: true,
			},
		},
	})
}

func scheduledActionChecks(ruleName, gracePeriod string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(resourceType+".test").Key("id").Exists(),
		check.That(resourceType+".test").Key("scheduled_actions_for_rule.#").HasValue("1"),
		check.That(resourceType+".test").
			Key("scheduled_actions_for_rule.0.rule_name").
			HasValue(ruleName),
		check.That(resourceType+".test").
			Key("scheduled_actions_for_rule.0.scheduled_action_configurations.#").
			HasValue("1"),
		check.That(resourceType+".test").
			Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.action_type").
			HasValue("block"),
		check.That(resourceType+".test").
			Key("scheduled_actions_for_rule.0.scheduled_action_configurations.0.grace_period_hours").
			HasValue(gracePeriod),
	)
}
