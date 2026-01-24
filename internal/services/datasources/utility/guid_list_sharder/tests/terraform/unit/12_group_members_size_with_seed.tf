# Test 12: Group Members - Size Strategy (With Seed)
#
# Purpose: Verify deterministic size-based distribution for group members
#
# Use Case: "Reproducible pilot group selection from department"
#
# Expected Behavior:
# - Exact shard sizes with Fisher-Yates shuffle
# - Reproducible distribution
# - Different seeds produce different member selections

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  shard_sizes   = [5, 10, -1]  # 5 pilot, 10 broader, rest for full
  strategy      = "size"
  seed          = "department-pilot-2024"
}

output "pilot_count" {
  description = "Members in pilot (should be exactly 5)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_count" {
  description = "Members in broader pilot (should be exactly 10)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_count" {
  description = "Members in full rollout (all remaining)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
