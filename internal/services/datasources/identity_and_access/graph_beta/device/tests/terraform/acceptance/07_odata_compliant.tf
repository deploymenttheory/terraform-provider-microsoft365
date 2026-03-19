data "microsoft365_graph_beta_identity_and_access_device" "test" {
  odata_query = "isCompliant eq true"
}
