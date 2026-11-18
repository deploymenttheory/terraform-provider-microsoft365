# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Connected Organization - Maximal (domain identity source)
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "acc-test-connected-organization-maximal-${random_string.suffix.result}"
  description  = "Connected organization managed by Terraform acceptance tests"
  state        = "proposed"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "${random_string.suffix.result}.192-0-2-1.sslip.io"
      domain_name  = "${random_string.suffix.result}.192-0-2-1.sslip.io"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
