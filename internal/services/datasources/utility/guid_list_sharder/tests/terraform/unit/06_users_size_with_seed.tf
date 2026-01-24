# Test 06: Users - Size Strategy (With Seed)
#
# Purpose: Verify deterministic size-based distribution with reproducibility
#
# Use Case: "Same pilot capacity across multiple rollouts with different users"
#
# Expected Behavior:
# - Exact shard sizes (25, 50, remainder)
# - Fisher-Yates shuffle with seed, then size-based split
# - Reproducible across Terraform runs

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_sizes   = [5, 10, -1]  # 5 pilot, 10 broader, rest for full
  strategy      = "size"
  seed          = "mfa-rollout-2024"
}

output "pilot_count" {
  description = "Users in pilot (should be exactly 5)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_count" {
  description = "Users in broader pilot (should be exactly 10)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_count" {
  description = "Users in full rollout (all remaining)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "distribution_note" {
  description = "Seed ensures reproducibility"
  value       = "Same seed produces identical distribution across runs, different seeds vary who's in pilot"
}
