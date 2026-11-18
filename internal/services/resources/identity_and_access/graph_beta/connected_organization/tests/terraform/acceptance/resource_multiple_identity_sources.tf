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
      display_name = "first-${random_string.suffix.result}.192-0-2-1.sslip.io"
      domain_name  = "first-${random_string.suffix.result}.192-0-2-1.sslip.io"
    },
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "second-${random_string.suffix.result}.192-0-2-1.sslip.io"
      domain_name  = "second-${random_string.suffix.result}.192-0-2-1.sslip.io"
    }
  ]
}
