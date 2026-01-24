# Test 17: Devices - Size Strategy (No Seed)
#
# Purpose: Verify absolute size-based distribution for devices
#
# Use Case: "Test updates on exactly 10 devices, then 50, then all"
#
# Expected Behavior:
# - Exact device counts per ring
# - Uses API order
# - Last ring gets all remaining devices

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows'"
  shard_sizes   = [6, 18, -1]  # 6 test ring, 18 pilot ring, rest for broad
  strategy      = "size"
  # No seed - uses API order
}

output "test_ring_count" {
  description = "Devices in test ring (should be exactly 6)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "pilot_ring_count" {
  description = "Devices in pilot ring (should be exactly 18)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "broad_ring_count" {
  description = "Devices in broad ring (all remaining)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
