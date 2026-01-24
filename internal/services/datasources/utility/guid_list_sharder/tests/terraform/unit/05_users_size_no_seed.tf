# ==============================================================================
# Test 05: Users - Size Strategy (No Seed)
#
# Purpose: Verify absolute size-based distribution using API order
#
# Use Case: "Support team can handle exactly 50 pilot users, then 200 broader"
#
# Expected Behavior:
# - Exact shard sizes (50, 200, remainder)
# - Uses API order (may change between Terraform runs)
# - Last shard gets all remaining GUIDs
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_sizes   = [10, 20, -1]  # 10 pilot, 20 broader, rest for full
  strategy      = "size"
  # No seed - uses API order (non-deterministic)
}

output "pilot_count" {
  description = "Users in pilot (should be exactly 10)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_count" {
  description = "Users in broader pilot (should be exactly 20)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_count" {
  description = "Users in full rollout (all remaining)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_users" {
  description = "Total users across all shards"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
