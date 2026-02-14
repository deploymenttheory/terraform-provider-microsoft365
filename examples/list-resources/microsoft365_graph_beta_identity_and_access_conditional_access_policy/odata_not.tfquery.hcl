# Find policies that are not disabled
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "not_disabled" {
  provider = microsoft365
  config {
    odata_filter = "state ne 'disabled'"
  }
}
