# Test 02: Users - Hash Strategy (With Seed)
#
# Purpose: Verify hash-based distribution with seed produces different
# distributions for different seeds (distributes pilot burden)
#
# Use Case: Running multiple independent rollouts where you want different
# users in pilot groups for each rollout to prevent pilot fatigue
#
# Expected Behavior:
# - Different seeds produce different distributions
# - Same seed always produces same distribution (reproducible)
# - User X might be in shard_0 for MFA but shard_2 for Windows Updates

data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type = "users"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "mfa-rollout-2024"
}

data "microsoft365_utility_guid_list_sharder" "windows_rollout" {
  resource_type = "users"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "hash"
  seed          = "windows-rollout-2024"  # Different seed = different distribution
}

output "mfa_shard_0_count" {
  description = "Users in MFA pilot (shard 0)"
  value       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
}

output "windows_shard_0_count" {
  description = "Users in Windows pilot (shard 0) - likely different users"
  value       = length(data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_0"])
}

output "distribution_note" {
  description = "Verification note"
  value       = "With different seeds, same users will be in different shards across rollouts"
}
