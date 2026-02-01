package graphBetaServicePrincipal_test

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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func TestAccResourceServicePrincipal_01_Minimal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Step 1: Creating service principal with minimal configuration")
				},
				Config: loadAccTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("service principal", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").Exists(),
					check.That(resourceType+".test_minimal").Key("account_enabled").Exists(),
					check.That(resourceType+".test_minimal").Key("service_principal_type").HasValue("Application"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceServicePrincipal_02_Maximal(t *testing.T) {
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
					testlog.StepAction(resourceType, "Step 1: Creating service principal with maximal configuration")
				},
				Config: loadAccTestTerraform("resource_02_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("service principal", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").Exists(),
					check.That(resourceType+".test_maximal").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("app_role_assignment_required").HasValue("true"),
					check.That(resourceType+".test_maximal").Key("description").HasValue("Maximal service principal configuration for testing"),
					check.That(resourceType+".test_maximal").Key("login_url").HasValue("https://login.example.com"),
					check.That(resourceType+".test_maximal").Key("notes").HasValue("Service principal for maximal acceptance testing"),
					check.That(resourceType+".test_maximal").Key("notification_email_addresses.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("preferred_single_sign_on_mode").HasValue("saml"),
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("2"),
					check.That(resourceType+".test_maximal").Key("service_principal_type").HasValue("Application"),
					check.That(resourceType+".test_maximal").Key("sign_in_audience").Exists(),
					check.That(resourceType+".test_maximal").Key("service_principal_names.#").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName:      resourceType + ".test_maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceServicePrincipal_03_Update(t *testing.T) {
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
					testlog.StepAction(resourceType, "Step 1: Creating service principal for update test")
				},
				Config: loadAccTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("service principal", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating service principal properties")
				},
				Config: `
resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-sp-minimal-${random_string.test_id.result}"
  description  = "Application for service principal minimal acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id                       = microsoft365_graph_beta_applications_application.test.app_id
  app_role_assignment_required = true
  tags                         = ["HideApp", "WindowsAzureActiveDirectoryIntegratedApp"]

  depends_on = [time_sleep.wait_for_app]
}`,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("service principal", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("app_role_assignment_required").HasValue("true"),
					check.That(resourceType+".test_minimal").Key("tags.#").HasValue("2"),
				),
			},
		},
	})
}
