# Find users with either of two job titles
list "microsoft365_graph_beta_users_user" "multiple_titles" {
  provider = microsoft365
  config {
    odata_filter = "jobTitle eq 'Manager' or jobTitle eq 'Director'"
  }
}
