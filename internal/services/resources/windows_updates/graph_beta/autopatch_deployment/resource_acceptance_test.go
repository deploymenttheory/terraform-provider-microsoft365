package graphBetaWindowsUpdatesAutopatchDeployment_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsUpdatesAutopatchDeploymentResource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_deployment"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	testResource = WindowsUpdatesAutopatchDeploymentResource.WindowsUpdateDeploymentTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceWindowsUpdateDeployment_01_FeatureUpdateDeployment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create feature update deployment with single monitoring rule")
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P7D"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("100"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import deployment")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_02_MinimalToFull(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create minimal deployment (no settings)")
				},
				Config: loadAcceptanceTestTerraform("02_minimal_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update to full settings with single monitoring rule (min to max)")
					testlog.WaitForConsistency(fmt.Sprintf("%s (settings update)", resourceType), 10*time.Second)
					time.Sleep(10 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P7D"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("100"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_03_FullToMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create full deployment with multiple monitoring rules (rollback + ineligible/offerFallback)")
				},
				Config: loadAcceptanceTestTerraform("03_feature_update_multiple_rules.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P14D"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("200"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("2"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update to minimal deployment (max to min)")
					testlog.WaitForConsistency(fmt.Sprintf("%s (settings removal)", resourceType), 10*time.Second)
					time.Sleep(10 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("02_minimal_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_04_MultipleMonitoringRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment with multiple monitoring rules (rollback + ineligible/offerFallback)")
				},
				Config: loadAcceptanceTestTerraform("03_feature_update_multiple_rules.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("2"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("200"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Reduce to single monitoring rule (multiple to single)")
					testlog.WaitForConsistency(fmt.Sprintf("%s (rule reduction)", resourceType), 10*time.Second)
					time.Sleep(10 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_05_RollbackAlertError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment with rollback/alertError monitoring rule")
				},
				Config: loadAcceptanceTestTerraform("05_rollback_alert_error.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P7D"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import deployment")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestAccResourceWindowsUpdateDeployment_06_IneligibleOfferFallback(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create deployment with ineligible/offerFallback monitoring rule (no threshold)")
				},
				Config: loadAcceptanceTestTerraform("06_ineligible_offer_fallback.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import deployment")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

