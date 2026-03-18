resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "test" {
  update_policy_id   = "983f03cd-03cd-983f-cd03-3f98cd033f98"
  catalog_entry_id   = "d0c03fbb-43b9-4dff-840b-974ef227384d"
  catalog_entry_type = "qualityUpdate"

  deployment_settings = {
    schedule = {
      start_date_time = "2026-03-11T00:00:00Z"
    }
  }
}
