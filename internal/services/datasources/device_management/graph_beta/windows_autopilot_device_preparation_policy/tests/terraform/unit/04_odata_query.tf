# Lookup Windows Autopilot device preparation policies with a custom OData query
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  odata_query = "isAssigned eq true"
}
