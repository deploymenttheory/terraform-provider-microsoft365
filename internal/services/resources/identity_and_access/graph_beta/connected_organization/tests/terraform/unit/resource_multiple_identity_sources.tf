resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-multiple-identity-sources"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "example.net"
      domain_name  = "example.net"
    },
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "example.org"
      domain_name  = "example.org"
    }
  ]
}
