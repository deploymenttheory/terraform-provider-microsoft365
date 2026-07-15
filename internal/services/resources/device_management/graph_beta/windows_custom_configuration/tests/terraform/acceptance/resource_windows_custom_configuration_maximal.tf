resource "random_string" "custom_config_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "custom_config_group_1" {
  display_name     = "acc-test-windows-custom-config-group-1-${random_string.custom_config_suffix.result}"
  mail_nickname    = "acc-test-custom-config-1-${random_string.custom_config_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "custom_config_group_2" {
  display_name     = "acc-test-windows-custom-config-group-2-${random_string.custom_config_suffix.result}"
  mail_nickname    = "acc-test-custom-config-2-${random_string.custom_config_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "custom_configuration_example" {
  display_name = "acc-test-windows-custom-config-${random_string.custom_config_suffix.result}"
  description  = "Example Windows custom configuration profile using OMA-URI settings"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ADMX Ingest"
      description  = "https://code.visualstudio.com/docs/enterprise/policies#_windows-group-policies"
      oma_uri      = "./Device/Vendor/MSFT/Policy/ConfigOperations/ADMXInstall/VSCode/Policy/VSCodeADMX"
      value        = <<-EOT
        <?xml version="1.0" encoding="utf-8"?>
        <policyDefinitions revision="1.1" schemaVersion="1.0">
          <policyNamespaces>
            <target prefix="VSCode" namespace="Microsoft.Policies.VSCode" />
          </policyNamespaces>
          <resources minRequiredRevision="1.0" />
          <supportedOn>
            <definitions>
              <definition name="Supported_1_67" displayName="$(string.Supported_1_67)" />
            </definitions>
          </supportedOn>
          <categories>
            <category displayName="$(string.Application)" name="Application" />
            <category displayName="$(string.Category_updateConfigurationTitle)" name="updateConfigurationTitle"><parentCategory ref="Application" /></category>
          </categories>
          <policies>
            <policy name="UpdateMode" class="Both" displayName="$(string.UpdateMode)" explainText="$(string.UpdateMode_updateMode)" key="Software\Policies\Microsoft\VSCode" presentation="$(presentation.UpdateMode)">
              <parentCategory ref="updateConfigurationTitle" />
              <supportedOn ref="Supported_1_67" />
              <elements>
                <enum id="UpdateMode" valueName="UpdateMode">
                  <item displayName="$(string.UpdateMode_none)"><value><string>none</string></value></item>
                  <item displayName="$(string.UpdateMode_manual)"><value><string>manual</string></value></item>
                </enum>
              </elements>
            </policy>
          </policies>
        </policyDefinitions>
      EOT
    },
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "UpdateMode"
      description  = "Set VSCode update mode to manual"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~updateConfigurationTitle/UpdateMode"
      value        = "<enabled/>\n<data id=\"UpdateMode\" value=\"manual\"/>"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "ExtensionsAutoUpdateDelay"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdateDelay"
      value        = "30"
    }
  ]

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.custom_config_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.custom_config_group_2.id
    }
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
