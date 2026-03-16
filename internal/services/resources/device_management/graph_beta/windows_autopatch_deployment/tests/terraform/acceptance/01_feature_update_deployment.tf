data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P1W"
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
