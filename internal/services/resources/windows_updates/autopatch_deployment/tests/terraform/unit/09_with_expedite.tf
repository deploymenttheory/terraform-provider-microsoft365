resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = "d0c03fbb-43b9-4dff-840b-974ef227384d"
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2024-01-15T10:00:00Z"
    }
    expedite = {
      is_expedited      = true
      is_readiness_test = false
    }
  }
}
