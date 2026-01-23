# Test 11: Group Members - Percentage Strategy (No Seed)
#
# Purpose: Verify percentage-based splitting of group members produces
# custom-sized subgroups using API order
#
# Use Case: Quick phased access to group resources with specific percentages
#
# Expected Behavior:
# - Subgroup sizes match specified percentages
# - Uses API order (may change between runs)
# - Last subgroup gets all remaining members

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "group_members"
  group_id          = "12345678-1234-1234-1234-123456789abc"
  odata_query       = "$filter=accountEnabled eq true"
  shard_percentages = [20, 30, 50]
  strategy          = "percentage"
  # No seed - uses API order
}

output "tier_1_count" {
  description = "Members in Tier 1 (~20%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "tier_2_count" {
  description = "Members in Tier 2 (~30%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "tier_3_count" {
  description = "Members in Tier 3 (~50%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_members" {
  description = "Total members distributed"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
