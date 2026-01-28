package graphBetaApplication_test

import (
	"fmt"
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

func TestAccResourceApplication_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating minimal application")
				},
				Config: testAccConfig01Minimal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-minimal-[a-z0-9]+$`)),
					check.That(resourceType+".test_minimal").Key("description").HasValue("Minimal acceptance test application"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName: resourceType + ".test_minimal",
				ImportState:  true,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs, ok := s.RootModule().Resources[resourceType+".test_minimal"]
				if !ok {
					return "", fmt.Errorf("resource not found: %s", resourceType+".test_minimal")
				}
				preventDuplicateNames := rs.Primary.Attributes["prevent_duplicate_names"]
				hardDelete := rs.Primary.Attributes["hard_delete"]
				return fmt.Sprintf("%s:prevent_duplicate_names=%s:hard_delete=%s", rs.Primary.ID, preventDuplicateNames, hardDelete), nil
			},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceApplication_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating maximal application")
				},
				Config: testAccConfig02Maximal(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-maximal-[a-z0-9]+$`)),
					check.That(resourceType+".test_maximal").Key("description").HasValue("Maximal acceptance test application with all fields configured"),
					check.That(resourceType+".test_maximal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_maximal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("terraform"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("acceptance-test"),
					check.That(resourceType+".test_maximal").Key("tags.*").ContainsTypeSetElement("maximal"),
					check.That(resourceType+".test_maximal").Key("owner_user_ids.#").HasValue("2"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Import state verification")
				},
				ResourceName: resourceType + ".test_maximal",
				ImportState:  true,
			ImportStateIdFunc: func(s *terraform.State) (string, error) {
				rs, ok := s.RootModule().Resources[resourceType+".test_maximal"]
				if !ok {
					return "", fmt.Errorf("resource not found: %s", resourceType+".test_maximal")
				}
				preventDuplicateNames := rs.Primary.Attributes["prevent_duplicate_names"]
				hardDelete := rs.Primary.Attributes["hard_delete"]
				return fmt.Sprintf("%s:prevent_duplicate_names=%s:hard_delete=%s", rs.Primary.ID, preventDuplicateNames, hardDelete), nil
			},
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceApplication_03_WebApp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating web application")
				},
				Config: testAccConfig03WebApp(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_web_app").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_web_app").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-web-app-[a-z0-9]+$`)),
					check.That(resourceType+".test_web_app").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_web_app").Key("web.home_page_url").Exists(),
					check.That(resourceType+".test_web_app").Key("web.redirect_uris.#").HasValue("1"),
				),
			},
		},
	})
}

func TestAccResourceApplication_04_SPA(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating SPA application")
				},
				Config: testAccConfig04SPA(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_spa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_spa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-spa-[a-z0-9]+$`)),
					check.That(resourceType+".test_spa").Key("sign_in_audience").HasValue("AzureADMultipleOrgs"),
					check.That(resourceType+".test_spa").Key("spa.redirect_uris.#").HasValue("2"),
				),
			},
		},
	})
}

func TestAccResourceApplication_05_PublicClient(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating public client application")
				},
				Config: testAccConfig05PublicClient(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_public_client").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_public_client").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-public-client-[a-z0-9]+$`)),
					check.That(resourceType+".test_public_client").Key("is_fallback_public_client").HasValue("true"),
					check.That(resourceType+".test_public_client").Key("public_client.redirect_uris.#").HasValue("2"),
				),
			},
		},
	})
}

func TestAccResourceApplication_06_Multitenant(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
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
					testlog.StepAction(resourceType, "Step 1: Creating multitenant application")
				},
				Config: testAccConfig06Multitenant(),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_multitenant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_multitenant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-multitenant-[a-z0-9]+$`)),
					check.That(resourceType+".test_multitenant").Key("sign_in_audience").HasValue("AzureADandPersonalMicrosoftAccount"),
					check.That(resourceType+".test_multitenant").Key("tags.#").HasValue("2"),
					check.That(resourceType+".test_multitenant").Key("web.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_multitenant").Key("spa.redirect_uris.#").HasValue("1"),
				),
			},
		},
	})
}

func testAccConfig01Minimal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_01_minimal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfig02Maximal() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_02_maximal.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfig03WebApp() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_03_web_app.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfig04SPA() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_04_spa.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfig05PublicClient() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_05_public_client.tf")
	if err != nil {
		panic(err)
	}
	return content
}

func testAccConfig06Multitenant() string {
	content, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_06_multitenant.tf")
	if err != nil {
		panic(err)
	}
	return content
}
