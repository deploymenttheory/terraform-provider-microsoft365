# Test 14: Devices - Hash Strategy (With Seed)
#
# Purpose: Verify hash-based distribution with seed produces different
# distributions for different update types
#
# Use Case: Running multiple independent update rollouts (Windows Updates,
# App Updates, Driver Updates) where you want different devices in early
# rings for each type
#
# Expected Behavior:
# - Different seeds produce different distributions
# - Device X might be in Ring 0 for Windows Updates but Ring 2 for App Updates
# - Distributes validation burden across different devices

data "microsoft365_utility_guid_list_sharder" "windows_updates" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "windows-updates-2024"
}

data "microsoft365_utility_guid_list_sharder" "app_updates" {
  resource_type = "devices"
  odata_query   = "$filter=operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "app-updates-2024"  # Different seed = different distribution
}

output "windows_ring_0_count" {
  description = "Devices in Windows Updates Ring 0 (early)"
  value       = length(data.microsoft365_utility_guid_list_sharder.windows_updates.shards["shard_0"])
}

output "app_ring_0_count" {
  description = "Devices in App Updates Ring 0 (likely different devices)"
  value       = length(data.microsoft365_utility_guid_list_sharder.app_updates.shards["shard_0"])
}

output "distribution_note" {
  description = "Verification note"
  value       = "Different seeds ensure same device isn't always in early ring across all update types"
}
