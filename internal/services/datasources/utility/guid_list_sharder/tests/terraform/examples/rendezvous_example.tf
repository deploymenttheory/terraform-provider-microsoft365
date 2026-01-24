# ==============================================================================
# Rendezvous Hashing Strategy Example
# ==============================================================================
# 
# This example demonstrates the Rendezvous (Highest Random Weight) strategy
# for distributing users across deployment rings with minimal disruption
# when ring counts change.
#
# Key Benefits:
# - Always deterministic (reproducible results)
# - Minimal disruption when adding/removing rings (~1/n users move)
# - Each user independently evaluates all rings
# - Superior stability compared to position-based strategies
#
# Use Case: Organizations that expect to add or remove deployment rings
# over time and want to minimize user movement between rings.
# ==============================================================================

terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

# Query all active users and distribute across 3 deployment rings
data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and userType eq 'Member'"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "mfa-rollout-2024" # Different seed = different distribution
}

# Create deployment ring groups
resource "microsoft365_graph_beta_groups_group" "ring_0_pilot" {
  display_name     = "MFA Rollout - Ring 0 (Pilot)"
  mail_nickname    = "mfa-ring-0-pilot"
  security_enabled = true
  
  members = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]
}

resource "microsoft365_graph_beta_groups_group" "ring_1_broad" {
  display_name     = "MFA Rollout - Ring 1 (Broad)"
  mail_nickname    = "mfa-ring-1-broad"
  security_enabled = true
  
  members = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"]
}

resource "microsoft365_graph_beta_groups_group" "ring_2_full" {
  display_name     = "MFA Rollout - Ring 2 (Full)"
  mail_nickname    = "mfa-ring-2-full"
  security_enabled = true
  
  members = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"]
}

# Monitor distribution
output "ring_distribution" {
  description = "Number of users in each deployment ring"
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
    total        = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]) + 
                   length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"]) + 
                   length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
  }
}

# ==============================================================================
# Testing Stability: Add a 4th Ring
# ==============================================================================
# 
# To test the stability benefit, change shard_count from 3 to 4.
# With Rendezvous, only ~25% of users will move (those assigned to new ring_3).
# With round-robin or percentage strategies, ~75% of users would move!
#
# Uncomment below and change shard_count to 4 above to test:

# resource "microsoft365_graph_beta_groups_group" "ring_3_extended" {
#   display_name     = "MFA Rollout - Ring 3 (Extended)"
#   mail_nickname    = "mfa-ring-3-extended"
#   security_enabled = true
#   
#   members = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_3"]
# }
