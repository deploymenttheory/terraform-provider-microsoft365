# Find policies that are either enabled or in report-only mode
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "active_policies" {
  provider = microsoft365
  config {
    odata_filter = "state eq 'enabled' or state eq 'enabledForReportingButNotEnforced'"
  }
}
