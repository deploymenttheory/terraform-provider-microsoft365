package graphBetaMobileAppSupersedence_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	win32 "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/win32_app"
)

// loadAcceptanceTestTerraform loads acceptance test Terraform configuration files.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceMobileAppSupersedence_01_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			t.Setenv("TF_VAR_installer_zip", prepareAcceptanceInstaller(t))
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: win32.ResourceName,
				TestResource: win32.Win32AppTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "01: Creating Win32 apps and a update supersedence") },
				Config:    loadAcceptanceTestTerraform("01_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(win32.ResourceName+`.apps.0`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps.1`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps.0`).
						Key("committed_content_version").
						IsNotEmpty(),
					check.That(win32.ResourceName+`.apps.1`).
						Key("committed_content_version").
						IsNotEmpty(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("source_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.0`).Key("id")),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.1`).Key("id")),
					check.That(resourceType+".test").Key("supersedence_type").HasValue("update"),
					check.That(resourceType+".test").Key("target_type").HasValue("child"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "01: Importing the supersedence") },
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("01_update.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceMobileAppSupersedence_02_ReplaceToUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			t.Setenv("TF_VAR_installer_zip", prepareAcceptanceInstaller(t))
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: win32.ResourceName,
				TestResource: win32.Win32AppTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "02: Creating Win32 apps and a replace supersedence") },
				Config:    loadAcceptanceTestTerraform("02_replace.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("source_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.0`).Key("id")),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.1`).Key("id")),
					check.That(resourceType+".test").Key("supersedence_type").HasValue("replace"),
					check.That(resourceType+".test").Key("timeouts.create").HasValue("5m"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "02: Changing the supersedence to update in place") },
				Config:    loadAcceptanceTestTerraform("01_update.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(
						resourceType+".test",
						plancheck.ResourceActionUpdate,
					),
				}},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("supersedence_type").HasValue("update"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "02: Importing the updated supersedence") },
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("01_update.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceMobileAppSupersedence_03_SiblingRemoval(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			t.Setenv("TF_VAR_installer_zip", prepareAcceptanceInstaller(t))
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
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
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{
				ResourceType: win32.ResourceName,
				TestResource: win32.Win32AppTestResource{},
			},
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "03: Creating Win32 apps and two supersedence relationships") },
				Config:    loadAcceptanceTestTerraform("03_siblings_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".sibling").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.1`).Key("id")),
					check.That(resourceType+".sibling").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps.2`).Key("id")),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "03: Removing one supersedence and retaining its sibling") },
				Config:    loadAcceptanceTestTerraform("03_siblings_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".sibling").ExistsInGraph(testResource),
					check.That(resourceType+".sibling").
						Key("supersedence_type").
						HasValue("replace"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "03: Importing the remaining supersedence") },
				ResourceName:            resourceType + ".sibling",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("03_siblings_step2.tf"), PlanOnly: true},
			{
				PreConfig: func() {
					testlog.StepAction(
						resourceType,
						"03: Removing the last supersedence while retaining the apps",
					)
				},
				Config: loadAcceptanceTestTerraform("03_siblings_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(win32.ResourceName+`.apps.0`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps.1`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps.2`).
						ExistsInGraph(win32.Win32AppTestResource{}),
				),
			},
		},
	})
}
