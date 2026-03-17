# ==============================================================================
# Feature Update — Ineligible Signal with Offer Fallback Action
# ==============================================================================
# Deploys a Windows 11 feature update with a monitoring rule that automatically
# offers Windows 10 22H2 as a fallback to devices that are ineligible for the
# Windows 11 update. This combination is only valid for Windows 11 feature
# update deployments.
#
# Note: the "ineligible" signal must always be paired with the "offerFallback"
# action, and no threshold is accepted for this combination — the fallback is
# offered to all ineligible devices unconditionally.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "ineligible_fallback" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    monitoring = {
      monitoring_rules = [
        {
          signal = "ineligible"
          action = "offerFallback"
        }
      ]
    }
  }
}
