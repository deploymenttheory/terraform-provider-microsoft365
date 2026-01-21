resource "random_string" "cryptography_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "cryptography" {
  name               = "acc-test-04-cryptography-${random_string.cryptography_suffix.result}"
  description        = "Acceptance test policy for cryptography settings"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"                        = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId                  = "device_vendor_msft_policy_config_cryptography_allowfipsalgorithmpolicy"
          settingInstanceTemplateReference     = null
          auditRuleInformation                 = null
          choiceSettingValue = {
            value                            = "device_vendor_msft_policy_config_cryptography_allowfipsalgorithmpolicy_1"
            settingValueTemplateReference    = null
            children                         = []
          }
        }
      }
    ]
  })
}
