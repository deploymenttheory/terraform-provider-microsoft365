# ==============================================================================
# Test 22: Special Case - No OData Filter (All Resources)
#
# Purpose: Verify behavior when no odata_filter is provided (gets ALL resources)
#
# Use Case: Splitting entire user/device/member population without filtering
#
# Expected Behavior:
# - Retrieves all resources of specified type
# - Distributes entire population according to strategy
# - Warning: May be slow for large tenants
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  # No odata_filter - gets ALL users
  shard_count             = 5
  strategy                = "round-robin"
  recalculate_on_next_run = true
}

output "shard_0_count" {
  description = "Users in shard 0 (from entire tenant)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard 1"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard 2"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_3_count" {
  description = "Users in shard 3"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}

output "shard_4_count" {
  description = "Users in shard 4"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_4"])
}

output "total_users" {
  description = "Total users across all shards (entire tenant)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_4"])
}

output "warning_note" {
  description = "Performance consideration"
  value       = "No filter retrieves ALL resources - may be slow for large tenants. Use OData filters to narrow scope."
}
