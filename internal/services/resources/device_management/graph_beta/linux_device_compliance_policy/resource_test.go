package graphBetaLinuxDeviceCompliancePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaLinuxDeviceCompliancePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_device_compliance_policy"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/linux_device_compliance_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *policyMocks.LinuxDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	policyMock := &policyMocks.LinuxDeviceCompliancePolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *policyMocks.LinuxDeviceCompliancePolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	policyMock := &policyMocks.LinuxDeviceCompliancePolicyMock{}
	policyMock.RegisterErrorMocks()
	return mockClient, policyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceLinuxDeviceCompliancePolicy_01_Scenario_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("001_scenario_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("name").HasValue("unit-test-linux-compliance-001"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_001", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_001",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_02_Scenario_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_002_scenario_maximal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("002_scenario_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("name").HasValue("unit-test-linux-compliance-002"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("settings_count").HasValue("8"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("device_encryption_required").HasValue("true"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("custom_compliance_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("password_policy_minimum_length").HasValue("12"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_002", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_002",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_03_Lifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_update/put_003_lifecycle_maximal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("name").HasValue("unit-test-linux-compliance-003"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "assignments.#"),
				),
			},
			{
				Config: loadUnitTestTerraform("003_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("name").HasValue("unit-test-linux-compliance-003"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("settings_count").HasValue("8"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("device_encryption_required").HasValue("true"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("custom_compliance_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("password_policy_minimum_length").HasValue("12"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_003", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_003",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_04_Lifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_002_scenario_maximal.json", "validate_update/put_004_lifecycle_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("name").HasValue("unit-test-linux-compliance-004"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("description").HasValue("Maximal Linux compliance policy"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("settings_count").HasValue("8"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("device_encryption_required").HasValue("true"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("custom_compliance_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.#").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.0.type").HasValue("ubuntu"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.0.minimum_version").HasValue("22.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.0.maximum_version").HasValue("24.04"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.1.type").HasValue("rhel"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.1.minimum_version").HasValue("8.0"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("distribution_allowed_distros.1.maximum_version").HasValue("9.9"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("password_policy_minimum_digits").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("password_policy_minimum_length").HasValue("12"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("password_policy_minimum_lowercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("password_policy_minimum_symbols").HasValue("1"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("password_policy_minimum_uppercase").HasValue("2"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.#").HasValue("2"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "assignments.#"),
				),
			},
			{
				Config: loadUnitTestTerraform("004_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("name").HasValue("unit-test-linux-compliance-004"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_004", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_004",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_05_AssignmentsMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("005_assignments_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("name").HasValue("unit-test-linux-compliance-005"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_005", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_005",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_06_AssignmentsMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("006_assignments_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("name").HasValue("unit-test-linux-compliance-006"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_006", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_006",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_07_AssignmentsLifecycle_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("name").HasValue("unit-test-linux-compliance-007"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				Config: loadUnitTestTerraform("007_assignments_lifecycle_minimal_to_maximal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("name").HasValue("unit-test-linux-compliance-007"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_007", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_007",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_08_AssignmentsLifecycle_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json", "validate_create/post_001_scenario_minimal.json"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-compliance-008"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("5"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_type": "include", "filter_id": "44444444-4444-4444-4444-444444444444"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "exclusionGroupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "allDevicesAssignmentTarget"}),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "allLicensedUsersAssignmentTarget"}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-compliance-008"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_uppercase"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.*", map[string]string{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_type": "none", "filter_id": "00000000-0000-0000-0000-000000000000"}),
				),
			},
			{
				Config: loadUnitTestTerraform("008_assignments_lifecycle_maximal_to_minimal_step_3.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("name").HasValue("unit-test-linux-compliance-008"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("platforms").HasValue("linux"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("technologies").HasValue("linuxMdm"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("created_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("last_modified_date_time").Exists(),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("device_encryption_required").HasValue("false"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008").Key("role_scope_tag_ids.0").HasValue("0"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "distribution_allowed_distros.#"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_digits"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_length"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_lowercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_symbols"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "password_policy_minimum_uppercase"),
					resource.TestCheckNoResourceAttr(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_008", "assignments.#"),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_008",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_09_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("009_error_scenario.tf"),
				ExpectError: regexp.MustCompile("Invalid Linux Device Compliance Policy data"),
			},
		},
	})
}

func TestUnitResourceLinuxDeviceCompliancePolicy_10_CustomComplianceUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	policyMock.ExpectedWrites = []string{"validate_create/post_010_custom_compliance.json", "validate_update/put_010_custom_compliance.json"}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		CheckDestroy:             func(_ *terraform.State) error { return policyMock.CheckDestroyed() },
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("010_custom_compliance_step_1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_required").HasValue("true"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_discovery_script").HasValue("55555555-5555-5555-5555-555555555555"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_rules").Exists(),
				),
			},
			{
				Config: loadUnitTestTerraform("010_custom_compliance_step_2.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_required").HasValue("true"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_discovery_script").HasValue("55555555-5555-5555-5555-555555555555"),
					check.That(graphBetaLinuxDeviceCompliancePolicy.ResourceName+".test_010").Key("custom_compliance_rules").Exists(),
				),
			},
			{
				ResourceName:            graphBetaLinuxDeviceCompliancePolicy.ResourceName + ".test_010",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}
