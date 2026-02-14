# List policies in report-only mode
list "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "report_only" {
  provider = microsoft365
  config {
    state_filter = "enabledForReportingButNotEnforced"
  }
}
