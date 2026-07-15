# Example: Windows Custom Configuration (OMA-URI) profile
# This example ingests a third party ADMX file (VSCode) and configures one of its policies.
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "vscode_policy" {
  display_name = "VSCode Policy"
  description  = "Configure Visual Studio Code enterprise policies via ADMX ingestion"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ADMX"
      description  = "https://code.visualstudio.com/docs/enterprise/policies#_windows-group-policies"
      oma_uri      = "./Device/Vendor/MSFT/Policy/ConfigOperations/ADMXInstall/VSCode/Policy/VSCodeADMX"
      value        = file("${path.module}/vscode.admx")
    },
    {
      odata_type   = "#microsoft.graph.omaSettingString"
      display_name = "ExtensionsAutoUpdate"
      description  = "Disable VSCode extension auto update"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/VSCode~Policy~Application~extensionsConfigurationTitle/ExtensionsAutoUpdate"
      value        = "<enabled/>\n<data id=\"ExtensionsAutoUpdate\" value=\"off\"/>"
    }
  ]

  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Example: other OMA setting value types
resource "microsoft365_graph_beta_device_management_windows_custom_configuration" "value_types" {
  display_name = "Windows Custom Configuration - value types"

  oma_settings = [
    {
      odata_type   = "#microsoft.graph.omaSettingInteger"
      display_name = "Integer setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/IntegerSetting"
      value        = "30"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingBoolean"
      display_name = "Boolean setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/BooleanSetting"
      value        = "true"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingFloatingPoint"
      display_name = "Floating point setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/FloatSetting"
      value        = "1.5"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingDateTime"
      display_name = "Date time setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/DateTimeSetting"
      value        = "2024-01-01T00:00:00Z"
    },
    {
      odata_type   = "#microsoft.graph.omaSettingBase64"
      display_name = "Base64 setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/Base64Setting"
      file_name    = "logo.png"
      value        = filebase64("${path.module}/logo.png")
    },
    {
      odata_type   = "#microsoft.graph.omaSettingStringXml"
      display_name = "String XML setting"
      oma_uri      = "./Device/Vendor/MSFT/Policy/Config/Example/XmlSetting"
      file_name    = "settings.xml"
      value        = file("${path.module}/settings.xml")
    }
  ]
}
