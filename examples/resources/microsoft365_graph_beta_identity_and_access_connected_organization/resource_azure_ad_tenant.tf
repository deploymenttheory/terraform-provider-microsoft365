# Connected organization for a Microsoft Entra tenant, identified by tenant ID.
# The tenant ID of a domain can be found with tenantRelationship: findTenantInformationByDomainName.
resource "microsoft365_graph_beta_identity_and_access_connected_organization" "contoso" {
  display_name = "Contoso"
  description  = "Contoso Microsoft Entra tenant"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.azureActiveDirectoryTenant"
      display_name = "Contoso"
      tenant_id    = "00000000-0000-0000-0000-000000000000"
    }
  ]
}
