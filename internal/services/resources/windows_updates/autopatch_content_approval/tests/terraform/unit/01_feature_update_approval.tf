resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "test" {
  update_policy_id    = "983f03cd-03cd-983f-cd03-3f98cd033f98"
  catalog_entry_id    = "c1dec151-c151-c1de-51c1-dec151c1dec1"
  catalog_entry_type  = "featureUpdate"

  deployment_settings = {
    schedule = {
      start_date_time = "2026-03-10T00:00:00Z"
      gradual_rollout = {
        end_date_time = "2026-03-20T00:00:00Z"
      }
    }
  }
}
