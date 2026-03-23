# ==============================================================================
# Scenario: Sharder → Security Groups → Group Member Assignments
#
# This is the primary integration pattern. Shard GUIDs from the sharder are
# distributed into separate Entra ID security groups. Downstream resources
# (Intune policies, Conditional Access, Windows Update rings, etc.) then target
# those groups rather than individual users or devices, following the standard
# Microsoft 365 group-based targeting model.
#
# Pattern:
#   sharder → creates shard sets of GUIDs
#   group   → one group per ring (created once, stable IDs)
#   group_member_assignment (for_each) → adds each shard's GUIDs to the group
#
# The for_each on group_member_assignment means each GUID gets its own resource
# instance tracked in state. Terraform will add/remove only the specific members
# that change — it will NOT recreate the group or remove members from the wrong ring.
#
# Use case: Phased MFA rollout across 4 rings, 10/20/30/40% of enabled members
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "mfa_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 20, 30, 40]
  strategy                = "percentage"
  seed                    = "mfa-2026"
  recalculate_on_next_run = false
}

# One security group per ring — created independently of the sharder so that
# their IDs remain stable even if shards are later recomputed.
resource "microsoft365_graph_beta_groups_group" "mfa_ring_0" {
  display_name     = "MFA Rollout - Ring 0 (10% Pilot)"
  mail_nickname    = "mfa-ring-0-pilot"
  description      = "Initial 10% pilot for MFA enforcement — validate before expanding"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_1" {
  display_name     = "MFA Rollout - Ring 1 (20% Broader)"
  mail_nickname    = "mfa-ring-1-broader"
  description      = "Second wave — 20% of users after Ring 0 validation"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_2" {
  display_name     = "MFA Rollout - Ring 2 (30% Broad)"
  mail_nickname    = "mfa-ring-2-broad"
  description      = "Third wave — 30% of users"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_3" {
  display_name     = "MFA Rollout - Ring 3 (40% Production)"
  mail_nickname    = "mfa-ring-3-production"
  description      = "Final wave — remaining 40% of users"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

# Populate each group using for_each so Terraform tracks each membership
# individually. Adding or removing users only touches the affected member
# resource instances — the groups themselves are never recreated.
resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_0_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_0.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_1_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_1.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_2_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_2.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_3_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_3.id
  member_id          = each.value
  member_object_type = "User"
}

output "mfa_group_ids" {
  description = "Security group object IDs — reference these in Intune or CA policies"
  value = {
    ring_0_pilot      = microsoft365_graph_beta_groups_group.mfa_ring_0.id
    ring_1_broader    = microsoft365_graph_beta_groups_group.mfa_ring_1.id
    ring_2_broad      = microsoft365_graph_beta_groups_group.mfa_ring_2.id
    ring_3_production = microsoft365_graph_beta_groups_group.mfa_ring_3.id
  }
}

output "mfa_ring_headcount" {
  description = "Number of users in each ring — locked from the point of first apply"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
  }
}
