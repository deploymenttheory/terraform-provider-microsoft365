resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  policy_id        = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name     = "Pilot Ring"
  description      = "Quality updates deployed to the pilot group after a 7-day deferral"
  is_paused        = false
  deferral_in_days = 7

  included_group_assignment = {
    assignments = [
      {
        group_id = "00000000-0000-0000-0000-000000000001"
      }
    ]
  }
}
