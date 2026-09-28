resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "unit-test-aosp-restrictions"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                    = "#microsoft.graph.aospDeviceOwnerDeviceConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"     = null,
    "deviceManagementApplicabilityRuleOsVersion"     = null,
    "deviceManagementApplicabilityRuleDeviceMode"    = null,
    "appsBlockInstallFromUnknownSources"             = null,
    "bluetoothBlocked"                               = null,
    "bluetoothBlockConfiguration"                    = null,
    "cameraBlocked"                                  = true,
    "factoryResetBlocked"                            = null,
    "passwordMinimumLength"                          = null,
    "passwordMinutesOfInactivityBeforeScreenTimeout" = null,
    "passwordRequiredType"                           = null,
    "passwordSignInFailureCountBeforeFactoryReset"   = null,
    "screenCaptureBlocked"                           = null,
    "securityAllowDebuggingFeatures"                 = null,
    "storageBlockExternalMedia"                      = null,
    "storageBlockUsbFileTransfer"                    = null,
    "wifiBlockEditConfigurations"                    = null
  })
}
