// Example: Reuseable Policy Setting Resource

resource "microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings" "example" {
  display_name = "epm certificate"
  description  = "Endpoint Privilege Management supports using reusable settings groups to manage the certificates in place of adding that certificate directly to an elevation rule"
  settings = jsonencode({
    "settings" : [
      {
        "id" : "0",
        "settingInstance" : {
          "@odata.type" : "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance",
          "settingDefinitionId" : "device_vendor_msft_policy_privilegemanagement_reusablesettings_certificatefile",
          "simpleSettingValue" : {
            "@odata.type" : "#microsoft.graph.deviceManagementConfigurationStringSettingValue",
            "value" : "-----BEGIN CERTIFICATE-----\r\nMIIDXzCCAkegAwIBAgtcpH\r\nWD9f\r\n-----END CERTIFICATE-----\r\n"
          }
        }
      }
    ]
  })

  # Timeouts configuration (optional)
  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}