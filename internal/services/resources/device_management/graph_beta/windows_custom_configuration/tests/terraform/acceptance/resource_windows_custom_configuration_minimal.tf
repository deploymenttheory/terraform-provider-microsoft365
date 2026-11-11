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

resource "time_sleep" "group_replication" {
  depends_on = [
    microsoft365_graph_beta_groups_group.custom_config_group_1,
    microsoft365_graph_beta_groups_group.custom_config_group_2,
  ]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "custom_configuration_example" {
  depends_on   = [time_sleep.group_replication]
  display_name = "acc-test-windows-custom-config-${random_string.custom_config_suffix.result}"
  description  = "Minimal Windows custom configuration profile"

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
    }
  ]

  role_scope_tag_ids = ["0"]


  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
