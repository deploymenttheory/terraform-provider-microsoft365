package graphBetaCrossTenantAccessDefaultSettings_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	crossTenantAccessDefaultSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_default_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *crossTenantAccessDefaultSettingsMocks.CrossTenantAccessDefaultSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	crossTenantAccessDefaultSettingsMock := &crossTenantAccessDefaultSettingsMocks.CrossTenantAccessDefaultSettingsMock{}
	crossTenantAccessDefaultSettingsMock.RegisterMocks()
	return mockClient, crossTenantAccessDefaultSettingsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *crossTenantAccessDefaultSettingsMocks.CrossTenantAccessDefaultSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	crossTenantAccessDefaultSettingsMock := &crossTenantAccessDefaultSettingsMocks.CrossTenantAccessDefaultSettingsMock{}
	crossTenantAccessDefaultSettingsMock.RegisterErrorMocks()
	return mockClient, crossTenantAccessDefaultSettingsMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestUnitResourceCrossTenantAccessDefaultSettings_01_Minimal tests the minimal
// configuration with only b2b_collaboration_outbound configured.
func TestUnitResourceCrossTenantAccessDefaultSettings_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),
					check.That(resourceType+".test").Key("restore_defaults_on_destroy").HasValue("true"),

					// B2B Collaboration Outbound - Users and Groups
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.targets.0.target").HasValue("AllUsers"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.targets.0.target_type").HasValue("user"),

					// B2B Collaboration Outbound - Applications
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.applications.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.applications.targets.0.target").HasValue("AllApplications"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.applications.targets.0.target_type").HasValue("application"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_02_B2BCollaborationInbound tests
// inbound B2B collaboration settings with AllUsers and AllApplications.
func TestUnitResourceCrossTenantAccessDefaultSettings_02_B2BCollaborationInbound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_b2b_collaboration_inbound.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// B2B Collaboration Inbound - Users and Groups
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.targets.0.target").HasValue("AllUsers"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.targets.0.target_type").HasValue("user"),

					// B2B Collaboration Inbound - Applications
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.applications.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.applications.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.applications.targets.0.target").HasValue("AllApplications"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.applications.targets.0.target_type").HasValue("application"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_03_B2BDirectConnect tests
// B2B direct connect inbound and outbound settings with blocked access.
func TestUnitResourceCrossTenantAccessDefaultSettings_03_B2BDirectConnect(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_03_b2b_direct_connect.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// B2B Direct Connect Inbound - Users and Groups
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.users_and_groups.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.users_and_groups.targets.0.target").HasValue("AllUsers"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.users_and_groups.targets.0.target_type").HasValue("user"),

					// B2B Direct Connect Inbound - Applications (Office365)
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target").HasValue("Office365"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target_type").HasValue("application"),

					// B2B Direct Connect Outbound - Users and Groups
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.0.target").HasValue("AllUsers"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.0.target_type").HasValue("user"),

					// B2B Direct Connect Outbound - Applications
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.applications.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.applications.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.applications.targets.0.target").HasValue("AllApplications"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.applications.targets.0.target_type").HasValue("application"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_04_InboundTrust tests
// inbound trust settings for MFA, compliant devices, and hybrid Azure AD joined devices.
func TestUnitResourceCrossTenantAccessDefaultSettings_04_InboundTrust(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_04_inbound_trust.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Inbound Trust Settings
					check.That(resourceType+".test").Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_compliant_device_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_hybrid_azure_ad_joined_device_accepted").HasValue("true"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_05_InvitationRedemption tests
// invitation redemption identity provider configuration with precedence order.
func TestUnitResourceCrossTenantAccessDefaultSettings_05_InvitationRedemption(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_05_invitation_redemption.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Invitation Redemption Configuration
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.primary_identity_provider_precedence_order.#").HasValue("3"),
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.primary_identity_provider_precedence_order.0").HasValue("azureActiveDirectory"),
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.primary_identity_provider_precedence_order.1").HasValue("externalFederation"),
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.primary_identity_provider_precedence_order.2").HasValue("socialIdentityProviders"),
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.fallback_identity_provider").HasValue("emailOneTimePasscode"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_06_TenantRestrictions tests
// tenant restrictions for users/groups and applications.
func TestUnitResourceCrossTenantAccessDefaultSettings_06_TenantRestrictions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_06_tenant_restrictions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Tenant Restrictions - Users and Groups
					check.That(resourceType+".test").Key("tenant_restrictions.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("tenant_restrictions.users_and_groups.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("tenant_restrictions.users_and_groups.targets.0.target").HasValue("AllUsers"),
					check.That(resourceType+".test").Key("tenant_restrictions.users_and_groups.targets.0.target_type").HasValue("user"),

					// Tenant Restrictions - Applications
					check.That(resourceType+".test").Key("tenant_restrictions.applications.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("tenant_restrictions.applications.targets.#").HasValue("1"),
					check.That(resourceType+".test").Key("tenant_restrictions.applications.targets.0.target").HasValue("AllApplications"),
					check.That(resourceType+".test").Key("tenant_restrictions.applications.targets.0.target_type").HasValue("application"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_07_AutomaticUserConsent tests
// automatic user consent settings (always false in default configuration).
func TestUnitResourceCrossTenantAccessDefaultSettings_07_AutomaticUserConsent(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_07_automatic_user_consent.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Automatic User Consent Settings
					check.That(resourceType+".test").Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(resourceType+".test").Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_08_Maximal tests
// a comprehensive configuration with all available settings configured.
func TestUnitResourceCrossTenantAccessDefaultSettings_08_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_08_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),
					check.That(resourceType+".test").Key("restore_defaults_on_destroy").HasValue("true"),

					// B2B Collaboration Inbound
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.applications.access_type").HasValue("allowed"),

					// B2B Collaboration Outbound
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.applications.access_type").HasValue("allowed"),

					// B2B Direct Connect Inbound
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.access_type").HasValue("blocked"),

					// B2B Direct Connect Outbound
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.applications.access_type").HasValue("blocked"),

					// Inbound Trust
					check.That(resourceType+".test").Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_compliant_device_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_hybrid_azure_ad_joined_device_accepted").HasValue("true"),

					// Invitation Redemption
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.primary_identity_provider_precedence_order.#").HasValue("3"),
					check.That(resourceType+".test").Key("invitation_redemption_identity_provider_configuration.fallback_identity_provider").HasValue("emailOneTimePasscode"),

					// Tenant Restrictions
					check.That(resourceType+".test").Key("tenant_restrictions.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("tenant_restrictions.applications.access_type").HasValue("blocked"),

					// Automatic User Consent
					check.That(resourceType+".test").Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(resourceType+".test").Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
		{
			ResourceName:      resourceType + ".test",
			ImportState:       true,
			ImportStateId:     "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true",
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"timeouts",
			},
		},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_09_ValidatorGUID tests
// that the target validator accepts valid GUIDs.
func TestUnitResourceCrossTenantAccessDefaultSettings_09_ValidatorGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_09_validator_guid.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Verify GUID targets are accepted
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.#").HasValue("2"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.0.target").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.1.target").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
				),
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessDefaultSettings_10_ValidatorOffice365 tests
// that the target validator accepts the special "Office365" value for applications.
func TestUnitResourceCrossTenantAccessDefaultSettings_10_ValidatorOffice365(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessDefaultSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessDefaultSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_10_validator_office365.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("crossTenantAccessDefaultSettings"),

					// Verify Office365 special value is accepted
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target").HasValue("Office365"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target_type").HasValue("application"),
				),
			},
		},
	})
}
