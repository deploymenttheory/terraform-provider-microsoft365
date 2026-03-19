# ==============================================================================
# Quality Update — Date-Driven Rollout
# ==============================================================================
# Deploys a cumulative quality (security) update using a date-driven gradual
# rollout. Devices are progressively offered the update until the specified
# end date, rather than rolling out a fixed number of devices per wave.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "quality_update" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2026-02-01T08:00:00Z"
      gradual_rollout = {
        end_date_time = "2026-03-01T08:00:00Z"
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 5
          action    = "pauseDeployment"
        }
      ]
    }
  }
}
