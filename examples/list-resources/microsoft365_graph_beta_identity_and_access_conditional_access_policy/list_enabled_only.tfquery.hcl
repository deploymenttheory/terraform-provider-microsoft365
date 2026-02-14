# List only enabled policies
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_only" {
  provider = microsoft365
  config {
    state_filter = "enabled"
  }
}
