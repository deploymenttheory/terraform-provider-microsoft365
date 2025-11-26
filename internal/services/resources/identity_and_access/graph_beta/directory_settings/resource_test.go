package graphBetaDirectorySettings_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDirectorySettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/directory_settings"
	directorySettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/directory_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaDirectorySettings.ResourceName
)

// setupMockEnvironment sets up the mock environment using centralized mocks
func setupMockEnvironment() (*mocks.Mocks, *directorySettingsMocks.DirectorySettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	directorySettingsMock := &directorySettingsMocks.DirectorySettingsMock{}
	directorySettingsMock.RegisterMocks()
	return mockClient, directorySettingsMock
}

// setupErrorMockEnvironment sets up the mock environment for error testing
func setupErrorMockEnvironment() (*mocks.Mocks, *directorySettingsMocks.DirectorySettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	directorySettingsMock := &directorySettingsMocks.DirectorySettingsMock{}
	directorySettingsMock.RegisterErrorMocks()
	return mockClient, directorySettingsMock
}

// TestDirectorySettingsResource_GroupUnifiedGuest tests Group.Unified.Guest template
func TestDirectorySettingsResource_GroupUnifiedGuest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupUnifiedGuest(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".group_unified_guest").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".group_unified_guest").Key("template_type").HasValue("Group.Unified.Guest"),
					check.That(resourceType+".group_unified_guest").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".group_unified_guest").Key("group_unified_guest.allow_to_add_guests").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".group_unified_guest",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigGroupUnifiedGuest() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_unified_guest_maximal.tf")
	if err != nil {
		panic("failed to load group_unified_guest config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_Application tests Application template
func TestDirectorySettingsResource_Application(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigApplication(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".application").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".application").Key("template_type").HasValue("Application"),
					check.That(resourceType+".application").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".application").Key("application.enable_access_check_for_privileged_application_updates").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".application",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigApplication() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_application_maximal.tf")
	if err != nil {
		panic("failed to load application config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_PasswordRuleSettings tests Password Rule Settings template
func TestDirectorySettingsResource_PasswordRuleSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigPasswordRuleSettings(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".password_rule_settings").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".password_rule_settings").Key("template_type").HasValue("Password Rule Settings"),
					check.That(resourceType+".password_rule_settings").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.banned_password_check_on_premises_mode").HasValue("Enforce"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.enable_banned_password_check_on_premises").HasValue("true"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.enable_banned_password_check").HasValue("true"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.lockout_duration_in_seconds").HasValue("120"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.lockout_threshold").HasValue("5"),
					check.That(resourceType+".password_rule_settings").Key("password_rule_settings.banned_password_list").HasValue("password123\tcompany123\tadmin123\twelcome123"),
				),
			},
			{
				ResourceName:      resourceType + ".password_rule_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigPasswordRuleSettings() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_password_rule_settings_maximal.tf")
	if err != nil {
		panic("failed to load password_rule_settings config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_GroupUnified tests Group.Unified template
func TestDirectorySettingsResource_GroupUnified(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigGroupUnified(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".group_unified").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".group_unified").Key("template_type").HasValue("Group.Unified"),
					check.That(resourceType+".group_unified").Key("overwrite_existing_settings").HasValue("true"),

					// Naming Policy Settings
					check.That(resourceType+".group_unified").Key("group_unified.prefix_suffix_naming_requirement").HasValue("GRP_[GroupName]_[Department]"),
					check.That(resourceType+".group_unified").Key("group_unified.custom_blocked_words_list").HasValue("CEO,President,Admin,Executive,Confidential"),
					check.That(resourceType+".group_unified").Key("group_unified.enable_ms_standard_blocked_words").HasValue("true"),

					// Group Creation Settings
					check.That(resourceType+".group_unified").Key("group_unified.enable_group_creation").HasValue("false"),
					check.That(resourceType+".group_unified").Key("group_unified.group_creation_allowed_group_id").HasValue("12345678-1234-1234-1234-123456789012"),

					// Guest Access Settings
					check.That(resourceType+".group_unified").Key("group_unified.allow_guests_to_access_groups").HasValue("false"),
					check.That(resourceType+".group_unified").Key("group_unified.allow_guests_to_be_group_owner").HasValue("false"),
					check.That(resourceType+".group_unified").Key("group_unified.allow_to_add_guests").HasValue("false"),
					check.That(resourceType+".group_unified").Key("group_unified.guest_usage_guidelines_url").HasValue("https://contoso.com/guest-guidelines"),

					// Classification Settings
					check.That(resourceType+".group_unified").Key("group_unified.classification_list").HasValue("Low,Medium,High,Confidential"),
					check.That(resourceType+".group_unified").Key("group_unified.default_classification").HasValue("Medium"),
					check.That(resourceType+".group_unified").Key("group_unified.classification_descriptions").HasValue("[{\"Value\":\"Low\",\"Description\":\"Low business impact\"},{\"Value\":\"Medium\",\"Description\":\"Medium business impact\"},{\"Value\":\"High\",\"Description\":\"High business impact\"},{\"Value\":\"Confidential\",\"Description\":\"Confidential information\"}]"),

					// Other Settings
					check.That(resourceType+".group_unified").Key("group_unified.usage_guidelines_url").HasValue("https://contoso.com/group-guidelines"),
					check.That(resourceType+".group_unified").Key("group_unified.enable_mip_labels").HasValue("true"),
					check.That(resourceType+".group_unified").Key("group_unified.new_unified_group_writeback_default").HasValue("false"),
				),
			},
			{
				ResourceName:      resourceType + ".group_unified",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigGroupUnified() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_group_unified_maximal.tf")
	if err != nil {
		panic("failed to load group_unified config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_ProhibitedNamesSettings tests Prohibited Names Settings template
func TestDirectorySettingsResource_ProhibitedNamesSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigProhibitedNamesSettings(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".prohibited_names_settings").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".prohibited_names_settings").Key("template_type").HasValue("Prohibited Names Settings"),
					check.That(resourceType+".prohibited_names_settings").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".prohibited_names_settings").Key("prohibited_names_settings.custom_blocked_sub_strings_list").HasValue("microsoft,windows,azure,office"),
					check.That(resourceType+".prohibited_names_settings").Key("prohibited_names_settings.custom_blocked_whole_words_list").HasValue("Microsoft,Windows,Azure,Office365,Outlook"),
				),
			},
			{
				ResourceName:      resourceType + ".prohibited_names_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigProhibitedNamesSettings() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_prohibited_names_settings_maximal.tf")
	if err != nil {
		panic("failed to load prohibited_names_settings config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_CustomPolicySettings tests Custom Policy Settings template
func TestDirectorySettingsResource_CustomPolicySettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigCustomPolicySettings(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".custom_policy_settings").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".custom_policy_settings").Key("template_type").HasValue("Custom Policy Settings"),
					check.That(resourceType+".custom_policy_settings").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".custom_policy_settings").Key("custom_policy_settings.custom_conditional_access_policy_url").HasValue("https://contoso.com/custom-ca-policy"),
				),
			},
			{
				ResourceName:      resourceType + ".custom_policy_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigCustomPolicySettings() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_custom_policy_settings_maximal.tf")
	if err != nil {
		panic("failed to load custom_policy_settings config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_ProhibitedNamesRestrictedSettings tests Prohibited Names Restricted Settings template
func TestDirectorySettingsResource_ProhibitedNamesRestrictedSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigProhibitedNamesRestrictedSettings(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".prohibited_names_restricted_settings").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".prohibited_names_restricted_settings").Key("template_type").HasValue("Prohibited Names Restricted Settings"),
					check.That(resourceType+".prohibited_names_restricted_settings").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".prohibited_names_restricted_settings").Key("prohibited_names_restricted_settings.custom_allowed_sub_strings_list").HasValue("contoso,fabrikam,northwind"),
					check.That(resourceType+".prohibited_names_restricted_settings").Key("prohibited_names_restricted_settings.custom_allowed_whole_words_list").HasValue("ContosoApp,FabrikamSolution,NorthwindTraders"),
					check.That(resourceType+".prohibited_names_restricted_settings").Key("prohibited_names_restricted_settings.do_not_validate_against_trademark").HasValue("true"),
				),
			},
			{
				ResourceName:      resourceType + ".prohibited_names_restricted_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigProhibitedNamesRestrictedSettings() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_prohibited_names_restricted_settings_maximal.tf")
	if err != nil {
		panic("failed to load prohibited_names_restricted_settings config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_ConsentPolicySettings tests Consent Policy Settings template
func TestDirectorySettingsResource_ConsentPolicySettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigConsentPolicySettings(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".consent_policy_settings").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".consent_policy_settings").Key("template_type").HasValue("Consent Policy Settings"),
					check.That(resourceType+".consent_policy_settings").Key("overwrite_existing_settings").HasValue("true"),
					check.That(resourceType+".consent_policy_settings").Key("consent_policy_settings.enable_group_specific_consent").HasValue("true"),
					check.That(resourceType+".consent_policy_settings").Key("consent_policy_settings.block_user_consent_for_risky_apps").HasValue("true"),
					check.That(resourceType+".consent_policy_settings").Key("consent_policy_settings.enable_admin_consent_requests").HasValue("true"),
					check.That(resourceType+".consent_policy_settings").Key("consent_policy_settings.constrain_group_specific_consent_to_members_of_group_id").HasValue("87654321-4321-4321-4321-210987654321"),
				),
			},
			{
				ResourceName:      resourceType + ".consent_policy_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConfigConsentPolicySettings() string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/resource_consent_policy_settings_maximal.tf")
	if err != nil {
		panic("failed to load consent_policy_settings config: " + err.Error())
	}
	return config
}

// TestDirectorySettingsResource_ErrorHandling tests error scenarios
func TestDirectorySettingsResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, directorySettingsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer directorySettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigGroupUnified(),
				ExpectError: regexp.MustCompile("MockError_BadRequest"),
			},
		},
	})
}
