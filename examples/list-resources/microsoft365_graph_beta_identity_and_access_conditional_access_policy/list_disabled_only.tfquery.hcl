# List only disabled policies
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "disabled_only" {
  provider = microsoft365
  config {
    state_filter = "disabled"
  }
}
