# ==============================================================================
# Example 1: Round-Robin Strategy
#
# Round-robin cycles through shards in order, guaranteeing perfect ±1 balance.
# It is the right choice when you need equal-sized rings and predictable sizes
# matter more than per-GUID determinism.
#
# Without a seed the distribution follows the API return order, which can change
# between runs as users are added or removed. With a seed the input list is
# shuffled deterministically before round-robin is applied, so the same tenant
# state always produces the same assignments — useful for A/B testing or when
# you need to communicate exact ring membership to stakeholders.
#
# recalculate_on_next_run = false (recommended default)
#   - On the very first apply, assignments are always computed regardless of
#     this value — you do not need a two-step "set true, then set false" dance.
#   - On all subsequent applies the stored assignments are returned from state
#     unchanged; no Graph API call is made, and new tenant members are ignored.
#   - Set to true only when you intentionally want to reshard (e.g. after a
#     large onboarding wave), then set back to false to re-lock.
# ==============================================================================

# Without seed — uses API return order; equal ring sizes but not reproducible
resource "microsoft365_utility_guid_list_sharder" "mfa_rollout_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4 # Ring 0 → Ring 3; equal ±1 split
  strategy                = "round-robin"
  recalculate_on_next_run = false
}

# With seed — deterministic; same tenant state always produces the same rings.
# Use a different seed per rollout (e.g. "windows-updates-2026" vs "mfa-2026")
# so that users who are in the pilot ring for one initiative are NOT always
# in the pilot ring for every other initiative (prevents pilot fatigue).
resource "microsoft365_utility_guid_list_sharder" "mfa_rollout_rings_seeded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4
  strategy                = "round-robin"
  seed                    = "mfa-rollout-2026" # Change per rollout to vary pilot ring population and reduce pilot fatigue
  recalculate_on_next_run = false
}

# Diagnostic outputs — useful during initial deployment for sanity checks
output "mfa_ring_distribution" {
  description = "Number of users assigned to each MFA rollout ring"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_3"])
    total  = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_3"])
  }
}
