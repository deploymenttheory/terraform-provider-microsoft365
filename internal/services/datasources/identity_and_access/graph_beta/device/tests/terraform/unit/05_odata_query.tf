# Lookup devices with OData query filter
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  odata_query = "operatingSystem eq 'Windows' and isCompliant eq true"
}
