# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Connected Organization - Multiple domain identity sources
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "acc-test-connected-organization-multi-${random_string.suffix.result}"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "wingtiptoys.com"
      domain_name  = "wingtiptoys.com"
    },
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "tailspintoys.com"
      domain_name  = "tailspintoys.com"
    }
  ]
}
