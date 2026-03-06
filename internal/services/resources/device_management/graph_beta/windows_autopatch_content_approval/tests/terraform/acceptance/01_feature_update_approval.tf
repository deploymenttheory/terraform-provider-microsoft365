resource "microsoft365_graph_beta_device_management_windows_autopatch_content_approval" "test" {
  update_policy_id    = "45a01ef3-fb4b-8c1d-2428-1f060464033c"
  catalog_entry_id    = "c1dec151-c151-c1de-51c1-dec151c1dec1"
  catalog_entry_type  = "featureUpdate"

  deployment_settings = {
    schedule = {
      start_date_time = "2026-04-01T00:00:00Z"
      gradual_rollout = {
        end_date_time = "2026-04-15T00:00:00Z"
      }
    }
  }
}
