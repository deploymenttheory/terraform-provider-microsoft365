# Test 06: Users - Percentage Strategy (With Seed)
#
# Purpose: Verify percentage-based distribution with seed produces custom-sized
# shards AND reproducible results
#
# Use Case: Structured phased rollout (10% → 30% → 60%) where you need the
# SAME users in each phase every time
#
# Expected Behavior:
# - Shard sizes match specified percentages
# - Deterministic shuffle before percentage split
# - Same seed = same phase membership every time

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type     = "users"
  odata_query       = "$filter=accountEnabled eq true and userType eq 'Member'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-phased-2024"  # Makes distribution reproducible
}

output "pilot_users" {
  description = "Pilot phase users (~10%, reproducible)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "pilot_count" {
  description = "Number of users in pilot phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "broader_pilot_count" {
  description = "Number of users in broader pilot phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "full_rollout_count" {
  description = "Number of users in full rollout phase"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "reproducibility_note" {
  description = "Verification note"
  value       = "With seed, same users will always be in same phases across Terraform runs"
}
