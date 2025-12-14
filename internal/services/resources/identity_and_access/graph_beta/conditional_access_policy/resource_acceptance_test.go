package graphBetaConditionalAccessPolicy_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/conditional_access_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaConditionalAccessPolicy.ResourceName

	// testResource is the test resource implementation for conditional access policies
	testResource = graphBetaConditionalAccessPolicy.ConditionalAccessPolicyTestResource{}
)

// CAD001: macOS Device Compliance
func TestAccConditionalAccessPolicyResource_CAD001(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD001 macOS device compliance policy")
				},
				Config: testAccConfigCAD001(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad001_macos_compliant").ExistsInGraph(testResource),
					check.That(resourceType+".cad001_macos_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad001_macos_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad001-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad001_macos_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),

					// Conditions - Users - Exclude Guests
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),

					// Conditions - Applications
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("macOS"),

					// Grant Controls
					check.That(resourceType+".cad001_macos_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad001_macos_compliant").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD001 policy")
				},
				ResourceName:            resourceType + ".cad001_macos_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD002: Windows Device Compliance
func TestAccConditionalAccessPolicyResource_CAD002(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD002 Windows device compliance policy")
				},
				Config: testAccConfigCAD002(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad002_windows_compliant").ExistsInGraph(testResource),
					check.That(resourceType+".cad002_windows_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad002_windows_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad002-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad002_windows_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),

					// Conditions - Users - Exclude Guests
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),

					// Conditions - Applications
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),

					// Grant Controls
					check.That(resourceType+".cad002_windows_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad002_windows_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad002_windows_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad002_windows_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD002 policy")
				},
				ResourceName:            resourceType + ".cad002_windows_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD003: CAD003-O365
func TestAccConditionalAccessPolicyResource_CAD003(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD003 policy")
				},
				Config: testAccConfigCAD003(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad003-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD003 policy")
				},
				ResourceName:            resourceType + ".cad003_mobile_compliant_or_app_protection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD004: CAD004-O365
func TestAccConditionalAccessPolicyResource_CAD004(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD004 policy")
				},
				Config: testAccConfigCAD004(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad004-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD004 policy")
				},
				ResourceName:            resourceType + ".cad004_browser_noncompliant_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD005: CAD005-O365
func TestAccConditionalAccessPolicyResource_CAD005(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD005 policy")
				},
				Config: testAccConfigCAD005(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad005-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.#").HasValue("5"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("iOS"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("macOS"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("linux"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("grant_controls.built_in_controls.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD005 policy")
				},
				ResourceName:            resourceType + ".cad005_block_unsupported_platforms",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD006: CAD006-O365
func TestAccConditionalAccessPolicyResource_CAD006(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD006 policy")
				},
				Config: testAccConfigCAD006(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad006-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("session_controls.application_enforced_restrictions.is_enabled").HasValue("true"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD006 policy")
				},
				ResourceName:            resourceType + ".cad006_session_block_download_unmanaged",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD007: CAD007-O365
func TestAccConditionalAccessPolicyResource_CAD007(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD007 policy")
				},
				Config: testAccConfigCAD007(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad007-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("7"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("days"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD007 policy")
				},
				ResourceName:            resourceType + ".cad007_mobile_signin_frequency",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD008: CAD008-All
func TestAccConditionalAccessPolicyResource_CAD008(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD008 policy")
				},
				Config: testAccConfigCAD008(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad008_browser_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad008-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("days"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD008 policy")
				},
				ResourceName:            resourceType + ".cad008_browser_signin_frequency",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD009: CAD009-All
func TestAccConditionalAccessPolicyResource_CAD009(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD009 policy")
				},
				Config: testAccConfigCAD009(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad009_disable_browser_persistence").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad009-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("session_controls.persistent_browser.mode").HasValue("never"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("session_controls.persistent_browser.is_enabled").HasValue("true"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD009 policy")
				},
				ResourceName:            resourceType + ".cad009_disable_browser_persistence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD010: CAD010-RJD
func TestAccConditionalAccessPolicyResource_CAD010(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD010 policy")
				},
				Config: testAccConfigCAD010(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad010_device_registration_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad010_device_registration_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad010-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad010_device_registration_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registerdevice"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD010 policy")
				},
				ResourceName:            resourceType + ".cad010_device_registration_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD011: CAD011-O365
func TestAccConditionalAccessPolicyResource_CAD011(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD011 policy")
				},
				Config: testAccConfigCAD011(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad011_linux_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad011_linux_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad011-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad011_linux_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_users.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_users.*").ContainsTypeSetElement("GuestsOrExternalUsers"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("linux"),
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD011 policy")
				},
				ResourceName:            resourceType + ".cad011_linux_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD012: CAD012-All
func TestAccConditionalAccessPolicyResource_CAD012(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD012 policy")
				},
				Config: testAccConfigCAD012(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad012_admin_compliant_access").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad012_admin_compliant_access").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad012-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad012_admin_compliant_access").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.include_roles.#").HasValue("26"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD012 policy")
				},
				ResourceName:            resourceType + ".cad012_admin_compliant_access",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD013: CAD013-Selected
func TestAccConditionalAccessPolicyResource_CAD013(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD013 policy")
				},
				Config: testAccConfigCAD013(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad013_selected_apps_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad013-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.#").HasValue("4"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("a4f2693f-129c-4b96-982b-2c364b8314d7"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("499b84ac-1321-427f-aa17-267ca6975798"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("996def3d-b36c-4153-8607-a6fd3c01b89f"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("797f4846-ba00-4fd7-ba43-dac1f8f63013"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD013 policy")
				},
				ResourceName:            resourceType + ".cad013_selected_apps_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD014: CAD014-O365
func TestAccConditionalAccessPolicyResource_CAD014(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD014 policy")
				},
				Config: testAccConfigCAD014(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad014-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD014 policy")
				},
				ResourceName:            resourceType + ".cad014_edge_app_protection_windows",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD015: CAD015-All
func TestAccConditionalAccessPolicyResource_CAD015(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD015 policy")
				},
				Config: testAccConfigCAD015(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad015-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("macOS"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD015 policy")
				},
				ResourceName:            resourceType + ".cad015_windows_macos_browser_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD016: CAD016-EXO_SPO_CloudPC
func TestAccConditionalAccessPolicyResource_CAD016(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD016 policy")
				},
				Config: testAccConfigCAD016(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad016_token_protection_windows").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad016_token_protection_windows").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad016-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad016_token_protection_windows").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.applications.include_applications.#").HasValue("5"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad016_token_protection_windows").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD016 policy")
				},
				ResourceName:            resourceType + ".cad016_token_protection_windows",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD017: CAD017-Selected
func TestAccConditionalAccessPolicyResource_CAD017(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD017 policy")
				},
				Config: testAccConfigCAD017(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad017-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("None"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD017 policy")
				},
				ResourceName:            resourceType + ".cad017_selected_mobile_app_protection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD018: CAD018-CloudPC
func TestAccConditionalAccessPolicyResource_CAD018(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD018 policy")
				},
				Config: testAccConfigCAD018(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad018-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.applications.include_applications.#").HasValue("4"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD018 policy")
				},
				ResourceName:            resourceType + ".cad018_cloudpc_mobile_app_protection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAD019: CAD019-Intune
func TestAccConditionalAccessPolicyResource_CAD019(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAD019 policy")
				},
				Config: testAccConfigCAD019(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cad019-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAD019 policy")
				},
				ResourceName:            resourceType + ".cad019_intune_enrollment_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Found 28 CAL/CAP/CAU tests to convert

// CAL001: CAL001-All
func TestAccConditionalAccessPolicyResource_CAL001(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL001 policy")
				},
				Config: testAccConfigCAL001(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal001_block_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal001_block_locations").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal001-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal001_block_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL001 policy")
				},
				ResourceName:            resourceType + ".cal001_block_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAL002: CAL002-RSI
func TestAccConditionalAccessPolicyResource_CAL002(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL002 policy")
				},
				Config: testAccConfigCAL002(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal002-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registersecurityinfo"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL002 policy")
				},
				ResourceName:            resourceType + ".cal002_mfa_registration_trusted_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAL003: CAL003-All
func TestAccConditionalAccessPolicyResource_CAL003(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL003 policy")
				},
				Config: testAccConfigCAL003(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal003-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL003 policy")
				},
				ResourceName:            resourceType + ".cal003_block_service_accounts_untrusted",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAL004: CAL004-All
func TestAccConditionalAccessPolicyResource_CAL004(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL004 policy")
				},
				Config: testAccConfigCAL004(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal004-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.users.include_roles.#").HasValue("26"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL004 policy")
				},
				ResourceName:            resourceType + ".cal004_block_admin_untrusted_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAL005: CAL005-Selected
func TestAccConditionalAccessPolicyResource_CAL005(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL005 policy")
				},
				Config: testAccConfigCAL005(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal005-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.exclude_applications.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.exclude_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL005 policy")
				},
				ResourceName:            resourceType + ".cal005_less_trusted_locations_compliant",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAL006: CAL006-All
func TestAccConditionalAccessPolicyResource_CAL006(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAL006 policy")
				},
				Config: testAccConfigCAL006(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cal006-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAL006 policy")
				},
				ResourceName:            resourceType + ".cal006_allow_only_specified_locations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAP001: CAP001-All
func TestAccConditionalAccessPolicyResource_CAP001(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAP001 policy")
				},
				Config: testAccConfigCAP001(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cap001_block_legacy_auth").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap001_block_legacy_auth").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cap001-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cap001_block_legacy_auth").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.client_app_types.*").ContainsTypeSetElement("other"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAP001 policy")
				},
				ResourceName:            resourceType + ".cap001_block_legacy_auth",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAP002: CAP002-All
func TestAccConditionalAccessPolicyResource_CAP002(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAP002 policy")
				},
				Config: testAccConfigCAP002(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cap002_block_exchange_activesync").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cap002-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.client_app_types.*").ContainsTypeSetElement("exchangeActiveSync"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAP002 policy")
				},
				ResourceName:            resourceType + ".cap002_block_exchange_activesync",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAP003: CAP003-All
func TestAccConditionalAccessPolicyResource_CAP003(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAP003 policy")
				},
				Config: testAccConfigCAP003(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cap003_block_device_code_flow").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap003_block_device_code_flow").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cap003-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cap003_block_device_code_flow").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.authentication_flows.transfer_methods").HasValue("deviceCodeFlow"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAP003 policy")
				},
				ResourceName:            resourceType + ".cap003_block_device_code_flow",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAP004: CAP004-All
func TestAccConditionalAccessPolicyResource_CAP004(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAP004 policy")
				},
				Config: testAccConfigCAP004(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cap004_block_auth_transfer").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap004_block_auth_transfer").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cap004-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cap004_block_auth_transfer").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.authentication_flows.transfer_methods").HasValue("authenticationTransfer"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAP004 policy")
				},
				ResourceName:            resourceType + ".cap004_block_auth_transfer",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU001: CAU001-All
func TestAccConditionalAccessPolicyResource_CAU001(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU001 policy")
				},
				Config: testAccConfigCAU001(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau001_guest_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau001_guest_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau001-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau001_guest_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU001 policy")
				},
				ResourceName:            resourceType + ".cau001_guest_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU001A: CAU001A-Windows Azure Active Directory
func TestAccConditionalAccessPolicyResource_CAU001A(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU001A policy")
				},
				Config: testAccConfigCAU001A(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau001a-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU001A policy")
				},
				ResourceName:            resourceType + ".cau001a_guest_mfa_azure_ad",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU002: CAU002-All
func TestAccConditionalAccessPolicyResource_CAU002(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU002 policy")
				},
				Config: testAccConfigCAU002(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau002_all_users_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau002_all_users_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau002-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau002_all_users_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.exclude_roles.#").HasValue("23"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau002_all_users_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau002_all_users_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU002 policy")
				},
				ResourceName:            resourceType + ".cau002_all_users_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU003: CAU003-Selected
func TestAccConditionalAccessPolicyResource_CAU003(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU003 policy")
				},
				Config: testAccConfigCAU003(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau003-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU003 policy")
				},
				ResourceName:            resourceType + ".cau003_block_unapproved_apps_guests",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU004: CAU004-Selected
func TestAccConditionalAccessPolicyResource_CAU004(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU004 policy")
				},
				Config: testAccConfigCAU004(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau004_mdca_route").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau004_mdca_route").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau004-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau004_mdca_route").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),
					check.That(resourceType+".cau004_mdca_route").Key("session_controls.cloud_app_security.cloud_app_security_type").HasValue("mcasConfigured"),
					check.That(resourceType+".cau004_mdca_route").Key("session_controls.cloud_app_security.is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU004 policy")
				},
				ResourceName:            resourceType + ".cau004_mdca_route",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU006: CAU006-All
func TestAccConditionalAccessPolicyResource_CAU006(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU006 policy")
				},
				Config: testAccConfigCAU006(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau006_signin_risk_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau006-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("medium"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU006 policy")
				},
				ResourceName:            resourceType + ".cau006_signin_risk_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU007: CAU007-All
func TestAccConditionalAccessPolicyResource_CAU007(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU007 policy")
				},
				Config: testAccConfigCAU007(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau007_user_risk_password_change").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau007_user_risk_password_change").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau007-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau007_user_risk_password_change").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("medium"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.operator").HasValue("AND"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("passwordChange"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU007 policy")
				},
				ResourceName:            resourceType + ".cau007_user_risk_password_change",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU008: CAU008-All
func TestAccConditionalAccessPolicyResource_CAU008(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU008 policy")
				},
				Config: testAccConfigCAU008(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau008-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.users.include_roles.#").HasValue("26"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000004"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU008 policy")
				},
				ResourceName:            resourceType + ".cau008_admin_phishing_resistant_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU009: CAU009-Management
func TestAccConditionalAccessPolicyResource_CAU009(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU009 policy")
				},
				Config: testAccConfigCAU009(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau009_admin_portals_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau009-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.applications.include_applications.#").HasValue("2"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("MicrosoftAdminPortals"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU009 policy")
				},
				ResourceName:            resourceType + ".cau009_admin_portals_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU010: CAU010-All
func TestAccConditionalAccessPolicyResource_CAU010(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU010 policy")
				},
				Config: testAccConfigCAU010(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau010_terms_of_use").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau010_terms_of_use").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau010-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau010_terms_of_use").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau010_terms_of_use").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau010_terms_of_use").Key("grant_controls.terms_of_use.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU010 policy")
				},
				ResourceName:            resourceType + ".cau010_terms_of_use",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU011: CAU011-All
func TestAccConditionalAccessPolicyResource_CAU011(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU011 policy")
				},
				Config: testAccConfigCAU011(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau011_block_unlicensed").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau011_block_unlicensed").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau011-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau011_block_unlicensed").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.exclude_users.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.exclude_users.*").ContainsTypeSetElement("GuestsOrExternalUsers"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU011 policy")
				},
				ResourceName:            resourceType + ".cau011_block_unlicensed",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU012: CAU012-RSI
func TestAccConditionalAccessPolicyResource_CAU012(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU012 policy")
				},
				Config: testAccConfigCAU012(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau012_security_info_registration_tap").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau012-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registersecurityinfo"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU012 policy")
				},
				ResourceName:            resourceType + ".cau012_security_info_registration_tap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU013: CAU013-All
func TestAccConditionalAccessPolicyResource_CAU013(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU013 policy")
				},
				Config: testAccConfigCAU013(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau013-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000004"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU013 policy")
				},
				ResourceName:            resourceType + ".cau013_all_users_phishing_resistant_mfa",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU014: CAU014-All
func TestAccConditionalAccessPolicyResource_CAU014(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU014 policy")
				},
				Config: testAccConfigCAU014(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau014-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.*").ContainsTypeSetElement("medium"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_applications.include_service_principals.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_applications.include_service_principals.*").ContainsTypeSetElement("14ddb4bd-2aee-4603-86d2-467e438cda0a"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU014 policy")
				},
				ResourceName:            resourceType + ".cau014_block_managed_identity_risk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU015: CAU015-All
func TestAccConditionalAccessPolicyResource_CAU015(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU015 policy")
				},
				Config: testAccConfigCAU015(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau015_block_high_signin_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau015-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.sign_in_risk_levels.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU015 policy")
				},
				ResourceName:            resourceType + ".cau015_block_high_signin_risk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU016: CAU016-All
func TestAccConditionalAccessPolicyResource_CAU016(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU016 policy")
				},
				Config: testAccConfigCAU016(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau016_block_high_user_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau016_block_high_user_risk").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau016-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau016_block_high_user_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.user_risk_levels.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU016 policy")
				},
				ResourceName:            resourceType + ".cau016_block_high_user_risk",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU017: CAU017-All
func TestAccConditionalAccessPolicyResource_CAU017(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU017 policy")
				},
				Config: testAccConfigCAU017(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau017_admin_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau017-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.users.include_roles.#").HasValue("26"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("10"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("hours"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU017 policy")
				},
				ResourceName:            resourceType + ".cau017_admin_signin_frequency",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU018: CAU018-All
func TestAccConditionalAccessPolicyResource_CAU018(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU018 policy")
				},
				Config: testAccConfigCAU018(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau018-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.users.include_roles.#").HasValue("25"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("session_controls.persistent_browser.is_enabled").HasValue("true"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("session_controls.persistent_browser.mode").HasValue("never"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU018 policy")
				},
				ResourceName:            resourceType + ".cau018_admin_disable_browser_persistence",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAU019: CAU019-Selected
func TestAccConditionalAccessPolicyResource_CAU019(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU019 policy")
				},
				Config: testAccConfigCAU019(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau019-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("5"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.exclude_applications.#").HasValue("10"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU019 policy")
				},
				ResourceName:            resourceType + ".cau019_allow_only_approved_apps_guests",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAAU001: Agent ID Resources - Block All
func TestAccConditionalAccessPolicyResource_CAAU001(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAAU001 agent ID resources policy")
				},
				Config: testAccConfigCAAU001(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".caau001_all").ExistsInGraph(testResource),
					check.That(resourceType+".caau001_all").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".caau001_all").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-caau001-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".caau001_all").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".caau001_all").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".caau001_all").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".caau001_all").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".caau001_all").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),
					check.That(resourceType+".caau001_all").Key("conditions.users.exclude_users.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.users.include_groups.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.users.exclude_groups.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.users.include_roles.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.users.exclude_roles.#").HasValue("0"),

					// Conditions - Applications
					check.That(resourceType+".caau001_all").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".caau001_all").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("AllAgentIdResources"),
					check.That(resourceType+".caau001_all").Key("conditions.applications.exclude_applications.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.applications.include_user_actions.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.applications.include_authentication_context_class_references.#").HasValue("0"),

					// Conditions - Client Applications
					check.That(resourceType+".caau001_all").Key("conditions.client_applications.include_agent_id_service_principals.#").HasValue("1"),
					check.That(resourceType+".caau001_all").Key("conditions.client_applications.include_agent_id_service_principals.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".caau001_all").Key("conditions.client_applications.exclude_agent_id_service_principals.#").HasValue("0"),

					// Conditions - Agent ID Risk Levels
					check.That(resourceType+".caau001_all").Key("conditions.agent_id_risk_levels.#").HasValue("2"),
					check.That(resourceType+".caau001_all").Key("conditions.agent_id_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".caau001_all").Key("conditions.agent_id_risk_levels.*").ContainsTypeSetElement("medium"),

					// Conditions - Other Risk Levels
					check.That(resourceType+".caau001_all").Key("conditions.sign_in_risk_levels.#").HasValue("0"),
					check.That(resourceType+".caau001_all").Key("conditions.service_principal_risk_levels.#").HasValue("0"),

					// Grant Controls
					check.That(resourceType+".caau001_all").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".caau001_all").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".caau001_all").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
					check.That(resourceType+".caau001_all").Key("grant_controls.custom_authentication_factors.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAAU001 policy")
				},
				ResourceName:            resourceType + ".caau001_all",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// CAAU002: Agent ID Resources - All Applications
func TestAccConditionalAccessPolicyResource_CAAU002(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAAU002 agent ID policy targeting all applications")
				},
				Config: testAccConfigCAAU002(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".caau002_o365").ExistsInGraph(testResource),
					check.That(resourceType+".caau002_o365").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".caau002_o365").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-caau002-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".caau002_o365").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".caau002_o365").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".caau002_o365").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".caau002_o365").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.exclude_users.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.include_groups.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.exclude_groups.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.include_roles.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.users.exclude_roles.#").HasValue("0"),

					// Conditions - Applications
					check.That(resourceType+".caau002_o365").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".caau002_o365").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".caau002_o365").Key("conditions.applications.exclude_applications.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.applications.include_user_actions.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.applications.include_authentication_context_class_references.#").HasValue("0"),

					// Conditions - Client Applications
					check.That(resourceType+".caau002_o365").Key("conditions.client_applications.include_agent_id_service_principals.#").HasValue("1"),
					check.That(resourceType+".caau002_o365").Key("conditions.client_applications.include_agent_id_service_principals.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".caau002_o365").Key("conditions.client_applications.exclude_agent_id_service_principals.#").HasValue("0"),

					// Conditions - Agent ID Risk Levels
					check.That(resourceType+".caau002_o365").Key("conditions.agent_id_risk_levels.#").HasValue("2"),
					check.That(resourceType+".caau002_o365").Key("conditions.agent_id_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".caau002_o365").Key("conditions.agent_id_risk_levels.*").ContainsTypeSetElement("medium"),

					// Conditions - Other Risk Levels
					check.That(resourceType+".caau002_o365").Key("conditions.sign_in_risk_levels.#").HasValue("0"),
					check.That(resourceType+".caau002_o365").Key("conditions.service_principal_risk_levels.#").HasValue("0"),

					// Grant Controls
					check.That(resourceType+".caau002_o365").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".caau002_o365").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".caau002_o365").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
					check.That(resourceType+".caau002_o365").Key("grant_controls.custom_authentication_factors.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAAU002 policy")
				},
				ResourceName:            resourceType + ".caau002_o365",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Test config loading functions
func testAccConfigCAD001() string {
	config := mocks.LoadTerraformConfigFile("resource_cad001-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD002() string {
	config := mocks.LoadTerraformConfigFile("resource_cad002-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD003() string {
	config := mocks.LoadTerraformConfigFile("resource_cad003-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD004() string {
	config := mocks.LoadTerraformConfigFile("resource_cad004-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD005() string {
	config := mocks.LoadTerraformConfigFile("resource_cad005-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD006() string {
	config := mocks.LoadTerraformConfigFile("resource_cad006-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD007() string {
	config := mocks.LoadTerraformConfigFile("resource_cad007-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD008() string {
	config := mocks.LoadTerraformConfigFile("resource_cad008-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD009() string {
	config := mocks.LoadTerraformConfigFile("resource_cad009-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD010() string {
	config := mocks.LoadTerraformConfigFile("resource_cad010-rjd.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD011() string {
	config := mocks.LoadTerraformConfigFile("resource_cad011-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD012() string {
	config := mocks.LoadTerraformConfigFile("resource_cad012-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD013() string {
	config := mocks.LoadTerraformConfigFile("resource_cad013-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD014() string {
	config := mocks.LoadTerraformConfigFile("resource_cad014-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD015() string {
	config := mocks.LoadTerraformConfigFile("resource_cad015-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD016() string {
	config := mocks.LoadTerraformConfigFile("resource_cad016-exo_spo_cloudpc.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD017() string {
	config := mocks.LoadTerraformConfigFile("resource_cad017-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD018() string {
	config := mocks.LoadTerraformConfigFile("resource_cad018-cloudpc.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAD019() string {
	config := mocks.LoadTerraformConfigFile("resource_cad019-intune.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL001() string {
	config := mocks.LoadTerraformConfigFile("resource_cal001-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL002() string {
	config := mocks.LoadTerraformConfigFile("resource_cal002-rsi.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL003() string {
	config := mocks.LoadTerraformConfigFile("resource_cal003-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL004() string {
	config := mocks.LoadTerraformConfigFile("resource_cal004-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL005() string {
	config := mocks.LoadTerraformConfigFile("resource_cal005-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAL006() string {
	config := mocks.LoadTerraformConfigFile("resource_cal006-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAP001() string {
	config := mocks.LoadTerraformConfigFile("resource_cap001-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAP002() string {
	config := mocks.LoadTerraformConfigFile("resource_cap002-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAP003() string {
	config := mocks.LoadTerraformConfigFile("resource_cap003-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAP004() string {
	config := mocks.LoadTerraformConfigFile("resource_cap004-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU001() string {
	config := mocks.LoadTerraformConfigFile("resource_cau001-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU001A() string {
	config := mocks.LoadTerraformConfigFile("resource_cau001a-windows_azure_active_directory.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU002() string {
	config := mocks.LoadTerraformConfigFile("resource_cau002-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU003() string {
	config := mocks.LoadTerraformConfigFile("resource_cau003-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU004() string {
	config := mocks.LoadTerraformConfigFile("resource_cau004-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU006() string {
	config := mocks.LoadTerraformConfigFile("resource_cau006-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU007() string {
	config := mocks.LoadTerraformConfigFile("resource_cau007-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU008() string {
	config := mocks.LoadTerraformConfigFile("resource_cau008-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU009() string {
	config := mocks.LoadTerraformConfigFile("resource_cau009-management.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU010() string {
	config := mocks.LoadTerraformConfigFile("resource_cau010-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU011() string {
	config := mocks.LoadTerraformConfigFile("resource_cau011-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU012() string {
	config := mocks.LoadTerraformConfigFile("resource_cau012-rsi.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU013() string {
	config := mocks.LoadTerraformConfigFile("resource_cau013-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU014() string {
	config := mocks.LoadTerraformConfigFile("resource_cau014-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU015() string {
	config := mocks.LoadTerraformConfigFile("resource_cau015-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU016() string {
	config := mocks.LoadTerraformConfigFile("resource_cau016-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU017() string {
	config := mocks.LoadTerraformConfigFile("resource_cau017-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU018() string {
	config := mocks.LoadTerraformConfigFile("resource_cau018-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAU019() string {
	config := mocks.LoadTerraformConfigFile("resource_cau019-selected.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAAU001() string {
	config := mocks.LoadTerraformConfigFile("resource_caau001-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccConfigCAAU002() string {
	config := mocks.LoadTerraformConfigFile("resource_caau002-o365.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccConditionalAccessPolicyResource_CAU020 tests the conditional access policy resource
// with insider risk levels for all users.
func TestAccConditionalAccessPolicyResource_CAU020(t *testing.T) {
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
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating CAU020 insider risk policy")
				},
				Config: testAccConfigCAU020(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".cau020_all").ExistsInGraph(testResource),
					check.That(resourceType+".cau020_all").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau020_all").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-cau020-[^:]+: .+ [a-z0-9]{8}$`)),
					check.That(resourceType+".cau020_all").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau020_all").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau020_all").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cau020_all").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau020_all").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau020_all").Key("conditions.users.exclude_users.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.users.include_groups.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cau020_all").Key("conditions.users.include_roles.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.users.exclude_roles.#").HasValue("0"),

					// Conditions - Applications
					check.That(resourceType+".cau020_all").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau020_all").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau020_all").Key("conditions.applications.exclude_applications.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.applications.include_user_actions.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.applications.include_authentication_context_class_references.#").HasValue("0"),

					// Conditions - Risk Levels
					check.That(resourceType+".cau020_all").Key("conditions.sign_in_risk_levels.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.user_risk_levels.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.service_principal_risk_levels.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.agent_id_risk_levels.#").HasValue("0"),
					check.That(resourceType+".cau020_all").Key("conditions.insider_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau020_all").Key("conditions.insider_risk_levels.*").ContainsTypeSetElement("moderate"),
					check.That(resourceType+".cau020_all").Key("conditions.insider_risk_levels.*").ContainsTypeSetElement("elevated"),

					// Grant Controls
					check.That(resourceType+".cau020_all").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau020_all").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau020_all").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
					check.That(resourceType+".cau020_all").Key("grant_controls.custom_authentication_factors.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing CAU020 policy")
				},
				Config:            testAccConfigCAU020(),
				ResourceName:      resourceType + ".cau020_all",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigCAU020() string {
	config := mocks.LoadTerraformConfigFile("resource_cau020-all.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}
