package utilityGuidListSharder

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
)

// =============================================================================
// Sharding Strategy Implementations
// =============================================================================

// shardByRoundRobin distributes GUIDs in circular order, guaranteeing equal shard sizes
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: shuffles using Fisher-Yates first, then applies round-robin (deterministic, reproducible)
func shardByRoundRobin(guids []string, shardCount int, seed string) [][]string {

	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)

	// Use deterministic shuffle if seed provided
	workingGuids := guids
	if seed != "" {
		workingGuids = shuffleWithSeed(guids, seed)
	}

	// Distribute in round-robin fashion
	for i, guid := range workingGuids {
		shardIndex := i % shardCount
		shards[shardIndex] = append(shards[shardIndex], guid)
	}

	return shards
}

// shardByPercentage distributes GUIDs according to specified percentages
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: shuffles using Fisher-Yates first, then applies percentage split (deterministic, reproducible)
func shardByPercentage(guids []string, percentages []int64, seed string) [][]string {
	totalGuids := len(guids)
	shardCount := len(percentages)
	shards := make([][]string, shardCount)

	if totalGuids == 0 {
		return shards
	}

	// Use deterministic shuffle if seed provided
	workingGuids := guids
	if seed != "" {
		workingGuids = shuffleWithSeed(guids, seed)
	}

	// Distribute by percentages
	currentIndex := 0
	for i, percentage := range percentages {
		var shardSize int
		if i == shardCount-1 {
			// Last shard gets all remaining GUIDs
			shardSize = totalGuids - currentIndex
		} else {
			shardSize = int(float64(totalGuids) * float64(percentage) / 100.0)
		}

		if currentIndex+shardSize > totalGuids {
			shardSize = totalGuids - currentIndex
		}

		shards[i] = workingGuids[currentIndex : currentIndex+shardSize]
		currentIndex += shardSize
	}

	return shards
}

// shardBySize distributes GUIDs according to specified absolute sizes
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: shuffles using Fisher-Yates first, then applies size-based split (deterministic, reproducible)
// Supports -1 in the last position to mean "all remaining GUIDs"
func shardBySize(guids []string, sizes []int64, seed string) [][]string {
	totalGuids := len(guids)
	shardCount := len(sizes)
	shards := make([][]string, shardCount)

	if totalGuids == 0 {
		return shards
	}

	// Use deterministic shuffle if seed provided
	workingGuids := guids
	if seed != "" {
		workingGuids = shuffleWithSeed(guids, seed)
	}

	// Distribute by sizes
	currentIndex := 0
	for i, size := range sizes {
		var shardSize int
		
		if size == -1 {
			// -1 means "all remaining GUIDs"
			shardSize = totalGuids - currentIndex
		} else {
			shardSize = int(size)
			
			// If we don't have enough GUIDs left, take what's available
			if currentIndex+shardSize > totalGuids {
				shardSize = totalGuids - currentIndex
			}
		}

		// Always initialize shard, even if empty
		// Why: nil slices become null in Terraform state, breaking HCL expressions like length()
		// Empty slices []string{} become empty sets (length 0) which work correctly in HCL
		if shardSize > 0 && currentIndex < totalGuids {
			shards[i] = workingGuids[currentIndex : currentIndex+shardSize]
			currentIndex += shardSize
		} else {
			shards[i] = []string{}
		}
	}

	return shards
}

// =============================================================================
// Seeding and Shuffle Helpers
// =============================================================================

// createSeededRNG creates a deterministic random number generator from a seed string
// Uses SHA-256 to convert string seed into int64 for reproducible randomization
func createSeededRNG(seed string) *rand.Rand {
	hash := sha256.Sum256([]byte(seed))
	seedValue := int64(binary.BigEndian.Uint64(hash[:8]))
	return rand.New(rand.NewSource(seedValue))
}

// shuffle performs Fisher-Yates shuffle on a copy of the input slice using provided RNG
// Returns shuffled copy without mutating the original slice
func shuffle(guids []string, rng *rand.Rand) []string {
	// Create a copy to avoid mutating original slice
	shuffled := make([]string, len(guids))
	copy(shuffled, guids)

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

// shuffleWithSeed combines seeding and shuffling for convenience
// Used by round-robin and percentage strategies when seed is provided for reproducible randomization
func shuffleWithSeed(guids []string, seed string) []string {
	rng := createSeededRNG(seed)
	return shuffle(guids, rng)
}

// shardByRendezvous distributes GUIDs using Highest Random Weight (HRW) algorithm
// Each GUID computes a score for every shard and is assigned to the shard with the highest score
// This provides superior stability when shard counts change - only ~1/n GUIDs move when adding a shard
// Always deterministic (reproducible across runs) - seed affects which shard wins for each GUID
func shardByRendezvous(guids []string, shardCount int, seed string) [][]string {
	if shardCount <= 0 {
		shardCount = 1
	}

	shards := make([][]string, shardCount)
	
	// Initialize all shards as empty slices to prevent nil
	// Why: nil slices become null in Terraform state, breaking HCL expressions like length()
	// Empty slices []string{} become empty sets (length 0) which work correctly in HCL
	for i := 0; i < shardCount; i++ {
		shards[i] = []string{}
	}

	// For each GUID, compute weight for every shard and assign to highest
	for _, guid := range guids {
		highestWeight := uint64(0)
		selectedShard := 0

		// Evaluate this GUID against all shards
		for shardIdx := 0; shardIdx < shardCount; shardIdx++ {
			// Combine GUID + shard identifier + seed for deterministic weight
			// Format: "guid:shard_N:seed" ensures each GUID-shard pair gets unique hash
			input := fmt.Sprintf("%s:shard_%d:%s", guid, shardIdx, seed)
			hash := sha256.Sum256([]byte(input))
			
			// Use first 8 bytes of hash as weight (uint64 for large range)
			weight := binary.BigEndian.Uint64(hash[:8])

			// Track shard with highest weight for this GUID
			if weight > highestWeight {
				highestWeight = weight
				selectedShard = shardIdx
			}
		}

		shards[selectedShard] = append(shards[selectedShard], guid)
	}

	return shards
}
