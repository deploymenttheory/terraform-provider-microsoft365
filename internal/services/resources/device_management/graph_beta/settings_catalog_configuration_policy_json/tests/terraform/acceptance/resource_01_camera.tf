resource "random_string" "camera_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "camera" {
  name               = "acc-test-01-camera-${random_string.camera_suffix.result}"
  description        = "Acceptance test policy for camera settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId              = "device_vendor_msft_policy_config_camera_allowcamera"
          settingInstanceTemplateReference = null
          auditRuleInformation             = null
          choiceSettingValue = {
            value                         = "device_vendor_msft_policy_config_camera_allowcamera_0"
            settingValueTemplateReference = null
            children                      = []
          }
        }
      }
    ]
  })
}
