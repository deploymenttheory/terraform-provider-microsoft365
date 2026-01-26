# Scenario 2: Group Members → Groups → Conditional Access Policy
# Use case: Roll out new CA policy to IT department in phases without nested groups

# Distribute existing IT department group members into 3 deployment rings
data "microsoft365_utility_guid_list_sharder" "it_dept_members" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc" # IT Department Group ID
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "it-mfa-policy-2024"
}

# Create deployment ring groups from IT department members
resource "microsoft365_graph_beta_group" "it_ring_0_pilot" {
  display_name     = "IT MFA Rollout - Ring 0 (Pilot)"
  mail_nickname    = "it-mfa-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "it_ring_1_validation" {
  display_name     = "IT MFA Rollout - Ring 1 (Validation)"
  mail_nickname    = "it-mfa-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "it_ring_2_full" {
  display_name     = "IT MFA Rollout - Ring 2 (Full)"
  mail_nickname    = "it-mfa-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_2"]
}

# Single CA policy targeting all IT deployment rings
resource "microsoft365_graph_beta_conditional_access_policy" "it_mfa_policy" {
  display_name = "Require MFA - IT Department (Phased Rollout)"
  state        = "enabled"

  conditions {
    users {
      include_groups = [
        microsoft365_graph_beta_group.it_ring_0_pilot.id,
        microsoft365_graph_beta_group.it_ring_1_validation.id,
        microsoft365_graph_beta_group.it_ring_2_full.id
      ]
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
