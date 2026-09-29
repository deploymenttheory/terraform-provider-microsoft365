resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-windows-delivery-optimization-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                               = "#microsoft.graph.windowsDeliveryOptimizationConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"                = null
    "deviceManagementApplicabilityRuleOsVersion"                = null
    "deviceManagementApplicabilityRuleDeviceMode"               = null
    "deliveryOptimizationMode"                                  = "userDefined"
    "restrictPeerSelectionBy"                                   = "notConfigured"
    "groupIdSource"                                             = null
    "bandwidthMode"                                             = null
    "backgroundDownloadFromHttpDelayInSeconds"                  = 0
    "foregroundDownloadFromHttpDelayInSeconds"                  = null
    "minimumRamAllowedToPeerInGigabytes"                        = null
    "minimumDiskSizeAllowedToPeerInGigabytes"                   = null
    "minimumFileSizeToCacheInMegabytes"                         = null
    "minimumBatteryPercentageAllowedToUpload"                   = null
    "modifyCacheLocation"                                       = null
    "maximumCacheAgeInDays"                                     = null
    "maximumCacheSize"                                          = null
    "vpnPeerCaching"                                            = "notConfigured"
    "cacheServerHostNames"                                      = []
    "cacheServerForegroundDownloadFallbackToHttpDelayInSeconds" = 0
    "cacheServerBackgroundDownloadFallbackToHttpDelayInSeconds" = 0
  })
}
