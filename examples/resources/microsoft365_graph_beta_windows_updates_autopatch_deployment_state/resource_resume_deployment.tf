# ==============================================================================
# Resume a Deployment (set to none / active)
# ==============================================================================
# Sets the deployment state to "none", which means the deployment is active and
# offering updates to devices. Use this to resume a previously paused deployment
# or to explicitly set a deployment into its offering state.

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

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "active" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "none"
}
