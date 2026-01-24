# ==============================================================================
# Test 02: Users - Round-Robin Strategy (With Seed)
#
# Purpose: Verify round-robin distribution with seed produces exactly equal
# shard sizes AND reproducible results
#
# Use Case: A/B testing, capacity planning, or when you need exact equal
# distribution that you can recreate
#
# Expected Behavior:
# - Exactly equal shard sizes (within ±1)
# - Deterministic shuffle before round-robin
# - Same seed = same distribution every time
# ==============================================================================


data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "ab-test-2024" # Makes distribution reproducible
}

output "group_a_count" {
  description = "Users in Group A (exactly 50% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "group_b_count" {
  description = "Users in Group B (exactly 50% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "is_balanced" {
  description = "Confirms equal split (difference should be 0 or 1)"
  value       = abs(length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) - length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]))
}

output "group_a_users" {
  description = "Group A user GUIDs (for verification)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}
