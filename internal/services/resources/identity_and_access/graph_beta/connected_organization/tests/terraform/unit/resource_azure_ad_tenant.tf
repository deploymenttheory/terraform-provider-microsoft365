resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "unit-test-connected-organization-azure-ad-tenant"
  description  = "Connected organization for a Microsoft Entra tenant"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.azureActiveDirectoryTenant"
      display_name = "Contoso"
      tenant_id    = "11111111-1111-1111-1111-111111111111"
    }
  ]
}
