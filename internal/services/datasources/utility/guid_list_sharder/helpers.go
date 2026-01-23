package utilityGuidListSharder

import (
	"crypto/sha256"
	"encoding/binary"
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
