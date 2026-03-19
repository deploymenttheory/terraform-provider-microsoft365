# ==============================================================================
# Feature Update — Rollback Signal with Alert Error Action
# ==============================================================================
# Deploys a feature update with a gradual rollout schedule and a monitoring
# rule that raises an alert (rather than pausing) when the rollback threshold
# is reached. Use this when you want visibility into rollback rates without
# automatically stopping the deployment.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "rollback_alert" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 100
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 10
          action    = "alertError"
        }
      ]
    }
  }
}
