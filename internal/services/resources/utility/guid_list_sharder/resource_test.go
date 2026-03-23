package utilityGuidListSharder_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGuidListSharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/utility/guid_list_sharder"
	guidListSharderMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/utility/guid_list_sharder/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = utilityGuidListSharder.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *guidListSharderMocks.GuidListSharderMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	guidListSharderMock := &guidListSharderMocks.GuidListSharderMock{}
	guidListSharderMock.RegisterMocks()
	return mockClient, guidListSharderMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// =============================================================================
// Users - Round-Robin Strategy Tests
// =============================================================================

// Test 01: Users - Round-Robin Strategy (No Seed)
func TestUnitResourceGuidListSharder_01_UsersRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_users_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("4"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 02: Users - Round-Robin Strategy (With Seed)
func TestUnitResourceGuidListSharder_02_UsersRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_users_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("2"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Users - Percentage Strategy Tests
// =============================================================================

// Test 03: Users - Percentage Strategy (No Seed)
func TestUnitResourceGuidListSharder_03_UsersPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_users_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 04: Users - Percentage Strategy (With Seed)
func TestUnitResourceGuidListSharder_04_UsersPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_users_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Users - Size Strategy Tests
// =============================================================================

// Test 05: Users - Size Strategy (No Seed)
func TestUnitResourceGuidListSharder_05_UsersSizeNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_users_size_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").HasValue("10"),
					check.That(resourceType+".test").Key("shards.shard_1.#").HasValue("20"),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("pilot_count", "10"),
					resource.TestCheckOutput("broader_count", "20"),
				),
			},
		},
	})
}

// Test 06: Users - Size Strategy (With Seed)
func TestUnitResourceGuidListSharder_06_UsersSizeWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_users_size_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").HasValue("5"),
					check.That(resourceType+".test").Key("shards.shard_1.#").HasValue("10"),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("pilot_count", "5"),
					resource.TestCheckOutput("broader_count", "10"),
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Round-Robin Strategy Tests
// =============================================================================

// Test 07: Group Members - Round-Robin Strategy (No Seed)
func TestUnitResourceGuidListSharder_07_GroupMembersRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_group_members_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("4"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 08: Group Members - Round-Robin Strategy (With Seed)
func TestUnitResourceGuidListSharder_08_GroupMembersRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_group_members_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("2"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Percentage Strategy Tests
// =============================================================================

// Test 09: Group Members - Percentage Strategy (No Seed)
func TestUnitResourceGuidListSharder_09_GroupMembersPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_group_members_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 10: Group Members - Percentage Strategy (With Seed)
func TestUnitResourceGuidListSharder_10_GroupMembersPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("10_group_members_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Size Strategy Tests
// =============================================================================

// Test 11: Group Members - Size Strategy (No Seed)
func TestUnitResourceGuidListSharder_11_GroupMembersSizeNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("11_group_members_size_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").HasValue("5"),
					check.That(resourceType+".test").Key("shards.shard_1.#").HasValue("15"),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("pilot_count", "5"),
					resource.TestCheckOutput("broader_count", "15"),
				),
			},
		},
	})
}

// Test 12: Group Members - Size Strategy (With Seed)
func TestUnitResourceGuidListSharder_12_GroupMembersSizeWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("12_group_members_size_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Round-Robin Strategy Tests
// =============================================================================

// Test 13: Devices - Round-Robin Strategy (No Seed)
func TestUnitResourceGuidListSharder_13_DevicesRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("13_devices_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("4"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 14: Devices - Round-Robin Strategy (With Seed)
func TestUnitResourceGuidListSharder_14_DevicesRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("14_devices_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("2"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Percentage Strategy Tests
// =============================================================================

// Test 15: Devices - Percentage Strategy (No Seed)
func TestUnitResourceGuidListSharder_15_DevicesPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("15_devices_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 16: Devices - Percentage Strategy (With Seed)
func TestUnitResourceGuidListSharder_16_DevicesPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("16_devices_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Size Strategy Tests
// =============================================================================

// Test 17: Devices - Size Strategy (No Seed)
func TestUnitResourceGuidListSharder_17_DevicesSizeNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("17_devices_size_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("3"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 18: Devices - Size Strategy (With Seed, multiple resources)
func TestUnitResourceGuidListSharder_18_DevicesSizeWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("18_devices_size_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".windows_updates").Key("id").Exists(),
					check.That(resourceType+".windows_updates").Key("shards.%").HasValue("3"),
					check.That(resourceType+".windows_updates").Key("shards.shard_0.#").HasValue("6"),
					check.That(resourceType+".windows_updates").Key("shards.shard_1.#").HasValue("12"),
					check.That(resourceType+".windows_updates").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".app_updates").Key("id").Exists(),
					check.That(resourceType+".app_updates").Key("shards.%").HasValue("3"),
					check.That(resourceType+".app_updates").Key("shards.shard_0.#").HasValue("6"),
					check.That(resourceType+".app_updates").Key("shards.shard_1.#").HasValue("12"),
					check.That(resourceType+".app_updates").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("windows_test_count", "6"),
					resource.TestCheckOutput("app_test_count", "6"),
				),
			},
		},
	})
}

// =============================================================================
// Integration Tests
// =============================================================================

// Test 19: Integration - Conditional Access Policy Assignment
func TestUnitResourceGuidListSharder_19_IntegrationConditionalAccess(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("19_integration_conditional_access.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".mfa_rollout").Key("id").Exists(),
					check.That(resourceType+".mfa_rollout").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("pilot_count", "3"),
					resource.TestCheckOutput("broader_count", "9"),
					resource.TestCheckOutput("full_count", "18"),
				),
			},
		},
	})
}

// Test 20: Integration - Group Assignment
func TestUnitResourceGuidListSharder_20_IntegrationGroupAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("20_integration_group_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".split_department").Key("id").Exists(),
					check.That(resourceType+".split_department").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("group_a_count", "7"),
					resource.TestCheckOutput("group_b_count", "7"),
					resource.TestCheckOutput("group_c_count", "6"),
				),
			},
		},
	})
}

// =============================================================================
// Special/Edge Case Tests
// =============================================================================

// Test 21: Special - Single Shard
func TestUnitResourceGuidListSharder_21_SpecialSingleShard(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("21_special_single_shard.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("1"),
					check.That(resourceType+".test").Key("shards.shard_0.#").HasValue("30"),
					resource.TestCheckOutput("user_count", "30"),
				),
			},
		},
	})
}

// Test 22: Special - No Filter
func TestUnitResourceGuidListSharder_22_SpecialNoFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("22_special_no_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("5"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_4.#").Exists(),
					resource.TestCheckOutput("total_users", "30"),
				),
			},
		},
	})
}

// =============================================================================
// Users - Rendezvous Strategy Tests
// =============================================================================

// Test 23: Users - Rendezvous Strategy (No Seed)
func TestUnitResourceGuidListSharder_23_UsersRendezvousNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("23_users_rendezvous_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("4"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 24: Users - Rendezvous Strategy (With Seed)
func TestUnitResourceGuidListSharder_24_UsersRendezvousWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("24_users_rendezvous_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("shards.%").HasValue("4"),
					check.That(resourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(resourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 25: Rendezvous Stability Test
func TestUnitResourceGuidListSharder_25_RendezvousStability(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("25_rendezvous_stability_test.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".baseline_3_shards").Key("id").Exists(),
					check.That(resourceType+".expanded_4_shards").Key("id").Exists(),
					check.That(resourceType+".baseline_3_shards").Key("shards.%").HasValue("3"),
					check.That(resourceType+".expanded_4_shards").Key("shards.%").HasValue("4"),
					resource.TestMatchOutput("stability_percentage", regexp.MustCompile(`^(7[0-9]|8[0-9]|9[0-9]|100)$`)),
					check.That(resourceType+".expanded_4_shards").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Validation Tests
// =============================================================================

func TestUnitResourceGuidListSharder_Validation_MissingGroupId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "group_members"
  shard_count             = 2
  strategy                = "round-robin"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile(`(?i)group_id.*required|attribute required`),
			},
		},
	})
}

func TestUnitResourceGuidListSharder_Validation_InvalidStrategy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "users"
  shard_count             = 2
  strategy                = "invalid"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile("(?i)invalid|unexpected|round-robin|percentage"),
			},
		},
	})
}

func TestUnitResourceGuidListSharder_Validation_InvalidResourceType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "invalid"
  shard_count             = 2
  strategy                = "round-robin"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile("(?i)invalid|unexpected|users|devices|group_members"),
			},
		},
	})
}

func TestUnitResourceGuidListSharder_Validation_BothShardingMethods(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "users"
  shard_count             = 3
  shard_percentages       = [50, 50]
  strategy                = "round-robin"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile("(?i)Invalid Attribute Combination|exactly one"),
			},
		},
	})
}

func TestUnitResourceGuidListSharder_Validation_InvalidPercentageSum(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "users"
  shard_percentages       = [50, 40]
  strategy                = "percentage"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile("(?i)sum.*100|Invalid List Sum"),
			},
		},
	})
}

func TestUnitResourceGuidListSharder_Validation_NegativeShardCount(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "users"
  shard_count             = -1
  strategy                = "round-robin"
  recalculate_on_next_run = true
}
`,
				ExpectError: regexp.MustCompile("(?i)value must be at least 1"),
			},
		},
	})
}
