package utilityGuidListSharder_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityGuidListSharder "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/guid_list_sharder"
	guidListSharderMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/guid_list_sharder/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	dataSourceType = "data." + utilityGuidListSharder.DataSourceName
)

func setupMockEnvironment() (*mocks.Mocks, *guidListSharderMocks.GuidListSharderMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	guidListSharderMock := &guidListSharderMocks.GuidListSharderMock{}
	guidListSharderMock.RegisterMocks()
	return mockClient, guidListSharderMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *guidListSharderMocks.GuidListSharderMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	guidListSharderMock := &guidListSharderMocks.GuidListSharderMock{}
	guidListSharderMock.RegisterErrorMocks()
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
// Verifies round-robin distribution without seed (API order-based)
func TestUnitGuidListSharderDataSource_01_UsersRoundRobinNoSeed(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("4"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_3.#").Exists(),
				),
			},
		},
	})
}

// Test 02: Users - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution with seed (deterministic shuffling)
func TestUnitGuidListSharderDataSource_02_UsersRoundRobinWithSeed(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Users - Percentage Strategy Tests
// =============================================================================

// Test 03: Users - Percentage Strategy (No Seed)
// Verifies percentage-based distribution without seed (API order-based)
func TestUnitGuidListSharderDataSource_03_UsersPercentageNoSeed(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(), // ~10%
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(), // ~30%
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(), // ~60%
				),
			},
		},
	})
}

// Test 04: Users - Percentage Strategy (With Seed)
// Verifies percentage-based distribution with seed (deterministic)
func TestUnitGuidListSharderDataSource_04_UsersPercentageWithSeed(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(), // ~25%
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(), // ~50%
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(), // ~25%
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Round-Robin Strategy Tests
// =============================================================================

// Test 05: Group Members - Round-Robin Strategy (No Seed)
// Verifies round-robin distribution for group members without seed
func TestUnitGuidListSharderDataSource_05_GroupMembersRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_group_members_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("4"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("5"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("5"),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").HasValue("5"),
					check.That(dataSourceType+".test").Key("shards.shard_3.#").HasValue("5"),
				),
			},
		},
	})
}

// Test 06: Group Members - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution for group members with seed
func TestUnitGuidListSharderDataSource_06_GroupMembersRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_group_members_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("10"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("10"),
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Percentage Strategy Tests
// =============================================================================

// Test 07: Group Members - Percentage Strategy (No Seed)
// Verifies percentage-based distribution for group members without seed
func TestUnitGuidListSharderDataSource_07_GroupMembersPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_group_members_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 08: Group Members - Percentage Strategy (With Seed)
// Verifies percentage-based distribution for group members with seed
func TestUnitGuidListSharderDataSource_08_GroupMembersPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_group_members_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Round-Robin Strategy Tests
// =============================================================================

// Test 09: Devices - Round-Robin Strategy (No Seed)
// Verifies round-robin distribution for devices without seed
func TestUnitGuidListSharderDataSource_09_DevicesRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_devices_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("4"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("6"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("6"),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").HasValue("6"),
					check.That(dataSourceType+".test").Key("shards.shard_3.#").HasValue("6"),
				),
			},
		},
	})
}

// Test 10: Devices - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution for devices with seed
func TestUnitGuidListSharderDataSource_10_DevicesRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("10_devices_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("12"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("12"),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Percentage Strategy Tests
// =============================================================================

// Test 11: Devices - Percentage Strategy (No Seed)
// Verifies percentage-based distribution for devices without seed
func TestUnitGuidListSharderDataSource_11_DevicesPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("11_devices_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 12: Devices - Percentage Strategy (With Seed)
// Verifies percentage-based distribution for devices with seed
func TestUnitGuidListSharderDataSource_12_DevicesPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("12_devices_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Integration Tests
// =============================================================================

// Test 13: Integration - Conditional Access Policy Assignment
// Verifies sharded GUIDs can be used in conditional access policy group assignments
func TestUnitGuidListSharderDataSource_13_IntegrationConditionalAccess(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("13_integration_conditional_access.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".mfa_rollout").Key("id").Exists(),
					check.That(dataSourceType+".mfa_rollout").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("pilot_count", "3"),
					resource.TestCheckOutput("broader_count", "9"),
					resource.TestCheckOutput("full_count", "18"),
				),
			},
		},
	})
}

// Test 14: Integration - Group Assignment
// Verifies sharded GUIDs can be used for group membership assignments
func TestUnitGuidListSharderDataSource_14_IntegrationGroupAssignment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("14_integration_group_assignment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".split_department").Key("id").Exists(),
					check.That(dataSourceType+".split_department").Key("shards.%").HasValue("3"),
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

// Test 15: Special - Single Shard
// Verifies behavior when creating a single shard (all GUIDs in one shard)
func TestUnitGuidListSharderDataSource_15_SpecialSingleShard(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("15_special_single_shard.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("1"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("30"),
					resource.TestCheckOutput("user_count", "30"),
				),
			},
		},
	})
}

// Test 16: Special - No Filter
// Verifies behavior when no OData query filter is applied (returns all resources)
func TestUnitGuidListSharderDataSource_16_SpecialNoFilter(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("16_special_no_filter.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").Exists(),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("5"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_3.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_4.#").Exists(),
					resource.TestCheckOutput("total_users", "30"),
				),
			},
		},
	})
}

// =============================================================================
// Validation/Error Tests
// =============================================================================

// Test: Validation - Missing group_id when resource_type is group_members
// Verifies schema validation rejects configuration when group_id is missing for group_members
func TestUnitGuidListSharderDataSource_Validation_MissingGroupId(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  shard_count   = 2
  strategy      = "round-robin"
}
`,
				ExpectError: regexp.MustCompile(`(?i)group_id.*required|attribute required`),
			},
		},
	})
}

// Test: Validation - Invalid strategy value
// Verifies schema validation rejects invalid strategy values
func TestUnitGuidListSharderDataSource_Validation_InvalidStrategy(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  shard_count   = 2
  strategy      = "invalid"
}
`,
				ExpectError: regexp.MustCompile("(?i)invalid|unexpected|round-robin|percentage"),
			},
		},
	})
}

// Test: Validation - Invalid resource_type value
// Verifies schema validation rejects invalid resource_type values
func TestUnitGuidListSharderDataSource_Validation_InvalidResourceType(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "invalid"
  shard_count   = 2
  strategy      = "round-robin"
}
`,
				ExpectError: regexp.MustCompile("(?i)invalid|unexpected|users|devices|group_members"),
			},
		},
	})
}

// Test: Validation - Both shard_count and shard_percentages provided
// Verifies schema validation rejects configuration with both shard_count and shard_percentages
func TestUnitGuidListSharderDataSource_Validation_BothShardingMethods(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "users"
  shard_count       = 3
  shard_percentages = [50, 50]
  strategy          = "round-robin"
}
`,
				ExpectError: regexp.MustCompile("(?i)Invalid Attribute Combination|exactly one"),
			},
		},
	})
}

// Test: Validation - Percentages don't sum to 100
// Verifies schema validation rejects shard_percentages that don't sum to 100
func TestUnitGuidListSharderDataSource_Validation_InvalidPercentageSum(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "users"
  shard_percentages = [50, 40]
  strategy          = "percentage"
}
`,
				ExpectError: regexp.MustCompile("(?i)sum.*100|Invalid List Sum"),
			},
		},
	})
}

// Test: Validation - Negative shard_count
// Verifies schema validation rejects negative shard_count
func TestUnitGuidListSharderDataSource_Validation_NegativeShardCount(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  shard_count   = -1
  strategy      = "round-robin"
}
`,
				ExpectError: regexp.MustCompile("(?i)value must be at least 1"),
			},
		},
	})
}
