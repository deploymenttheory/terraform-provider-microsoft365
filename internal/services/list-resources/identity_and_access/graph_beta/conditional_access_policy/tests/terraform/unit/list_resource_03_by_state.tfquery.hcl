provider "microsoft365" {}

# List Conditional Access policies filtered by state
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "by_state" {
  provider = microsoft365
  config {
    state_filter = "enabled"
  }
}
