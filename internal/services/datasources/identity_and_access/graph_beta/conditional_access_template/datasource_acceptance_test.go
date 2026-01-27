package graphBetaConditionalAccessTemplate_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaConditionalAccessPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/conditional_access_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// policyResourceType is the resource type name for conditional access policies
	policyResourceType = graphBetaConditionalAccessPolicy.ResourceName

	// testResource is the test resource implementation for conditional access policies
	testResource = graphBetaConditionalAccessPolicy.ConditionalAccessPolicyTestResource{}
)

func TestAccDatasourceConditionalAccessTemplate_01_ByTemplateId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAccTestTerraform("01_by_template_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Top level attributes
					check.That("data."+dataSourceType+".by_template_id").Key("template_id").HasValue("c7503427-338e-4c5e-902d-abe252abfb43"),
					check.That("data."+dataSourceType+".by_template_id").Key("name").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("description").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("scenarios.#").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("id").IsSet(),

					// Details object
					check.That("data."+dataSourceType+".by_template_id").Key("details").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls").IsSet(),

					// Details - Conditions - Applications
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.include_applications").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.applications.exclude_applications").IsSet(),

					// Details - Conditions - Users
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_users").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_users").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_groups").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_groups").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.include_roles").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.conditions.users.exclude_roles").IsSet(),

					// Details - Grant Controls
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.operator").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.built_in_controls").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.custom_authentication_factors").IsSet(),
					check.That("data."+dataSourceType+".by_template_id").Key("details.grant_controls.terms_of_use").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_02_ByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config: loadAccTestTerraform("02_by_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Top level attributes
					check.That("data."+dataSourceType+".by_name").Key("name").HasValue("Require multifactor authentication for admins"),
					check.That("data."+dataSourceType+".by_name").Key("template_id").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("description").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("scenarios.#").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("id").IsSet(),

					// Details object
					check.That("data."+dataSourceType+".by_name").Key("details").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls").IsSet(),

					// Details - Conditions - Applications
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.include_applications").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.applications.exclude_applications").IsSet(),

					// Details - Conditions - Users
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_users").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_users").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_groups").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_groups").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.include_roles").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.conditions.users.exclude_roles").IsSet(),

					// Details - Grant Controls
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.operator").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.built_in_controls").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.custom_authentication_factors").IsSet(),
					check.That("data."+dataSourceType+".by_name").Key("details.grant_controls.terms_of_use").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_03_FuzzyMatchSuggestion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      loadAccTestTerraform("03_invalid_name_fuzzy.tf"),
				ExpectError: regexp.MustCompile("(?s)Invalid Template Name.*No conditional access template found with name: Require MFA for admin.*Did you mean one of these"),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_04_MFAAdmins(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from MFA for admins template")
				},
				Config: loadAccTestTerraform("04_require_mfa_for_admins.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_admins").Key("name").HasValue("Require multifactor authentication for admins"),
					check.That("data."+dataSourceType+".mfa_admins").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_admins").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-admins-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_admins").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_admins").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_05_BlockLegacyAuth(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from block legacy authentication template")
				},
				Config: loadAccTestTerraform("05_block_legacy_authentication.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".block_legacy_auth").Key("name").HasValue("Block legacy authentication"),
					check.That("data."+dataSourceType+".block_legacy_auth").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_block_legacy_auth").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-block-legacy-authentication-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_block_legacy_auth").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_block_legacy_auth").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_06_SecuringSecurityInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from securing security info registration template")
				},
				Config: loadAccTestTerraform("06_securing_security_info_registration.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".securing_security_info").Key("name").HasValue("Securing security info registration"),
					check.That("data."+dataSourceType+".securing_security_info").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_securing_security_info").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-securing-security-info-registration-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_securing_security_info").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_securing_security_info").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_07_MFAAllUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MFA for all users template")
				},
				Config: loadAccTestTerraform("07_require_mfa_for_all_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_all_users").Key("name").HasValue("Require multifactor authentication for all users"),
					check.That("data."+dataSourceType+".mfa_all_users").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_all_users").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-all-users-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_all_users").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_all_users").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_08_MFAGuestAccess(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MFA for guest access template")
				},
				Config: loadAccTestTerraform("08_require_mfa_for_guest_access.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_guest_access").Key("name").HasValue("Require multifactor authentication for guest access"),
					check.That("data."+dataSourceType+".mfa_guest_access").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_guest_access").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-guest-access-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_guest_access").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_guest_access").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_09_MFAAzureManagement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MFA for Azure management template")
				},
				Config: loadAccTestTerraform("09_require_mfa_for_azure_management.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_azure_management").Key("name").HasValue("Require multifactor authentication for Azure management"),
					check.That("data."+dataSourceType+".mfa_azure_management").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_azure_management").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-azure-management-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_azure_management").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_azure_management").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_10_MFARiskySignins(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MFA for risky sign-ins template")
				},
				Config: loadAccTestTerraform("10_require_mfa_for_risky_signins.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_risky_signins").Key("name").HasValue("Require multifactor authentication for risky sign-ins"),
					check.That("data."+dataSourceType+".mfa_risky_signins").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_risky_signins").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-risky-signins-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_risky_signins").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_risky_signins").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_11_PasswordChangeHighRisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require password change for high-risk users template")
				},
				Config: loadAccTestTerraform("11_require_password_change_high_risk_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".password_change_high_risk").Key("name").HasValue("Require password change for high-risk users"),
					check.That("data."+dataSourceType+".password_change_high_risk").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_password_change_high_risk").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-password-change-for-high-risk-users-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_password_change_high_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_password_change_high_risk").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_12_CompliantDeviceAdmins(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require compliant or hybrid device for admins template")
				},
				Config: loadAccTestTerraform("12_require_compliant_or_hybrid_device_for_admins.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".compliant_device_admins").Key("name").HasValue("Require compliant or hybrid Azure AD joined device for admins"),
					check.That("data."+dataSourceType+".compliant_device_admins").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_compliant_device_admins").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-compliant-or-hybrid-device-for-admins-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_compliant_device_admins").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_compliant_device_admins").Key("id").IsSet(),
				),
			},
		},
	})
}

// Test 13 skipped: "Block access for unknown or unsupported device platform" template has errors
// (contradictory platform conditions and incompatible grant controls)

func TestAccDatasourceConditionalAccessTemplate_14_CompliantDeviceOrMFAAllUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require compliant device or MFA for all users template")
				},
				Config: loadAccTestTerraform("14_require_compliant_device_or_mfa_all_users.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".compliant_device_or_mfa_all_users").Key("name").HasValue("Require compliant or hybrid Azure AD joined device or multifactor authentication for all users"),
					check.That("data."+dataSourceType+".compliant_device_or_mfa_all_users").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_compliant_device_or_mfa_all_users").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-compliant-device-or-mfa-for-all-users-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_compliant_device_or_mfa_all_users").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_compliant_device_or_mfa_all_users").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_15_AppEnforcedRestrictionsO365(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from use application enforced restrictions for O365 apps template")
				},
				Config: loadAccTestTerraform("15_use_application_enforced_restrictions_o365.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".app_enforced_restrictions_o365").Key("name").HasValue("Use application enforced restrictions for O365 apps"),
					check.That("data."+dataSourceType+".app_enforced_restrictions_o365").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_app_enforced_restrictions_o365").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-use-application-enforced-restrictions-for-o365-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_app_enforced_restrictions_o365").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_app_enforced_restrictions_o365").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_16_PhishingResistantMFAAdmins(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require phishing-resistant MFA for admins template")
				},
				Config: loadAccTestTerraform("16_require_phishing_resistant_mfa_for_admins.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".phishing_resistant_mfa_admins").Key("name").HasValue("Require phishing-resistant multifactor authentication for admins"),
					check.That("data."+dataSourceType+".phishing_resistant_mfa_admins").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_phishing_resistant_mfa_admins").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-phishing-resistant-mfa-for-admins-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_phishing_resistant_mfa_admins").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_phishing_resistant_mfa_admins").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_17_MFAAdminPortals(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MFA for Microsoft admin portals template")
				},
				Config: loadAccTestTerraform("17_require_mfa_for_admin_portals.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mfa_admin_portals").Key("name").HasValue("Require multifactor authentication for Microsoft admin portals"),
					check.That("data."+dataSourceType+".mfa_admin_portals").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mfa_admin_portals").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mfa-for-admin-portals-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mfa_admin_portals").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mfa_admin_portals").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_18_BlockO365InsiderRisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from block access to Office365 apps for users with insider risk template")
				},
				Config: loadAccTestTerraform("18_block_access_o365_insider_risk.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".block_o365_insider_risk").Key("name").HasValue("Block access to Office365 apps for users with insider risk"),
					check.That("data."+dataSourceType+".block_o365_insider_risk").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_block_o365_insider_risk").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-block-access-o365-insider-risk-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_block_o365_insider_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_block_o365_insider_risk").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_19_RequireMDMCompliantDevice(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from require MDM-enrolled and compliant device template")
				},
				Config: loadAccTestTerraform("19_require_mdm_enrolled_compliant_device.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".mdm_compliant_device").Key("name").HasValue("Require MDM-enrolled and compliant device to access cloud apps for all users (Preview)"),
					check.That("data."+dataSourceType+".mdm_compliant_device").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_mdm_compliant_device").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-require-mdm-enrolled-compliant-device-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_mdm_compliant_device").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_mdm_compliant_device").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_20_SecureAccountRecoveryIdentityVerification(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from secure account recovery with identity verification template")
				},
				Config: loadAccTestTerraform("20_secure_account_recovery_identity_verification.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".secure_account_recovery").Key("name").HasValue("Secure account recovery with identity verification (Preview)"),
					check.That("data."+dataSourceType+".secure_account_recovery").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_secure_account_recovery").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-secure-account-recovery-identity-verification-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_secure_account_recovery").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_secure_account_recovery").Key("id").IsSet(),
				),
			},
		},
	})
}

func TestAccDatasourceConditionalAccessTemplate_21_BlockHighRiskAgentIdentities(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			policyResourceType,
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
					testlog.StepAction(policyResourceType, "Creating conditional access policy from block high risk agent identities template")
				},
				Config: loadAccTestTerraform("21_block_high_risk_agent_identities.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("conditional access policy", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That("data."+dataSourceType+".block_high_risk_agents").Key("name").HasValue("Block high risk agent identities from accessing resources"),
					check.That("data."+dataSourceType+".block_high_risk_agents").Key("template_id").IsSet(),
					check.That(policyResourceType+".from_template_block_high_risk_agents").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-ca-policy-template-block-high-risk-agent-identities-[a-z0-9]{8}$`)),
					check.That(policyResourceType+".from_template_block_high_risk_agents").Key("state").HasValue("enabledForReportingButNotEnforced"),
					check.That(policyResourceType+".from_template_block_high_risk_agents").Key("id").IsSet(),
				),
			},
		},
	})
}

// Helper function to load test configs from acceptance directory
func loadAccTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}
