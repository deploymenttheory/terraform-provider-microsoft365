package graphBetaWindowsAutopilotDevicePreparationPolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/device_management/graph_beta/windows_autopilot_device_preparation_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	pMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	pMock.RegisterMocks()
	return mockClient, pMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.WindowsAutopilotDevicePreparationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	pMock := &policyMocks.WindowsAutopilotDevicePreparationPolicyMock{}
	pMock.RegisterErrorMocks()
	return mockClient, pMock
}

// TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_01_ListAll verifies that listing returns
// only policies created from an Autopilot device preparation template. The fixture contains a third,
// unrelated settings catalog policy which must be excluded.
func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_01_ListAll(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_list_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_all").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("2"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.name").HasValue("Autopilot Device Preparation - User Driven"),
					check.That(dataSourceType+".test").Key("items.1.id").HasValue("22222222-2222-2222-2222-222222222222"),
					check.That(dataSourceType+".test").Key("items.1.name").HasValue("Autopilot Device Preparation - Automatic"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_02_ByPolicyId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_by_policy_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("policy_id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.name").HasValue("Autopilot Device Preparation - User Driven"),
					check.That(dataSourceType+".test").Key("items.0.description").HasValue("User driven Autopilot device preparation policy"),
					check.That(dataSourceType+".test").Key("items.0.setting_count").HasValue("8"),
					check.That(dataSourceType+".test").Key("items.0.is_assigned").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.0.role_scope_tag_ids.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.role_scope_tag_ids.0").HasValue("0"),
					check.That(dataSourceType+".test").Key("items.0.template_reference.template_id").HasValue("80d33118-b7b4-40d8-b15f-81be745e053f_1"),
					check.That(dataSourceType+".test").Key("items.0.template_reference.template_family").HasValue("enrollmentConfiguration"),
					check.That(dataSourceType+".test").Key("items.0.template_reference.deployment_mode").HasValue("user_driven"),
				),
			},
		},
	})
}

// TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_03_ByName verifies the exact name match
// and, critically, that an unrelated settings catalog policy sharing the same name is not returned.
func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_03_ByName(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_by_name.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("name").HasValue("Autopilot Device Preparation - User Driven"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.template_reference.deployment_mode").HasValue("user_driven"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_04_ByODataQuery(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_odata_query.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("odata_query").HasValue("isAssigned eq true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("items.0.id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(dataSourceType+".test").Key("items.0.is_assigned").HasValue("true"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_05_NotADevicePreparationPolicy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("05_not_a_device_preparation_policy.tf"),
				ExpectError: regexp.MustCompile(`(?s)is\s+not\s+a\s+Windows\s+Autopilot\s+device\s+preparation\s+policy`),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_06_NameNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("06_name_not_found.tf"),
				ExpectError: regexp.MustCompile(`(?s)No\s+Windows\s+Autopilot\s+device\s+preparation\s+policy\s+found\s+with\s+name:\s+Policy\s+That\s+Does\s+Not\s+Exist`),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_07_ConflictingSearchMethods(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("07_conflicting_search_methods.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_08_NoSearchMethod(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("08_no_search_method.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_09_WithAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_with_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("list_assignments").HasValue("true"),
					check.That(dataSourceType+".test").Key("items.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("assignments.#").HasValue("1"),
					check.That(dataSourceType+".test").Key("assignments.0.type").HasValue("groupAssignmentTarget"),
					check.That(dataSourceType+".test").Key("assignments.0.group_id").HasValue("aaaaaaaa-0000-1111-2222-333333333333"),
					check.That(dataSourceType+".test").Key("assignments.0.filter_id").HasValue("bbbbbbbb-0000-1111-2222-444444444444"),
					check.That(dataSourceType+".test").Key("assignments.0.filter_type").HasValue("include"),
				),
			},
		},
	})
}

func TestUnitDatasourceWindowsAutopilotDevicePreparationPolicy_10_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, pMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer pMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_list_all.tf"),
				ExpectError: regexp.MustCompile("Forbidden|403|insufficient|privileges"),
			},
		},
	})
}
