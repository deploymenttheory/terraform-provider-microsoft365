# Test 16: Devices - Round-Robin Strategy (With Seed)
#
# Purpose: Verify round-robin distribution with seed produces exactly equal
# ring sizes AND reproducible results
#
# Use Case: Device capacity testing or load balancing where you need exact
# equal distribution that you can recreate
#
# Expected Behavior:
# - Exactly equal ring sizes (within ±1)
# - Deterministic shuffle before round-robin
# - Same seed = same device assignments every time

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "device-load-balancing-2024"  # Makes distribution reproducible
}

output "pool_a_devices" {
  description = "Devices in Pool A (exactly 50% ±1, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "pool_b_devices" {
  description = "Devices in Pool B (exactly 50% ±1, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]
}

output "pool_a_count" {
  description = "Number of devices in Pool A"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "pool_b_count" {
  description = "Number of devices in Pool B"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "is_balanced" {
  description = "Confirms equal split (should be 0 or 1)"
  value       = abs(length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) - length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]))
}
