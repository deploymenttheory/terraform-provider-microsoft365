resource "microsoft365_graph_beta_applications_oauth2_permission_grant" "test" {
  client_id    = "11111111-1111-4111-8111-111111111111"
  resource_id  = "22222222-2222-4222-8222-222222222222"
  consent_type = "Principal"
  principal_id = "33333333-3333-4333-8333-333333333333"
  scopes       = ["Probe.Read"]
}
