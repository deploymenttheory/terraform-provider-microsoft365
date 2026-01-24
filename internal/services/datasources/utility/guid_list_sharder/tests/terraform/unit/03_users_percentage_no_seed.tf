# ==============================================================================
# Test 03: Users - Percentage Strategy (No Seed)
#
# Purpose: Verify percentage-based distribution produces custom-sized shards
# using API order (non-deterministic between runs)
#
# Use Case: Quick phased rollout (10% → 30% → 60%) where reproducibility
# isn't needed
#
# Expected Behavior:
# - Shard sizes match specified percentages
# - Uses API order (may change between Terraform runs)
# - Last shard gets all remaining users
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  # No seed - uses API order (non-deterministic)
}

output "pilot_count" {
  description = "Users in pilot phase (~10%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_pilot_count" {
  description = "Users in broader pilot phase (~30%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_rollout_count" {
  description = "Users in full rollout phase (~60%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_users" {
  description = "Total users distributed"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "pilot_percentage" {
  description = "Actual pilot percentage"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) / (length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])) * 100
}
