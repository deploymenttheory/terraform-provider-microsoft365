resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "test" {
  template_id  = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name = "unit-test-endpoint-security-intent-invalid-json"

  settings = [
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerRecoveryPasswordRotation"
      setting_type  = "string"
      value_json    = "enabledForAzureAd"
    }
  ]
}
