data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "isPrivileged eq true"
}

