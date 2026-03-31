resource "microsoft365_graph_beta_windows_updates_update_policy" "test" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true
}
