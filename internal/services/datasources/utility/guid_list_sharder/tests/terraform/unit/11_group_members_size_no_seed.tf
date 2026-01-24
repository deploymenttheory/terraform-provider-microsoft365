# Test 11: Group Members - Size Strategy (No Seed)
#
# Purpose: Verify absolute size-based distribution for group members
#
# Use Case: "Split large group into fixed-size pilot subgroups"
#
# Expected Behavior:
# - Exact shard sizes from group membership
# - Uses API order
# - Last shard gets all remaining members

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  shard_sizes   = [5, 15, -1]  # 5 pilot, 15 broader, rest for full
  strategy      = "size"
  # No seed - uses API order
}

output "pilot_count" {
  description = "Members in pilot (should be exactly 5)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_count" {
  description = "Members in broader pilot (should be exactly 15)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_count" {
  description = "Members in full rollout (all remaining)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
