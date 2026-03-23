# ==============================================================================
# Example 6: Managing the recalculate_on_next_run Lifecycle
#
# This example demonstrates the recommended workflow for locking and unlocking
# shard assignments over the lifetime of a rollout.
#
# Day 1 — Initial deployment
#   Set recalculate_on_next_run = false. The first terraform apply always
#   computes assignments from scratch regardless of this flag (there is no
#   prior state). Assignments are then locked.
#
# Ongoing — Steady state
#   Leave recalculate_on_next_run = false. Plans are fast (no API call), and
#   new users joining the tenant do not cause membership churn. The policy
#   targets exactly the same population on every apply.
#
# Intentional reshard — e.g. after a large onboarding wave or shard_count change
#   1. Set recalculate_on_next_run = true
#   2. Run terraform apply — Update re-queries Graph and recomputes shards
#   3. Set recalculate_on_next_run = false again to re-lock
#
# The distinction between Read and Update:
#   recalculate = false + no config change → Read returns cached state (no API call)
#   recalculate = false + config change    → Update preserves existing shards, saves new config
#   recalculate = true  + any apply        → Update re-queries Graph and reshards
# ==============================================================================

# Phase 1: Initial deployment — flag=false, assignments computed on first apply
resource "microsoft365_utility_guid_list_sharder" "compliance_rollout" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "compliance-2026"
  recalculate_on_next_run = false # Safe from day one; first apply always computes
}

# Phase 2 (illustrative — replace the above block in practice):
# After a large onboarding wave you decide to rebalance.
# Step 1: uncomment recalculate = true and apply
# Step 2: revert to false and apply again to lock
#
# resource "microsoft365_utility_guid_list_sharder" "compliance_rollout" {
#   resource_type           = "users"
#   odata_filter            = "accountEnabled eq true and userType eq 'Member'"
#   shard_count             = 3
#   strategy                = "round-robin"
#   seed                    = "compliance-2026"
#   recalculate_on_next_run = true    # Triggers reshard on next apply; set back to false after
# }

output "compliance_rollout_status" {
  description = "Ring membership counts and whether assignments are currently locked"
  value = {
    ring_0_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_0"])
    ring_1_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_1"])
    ring_2_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_2"])
    assignments_locked = !microsoft365_utility_guid_list_sharder.compliance_rollout.recalculate_on_next_run
  }
}
