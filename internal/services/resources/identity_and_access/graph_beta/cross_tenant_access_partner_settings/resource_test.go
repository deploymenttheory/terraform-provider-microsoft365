package graphBetaCrossTenantAccessPartnerSettings_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	crossTenantAccessPartnerSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/cross_tenant_access_partner_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const resourceType = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings"

func setupMockEnvironment() (*mocks.Mocks, *crossTenantAccessPartnerSettingsMocks.CrossTenantAccessPartnerSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	crossTenantAccessPartnerSettingsMock := &crossTenantAccessPartnerSettingsMocks.CrossTenantAccessPartnerSettingsMock{}
	crossTenantAccessPartnerSettingsMock.RegisterMocks()
	return mockClient, crossTenantAccessPartnerSettingsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *crossTenantAccessPartnerSettingsMocks.CrossTenantAccessPartnerSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	crossTenantAccessPartnerSettingsMock := &crossTenantAccessPartnerSettingsMocks.CrossTenantAccessPartnerSettingsMock{}
	crossTenantAccessPartnerSettingsMock.RegisterErrorMocks()
	return mockClient, crossTenantAccessPartnerSettingsMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// TestUnitResourceCrossTenantAccessPartnerSettings_01_Minimal tests the minimal
// configuration with only b2b_collaboration_outbound configured.
func TestUnitResourceCrossTenantAccessPartnerSettings_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),

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
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_02_B2BCollaborationInbound tests
// inbound B2B collaboration settings with AllUsers and AllApplications.
func TestUnitResourceCrossTenantAccessPartnerSettings_02_B2BCollaborationInbound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_b2b_collaboration_inbound.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

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
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_03_B2BDirectConnect tests
// B2B direct connect inbound and outbound settings with blocked access.
func TestUnitResourceCrossTenantAccessPartnerSettings_03_B2BDirectConnect(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_03_b2b_direct_connect.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

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
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_04_InboundTrust tests
// inbound trust settings for MFA, compliant devices, and hybrid Azure AD joined devices.
func TestUnitResourceCrossTenantAccessPartnerSettings_04_InboundTrust(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_04_inbound_trust.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Inbound Trust Settings
					check.That(resourceType+".test").Key("inbound_trust.is_mfa_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_compliant_device_accepted").HasValue("true"),
					check.That(resourceType+".test").Key("inbound_trust.is_hybrid_azure_ad_joined_device_accepted").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_05_TenantRestrictions tests
// tenant restrictions for users/groups and applications.
func TestUnitResourceCrossTenantAccessPartnerSettings_05_TenantRestrictions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_05_tenant_restrictions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

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
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_06_AutomaticUserConsent tests
// automatic user consent settings.
func TestUnitResourceCrossTenantAccessPartnerSettings_06_AutomaticUserConsent(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_06_automatic_user_consent.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Automatic User Consent Settings
					check.That(resourceType+".test").Key("automatic_user_consent_settings.inbound_allowed").HasValue("false"),
					check.That(resourceType+".test").Key("automatic_user_consent_settings.outbound_allowed").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_07_Maximal tests
// a comprehensive configuration with all available settings configured.
func TestUnitResourceCrossTenantAccessPartnerSettings_07_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_07_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("is_service_provider").HasValue("true"),

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
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_08_ValidatorGUID tests
// that the target validator accepts valid GUIDs.
func TestUnitResourceCrossTenantAccessPartnerSettings_08_ValidatorGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_08_validator_guid.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Verify GUID targets are accepted
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.#").HasValue("2"),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.0.target").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
					check.That(resourceType+".test").Key("b2b_direct_connect_outbound.users_and_groups.targets.1.target").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)),
				),
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_09_ValidatorOffice365 tests
// that the target validator accepts the special "Office365" value for applications.
func TestUnitResourceCrossTenantAccessPartnerSettings_09_ValidatorOffice365(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_09_validator_office365.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Verify Office365 special value is accepted
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target").HasValue("Office365"),
					check.That(resourceType+".test").Key("b2b_direct_connect_inbound.applications.targets.0.target_type").HasValue("application"),
				),
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_10_Update tests
// updating a partner configuration from minimal to more complex settings.
func TestUnitResourceCrossTenantAccessPartnerSettings_10_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("allowed"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_10_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Verify updated settings
					check.That(resourceType+".test").Key("b2b_collaboration_inbound.users_and_groups.access_type").HasValue("blocked"),
					check.That(resourceType+".test").Key("b2b_collaboration_outbound.users_and_groups.access_type").HasValue("blocked"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "12345678-1234-1234-1234-123456789012",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_11_HardDelete tests
// the hard_delete functionality for permanent deletion.
func TestUnitResourceCrossTenantAccessPartnerSettings_11_HardDelete(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_11_hard_delete.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("tenant_id").HasValue("12345678-1234-1234-1234-123456789012"),
					check.That(resourceType+".test").Key("hard_delete").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     "12345678-1234-1234-1234-123456789012:hard_delete=true",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// TestUnitResourceCrossTenantAccessPartnerSettings_Error tests
// error handling for API failures.
func TestUnitResourceCrossTenantAccessPartnerSettings_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, crossTenantAccessPartnerSettingsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer crossTenantAccessPartnerSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_01_minimal.tf"),
				ExpectError: regexp.MustCompile(`Internal Server Error|Mock error for testing`),
			},
		},
	})
}
