# ==============================================================================
# Example 4: Rendezvous (Highest Random Weight) Strategy
#
# Rendezvous assigns each GUID independently using a per-(guid, shard, seed)
# hash score. The shard with the highest score wins for each GUID. Because
# each GUID's assignment is computed independently, adding or removing a shard
# only moves the GUIDs that must move: approximately 1/n of the population
# when going from n to n+1 shards, versus ~75% with position-based strategies.
#
# This stability makes rendezvous the right choice when:
#   - The ring count is expected to grow during the rollout lifecycle
#   - Minimising reassignment churn matters (e.g. avoiding repeated policy
#     application to devices that already received an earlier ring's config)
#   - You need per-GUID determinism without caring about equal-sized rings
#
# Distribution is probabilistic: with 12 users and 4 shards you expect ~3 per
# shard, but natural hash variance means actual counts differ. Assertions on
# per-shard size are not meaningful — assert only total count.
#
# The seed affects which shard each GUID is assigned to; an empty seed is used
# internally when omitted, which is still deterministic.
# ==============================================================================

# Starting with 3 rings — ring count may grow later without excessive churn
resource "microsoft365_utility_guid_list_sharder" "compliance_rings_initial" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 3
  strategy                = "rendezvous"
  seed                    = "compliance-baseline-2026"
  recalculate_on_next_run = false
}

# After expanding to 4 rings: only ~25% of users move to the new shard_3.
# The remaining ~75% stay exactly where they were — no unnecessary disruption.
# To expand: change shard_count to 4, set recalculate_on_next_run = true,
# apply, then set recalculate_on_next_run back to false.
resource "microsoft365_utility_guid_list_sharder" "compliance_rings_expanded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4 # One extra ring added; only ~25% of users move
  strategy                = "rendezvous"
  seed                    = "compliance-baseline-2026" # Same seed preserves prior assignments
  recalculate_on_next_run = false
}

output "compliance_ring_distribution" {
  description = "Approximate per-ring user counts (rendezvous is probabilistic — totals are authoritative)"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_3"])
    total  = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_3"])
  }
}
