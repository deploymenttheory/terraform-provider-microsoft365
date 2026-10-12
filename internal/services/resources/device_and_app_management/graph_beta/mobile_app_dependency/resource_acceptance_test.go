package graphBetaMobileAppDependency_test

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

func TestAccResourceMobileAppDependency_01_Detect(t *testing.T) {
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
				PreConfig: func() { testlog.StepAction(resourceType, "01: Creating Win32 apps and a detect dependency") },
				Config:    loadAcceptanceTestTerraform("01_detect.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(win32.ResourceName+`.apps["main"]`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps["runtime"]`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps["main"]`).
						Key("committed_content_version").
						IsNotEmpty(),
					check.That(win32.ResourceName+`.apps["runtime"]`).
						Key("committed_content_version").
						IsNotEmpty(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("source_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["main"]`).Key("id")),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["runtime"]`).Key("id")),
					check.That(resourceType+".test").Key("dependency_type").HasValue("detect"),
					check.That(resourceType+".test").Key("target_type").HasValue("child"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "01: Importing the dependency") },
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("01_detect.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceMobileAppDependency_02_AutoInstallToDetect(t *testing.T) {
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
				PreConfig: func() { testlog.StepAction(resourceType, "02: Creating Win32 apps and an auto-install dependency") },
				Config:    loadAcceptanceTestTerraform("02_auto_install.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("source_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["main"]`).Key("id")),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["runtime"]`).Key("id")),
					check.That(resourceType+".test").Key("dependency_type").HasValue("autoInstall"),
					check.That(resourceType+".test").Key("timeouts.create").HasValue("5m"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "02: Changing the dependency to detect in place") },
				Config:    loadAcceptanceTestTerraform("01_detect.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(
						resourceType+".test",
						plancheck.ResourceActionUpdate,
					),
				}},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("dependency_type").HasValue("detect"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "02: Importing the updated dependency") },
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("01_detect.tf"), PlanOnly: true},
		},
	})
}

func TestAccResourceMobileAppDependency_03_SiblingRemoval(t *testing.T) {
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
				PreConfig: func() { testlog.StepAction(resourceType, "03: Creating Win32 apps and two dependencies") },
				Config:    loadAcceptanceTestTerraform("03_siblings_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".sibling").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["runtime"]`).Key("id")),
					check.That(resourceType+".sibling").
						Key("target_id").
						MatchesOtherKey(check.That(win32.ResourceName+`.apps["second_runtime"]`).Key("id")),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "03: Removing one dependency and retaining its sibling") },
				Config:    loadAcceptanceTestTerraform("03_siblings_step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".sibling").ExistsInGraph(testResource),
					check.That(resourceType+".sibling").
						Key("dependency_type").
						HasValue("autoInstall"),
				),
			},
			{
				PreConfig:               func() { testlog.StepAction(resourceType, "03: Importing the remaining dependency") },
				ResourceName:            resourceType + ".sibling",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("03_siblings_step2.tf"), PlanOnly: true},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "03: Removing the last dependency while retaining the apps") },
				Config:    loadAcceptanceTestTerraform("03_siblings_step3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(win32.ResourceName+`.apps["main"]`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps["runtime"]`).
						ExistsInGraph(win32.Win32AppTestResource{}),
					check.That(win32.ResourceName+`.apps["second_runtime"]`).
						ExistsInGraph(win32.Win32AppTestResource{}),
				),
			},
		},
	})
}
