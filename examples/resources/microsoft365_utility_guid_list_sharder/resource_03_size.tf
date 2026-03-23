# ==============================================================================
# Example 3: Size Strategy
#
# Size distributes GUIDs by absolute counts rather than ratios. This is the
# right choice when stakeholders specify requirements in headcount rather than
# percentages: "we need exactly 50 users in the pilot, 200 in the broader wave,
# and everyone else in the final ring".
#
# Use -1 as the last value to mean "all remaining GUIDs". Only the last element
# may be -1. Without -1 the last shard is exactly the specified size and any
# remaining GUIDs are discarded — an intentional way to cap ring sizes.
#
# As with percentage, the seed controls membership assignment not shard size;
# sizes are always exactly as specified regardless of seed.
# ==============================================================================

# Absolute-size rings — last shard captures all remaining users via -1 sentinel
resource "microsoft365_utility_guid_list_sharder" "windows_update_rings" {
  resource_type           = "devices"
  odata_filter            = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_sizes             = [50, 200, -1] # Exactly 50 → 200 → all remaining
  strategy                = "size"
  seed                    = "windows-updates-2026"
  recalculate_on_next_run = false
}

# Capped rings — useful when a pilot must not exceed a fixed headcount.
# No -1 means users beyond the sum of sizes are deliberately excluded.
# The fourth shard would receive the overflow if added later.
resource "microsoft365_utility_guid_list_sharder" "it_pilot_capped" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and department eq 'IT'"
  shard_sizes             = [10, 25] # Hard cap: pilot=10, validation=25, rest not yet targeted
  strategy                = "size"
  seed                    = "it-dept-pilot-2026"
  recalculate_on_next_run = false
}

output "windows_update_ring_sizes" {
  description = "Device count per Windows Update ring (ring_2 = all remaining devices)"
  value = {
    ring_0_validation = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"])
    ring_1_pilot      = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"])
    ring_2_production = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"])
  }
}
