# ==============================================================================
# Test 18: Devices - Size Strategy (With Seed)
#
# Purpose: Verify deterministic size-based distribution for devices with different seeds
#
# Use Case: "Different devices in test ring for Windows vs App updates"
#
# Expected Behavior:
# - Exact device counts per ring
# - Fisher-Yates shuffle with seed
# - Different seeds distribute update burden across devices
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "windows_updates" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows'"
  shard_sizes   = [6, 12, -1]  # 6 test, 12 pilot, rest for broad
  strategy      = "size"
  seed          = "windows-updates-2024"
}

data "microsoft365_utility_guid_list_sharder" "app_updates" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows'"
  shard_sizes   = [6, 12, -1]  # Same sizes, different seed
  strategy      = "size"
  seed          = "app-updates-2024"
}

output "windows_test_count" {
  description = "Devices in Windows test ring"
  value       = length(data.microsoft365_utility_guid_list_sharder.windows_updates.shards["shard_0"])
}

output "app_test_count" {
  description = "Devices in App test ring"
  value       = length(data.microsoft365_utility_guid_list_sharder.app_updates.shards["shard_0"])
}

output "distribution_note" {
  description = "Different seeds ensure same device isn't always in test ring"
  value       = "Different seeds mean different devices get early updates for each rollout type"
}
