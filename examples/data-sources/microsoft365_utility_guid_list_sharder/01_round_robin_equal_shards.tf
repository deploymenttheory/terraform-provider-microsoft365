# Round-Robin Distribution: Perfect equal distribution across 4 deployment rings
# Use case: Equal-sized pilot, validation, pre-production, and production rings

data "microsoft365_utility_guid_list_sharder" "mfa_rings" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-rollout-2026" # Optional: ensures reproducible distribution
}

# Ring 0: Pilot (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_0" {
  display_name = "MFA Required - Ring 0 (Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Ring 1: Validation (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_1" {
  display_name = "MFA Required - Ring 1 (Validation)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Ring 2: Pre-Production (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_2" {
  display_name = "MFA Required - Ring 2 (Pre-Production)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Ring 3: Production (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_3" {
  display_name = "MFA Required - Ring 3 (Production)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Verify perfect distribution
output "ring_distribution" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"])
    ring_3_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
    total_users  = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
  }
  description = "Round-robin guarantees perfect ±1 GUID balance across all rings"
}
