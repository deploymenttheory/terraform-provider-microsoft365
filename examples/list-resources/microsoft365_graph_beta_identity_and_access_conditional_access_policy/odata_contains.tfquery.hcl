# Find policies with specific text in display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "contains_guest" {
  provider = microsoft365
  config {
    odata_filter = "contains(displayName, 'Guest')"
  }
}
