data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "test_policy_02" {
}

resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "test" {
  update_policy_id   = microsoft365_graph_beta_windows_updates_update_policy.test_policy_02.id
  catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
  catalog_entry_type = "featureUpdate"

  deployment_settings = {
    schedule = {
      start_date_time = "2026-05-01T00:00:00Z"
      gradual_rollout = {
        end_date_time = "2026-05-15T00:00:00Z"
      }
    }
  }
}
