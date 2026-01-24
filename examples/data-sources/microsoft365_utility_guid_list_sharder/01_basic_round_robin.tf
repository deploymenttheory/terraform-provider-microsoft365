# Basic Round-Robin: Equal distribution across shards
# Perfect ±1 balance guaranteed

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-rollout-2024"
}

# Output shard counts
output "distribution" {
  value = {
    shard_0 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    shard_1 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    shard_2 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
    shard_3 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_3"])
  }
}
