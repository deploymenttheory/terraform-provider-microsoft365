resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members" "test" {
  audience_id = "445aeeb2-0ccd-458d-8f0d-101c5678eff2"
  member_type = "azureADDevice"

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  exclusions = [
    "00000000-0000-0000-0000-000000000003"
  ]
}
