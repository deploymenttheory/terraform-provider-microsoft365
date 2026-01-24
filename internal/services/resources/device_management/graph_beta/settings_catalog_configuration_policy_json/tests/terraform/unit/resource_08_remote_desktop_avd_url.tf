resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "remote_desktop_avd_url" {
  name               = "Test Remote Desktop AVD URL Policy"
  description        = "Test policy for remote desktop AVD URL"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                    = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
          settingDefinitionId              = "user_vendor_msft_policy_config_remotedesktop_autosubscription"
          settingInstanceTemplateReference = null
          simpleSettingCollectionValue = [
            {
              "@odata.type"                 = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
              value                         = "https://rdweb.wvd.microsoft.com/api/arm/feeddiscovery"
              settingValueTemplateReference = null
            }
          ]
        }
      }
    ]
  })
}
