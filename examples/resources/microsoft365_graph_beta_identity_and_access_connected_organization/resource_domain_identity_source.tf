# Connected organization for an external domain.
# The domain must resolve publicly and can belong to only one connected organization in the tenant.
resource "microsoft365_graph_beta_identity_and_access_connected_organization" "fabrikam" {
  display_name = "Fabrikam"
  description  = "Partner organization allowed to request access packages"
  state        = "configured"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "fabrikam.com"
      domain_name  = "fabrikam.com"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
