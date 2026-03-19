resource "microsoft365_graph_beta_windows_updates_autopatch_operational_insights_connection" "example" {
  azure_resource_group_name = "my-resource-group"
  azure_subscription_id     = "12345678-1234-1234-1234-123456789012"
  workspace_name            = "my-log-analytics-workspace"
}
