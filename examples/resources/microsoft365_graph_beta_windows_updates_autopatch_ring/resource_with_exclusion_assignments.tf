resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  policy_id    = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name = "Broad Ring"
  description  = "Quality updates deployed to all devices, excluding the VIP group"
  is_paused    = false
  deferral_in_days = 14

  excluded_group_assignment = {
    assignments = [
      {
        group_id = "00000000-0000-0000-0000-000000000002"
      }
    ]
  }
}
