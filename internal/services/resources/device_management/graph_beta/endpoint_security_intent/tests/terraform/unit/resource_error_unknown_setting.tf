resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "test" {
  template_id  = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name = "unit-test-endpoint-security-intent-unknown-setting"

  settings = [
    {
      definition_id = "deviceConfiguration--notATemplateSetting"
      setting_type  = "boolean"
      value_json    = jsonencode(true)
    }
  ]
}
