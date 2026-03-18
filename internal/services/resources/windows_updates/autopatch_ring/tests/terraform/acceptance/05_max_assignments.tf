resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "included_group_1" {
  display_name     = "acc-test-ring-inc-1-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-ring-inc-1-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "included_group_2" {
  display_name     = "acc-test-ring-inc-2-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-ring-inc-2-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.included_group_1,
    microsoft365_graph_beta_groups_group.included_group_2,
  ]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "base" {
  display_name = "acc-test-ring-policy-${random_string.suffix.result}"
  description  = "Acceptance test - ring base policy"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "test" {
  depends_on = [time_sleep.wait_for_groups]

  policy_id    = microsoft365_graph_beta_windows_updates_autopatch_policy.base.id
  display_name = "Acc Test Ring 05 Max Assignments"
  description  = "Acceptance test ring - maximum assignments (2 groups)"
  is_paused    = false
  deferral_in_days = 7

  included_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.included_group_1.id
      },
      {
        group_id = microsoft365_graph_beta_groups_group.included_group_2.id
      }
    ]
  }

}
