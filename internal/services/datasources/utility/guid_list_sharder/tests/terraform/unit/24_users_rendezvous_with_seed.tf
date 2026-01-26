# ==============================================================================
# Test 24: Users - Rendezvous (HRW) Strategy (With Seed)
#
# Purpose: Verify Rendezvous Hashing with explicit seed produces different
# distribution than empty seed while maintaining balance
#
# Use Case: Multiple independent rollouts where you want different users in
# pilot groups. MFA rollout uses seed "mfa-2024", Windows Updates uses
# "windows-2024" - same user gets different ring assignments
#
# Expected Behavior:
# - Approximately equal shard sizes (~probabilistic balance, not perfect ±1)
# - Always deterministic (reproducible across runs with same seed)
# - Different seed = different distribution (distributes pilot burden)
# - Minimal disruption when shard count changes (only ~1/n GUIDs move)
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_filter  = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "rendezvous"
  seed          = "deployment-ring-2024" # Different seeds = different distributions
}

output "shard_0_count" {
  description = "Users in shard 0 (should be ~25%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard 1 (should be ~25%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard 2 (should be ~25%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_3_count" {
  description = "Users in shard 3 (should be ~25%)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}

output "size_variance" {
  description = "Max difference between largest and smallest shard (may be higher than round-robin)"
  value       = max(length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])) - min(length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]), length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"]))
}

output "total_users" {
  description = "Total users distributed (should equal sum of all shards)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}
