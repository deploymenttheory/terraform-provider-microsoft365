resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "test" {
  policy_id        = "983f03cd-03cd-983f-cd03-3f98cd033f98"
  display_name     = "Test Ring"
  description      = "A test ring for unit tests"
  is_paused        = false
  deferral_in_days = 7

  included_group_assignment = {
    assignments = [
      {
        group_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
      }
    ]
  }

  excluded_group_assignment = {
    assignments = [
      {
        group_id = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
      }
    ]
  }
}
