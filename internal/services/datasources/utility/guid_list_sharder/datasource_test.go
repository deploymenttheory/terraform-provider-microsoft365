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
// Users - Hash Strategy Tests
// =============================================================================

// Test 01: Users - Hash Strategy (No Seed)
// Verifies hash-based distribution without seed produces consistent distribution
func TestUnitGuidListSharderDataSource_01_UsersHashNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_users_hash_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-3-hash$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
					resource.TestCheckOutput("total_users", "30"),
				),
			},
		},
	})
}

// Test 02: Users - Hash Strategy (With Seed)
// Verifies hash-based distribution with seed produces different distributions for different seeds
func TestUnitGuidListSharderDataSource_02_UsersHashWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_users_hash_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".mfa_rollout").Key("id").MatchesRegex(regexp.MustCompile(`^users-3-hash$`)),
					check.That(dataSourceType+".mfa_rollout").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".windows_rollout").Key("id").MatchesRegex(regexp.MustCompile(`^users-3-hash$`)),
					check.That(dataSourceType+".windows_rollout").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("distribution_note", "With different seeds, same users will be in different shards across rollouts"),
				),
			},
		},
	})
}

// =============================================================================
// Users - Round-Robin Strategy Tests
// =============================================================================

// Test 03: Users - Round-Robin Strategy (No Seed)
// Verifies round-robin distribution without seed (API order-based)
func TestUnitGuidListSharderDataSource_03_UsersRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_users_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-3-round-robin$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("10"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("10"),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").HasValue("10"),
				),
			},
		},
	})
}

// Test 04: Users - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution with seed (deterministic hash-based ordering)
func TestUnitGuidListSharderDataSource_04_UsersRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_users_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-4-round-robin$`)),
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

// =============================================================================
// Users - Percentage Strategy Tests
// =============================================================================

// Test 05: Users - Percentage Strategy (No Seed)
// Verifies percentage-based distribution without seed (API order-based)
func TestUnitGuidListSharderDataSource_05_UsersPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_users_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-\d+-percentage$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("4"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(), // ~10%
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(), // ~20%
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(), // ~30%
					check.That(dataSourceType+".test").Key("shards.shard_3.#").Exists(), // ~40%
				),
			},
		},
	})
}

// Test 06: Users - Percentage Strategy (With Seed)
// Verifies percentage-based distribution with seed (deterministic)
func TestUnitGuidListSharderDataSource_06_UsersPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_users_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-\d+-percentage$`)),
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
// Group Members - Hash Strategy Tests
// =============================================================================

// Test 07: Group Members - Hash Strategy (No Seed)
// Verifies hash-based distribution for group members without seed
func TestUnitGuidListSharderDataSource_07_GroupMembersHashNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_group_members_hash_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-2-hash$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					resource.TestCheckOutput("total_members", "20"),
				),
			},
		},
	})
}

// Test 08: Group Members - Hash Strategy (With Seed)
// Verifies hash-based distribution for group members with seed
func TestUnitGuidListSharderDataSource_08_GroupMembersHashWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_group_members_hash_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-2-hash$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Group Members - Round-Robin Strategy Tests
// =============================================================================

// Test 09: Group Members - Round-Robin Strategy (No Seed)
// Verifies round-robin distribution for group members without seed
func TestUnitGuidListSharderDataSource_09_GroupMembersRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_group_members_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-4-round-robin$`)),
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

// Test 10: Group Members - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution for group members with seed
func TestUnitGuidListSharderDataSource_10_GroupMembersRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("10_group_members_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-4-round-robin$`)),
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

// =============================================================================
// Group Members - Percentage Strategy Tests
// =============================================================================

// Test 11: Group Members - Percentage Strategy (No Seed)
// Verifies percentage-based distribution for group members without seed
func TestUnitGuidListSharderDataSource_11_GroupMembersPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("11_group_members_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-\d+-percentage$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 12: Group Members - Percentage Strategy (With Seed)
// Verifies percentage-based distribution for group members with seed
func TestUnitGuidListSharderDataSource_12_GroupMembersPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("12_group_members_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^group_members-\d+-percentage$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Hash Strategy Tests
// =============================================================================

// Test 13: Devices - Hash Strategy (No Seed)
// Verifies hash-based distribution for devices without seed
func TestUnitGuidListSharderDataSource_13_DevicesHashNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("13_devices_hash_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^devices-3-hash$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").Exists(),
				),
			},
		},
	})
}

// Test 14: Devices - Hash Strategy (With Seed)
// Verifies hash-based distribution for devices with seed
func TestUnitGuidListSharderDataSource_14_DevicesHashWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("14_devices_hash_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".windows_updates").Key("id").MatchesRegex(regexp.MustCompile(`^devices-3-hash$`)),
					check.That(dataSourceType+".windows_updates").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".app_updates").Key("id").MatchesRegex(regexp.MustCompile(`^devices-3-hash$`)),
					check.That(dataSourceType+".app_updates").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("distribution_note", "Different seeds ensure same device isn't always in early ring across all update types"),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Round-Robin Strategy Tests
// =============================================================================

// Test 15: Devices - Round-Robin Strategy (No Seed)
// Verifies round-robin distribution for devices without seed
func TestUnitGuidListSharderDataSource_15_DevicesRoundRobinNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("15_devices_round_robin_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^devices-4-round-robin$`)),
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

// Test 16: Devices - Round-Robin Strategy (With Seed)
// Verifies round-robin distribution for devices with seed
func TestUnitGuidListSharderDataSource_16_DevicesRoundRobinWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("16_devices_round_robin_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^devices-3-round-robin$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("3"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("8"),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").HasValue("8"),
					check.That(dataSourceType+".test").Key("shards.shard_2.#").HasValue("8"),
				),
			},
		},
	})
}

// =============================================================================
// Devices - Percentage Strategy Tests
// =============================================================================

// Test 17: Devices - Percentage Strategy (No Seed)
// Verifies percentage-based distribution for devices without seed
func TestUnitGuidListSharderDataSource_17_DevicesPercentageNoSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("17_devices_percentage_no_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^devices-\d+-percentage$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// Test 18: Devices - Percentage Strategy (With Seed)
// Verifies percentage-based distribution for devices with seed
func TestUnitGuidListSharderDataSource_18_DevicesPercentageWithSeed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, guidListSharderMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer guidListSharderMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("18_devices_percentage_with_seed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^devices-\d+-percentage$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("2"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").Exists(),
					check.That(dataSourceType+".test").Key("shards.shard_1.#").Exists(),
				),
			},
		},
	})
}

// =============================================================================
// Integration Tests
// =============================================================================

// Test 19: Integration - Conditional Access Policy Assignment
// Verifies sharded GUIDs can be used in conditional access policy group assignments
func TestUnitGuidListSharderDataSource_19_IntegrationConditionalAccess(t *testing.T) {
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
					check.That(dataSourceType+".phased_users").Key("id").Exists(),
					check.That(dataSourceType+".phased_users").Key("shards.%").HasValue("3"),
					resource.TestCheckOutput("phase_1_users_count", "10"),
					resource.TestCheckOutput("phase_2_users_count", "10"),
					resource.TestCheckOutput("phase_3_users_count", "10"),
				),
			},
		},
	})
}

// Test 20: Integration - Group Assignment
// Verifies sharded GUIDs can be used for group membership assignments
func TestUnitGuidListSharderDataSource_20_IntegrationGroupAssignment(t *testing.T) {
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
					check.That(dataSourceType+".pilot_users").Key("id").Exists(),
					check.That(dataSourceType+".pilot_users").Key("shards.%").HasValue("2"),
					resource.TestCheckOutput("pilot_group_size", "6"),
					resource.TestCheckOutput("production_group_size", "24"),
				),
			},
		},
	})
}

// =============================================================================
// Special/Edge Case Tests
// =============================================================================

// Test 21: Special - Single Shard
// Verifies behavior when creating a single shard (all GUIDs in one shard)
func TestUnitGuidListSharderDataSource_21_SpecialSingleShard(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-1-hash$`)),
					check.That(dataSourceType+".test").Key("shards.%").HasValue("1"),
					check.That(dataSourceType+".test").Key("shards.shard_0.#").HasValue("30"),
					resource.TestCheckOutput("user_count", "30"),
				),
			},
		},
	})
}

// Test 22: Special - No Filter
// Verifies behavior when no OData query filter is applied (returns all resources)
func TestUnitGuidListSharderDataSource_22_SpecialNoFilter(t *testing.T) {
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
					check.That(dataSourceType+".test").Key("id").MatchesRegex(regexp.MustCompile(`^users-5-hash$`)),
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
  strategy      = "hash"
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
				ExpectError: regexp.MustCompile("(?i)invalid|unexpected|hash|round-robin|percentage"),
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
  strategy      = "hash"
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
  strategy          = "hash"
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
  strategy      = "hash"
}
`,
				ExpectError: regexp.MustCompile("(?i)value must be at least 1"),
			},
		},
	})
}
