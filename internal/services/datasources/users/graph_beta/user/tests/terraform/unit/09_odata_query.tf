data "microsoft365_graph_beta_users_user" "test" {
  odata_query = "accountEnabled eq true and userType eq 'Member'"
}
