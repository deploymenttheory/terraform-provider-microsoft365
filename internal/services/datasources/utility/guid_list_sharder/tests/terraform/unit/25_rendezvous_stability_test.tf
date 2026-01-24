# ==============================================================================
# Test 25: Rendezvous Stability Test - Proving Minimal Disruption
#
# Purpose: Prove that Rendezvous Hashing has superior stability compared to
# position-based strategies when shard count changes
#
# Hypothesis: When increasing from 3 to 4 shards, only ~25% of GUIDs should
# move (those assigned to the new shard_3). Round-robin would cause ~75% to move!
#
# Test Design:
# - Create TWO datasources side-by-side with identical GUIDs
# - First: 3 shards (baseline)
# - Second: 4 shards (expanded)
# - Calculate: How many GUIDs stayed in their original shard?
# - Assert: Movement should be <= 30% (theoretical: ~25%)
# ==============================================================================

# Baseline: 3-shard distribution
data "microsoft365_utility_guid_list_sharder" "baseline_3_shards" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "stability-test-2024"
}

# Expanded: 4-shard distribution (same users, same seed, +1 shard)
data "microsoft365_utility_guid_list_sharder" "expanded_4_shards" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "rendezvous"
  seed          = "stability-test-2024"
}

# ==============================================================================
# Stability Calculations
# ==============================================================================

# Count how many GUIDs remained in shard_0
output "shard_0_stable_count" {
  description = "GUIDs that stayed in shard_0 (3-shard → 4-shard)"
  value       = length(setintersection(
    data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_0"],
    data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_0"]
  ))
}

# Count how many GUIDs remained in shard_1
output "shard_1_stable_count" {
  description = "GUIDs that stayed in shard_1 (3-shard → 4-shard)"
  value       = length(setintersection(
    data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_1"],
    data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_1"]
  ))
}

# Count how many GUIDs remained in shard_2
output "shard_2_stable_count" {
  description = "GUIDs that stayed in shard_2 (3-shard → 4-shard)"
  value       = length(setintersection(
    data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_2"],
    data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_2"]
  ))
}

# Total GUIDs that didn't move
output "total_stable_guids" {
  description = "Total GUIDs that stayed in same shard number"
  value = length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_0"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_0"])) + length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_1"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_1"])) + length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_2"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_2"]))
}

# Total GUIDs being distributed
output "total_guids" {
  description = "Total GUIDs in the dataset"
  value = length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_2"])
}

# Stability percentage (what we're proving!)
output "stability_percentage" {
  description = "% of GUIDs that stayed in same shard (target: >=70%, proves <30% moved)"
  value = floor(((length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_0"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_0"])) + length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_1"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_1"])) + length(setintersection(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_2"], data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_2"]))) / (length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.baseline_3_shards.shards["shard_2"]))) * 100)
}

# New shard_3 size (should be ~25% of total)
output "new_shard_3_count" {
  description = "Size of new shard_3 (should be ~25% of total)"
  value       = length(data.microsoft365_utility_guid_list_sharder.expanded_4_shards.shards["shard_3"])
}

# Verify determinism: Same config = same ID
output "baseline_id" {
  value = data.microsoft365_utility_guid_list_sharder.baseline_3_shards.id
}

output "expanded_id" {
  value = data.microsoft365_utility_guid_list_sharder.expanded_4_shards.id
}
