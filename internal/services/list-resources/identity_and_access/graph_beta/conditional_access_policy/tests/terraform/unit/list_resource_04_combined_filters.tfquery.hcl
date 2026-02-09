provider "microsoft365" {}

# List Conditional Access policies with combined filters
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "combined" {
  provider = microsoft365
  config {
    display_name_filter = "MFA"
    state_filter        = "enabled"
  }
}
