# Test 17: Devices - Percentage Strategy (No Seed)
#
# Purpose: Verify percentage-based distribution produces industry-standard
# Windows Update rings using API order
#
# Use Case: Quick setup of standard update rings (5% canary, 15% early, 80% broad)
#
# Expected Behavior:
# - Ring sizes match specified percentages
# - Uses API order (may change between runs)
# - Last ring gets all remaining devices

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "devices"
  odata_query       = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  # No seed - uses API order
}

output "canary_ring_count" {
  description = "Devices in Canary Ring (~5%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "early_ring_count" {
  description = "Devices in Early Ring (~15%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "broad_ring_count" {
  description = "Devices in Broad Ring (~80%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "canary_devices" {
  description = "Device GUIDs in Canary Ring"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "canary_percentage" {
  description = "Actual canary ring percentage"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) / (length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])) * 100
}
