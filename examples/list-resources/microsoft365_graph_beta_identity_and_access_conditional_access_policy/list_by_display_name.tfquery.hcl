# List policies with "MFA" in the display name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "mfa_policies" {
  provider = microsoft365
  config {
    display_name_filter = "MFA"
  }
}
