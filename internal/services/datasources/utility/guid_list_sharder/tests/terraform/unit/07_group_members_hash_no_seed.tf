# Test 07: Group Members - Hash Strategy (No Seed)
#
# Purpose: Verify hash-based splitting of group members without seed produces
# consistent split across all instances
#
# Use Case: Split an existing group into multiple subgroups with consistent
# membership that's identical everywhere
#
# Expected Behavior:
# - Same member always goes to same subgroup
# - Consistent across all instances and Terraform runs
# - Approximately equal subgroup sizes

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  # No seed - consistent split everywhere
}

output "subgroup_0_members" {
  description = "Members for subgroup 0 (consistent everywhere)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "subgroup_0_count" {
  description = "Number of members in subgroup 0"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "subgroup_1_count" {
  description = "Number of members in subgroup 1"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "subgroup_2_count" {
  description = "Number of members in subgroup 2"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
