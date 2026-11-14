# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Connected Organization - Minimal with a different identity source (forces replacement)
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "acc-test-connected-organization-${random_string.suffix.result}"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "adatum.com"
      domain_name  = "adatum.com"
    }
  ]
}
