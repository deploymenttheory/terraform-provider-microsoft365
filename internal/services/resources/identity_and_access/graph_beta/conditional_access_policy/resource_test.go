package graphBetaConditionalAccessPolicy_test

import (
	"regexp"
	"testing"

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

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "display_name", "CAD001-O365: Grant macOS access for All users when Modern Auth Clients and Compliant-v1.1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "state", "enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.client_app_types.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.client_app_types.*", "mobileAppsAndDesktopClients"),

					// Conditions - Users
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.include_users.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.include_users.*", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_groups.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_groups.*", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_groups.*", "33333333-3333-3333-3333-333333333333"),

					// Conditions - Users - Exclude Guests
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#", "6"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "internalGuest"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bCollaborationGuest"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bCollaborationMember"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bDirectConnectUser"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "otherExternalUser"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "serviceProvider"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind", "all"),

					// Conditions - Applications
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.applications.include_applications.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.applications.include_applications.*", "Office365"),

					// Conditions - Platforms
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.platforms.include_platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "conditions.platforms.include_platforms.*", "macOS"),

					// Grant Controls
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "grant_controls.operator", "OR"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "grant_controls.built_in_controls.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad001_macos_compliant", "grant_controls.built_in_controls.*", "compliantDevice"),
				),
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
					testCheckExists("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "display_name", "CAD002-O365: Grant Windows access for All users when Modern Auth Clients and Compliant-v1.1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "state", "enabledForReportingButNotEnforced"),

					// Conditions - Client App Types
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.client_app_types.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.client_app_types.*", "mobileAppsAndDesktopClients"),

					// Conditions - Users
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.include_users.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.include_users.*", "All"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_groups.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_groups.*", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_groups.*", "33333333-3333-3333-3333-333333333333"),

					// Conditions - Users - Exclude Guests
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.#", "6"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "internalGuest"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bCollaborationGuest"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bCollaborationMember"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "b2bDirectConnectUser"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "otherExternalUser"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.guest_or_external_user_types.*", "serviceProvider"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.users.exclude_guests_or_external_users.external_tenants.membership_kind", "all"),

					// Conditions - Applications
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.applications.include_applications.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.applications.include_applications.*", "Office365"),

					// Conditions - Platforms
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.platforms.include_platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "conditions.platforms.include_platforms.*", "windows"),

					// Grant Controls
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "grant_controls.operator", "OR"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "grant_controls.built_in_controls.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "grant_controls.built_in_controls.*", "compliantDevice"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_identity_and_access_conditional_access_policy.cad002_windows_compliant", "grant_controls.built_in_controls.*", "domainJoinedDevice"),
				),
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
