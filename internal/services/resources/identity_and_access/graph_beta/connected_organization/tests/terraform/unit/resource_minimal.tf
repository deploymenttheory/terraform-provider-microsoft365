resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-minimal"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "example.com"
      domain_name  = "example.com"
    }
  ]
}
