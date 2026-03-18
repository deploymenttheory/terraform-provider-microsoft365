resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Groups — included and excluded deployment audience
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "included_group" {
  display_name     = "wu-ring-included-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "wu-ring-included-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "excluded_group" {
  display_name     = "wu-ring-excluded-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "wu-ring-excluded-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

# ==============================================================================
# Wait for groups to propagate before assigning them to the ring
# ==============================================================================

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.included_group,
    microsoft365_graph_beta_groups_group.excluded_group,
  ]
  create_duration = "30s"
}

# ==============================================================================
# Parent autopatch policy
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "example" {
  display_name = "Quality Update Policy - ${random_string.suffix.result}"
  description  = "Policy managing quality update rings"
}

# ==============================================================================
# Ring — full configuration with both included and excluded assignments,
# hotpatch enabled, deferral, and paused state
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  depends_on = [time_sleep.wait_for_groups]

  policy_id          = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name       = "Production Ring - ${random_string.suffix.result}"
  description        = "Quality updates with 14-day deferral for the production audience"
  is_paused          = false
  deferral_in_days   = 14
  is_hotpatch_enabled = false

  included_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.included_group.id
      }
    ]
  }

  excluded_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.excluded_group.id
      }
    ]
  }
}
