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

func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

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
				Config: loadAccTestTerraform("resource_01_minimal.tf"),
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
				Config: loadAccTestTerraform("resource_02_maximal.tf"),
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating web application")
				},
				Config: loadAccTestTerraform("resource_03_web_app.tf"),
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating SPA application")
				},
				Config: loadAccTestTerraform("resource_04_spa.tf"),
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
				Config: loadAccTestTerraform("resource_05_public_client.tf"),
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating multitenant application")
				},
				Config: loadAccTestTerraform("resource_06_multitenant.tf"),
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

func TestAccResourceApplication_07_MinimalToMaximal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating minimal application")
				},
				Config: loadAccTestTerraform("resource_07_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal_to_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-min-to-max-[a-z0-9]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("description").HasValue("Minimal to maximal test application - step 1"),
					check.That(resourceType+".test_minimal_to_maximal").Key("owner_user_ids.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to maximal configuration")
				},
				Config: loadAccTestTerraform("resource_07_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_minimal_to_maximal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-min-to-max-[a-z0-9]+$`)),
					check.That(resourceType+".test_minimal_to_maximal").Key("description").HasValue("Minimal to maximal test application - step 2"),
					check.That(resourceType+".test_minimal_to_maximal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_minimal_to_maximal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_minimal_to_maximal").Key("tags.*").ContainsTypeSetElement("terraform"),
					check.That(resourceType+".test_minimal_to_maximal").Key("tags.*").ContainsTypeSetElement("acceptance-test"),
					check.That(resourceType+".test_minimal_to_maximal").Key("tags.*").ContainsTypeSetElement("maximal"),
					check.That(resourceType+".test_minimal_to_maximal").Key("owner_user_ids.#").HasValue("2"),
					check.That(resourceType+".test_minimal_to_maximal").Key("group_membership_claims.#").HasValue("1"),
					check.That(resourceType+".test_minimal_to_maximal").Key("notes").HasValue("This is a test application for acceptance testing"),
					check.That(resourceType+".test_minimal_to_maximal").Key("is_device_only_auth_supported").HasValue("false"),
					check.That(resourceType+".test_minimal_to_maximal").Key("is_fallback_public_client").HasValue("false"),
					check.That(resourceType+".test_minimal_to_maximal").Key("service_management_reference").HasValue("https://contoso.com/app-management"),
					check.That(resourceType+".test_minimal_to_maximal").Key("api.accept_mapped_claims").HasValue("true"),
					check.That(resourceType+".test_minimal_to_maximal").Key("api.requested_access_token_version").HasValue("2"),
					check.That(resourceType+".test_minimal_to_maximal").Key("app_roles.#").HasValue("3"),
					check.That(resourceType+".test_minimal_to_maximal").Key("info.marketing_url").HasValue("https://contoso.com/marketing"),
					check.That(resourceType+".test_minimal_to_maximal").Key("info.privacy_statement_url").HasValue("https://contoso.com/privacy"),
					check.That(resourceType+".test_minimal_to_maximal").Key("info.support_url").HasValue("https://contoso.com/support"),
					check.That(resourceType+".test_minimal_to_maximal").Key("info.terms_of_service_url").HasValue("https://contoso.com/terms"),
					check.That(resourceType+".test_minimal_to_maximal").Key("parental_control_settings.countries_blocked_for_minors.#").HasValue("2"),
					check.That(resourceType+".test_minimal_to_maximal").Key("parental_control_settings.legal_age_group_rule").HasValue("Allow"),
					check.That(resourceType+".test_minimal_to_maximal").Key("public_client.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_minimal_to_maximal").Key("spa.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_minimal_to_maximal").Key("web.home_page_url").HasValue("https://contoso.com"),
					check.That(resourceType+".test_minimal_to_maximal").Key("web.logout_url").HasValue("https://contoso.com/logout"),
					check.That(resourceType+".test_minimal_to_maximal").Key("web.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_minimal_to_maximal").Key("web.implicit_grant_settings.enable_access_token_issuance").HasValue("false"),
					check.That(resourceType+".test_minimal_to_maximal").Key("web.implicit_grant_settings.enable_id_token_issuance").HasValue("true"),
				),
			},
		},
	})
}

func TestAccResourceApplication_08_MaximalToMinimal(t *testing.T) {
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
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 1: Creating maximal application")
				},
				Config: loadAccTestTerraform("resource_08_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal_to_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-max-to-min-[a-z0-9]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("description").HasValue("Maximal to minimal test application - step 1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_maximal_to_minimal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_maximal_to_minimal").Key("tags.*").ContainsTypeSetElement("terraform"),
					check.That(resourceType+".test_maximal_to_minimal").Key("tags.*").ContainsTypeSetElement("acceptance-test"),
					check.That(resourceType+".test_maximal_to_minimal").Key("tags.*").ContainsTypeSetElement("maximal"),
					check.That(resourceType+".test_maximal_to_minimal").Key("owner_user_ids.#").HasValue("2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("group_membership_claims.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("notes").HasValue("This is a test application for acceptance testing"),
					check.That(resourceType+".test_maximal_to_minimal").Key("is_device_only_auth_supported").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("is_fallback_public_client").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("service_management_reference").HasValue("https://contoso.com/app-management"),
					check.That(resourceType+".test_maximal_to_minimal").Key("api.accept_mapped_claims").HasValue("true"),
					check.That(resourceType+".test_maximal_to_minimal").Key("api.requested_access_token_version").HasValue("2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("app_roles.#").HasValue("3"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.marketing_url").HasValue("https://contoso.com/marketing"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.privacy_statement_url").HasValue("https://contoso.com/privacy"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.support_url").HasValue("https://contoso.com/support"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.terms_of_service_url").HasValue("https://contoso.com/terms"),
					check.That(resourceType+".test_maximal_to_minimal").Key("parental_control_settings.countries_blocked_for_minors.#").HasValue("2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("parental_control_settings.legal_age_group_rule").HasValue("Allow"),
					check.That(resourceType+".test_maximal_to_minimal").Key("public_client.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("spa.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.home_page_url").HasValue("https://contoso.com"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.logout_url").HasValue("https://contoso.com/logout"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.implicit_grant_settings.enable_access_token_issuance").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.implicit_grant_settings.enable_id_token_issuance").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Step 2: Updating to minimal configuration")
				},
				Config: loadAccTestTerraform("resource_08_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						testlog.WaitForConsistency("application", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".test_maximal_to_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("app_id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-app-max-to-min-[a-z0-9]+$`)),
					check.That(resourceType+".test_maximal_to_minimal").Key("description").HasValue("Maximal to minimal test application - step 2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("owner_user_ids.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("sign_in_audience").HasValue("AzureADMyOrg"),
					check.That(resourceType+".test_maximal_to_minimal").Key("tags.#").HasValue("3"),
					check.That(resourceType+".test_maximal_to_minimal").Key("group_membership_claims.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("notes").HasValue("This is a test application for acceptance testing"),
					check.That(resourceType+".test_maximal_to_minimal").Key("service_management_reference").HasValue("https://contoso.com/app-management"),
					check.That(resourceType+".test_maximal_to_minimal").Key("is_device_only_auth_supported").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("is_fallback_public_client").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("api.accept_mapped_claims").HasValue("true"),
					check.That(resourceType+".test_maximal_to_minimal").Key("api.requested_access_token_version").HasValue("2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("app_roles.#").HasValue("3"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.marketing_url").HasValue("https://contoso.com/marketing"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.privacy_statement_url").HasValue("https://contoso.com/privacy"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.support_url").HasValue("https://contoso.com/support"),
					check.That(resourceType+".test_maximal_to_minimal").Key("info.terms_of_service_url").HasValue("https://contoso.com/terms"),
					check.That(resourceType+".test_maximal_to_minimal").Key("parental_control_settings.countries_blocked_for_minors.#").HasValue("2"),
					check.That(resourceType+".test_maximal_to_minimal").Key("parental_control_settings.legal_age_group_rule").HasValue("Allow"),
					check.That(resourceType+".test_maximal_to_minimal").Key("public_client.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("spa.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.home_page_url").HasValue("https://contoso.com"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.logout_url").HasValue("https://contoso.com/logout"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.redirect_uris.#").HasValue("1"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.implicit_grant_settings.enable_access_token_issuance").HasValue("false"),
					check.That(resourceType+".test_maximal_to_minimal").Key("web.implicit_grant_settings.enable_id_token_issuance").HasValue("true"),
				),
			},
		},
	})
}

