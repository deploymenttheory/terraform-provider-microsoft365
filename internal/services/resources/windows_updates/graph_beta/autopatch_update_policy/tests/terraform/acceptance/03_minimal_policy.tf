resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "test" {
  audience_id        = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id
  compliance_changes = true
}
