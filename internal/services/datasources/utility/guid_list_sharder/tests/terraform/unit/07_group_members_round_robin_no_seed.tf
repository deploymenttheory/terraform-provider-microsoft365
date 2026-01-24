# ==============================================================================
# Test 07: Group Members - Round-Robin Strategy (No Seed)
#
# Purpose: Verify round-robin splitting of group members produces exactly
# equal subgroup sizes using API order
#
# Use Case: Quick one-time split of a group into equal parts
#
# Expected Behavior:
# - Exactly equal subgroup sizes (within ±1)
# - Uses API order (may change between runs)
# - Fast processing
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  # No seed - uses API order
}

output "subgroup_0_count" {
  description = "Members in subgroup 0 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "subgroup_1_count" {
  description = "Members in subgroup 1 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "subgroup_2_count" {
  description = "Members in subgroup 2 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "subgroup_3_count" {
  description = "Members in subgroup 3 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}
