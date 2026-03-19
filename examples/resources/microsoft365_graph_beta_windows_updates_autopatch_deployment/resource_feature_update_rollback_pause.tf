# ==============================================================================
# Feature Update — Rollback Signal with Pause Deployment Action
# ==============================================================================
# Deploys a feature update with a rate-driven gradual rollout schedule and a
# monitoring rule that pauses the deployment if too many devices roll back.
# This is the most common monitoring configuration for feature update deployments.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "rollback_pause" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2025-02-01T08:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 100
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
