data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "isPrivileged eq true"
  odata_top     = 5
  odata_skip    = 0
  odata_count   = true
  odata_orderby = "displayName desc"
}

