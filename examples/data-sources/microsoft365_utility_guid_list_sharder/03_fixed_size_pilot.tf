# Size-Based Distribution: Fixed pilot group sizes
# Use case: Compliance requires exactly 50 pilot users, 100 validation users

data "microsoft365_utility_guid_list_sharder" "compliance_pilot" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1] # 50 pilot, 100 validation, remainder for broad
  strategy      = "size"
  seed          = "compliance-pilot-2026"
}

# Pilot Group: Exactly 50 users
resource "microsoft365_graph_beta_group" "compliance_pilot" {
  display_name     = "Compliance Policy - Pilot (50 users)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"]
}

# Validation Group: Exactly 100 users
resource "microsoft365_graph_beta_group" "compliance_validation" {
  display_name     = "Compliance Policy - Validation (100 users)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"]
}

# Broad Deployment: All remaining IT users
resource "microsoft365_graph_beta_group" "compliance_broad" {
  display_name     = "Compliance Policy - Broad (All Remaining)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"]
}

# Conditional Access Policy targeting pilot group
resource "microsoft365_graph_beta_conditional_access_policy" "compliance_policy_pilot" {
  display_name = "Device Compliance Required - Pilot"
  state        = "enabled"

  conditions {
    users {
      include_groups = [microsoft365_graph_beta_group.compliance_pilot.id]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["compliantDevice"]
  }
}

# Verify exact counts
output "pilot_group_sizes" {
  value = {
    pilot_count      = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"])
    validation_count = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"])
    broad_count      = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"])
    total_it_users   = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"])
  }
  description = "Pilot should be exactly 50, Validation exactly 100, Broad gets remainder"
}
