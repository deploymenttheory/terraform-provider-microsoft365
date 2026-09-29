---
page_title: "microsoft365_graph_beta_device_management_device_configuration_templates_json Resource - terraform-provider-microsoft365"
subcategory: "Device Management"
description: |-
  Manages cross-platform Microsoft Intune device configuration templates through /deviceManagement/deviceConfigurations. Profile metadata uses standard Terraform attributes. Template-specific settings, including the root @odata.type, use JSON. Assignments are managed separately in this resource. This resource is separate from Settings Catalog (/configurationPolicies). Do not manage the same profile with both a typed resource and this JSON resource. See device configurations https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta.
---

# microsoft365_graph_beta_device_management_device_configuration_templates_json (Resource)

Manages cross-platform Microsoft Intune device configuration templates through `/deviceManagement/deviceConfigurations`. Profile metadata uses standard Terraform attributes. Template-specific settings, including the root `@odata.type`, use JSON. Assignments are managed separately in this resource. This resource is separate from Settings Catalog (`/configurationPolicies`). Do not manage the same profile with both a typed resource and this JSON resource. See [device configurations](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta).

## Microsoft Documentation

- [deviceConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta)
- [Get deviceConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-deviceconfiguration-get?view=graph-rest-beta)
- [Assign deviceConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-deviceconfiguration-assign?view=graph-rest-beta)
- [Get an OMA setting plaintext value](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-deviceconfiguration-getomasettingplaintextvalue?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**

- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementConfiguration.ReadWrite.All`

**Optional:**

- `Group.ReadWrite.All` (when creating the security groups included in the assignment examples)

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| Unreleased | Experimental | Initial release |

## Important Notes

- **Complete settings:** Include the complete writable settings JSON, including defaults, nulls, empty values, and the root `@odata.type`. Configure `display_name`, `description`, `role_scope_tag_ids`, and `assignments` as separate resource attributes.
- **Updates:** Supply explicit reset values for ordinary settings. Omit a certificate binding to remove its relationship. Changing the root profile type requires Terraform's `-replace` option.
- **XML and binary values:** Configure `omaSettingStringXml.value` as cleartext XML. Binary OMA values, certificate content, and custom Apple payloads retain their Base64 representation.
- **Example values:** Replace the sample certificates, service URLs, network names, application identifiers, and OMA-URIs with values for your environment. The certificate examples include their related profile resources in the same file.
- **Assignments:** The group examples create their security groups before the profile and wait for directory consistency. The all-targets example uses an empty custom payload. Profiles without assignments are not deployed to users or devices.
- **Import:** Wi-Fi pre-shared keys cannot be recovered on import and must be supplied again. Secret settings are stored in Terraform state.

## Example Usage

Each example below is a separate configuration derived from the acceptance scenarios. Select the file for the platform and profile you want to manage.

### Android

#### General Device Configuration

```terraform
# Android General Device Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_general_configuration" {
  display_name       = "Android General Device Configuration"
  description        = "Android general device configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                              = "#microsoft.graph.androidDeviceOwnerGeneralDeviceConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"               = null,
    "deviceManagementApplicabilityRuleOsVersion"               = null,
    "deviceManagementApplicabilityRuleDeviceMode"              = null,
    "accountsBlockModification"                                = null,
    "appsAllowInstallFromUnknownSources"                       = null,
    "appsAutoUpdatePolicy"                                     = null,
    "appsDefaultPermissionPolicy"                              = null,
    "appsRecommendSkippingFirstUseHints"                       = null,
    "bluetoothBlockConfiguration"                              = null,
    "bluetoothBlockContactSharing"                             = null,
    "cameraBlocked"                                            = true,
    "cellularBlockWiFiTethering"                               = null,
    "certificateCredentialConfigurationDisabled"               = null,
    "crossProfilePoliciesAllowCopyPaste"                       = null,
    "crossProfilePoliciesAllowDataSharing"                     = null,
    "crossProfilePoliciesShowWorkContactsInPersonalProfile"    = null,
    "microsoftLauncherConfigurationEnabled"                    = null,
    "microsoftLauncherCustomWallpaperEnabled"                  = null,
    "microsoftLauncherCustomWallpaperImageUrl"                 = null,
    "microsoftLauncherCustomWallpaperAllowUserModification"    = null,
    "microsoftLauncherFeedEnabled"                             = null,
    "microsoftLauncherFeedAllowUserModification"               = null,
    "microsoftLauncherDockPresenceConfiguration"               = null,
    "microsoftLauncherDockPresenceAllowUserModification"       = null,
    "microsoftLauncherSearchBarPlacementConfiguration"         = null,
    "enrollmentProfile"                                        = "notConfigured",
    "dataRoamingBlocked"                                       = null,
    "dateTimeConfigurationBlocked"                             = null,
    "detailedHelpText"                                         = null,
    "deviceOwnerLockScreenMessage"                             = null,
    "securityCommonCriteriaModeEnabled"                        = null,
    "factoryResetDeviceAdministratorEmails"                    = [],
    "factoryResetBlocked"                                      = null,
    "globalProxy"                                              = null,
    "googleAccountsBlocked"                                    = null,
    "kioskCustomizationDeviceSettingsBlocked"                  = null,
    "kioskCustomizationPowerButtonActionsBlocked"              = null,
    "kioskCustomizationStatusBar"                              = null,
    "kioskCustomizationSystemErrorWarnings"                    = null,
    "kioskCustomizationSystemNavigation"                       = null,
    "kioskModeScreenSaverConfigurationEnabled"                 = null,
    "kioskModeScreenSaverImageUrl"                             = null,
    "kioskModeScreenSaverDisplayTimeInSeconds"                 = null,
    "kioskModeScreenSaverStartDelayInSeconds"                  = null,
    "kioskModeScreenSaverDetectMediaDisabled"                  = null,
    "kioskModeWallpaperUrl"                                    = null,
    "kioskModeExitCode"                                        = null,
    "isKioskModeExitCodeSet"                                   = false,
    "kioskModeVirtualHomeButtonEnabled"                        = null,
    "kioskModeVirtualHomeButtonType"                           = null,
    "kioskModeBluetoothConfigurationEnabled"                   = null,
    "kioskModeWiFiConfigurationEnabled"                        = null,
    "kioskModeFlashlightConfigurationEnabled"                  = null,
    "kioskModeMediaVolumeConfigurationEnabled"                 = null,
    "kioskModeShowDeviceInfo"                                  = null,
    "kioskModeManagedSettingsEntryDisabled"                    = null,
    "kioskModeDebugMenuEasyAccessEnabled"                      = null,
    "kioskModeShowAppNotificationBadge"                        = null,
    "kioskModeScreenOrientation"                               = null,
    "kioskModeIconSize"                                        = null,
    "kioskModeFolderIcon"                                      = null,
    "kioskModeWifiAllowedSsids"                                = [],
    "kioskModeAppOrderEnabled"                                 = null,
    "kioskModeAppsInFolderOrderedByName"                       = null,
    "kioskModeGridHeight"                                      = null,
    "kioskModeGridWidth"                                       = null,
    "kioskModeLockHomeScreen"                                  = null,
    "kioskModeManagedHomeScreenAutoSignout"                    = null,
    "kioskModeManagedHomeScreenInactiveSignOutDelayInSeconds"  = null,
    "kioskModeManagedHomeScreenInactiveSignOutNoticeInSeconds" = null,
    "kioskModeManagedHomeScreenPinComplexity"                  = null,
    "kioskModeManagedHomeScreenPinRequired"                    = null,
    "kioskModeManagedHomeScreenPinRequiredToResume"            = null,
    "kioskModeManagedHomeScreenSignInBackground"               = null,
    "kioskModeManagedHomeScreenSignInBrandingLogo"             = null,
    "kioskModeManagedHomeScreenSignInEnabled"                  = null,
    "kioskModeUseManagedHomeScreenApp"                         = "notConfigured",
    "microphoneForceMute"                                      = null,
    "networkEscapeHatchAllowed"                                = null,
    "nfcBlockOutgoingBeam"                                     = null,
    "passwordBlockKeyguard"                                    = null,
    "passwordBlockKeyguardFeatures"                            = [],
    "passwordExpirationDays"                                   = null,
    "passwordMinimumLength"                                    = null,
    "passwordMinimumLetterCharacters"                          = null,
    "passwordMinimumLowerCaseCharacters"                       = null,
    "passwordMinimumNonLetterCharacters"                       = null,
    "passwordMinimumNumericCharacters"                         = null,
    "passwordMinimumSymbolCharacters"                          = null,
    "passwordMinimumUpperCaseCharacters"                       = null,
    "passwordMinutesOfInactivityBeforeScreenTimeout"           = null,
    "passwordPreviousPasswordCountToBlock"                     = null,
    "passwordRequiredType"                                     = null,
    "passwordRequireUnlock"                                    = null,
    "passwordSignInFailureCountBeforeFactoryReset"             = null,
    "playStoreMode"                                            = null,
    "screenCaptureBlocked"                                     = null,
    "securityDeveloperSettingsEnabled"                         = null,
    "securityRequireVerifyApps"                                = null,
    "shortHelpText"                                            = null,
    "statusBarBlocked"                                         = null,
    "stayOnModes"                                              = [],
    "storageAllowUsb"                                          = null,
    "storageBlockExternalMedia"                                = null,
    "storageBlockUsbFileTransfer"                              = null,
    "systemUpdateWindowStartMinutesAfterMidnight"              = null,
    "systemUpdateWindowEndMinutesAfterMidnight"                = null,
    "systemUpdateInstallType"                                  = null,
    "systemWindowsBlocked"                                     = null,
    "usersBlockAdd"                                            = null,
    "usersBlockRemove"                                         = null,
    "volumeBlockAdjustment"                                    = null,
    "vpnAlwaysOnLockdownMode"                                  = null,
    "vpnAlwaysOnPackageIdentifier"                             = null,
    "wifiBlockEditConfigurations"                              = null,
    "wifiBlockEditPolicyDefinedConfigurations"                 = null,
    "personalProfileAppsAllowInstallFromUnknownSources"        = null,
    "personalProfileCameraBlocked"                             = null,
    "personalProfileScreenCaptureBlocked"                      = null,
    "personalProfilePlayStoreMode"                             = null,
    "workProfilePasswordExpirationDays"                        = null,
    "workProfilePasswordMinimumLength"                         = null,
    "workProfilePasswordMinimumNumericCharacters"              = null,
    "workProfilePasswordMinimumNonLetterCharacters"            = null,
    "workProfilePasswordMinimumLetterCharacters"               = null,
    "workProfilePasswordMinimumLowerCaseCharacters"            = null,
    "workProfilePasswordMinimumUpperCaseCharacters"            = null,
    "workProfilePasswordMinimumSymbolCharacters"               = null,
    "workProfilePasswordPreviousPasswordCountToBlock"          = null,
    "workProfilePasswordSignInFailureCountBeforeFactoryReset"  = null,
    "workProfilePasswordRequiredType"                          = null,
    "workProfilePasswordRequireUnlock"                         = null,
    "locateDeviceUserlessDisabled"                             = null,
    "locateDeviceLostModeEnabled"                              = null,
    "shareDeviceLocationDisabled"                              = null,
    "deviceLocationMode"                                       = null,
    "azureAdSharedDeviceDataClearApps"                         = [],
    "kioskModeApps"                                            = [],
    "kioskModeManagedFolders"                                  = [],
    "kioskModeAppPositions"                                    = [],
    "kioskModeManagedHomeScreenAppSettings"                    = [],
    "systemUpdateFreezePeriods"                                = [],
    "personalProfilePersonalApplications"                      = [],
    "androidDeviceOwnerDelegatedScopeAppSettings"              = []
  })
}
```

#### Android Enterprise Wi-Fi

```terraform
# Android Enterprise Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_wifi" {
  display_name       = "Android Enterprise Wi-Fi"
  description        = "Android android enterprise wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = null,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "proxyExclusionList"                          = null,
    "macAddressRandomizationMode"                 = null
  })
}
```

#### Work Profile Wi-Fi

```terraform
# Android Work Profile Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_wifi" {
  display_name       = "Android Work Profile Wi-Fi"
  description        = "Android work profile wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026",
    "proxySettings"                               = "none",
    "proxyAutomaticConfigurationUrl"              = null
  })
}
```

#### AOSP Wi-Fi

```terraform
# Android AOSP Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_wifi" {
  display_name       = "Android AOSP Wi-Fi"
  description        = "Android aosp wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.aospDeviceOwnerWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = null,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026",
    "proxySetting"                                = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "proxyExclusionList"                          = []
  })
}
```

#### Android Enterprise Wi-Fi with Certificates

```terraform
# Android Enterprise Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_certificate_wifi_android_owner_root" {
  display_name       = "example-android-owner-root"
  description        = "Android android enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_certificate_wifi_scep_certificate" {
  display_name       = "example-android-owner-scep"
  description        = "Android android enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "certificateAccessType"              = null,
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames"  = [],
    "silentCertificateAccessDetails" = [],
    "rootCertificate@odata.bind"     = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_enterprise_certificate_wifi_android_owner_root.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_certificate_wifi" {
  display_name       = "Android Enterprise Wi-Fi with Certificates"
  description        = "Android android enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                           = "#microsoft.graph.androidDeviceOwnerEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"            = null,
    "deviceManagementApplicabilityRuleOsVersion"            = null,
    "deviceManagementApplicabilityRuleDeviceMode"           = null,
    "networkName"                                           = "Provider Enterprise WiFi",
    "ssid"                                                  = "Provider-Unassigned-EAP",
    "connectAutomatically"                                  = false,
    "connectWhenNetworkNameIsHidden"                        = null,
    "wiFiSecurityType"                                      = "wpaEnterprise",
    "preSharedKey"                                          = null,
    "proxySettings"                                         = "none",
    "proxyManualAddress"                                    = null,
    "proxyManualPort"                                       = null,
    "proxyAutomaticConfigurationUrl"                        = null,
    "proxyExclusionList"                                    = null,
    "macAddressRandomizationMode"                           = null,
    "eapType"                                               = "eapTls",
    "trustedServerCertificateNames"                         = [],
    "authenticationMethod"                                  = "certificate",
    "innerAuthenticationProtocolForEapTtls"                 = null,
    "innerAuthenticationProtocolForPeap"                    = null,
    "outerIdentityPrivacyTemporaryValue"                    = null,
    "rootCertificateForServerValidation@odata.bind"         = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_enterprise_certificate_wifi_android_owner_root.id}')",
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_enterprise_certificate_wifi_scep_certificate.id}')"
  })
}
```

#### Work Profile Wi-Fi with Certificates

```terraform
# Android Work Profile Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_certificate_wifi_android_work_root" {
  display_name       = "example-android-work-root"
  description        = "Android work profile wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_certificate_wifi_scep_certificate" {
  display_name       = "example-android-work-scep"
  description        = "Android work profile wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_work_profile_certificate_wifi_android_work_root.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_certificate_wifi" {
  display_name       = "Android Work Profile Wi-Fi with Certificates"
  description        = "Android work profile wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                           = "#microsoft.graph.androidWorkProfileEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"            = null,
    "deviceManagementApplicabilityRuleOsVersion"            = null,
    "deviceManagementApplicabilityRuleDeviceMode"           = null,
    "networkName"                                           = "Provider Enterprise WiFi",
    "ssid"                                                  = "Provider-Unassigned-EAP",
    "connectAutomatically"                                  = false,
    "connectWhenNetworkNameIsHidden"                        = false,
    "wiFiSecurityType"                                      = "wpaEnterprise",
    "preSharedKey"                                          = null,
    "proxySettings"                                         = "none",
    "proxyAutomaticConfigurationUrl"                        = null,
    "eapType"                                               = "eapTls",
    "trustedServerCertificateNames"                         = [],
    "authenticationMethod"                                  = "certificate",
    "innerAuthenticationProtocolForEapTtls"                 = null,
    "innerAuthenticationProtocolForPeap"                    = null,
    "outerIdentityPrivacyTemporaryValue"                    = null,
    "rootCertificateForServerValidation@odata.bind"         = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_work_profile_certificate_wifi_android_work_root.id}')",
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_work_profile_certificate_wifi_scep_certificate.id}')"
  })
}
```

#### AOSP Wi-Fi with Certificates

```terraform
# Android AOSP Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_certificate_wifi_aosp_root" {
  display_name       = "example-aosp-root"
  description        = "Android aosp wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.aospDeviceOwnerTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_certificate_wifi_scep_certificate" {
  display_name       = "example-aosp-scep"
  description        = "Android aosp wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.aospDeviceOwnerScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_aosp_certificate_wifi_aosp_root.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_certificate_wifi" {
  display_name       = "Android AOSP Wi-Fi with Certificates"
  description        = "Android aosp wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                           = "#microsoft.graph.aospDeviceOwnerEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"            = null,
    "deviceManagementApplicabilityRuleOsVersion"            = null,
    "deviceManagementApplicabilityRuleDeviceMode"           = null,
    "networkName"                                           = "Provider Enterprise WiFi",
    "ssid"                                                  = "Provider-Unassigned-EAP",
    "connectAutomatically"                                  = false,
    "connectWhenNetworkNameIsHidden"                        = null,
    "wiFiSecurityType"                                      = "wpaEnterprise",
    "preSharedKey"                                          = null,
    "proxySetting"                                          = "none",
    "proxyManualAddress"                                    = null,
    "proxyManualPort"                                       = null,
    "proxyAutomaticConfigurationUrl"                        = null,
    "proxyExclusionList"                                    = [],
    "eapType"                                               = "eapTls",
    "trustedServerCertificateNames"                         = [],
    "authenticationMethod"                                  = "certificate",
    "innerAuthenticationProtocolForEapTtls"                 = null,
    "innerAuthenticationProtocolForPeap"                    = null,
    "outerIdentityPrivacyTemporaryValue"                    = null,
    "rootCertificateForServerValidation@odata.bind"         = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_aosp_certificate_wifi_aosp_root.id}')",
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_aosp_certificate_wifi_scep_certificate.id}')"
  })
}
```

#### Android Enterprise VPN

```terraform
# Android Enterprise VPN.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_vpn_android_owner_root" {
  display_name       = "example-android-owner-root"
  description        = "Android android enterprise vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_vpn_scep_certificate" {
  display_name       = "example-android-owner-scep"
  description        = "Android android enterprise vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "certificateAccessType"              = null,
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames"  = [],
    "silentCertificateAccessDetails" = [],
    "rootCertificate@odata.bind"     = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_enterprise_vpn_android_owner_root.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_enterprise_vpn" {
  display_name       = "Android Enterprise VPN"
  description        = "Android android enterprise vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerVpnConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "authenticationMethod"                        = "certificate",
    "connectionName"                              = "Provider VPN Test",
    "role"                                        = null,
    "realm"                                       = null,
    "connectionType"                              = "ciscoAnyConnect",
    "proxyServer"                                 = null,
    "targetedPackageIds"                          = [],
    "alwaysOn"                                    = null,
    "alwaysOnLockdown"                            = null,
    "lockdownExclusionList"                       = [],
    "microsoftTunnelSiteId"                       = null,
    "proxyExclusionList"                          = [],
    "servers" = [
      {
        "description"     = "Provider unassigned test",
        "address"         = "vpn.example.invalid",
        "isDefaultServer" = true
      }
    ],
    "targetedMobileApps"             = [],
    "customData"                     = [],
    "customKeyValueData"             = [],
    "identityCertificate@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.android_enterprise_vpn_scep_certificate.id}')"
  })
}
```

#### Work Profile PKCS Certificate

```terraform
# Android Work Profile PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_profile_pkcs" {
  display_name       = "Android Work Profile PKCS Certificate"
  description        = "Android work profile pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfilePkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "subjectAlternativeNameType"                  = "none",
    "certificationAuthority"                      = "ca.example.invalid",
    "certificationAuthorityName"                  = "Provider Test CA",
    "certificateTemplateName"                     = "ProviderTest",
    "subjectAlternativeNameFormatString"          = null,
    "subjectNameFormatString"                     = "CN={{DeviceId}}",
    "certificateStore"                            = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = []
  })
}
```

#### AOSP Device Restrictions

```terraform
# Android AOSP Device Restrictions.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_aosp_restrictions" {
  display_name       = "Android AOSP Device Restrictions"
  description        = "Android aosp device restrictions example"
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
```

#### Clear a Wi-Fi Pre-shared Key

```terraform
# Android Clear a Wi-Fi Pre-shared Key.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_wifi_clear_pre_shared_key" {
  display_name       = "Android Clear a Wi-Fi Pre-shared Key"
  description        = "Android clear a wi-fi pre-shared key example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidDeviceOwnerWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = null,
    "wiFiSecurityType"                            = "wpaPersonal",
    "preSharedKey"                                = null,
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "proxyExclusionList"                          = null,
    "macAddressRandomizationMode"                 = null
  })
}
```

### iOS / iPadOS

#### General Device Configuration

```terraform
# iOS / iPadOS General Device Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_general_configuration" {
  display_name       = "iOS / iPadOS General Device Configuration"
  description        = "iOS / iPadOS general device configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                    = "#microsoft.graph.iosGeneralDeviceConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"     = null,
    "deviceManagementApplicabilityRuleOsVersion"     = null,
    "deviceManagementApplicabilityRuleDeviceMode"    = null,
    "accountBlockModification"                       = false,
    "activationLockAllowWhenSupervised"              = false,
    "airDropBlocked"                                 = true,
    "airDropForceUnmanagedDropTarget"                = false,
    "airPlayForcePairingPasswordForOutgoingRequests" = false,
    "appleWatchBlockPairing"                         = false,
    "appleWatchForceWristDetection"                  = false,
    "appleNewsBlocked"                               = false,
    "appsVisibilityListType"                         = "none",
    "appStoreBlockAutomaticDownloads"                = false,
    "appStoreBlocked"                                = false,
    "appStoreBlockInAppPurchases"                    = false,
    "appStoreBlockUIAppInstallation"                 = false,
    "appStoreRequirePassword"                        = false,
    "autoFillForceAuthentication"                    = false,
    "bluetoothBlockModification"                     = false,
    "cameraBlocked"                                  = true,
    "cellularBlockDataRoaming"                       = false,
    "cellularBlockGlobalBackgroundFetchWhileRoaming" = false,
    "cellularBlockPerAppDataModification"            = false,
    "cellularBlockPersonalHotspot"                   = false,
    "cellularBlockPlanModification"                  = false,
    "cellularBlockVoiceRoaming"                      = false,
    "certificatesBlockUntrustedTlsCertificates"      = false,
    "classroomAppBlockRemoteScreenObservation"       = false,
    "classroomAppForceUnpromptedScreenObservation"   = false,
    "classroomForceAutomaticallyJoinClasses"         = false,
    "classroomForceUnpromptedAppAndDeviceLock"       = false,
    "compliantAppListType"                           = "none",
    "configurationProfileBlockChanges"               = false,
    "definitionLookupBlocked"                        = false,
    "deviceBlockEnableRestrictions"                  = false,
    "deviceBlockEraseContentAndSettings"             = false,
    "deviceBlockNameModification"                    = false,
    "diagnosticDataBlockSubmission"                  = false,
    "diagnosticDataBlockSubmissionModification"      = false,
    "documentsBlockManagedDocumentsInUnmanagedApps"  = false,
    "documentsBlockUnmanagedDocumentsInManagedApps"  = false,
    "emailInDomainSuffixes"                          = [],
    "enterpriseAppBlockTrust"                        = false,
    "enterpriseAppBlockTrustModification"            = false,
    "esimBlockModification"                          = false,
    "faceTimeBlocked"                                = false,
    "findMyFriendsBlocked"                           = false,
    "gamingBlockGameCenterFriends"                   = false,
    "gamingBlockMultiplayer"                         = false,
    "gameCenterBlocked"                              = false,
    "hostPairingBlocked"                             = false,
    "iBooksStoreBlocked"                             = false,
    "iBooksStoreBlockErotica"                        = false,
    "iCloudBlockActivityContinuation"                = false,
    "iCloudBlockBackup"                              = false,
    "iCloudBlockDocumentSync"                        = false,
    "iCloudBlockManagedAppsSync"                     = false,
    "iCloudBlockPhotoLibrary"                        = false,
    "iCloudBlockPhotoStreamSync"                     = false,
    "iCloudBlockSharedPhotoStream"                   = false,
    "iCloudRequireEncryptedBackup"                   = false,
    "iTunesBlockExplicitContent"                     = false,
    "iTunesBlockMusicService"                        = false,
    "iTunesBlockRadio"                               = false,
    "keyboardBlockAutoCorrect"                       = false,
    "keyboardBlockDictation"                         = false,
    "keyboardBlockPredictive"                        = false,
    "keyboardBlockShortcuts"                         = false,
    "keyboardBlockSpellCheck"                        = false,
    "kioskModeAllowAssistiveSpeak"                   = false,
    "kioskModeAllowAssistiveTouchSettings"           = false,
    "kioskModeAllowAutoLock"                         = false,
    "kioskModeBlockAutoLock"                         = false,
    "kioskModeAllowColorInversionSettings"           = false,
    "kioskModeAllowRingerSwitch"                     = false,
    "kioskModeBlockRingerSwitch"                     = false,
    "kioskModeAllowScreenRotation"                   = false,
    "kioskModeBlockScreenRotation"                   = false,
    "kioskModeAllowSleepButton"                      = false,
    "kioskModeBlockSleepButton"                      = false,
    "kioskModeAllowTouchscreen"                      = false,
    "kioskModeBlockTouchscreen"                      = false,
    "kioskModeEnableVoiceControl"                    = false,
    "kioskModeAllowVoiceControlModification"         = false,
    "kioskModeAllowVoiceOverSettings"                = false,
    "kioskModeAllowVolumeButtons"                    = false,
    "kioskModeBlockVolumeButtons"                    = false,
    "kioskModeAllowZoomSettings"                     = false,
    "kioskModeAppStoreUrl"                           = null,
    "kioskModeBuiltInAppId"                          = null,
    "kioskModeRequireAssistiveTouch"                 = false,
    "kioskModeRequireColorInversion"                 = false,
    "kioskModeRequireMonoAudio"                      = false,
    "kioskModeRequireVoiceOver"                      = false,
    "kioskModeRequireZoom"                           = false,
    "kioskModeManagedAppId"                          = null,
    "lockScreenBlockControlCenter"                   = false,
    "lockScreenBlockNotificationView"                = false,
    "lockScreenBlockPassbook"                        = false,
    "lockScreenBlockTodayView"                       = false,
    "mediaContentRatingAustralia"                    = null,
    "mediaContentRatingCanada"                       = null,
    "mediaContentRatingFrance"                       = null,
    "mediaContentRatingGermany"                      = null,
    "mediaContentRatingIreland"                      = null,
    "mediaContentRatingJapan"                        = null,
    "mediaContentRatingNewZealand"                   = null,
    "mediaContentRatingUnitedKingdom"                = null,
    "mediaContentRatingUnitedStates"                 = null,
    "mediaContentRatingApps"                         = "allAllowed",
    "messagesBlocked"                                = false,
    "notificationsBlockSettingsModification"         = false,
    "passcodeBlockFingerprintUnlock"                 = false,
    "passcodeBlockFingerprintModification"           = false,
    "passcodeBlockModification"                      = false,
    "passcodeBlockSimple"                            = false,
    "passcodeExpirationDays"                         = null,
    "passcodeMinimumLength"                          = null,
    "passcodeMinutesOfInactivityBeforeLock"          = null,
    "passcodeMinutesOfInactivityBeforeScreenTimeout" = null,
    "passcodeMinimumCharacterSetCount"               = null,
    "passcodePreviousPasscodeBlockCount"             = null,
    "passcodeSignInFailureCountBeforeWipe"           = null,
    "passcodeRequiredType"                           = "deviceDefault",
    "passcodeRequired"                               = false,
    "podcastsBlocked"                                = false,
    "proximityBlockSetupToNewDevice"                 = false,
    "safariBlockAutofill"                            = false,
    "safariBlockJavaScript"                          = false,
    "safariBlockPopups"                              = false,
    "safariBlocked"                                  = false,
    "safariCookieSettings"                           = "browserDefault",
    "safariManagedDomains"                           = [],
    "safariPasswordAutoFillDomains"                  = [],
    "safariRequireFraudWarning"                      = false,
    "screenCaptureBlocked"                           = false,
    "siriBlocked"                                    = false,
    "siriBlockedWhenLocked"                          = false,
    "siriBlockUserGeneratedContent"                  = false,
    "siriRequireProfanityFilter"                     = false,
    "softwareUpdatesEnforcedDelayInDays"             = null,
    "softwareUpdatesForceDelayed"                    = false,
    "spotlightBlockInternetResults"                  = false,
    "voiceDialingBlocked"                            = false,
    "wallpaperBlockModification"                     = false,
    "wiFiConnectOnlyToConfiguredNetworks"            = false,
    "classroomForceRequestPermissionToLeaveClasses"  = false,
    "keychainBlockCloudSync"                         = false,
    "pkiBlockOTAUpdates"                             = false,
    "privacyForceLimitAdTracking"                    = false,
    "enterpriseBookBlockBackup"                      = false,
    "enterpriseBookBlockMetadataSync"                = false,
    "airPrintBlocked"                                = false,
    "airPrintBlockCredentialsStorage"                = false,
    "airPrintForceTrustedTLS"                        = false,
    "airPrintBlockiBeaconDiscovery"                  = false,
    "filesNetworkDriveAccessBlocked"                 = false,
    "filesUsbDriveAccessBlocked"                     = false,
    "wifiPowerOnForced"                              = false,
    "blockSystemAppRemoval"                          = false,
    "vpnBlockCreation"                               = false,
    "appRemovalBlocked"                              = false,
    "usbRestrictedModeBlocked"                       = false,
    "passwordBlockAutoFill"                          = false,
    "passwordBlockProximityRequests"                 = false,
    "passwordBlockAirDropSharing"                    = false,
    "dateAndTimeForceSetAutomatically"               = false,
    "contactsAllowManagedToUnmanagedWrite"           = false,
    "contactsAllowUnmanagedToManagedRead"            = false,
    "cellularBlockPersonalHotspotModification"       = false,
    "continuousPathKeyboardBlocked"                  = false,
    "findMyDeviceInFindMyAppBlocked"                 = false,
    "findMyFriendsInFindMyAppBlocked"                = false,
    "iTunesBlocked"                                  = false,
    "sharedDeviceBlockTemporarySessions"             = false,
    "appClipsBlocked"                                = false,
    "applePersonalizedAdsBlocked"                    = false,
    "nfcBlocked"                                     = false,
    "autoUnlockBlocked"                              = false,
    "unpairedExternalBootToRecoveryAllowed"          = false,
    "onDeviceOnlyDictationForced"                    = false,
    "wiFiConnectToAllowedNetworksOnlyForced"         = false,
    "onDeviceOnlyTranslationForced"                  = false,
    "managedPasteboardRequired"                      = false,
    "iCloudPrivateRelayBlocked"                      = false,
    "kioskModeAppType"                               = "notConfigured",
    "appsSingleAppModeList"                          = [],
    "appsVisibilityList"                             = [],
    "compliantAppsList"                              = [],
    "networkUsageRules"                              = []
  })
}
```

#### Device Features

```terraform
# iOS / iPadOS Device Features.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_device_features" {
  display_name       = "iOS / iPadOS Device Features"
  description        = "iOS / iPadOS device features example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosDeviceFeaturesConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "assetTagTemplate"                            = null,
    "contentFilterSettings"                       = null,
    "lockScreenFootnote"                          = null,
    "homeScreenGridWidth"                         = null,
    "homeScreenGridHeight"                        = null,
    "singleSignOnSettings"                        = null,
    "wallpaperDisplayLocation"                    = "notConfigured",
    "wallpaperImage"                              = null,
    "singleSignOnExtension"                       = null,
    "iosSingleSignOnExtension"                    = null,
    "airPrintDestinations"                        = [],
    "homeScreenDockIcons"                         = [],
    "homeScreenPages" = [
      {
        "displayName" = "Page 1",
        "icons" = [
          {
            "@odata.type" = "#microsoft.graph.iosHomeScreenApp",
            "displayName" = "Safari",
            "bundleID"    = "com.apple.mobilesafari",
            "isWebClip"   = false
          }
        ]
      }
    ],
    "notificationSettings" = []
  })
}
```

#### Custom Configuration

```terraform
# iOS / iPadOS Custom Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_custom_configuration" {
  display_name       = "iOS / iPadOS Custom Configuration"
  description        = "iOS / iPadOS custom configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })
}
```

#### Trusted Root Certificate

```terraform
# iOS / iPadOS Trusted Root Certificate.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_trusted_root_certificate" {
  display_name       = "iOS / iPadOS Trusted Root Certificate"
  description        = "iOS / iPadOS trusted root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}
```

#### SCEP Certificate

```terraform
# iOS / iPadOS SCEP Certificate.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_scep_certificate_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS scep certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_scep_certificate" {
  display_name       = "iOS / iPadOS SCEP Certificate"
  description        = "iOS / iPadOS scep certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_scep_certificate_root_certificate.id}')"
  })
}
```

#### Home Screen with Nested Folders

```terraform
# iOS / iPadOS Home Screen with Nested Folders.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_home_screen_folders" {
  display_name       = "iOS / iPadOS Home Screen with Nested Folders"
  description        = "iOS / iPadOS home screen with nested folders example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosDeviceFeaturesConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "assetTagTemplate"                            = null,
    "contentFilterSettings"                       = null,
    "lockScreenFootnote"                          = null,
    "homeScreenGridWidth"                         = null,
    "homeScreenGridHeight"                        = null,
    "singleSignOnSettings"                        = null,
    "wallpaperDisplayLocation"                    = "notConfigured",
    "wallpaperImage"                              = null,
    "singleSignOnExtension"                       = null,
    "iosSingleSignOnExtension"                    = null,
    "airPrintDestinations"                        = [],
    "homeScreenDockIcons"                         = [],
    "homeScreenPages" = [
      {
        "displayName" = "Page 1",
        "icons" = [
          {
            "@odata.type" = "#microsoft.graph.iosHomeScreenFolder",
            "displayName" = "Apps",
            "pages" = [
              {
                "displayName" = "Folder page",
                "apps" = [
                  {
                    "displayName" = "Safari",
                    "bundleID"    = "com.apple.mobilesafari",
                    "isWebClip"   = false
                  }
                ]
              }
            ]
          }
        ]
      }
    ],
    "notificationSettings" = []
  })
}
```

#### Wi-Fi

```terraform
# iOS / iPadOS Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_wifi" {
  display_name       = "iOS / iPadOS Wi-Fi"
  description        = "iOS / iPadOS wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaPersonal",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "disableMacAddressRandomization"              = null,
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026"
  })
}
```

#### Enterprise Wi-Fi with Certificates

```terraform
# iOS / iPadOS Enterprise Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_enterprise_wifi_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_enterprise_wifi_scep_certificate" {
  display_name       = "example-ios-scep"
  description        = "iOS / iPadOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_enterprise_wifi_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_enterprise_wifi" {
  display_name       = "iOS / iPadOS Enterprise Wi-Fi with Certificates"
  description        = "iOS / iPadOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider Enterprise WiFi",
    "ssid"                                        = "Provider-Unassigned-EAP",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaEnterprise",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "disableMacAddressRandomization"              = null,
    "preSharedKey"                                = null,
    "eapType"                                     = "eapTls",
    "eapFastConfiguration"                        = null,
    "trustedServerCertificateNames"               = [],
    "authenticationMethod"                        = "certificate",
    "innerAuthenticationProtocolForEapTtls"       = null,
    "outerIdentityPrivacyTemporaryValue"          = null,
    "usernameFormatString"                        = null,
    "passwordFormatString"                        = null,
    "rootCertificatesForServerValidation@odata.bind" = [
      "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_enterprise_wifi_root_certificate.id}')"
    ],
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_enterprise_wifi_scep_certificate.id}')"
  })
}
```

#### VPN

```terraform
# iOS / iPadOS VPN.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_vpn_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_vpn_scep_certificate" {
  display_name       = "example-ios-scep"
  description        = "iOS / iPadOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_vpn_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_vpn" {
  display_name       = "iOS / iPadOS VPN"
  description        = "iOS / iPadOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosVpnConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "connectionName"                              = "Provider VPN Test",
    "connectionType"                              = "ciscoAnyConnect",
    "loginGroupOrDomain"                          = null,
    "role"                                        = null,
    "realm"                                       = null,
    "identifier"                                  = null,
    "enableSplitTunneling"                        = false,
    "authenticationMethod"                        = "certificate",
    "enablePerApp"                                = null,
    "safariDomains"                               = [],
    "providerType"                                = null,
    "associatedDomains"                           = [],
    "excludedDomains"                             = [],
    "disableOnDemandUserOverride"                 = null,
    "disconnectOnIdle"                            = null,
    "disconnectOnIdleTimerInSeconds"              = null,
    "includeAllNetworks"                          = null,
    "excludeLocalNetworks"                        = null,
    "proxyServer"                                 = null,
    "optInToDeviceIdSharing"                      = null,
    "userDomain"                                  = null,
    "strictEnforcement"                           = null,
    "cloudName"                                   = null,
    "excludeList"                                 = [],
    "microsoftTunnelSiteId"                       = null,
    "server" = {
      "description"     = "Provider unassigned test",
      "address"         = "vpn.example.invalid",
      "isDefaultServer" = true
    },
    "customData"                     = [],
    "customKeyValueData"             = [],
    "onDemandRules"                  = [],
    "targetedMobileApps"             = [],
    "identityCertificate@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_vpn_scep_certificate.id}')"
  })
}
```

#### Email

```terraform
# iOS / iPadOS Email.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_email_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS email example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_email_scep_certificate" {
  display_name       = "example-ios-scep"
  description        = "iOS / iPadOS email example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_email_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_email" {
  display_name       = "iOS / iPadOS Email"
  description        = "iOS / iPadOS email example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                   = "#microsoft.graph.iosEasEmailProfileConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"    = null,
    "deviceManagementApplicabilityRuleOsVersion"    = null,
    "deviceManagementApplicabilityRuleDeviceMode"   = null,
    "usernameSource"                                = "userPrincipalName",
    "usernameAADSource"                             = null,
    "userDomainNameSource"                          = null,
    "customDomainName"                              = null,
    "accountName"                                   = "Provider Email Test",
    "authenticationMethod"                          = "certificate",
    "blockMovingMessagesToOtherEmailAccounts"       = null,
    "blockSendingEmailFromThirdPartyApps"           = null,
    "blockSyncingRecentlyUsedEmailAddresses"        = null,
    "durationOfEmailToSync"                         = "userDefined",
    "emailAddressSource"                            = "userPrincipalName",
    "easServices"                                   = null,
    "easServicesUserOverrideEnabled"                = null,
    "hostName"                                      = "mail.example.invalid",
    "requireSmime"                                  = null,
    "smimeEnablePerMessageSwitch"                   = null,
    "smimeEncryptByDefaultEnabled"                  = null,
    "smimeSigningEnabled"                           = null,
    "smimeSigningUserOverrideEnabled"               = null,
    "smimeEncryptByDefaultUserOverrideEnabled"      = null,
    "smimeSigningCertificateUserOverrideEnabled"    = null,
    "smimeEncryptionCertificateUserOverrideEnabled" = null,
    "requireSsl"                                    = true,
    "useOAuth"                                      = null,
    "signingCertificateType"                        = null,
    "encryptionCertificateType"                     = null,
    "perAppVPNProfileId"                            = null,
    "identityCertificate@odata.bind"                = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_email_scep_certificate.id}')"
  })
}
```

#### PKCS Certificate

```terraform
# iOS / iPadOS PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_pkcs_certificate" {
  display_name       = "iOS / iPadOS PKCS Certificate"
  description        = "iOS / iPadOS pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosPkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "certificationAuthority"                      = "ca.example.invalid",
    "certificationAuthorityName"                  = "Provider Test CA",
    "certificateTemplateName"                     = "ProviderTest",
    "subjectAlternativeNameFormatString"          = null,
    "subjectNameFormatString"                     = "CN={{DeviceId}}",
    "certificateStore"                            = "user",
    "customSubjectAlternativeNames"               = []
  })
}
```

#### Reset General Restrictions

```terraform
# iOS / iPadOS Reset General Restrictions.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_reset_restrictions" {
  display_name       = "iOS / iPadOS Reset General Restrictions"
  description        = "iOS / iPadOS reset general restrictions example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                    = "#microsoft.graph.iosGeneralDeviceConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"     = null,
    "deviceManagementApplicabilityRuleOsVersion"     = null,
    "deviceManagementApplicabilityRuleDeviceMode"    = null,
    "accountBlockModification"                       = false,
    "activationLockAllowWhenSupervised"              = false,
    "airDropBlocked"                                 = false,
    "airDropForceUnmanagedDropTarget"                = false,
    "airPlayForcePairingPasswordForOutgoingRequests" = false,
    "appleWatchBlockPairing"                         = false,
    "appleWatchForceWristDetection"                  = false,
    "appleNewsBlocked"                               = false,
    "appsVisibilityListType"                         = "none",
    "appStoreBlockAutomaticDownloads"                = false,
    "appStoreBlocked"                                = false,
    "appStoreBlockInAppPurchases"                    = false,
    "appStoreBlockUIAppInstallation"                 = false,
    "appStoreRequirePassword"                        = false,
    "autoFillForceAuthentication"                    = false,
    "bluetoothBlockModification"                     = false,
    "cameraBlocked"                                  = false,
    "cellularBlockDataRoaming"                       = false,
    "cellularBlockGlobalBackgroundFetchWhileRoaming" = false,
    "cellularBlockPerAppDataModification"            = false,
    "cellularBlockPersonalHotspot"                   = false,
    "cellularBlockPlanModification"                  = false,
    "cellularBlockVoiceRoaming"                      = false,
    "certificatesBlockUntrustedTlsCertificates"      = false,
    "classroomAppBlockRemoteScreenObservation"       = false,
    "classroomAppForceUnpromptedScreenObservation"   = false,
    "classroomForceAutomaticallyJoinClasses"         = false,
    "classroomForceUnpromptedAppAndDeviceLock"       = false,
    "compliantAppListType"                           = "none",
    "configurationProfileBlockChanges"               = false,
    "definitionLookupBlocked"                        = false,
    "deviceBlockEnableRestrictions"                  = false,
    "deviceBlockEraseContentAndSettings"             = false,
    "deviceBlockNameModification"                    = false,
    "diagnosticDataBlockSubmission"                  = false,
    "diagnosticDataBlockSubmissionModification"      = false,
    "documentsBlockManagedDocumentsInUnmanagedApps"  = false,
    "documentsBlockUnmanagedDocumentsInManagedApps"  = false,
    "emailInDomainSuffixes"                          = [],
    "enterpriseAppBlockTrust"                        = false,
    "enterpriseAppBlockTrustModification"            = false,
    "esimBlockModification"                          = false,
    "faceTimeBlocked"                                = false,
    "findMyFriendsBlocked"                           = false,
    "gamingBlockGameCenterFriends"                   = false,
    "gamingBlockMultiplayer"                         = false,
    "gameCenterBlocked"                              = false,
    "hostPairingBlocked"                             = false,
    "iBooksStoreBlocked"                             = false,
    "iBooksStoreBlockErotica"                        = false,
    "iCloudBlockActivityContinuation"                = false,
    "iCloudBlockBackup"                              = false,
    "iCloudBlockDocumentSync"                        = false,
    "iCloudBlockManagedAppsSync"                     = false,
    "iCloudBlockPhotoLibrary"                        = false,
    "iCloudBlockPhotoStreamSync"                     = false,
    "iCloudBlockSharedPhotoStream"                   = false,
    "iCloudRequireEncryptedBackup"                   = false,
    "iTunesBlockExplicitContent"                     = false,
    "iTunesBlockMusicService"                        = false,
    "iTunesBlockRadio"                               = false,
    "keyboardBlockAutoCorrect"                       = false,
    "keyboardBlockDictation"                         = false,
    "keyboardBlockPredictive"                        = false,
    "keyboardBlockShortcuts"                         = false,
    "keyboardBlockSpellCheck"                        = false,
    "kioskModeAllowAssistiveSpeak"                   = false,
    "kioskModeAllowAssistiveTouchSettings"           = false,
    "kioskModeAllowAutoLock"                         = false,
    "kioskModeBlockAutoLock"                         = false,
    "kioskModeAllowColorInversionSettings"           = false,
    "kioskModeAllowRingerSwitch"                     = false,
    "kioskModeBlockRingerSwitch"                     = false,
    "kioskModeAllowScreenRotation"                   = false,
    "kioskModeBlockScreenRotation"                   = false,
    "kioskModeAllowSleepButton"                      = false,
    "kioskModeBlockSleepButton"                      = false,
    "kioskModeAllowTouchscreen"                      = false,
    "kioskModeBlockTouchscreen"                      = false,
    "kioskModeEnableVoiceControl"                    = false,
    "kioskModeAllowVoiceControlModification"         = false,
    "kioskModeAllowVoiceOverSettings"                = false,
    "kioskModeAllowVolumeButtons"                    = false,
    "kioskModeBlockVolumeButtons"                    = false,
    "kioskModeAllowZoomSettings"                     = false,
    "kioskModeAppStoreUrl"                           = null,
    "kioskModeBuiltInAppId"                          = null,
    "kioskModeRequireAssistiveTouch"                 = false,
    "kioskModeRequireColorInversion"                 = false,
    "kioskModeRequireMonoAudio"                      = false,
    "kioskModeRequireVoiceOver"                      = false,
    "kioskModeRequireZoom"                           = false,
    "kioskModeManagedAppId"                          = null,
    "lockScreenBlockControlCenter"                   = false,
    "lockScreenBlockNotificationView"                = false,
    "lockScreenBlockPassbook"                        = false,
    "lockScreenBlockTodayView"                       = false,
    "mediaContentRatingAustralia"                    = null,
    "mediaContentRatingCanada"                       = null,
    "mediaContentRatingFrance"                       = null,
    "mediaContentRatingGermany"                      = null,
    "mediaContentRatingIreland"                      = null,
    "mediaContentRatingJapan"                        = null,
    "mediaContentRatingNewZealand"                   = null,
    "mediaContentRatingUnitedKingdom"                = null,
    "mediaContentRatingUnitedStates"                 = null,
    "mediaContentRatingApps"                         = "allAllowed",
    "messagesBlocked"                                = false,
    "notificationsBlockSettingsModification"         = false,
    "passcodeBlockFingerprintUnlock"                 = false,
    "passcodeBlockFingerprintModification"           = false,
    "passcodeBlockModification"                      = false,
    "passcodeBlockSimple"                            = false,
    "passcodeExpirationDays"                         = null,
    "passcodeMinimumLength"                          = null,
    "passcodeMinutesOfInactivityBeforeLock"          = null,
    "passcodeMinutesOfInactivityBeforeScreenTimeout" = null,
    "passcodeMinimumCharacterSetCount"               = null,
    "passcodePreviousPasscodeBlockCount"             = null,
    "passcodeSignInFailureCountBeforeWipe"           = null,
    "passcodeRequiredType"                           = "deviceDefault",
    "passcodeRequired"                               = false,
    "podcastsBlocked"                                = false,
    "proximityBlockSetupToNewDevice"                 = false,
    "safariBlockAutofill"                            = false,
    "safariBlockJavaScript"                          = false,
    "safariBlockPopups"                              = false,
    "safariBlocked"                                  = false,
    "safariCookieSettings"                           = "browserDefault",
    "safariManagedDomains"                           = [],
    "safariPasswordAutoFillDomains"                  = [],
    "safariRequireFraudWarning"                      = false,
    "screenCaptureBlocked"                           = false,
    "siriBlocked"                                    = false,
    "siriBlockedWhenLocked"                          = false,
    "siriBlockUserGeneratedContent"                  = false,
    "siriRequireProfanityFilter"                     = false,
    "softwareUpdatesEnforcedDelayInDays"             = null,
    "softwareUpdatesForceDelayed"                    = false,
    "spotlightBlockInternetResults"                  = false,
    "voiceDialingBlocked"                            = false,
    "wallpaperBlockModification"                     = false,
    "wiFiConnectOnlyToConfiguredNetworks"            = false,
    "classroomForceRequestPermissionToLeaveClasses"  = false,
    "keychainBlockCloudSync"                         = false,
    "pkiBlockOTAUpdates"                             = false,
    "privacyForceLimitAdTracking"                    = false,
    "enterpriseBookBlockBackup"                      = false,
    "enterpriseBookBlockMetadataSync"                = false,
    "airPrintBlocked"                                = false,
    "airPrintBlockCredentialsStorage"                = false,
    "airPrintForceTrustedTLS"                        = false,
    "airPrintBlockiBeaconDiscovery"                  = false,
    "filesNetworkDriveAccessBlocked"                 = false,
    "filesUsbDriveAccessBlocked"                     = false,
    "wifiPowerOnForced"                              = false,
    "blockSystemAppRemoval"                          = false,
    "vpnBlockCreation"                               = false,
    "appRemovalBlocked"                              = false,
    "usbRestrictedModeBlocked"                       = false,
    "passwordBlockAutoFill"                          = false,
    "passwordBlockProximityRequests"                 = false,
    "passwordBlockAirDropSharing"                    = false,
    "dateAndTimeForceSetAutomatically"               = false,
    "contactsAllowManagedToUnmanagedWrite"           = false,
    "contactsAllowUnmanagedToManagedRead"            = false,
    "cellularBlockPersonalHotspotModification"       = false,
    "continuousPathKeyboardBlocked"                  = false,
    "findMyDeviceInFindMyAppBlocked"                 = false,
    "findMyFriendsInFindMyAppBlocked"                = false,
    "iTunesBlocked"                                  = false,
    "sharedDeviceBlockTemporarySessions"             = false,
    "appClipsBlocked"                                = false,
    "applePersonalizedAdsBlocked"                    = false,
    "nfcBlocked"                                     = false,
    "autoUnlockBlocked"                              = false,
    "unpairedExternalBootToRecoveryAllowed"          = false,
    "onDeviceOnlyDictationForced"                    = false,
    "wiFiConnectToAllowedNetworksOnlyForced"         = false,
    "onDeviceOnlyTranslationForced"                  = false,
    "managedPasteboardRequired"                      = false,
    "iCloudPrivateRelayBlocked"                      = false,
    "kioskModeAppType"                               = "notConfigured",
    "appsSingleAppModeList"                          = [],
    "appsVisibilityList"                             = [],
    "compliantAppsList"                              = [],
    "networkUsageRules"                              = []
  })
}
```

#### Custom Configuration Using a JSON Heredoc

```terraform
# iOS / iPadOS Custom Configuration Using a JSON Heredoc.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_custom_json_heredoc" {
  display_name = "iOS / iPadOS Custom Configuration Using a JSON Heredoc"
  description  = "iOS / iPadOS custom configuration using a json heredoc example"
  settings     = <<JSON
{
  "@odata.type": "#microsoft.graph.iosCustomConfiguration",
  "deviceManagementApplicabilityRuleOsEdition": null,
  "deviceManagementApplicabilityRuleOsVersion": null,
  "deviceManagementApplicabilityRuleDeviceMode": null,
  "payloadName": "Wire Study",
  "payloadFileName": "wire-study.mobileconfig",
  "payload": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
}
JSON
}
```

#### Minimal Group Assignment

```terraform
# iOS / iPadOS Minimal Group Assignment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_minimal_assignments" {
  depends_on         = [time_sleep.ios_minimal_assignments_directory_consistency]
  display_name       = "iOS / iPadOS Minimal Group Assignment"
  description        = "iOS / iPadOS minimal group assignment example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.ios_minimal_assignments_include.id
      filter_type = "none"
    }
  ]
}

resource "microsoft365_graph_beta_groups_group" "ios_minimal_assignments_include" {
  display_name     = "example-ios-minimal-assignments-include"
  mail_nickname    = "example-ios-minimal-assignments-include"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ios_minimal_assignments_exclude" {
  display_name     = "example-ios-minimal-assignments-exclude"
  mail_nickname    = "example-ios-minimal-assignments-exclude"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "ios_minimal_assignments_directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.ios_minimal_assignments_include, microsoft365_graph_beta_groups_group.ios_minimal_assignments_exclude]
  create_duration = "30s"
}
```

#### All Assignment Target Types

```terraform
# iOS / iPadOS All Assignment Target Types.
# The empty custom payload demonstrates all four assignment target types.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_maximal_assignments" {
  depends_on         = [time_sleep.ios_maximal_assignments_directory_consistency]
  display_name       = "iOS / iPadOS All Assignment Target Types"
  description        = "iOS / iPadOS all assignment target types example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.ios_maximal_assignments_include.id
      filter_type = "none"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.ios_maximal_assignments_exclude.id
      filter_type = "none"
    }
  ]
}

resource "microsoft365_graph_beta_groups_group" "ios_maximal_assignments_include" {
  display_name     = "example-ios-maximal-assignments-include"
  mail_nickname    = "example-ios-maximal-assignments-include"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ios_maximal_assignments_exclude" {
  display_name     = "example-ios-maximal-assignments-exclude"
  mail_nickname    = "example-ios-maximal-assignments-exclude"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "ios_maximal_assignments_directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.ios_maximal_assignments_include, microsoft365_graph_beta_groups_group.ios_maximal_assignments_exclude]
  create_duration = "30s"
}
```

#### Filtered Group and Exclusion Assignments

```terraform
# iOS / iPadOS Filtered Group and Exclusion Assignments.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_filtered_assignments" {
  depends_on         = [time_sleep.ios_filtered_assignments_directory_consistency]
  display_name       = "iOS / iPadOS Filtered Group and Exclusion Assignments"
  description        = "iOS / iPadOS filtered group and exclusion assignments example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.ios_filtered_assignments_include.id
      filter_type = "include"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.ios_filtered_assignments_filter.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.ios_filtered_assignments_exclude.id
    }
  ]
}

resource "microsoft365_graph_beta_groups_group" "ios_filtered_assignments_include" {
  display_name     = "example-ios-filtered-assignments-include"
  mail_nickname    = "example-ios-filtered-assignments-include"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ios_filtered_assignments_exclude" {
  display_name     = "example-ios-filtered-assignments-exclude"
  mail_nickname    = "example-ios-filtered-assignments-exclude"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_assignment_filter" "ios_filtered_assignments_filter" {
  display_name                      = "example-ios-filtered-assignments-filter"
  platform                          = "iOS"
  rule                              = "(device.deviceName -startsWith \"ProviderTest\")"
  assignment_filter_management_type = "devices"
}

resource "time_sleep" "ios_filtered_assignments_directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.ios_filtered_assignments_include, microsoft365_graph_beta_groups_group.ios_filtered_assignments_exclude]
  create_duration = "30s"
}
```

#### Clear Assignments

```terraform
# iOS / iPadOS Clear Assignments.
# Use assignments = [] on an existing profile to remove every assignment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_clear_assignments" {
  depends_on         = [time_sleep.ios_clear_assignments_directory_consistency]
  display_name       = "iOS / iPadOS Clear Assignments"
  description        = "iOS / iPadOS clear assignments example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K"
  })

  assignments = []
}

resource "microsoft365_graph_beta_groups_group" "ios_clear_assignments_include" {
  display_name     = "example-ios-clear-assignments-include"
  mail_nickname    = "example-ios-clear-assignments-include"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ios_clear_assignments_exclude" {
  display_name     = "example-ios-clear-assignments-exclude"
  mail_nickname    = "example-ios-clear-assignments-exclude"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_assignment_filter" "ios_clear_assignments_filter" {
  display_name                      = "example-ios-clear-assignments-filter"
  platform                          = "iOS"
  rule                              = "(device.deviceName -startsWith \"ProviderTest\")"
  assignment_filter_management_type = "devices"
}

resource "time_sleep" "ios_clear_assignments_directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.ios_clear_assignments_include, microsoft365_graph_beta_groups_group.ios_clear_assignments_exclude]
  create_duration = "30s"
}
```

#### Replace an SCEP Root Certificate

```terraform
# iOS / iPadOS Replace an SCEP Root Certificate.
# Replace the sample public certificate with your own root certificate.
# Change rootCertificate@odata.bind to switch the SCEP profile to the additional root.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_replace_root_certificate_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS replace an scep root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_replace_root_certificate_additional_root_certificate" {
  display_name       = "example-ios-second-root"
  description        = "iOS / iPadOS replace an scep root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_replace_root_certificate" {
  display_name       = "iOS / iPadOS Replace an SCEP Root Certificate"
  description        = "iOS / iPadOS replace an scep root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_replace_root_certificate_additional_root_certificate.id}')"
  })
}
```

#### Enterprise Wi-Fi with Multiple Root Certificates

```terraform
# iOS / iPadOS Enterprise Wi-Fi with Multiple Root Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_multiple_root_certificates_root_certificate" {
  display_name       = "example-ios-root"
  description        = "iOS / iPadOS enterprise wi-fi with multiple root certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_multiple_root_certificates_additional_root_certificate" {
  display_name       = "example-ios-second-root"
  description        = "iOS / iPadOS enterprise wi-fi with multiple root certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_multiple_root_certificates_scep_certificate" {
  display_name       = "example-ios-scep"
  description        = "iOS / iPadOS enterprise wi-fi with multiple root certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.iosScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_multiple_root_certificates_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "ios_multiple_root_certificates" {
  display_name       = "iOS / iPadOS Enterprise Wi-Fi with Multiple Root Certificates"
  description        = "iOS / iPadOS enterprise wi-fi with multiple root certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                           = "#microsoft.graph.iosEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"            = null,
    "deviceManagementApplicabilityRuleOsVersion"            = null,
    "deviceManagementApplicabilityRuleDeviceMode"           = null,
    "networkName"                                           = "Provider Enterprise WiFi",
    "ssid"                                                  = "Provider-Unassigned-EAP",
    "connectAutomatically"                                  = false,
    "connectWhenNetworkNameIsHidden"                        = false,
    "wiFiSecurityType"                                      = "wpaEnterprise",
    "proxySettings"                                         = "none",
    "proxyManualAddress"                                    = null,
    "proxyManualPort"                                       = null,
    "proxyAutomaticConfigurationUrl"                        = null,
    "disableMacAddressRandomization"                        = null,
    "preSharedKey"                                          = null,
    "eapType"                                               = "eapTls",
    "eapFastConfiguration"                                  = null,
    "trustedServerCertificateNames"                         = [],
    "authenticationMethod"                                  = "certificate",
    "innerAuthenticationProtocolForEapTtls"                 = null,
    "outerIdentityPrivacyTemporaryValue"                    = null,
    "usernameFormatString"                                  = null,
    "passwordFormatString"                                  = null,
    "rootCertificatesForServerValidation@odata.bind"        = ["https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_multiple_root_certificates_root_certificate.id}')", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_multiple_root_certificates_additional_root_certificate.id}')"],
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.ios_multiple_root_certificates_scep_certificate.id}')"
  })
}
```

### macOS

#### Custom Configuration

```terraform
# macOS Custom Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_custom_configuration" {
  display_name       = "macOS Custom Configuration"
  description        = "macOS custom configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSCustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "payloadName"                                 = "Wire Study",
    "payloadFileName"                             = "wire-study.mobileconfig",
    "payload"                                     = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHBsaXN0IFBVQkxJQyAiLS8vQXBwbGUvL0RURCBQTElTVCAxLjAvL0VOIiAiaHR0cDovL3d3dy5hcHBsZS5jb20vRFREcy9Qcm9wZXJ0eUxpc3QtMS4wLmR0ZCI+CjxwbGlzdCB2ZXJzaW9uPSIxLjAiPgo8ZGljdD4KCTxrZXk+UGF5bG9hZENvbnRlbnQ8L2tleT4KCTxhcnJheS8+Cgk8a2V5PlBheWxvYWREaXNwbGF5TmFtZTwva2V5PgoJPHN0cmluZz5XaXJlIFN0dWR5PC9zdHJpbmc+Cgk8a2V5PlBheWxvYWRJZGVudGlmaWVyPC9rZXk+Cgk8c3RyaW5nPmNvbS5leGFtcGxlLmNvZGV4LndpcmUtc3R1ZHk8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFR5cGU8L2tleT4KCTxzdHJpbmc+Q29uZmlndXJhdGlvbjwvc3RyaW5nPgoJPGtleT5QYXlsb2FkVVVJRDwva2V5PgoJPHN0cmluZz45ZDBhYjBmOS04YzBjLTQwMzUtOGFlOC03OTA3ZDI2MTc3M2I8L3N0cmluZz4KCTxrZXk+UGF5bG9hZFZlcnNpb248L2tleT4KCTxpbnRlZ2VyPjE8L2ludGVnZXI+CjwvZGljdD4KPC9wbGlzdD4K",
    "deploymentChannel"                           = "deviceChannel"
  })
}
```

#### Trusted Root Certificate

```terraform
# macOS Trusted Root Certificate.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_trusted_root_certificate" {
  display_name       = "macOS Trusted Root Certificate"
  description        = "macOS trusted root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "deploymentChannel"                           = null
  })
}
```

#### SCEP Certificate

```terraform
# macOS SCEP Certificate.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_scep_certificate_root_certificate" {
  display_name       = "example-macos-root"
  description        = "macOS scep certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "deploymentChannel"                           = null
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_scep_certificate" {
  display_name       = "macOS SCEP Certificate"
  description        = "macOS scep certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = null,
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "allowAllAppsAccess"                 = null,
    "deploymentChannel"                  = null,
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_scep_certificate_root_certificate.id}')"
  })
}
```

#### Wi-Fi

```terraform
# macOS Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_wifi" {
  display_name       = "macOS Wi-Fi"
  description        = "macOS wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider WiFi Study",
    "ssid"                                        = "Provider-Unassigned-Test",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaPersonal",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "deploymentChannel"                           = null,
    "wifiRequirePhysicalMacAddressEnabled"        = null,
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026"
  })
}
```

#### Enterprise Wi-Fi with Certificates

```terraform
# macOS Enterprise Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_enterprise_wifi_root_certificate" {
  display_name       = "example-macos-root"
  description        = "macOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "deploymentChannel"                           = null
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_enterprise_wifi_scep_certificate" {
  display_name       = "example-macos-scep"
  description        = "macOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "allowAllAppsAccess"                 = null,
    "deploymentChannel"                  = null,
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_enterprise_wifi_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_enterprise_wifi" {
  display_name       = "macOS Enterprise Wi-Fi with Certificates"
  description        = "macOS enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSEnterpriseWiFiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "networkName"                                 = "Provider Enterprise WiFi",
    "ssid"                                        = "Provider-Unassigned-EAP",
    "connectAutomatically"                        = false,
    "connectWhenNetworkNameIsHidden"              = false,
    "wiFiSecurityType"                            = "wpaEnterprise",
    "proxySettings"                               = "none",
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "deploymentChannel"                           = null,
    "wifiRequirePhysicalMacAddressEnabled"        = null,
    "preSharedKey"                                = null,
    "eapType"                                     = "eapTls",
    "eapFastConfiguration"                        = null,
    "trustedServerCertificateNames"               = [],
    "authenticationMethod"                        = "certificate",
    "innerAuthenticationProtocolForEapTtls"       = null,
    "outerIdentityPrivacyTemporaryValue"          = null,
    "rootCertificatesForServerValidation@odata.bind" = [
      "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_enterprise_wifi_root_certificate.id}')"
    ],
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_enterprise_wifi_scep_certificate.id}')"
  })
}
```

#### VPN

```terraform
# macOS VPN.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_vpn_root_certificate" {
  display_name       = "example-macos-root"
  description        = "macOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSTrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "deploymentChannel"                           = null
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_vpn_scep_certificate" {
  display_name       = "example-macos-scep"
  description        = "macOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSScepCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "allowAllAppsAccess"                 = null,
    "deploymentChannel"                  = null,
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_vpn_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_vpn" {
  display_name       = "macOS VPN"
  description        = "macOS vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSVpnConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "connectionName"                              = "Provider VPN Test",
    "connectionType"                              = "ciscoAnyConnect",
    "loginGroupOrDomain"                          = null,
    "role"                                        = null,
    "realm"                                       = null,
    "identifier"                                  = null,
    "enableSplitTunneling"                        = false,
    "authenticationMethod"                        = "certificate",
    "enablePerApp"                                = null,
    "safariDomains"                               = [],
    "providerType"                                = null,
    "associatedDomains"                           = [],
    "excludedDomains"                             = [],
    "disableOnDemandUserOverride"                 = null,
    "disconnectOnIdle"                            = null,
    "disconnectOnIdleTimerInSeconds"              = null,
    "includeAllNetworks"                          = null,
    "excludeLocalNetworks"                        = null,
    "proxyServer"                                 = null,
    "optInToDeviceIdSharing"                      = null,
    "deploymentChannel"                           = null,
    "server" = {
      "description"     = "Provider unassigned test",
      "address"         = "vpn.example.invalid",
      "isDefaultServer" = true
    },
    "customData"                     = [],
    "customKeyValueData"             = [],
    "onDemandRules"                  = [],
    "identityCertificate@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.macos_vpn_scep_certificate.id}')"
  })
}
```

#### PKCS Certificate

```terraform
# macOS PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_pkcs_certificate" {
  display_name       = "macOS PKCS Certificate"
  description        = "macOS pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSPkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "certificationAuthority"                      = "ca.example.invalid",
    "certificationAuthorityName"                  = "Provider Test CA",
    "certificateTemplateName"                     = "ProviderTest",
    "subjectAlternativeNameFormatString"          = null,
    "subjectNameFormatString"                     = "CN={{DeviceId}}",
    "certificateStore"                            = "user",
    "allowAllAppsAccess"                          = null,
    "deploymentChannel"                           = null,
    "customSubjectAlternativeNames"               = []
  })
}
```

#### Application Preferences

```terraform
# macOS Application Preferences.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_preferences" {
  display_name       = "macOS Application Preferences"
  description        = "macOS application preferences example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSCustomAppConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "bundleId"                                    = "com.example.provider-wire-study",
    "fileName"                                    = "preferences.plist",
    "configurationXml"                            = "PD94bWwgdmVyc2lvbj0iMS4wIj8+PHBsaXN0IHZlcnNpb249IjEuMCI+PGRpY3Q+PGtleT5FeGFtcGxlRW5hYmxlZDwva2V5Pjx0cnVlLz48L2RpY3Q+PC9wbGlzdD4="
  })
}
```

#### Software Updates

```terraform
# macOS Software Updates.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "macos_software_updates" {
  display_name       = "macOS Software Updates"
  description        = "macOS software updates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.macOSSoftwareUpdateConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "criticalUpdateBehavior"                      = "notConfigured",
    "configDataUpdateBehavior"                    = "notConfigured",
    "firmwareUpdateBehavior"                      = "notConfigured",
    "allOtherUpdateBehavior"                      = "notConfigured",
    "updateScheduleType"                          = "alwaysUpdate",
    "updateTimeWindowUtcOffsetInMinutes"          = null,
    "maxUserDeferralsCount"                       = null,
    "priority"                                    = null,
    "customUpdateTimeWindows"                     = []
  })
}
```

### Windows

#### Custom Configuration

```terraform
# Windows Custom Configuration.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_custom_configuration" {
  display_name       = "Windows Custom Configuration"
  description        = "Windows custom configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingInteger",
        "displayName" = "Allow Camera",
        "description" = "Wire Study",
        "omaUri"      = "./Device/Vendor/MSFT/Policy/Config/Camera/AllowCamera",
        "value"       = 1
      }
    ]
  })
}
```

#### Trusted Root Certificate

```terraform
# Windows Trusted Root Certificate.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_trusted_root_certificate" {
  display_name       = "Windows Trusted Root Certificate"
  description        = "Windows trusted root certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows81TrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "destinationStore"                            = "computerCertStoreRoot"
  })
}
```

#### Encrypted OMA Setting

```terraform
# Windows Encrypted OMA Setting.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_encrypted_oma_setting" {
  display_name       = "Windows Encrypted OMA Setting"
  description        = "Windows encrypted oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingString",
        "displayName" = "Provider test string",
        "description" = "Disposable unassigned test",
        "omaUri"      = "./Device/Vendor/MSFT/Policy/Config/Experience/ConfigureWindowsSpotlightOnLockScreen",
        "value"       = "provider-test-string"
      }
    ]
  })
}
```

#### Integer OMA Setting

```terraform
# Windows Integer OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_integer" {
  display_name       = "Windows Integer OMA Setting"
  description        = "Windows integer oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingInteger",
        "displayName" = "Provider Integer",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Integer",
        "value"       = 42
      }
    ]
  })
}
```

#### Boolean OMA Setting

```terraform
# Windows Boolean OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_boolean" {
  display_name       = "Windows Boolean OMA Setting"
  description        = "Windows boolean oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingBoolean",
        "displayName" = "Provider Boolean",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Boolean",
        "value"       = true
      }
    ]
  })
}
```

#### Floating-point OMA Setting

```terraform
# Windows Floating-point OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_floating_point" {
  display_name       = "Windows Floating-point OMA Setting"
  description        = "Windows floating-point oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingFloatingPoint",
        "displayName" = "Provider FloatingPoint",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/FloatingPoint",
        "value"       = 1.5
      }
    ]
  })
}
```

#### Date/Time OMA Setting

```terraform
# Windows Date/Time OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_date_time" {
  display_name       = "Windows Date/Time OMA Setting"
  description        = "Windows date/time oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingDateTime",
        "displayName" = "Provider DateTime",
        "description" = "Disposable value test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/DateTime",
        "value"       = "2026-01-01T00:00:00Z"
      }
    ]
  })
}
```

#### String OMA Setting

```terraform
# Windows String OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_string" {
  display_name       = "Windows String OMA Setting"
  description        = "Windows string oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingString",
        "displayName" = "Provider String",
        "description" = "Disposable encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/String",
        "value"       = "Synthetic string"
      }
    ]
  })
}
```

#### Binary Base64 OMA Setting

```terraform
# Windows Binary Base64 OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_base64" {
  display_name       = "Windows Binary Base64 OMA Setting"
  description        = "Windows binary base64 oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingBase64",
        "displayName" = "Provider Base64",
        "description" = "Disposable encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Base64",
        "fileName"    = null,
        "value"       = "U3ludGhldGljIGJpbmFyeQ=="
      }
    ]
  })
}
```

#### Cleartext XML OMA Setting

```terraform
# Windows Cleartext XML OMA Setting.
# Replace the example OMA-URI with the CSP URI for your environment.
# XML values are cleartext; the provider Base64-encodes XML when sending it to Graph.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_xml" {
  display_name       = "Windows Cleartext XML OMA Setting"
  description        = "Windows cleartext xml oma setting example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingStringXml",
        "displayName" = "Provider StringXml",
        "description" = "Disposable encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/StringXml",
        "fileName"    = null,
        "value"       = "<test>\n  <name>é 日本語</name>\n</test>\n"
      }
    ]
  })
}
```

#### Mixed OMA Value Types

```terraform
# Windows Mixed OMA Value Types.
# Replace the example OMA-URI with the CSP URI for your environment.
# XML values are cleartext; the provider Base64-encodes XML when sending it to Graph.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_oma_mixed_values" {
  display_name       = "Windows Mixed OMA Value Types"
  description        = "Windows mixed oma value types example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10CustomConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "omaSettings" = [
      {
        "@odata.type" = "#microsoft.graph.omaSettingFloatingPoint",
        "displayName" = "Provider FloatingPoint",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/FloatingPoint",
        "value"       = 1.5
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingBoolean",
        "displayName" = "Provider Boolean",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Boolean",
        "value"       = false
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingInteger",
        "displayName" = "Provider Integer",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Integer",
        "value"       = 7
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingBase64",
        "displayName" = "Provider Base64",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/Base64",
        "value"       = "AAECA//+/Q==",
        "fileName"    = null
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingString",
        "displayName" = "Provider String",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/String",
        "value"       = "test"
      },
      {
        "@odata.type" = "#microsoft.graph.omaSettingStringXml",
        "displayName" = "Provider StringXml",
        "description" = "Mixed OMA encoding test",
        "omaUri"      = "./Device/Vendor/MSFT/Test/StringXml",
        "value"       = "<test>\n  <name>é 日本語</name>\n</test>\n",
        "fileName"    = null
      }
    ]
  })
}
```

#### Windows Update for Business

```terraform
# Windows Update for Business.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_update_for_business" {
  display_name       = "Windows Update for Business"
  description        = "Windows windows update for business example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsUpdateForBusinessConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "deliveryOptimizationMode"                    = "userDefined",
    "prereleaseFeatures"                          = "userDefined",
    "automaticUpdateMode"                         = "userDefined",
    "microsoftUpdateServiceAllowed"               = false,
    "driversExcluded"                             = true,
    "installationSchedule"                        = null,
    "qualityUpdatesDeferralPeriodInDays"          = 0,
    "featureUpdatesDeferralPeriodInDays"          = 7,
    "qualityUpdatesPaused"                        = false,
    "featureUpdatesPaused"                        = false,
    "businessReadyUpdatesOnly"                    = "userDefined",
    "skipChecksBeforeRestart"                     = false,
    "updateWeeks"                                 = null,
    "featureUpdatesRollbackWindowInDays"          = null,
    "engagedRestartDeadlineInDays"                = null,
    "engagedRestartSnoozeScheduleInDays"          = null,
    "engagedRestartTransitionScheduleInDays"      = null,
    "deadlineForFeatureUpdatesInDays"             = null,
    "deadlineForQualityUpdatesInDays"             = null,
    "deadlineGracePeriodInDays"                   = null,
    "postponeRebootUntilAfterDeadline"            = null,
    "autoRestartNotificationDismissal"            = "notConfigured",
    "scheduleRestartWarningInHours"               = null,
    "scheduleImminentRestartWarningInMinutes"     = null,
    "userPauseAccess"                             = "notConfigured",
    "userWindowsUpdateScanAccess"                 = "notConfigured",
    "updateNotificationLevel"                     = "notConfigured",
    "allowWindows11Upgrade"                       = false
  })
}
```

#### Wi-Fi

```terraform
# Windows Wi-Fi.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_wifi" {
  display_name       = "Windows Wi-Fi"
  description        = "Windows wi-fi example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windowsWifiConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "wifiSecurityType"                            = "wpaPersonal",
    "meteredConnectionLimit"                      = null,
    "ssid"                                        = "Provider-Unassigned-Test",
    "networkName"                                 = "Provider WiFi Study",
    "connectAutomatically"                        = false,
    "connectToPreferredNetwork"                   = null,
    "connectWhenNetworkNameIsHidden"              = null,
    "proxySetting"                                = null,
    "proxyManualAddress"                          = null,
    "proxyManualPort"                             = null,
    "proxyAutomaticConfigurationUrl"              = null,
    "forceFIPSCompliance"                         = null,
    "preSharedKey"                                = "Synthetic-Wire-Study-Only-2026"
  })
}
```

#### Enterprise Wi-Fi with Certificates

```terraform
# Windows Enterprise Wi-Fi with Certificates.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_enterprise_wifi_root_certificate" {
  display_name       = "example-windows-root"
  description        = "Windows enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows81TrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "destinationStore"                            = "computerCertStoreRoot"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_enterprise_wifi_scep_certificate" {
  display_name       = "example-windows-scep"
  description        = "Windows enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows81SCEPCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "keyStorageProvider"                          = "useSoftwareKsp",
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.windows_enterprise_wifi_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_enterprise_wifi" {
  display_name       = "Windows Enterprise Wi-Fi with Certificates"
  description        = "Windows enterprise wi-fi with certificates example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                  = "#microsoft.graph.windowsWifiEnterpriseEAPConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"   = null,
    "deviceManagementApplicabilityRuleOsVersion"   = null,
    "deviceManagementApplicabilityRuleDeviceMode"  = null,
    "wifiSecurityType"                             = "wpaEnterprise",
    "meteredConnectionLimit"                       = null,
    "ssid"                                         = "Provider-Unassigned-EAP",
    "networkName"                                  = "Provider Enterprise WiFi",
    "connectAutomatically"                         = false,
    "connectToPreferredNetwork"                    = null,
    "connectWhenNetworkNameIsHidden"               = null,
    "proxySetting"                                 = null,
    "proxyManualAddress"                           = null,
    "proxyManualPort"                              = null,
    "proxyAutomaticConfigurationUrl"               = null,
    "forceFIPSCompliance"                          = null,
    "preSharedKey"                                 = null,
    "networkSingleSignOn"                          = null,
    "maximumAuthenticationTimeoutInSeconds"        = null,
    "userBasedVirtualLan"                          = null,
    "promptForAdditionalAuthenticationCredentials" = null,
    "enablePairwiseMasterKeyCaching"               = null,
    "maximumPairwiseMasterKeyCacheTimeInMinutes"   = null,
    "maximumNumberOfPairwiseMasterKeysInCache"     = null,
    "enablePreAuthentication"                      = null,
    "maximumPreAuthenticationAttempts"             = null,
    "eapType"                                      = "eapTls",
    "trustedServerCertificateNames"                = [],
    "authenticationMethod"                         = "certificate",
    "innerAuthenticationProtocolForEAPTTLS"        = null,
    "outerIdentityPrivacyTemporaryValue"           = null,
    "requireCryptographicBinding"                  = null,
    "performServerValidation"                      = null,
    "disableUserPromptForServerValidation"         = null,
    "authenticationPeriodInSeconds"                = null,
    "authenticationRetryDelayPeriodInSeconds"      = null,
    "eapolStartPeriodInSeconds"                    = null,
    "maximumEAPOLStartMessages"                    = null,
    "maximumAuthenticationFailures"                = null,
    "cacheCredentials"                             = null,
    "authenticationType"                           = null,
    "rootCertificatesForServerValidation@odata.bind" = [
      "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.windows_enterprise_wifi_root_certificate.id}')"
    ],
    "identityCertificateForClientAuthentication@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.windows_enterprise_wifi_scep_certificate.id}')"
  })
}
```

#### VPN

```terraform
# Windows VPN.
# Replace the sample public certificate with your own root certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_vpn_root_certificate" {
  display_name       = "example-windows-root"
  description        = "Windows vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows81TrustedRootCertificate",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "trustedRootCertificate"                      = "MIIDFjCCAf6gAwIBAgICAhgwDQYJKoZIhvcNAQELBQAwKzEpMCcGA1UEAxMgVGVycmFmb3JtIE1pY3Jvc29mdDM2NSBUZXN0IFJvb3QwIBcNMjAwMTAxMDAwMDAwWhgPMjA1MDAxMDEwMDAwMDBaMCsxKTAnBgNVBAMTIFRlcnJhZm9ybSBNaWNyb3NvZnQzNjUgVGVzdCBSb290MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt0FBlUiuQItKihIPdimOfiHkajmQgMvC/tztuj1Sis7rPeChIMVbexXSNHUKBHST6ti6ck91yLy2kUll9HLlM8Y/GczP7CRaIf2xv+8/In7YEuhZJ2E+yA/q8ZQZYRcUXYpYWtuutWBjKv/G1oh2l5IzGVePAn7gRpm79gVDsyw/cuhNtj2CQAGcScUCP2Yee+Dgc0MLrFLCvNtrVc0iosWARBLL2OataXzrSmoSQm29TyJDsZ+Qr8kanHbQICUp2mmyMVKr2uxd0JzkJKxNTrNB8dVhRUytKHz3BE0QkKOz+RFqsLjixx8InoyhQLj6FK4mhu5KPP3P2kdFDQin4QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+orhr+LWDMSsmhI7Le0bRV984wEwDQYJKoZIhvcNAQELBQADggEBAImrXrkNyTG/UL9fe5jiLyYExt1CfZcz6T0SK6F7zAnvK2ggJfxK92d99u/v0x6OVj2rLNxYE8Uf3tO+2xjo9hwLsO5bXzD9pqc+UeVDGfNiSUWyt4bkvG7typSl0VUbSw1OgeVYI5Gr4byOGZ8ph9DOuc8beQSJ5VlwV+PROK6rtm+QJKiLNeQElZYbqrjgtQvw1P9lNocw0AFcoa19PKIjd8ARf9NHr3GRt1JE5EbCLsjZcwpIV3SKf3cDZ6YmyO4KHPwwDKu/4RMvG2cm+Omqa80/eIKdZYz2yg2kiABGT0B/t9ZcgOOyjNjVwIBB8Ed0nTu+CEc9SVPBZ6B2ITM=",
    "certFileName"                                = "wire-study.cer",
    "destinationStore"                            = "computerCertStoreRoot"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_vpn_scep_certificate" {
  display_name       = "example-windows-scep"
  description        = "Windows vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows81SCEPCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "keyStorageProvider"                          = "useSoftwareKsp",
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "scepServerUrls" = [
      "https://scep.example.invalid/certsrv/mscep/mscep.dll"
    ],
    "subjectNameFormatString"            = "CN={{DeviceId}}",
    "keyUsage"                           = "keyEncipherment,digitalSignature",
    "keySize"                            = "size2048",
    "hashAlgorithm"                      = "sha2",
    "subjectAlternativeNameFormatString" = null,
    "certificateStore"                   = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = [],
    "rootCertificate@odata.bind"    = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.windows_vpn_root_certificate.id}')"
  })
}

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_vpn" {
  display_name       = "Windows VPN"
  description        = "Windows vpn example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10VpnConfiguration",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "connectionName"                              = "Provider VPN Test",
    "customXml"                                   = null,
    "profileTarget"                               = null,
    "connectionType"                              = "ciscoAnyConnect",
    "enableSplitTunneling"                        = false,
    "enableAlwaysOn"                              = null,
    "enableDeviceTunnel"                          = null,
    "enableDnsRegistration"                       = null,
    "dnsSuffixes"                                 = [],
    "microsoftTunnelSiteId"                       = null,
    "authenticationMethod"                        = "certificate",
    "rememberUserCredentials"                     = false,
    "enableConditionalAccess"                     = false,
    "enableSingleSignOnWithAlternateCertificate"  = false,
    "singleSignOnEku"                             = null,
    "singleSignOnIssuerHash"                      = null,
    "eapXml"                                      = null,
    "proxyServer"                                 = null,
    "onlyAssociatedAppsCanUseConnection"          = null,
    "windowsInformationProtectionDomain"          = null,
    "trustedNetworkDomains"                       = [],
    "cryptographySuite"                           = null,
    "servers" = [
      {
        "description"     = "Provider unassigned test",
        "address"         = "vpn.example.invalid",
        "isDefaultServer" = true
      }
    ],
    "associatedApps"                 = [],
    "trafficRules"                   = [],
    "routes"                         = [],
    "dnsRules"                       = [],
    "identityCertificate@odata.bind" = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('${microsoft365_graph_beta_device_management_device_configuration_templates_json.windows_vpn_scep_certificate.id}')"
  })
}
```

#### PKCS Certificate

```terraform
# Windows PKCS Certificate.

resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "windows_pkcs_certificate" {
  display_name       = "Windows PKCS Certificate"
  description        = "Windows pkcs certificate example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.windows10PkcsCertificateProfile",
    "deviceManagementApplicabilityRuleOsEdition"  = null,
    "deviceManagementApplicabilityRuleOsVersion"  = null,
    "deviceManagementApplicabilityRuleDeviceMode" = null,
    "renewalThresholdPercentage"                  = 20,
    "keyStorageProvider"                          = "useSoftwareKsp",
    "subjectNameFormat"                           = "custom",
    "subjectAlternativeNameType"                  = "none",
    "certificateValidityPeriodValue"              = 1,
    "certificateValidityPeriodScale"              = "years",
    "certificationAuthority"                      = "ca.example.invalid",
    "certificationAuthorityName"                  = "Provider Test CA",
    "certificateTemplateName"                     = "ProviderTest",
    "subjectAlternativeNameFormatString"          = null,
    "subjectNameFormatString"                     = "CN={{DeviceId}}",
    "certificateStore"                            = "user",
    "extendedKeyUsages" = [
      {
        "name"             = "Client Authentication",
        "objectIdentifier" = "1.3.6.1.5.5.7.3.2"
      }
    ],
    "customSubjectAlternativeNames" = []
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Display name of the device configuration profile.
- `settings` (String) Complete writable template settings as a JSON object, including the root `@odata.type`. Use `jsonencode()`, a JSON heredoc, or `file("profile.json")`. Configure metadata through `display_name`, `description`, and `role_scope_tag_ids`; configure assignments through `assignments`. Exclude read-only response metadata (`id`, timestamps, `version`, `supportsScopeTags`) and OData response annotations. This endpoint uses a flat profile object. Unlike Settings Catalog, there is no nested `settings` collection. Include the complete writable settings returned by Graph, including default, null, and empty values. The provider reads the full remote settings; it does not project the response onto previously configured keys. Updates use PATCH, so use explicit API-supported reset values to clear settings. To change the root `@odata.type`, recreate the profile with Terraform’s `-replace` option. Nested OData discriminators and relationship bindings remain JSON. Use full Graph URLs for `@odata.bind` references; omit a binding to remove it. The provider reads certificate relationships separately and updates them through `$ref`. Windows encrypted OMA values are recovered through the plaintext endpoint; `omaSettingStringXml.value` accepts cleartext XML, which the provider Base64-encodes only when sending requests. Binary `omaSettingBase64` values and certificate content remain Base64-encoded. Wi-Fi pre-shared keys cannot be recovered on import and must be supplied again. Windows Update derived pause/rollback fields and OMA encryption metadata are read-only and must be excluded. Secret values are stored in Terraform state. The API determines which device configuration template types are supported; see the examples for tested iOS/iPadOS, macOS, Windows, and Android profiles.

### Optional

- `assignments` (Attributes Set) Assignments owned by this profile. Removing all entries unassigns the profile. Do not manage its assignments with another resource. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Description of the device configuration profile. Maximum length is 1500 characters.
- `role_scope_tag_ids` (Set of String) Set of Intune scope tag IDs. Defaults to the default scope tag, `0`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The Intune device configuration ID.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment. Required when filter_type is 'include' or 'exclude'. Should be omitted when filter_type is 'none'.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax. Match the resource label to the example you are using.

```shell
# {resource_id} is the device configuration profile ID.
terraform import microsoft365_graph_beta_device_management_device_configuration_templates_json.android_general_configuration {resource_id}
```
