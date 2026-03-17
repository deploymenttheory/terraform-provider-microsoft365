# ==============================================================================
# Minimal Deployment (no settings)
# ==============================================================================
# Creates a Windows Update deployment with only the required content block.
# No schedule or monitoring rules are configured. The deployment will be
# created in a pending state and can have settings added later by updating
# the resource. Once settings are applied they cannot be modified in place —
# changes will require the resource to be replaced.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "minimal" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }
}
