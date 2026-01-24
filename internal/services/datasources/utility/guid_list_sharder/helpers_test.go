package utilityGuidListSharder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Test Data Helpers
// =============================================================================

// generateTestGUIDs creates a predictable set of GUIDs for testing
func generateTestGUIDs(count int) []string {
	guids := make([]string, count)
	for i := 0; i < count; i++ {
		// Generate predictable GUIDs: 00000000-0000-0000-0000-000000000000 to 00000000-0000-0000-0000-000000000099
		guids[i] = fmt.Sprintf("00000000-0000-0000-0000-%012d", i)
	}
	return guids
}

// Helper to count total GUIDs across all shards
func countTotalGUIDs(shards [][]string) int {
	total := 0
	for _, shard := range shards {
		total += len(shard)
	}
	return total
}

// Helper to check if a GUID exists in any shard
func containsGUID(shards [][]string, guid string) bool {
	for _, shard := range shards {
		for _, g := range shard {
			if g == guid {
				return true
			}
		}
	}
	return false
}

// Helper to find which shard contains a specific GUID
func findGUIDShard(shards [][]string, guid string) int {
	for i, shard := range shards {
		for _, g := range shard {
			if g == guid {
				return i
			}
		}
	}
	return -1 // Not found
}

// =============================================================================
// createSeededRNG Tests
// =============================================================================

func TestCreateSeededRNG_Deterministic(t *testing.T) {
	seed := "test-seed"

	rng1 := createSeededRNG(seed)
	rng2 := createSeededRNG(seed)

	// Generate some random numbers to verify they're identical
	for i := 0; i < 10; i++ {
		val1 := rng1.Intn(1000)
		val2 := rng2.Intn(1000)
		assert.Equal(t, val1, val2, "Same seed should produce identical random sequences")
	}
}

func TestCreateSeededRNG_DifferentSeeds(t *testing.T) {
	rng1 := createSeededRNG("seed1")
	rng2 := createSeededRNG("seed2")

	// Generate some random numbers - they should be different
	differentCount := 0
	for i := 0; i < 10; i++ {
		val1 := rng1.Intn(1000)
		val2 := rng2.Intn(1000)
		if val1 != val2 {
			differentCount++
		}
	}

	assert.Greater(t, differentCount, 0, "Different seeds should produce different random sequences")
}

// =============================================================================
// shuffle Tests
// =============================================================================

func TestShuffle_EmptyList(t *testing.T) {
	guids := []string{}
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(guids, rng)

	assert.Empty(t, shuffled, "Expected empty list")
}

func TestShuffle_SingleItem(t *testing.T) {
	guids := []string{"12345678-1234-1234-1234-123456789abc"}
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(guids, rng)

	require.Len(t, shuffled, 1, "Expected 1 item")
	assert.Equal(t, guids[0], shuffled[0], "Single item should remain unchanged")
}

func TestShuffle_MultipleItems(t *testing.T) {
	guids := generateTestGUIDs(10)
	rng := createSeededRNG("test-seed")
	shuffled := shuffle(guids, rng)

	assert.Len(t, shuffled, len(guids), "Shuffled list should have same length as input")

	// Verify all original items are present (no loss or duplication)
	for _, guid := range guids {
		assert.Contains(t, shuffled, guid, "Original GUID should be in shuffled list")
	}
}

func TestShuffle_DoesNotMutateOriginal(t *testing.T) {
	original := generateTestGUIDs(10)
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	rng := createSeededRNG("test-seed")
	_ = shuffle(original, rng)

	assert.Equal(t, originalCopy, original, "Original slice should not be mutated")
}

func TestShuffle_Deterministic(t *testing.T) {
	guids := generateTestGUIDs(20)

	rng1 := createSeededRNG("test-seed")
	shuffled1 := shuffle(guids, rng1)

	rng2 := createSeededRNG("test-seed")
	shuffled2 := shuffle(guids, rng2)

	require.Len(t, shuffled1, len(shuffled2), "Shuffled lists should have same length")
	assert.Equal(t, shuffled1, shuffled2, "Same RNG seed should produce identical shuffle order")
}

// =============================================================================
// shuffleWithSeed Tests (Integration)
// =============================================================================

func TestShuffleWithSeed_EmptyList(t *testing.T) {
	guids := []string{}
	shuffled := shuffleWithSeed(guids, "test-seed")

	assert.Empty(t, shuffled, "Expected empty list")
}

func TestShuffleWithSeed_SingleItem(t *testing.T) {
	guids := []string{"12345678-1234-1234-1234-123456789abc"}
	shuffled := shuffleWithSeed(guids, "test-seed")

	require.Len(t, shuffled, 1, "Expected 1 item")
	assert.Equal(t, guids[0], shuffled[0], "Single item should remain unchanged")
}

func TestShuffleWithSeed_MultipleItems(t *testing.T) {
	guids := generateTestGUIDs(10)
	shuffled := shuffleWithSeed(guids, "test-seed")

	assert.Len(t, shuffled, len(guids), "Shuffled list should have same length as input")

	// Verify all original items are present (no loss or duplication)
	for _, guid := range guids {
		assert.Contains(t, shuffled, guid, "Original GUID should be in shuffled list")
	}
}

func TestShuffleWithSeed_Deterministic(t *testing.T) {
	guids := generateTestGUIDs(20)
	seed := "test-seed"

	shuffled1 := shuffleWithSeed(guids, seed)
	shuffled2 := shuffleWithSeed(guids, seed)

	require.Len(t, shuffled1, len(shuffled2), "Shuffled lists should have same length")
	assert.Equal(t, shuffled1, shuffled2, "Same seed should produce identical shuffle order")
}

func TestShuffleWithSeed_DifferentSeeds(t *testing.T) {
	guids := generateTestGUIDs(20)

	shuffled1 := shuffleWithSeed(guids, "seed1")
	shuffled2 := shuffleWithSeed(guids, "seed2")

	assert.NotEqual(t, shuffled1, shuffled2, "Different seeds should produce different shuffle orders")
}

func TestShuffleWithSeed_DoesNotMutateOriginal(t *testing.T) {
	original := generateTestGUIDs(10)
	originalCopy := make([]string, len(original))
	copy(originalCopy, original)

	_ = shuffleWithSeed(original, "test-seed")

	assert.Equal(t, originalCopy, original, "Original slice should not be mutated")
}

// =============================================================================
// shardByRoundRobin Tests - Perfect Distribution Verification
// =============================================================================

func TestShardByRoundRobin_EmptyList(t *testing.T) {
	guids := []string{}
	shards := shardByRoundRobin(guids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardByRoundRobin_SingleShard(t *testing.T) {
	guids := generateTestGUIDs(10)
	shards := shardByRoundRobin(guids, 1, "")

	require.Len(t, shards, 1, "Expected 1 shard")
	assert.Len(t, shards[0], len(guids), "All GUIDs should be in single shard")
}

// Test perfect distribution (no variance)
func TestShardByRoundRobin_PerfectDistribution(t *testing.T) {
	// Test with exactly divisible count
	guids := generateTestGUIDs(30)
	shards := shardByRoundRobin(guids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// Each shard should have exactly 10 GUIDs
	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have exactly 10 GUIDs", i)
	}
}

// Test perfect distribution with remainder
func TestShardByRoundRobin_PerfectDistribution_WithRemainder(t *testing.T) {
	// 31 GUIDs / 3 shards = 10, 10, 11
	guids := generateTestGUIDs(31)
	shards := shardByRoundRobin(guids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	// Shards should be within ±1 of each other
	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Adjacent shards should differ by at most 1")
	}

	// Total should be 31
	total := countTotalGUIDs(shards)
	assert.Equal(t, 31, total, "Total should be 31")
}

// Test realistic perfect distribution (512 users, 3 shards)
func TestShardByRoundRobin_RealisticDistribution_512Users(t *testing.T) {
	guids := generateTestGUIDs(512)
	shards := shardByRoundRobin(guids, 3, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// 512 / 3 = 170 remainder 2, so distribution should be: 171, 171, 170
	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}

	t.Logf("Shard counts: %v", counts)

	// All shards should be within 1 of each other
	for i := 0; i < len(counts)-1; i++ {
		diff := counts[i] - counts[i+1]
		assert.LessOrEqual(t, abs(diff), 1, "Shards should differ by at most 1")
	}

	total := countTotalGUIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

// Test that WITHOUT seed, order is based on input (API order)
func TestShardByRoundRobin_NoSeed_UsesInputOrder(t *testing.T) {
	guids := generateTestGUIDs(9)
	shards := shardByRoundRobin(guids, 3, "")

	// Without seed, round-robin uses input order
	// GUID 0 → shard 0, GUID 1 → shard 1, GUID 2 → shard 2, GUID 3 → shard 0, ...

	assert.Equal(t, guids[0], shards[0][0], "First GUID should be in shard 0, position 0")
	assert.Equal(t, guids[1], shards[1][0], "Second GUID should be in shard 1, position 0")
	assert.Equal(t, guids[2], shards[2][0], "Third GUID should be in shard 2, position 0")
	assert.Equal(t, guids[3], shards[0][1], "Fourth GUID should be in shard 0, position 1")
}

// Test that WITH seed, order is shuffled first, then round-robin
func TestShardByRoundRobin_WithSeed_Deterministic(t *testing.T) {
	guids := generateTestGUIDs(100)
	seed := "test-seed"

	shards1 := shardByRoundRobin(guids, 3, seed)
	shards2 := shardByRoundRobin(guids, 3, seed)

	// Verify each shard has identical contents with same seed
	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should have identical GUIDs and order with same seed", i)
	}
}

// Test that different seeds produce different distributions
func TestShardByRoundRobin_DifferentSeeds_DifferentDistributions(t *testing.T) {
	guids := generateTestGUIDs(100)

	shardsNoSeed := shardByRoundRobin(guids, 3, "")
	shardsSeed1 := shardByRoundRobin(guids, 3, "seed1")
	shardsSeed2 := shardByRoundRobin(guids, 3, "seed2")

	// Count how many GUIDs are in different shards
	differentFromNoSeed := 0
	differentBetweenSeeds := 0

	for _, guid := range guids {
		noSeedShard := findGUIDShard(shardsNoSeed, guid)
		seed1Shard := findGUIDShard(shardsSeed1, guid)
		seed2Shard := findGUIDShard(shardsSeed2, guid)

		if noSeedShard != seed1Shard {
			differentFromNoSeed++
		}
		if seed1Shard != seed2Shard {
			differentBetweenSeeds++
		}
	}

	assert.Greater(t, differentFromNoSeed, 50, "At least 50%% of GUIDs should be in different shards (no seed vs seed)")
	assert.Greater(t, differentBetweenSeeds, 50, "At least 50%% of GUIDs should be in different shards (seed1 vs seed2)")
}

// Test round-robin with 10 shards to verify scalability
func TestShardByRoundRobin_10Shards(t *testing.T) {
	guids := generateTestGUIDs(100)
	shards := shardByRoundRobin(guids, 10, "")

	require.Len(t, shards, 10, "Expected 10 shards")

	// 100 / 10 = 10 exactly, so each shard should have exactly 10 GUIDs
	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have 10 GUIDs", i)
	}

	total := countTotalGUIDs(shards)
	assert.Equal(t, 100, total, "All 100 GUIDs should be distributed")

	// Verify each GUID appears exactly once
	for _, guid := range guids {
		assert.True(t, containsGUID(shards, guid), "GUID %s should be in a shard", guid)
	}
}

// =============================================================================
// shardByPercentage Tests - Precise Percentage Verification
// =============================================================================

func TestShardByPercentage_EmptyList(t *testing.T) {
	guids := []string{}
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(guids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

// Test that percentages are accurately applied
func TestShardByPercentage_AccuratePercentages_100Users(t *testing.T) {
	guids := generateTestGUIDs(100)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(guids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// With 100 GUIDs: 10% = 10, 30% = 30, 60% = 60
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 GUIDs (10%%)")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 GUIDs (30%%)")
	assert.Len(t, shards[2], 60, "Shard 2 should have 60 GUIDs (60%%)")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 100, total, "Total should be 100")
}

// Test realistic percentages (512 users, 10/30/60 split)
func TestShardByPercentage_RealisticPercentages_512Users(t *testing.T) {
	guids := generateTestGUIDs(512)
	percentages := []int64{10, 30, 60}
	shards := shardByPercentage(guids, percentages, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	// Expected: 10% = 51.2 → 51, 30% = 153.6 → 153, 60% = remainder (308)
	// Last shard gets all remaining

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	t.Logf("Shard counts: %v", counts)

	// Shard 0: ~10% = ~51
	assert.InDelta(t, 51, counts[0], 2, "Shard 0 should have ~51 GUIDs (10%%)")

	// Shard 1: ~30% = ~153
	assert.InDelta(t, 154, counts[1], 2, "Shard 1 should have ~154 GUIDs (30%%)")

	// Shard 2: Gets remainder
	expectedRemainder := 512 - counts[0] - counts[1]
	assert.Equal(t, expectedRemainder, counts[2], "Shard 2 should get all remaining GUIDs")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 512, total, "Total should be 512")
}

// Test that last shard gets all remaining GUIDs (no loss)
func TestShardByPercentage_LastShardGetsRemainder(t *testing.T) {
	guids := generateTestGUIDs(103) // Odd number to ensure remainder
	percentages := []int64{10, 20, 70}
	shards := shardByPercentage(guids, percentages, "")

	// Verify all GUIDs are accounted for
	total := countTotalGUIDs(shards)
	assert.Equal(t, 103, total, "All 103 GUIDs should be distributed")

	// Verify each GUID appears exactly once
	for _, guid := range guids {
		assert.True(t, containsGUID(shards, guid), "GUID should be in a shard")
	}
}

// Test that WITHOUT seed, order is based on input
func TestShardByPercentage_NoSeed_UsesInputOrder(t *testing.T) {
	guids := generateTestGUIDs(10)
	percentages := []int64{20, 30, 50}
	shards := shardByPercentage(guids, percentages, "")

	// Without seed: first 20% (2 GUIDs) go to shard 0, next 30% (3 GUIDs) to shard 1, rest to shard 2

	// Shard 0 should have first 2 GUIDs
	assert.Contains(t, shards[0], guids[0], "Shard 0 should contain first GUID")
	assert.Contains(t, shards[0], guids[1], "Shard 0 should contain second GUID")
	assert.Len(t, shards[0], 2, "Shard 0 should have 2 GUIDs")

	// Shard 1 should have next 3 GUIDs
	assert.Contains(t, shards[1], guids[2], "Shard 1 should contain third GUID")
	assert.Len(t, shards[1], 3, "Shard 1 should have 3 GUIDs")

	// Shard 2 should have remaining 5 GUIDs
	assert.Len(t, shards[2], 5, "Shard 2 should have 5 GUIDs")
}

// Test that WITH seed is deterministic
func TestShardByPercentage_WithSeed_Deterministic(t *testing.T) {
	guids := generateTestGUIDs(100)
	percentages := []int64{10, 30, 60}
	seed := "test-seed"

	shards1 := shardByPercentage(guids, percentages, seed)
	shards2 := shardByPercentage(guids, percentages, seed)

	// Verify each shard has identical contents with same seed
	for i := 0; i < 3; i++ {
		assert.Equal(t, shards1[i], shards2[i], "Shard %d should have identical GUIDs and order with same seed", i)
	}
}

// Test that different seeds produce different distributions
func TestShardByPercentage_DifferentSeeds_DifferentDistributions(t *testing.T) {
	guids := generateTestGUIDs(100)
	percentages := []int64{10, 30, 60}

	shardsNoSeed := shardByPercentage(guids, percentages, "")
	shardsSeed1 := shardByPercentage(guids, percentages, "seed1")
	shardsSeed2 := shardByPercentage(guids, percentages, "seed2")

	// Count how many GUIDs are in different shards
	differentFromNoSeed := 0
	differentBetweenSeeds := 0

	for _, guid := range guids {
		noSeedShard := findGUIDShard(shardsNoSeed, guid)
		seed1Shard := findGUIDShard(shardsSeed1, guid)
		seed2Shard := findGUIDShard(shardsSeed2, guid)

		if noSeedShard != seed1Shard {
			differentFromNoSeed++
		}
		if seed1Shard != seed2Shard {
			differentBetweenSeeds++
		}
	}

	assert.Greater(t, differentFromNoSeed, 30, "At least 30%% of GUIDs should be in different shards (no seed vs seed)")
	assert.Greater(t, differentBetweenSeeds, 30, "At least 30%% of GUIDs should be in different shards (seed1 vs seed2)")
}

// Test percentage with 10 shards to verify scalability
func TestShardByPercentage_10Shards(t *testing.T) {
	guids := generateTestGUIDs(100)
	percentages := []int64{10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	shards := shardByPercentage(guids, percentages, "")

	require.Len(t, shards, 10, "Expected 10 shards")

	// Each shard should have 10 GUIDs (10% of 100)
	for i, shard := range shards {
		assert.Len(t, shard, 10, "Shard %d should have 10 GUIDs (10%%)", i)
	}

	total := countTotalGUIDs(shards)
	assert.Equal(t, 100, total, "All 100 GUIDs should be distributed")

	// Verify each GUID appears exactly once
	for _, guid := range guids {
		assert.True(t, containsGUID(shards, guid), "GUID %s should be in a shard", guid)
	}
}

// =============================================================================
// shardBySize Tests
// =============================================================================

func TestShardBySize_EmptyList(t *testing.T) {
	guids := []string{}
	sizes := []int64{10, 20, -1}
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	for i, shard := range shards {
		assert.Empty(t, shard, "Shard %d should be empty", i)
	}
}

func TestShardBySize_ExactSizes(t *testing.T) {
	guids := generateTestGUIDs(100)
	sizes := []int64{50, 30, 20}
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 GUIDs")
	assert.Len(t, shards[1], 30, "Shard 1 should have 30 GUIDs")
	assert.Len(t, shards[2], 20, "Shard 2 should have 20 GUIDs")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 100, total, "All 100 GUIDs should be distributed")
}

func TestShardBySize_WithRemainder(t *testing.T) {
	guids := generateTestGUIDs(1000)
	sizes := []int64{50, 200, -1}
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 GUIDs")
	assert.Len(t, shards[1], 200, "Shard 1 should have 200 GUIDs")
	assert.Len(t, shards[2], 750, "Shard 2 should have 750 GUIDs (remainder)")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 1000, total, "All 1000 GUIDs should be distributed")
}

func TestShardBySize_ZeroRemainder(t *testing.T) {
	// Edge case: sizes sum to exactly total GUIDs, leaving 0 for -1 shard
	guids := generateTestGUIDs(30)
	sizes := []int64{10, 20, -1} // 10+20=30, so 0 remaining
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 10, "Shard 0 should have 10 GUIDs")
	assert.Len(t, shards[1], 20, "Shard 1 should have 20 GUIDs")

	// Critical: must be empty slice []string{}, NOT nil
	// nil becomes null in Terraform state, breaking HCL expressions like length()
	assert.NotNil(t, shards[2], "Shard 2 must not be nil (would become null in Terraform)")
	assert.Len(t, shards[2], 0, "Shard 2 should have 0 GUIDs (exact match, no remainder)")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 30, total, "All 30 GUIDs should be distributed")
}

func TestShardBySize_NotEnoughGUIDs(t *testing.T) {
	guids := generateTestGUIDs(100)
	sizes := []int64{50, 200, -1} // Request 250+ but only have 100
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 50, "Shard 0 should have 50 GUIDs")
	assert.Len(t, shards[1], 50, "Shard 1 should have 50 GUIDs (only 50 left)")

	// Must be empty slice, not nil
	assert.NotNil(t, shards[2], "Shard 2 must not be nil")
	assert.Len(t, shards[2], 0, "Shard 2 should be empty (no GUIDs left)")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 100, total, "All 100 GUIDs should be distributed")
}

func TestShardBySize_NoSeed_UsesInputOrder(t *testing.T) {
	guids := generateTestGUIDs(10)
	sizes := []int64{3, 4, 3}
	shards := shardBySize(guids, sizes, "")

	// Without seed: first 3 GUIDs go to shard 0, next 4 to shard 1, last 3 to shard 2
	assert.Contains(t, shards[0], guids[0], "Shard 0 should contain first GUID")
	assert.Contains(t, shards[0], guids[1], "Shard 0 should contain second GUID")
	assert.Contains(t, shards[0], guids[2], "Shard 0 should contain third GUID")
	assert.Len(t, shards[0], 3, "Shard 0 should have 3 GUIDs")

	assert.Contains(t, shards[1], guids[3], "Shard 1 should contain fourth GUID")
	assert.Len(t, shards[1], 4, "Shard 1 should have 4 GUIDs")

	assert.Len(t, shards[2], 3, "Shard 2 should have 3 GUIDs")
}

func TestShardBySize_WithSeed_Deterministic(t *testing.T) {
	guids := generateTestGUIDs(100)
	sizes := []int64{25, 50, -1}
	seed := "test-seed"

	shards1 := shardBySize(guids, sizes, seed)
	shards2 := shardBySize(guids, sizes, seed)

	// Verify each shard has identical contents with same seed
	for i := 0; i < 3; i++ {
		assert.ElementsMatch(t, shards1[i], shards2[i], "Shard %d should be identical with same seed", i)
	}
}

func TestShardBySize_DifferentSeeds_DifferentDistributions(t *testing.T) {
	guids := generateTestGUIDs(100)
	sizes := []int64{20, 30, -1}

	shardsNoSeed := shardBySize(guids, sizes, "")
	shardsSeed1 := shardBySize(guids, sizes, "seed1")
	shardsSeed2 := shardBySize(guids, sizes, "seed2")

	// Count how many GUIDs are in different shards
	differentFromNoSeed := 0
	differentBetweenSeeds := 0

	for _, guid := range guids {
		noSeedShard := findGUIDShard(shardsNoSeed, guid)
		seed1Shard := findGUIDShard(shardsSeed1, guid)
		seed2Shard := findGUIDShard(shardsSeed2, guid)

		if noSeedShard != seed1Shard {
			differentFromNoSeed++
		}
		if seed1Shard != seed2Shard {
			differentBetweenSeeds++
		}
	}

	assert.Greater(t, differentFromNoSeed, 30, "At least 30%% of GUIDs should be in different shards (no seed vs seed)")
	assert.Greater(t, differentBetweenSeeds, 30, "At least 30%% of GUIDs should be in different shards (seed1 vs seed2)")
}

func TestShardBySize_RealisticScenario_512Users(t *testing.T) {
	guids := generateTestGUIDs(512)
	sizes := []int64{50, 200, -1} // pilot, broader, full
	shards := shardBySize(guids, sizes, "")

	require.Len(t, shards, 3, "Expected 3 shards")

	counts := []int{len(shards[0]), len(shards[1]), len(shards[2])}
	t.Logf("Shard counts: %v", counts)

	assert.Equal(t, 50, counts[0], "Pilot shard should have exactly 50 GUIDs")
	assert.Equal(t, 200, counts[1], "Broader shard should have exactly 200 GUIDs")
	assert.Equal(t, 262, counts[2], "Full shard should have 262 GUIDs (remainder)")

	total := countTotalGUIDs(shards)
	assert.Equal(t, 512, total, "All 512 GUIDs should be distributed")
}

func TestShardBySize_AllRemainderInMiddle_ShouldFail(t *testing.T) {
	// This test documents expected behavior if -1 is NOT in last position
	// Implementation should handle gracefully or validation should prevent
	guids := generateTestGUIDs(100)
	sizes := []int64{20, -1, 30} // -1 not in last position

	shards := shardBySize(guids, sizes, "")

	// Current implementation: -1 takes all remaining, leaving 0 for last shard
	require.Len(t, shards, 3, "Expected 3 shards")
	assert.Len(t, shards[0], 20, "Shard 0 should have 20 GUIDs")
	assert.Len(t, shards[1], 80, "Shard 1 should have 80 GUIDs (all remaining)")
	assert.Len(t, shards[2], 0, "Shard 2 should be empty (no GUIDs left)")
}

// Helper function for absolute value
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
