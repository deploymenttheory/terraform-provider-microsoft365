resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members" "test" {
  audience_id = "3ac1bab3-c634-4377-8290-6c9b729dd9a1"
  member_type = "azureADDevice"

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
}
