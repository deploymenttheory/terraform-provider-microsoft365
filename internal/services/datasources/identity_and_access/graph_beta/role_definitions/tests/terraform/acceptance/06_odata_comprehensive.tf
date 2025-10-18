data "microsoft365_graph_beta_identity_and_access_role_definitions" "odata_comprehensive" {
  filter_type   = "odata"
  odata_filter  = "isBuiltIn eq true"
  odata_count   = true
  odata_orderby = "displayName"
}
