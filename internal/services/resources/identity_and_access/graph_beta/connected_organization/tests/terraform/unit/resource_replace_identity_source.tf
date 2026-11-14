resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-replace-identity-source"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "example.org"
      domain_name  = "example.org"
    }
  ]
}
