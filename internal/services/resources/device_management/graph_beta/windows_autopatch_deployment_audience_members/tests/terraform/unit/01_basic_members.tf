resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience_members" "test" {
  audience_id = "test-audience-id-001"
  member_type = "azureADDevice"

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
}
