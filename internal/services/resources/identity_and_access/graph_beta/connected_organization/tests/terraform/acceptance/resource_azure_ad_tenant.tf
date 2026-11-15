# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Connected Organization - Microsoft Entra tenant identity source
# ==============================================================================

# The tenant ID is the public Microsoft corporate tenant. A connected organization
# only references it; nothing is changed in that tenant.
resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "acc-test-connected-organization-tenant-${random_string.suffix.result}"
  description  = "Connected organization for a Microsoft Entra tenant"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.azureActiveDirectoryTenant"
      display_name = "acc-test-tenant-${random_string.suffix.result}"
      tenant_id    = "72f988bf-86f1-41af-91ab-2d7cd011db47"
    }
  ]
}
