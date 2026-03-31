resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = "f341705b-0b15-4ce3-aaf2-6a1681d78606"
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2024-01-15T10:00:00Z"
    }
    content_applicability = {
      safeguard = {
        disabled_safeguard_profiles = [
          {
            category = "likelyIssues"
          }
        ]
      }
    }
  }
}
