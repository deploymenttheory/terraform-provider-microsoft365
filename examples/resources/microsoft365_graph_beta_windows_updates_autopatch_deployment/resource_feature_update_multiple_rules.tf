# ==============================================================================
# Feature Update — Multiple Monitoring Rules
# ==============================================================================
# Deploys a feature update with two monitoring rules configured simultaneously:
#   1. Pause the deployment if more than 5% of devices roll back.
#   2. Offer Windows 10 22H2 as a fallback to any device ineligible for Windows 11.
#
# Multiple monitoring rules allow you to respond to different failure signals
# independently within a single deployment.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "multiple_rules" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P14D"
        devices_per_offer       = 200
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 5
          action    = "pauseDeployment"
        },
        {
          signal = "ineligible"
          action = "offerFallback"
        }
      ]
    }
  }
}
