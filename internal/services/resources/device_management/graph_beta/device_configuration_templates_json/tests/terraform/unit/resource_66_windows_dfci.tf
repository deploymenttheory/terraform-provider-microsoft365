resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-windows-dfci"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10DeviceFirmwareConfigurationInterface"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "changeUefiSettingsPermission"                = "notConfiguredOnly"
    "virtualizationOfCpuAndIO"                    = "notConfigured"
    "cameras"                                     = "notConfigured"
    "microphonesAndSpeakers"                      = "notConfigured"
    "radios"                                      = "notConfigured"
    "bootFromExternalMedia"                       = "notConfigured"
    "bootFromBuiltInNetworkAdapters"              = "notConfigured"
    "windowsPlatformBinaryTable"                  = "notConfigured"
    "simultaneousMultiThreading"                  = "notConfigured"
    "frontCamera"                                 = "notConfigured"
    "rearCamera"                                  = "notConfigured"
    "infraredCamera"                              = "notConfigured"
    "microphone"                                  = "notConfigured"
    "bluetooth"                                   = "notConfigured"
    "wirelessWideAreaNetwork"                     = "notConfigured"
    "nearFieldCommunication"                      = "notConfigured"
    "wiFi"                                        = "notConfigured"
    "usbTypeAPort"                                = "notConfigured"
    "sdCard"                                      = "notConfigured"
    "wakeOnLAN"                                   = "notConfigured"
    "wakeOnPower"                                 = "notConfigured"
  })
}
