resource "microsoft365_graph_beta_applications_oauth2_permission_grant" "test" {
  client_id    = "%[1]s"
  resource_id  = "%[2]s"
  consent_type = "AllPrincipals"
  scopes       = ["Probe.Read"]
}
