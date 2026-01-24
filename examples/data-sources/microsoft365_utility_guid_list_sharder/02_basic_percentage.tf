# Basic Percentage: Custom ratio distribution
# Common pattern: 10% pilot, 30% broader, 60% full

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "ca-rollout-2024"
}

# Output shard counts
output "distribution" {
  value = {
    pilot_10pct   = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    broader_30pct = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    full_60pct    = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
  }
}
