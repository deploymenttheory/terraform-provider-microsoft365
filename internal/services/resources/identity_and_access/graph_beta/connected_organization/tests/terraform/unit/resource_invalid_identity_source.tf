resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-invalid"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "example.com"
      tenant_id    = "11111111-1111-1111-1111-111111111111"
    }
  ]
}
