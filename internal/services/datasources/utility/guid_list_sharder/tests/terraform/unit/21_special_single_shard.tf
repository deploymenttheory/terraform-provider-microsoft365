# Test 21: Special Case - Single Shard (All Users in One Set)
#
# Purpose: Verify edge case where shard_count = 1 produces a single set
# containing all users
#
# Use Case: Converting a query result into a set for use with resources
# without actually splitting
#
# Expected Behavior:
# - All GUIDs in shard_0
# - Useful for filtering + converting to set format

data "microsoft365_utility_guid_list_sharder" "test" {
  resource_type = "users"
  odata_query   = "$filter=accountEnabled eq true and userType eq 'Member'"
  shard_count   = 1
  strategy      = "hash"
  # No split - all users in one shard
}

output "all_users" {
  description = "All users in a single set (shard_0)"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "user_count" {
  description = "Total number of users"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "use_case_note" {
  description = "When to use single shard"
  value       = "Useful for converting filtered query results into set format without splitting"
}
