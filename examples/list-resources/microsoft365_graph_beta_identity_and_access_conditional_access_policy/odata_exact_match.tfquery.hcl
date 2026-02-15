# Find specific policy by exact display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "exact_name" {
  provider = microsoft365
  config {
    odata_filter = "displayName eq 'Require MFA for Administrators'"
  }
}
