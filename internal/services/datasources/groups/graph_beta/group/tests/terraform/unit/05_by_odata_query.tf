data "microsoft365_graph_beta_groups_group" "test" {
  odata_query = "displayName eq 'IT Security Team' and securityEnabled eq true"
}
