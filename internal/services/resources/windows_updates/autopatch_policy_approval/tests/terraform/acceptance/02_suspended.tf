data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "base" {
  display_name = "acc-test-policy-approval-01"
  description  = "Acceptance test - policy approval base policy"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy_approval" "test" {
  policy_id        = microsoft365_graph_beta_windows_updates_autopatch_policy.base.id
  catalog_entry_id = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
  status           = "suspended"
}
