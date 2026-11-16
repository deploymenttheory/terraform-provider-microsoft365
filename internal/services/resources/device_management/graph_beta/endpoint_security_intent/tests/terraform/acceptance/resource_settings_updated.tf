resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Empty security group created for this test only. Assigning the intent to it
# targets no users or devices.
resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-esi-group-${random_string.suffix.result}"
  mail_nickname    = "acc-test-esi-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "test" {
  template_id        = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name       = "acc-test-esi-updated-${random_string.suffix.result}"
  description        = "Acceptance test BitLocker intent with updated settings"
  role_scope_tag_ids = ["0"]

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
    }
  ]

  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test.id
    }
  ]
}
