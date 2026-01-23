# Test 13: Devices - Hash Strategy (No Seed)
#
# Purpose: Verify hash-based distribution of devices without seed produces
# consistent distribution across all instances
#
# Use Case: Creating standard device groups (e.g., update rings) that should
# be identical across all policies
#
# Expected Behavior:
# - Same device always goes to same ring
# - Consistent across all instances and Terraform runs
# - Approximately equal ring sizes

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  # No seed - ensures identical rings everywhere
}

output "ring_0_count" {
  description = "Devices in Ring 0"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "ring_1_count" {
  description = "Devices in Ring 1"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "ring_2_count" {
  description = "Devices in Ring 2"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "ring_0_devices" {
  description = "Device GUIDs in Ring 0"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}
