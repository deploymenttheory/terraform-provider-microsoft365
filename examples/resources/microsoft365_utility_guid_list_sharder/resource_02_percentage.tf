# ==============================================================================
# Example 2: Percentage Strategy
#
# Percentage distributes GUIDs according to the ratios you specify. The sizes
# are arithmetically deterministic — [10, 30, 60] always produces shards of
# 10%, 30%, and 60% regardless of whether a seed is supplied. The seed only
# controls *which* users land in each shard, not how many.
#
# This is the most natural fit for standard phased rollouts where stakeholders
# think in percentages: "start with 10%, validate, then expand to 30%, then all".
#
# shard_percentages must sum to exactly 100.
# Each shard is named shard_0, shard_1, ... matching the list index.
# ==============================================================================

# Without seed — order follows the Graph API response; not reproducible across runs
resource "microsoft365_utility_guid_list_sharder" "ca_policy_phased" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60] # Pilot → Broader → Full rollout
  strategy                = "percentage"
  recalculate_on_next_run = false
}

# With seed — deterministic; pilot group is the same every run for this rollout.
# Using a rollout-specific seed means the 10% pilot for conditional access is
# drawn from a different slice of the tenant than the 10% pilot for MFA —
# distributing the pilot burden across the organisation.
resource "microsoft365_utility_guid_list_sharder" "ca_policy_phased_seeded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "ca-policy-rollout-2026"
  recalculate_on_next_run = false
}

output "ca_policy_phase_sizes" {
  description = "Headcount per deployment phase — verify ratios before enabling the policy"
  value = {
    pilot_10pct   = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_0"])
    broader_30pct = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_1"])
    full_60pct    = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_2"])
  }
}
