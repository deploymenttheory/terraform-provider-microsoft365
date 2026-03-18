resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "base" {
  display_name = "acc-test-ring-policy-${random_string.suffix.result}"
  description  = "Acceptance test - ring base policy"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "test" {
  policy_id    = microsoft365_graph_beta_windows_updates_autopatch_policy.base.id
  display_name = "Acc Test Ring 03 No Assignments"
  description  = "Acceptance test ring - no group assignments"
  is_paused    = false
  deferral_in_days = 0

}
