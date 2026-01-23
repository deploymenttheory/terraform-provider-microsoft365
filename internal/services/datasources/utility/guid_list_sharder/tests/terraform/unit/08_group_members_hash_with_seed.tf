# Test 08: Group Members - Hash Strategy (With Seed)
#
# Purpose: Verify hash-based splitting with seed produces different splits
# for different seeds
#
# Use Case: Split same group differently for different purposes (e.g., pilot
# programs where you don't want same members always being guinea pigs)
#
# Expected Behavior:
# - Different seeds produce different member distributions
# - Same seed always produces same split (reproducible)
# - Member X might be in subgroup 0 for initiative A but subgroup 2 for initiative B

data "microsoft365_utility_guid_list_sharder" "initiative_a" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "initiative-a-2024"
}

data "microsoft365_utility_guid_list_sharder" "initiative_b" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "initiative-b-2024"  # Different seed = different split
}

output "initiative_a_pilot_count" {
  description = "Members in Initiative A pilot group"
  value       = length(data.microsoft365_utility_guid_list_sharder.initiative_a.shards["shard_0"])
}

output "initiative_b_pilot_count" {
  description = "Members in Initiative B pilot group (likely different members)"
  value       = length(data.microsoft365_utility_guid_list_sharder.initiative_b.shards["shard_0"])
}
