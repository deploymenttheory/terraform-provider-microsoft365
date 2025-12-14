package graphBetaConditionalAccessPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	conditionalAccessPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/conditional_access_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *conditionalAccessPolicyMocks.ConditionalAccessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	conditionalAccessPolicyMock := &conditionalAccessPolicyMocks.ConditionalAccessPolicyMock{}
	conditionalAccessPolicyMock.RegisterMocks()
	return mockClient, conditionalAccessPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *conditionalAccessPolicyMocks.ConditionalAccessPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	conditionalAccessPolicyMock := &conditionalAccessPolicyMocks.ConditionalAccessPolicyMock{}
	conditionalAccessPolicyMock.RegisterErrorMocks()
	return mockClient, conditionalAccessPolicyMock
}

// CAD001: macOS Device Compliance
func TestConditionalAccessPolicyResource_CAD001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD001(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad001_macos_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad001_macos_compliant").Key("display_name").HasValue("CAD001-O365: Grant macOS access for All users when Modern Auth Clients and Compliant-v1.1"),
					check.That(resourceType+".cad001_macos_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad001_macos_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

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
				ResourceName:      resourceType + ".cad001_macos_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD002: Windows Device Compliance
func TestConditionalAccessPolicyResource_CAD002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD002(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad002_windows_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad002_windows_compliant").Key("display_name").HasValue("CAD002-O365: Grant Windows access for All users when Modern Auth Clients and Compliant-v1.1"),
					check.That(resourceType+".cad002_windows_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad002_windows_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

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
				ResourceName:      resourceType + ".cad002_windows_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD003: iOS and Android Device Compliance or App Protection
func TestConditionalAccessPolicyResource_CAD003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD003(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("display_name").HasValue("CAD003-O365: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.3"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Users - Exclude Guests
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),

					// Conditions - Applications
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),

					// Grant Controls
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad003_mobile_compliant_or_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				ResourceName:      resourceType + ".cad003_mobile_compliant_or_app_protection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD004: Require MFA on Non-Compliant Devices via Browser
func TestConditionalAccessPolicyResource_CAD004(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD004(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("display_name").HasValue("CAD004-O365: Grant Require MFA for All users when Browser and Non-Compliant-v1.3"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Device Filter
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Grant Controls - Authentication Strength
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad004_browser_noncompliant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				ResourceName:      resourceType + ".cad004_browser_noncompliant_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD005: Block Unsupported Device Platforms
func TestConditionalAccessPolicyResource_CAD005(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD005(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("display_name").HasValue("CAD005-O365: Block access for unsupported device platforms for All users when Modern Auth Clients-v1.1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("all"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.#").HasValue("5"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("iOS"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("macOS"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("conditions.platforms.exclude_platforms.*").ContainsTypeSetElement("linux"),

					// Grant Controls
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad005_block_unsupported_platforms").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cad005_block_unsupported_platforms",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD006: Block Downloads on Unmanaged Devices
func TestConditionalAccessPolicyResource_CAD006(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD006(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("display_name").HasValue("CAD006-O365: Session block download on unmanaged device for All users when Browser and Modern App Clients and Non-Compliant-v1.5"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Device Filter
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Session Controls
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("session_controls.application_enforced_restrictions.is_enabled").HasValue("true"),

					// Grant Controls
					check.That(resourceType+".cad006_session_block_download_unmanaged").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cad006_session_block_download_unmanaged",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD007: Sign-in Frequency for Mobile Apps on Non-Compliant Devices
func TestConditionalAccessPolicyResource_CAD007(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD007(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("display_name").HasValue("CAD007-O365: Session set Sign-in Frequency for Apps for All users when Modern Auth Clients and Non-Compliant-v1.2"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),

					// Conditions - Device Filter
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Session Controls - Sign-in Frequency
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("7"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("days"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),

					// Grant Controls
					check.That(resourceType+".cad007_mobile_signin_frequency").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cad007_mobile_signin_frequency",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD008: Sign-in Frequency for Browser on Non-Compliant Devices
func TestConditionalAccessPolicyResource_CAD008(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD008(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad008_browser_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("display_name").HasValue("CAD008-All: Session set Sign-in Frequency for All users when Browser and Non-Compliant-v1.1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Device Filter
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Session Controls - Sign-in Frequency
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("1"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("days"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cad008_browser_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),

					// Grant Controls
					check.That(resourceType+".cad008_browser_signin_frequency").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cad008_browser_signin_frequency",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD009: Disable Browser Persistence on Non-Compliant Devices
func TestConditionalAccessPolicyResource_CAD009(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD009(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad009_disable_browser_persistence").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("display_name").HasValue("CAD009-All: Session disable browser persistence for All users when Browser and Non-Compliant-v1.2"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Device Filter
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Session Controls - Persistent Browser
					check.That(resourceType+".cad009_disable_browser_persistence").Key("session_controls.persistent_browser.mode").HasValue("never"),
					check.That(resourceType+".cad009_disable_browser_persistence").Key("session_controls.persistent_browser.is_enabled").HasValue("true"),

					// Grant Controls
					check.That(resourceType+".cad009_disable_browser_persistence").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cad009_disable_browser_persistence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD010: Require MFA for Device Registration/Join
func TestConditionalAccessPolicyResource_CAD010(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD010(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad010_device_registration_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad010_device_registration_mfa").Key("display_name").HasValue("CAD010-RJD: Require MFA for device join or registration when Browser and Modern Auth Clients-v1.1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications (User Actions)
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registerdevice"),

					// Grant Controls
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad010_device_registration_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				ResourceName:      resourceType + ".cad010_device_registration_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD011: Linux Device Compliance
func TestConditionalAccessPolicyResource_CAD011(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD011(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad011_linux_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad011_linux_compliant").Key("display_name").HasValue("CAD011-O365: Grant Linux access for All users when Modern Auth Clients and Compliant-v1.0"),
					check.That(resourceType+".cad011_linux_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_users.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_users.*").ContainsTypeSetElement("GuestsOrExternalUsers"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("linux"),

					// Grant Controls
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad011_linux_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
				),
			},
			{
				ResourceName:      resourceType + ".cad011_linux_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD012: Admin Access on Compliant Devices
func TestConditionalAccessPolicyResource_CAD012(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD012(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad012_admin_compliant_access").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad012_admin_compliant_access").Key("display_name").HasValue("CAD012-All: Grant access for Admin users when Browser and Modern Auth Clients and Compliant-v1.1"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (includes 26 admin roles)
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.include_roles.#").HasValue("26"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad012_admin_compliant_access").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				ResourceName:      resourceType + ".cad012_admin_compliant_access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD013: Selected Apps - Compliant Device Requirement
func TestConditionalAccessPolicyResource_CAD013(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD013(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad013_selected_apps_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("display_name").HasValue("CAD013-Selected: Grant access for All users when Browser and Modern Auth Clients and Compliant-v1.0"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),

					// Conditions - Applications (4 specific apps)
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.#").HasValue("4"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("a4f2693f-129c-4b96-982b-2c364b8314d7"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("499b84ac-1321-427f-aa17-267ca6975798"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("996def3d-b36c-4153-8607-a6fd3c01b89f"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("797f4846-ba00-4fd7-ba43-dac1f8f63013"),

					// Conditions - Platforms
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("all"),

					// Grant Controls
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad013_selected_apps_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				ResourceName:      resourceType + ".cad013_selected_apps_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD014: Edge App Protection on Windows
func TestConditionalAccessPolicyResource_CAD014(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD014(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("display_name").HasValue("CAD014-O365: Require App Protection Policy for Edge on Windows for All users when Browser and Non-Compliant-v1.0"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Platforms
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),

					// Conditions - Device Filter
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Grant Controls
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad014_edge_app_protection_windows").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				ResourceName:      resourceType + ".cad014_edge_app_protection_windows",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD015: Compliant Device for Windows and macOS Browser Access
func TestConditionalAccessPolicyResource_CAD015(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD015(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("display_name").HasValue("CAD015-All: Grant access for All users when Browser and Modern Auth Clients and Compliant on Windows and macOS-v1.0"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Platforms
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("macOS"),

					// Grant Controls
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad015_windows_macos_browser_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				ResourceName:      resourceType + ".cad015_windows_macos_browser_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD016: Token Protection for EXO/SPO/CloudPC on Windows
func TestConditionalAccessPolicyResource_CAD016(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD016(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad016_token_protection_windows").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad016_token_protection_windows").Key("display_name").HasValue("CAD016-EXO_SPO_CloudPC: Require token protection when Modern Auth Clients on Windows-v1.2"),
					check.That(resourceType+".cad016_token_protection_windows").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Users - Exclude Guests
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),

					// Conditions - Applications (5 specific applications)
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.applications.include_applications.#").HasValue("5"),

					// Conditions - Platforms
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.platforms.include_platforms.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("windows"),

					// Grant Controls
					check.That(resourceType+".cad016_token_protection_windows").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad016_token_protection_windows").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cad016_token_protection_windows").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cad016_token_protection_windows",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD017: Selected Apps - Mobile App Protection or Compliance
func TestConditionalAccessPolicyResource_CAD017(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD017(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("display_name").HasValue("CAD017-Selected: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Users - Exclude Guests
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("serviceProvider"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind").HasValue("all"),

					// Conditions - Applications
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("None"),

					// Conditions - Platforms
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),

					// Grant Controls
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad017_selected_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				ResourceName:      resourceType + ".cad017_selected_mobile_app_protection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD018: Cloud PC - Mobile App Protection or Compliance
func TestConditionalAccessPolicyResource_CAD018(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD018(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("display_name").HasValue("CAD018-CloudPC: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.0"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications (4 specific applications)
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.applications.include_applications.#").HasValue("4"),

					// Conditions - Platforms
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("android"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("conditions.platforms.include_platforms.*").ContainsTypeSetElement("iOS"),

					// Grant Controls
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cad018_cloudpc_mobile_app_protection").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantApplication"),
				),
			},
			{
				ResourceName:      resourceType + ".cad018_cloudpc_mobile_app_protection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAD019: Intune Enrollment - MFA and Sign-in Frequency
func TestConditionalAccessPolicyResource_CAD019(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAD019(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("display_name").HasValue("CAD019-Intune: Require MFA and set sign-in frequency to every time-v1.0"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.exclude_groups.#").HasValue("2"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.users.exclude_groups.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),

					// Conditions - Applications (1 specific application - Intune)
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),

					// Grant Controls - Authentication Strength
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),

					// Session Controls - Sign-in Frequency
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cad019_intune_enrollment_mfa").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".cad019_intune_enrollment_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL001: Block Specified Locations
func TestConditionalAccessPolicyResource_CAL001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL001(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal001_block_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal001_block_locations").Key("display_name").HasValue("CAL001-All: Block specified locations for All users when Browser and Modern Auth Clients-v1.1"),
					check.That(resourceType+".cal001_block_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.users.exclude_groups.#").HasValue("2"),

					// Conditions - Applications
					check.That(resourceType+".cal001_block_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Locations (deduplicated to 1)
					check.That(resourceType+".cal001_block_locations").Key("conditions.locations.include_locations.#").HasValue("1"),

					// Grant Controls
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal001_block_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cal001_block_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL002: MFA Registration from Trusted Locations Only
func TestConditionalAccessPolicyResource_CAL002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL002(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("display_name").HasValue("CAL002-RSI: Require MFA registration from trusted locations only for All users when Browser and Modern Auth Clients-v1.4"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications (User Actions)
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registersecurityinfo"),

					// Conditions - Locations
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),

					// Grant Controls
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal002_mfa_registration_trusted_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cal002_mfa_registration_trusted_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL003: Block Service Accounts from Non-Trusted Locations
func TestConditionalAccessPolicyResource_CAL003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL003(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("display_name").HasValue("CAL003-All: Block Access for Specified Service Accounts except from Provided Trusted Locations when Browser and Modern Auth Clients-v1.1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),

					// Conditions - Applications
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Locations
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),

					// Grant Controls
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal003_block_service_accounts_untrusted").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cal003_block_service_accounts_untrusted",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL004: Block Admin Access from Non-Trusted Locations
func TestConditionalAccessPolicyResource_CAL004(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL004(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("display_name").HasValue("CAL004-All: Block access for Admins from non-trusted locations when Browser and Modern Auth Clients-v1.2"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (26 unique admin roles)
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.users.include_roles.#").HasValue("26"),

					// Conditions - Applications
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Locations
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),

					// Grant Controls
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal004_block_admin_untrusted_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cal004_block_admin_untrusted_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL005: Less-Trusted Locations Require Compliance
func TestConditionalAccessPolicyResource_CAL005(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL005(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("display_name").HasValue("CAL005-Selected: Grant access for All users on less-trusted locations when Browser and Modern Auth Clients and Compliant-v1.0"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.exclude_applications.#").HasValue("1"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.applications.exclude_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Locations (deduplicated to 1)
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("conditions.locations.include_locations.#").HasValue("1"),

					// Grant Controls
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("compliantDevice"),
					check.That(resourceType+".cal005_less_trusted_locations_compliant").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("domainJoinedDevice"),
				),
			},
			{
				ResourceName:      resourceType + ".cal005_less_trusted_locations_compliant",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAL006: Allow Access Only from Specified Locations
func TestConditionalAccessPolicyResource_CAL006(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAL006(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("display_name").HasValue("CAL006-All: Only Allow Access from specified locations for specific accounts when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.users.include_groups.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.users.include_groups.*").ContainsTypeSetElement("77777777-7777-7777-7777-777777777777"),

					// Conditions - Applications
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Locations
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					// Exclude locations deduplicated to 1
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("conditions.locations.exclude_locations.#").HasValue("1"),

					// Grant Controls
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cal006_allow_only_specified_locations").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cal006_allow_only_specified_locations",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAP001: Block Legacy Authentication
func TestConditionalAccessPolicyResource_CAP001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAP001(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cap001_block_legacy_auth").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap001_block_legacy_auth").Key("display_name").HasValue("CAP001-All: Block Legacy Authentication for All users when OtherClients-v1.0"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.client_app_types.*").ContainsTypeSetElement("other"),

					// Conditions - Users
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap001_block_legacy_auth").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cap001_block_legacy_auth",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAP002: Block Exchange ActiveSync
func TestConditionalAccessPolicyResource_CAP002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAP002(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cap002_block_exchange_activesync").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("display_name").HasValue("CAP002-All: Block Exchange ActiveSync Clients for All users-v1.1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.client_app_types.*").ContainsTypeSetElement("exchangeActiveSync"),

					// Conditions - Users
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap002_block_exchange_activesync").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cap002_block_exchange_activesync",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAP003: Block Device Code Flow
func TestConditionalAccessPolicyResource_CAP003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAP003(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cap003_block_device_code_flow").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap003_block_device_code_flow").Key("display_name").HasValue("CAP003-All: Block device code authentication flow-v1.0"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Authentication Flows
					check.That(resourceType+".cap003_block_device_code_flow").Key("conditions.authentication_flows.transfer_methods").HasValue("deviceCodeFlow"),

					// Grant Controls
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap003_block_device_code_flow").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cap003_block_device_code_flow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAP004: Block Authentication Transfer
func TestConditionalAccessPolicyResource_CAP004(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAP004(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cap004_block_auth_transfer").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cap004_block_auth_transfer").Key("display_name").HasValue("CAP004-All: Block authentication transfer-v1.0"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Authentication Flows
					check.That(resourceType+".cap004_block_auth_transfer").Key("conditions.authentication_flows.transfer_methods").HasValue("authenticationTransfer"),

					// Grant Controls
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cap004_block_auth_transfer").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cap004_block_auth_transfer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU001: Require MFA for Guest Users
func TestConditionalAccessPolicyResource_CAU001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU001(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau001_guest_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau001_guest_mfa").Key("display_name").HasValue("CAU001-All: Grant Require MFA for guests when Browser and Modern Auth Clients-v1.1"),
					check.That(resourceType+".cau001_guest_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (include_guests_or_external_users)
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),

					// Conditions - Applications
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau001_guest_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau001_guest_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				ResourceName:      resourceType + ".cau001_guest_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU001A: Require MFA for Guests - Windows Azure AD
func TestConditionalAccessPolicyResource_CAU001A(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU001A(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("display_name").HasValue("CAU001A-Windows Azure Active Directory: Grant Require MFA for guests when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (include_guests_or_external_users)
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),

					// Conditions - Applications (deduplicated to 1)
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("conditions.applications.include_applications.#").HasValue("1"),

					// Grant Controls
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau001a_guest_mfa_azure_ad").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				ResourceName:      resourceType + ".cau001a_guest_mfa_azure_ad",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU002: Require MFA for All Users
func TestConditionalAccessPolicyResource_CAU002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU002(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau002_all_users_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau002_all_users_mfa").Key("display_name").HasValue("CAU002-All: Grant Require MFA for All users when Browser and Modern Auth Clients-v1.5"),
					check.That(resourceType+".cau002_all_users_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (23 unique exclude_roles)
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.users.exclude_roles.#").HasValue("23"),

					// Conditions - Applications
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau002_all_users_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls - Authentication Strength
					check.That(resourceType+".cau002_all_users_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau002_all_users_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000002"),
				),
			},
			{
				ResourceName:      resourceType + ".cau002_all_users_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU003: Block Unapproved Apps for Guests
func TestConditionalAccessPolicyResource_CAU003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU003(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("display_name").HasValue("CAU003-Selected: Block unapproved apps for guests when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (include_guests_or_external_users)
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("6"),

					// Conditions - Applications (deduplicated to 1)
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("conditions.applications.include_applications.#").HasValue("1"),

					// Grant Controls
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau003_block_unapproved_apps_guests").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau003_block_unapproved_apps_guests",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU004: Route Through Microsoft Defender for Cloud Apps
func TestConditionalAccessPolicyResource_CAU004(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU004(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau004_mdca_route").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau004_mdca_route").Key("display_name").HasValue("CAU004-Selected: Session route through MDCA for All users when Browser on Non-Compliant-v1.2"),
					check.That(resourceType+".cau004_mdca_route").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau004_mdca_route").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users
					check.That(resourceType+".cau004_mdca_route").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cau004_mdca_route").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("Office365"),

					// Conditions - Devices (device_filter)
					check.That(resourceType+".cau004_mdca_route").Key("conditions.devices.device_filter.mode").HasValue("exclude"),
					check.That(resourceType+".cau004_mdca_route").Key("conditions.devices.device_filter.rule").HasValue("device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""),

					// Session Controls
					check.That(resourceType+".cau004_mdca_route").Key("session_controls.cloud_app_security.cloud_app_security_type").HasValue("mcasConfigured"),
					check.That(resourceType+".cau004_mdca_route").Key("session_controls.cloud_app_security.is_enabled").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".cau004_mdca_route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU006: MFA for Medium/High Sign-in Risk
func TestConditionalAccessPolicyResource_CAU006(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU006(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau006_signin_risk_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("display_name").HasValue("CAU006-All: Grant access for Medium and High Risk Sign-in for All Users when Browser and Modern Auth Clients require MFA-v1.4"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Risk Levels
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("medium"),

					// Conditions - Users
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),

					// Session Controls
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cau006_signin_risk_mfa").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".cau006_signin_risk_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU007: Password Change for Medium/High User Risk
func TestConditionalAccessPolicyResource_CAU007(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU007(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau007_user_risk_password_change").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau007_user_risk_password_change").Key("display_name").HasValue("CAU007-All: Grant access for Medium and High Risk Users for All Users when Browser and Modern Auth Clients require PWD reset-v1.3"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Risk Levels
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("medium"),

					// Conditions - Users
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls (AND operator with mfa + passwordChange)
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.operator").HasValue("AND"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.#").HasValue("2"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("passwordChange"),

					// Session Controls
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
					check.That(resourceType+".cau007_user_risk_password_change").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".cau007_user_risk_password_change",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU008: Phishing-Resistant MFA for Admins
func TestConditionalAccessPolicyResource_CAU008(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU008(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("display_name").HasValue("CAU008-All: Grant Require Phishing Resistant MFA for Admins when Browser and Modern Auth Clients-v1.4"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (26 unique include_roles)
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.users.include_roles.#").HasValue("26"),

					// Conditions - Applications
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls - Authentication Strength (phishing_resistant_mfa)
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau008_admin_phishing_resistant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000004"),
				),
			},
			{
				ResourceName:      resourceType + ".cau008_admin_phishing_resistant_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU009: Require MFA for Admin Portals
func TestConditionalAccessPolicyResource_CAU009(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU009(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau009_admin_portals_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("display_name").HasValue("CAU009-Management: Grant Require MFA for Admin Portals for All Users when Browser and Modern Auth Clients-v1.2"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications (deduplicated to 2)
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.applications.include_applications.#").HasValue("2"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("MicrosoftAdminPortals"),

					// Grant Controls
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau009_admin_portals_mfa").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),
				),
			},
			{
				ResourceName:      resourceType + ".cau009_admin_portals_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU010: Require Terms of Use
func TestConditionalAccessPolicyResource_CAU010(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU010(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau010_terms_of_use").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau010_terms_of_use").Key("display_name").HasValue("CAU010-All: Grant Require ToU for All Users when Browser and Modern Auth Clients-v1.2"),
					check.That(resourceType+".cau010_terms_of_use").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau010_terms_of_use").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls - Terms of Use
					check.That(resourceType+".cau010_terms_of_use").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau010_terms_of_use").Key("grant_controls.terms_of_use.#").HasValue("1"),
				),
			},
			{
				ResourceName:      resourceType + ".cau010_terms_of_use",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU011: Block Unlicensed Users
func TestConditionalAccessPolicyResource_CAU011(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU011(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau011_block_unlicensed").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau011_block_unlicensed").Key("display_name").HasValue("CAU011-All: Block access for All users except licensed when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau011_block_unlicensed").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.exclude_users.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.users.exclude_users.*").ContainsTypeSetElement("GuestsOrExternalUsers"),

					// Conditions - Applications
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau011_block_unlicensed").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau011_block_unlicensed",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU012: Security Info Registration with TAP
func TestConditionalAccessPolicyResource_CAU012(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU012(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau012_security_info_registration_tap").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("display_name").HasValue("CAU012-RSI: Combined Security Info Registration with TAP-v1.1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Users
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.users.include_users.*").ContainsTypeSetElement("All"),

					// Conditions - Applications (User Action)
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.applications.include_user_actions.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.applications.include_user_actions.*").ContainsTypeSetElement("urn:user:registersecurityinfo"),

					// Conditions - Locations
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.include_locations.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.include_locations.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.exclude_locations.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("conditions.locations.exclude_locations.*").ContainsTypeSetElement("AllTrusted"),

					// Grant Controls
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("mfa"),

					// Session Controls - Sign-in Frequency
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau012_security_info_registration_tap").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("everyTime"),
				),
			},
			{
				ResourceName:      resourceType + ".cau012_security_info_registration_tap",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU013: Phishing-Resistant MFA for All Users
func TestConditionalAccessPolicyResource_CAU013(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU013(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("display_name").HasValue("CAU013-All: Grant Require phishing resistant MFA for All users when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.users.include_groups.#").HasValue("1"),

					// Conditions - Applications
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls - Authentication Strength (phishing_resistant_mfa)
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau013_all_users_phishing_resistant_mfa").Key("grant_controls.authentication_strength.id").HasValue("00000000-0000-0000-0000-000000000004"),
				),
			},
			{
				ResourceName:      resourceType + ".cau013_all_users_phishing_resistant_mfa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU014: Block Managed Identity with Medium/High Sign-in Risk
func TestConditionalAccessPolicyResource_CAU014(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU014(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("display_name").HasValue("CAU014-All: Block Managed Identity when Sign in Risk is Medium or High-v1.0"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("all"),

					// Conditions - Service Principal Risk Levels
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.#").HasValue("2"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.*").ContainsTypeSetElement("high"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.service_principal_risk_levels.*").ContainsTypeSetElement("medium"),

					// Conditions - Users
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.users.include_users.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.users.include_users.*").ContainsTypeSetElement("None"),

					// Conditions - Applications
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Conditions - Client Applications
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_applications.include_service_principals.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("conditions.client_applications.include_service_principals.*").ContainsTypeSetElement("14ddb4bd-2aee-4603-86d2-467e438cda0a"),

					// Grant Controls
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau014_block_managed_identity_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau014_block_managed_identity_risk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU015: Block High Sign-in Risk
func TestConditionalAccessPolicyResource_CAU015(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU015(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau015_block_high_signin_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("display_name").HasValue("CAU015-All: Block access for High Risk Sign-in for All Users when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Sign-in Risk Levels
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.sign_in_risk_levels.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.sign_in_risk_levels.*").ContainsTypeSetElement("high"),

					// Conditions - Users
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.users.include_groups.#").HasValue("1"),

					// Conditions - Applications
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau015_block_high_signin_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau015_block_high_signin_risk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU016: Block High User Risk
func TestConditionalAccessPolicyResource_CAU016(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU016(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau016_block_high_user_risk").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau016_block_high_user_risk").Key("display_name").HasValue("CAU016-All: Block access for High Risk Users for All Users when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - User Risk Levels
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.user_risk_levels.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.user_risk_levels.*").ContainsTypeSetElement("high"),

					// Conditions - Users
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.users.include_groups.#").HasValue("1"),

					// Conditions - Applications
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Grant Controls
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau016_block_high_user_risk").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau016_block_high_user_risk",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU017: Admin Sign-in Frequency
func TestConditionalAccessPolicyResource_CAU017(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU017(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau017_admin_signin_frequency").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("display_name").HasValue("CAU017-All: Session set Sign-in Frequency for Admins when Browser-v1.0"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users (26 unique include_roles)
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.users.include_roles.#").HasValue("26"),

					// Conditions - Applications
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Session Controls - Sign-in Frequency
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.is_enabled").HasValue("true"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.authentication_type").HasValue("primaryAndSecondaryAuthentication"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.frequency_interval").HasValue("timeBased"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.value").HasValue("10"),
					check.That(resourceType+".cau017_admin_signin_frequency").Key("session_controls.sign_in_frequency.type").HasValue("hours"),

					// Grant Controls (no built-in controls for session-only policy)
					check.That(resourceType+".cau017_admin_signin_frequency").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cau017_admin_signin_frequency",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU018: Disable Browser Persistence for Admins
func TestConditionalAccessPolicyResource_CAU018(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU018(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("display_name").HasValue("CAU018-All: Session disable browser persistence for Admins when Browser-v1.0"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.client_app_types.#").HasValue("1"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),

					// Conditions - Users (25 unique include_roles)
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.users.include_roles.#").HasValue("25"),

					// Conditions - Applications
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),

					// Session Controls - Persistent Browser
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("session_controls.persistent_browser.is_enabled").HasValue("true"),
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("session_controls.persistent_browser.mode").HasValue("never"),

					// Grant Controls (no built-in controls for session-only policy)
					check.That(resourceType+".cau018_admin_disable_browser_persistence").Key("grant_controls.operator").HasValue("OR"),
				),
			},
			{
				ResourceName:      resourceType + ".cau018_admin_disable_browser_persistence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// CAU019: Allow Only Approved Apps for Guests
func TestConditionalAccessPolicyResource_CAU019(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, conditionalAccessPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer conditionalAccessPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCAU019(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("display_name").HasValue("CAU019-Selected: Only allow approved apps for guests when Browser and Modern Auth Clients-v1.0"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("state").HasValue("enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.#").HasValue("2"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("browser"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.client_app_types.*").ContainsTypeSetElement("mobileAppsAndDesktopClients"),

					// Conditions - Users (include guests)
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.#").HasValue("5"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("internalGuest"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationGuest"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bCollaborationMember"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("b2bDirectConnectUser"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.users.include_guests_or_external_users.guest_or_external_user_types.*").ContainsTypeSetElement("otherExternalUser"),

					// Conditions - Applications (All apps except 10 approved ones)
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.include_applications.#").HasValue("1"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.include_applications.*").ContainsTypeSetElement("All"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("conditions.applications.exclude_applications.#").HasValue("10"),

					// Grant Controls
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.operator").HasValue("OR"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.built_in_controls.#").HasValue("1"),
					check.That(resourceType+".cau019_allow_only_approved_apps_guests").Key("grant_controls.built_in_controls.*").ContainsTypeSetElement("block"),
				),
			},
			{
				ResourceName:      resourceType + ".cau019_allow_only_approved_apps_guests",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Configuration helper functions
func testConfigCAD001() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad001-o365.tf")
	if err != nil {
		panic("failed to load CAD001 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD002() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad002-o365.tf")
	if err != nil {
		panic("failed to load CAD002 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD003() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad003-o365.tf")
	if err != nil {
		panic("failed to load CAD003 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD004() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad004-o365.tf")
	if err != nil {
		panic("failed to load CAD004 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD005() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad005-o365.tf")
	if err != nil {
		panic("failed to load CAD005 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD006() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad006-o365.tf")
	if err != nil {
		panic("failed to load CAD006 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD007() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad007-o365.tf")
	if err != nil {
		panic("failed to load CAD007 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD008() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad008-all.tf")
	if err != nil {
		panic("failed to load CAD008 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD009() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad009-all.tf")
	if err != nil {
		panic("failed to load CAD009 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD010() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad010-rjd.tf")
	if err != nil {
		panic("failed to load CAD010 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD011() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad011-o365.tf")
	if err != nil {
		panic("failed to load CAD011 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD012() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad012-all.tf")
	if err != nil {
		panic("failed to load CAD012 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD013() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad013-selected.tf")
	if err != nil {
		panic("failed to load CAD013 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD014() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad014-o365.tf")
	if err != nil {
		panic("failed to load CAD014 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD015() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad015-all.tf")
	if err != nil {
		panic("failed to load CAD015 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD016() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad016-exo_spo_cloudpc.tf")
	if err != nil {
		panic("failed to load CAD016 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD017() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad017-selected.tf")
	if err != nil {
		panic("failed to load CAD017 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD018() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad018-cloudpc.tf")
	if err != nil {
		panic("failed to load CAD018 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAD019() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cad019-intune.tf")
	if err != nil {
		panic("failed to load CAD019 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL001() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal001-all.tf")
	if err != nil {
		panic("failed to load CAL001 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL002() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal002-rsi.tf")
	if err != nil {
		panic("failed to load CAL002 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL003() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal003-all.tf")
	if err != nil {
		panic("failed to load CAL003 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL004() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal004-all.tf")
	if err != nil {
		panic("failed to load CAL004 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL005() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal005-selected.tf")
	if err != nil {
		panic("failed to load CAL005 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAL006() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cal006-all.tf")
	if err != nil {
		panic("failed to load CAL006 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAP001() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cap001-all.tf")
	if err != nil {
		panic("failed to load CAP001 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAP002() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cap002-all.tf")
	if err != nil {
		panic("failed to load CAP002 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAP003() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cap003-all.tf")
	if err != nil {
		panic("failed to load CAP003 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAP004() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cap004-all.tf")
	if err != nil {
		panic("failed to load CAP004 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU001() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau001-all.tf")
	if err != nil {
		panic("failed to load CAU001 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU001A() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau001a-windows_azure_active_directory.tf")
	if err != nil {
		panic("failed to load CAU001A config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU002() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau002-all.tf")
	if err != nil {
		panic("failed to load CAU002 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU003() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau003-selected.tf")
	if err != nil {
		panic("failed to load CAU003 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU004() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau004-selected.tf")
	if err != nil {
		panic("failed to load CAU004 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU006() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau006-all.tf")
	if err != nil {
		panic("failed to load CAU006 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU007() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau007-all.tf")
	if err != nil {
		panic("failed to load CAU007 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU008() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau008-all.tf")
	if err != nil {
		panic("failed to load CAU008 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU009() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau009-management.tf")
	if err != nil {
		panic("failed to load CAU009 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU010() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau010-all.tf")
	if err != nil {
		panic("failed to load CAU010 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU011() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau011-all.tf")
	if err != nil {
		panic("failed to load CAU011 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU012() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau012-rsi.tf")
	if err != nil {
		panic("failed to load CAU012 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU013() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau013-all.tf")
	if err != nil {
		panic("failed to load CAU013 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU014() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau014-all.tf")
	if err != nil {
		panic("failed to load CAU014 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU015() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau015-all.tf")
	if err != nil {
		panic("failed to load CAU015 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU016() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau016-all.tf")
	if err != nil {
		panic("failed to load CAU016 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU017() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau017-all.tf")
	if err != nil {
		panic("failed to load CAU017 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU018() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau018-all.tf")
	if err != nil {
		panic("failed to load CAU018 config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigCAU019() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_cau019-selected.tf")
	if err != nil {
		panic("failed to load CAU019 config: " + err.Error())
	}
	return unitTestConfig
}
