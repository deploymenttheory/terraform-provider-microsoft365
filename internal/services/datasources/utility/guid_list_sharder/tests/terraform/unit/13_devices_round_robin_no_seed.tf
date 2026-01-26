# ==============================================================================
# Test 13: Devices - Round-Robin Strategy (No Seed)
#
# Purpose: Verify round-robin distribution produces exactly equal ring sizes
# using API order
#
# Use Case: Quick one-time equal split of devices for capacity testing
#
# Expected Behavior:
# - Exactly equal ring sizes (within ±1)
# - Uses API order (may change between runs)
# - Fast processing
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "devices"
  odata_filter  = "operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  # No seed - uses API order
}

output "ring_0_count" {
  description = "Devices in Ring 0 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "ring_1_count" {
  description = "Devices in Ring 1 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "ring_2_count" {
  description = "Devices in Ring 2 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "ring_3_count" {
  description = "Devices in Ring 3 (exactly 25% ±1)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}

output "total_devices" {
  description = "Total devices distributed"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}
