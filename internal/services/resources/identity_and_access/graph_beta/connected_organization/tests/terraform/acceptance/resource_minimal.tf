# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Connected Organization - Minimal (domain identity source)
# ==============================================================================

# Graph only accepts a domain identity source for a domain that resolves publicly
# (it returns 404 IdentityProviderNotFound otherwise) and allows each domain in only
# one connected organization per tenant, so acceptance fixtures use fixed public
# sample domains rather than random subdomains.
resource "microsoft365_graph_beta_identity_and_access_connected_organization" "test" {
  display_name = "acc-test-connected-organization-${random_string.suffix.result}"

  identity_sources = [
    {
      odata_type   = "#microsoft.graph.domainIdentitySource"
      display_name = "fabrikam.com"
      domain_name  = "fabrikam.com"
    }
  ]
}
