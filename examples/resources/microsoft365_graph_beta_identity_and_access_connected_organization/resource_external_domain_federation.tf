# Connected organization for a domain federated through a SAML/WS-Fed identity provider.
# The federation must already be configured in this tenant
# (directory/federationConfigurations); otherwise Microsoft Graph returns IdentityProviderNotFound.
resource "microsoft365_graph_beta_identity_and_access_connected_organization" "woodgrove" {
  display_name = "Woodgrove"
  state        = "configured"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.externalDomainFederation"
      display_name = "Woodgrove"
      domain_name  = "woodgrovebank.com"
      issuer_uri   = "https://idp.woodgrovebank.com/saml"
    }
  ]
}
