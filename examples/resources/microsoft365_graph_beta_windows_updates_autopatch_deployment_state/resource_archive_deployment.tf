# ==============================================================================
# Archive a Deployment
# ==============================================================================
# Archives a deployment, permanently stopping it from offering updates to any
# further devices. Archived deployments are read-only and cannot be resumed.
# Use this when a deployment is complete or no longer needed.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "example" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
    catalog_entry_type = "qualityUpdate"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "archived" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "archived"
}
