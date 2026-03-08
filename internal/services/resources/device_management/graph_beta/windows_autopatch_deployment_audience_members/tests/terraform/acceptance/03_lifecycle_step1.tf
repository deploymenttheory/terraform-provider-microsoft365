# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - Groups for audience members
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "member_group_1" {
  display_name     = "acc-test-lifecycle-member-1-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-lc-member-1-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_2" {
  display_name     = "acc-test-lifecycle-member-2-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-lc-member-2-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_3" {
  display_name     = "acc-test-lifecycle-member-3-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-lc-member-3-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

# ==============================================================================
# Time Sleep - Wait for groups to propagate
# ==============================================================================

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.member_group_1,
    microsoft365_graph_beta_groups_group.member_group_2,
    microsoft365_graph_beta_groups_group.member_group_3
  ]

  create_duration = "30s"
}

# ==============================================================================
# Deployment Audience (Container)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience" "test" {
  depends_on = [time_sleep.wait_for_groups]
}

# ==============================================================================
# Time Sleep - Wait for audience to propagate
# ==============================================================================

resource "time_sleep" "wait_for_audience" {
  depends_on = [
    microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience.test
  ]

  create_duration = "30s"
}

# ==============================================================================
# Deployment Audience Members - Step 1: Initial 2 members
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience_members" "test" {
  depends_on = [time_sleep.wait_for_audience]

  audience_id = microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience.test.id
  member_type = "updatableAssetGroup"

  members = [
    microsoft365_graph_beta_groups_group.member_group_1.id,
    microsoft365_graph_beta_groups_group.member_group_2.id
  ]
}
