resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-external-domain-federation"
  description  = "Connected organization for a SAML/WS-Fed federated domain"
  state        = "configured"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.externalDomainFederation"
      display_name = "Fabrikam"
      domain_name  = "fabrikam.example"
      issuer_uri   = "https://idp.fabrikam.example/saml"
    }
  ]
}
