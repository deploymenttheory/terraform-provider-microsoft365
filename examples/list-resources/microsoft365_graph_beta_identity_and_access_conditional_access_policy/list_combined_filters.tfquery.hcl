# List enabled policies with "Admin" in the name
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "enabled_admin_policies" {
  provider = microsoft365
  config {
    display_name_filter = "Admin"
    state_filter        = "enabled"
  }
}
