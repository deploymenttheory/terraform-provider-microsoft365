# ==============================================================================
# Test 23: Users - Rendezvous (HRW) Strategy (No Seed)
#
# Purpose: Verify Rendezvous Hashing produces balanced distribution using
# Highest Random Weight algorithm without explicit seed
#
# Use Case: Deterministic assignment where each GUID independently evaluates
# all shards and picks the one with highest hash weight
#
# Expected Behavior:
# - Approximately equal shard sizes (~probabilistic balance, not perfect ±1)
# - Always deterministic (reproducible across runs even without seed)
# - Minimal disruption when shard count changes (only ~1/n GUIDs move)
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true"
  shard_count             = 4
  strategy                = "rendezvous"
  recalculate_on_next_run = true
  seed                    = "" # Empty seed still deterministic
}

output "shard_0_count" {
  description = "Users in shard 0 (should be ~25%)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard 1 (should be ~25%)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard 2 (should be ~25%)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_3_count" {
  description = "Users in shard 3 (should be ~25%)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}

output "size_variance" {
  description = "Max difference between largest and smallest shard"
  value       = max(length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])) - min(length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]), length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"]))
}

output "total_users" {
  description = "Total users distributed (should equal sum of all shards)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}
