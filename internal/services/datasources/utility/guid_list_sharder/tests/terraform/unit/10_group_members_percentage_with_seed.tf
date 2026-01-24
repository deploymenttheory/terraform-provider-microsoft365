# Test 12: Group Members - Percentage Strategy (With Seed)
#
# Purpose: Verify percentage-based splitting with seed produces custom-sized
# subgroups AND reproducible results
#
# Use Case: Phased access to group resources with specific percentages where
# you need same members in each tier every time
#
# Expected Behavior:
# - Subgroup sizes match specified percentages
# - Deterministic shuffle before percentage split
# - Same seed = same tier membership every time

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "group_members"
  group_id          = "12345678-1234-1234-1234-123456789abc"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [20, 30, 50]
  strategy          = "percentage"
  seed              = "resource-access-tiers-2024"  # Makes tiers reproducible
}

output "tier_1_members" {
  description = "Tier 1 members (~20%, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "tier_1_count" {
  description = "Number of members in Tier 1"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "tier_2_count" {
  description = "Number of members in Tier 2"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "tier_3_count" {
  description = "Number of members in Tier 3"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
