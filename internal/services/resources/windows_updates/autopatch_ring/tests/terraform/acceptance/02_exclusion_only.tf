resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "excluded_group" {
  display_name     = "acc-test-ring-excluded-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-ring-excluded-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on      = [microsoft365_graph_beta_groups_group.excluded_group]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "base" {
  display_name = "acc-test-ring-policy-${random_string.suffix.result}"
  description  = "Acceptance test - ring base policy"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "test" {
  depends_on = [time_sleep.wait_for_groups]

  policy_id        = microsoft365_graph_beta_windows_updates_autopatch_policy.base.id
  display_name     = "Acc Test Ring 02 Exclusion Only"
  description      = "Acceptance test ring - exclusion assignments only"
  is_paused        = false
  deferral_in_days = 7

  excluded_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.excluded_group.id
      }
    ]
  }
}
