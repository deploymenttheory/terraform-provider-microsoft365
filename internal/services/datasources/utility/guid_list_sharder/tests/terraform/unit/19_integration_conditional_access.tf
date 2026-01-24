# ==============================================================================
# Test 19: Integration - Conditional Access Policy Usage
#
# Purpose: Demonstrate how shards integrate directly with Conditional Access
# policies for progressive MFA rollout
#
# Use Case: Real-world MFA deployment pattern
#
# Note: This is a demonstration. Actual resource creation would happen in
# acceptance tests, not unit tests
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-rollout-2024"
}

# Example: Phase 1 - Pilot (10%)
# resource "microsoft365_graph_beta_conditional_access_policy" "mfa_pilot" {
#   display_name = "MFA Required - Phase 1 Pilot"
#   state        = "enabled"
#   
#   conditions {
#     users {
#       include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]
#     }
#   }
#   
#   grant_controls {
#     built_in_controls = ["mfa"]
#     operator          = "OR"
#   }
# }

# Example: Phase 2 - Broader Pilot (30%)
# resource "microsoft365_graph_beta_conditional_access_policy" "mfa_broader" {
#   display_name = "MFA Required - Phase 2 Broader"
#   state        = "enabledForReportingButNotEnforced"  # Test first
#   
#   conditions {
#     users {
#       include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"]
#     }
#   }
#   
#   grant_controls {
#     built_in_controls = ["mfa"]
#     operator          = "OR"
#   }
# }

output "pilot_users" {
  description = "Users in Phase 1 pilot (ready for CA policy)"
  value       = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]
}

output "pilot_count" {
  description = "Number of users in pilot phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
}

output "broader_count" {
  description = "Number of users in broader pilot phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"])
}

output "full_count" {
  description = "Number of users in full rollout phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
}

output "usage_note" {
  description = "Integration guidance"
  value       = "Shard outputs are sets - can be used directly in conditions.users.include_users"
}
