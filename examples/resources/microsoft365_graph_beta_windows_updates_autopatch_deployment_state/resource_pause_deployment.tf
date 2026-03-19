# ==============================================================================
# Pause a Deployment
# ==============================================================================
# Pauses an active autopatch deployment. The deployment must already exist,
# managed by the microsoft365_graph_beta_windows_updates_autopatch_deployment
# resource. Pausing stops the deployment from offering the update to additional
# devices while preserving progress on devices already receiving it.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "example" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "paused" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "paused"
}
