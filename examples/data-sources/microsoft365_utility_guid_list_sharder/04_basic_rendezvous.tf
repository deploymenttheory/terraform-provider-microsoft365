# Basic Rendezvous: Stable distribution when shard count changes
# Seed is REQUIRED for rendezvous strategy
# Only ~1/n GUIDs move when adding shards

data "microsoft365_utility_guid_list_sharder" "users_stable" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "stable-deployment-2024"
}

# Output shard counts
output "distribution" {
  value = {
    ring_0 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_0"])
    ring_1 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_1"])
    ring_2 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_2"])
  }
}

# When you change shard_count from 3 to 4:
# - Only ~25% of users will move to new ring_3
# - ~75% stay in their original ring
# - Compare with round-robin/percentage: ~75% would move
