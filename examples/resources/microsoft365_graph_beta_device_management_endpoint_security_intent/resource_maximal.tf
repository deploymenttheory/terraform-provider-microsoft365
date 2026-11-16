# BitLocker disk encryption intent with managed settings and assignments
resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "bitlocker_maximal" {
  template_id        = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name       = "BitLocker - encrypt devices"
  description        = "Require BitLocker and rotate recovery passwords for Entra joined devices"
  role_scope_tag_ids = ["0"]

  settings = [
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerEncryptDevice"
      setting_type  = "boolean"
      value_json    = jsonencode(true)
    },
    {
      definition_id = "deviceConfiguration--windows10EndpointProtectionConfiguration_bitLockerRecoveryPasswordRotation"
      setting_type  = "string"
      value_json    = jsonencode("enabledForAzureAd")
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
