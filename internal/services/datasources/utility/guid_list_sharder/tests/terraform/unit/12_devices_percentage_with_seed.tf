# Test 18: Devices - Percentage Strategy (With Seed)
#
# Purpose: Verify percentage-based distribution with seed produces
# industry-standard Windows Update rings with reproducible membership
#
# Use Case: Production Windows Update rings (5% canary, 15% early, 80% broad)
# where you need same devices in each ring every time
#
# Expected Behavior:
# - Ring sizes match specified percentages
# - Deterministic shuffle before percentage split
# - Same seed = same ring membership every time

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "devices"
  odata_query       = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "windows-update-rings-2024"  # Makes rings reproducible
}

output "canary_ring_devices" {
  description = "Devices in Canary Ring (~5%, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "canary_ring_count" {
  description = "Number of devices in Canary Ring"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "early_ring_count" {
  description = "Number of devices in Early Ring"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "broad_ring_count" {
  description = "Number of devices in Broad Ring"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_devices" {
  description = "Total devices across all rings"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "reproducibility_note" {
  description = "Verification note"
  value       = "With seed, same devices will always be in same rings across Terraform runs"
}
