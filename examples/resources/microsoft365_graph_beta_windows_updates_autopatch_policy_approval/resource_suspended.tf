# ==============================================================================
# Suspend an Approval
# ==============================================================================
# Creates a policy approval in the suspended state, preventing the update from
# being deployed to devices. Use this when you want to register an approval but
# hold off on deployment until a later time by updating the status to "approved".

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "example" {
  display_name = "My Quality Update Policy"
  description  = "Policy for approving quality updates"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy_approval" "suspended" {
  policy_id       = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  catalog_entry_id = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
  status          = "suspended"
}
