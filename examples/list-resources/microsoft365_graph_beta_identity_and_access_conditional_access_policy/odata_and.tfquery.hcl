# Find enabled policies with specific display name pattern
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_with_pattern" {
  provider = microsoft365
  config {
    odata_filter = "state eq 'enabled' and contains(displayName, 'Baseline')"
  }
}
