package utilityGuidListSharder

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

// =============================================================================
// Sharding Strategy Implementations
// =============================================================================

// shardByHash distributes GUIDs using hash-based assignment
// Without seed: consistent distribution based purely on GUID value (same everywhere)
// With seed: distribution varies per seed value (different rollouts get different distributions)
func shardByHash(guids []string, shardCount int, seed string) [][]string {
	shards := make([][]string, shardCount)

	for _, guid := range guids {
		hashValue := computeGUIDHash(guid, seed)
		shardIndex := int(hashValue % uint64(shardCount))
		shards[shardIndex] = append(shards[shardIndex], guid)
	}

	return shards
}

// shardByRoundRobin distributes GUIDs in circular order, guaranteeing equal shard sizes
// Without seed: uses API order (non-deterministic, may change between runs)
// With seed: sorts GUIDs by hash first, then applies round-robin (deterministic, reproducible)
func shardByRoundRobin(guids []string, shardCount int, seed string) [][]string {
	shards := make([][]string, shardCount)

	// Use deterministic ordering if seed provided
	workingGuids := guids
	if seed != "" {
		workingGuids = sortGUIDsByHash(guids, seed)
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
// With seed: sorts GUIDs by hash first, then applies percentage split (deterministic, reproducible)
func shardByPercentage(guids []string, percentages []int64, seed string) [][]string {
	totalGuids := len(guids)
	shardCount := len(percentages)
	shards := make([][]string, shardCount)

	if totalGuids == 0 {
		return shards
	}

	// Use deterministic ordering if seed provided
	workingGuids := guids
	if seed != "" {
		workingGuids = sortGUIDsByHash(guids, seed)
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
// Deterministic Ordering Helper
// =============================================================================

// computeGUIDHash computes a deterministic hash of a GUID with optional seed
// This is the core function that makes all strategies reproducible when a seed is provided
func computeGUIDHash(guid string, seed string) uint64 {
	var input string
	if seed != "" {
		input = seed + ":" + guid
	} else {
		input = guid
	}
	hash := sha256.Sum256([]byte(input))
	return binary.BigEndian.Uint64(hash[:8])
}

// guidWithHash pairs a GUID with its computed hash for sorting
type guidWithHash struct {
	guid string
	hash uint64
}

// sortGUIDsByHash returns GUIDs sorted by their hash values for deterministic ordering
// Used by round-robin and percentage strategies when seed is provided
func sortGUIDsByHash(guids []string, seed string) []string {
	pairs := make([]guidWithHash, len(guids))
	for i, guid := range guids {
		pairs[i] = guidWithHash{
			guid: guid,
			hash: computeGUIDHash(guid, seed),
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].hash < pairs[j].hash
	})

	sorted := make([]string, len(pairs))
	for i, pair := range pairs {
		sorted[i] = pair.guid
	}

	return sorted
}
