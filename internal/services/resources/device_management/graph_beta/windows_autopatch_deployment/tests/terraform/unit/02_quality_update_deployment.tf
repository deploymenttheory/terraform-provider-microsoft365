resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = "q123456-quality-update-id"
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        end_date_time = "2024-02-01T10:00:00Z"
      }
    }
  }
}
