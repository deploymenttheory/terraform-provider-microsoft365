resource "random_string" "notifications_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "notifications" {
  name               = "acc-test-05-notifications-${random_string.notifications_suffix.result}"
  description        = "Acceptance test policy for notifications settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "user_vendor_msft_policy_config_notifications_disallownotificationmirroring"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "user_vendor_msft_policy_config_notifications_disallownotificationmirroring_1"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })
}
