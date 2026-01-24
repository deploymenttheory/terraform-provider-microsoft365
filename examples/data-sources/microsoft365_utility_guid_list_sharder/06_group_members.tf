# Group Members Distribution: Shard members of an existing Entra ID group
# Use case: Deploy policy to IT department in phases without creating additional nested groups

data "microsoft365_utility_guid_list_sharder" "it_dept_phases" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc" # IT Department Group ID
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "it-dept-pilot-2026"
}

# Phase 1: IT Pilot (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_1" {
  display_name = "New IT Policy - Phase 1 (IT Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Phase 2: IT Validation (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_2" {
  display_name = "New IT Policy - Phase 2 (IT Validation)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Phase 3: IT Full (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_3" {
  display_name = "New IT Policy - Phase 3 (IT Full)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Monitor IT department phase distribution
output "it_dept_phase_distribution" {
  value = {
    phase_1_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"])
    phase_2_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"])
    phase_3_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"])
    total_it_members = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"])
  }
  description = "Equal distribution across 3 phases (perfect ±1 balance with round-robin)"
}
