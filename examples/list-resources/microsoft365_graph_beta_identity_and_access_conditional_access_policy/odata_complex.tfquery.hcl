# Complex query combining multiple conditions
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "complex" {
  provider = microsoft365
  config {
    odata_filter = "(state eq 'enabled' or state eq 'enabledForReportingButNotEnforced') and (contains(displayName, 'MFA') or contains(displayName, 'Admin'))"
  }
}
