provider "microsoft365" {}

# List Conditional Access policies filtered by display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "by_display_name" {
  provider = microsoft365
  config {
    display_name_filter = "MFA"
  }
}
