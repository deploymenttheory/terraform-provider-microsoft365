# Legacy deviceConfigurations API. The current UI uses a configurationPolicies template.
# Windows Delivery Optimization.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_delivery_optimization" {
  display_name       = "Windows Delivery Optimization"
  description        = "Template configuration example"
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
    "backgroundDownloadFromHttpDelayInSeconds"                  = 30
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
