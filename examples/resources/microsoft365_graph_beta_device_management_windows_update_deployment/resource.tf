# Example 1: Feature update deployment with rate-driven rollout
resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "feature_update" {
  content = {
    catalog_entry_id   = "{catalog_entry_id}"
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2024-03-01T10:00:00Z"
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

# Example 2: Quality update deployment with date-driven rollout
resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "quality_update" {
  content = {
    catalog_entry_id   = "{catalog_entry_id}"
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        end_date_time = "2024-03-15T10:00:00Z"
      }
    }
  }
}

# Example 3: Minimal deployment
resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "minimal" {
  content = {
    catalog_entry_id   = "{catalog_entry_id}"
    catalog_entry_type = "featureUpdate"
  }
}

# Example 4: Paused deployment
resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "paused" {
  content = {
    catalog_entry_id   = "{catalog_entry_id}"
    catalog_entry_type = "featureUpdate"
  }

  state = {
    requested_value = "paused"
  }
}

# Example 5: Using data source to get catalog entry
data "microsoft365_graph_beta_device_management_windows_update_catalog_item" "latest_feature_update" {
  filter_type  = "odata"
  odata_filter = "isof('microsoft.graph.windowsUpdates.featureUpdateCatalogEntry')"
  odata_top    = 1
}

resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment" "from_datasource" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_device_management_windows_update_catalog_item.latest_feature_update.items[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P1D"
        devices_per_offer       = 50
      }
    }
  }
}
