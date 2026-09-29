resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name = "invalid"
  settings     = "{\"@odata.type\": \"#microsoft.graph.iosCustomConfiguration\", \"displayName\": \"duplicate\"}"
}
