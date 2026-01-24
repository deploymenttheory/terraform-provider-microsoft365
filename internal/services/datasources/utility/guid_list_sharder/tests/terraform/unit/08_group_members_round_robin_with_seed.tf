# Test 10: Group Members - Round-Robin Strategy (With Seed)
#
# Purpose: Verify round-robin splitting with seed produces exactly equal
# subgroup sizes AND reproducible results
#
# Use Case: Split group equally with ability to recreate exact same split
#
# Expected Behavior:
# - Exactly equal subgroup sizes (within ±1)
# - Deterministic shuffle before round-robin
# - Same seed = same member assignments every time

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  odata_query   = "accountEnabled eq true"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "group-split-2024"  # Makes split reproducible
}

output "team_a_members" {
  description = "Members for Team A (exactly 50% ±1, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "team_b_members" {
  description = "Members for Team B (exactly 50% ±1, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]
}

output "team_a_count" {
  description = "Number of members in Team A"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "team_b_count" {
  description = "Number of members in Team B"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}
