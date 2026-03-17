# ==============================================================================
# Approve a Quality Update
# ==============================================================================
# Approves a specific quality update catalog entry within a Windows Update
# policy. The policy must already exist. Once approved, the policy will deploy
# the update to devices managed by the policy audience.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "example" {
  display_name = "My Quality Update Policy"
  description  = "Policy for approving quality updates"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy_approval" "approved" {
  policy_id       = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  catalog_entry_id = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
  status          = "approved"
}
