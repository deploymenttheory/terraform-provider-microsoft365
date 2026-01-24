# Basic Size: Absolute count distribution
# Use -1 for "all remaining" in last position

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1]
  strategy      = "size"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1]
  strategy      = "size"
  seed          = "it-pilot-2024"
}

# Output shard counts
output "distribution" {
  value = {
    pilot_exact_50      = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    validation_exact_100 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    broad_all_remaining = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
  }
}
