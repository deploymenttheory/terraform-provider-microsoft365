# Test 01: Users - Hash Strategy (No Seed)
#
# Purpose: Verify hash-based distribution without seed produces consistent
# distribution across all instances (same GUID always goes to same shard)
#
# Use Case: Creating standard user tiers that should be identical across
# all policies and all Terraform runs
#
# Expected Behavior:
# - Approximately equal shard sizes
# - Same distribution in all instances with same shard_count
# - Deterministic and reproducible

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  # No seed - ensures identical distribution everywhere
}

output "shard_0_count" {
  description = "Users in shard 0"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard 1"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard 2"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_users" {
  description = "Total users distributed"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
