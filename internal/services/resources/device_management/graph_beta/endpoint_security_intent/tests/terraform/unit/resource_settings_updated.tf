resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "test" {
  template_id        = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name       = "unit-test-endpoint-security-intent-maximal"
  description        = "Unit test BitLocker intent"
  role_scope_tag_ids = ["0", "1"]

  settings = [
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerEncryptDevice"
      setting_type  = "boolean"
      value_json    = jsonencode(false)
    },
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerRecoveryPasswordRotation"
      setting_type  = "string"
      value_json    = jsonencode("enabledForAzureAdAndHybrid")
    },
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerAllowStandardUserEncryption"
      setting_type  = "boolean"
      value_json    = jsonencode(true)
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "33333333-3333-3333-3333-333333333333"
      filter_type = "include"
    }
  ]
}
