---
page_title: "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Manages Android managed store app configurations in Microsoft Intune using the /deviceAppManagement/mobileAppConfigurations endpoint. Use app configuration policies in Microsoft Intune to provide custom configuration settings for Android apps from the managed Google Play store. These configuration settings allow an app to be customized based on the app supplier's direction using Android Enterprise managed configurations. Learn more here: https://learn.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android
---

# microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy (Resource)

Manages Android managed store app configurations in Microsoft Intune using the `/deviceAppManagement/mobileAppConfigurations` endpoint. Use app configuration policies in Microsoft Intune to provide custom configuration settings for Android apps from the managed Google Play store. These configuration settings allow an app to be customized based on the app supplier's direction using Android Enterprise managed configurations. Learn more here: https://learn.microsoft.com/en-us/mem/intune/apps/app-configuration-policies-use-android

## Microsoft Documentation

- [iosMobileAppConfiguration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreappconfiguration?view=graph-rest-beta)
- [Create AndroidManagedStoreAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-androidmanagedstoreappconfiguration-create?view=graph-rest-beta)
- [Update AndroidManagedStoreAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-androidmanagedstoreappconfiguration-update?view=graph-rest-beta)
- [Delete AndroidManagedStoreAppConfiguration](https://learn.microsoft.com/en-us/graph/api/intune-apps-androidmanagedstoreappconfiguration-delete?view=graph-rest-beta)
- [App Configuration Policies for Managed Android Devices](https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-android)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`, `DeviceManagementApps.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.36.0-alpha | Experimental | Initial release |

## Example Usage

### Minimal Configuration

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "minimal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-minimal"
  description          = "Acceptance test Android managed store app configuration"
  targeted_mobile_apps = ["9711516a-f6f8-4953-ad1f-45920ef34dda"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.officehubrow"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.officehubrow",
    "managedProperty" : [
      {
        "key" : "test.key",
        "valueString" : "test-value"
      }
    ]
  })
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true
}
```

### Microsoft Authenticator

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_authenticator_maximal" {
  display_name           = "acc-test-android-managed-device-app-configuration-policy-microsoft-authenticator-maximal"
  description            = ""
  targeted_mobile_apps   = ["0e8ea6ec-63bb-436f-a89e-3adc475eb628"]
  role_scope_tag_ids     = ["0"]
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true

  package_id = "app:com.azure.authenticator"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.azure.authenticator",
    "managedProperty" : [
      {
        "key" : "suppress_camera_consent",
        "valueBool" : false
      },
      {
        "key" : "preferred_auth_config",
        "valueString" : "thing"
      },
      {
        "key" : "sharedDeviceRegistrationToken",
        "valueString" : "thing"
      },
      {
        "key" : "sharedDeviceTenantId",
        "valueString" : "thing"
      },
      {
        "key" : "sharedDeviceRegistrationPrefillUpn",
        "valueString" : "thing"
      },
      {
        "key" : "sharedDeviceMode",
        "valueBool" : false
      }
    ]
  })

  permission_actions = [
    {
      permission = "android.permission-group.NEARBY_DEVICES"
      action     = "prompt"
    },
    {
      permission = "android.permission.NEARBY_WIFI_DEVICES"
      action     = "prompt"
    },
    {
      permission = "android.permission.BLUETOOTH_CONNECT"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_AUDIO"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_IMAGES"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_VIDEO"
      action     = "prompt"
    },
    {
      permission = "android.permission.POST_NOTIFICATIONS"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_EXTERNAL_STORAGE"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_EXTERNAL_STORAGE"
      action     = "autoDeny"
    },
    {
      permission = "android.permission.RECEIVE_MMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECEIVE_WAP_PUSH"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_SMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECEIVE_SMS"
      action     = "autoGrant"
    },
    {
      permission = "android.permission.SEND_SMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.BODY_SENSORS_BACKGROUND"
      action     = "prompt"
    },
    {
      permission = "android.permission.BODY_SENSORS"
      action     = "prompt"
    },
    {
      permission = "android.permission.PROCESS_OUTGOING_CALLS"
      action     = "prompt"
    },
    {
      permission = "android.permission.USE_SIP"
      action     = "prompt"
    },
    {
      permission = "android.permission.ADD_VOICEMAIL"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_CALL_LOG"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CALL_LOG"
      action     = "prompt"
    },
    {
      permission = "android.permission.CALL_PHONE"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_PHONE_STATE"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECORD_AUDIO"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_BACKGROUND_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_COARSE_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_FINE_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.GET_ACCOUNTS"
      action     = "autoGrant"
    },
    {
      permission = "android.permission.WRITE_CONTACTS"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CONTACTS"
      action     = "prompt"
    },
    {
      permission = "android.permission.CAMERA"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_CALENDAR"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CALENDAR"
      action     = "prompt"
    }
  ]
}
```

### Microsoft 365 Copilot

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_365_copilot_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-365-copilot-maximal"
  description          = ""
  targeted_mobile_apps = ["9711516a-f6f8-4953-ad1f-45920ef34dda"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.officehubrow"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.officehubrow",
    "managedProperty" : [
      {
        "key" : "com.microsoft.office.officemobile.BingChatEnterprise.IsAllowed",
        "valueBool" : true
      },
      {
        "key" : "com.microsoft.office.officemobile.TeamsApps.IsAllowed",
        "valueBool" : true
      },
      {
        "key" : "com.microsoft.office.NotesCreationEnabled",
        "valueBool" : true
      },
      {
        "key" : "com.microsoft.intune.mam.AllowedAccountUPNs",
        "valueString" : "{{EmployeeID}}"
      }
    ]
  })
  profile_applicability  = "androidWorkProfile"
  connected_apps_enabled = true

  permission_actions = [
    {
      permission = "android.permission-group.NEARBY_DEVICES"
      action     = "prompt"
    },
    {
      permission = "android.permission.NEARBY_WIFI_DEVICES"
      action     = "autoDeny"
    },
    {
      permission = "android.permission.BLUETOOTH_CONNECT"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_AUDIO"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_IMAGES"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_MEDIA_VIDEO"
      action     = "prompt"
    },
    {
      permission = "android.permission.POST_NOTIFICATIONS"
      action     = "autoGrant"
    },
    {
      permission = "android.permission.WRITE_EXTERNAL_STORAGE"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_EXTERNAL_STORAGE"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECEIVE_MMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECEIVE_WAP_PUSH"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_SMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECEIVE_SMS"
      action     = "autoGrant"
    },
    {
      permission = "android.permission.SEND_SMS"
      action     = "prompt"
    },
    {
      permission = "android.permission.BODY_SENSORS_BACKGROUND"
      action     = "prompt"
    },
    {
      permission = "android.permission.BODY_SENSORS"
      action     = "prompt"
    },
    {
      permission = "android.permission.PROCESS_OUTGOING_CALLS"
      action     = "prompt"
    },
    {
      permission = "android.permission.USE_SIP"
      action     = "prompt"
    },
    {
      permission = "android.permission.ADD_VOICEMAIL"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_CALL_LOG"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CALL_LOG"
      action     = "prompt"
    },
    {
      permission = "android.permission.CALL_PHONE"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_PHONE_STATE"
      action     = "prompt"
    },
    {
      permission = "android.permission.RECORD_AUDIO"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_BACKGROUND_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_COARSE_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.ACCESS_FINE_LOCATION"
      action     = "prompt"
    },
    {
      permission = "android.permission.GET_ACCOUNTS"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_CONTACTS"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CONTACTS"
      action     = "prompt"
    },
    {
      permission = "android.permission.CAMERA"
      action     = "prompt"
    },
    {
      permission = "android.permission.WRITE_CALENDAR"
      action     = "prompt"
    },
    {
      permission = "android.permission.READ_CALENDAR"
      action     = "prompt"
    }
  ]
}
```

### Managed Home Screen (Kiosk Mode)

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "managed_home_screen_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-managed-home-screen-maximal"
  description          = ""
  targeted_mobile_apps = ["86263173-c38d-491f-b090-b2a9dbfeb09f"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.launcher.enterprise"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.launcher.enterprise",
    "managedProperty" : [
      { "key" : "offline_work_time_before_required_sign_in", "valueInteger" : 60 },
      { "key" : "show_alarm_volume_control", "valueBool" : false },
      { "key" : "show_notification_volume_control", "valueBool" : false },
      { "key" : "show_ring_volume_control", "valueBool" : false },
      { "key" : "show_call_volume_control", "valueBool" : false },
      { "key" : "show_ringtone_selector", "valueBool" : false },
      { "key" : "show_autorotate_toggle", "valueBool" : false },
      { "key" : "show_brightness_slider", "valueBool" : false },
      { "key" : "show_adaptive_brightness_toggle", "valueBool" : false },
      { "key" : "minimum_inactive_time_before_session_pin_required", "valueInteger" : 0 },
      { "key" : "header_name_style", "valueString" : "Display Name" },
      { "key" : "header_secondary_element", "valueString" : "Serial Number" },
      { "key" : "header_primary_element", "valueString" : "Tenant Name" },
      { "key" : "fixed_time_to_give_user_notice", "valueInteger" : 30 },
      { "key" : "fixed_time_to_signout", "valueInteger" : 12 },
      { "key" : "enable_fixed_signout", "valueBool" : false },
      { "key" : "block_pinning_browser_web_pages_to_MHS", "valueBool" : false },
      { "key" : "amount_of_time_before_try_exit_PIN_again", "valueInteger" : 0 },
      { "key" : "max_number_of_attempts_for_exit_PIN", "valueInteger" : 0 },
      { "key" : "minimum_length_for_session_PIN", "valueInteger" : 1 },
      { "key" : "max_number_of_attempts_for_session_PIN", "valueInteger" : 0 },
      { "key" : "custom_privacy_statement_url", "valueString" : "thing" },
      { "key" : "custom_privacy_statement_title", "valueString" : "thing" },
      { "key" : "enable_language_setting", "valueBool" : false },
      { "key" : "enable_PIN_to_resume", "valueBool" : false },
      { "key" : "auto_signout_time_to_give_user_notice", "valueInteger" : 60 },
      { "key" : "inactive_time_to_signout", "valueInteger" : 300 },
      { "key" : "enable_auto_signout", "valueBool" : false },
      { "key" : "session_PIN_complexity", "valueString" : "simple" },
      { "key" : "enable_session_PIN", "valueBool" : false },
      { "key" : "signin_screen_branding_logo", "valueString" : "thing" },
      { "key" : "enable_corporate_logo", "valueBool" : true },
      { "key" : "signin_screen_wallpaper", "valueString" : "thing" },
      { "key" : "signin_type", "valueString" : "AAD" },
      { "key" : "enable_mhs_signin", "valueBool" : false },
      { "key" : "theme_color", "valueString" : "light" },
      { "key" : "max_absolute_time_outside_MHS", "valueInteger" : 600 },
      { "key" : "max_inactive_time_outside_MHS", "valueInteger" : 180 },
      { "key" : "enable_max_absolute_time_outside_MHS", "valueBool" : false },
      { "key" : "enable_max_inactive_time_outside_MHS", "valueBool" : false },
      { "key" : "enable_easy_access_debugmenu", "valueBool" : false },
      { "key" : "apps_in_folder_ordered_by_name", "valueBool" : true },
      { "key" : "app_order_enabled", "valueBool" : false },
      { "key" : "login_hint_text", "valueString" : "thing" },
      { "key" : "domain_name", "valueString" : "thing" },
      { "key" : "show_device_name", "valueBool" : false },
      { "key" : "show_device_info_setting", "valueBool" : false },
      { "key" : "show_volume_setting", "valueBool" : false },
      { "key" : "show_flashlight_setting", "valueBool" : false },
      { "key" : "show_bluetooth_setting", "valueBool" : false },
      { "key" : "show_managed_setting", "valueBool" : true },
      { "key" : "show_wifi_setting", "valueBool" : false },
      { "key" : "exit_lock_task_mode_code", "valueString" : "thing" },
      { "key" : "show_virtual_status_bar", "valueBool" : false },
      { "key" : "virtual_app_switcher_type", "valueString" : "float" },
      { "key" : "virtual_home_type", "valueString" : "swipe_up" },
      { "key" : "show_virtual_home", "valueBool" : false },
      { "key" : "media_detect_before_screen_saver", "valueBool" : true },
      { "key" : "inactive_time_to_show_screen_saver", "valueInteger" : 30 },
      { "key" : "screen_saver_show_time", "valueInteger" : 0 },
      { "key" : "screen_saver_image", "valueString" : "http://thing" },
      { "key" : "show_screen_saver", "valueBool" : false },
      { "key" : "enable_wifi_allowlist", "valueBool" : false },
      { "key" : "screen_orientation", "valueInteger" : 1 },
      { "key" : "app_folder_icon", "valueInteger" : 0 },
      { "key" : "icon_size", "valueInteger" : 2 },
      { "key" : "wallpaper", "valueString" : "default" },
      { "key" : "lock_home_screen", "valueBool" : true },
      { "key" : "show_notification_badge", "valueBool" : false },
      { "key" : "grid_size", "valueString" : "Auto" }
    ]
  })
  profile_applicability  = "default"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Defender for Endpoint

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_defender_antivirus_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-defender-antivirus-maximal"
  description          = ""
  targeted_mobile_apps = ["ada90f6b-fb7a-457f-ba46-153947294a3f"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.scmx"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.scmx",
    "managedProperty" : [
      { "key" : "TunnelGuestTenantId", "valueString" : "NONE" },
      { "key" : "EnableNonAPKFileScan", "valueInteger" : 0 },
      { "key" : "GlobalSecureAccessPrivateChannel", "valueInteger" : -1 },
      { "key" : "GlobalSecureAccessPA", "valueInteger" : 2 },
      { "key" : "EnableGSA", "valueInteger" : 0 },
      { "key" : "UserUPN", "valueString" : "NONE" },
      { "key" : "EnableLowTouchOnboarding", "valueInteger" : 0 },
      { "key" : "DisableSignOut", "valueInteger" : 1 },
      { "key" : "DefenderDeviceTag", "valueString" : "NONE" },
      { "key" : "DefenderNetworkProtectionPrivacy", "valueInteger" : 1 },
      { "key" : "DefenderCertificateDetection", "valueInteger" : 0 },
      { "key" : "DefenderOpenNetworkDetection", "valueInteger" : 0 },
      { "key" : "DefenderNetworkProtectionAutoRemediation", "valueInteger" : 1 },
      { "key" : "DefenderEndUserTrustFlowEnable", "valueInteger" : 0 },
      { "key" : "DefenderNetworkProtectionEnable", "valueInteger" : 0 },
      { "key" : "DefenderSendFeedback", "valueInteger" : 1 },
      { "key" : "DefenderAllowlistedCACertificates", "valueString" : "NONE" },
      { "key" : "DefenderTVMPrivacyMode-PP", "valueInteger" : 1 },
      { "key" : "DefenderTVMPrivacyMode", "valueInteger" : 0 },
      { "key" : "DefenderExcludeAppInReport-PP", "valueInteger" : 1 },
      { "key" : "DefenderExcludeAppInReport", "valueInteger" : 0 },
      { "key" : "DefenderExcludeURLInReport-PP", "valueInteger" : 1 },
      { "key" : "DefenderExcludeURLInReport", "valueInteger" : 0 },
      { "key" : "defendertoggle_PP", "valueInteger" : 0 },
      { "key" : "defendertoggle", "valueInteger" : 1 },
      { "key" : "vpn", "valueInteger" : 1 },
      { "key" : "antiphishing", "valueInteger" : 1 }
    ]
  })
  profile_applicability  = "default"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Edge Browser

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_edge_browser_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-edge-browser-maximal"
  description          = ""
  targeted_mobile_apps = ["57eb943a-1542-4c9b-b220-ed09b95a50a2"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.emmx"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.emmx",
    "managedProperty" : [
      { "key" : "EdgeShowBottomBarInKioskMode", "valueBool" : false },
      { "key" : "EdgeShowAddressBarInKioskMode", "valueBool" : false },
      { "key" : "EdgeRestoreBrowsingOption", "valueInteger" : 1 },
      { "key" : "EdgeNewTabPageLayoutUserSelectable", "valueBool" : true },
      { "key" : "EdgeNewTabPageLayoutCustom", "valueString" : "thing" },
      { "key" : "EdgeNewTabPageLayout", "valueString" : "thing" },
      { "key" : "EdgeMyApps", "valueBool" : false },
      { "key" : "EdgeLockedViewModeEnabled", "valueBool" : false },
      { "key" : "EdgeLockedViewModeAllowedActions", "valueString" : "thing" },
      { "key" : "EdgeEnableKioskMode", "valueBool" : false },
      { "key" : "EdgeDisabledFeatures", "valueString" : "thing" },
      { "key" : "EdgeDisableShareUsageData", "valueBool" : false },
      { "key" : "EdgeDisableShareBrowsingHistory", "valueBool" : true },
      { "key" : "EdgeDefaultHTTPS", "valueBool" : false },
      { "key" : "EdgeBrandLogo", "valueBool" : true },
      { "key" : "EdgeBrandColor", "valueBool" : false },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" },
      { "key" : "WebUsbBlockedForUrls", "valueString" : "thing" },
      { "key" : "WebUsbAskForUrls", "valueString" : "thing" },
      { "key" : "WebUsbAllowDevicesForUrls", "valueString" : "thing" },
      { "key" : "WebRtcUdpPortRange", "valueString" : "thing" },
      { "key" : "WebAuthenticationRemoteDesktopAllowedOrigins", "valueString" : "thing" },
      { "key" : "VirtualKeyboardResizesLayoutByDefault", "valueBool" : false },
      { "key" : "URLBlocklist", "valueString" : "thing" },
      { "key" : "URLAllowlist", "valueString" : "thing" },
      { "key" : "TranslateEnabled", "valueBool" : false },
      { "key" : "SyncDisabled", "valueBool" : false },
      { "key" : "SmartScreenEnabled", "valueBool" : false },
      { "key" : "ServiceWorkerToControlSrcdocIframeEnabled", "valueBool" : false },
      { "key" : "SensorsBlockedForUrls", "valueString" : "thing" },
      { "key" : "SensorsAllowedForUrls", "valueString" : "thing" },
      { "key" : "SearchSuggestEnabled", "valueBool" : false },
      { "key" : "ScrollToTextFragmentEnabled", "valueBool" : false },
      { "key" : "SavingBrowserHistoryDisabled", "valueBool" : false },
      { "key" : "SSLErrorOverrideAllowedForOrigins", "valueString" : "thing" },
      { "key" : "SSLErrorOverrideAllowed", "valueBool" : false },
      { "key" : "RelatedWebsiteSetsOverrides", "valueString" : "thing" },
      { "key" : "RelatedWebsiteSetsEnabled", "valueBool" : false },
      { "key" : "ReduceAcceptLanguageEnabled", "valueBool" : false },
      { "key" : "QRCodeGeneratorEnabled", "valueBool" : false },
      { "key" : "ProxySettings", "valueString" : "thing" },
      { "key" : "PromptForDownloadLocation", "valueBool" : false },
      { "key" : "PrintingEnabled", "valueBool" : false },
      { "key" : "PostQuantumKeyAgreementEnabled", "valueBool" : false },
      { "key" : "PopupsBlockedForUrls", "valueString" : "thing" },
      { "key" : "PopupsAllowedForUrls", "valueString" : "thing" },
      { "key" : "PaymentMethodQueryEnabled", "valueBool" : false },
      { "key" : "PasswordManagerEnabled", "valueBool" : false },
      { "key" : "PartitionedBlobUrlUsage", "valueBool" : false },
      { "key" : "OverrideSecurityRestrictionsOnInsecureOrigin", "valueString" : "thing" },
      { "key" : "OverlayPermissionDetectionEnabled", "valueBool" : false },
      { "key" : "NtlmV2Enabled", "valueBool" : false },
      { "key" : "NewTabPageLocation", "valueString" : "thing" },
      { "key" : "NetworkPredictionOptions", "valueString" : "0" },
      { "key" : "MobileSiteForceForUrls", "valueString" : "thing" },
      { "key" : "ManagedFavorites", "valueString" : "thing" },
      { "key" : "JavaScriptOptimizerBlockedForSites", "valueString" : "thing" },
      { "key" : "JavaScriptOptimizerAllowedForSites", "valueString" : "thing" },
      { "key" : "JavaScriptJitBlockedForSites", "valueString" : "thing" },
      { "key" : "JavaScriptJitAllowedForSites", "valueString" : "thing" },
      { "key" : "JavaScriptBlockedForUrls", "valueString" : "thing" },
      { "key" : "JavaScriptAllowedForUrls", "valueString" : "thing" },
      { "key" : "InPrivateModeAvailability", "valueString" : "0" },
      { "key" : "ImportPasswordsDisabled", "valueBool" : false },
      { "key" : "IdleTimeoutActions", "valueStringArray" : ["close_browsers", "clear_browsing_history", "clear_download_history", "clear_cookies_and_other_site_data", "clear_cached_images_and_files", "clear_password_signin", "clear_autofill", "clear_site_settings", "reload_pages", "close_tabs"] },
      { "key" : "IdleTimeout", "valueInteger" : 1 },
      { "key" : "IPv6ReachabilityOverrideEnabled", "valueBool" : false },
      { "key" : "HttpsOnlyMode", "valueString" : "allowed" },
      { "key" : "HttpAllowlist", "valueString" : "thing" },
      { "key" : "HideFirstRunExperience", "valueBool" : false },
      { "key" : "HappyEyeballsV3Enabled", "valueBool" : false },
      { "key" : "HSTSPolicyBypassList", "valueString" : "thing" },
      { "key" : "ForceYouTubeRestrict", "valueString" : "1" },
      { "key" : "ForcePermissionPolicyUnloadDefaultEnabled", "valueBool" : false },
      { "key" : "ForceGoogleSafeSearch", "valueBool" : false },
      { "key" : "ExtensionSettings", "valueString" : "thing" },
      { "key" : "ExperimentationAndConfigurationServiceControl", "valueString" : "0" },
      { "key" : "EncryptedClientHelloEnabled", "valueBool" : false },
      { "key" : "EnableMediaRouter", "valueBool" : false },
      { "key" : "EditFavoritesEnabled", "valueBool" : false },
      { "key" : "EdgeSharedDeviceSupportEnabled", "valueBool" : false },
      { "key" : "EdgeOneAuthProxy", "valueString" : "thing" },
      { "key" : "EdgeCopilotEnabled", "valueBool" : false },
      { "key" : "EdgeBlockSignInEnabled", "valueBool" : false },
      { "key" : "EdgeAllowedAccountUPN", "valueString" : "thing" },
      { "key" : "EdgeAllowedAccountOnly", "valueBool" : false },
      { "key" : "DownloadRestrictions", "valueString" : "1" },
      { "key" : "DisabledMiniApps", "valueStringArray" : ["featured_feature", "sydney_chat", "download_manager", "pdf_reader", "news_feed", "translator", "converter", "money", "images", "covid", "games", "deals", "wallet", "weather", "gaokao", "weather_etree", "weather_map", "weather_life"] },
      { "key" : "DisableAuthNegotiateCnameLookup", "valueBool" : false },
      { "key" : "DesktopSiteForceForUrls", "valueString" : "thing" },
      { "key" : "DefaultWebUsbGuardSetting", "valueString" : "2" },
      { "key" : "DefaultWebBluetoothGuardSetting", "valueString" : "3" },
      { "key" : "DefaultSensorsSetting", "valueString" : "1" },
      { "key" : "DefaultSearchProviderSuggestURL", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderSearchURL", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderName", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderImageURLPostParams", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderImageURL", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderEncodings", "valueString" : "thing" },
      { "key" : "DefaultSearchProviderEnabled", "valueBool" : false },
      { "key" : "DefaultPopupsSetting", "valueString" : "1" },
      { "key" : "DefaultJavaScriptSetting", "valueString" : "1" },
      { "key" : "DefaultJavaScriptOptimizerSetting", "valueString" : "1" },
      { "key" : "DefaultJavaScriptJitSetting", "valueString" : "2" },
      { "key" : "DefaultGeolocationSetting", "valueString" : "2" },
      { "key" : "DefaultDesktopSiteSetting", "valueString" : "2" },
      { "key" : "DefaultCookiesSetting", "valueString" : "4" },
      { "key" : "DefaultBrowserSettingEnabled", "valueBool" : false },
      { "key" : "DataUrlInSvgUseEnabled", "valueBool" : false },
      { "key" : "DataURLWhitespacePreservationEnabled", "valueBool" : false },
      { "key" : "CookiesSessionOnlyForUrls", "valueString" : "thing" },
      { "key" : "CookiesBlockedForUrls", "valueString" : "thing" },
      { "key" : "CookiesAllowedForUrls", "valueString" : "thing" },
      { "key" : "CertificateTransparencyEnforcementDisabledForUrls", "valueString" : "thing" },
      { "key" : "CertificateTransparencyEnforcementDisabledForCas", "valueString" : "thing" },
      { "key" : "CORSNonWildcardRequestHeadersSupport", "valueBool" : false },
      { "key" : "CAPlatformIntegrationEnabled", "valueBool" : false },
      { "key" : "CAHintCertificates", "valueString" : "thing" },
      { "key" : "CADistrustedCertificates", "valueString" : "thing" },
      { "key" : "CACertificatesWithConstraints", "valueString" : "thing" },
      { "key" : "CACertificates", "valueString" : "thing" },
      { "key" : "BuiltInDnsClientEnabled", "valueBool" : false },
      { "key" : "BlockThirdPartyCookies", "valueBool" : false },
      { "key" : "BiometricAuthenticationBeforeFilling", "valueBool" : false },
      { "key" : "BackForwardCacheEnabled", "valueBool" : false },
      { "key" : "AutofillCreditCardEnabled", "valueBool" : false },
      { "key" : "AutofillAddressEnabled", "valueBool" : false },
      { "key" : "AuthServerAllowlist", "valueString" : "thing" },
      { "key" : "AuthSchemes", "valueString" : "thing" },
      { "key" : "AuthNegotiateDelegateAllowlist", "valueString" : "thing" },
      { "key" : "AuthAndroidNegotiateAccountType", "valueString" : "thing" },
      { "key" : "AlternateErrorPagesEnabled", "valueBool" : false },
      { "key" : "AllowedDomainsForApps", "valueString" : "thing" },
      { "key" : "AllowWebAuthnWithBrokenTlsCerts", "valueBool" : false },
      { "key" : "AllowBackForwardCacheForCacheControlNoStorePageEnabled", "valueBool" : false },
      { "key" : "AllHttpAuthSchemesAllowedForOrigins", "valueString" : "thing" },
      { "key" : "AddressBarWorkSearchResultsEnabled", "valueBool" : false },
      { "key" : "AdditionalDnsQueryTypesEnabled", "valueBool" : false },
      { "key" : "AccessibilityPerformanceFilteringAllowed", "valueBool" : false },
      { "key" : "AccessControlAllowMethodsInCORSPreflightSpecConformant", "valueBool" : false }
    ]
  })
  profile_applicability  = "androidWorkProfile"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Excel

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_excel_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-excel-maximal"
  description          = ""
  targeted_mobile_apps = ["b50fb1cb-98c3-4c35-8de9-e37a5d71fa36"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.excel"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.excel",
    "managedProperty" : [
      { "key" : "com.microsoft.office.officemobile.BingChatEnterprise.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.officemobile.TeamsApps.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.NotesCreationEnabled", "valueBool" : true },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft OneDrive

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_onedrive_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-onedrive-maximal"
  description          = ""
  targeted_mobile_apps = ["970c9b4a-4879-4b6b-985e-693167bff8f6"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.skydrive"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.skydrive",
    "managedProperty" : [
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft OneNote

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_onenote_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-onenote-maximal"
  description          = ""
  targeted_mobile_apps = ["f400d6e7-08de-4267-8db2-035253751022"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.onenote"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.onenote",
    "managedProperty" : [
      { "key" : "com.microsoft.office.officemobile.BingChatEnterprise.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.officemobile.TeamsApps.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.NotesCreationEnabled", "valueBool" : true },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidWorkProfile"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Outlook

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_outlook_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-outlook-maximal"
  description          = ""
  targeted_mobile_apps = ["df3baafe-df9e-43c7-9bda-8c59f0e9c2ed"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.outlook"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.outlook",
    "managedProperty" : [
      { "key" : "com.microsoft.outlook.EmailProfile.AccountType", "valueString" : "ModernAuth" },
      { "key" : "com.microsoft.outlook.EmailProfile.EmailUPN", "valueString" : "{{userprincipalname}}" },
      { "key" : "com.microsoft.outlook.EmailProfile.EmailAddress", "valueString" : "{{userprincipalname}}" },
      { "key" : "IntuneMAMAllowedAccountsOnly", "valueString" : "Enabled" },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "{{userprincipalname}}" },
      { "key" : "com.microsoft.outlook.Mail.FocusedInbox", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Contacts.LocalSyncEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.OfficeFeedEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SuggestedRepliesEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.ExternalRecipientsToolTipEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.DefaultSignatureEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.BlockExternalImagesEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.OrganizeByThreadEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.PlayMyEmailsEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Settings.ThemesEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Calendar.NativeSyncEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.TextPredictionsEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.UserChangeAllowed", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.SignAllMail.UserChangeAllowed", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.EncryptAllMail.UserChangeAllowed", "valueBool" : true },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.LDAPHostName", "valueString" : "http://some_url" },
      { "key" : "com.microsoft.outlook.Mail.SMIMEEnabled.CertsFromIntune", "valueBool" : false }
    ]
  })
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft PowerPoint

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_powerpoint_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-powerpoint-maximal"
  description          = ""
  targeted_mobile_apps = ["fc4750a7-5f21-450b-a33c-0a8640356144"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.powerpoint"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.powerpoint",
    "managedProperty" : [
      { "key" : "com.microsoft.office.officemobile.BingChatEnterprise.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.officemobile.TeamsApps.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.NotesCreationEnabled", "valueBool" : true },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidWorkProfile"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Teams

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_teams_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-teams-maximal"
  description          = ""
  targeted_mobile_apps = ["33029352-f792-4507-b963-ab2441a0c5f0"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.teams"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.teams",
    "managedProperty" : [
      { "key" : "enable_numeric_emp_id_keypad", "valueBool" : false },
      { "key" : "preferred_auth_config", "valueString" : "thing" },
      { "key" : "domain_name", "valueString" : "thing" },
      { "key" : "com.microsoft.teams.forceLoginUsingPassword", "valueBool" : false },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidWorkProfile"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

### Microsoft Word

```terraform
resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_word_maximal" {
  display_name         = "acc-test-android-managed-device-app-configuration-policy-microsoft-word-maximal"
  description          = ""
  targeted_mobile_apps = ["276a0f04-d4bf-4772-9de4-5927f1f9d5ca"]
  role_scope_tag_ids   = ["0"]

  package_id = "app:com.microsoft.office.word"
  payload_json = jsonencode({
    "kind" : "androidenterprise#managedConfiguration",
    "productId" : "app:com.microsoft.office.word",
    "managedProperty" : [
      { "key" : "com.microsoft.office.officemobile.BingChatEnterprise.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.officemobile.TeamsApps.IsAllowed", "valueBool" : true },
      { "key" : "com.microsoft.office.NotesCreationEnabled", "valueBool" : true },
      { "key" : "com.microsoft.intune.mam.AllowedAccountUPNs", "valueString" : "thing" }
    ]
  })
  profile_applicability  = "androidDeviceOwner"
  connected_apps_enabled = true

  permission_actions = [
    { permission = "android.permission-group.NEARBY_DEVICES", action = "prompt" },
    { permission = "android.permission.NEARBY_WIFI_DEVICES", action = "prompt" },
    { permission = "android.permission.BLUETOOTH_CONNECT", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_AUDIO", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_IMAGES", action = "prompt" },
    { permission = "android.permission.READ_MEDIA_VIDEO", action = "prompt" },
    { permission = "android.permission.POST_NOTIFICATIONS", action = "prompt" },
    { permission = "android.permission.WRITE_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.READ_EXTERNAL_STORAGE", action = "prompt" },
    { permission = "android.permission.RECEIVE_MMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_WAP_PUSH", action = "prompt" },
    { permission = "android.permission.READ_SMS", action = "prompt" },
    { permission = "android.permission.RECEIVE_SMS", action = "prompt" },
    { permission = "android.permission.SEND_SMS", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS_BACKGROUND", action = "prompt" },
    { permission = "android.permission.BODY_SENSORS", action = "prompt" },
    { permission = "android.permission.PROCESS_OUTGOING_CALLS", action = "prompt" },
    { permission = "android.permission.USE_SIP", action = "prompt" },
    { permission = "android.permission.ADD_VOICEMAIL", action = "prompt" },
    { permission = "android.permission.WRITE_CALL_LOG", action = "prompt" },
    { permission = "android.permission.READ_CALL_LOG", action = "prompt" },
    { permission = "android.permission.CALL_PHONE", action = "prompt" },
    { permission = "android.permission.READ_PHONE_STATE", action = "prompt" },
    { permission = "android.permission.RECORD_AUDIO", action = "prompt" },
    { permission = "android.permission.ACCESS_BACKGROUND_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_COARSE_LOCATION", action = "prompt" },
    { permission = "android.permission.ACCESS_FINE_LOCATION", action = "prompt" },
    { permission = "android.permission.GET_ACCOUNTS", action = "prompt" },
    { permission = "android.permission.WRITE_CONTACTS", action = "prompt" },
    { permission = "android.permission.READ_CONTACTS", action = "prompt" },
    { permission = "android.permission.CAMERA", action = "prompt" },
    { permission = "android.permission.WRITE_CALENDAR", action = "prompt" },
    { permission = "android.permission.READ_CALENDAR", action = "prompt" }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the Android mobile app configuration
- `package_id` (String) The package ID of the Android app (e.g., `app:com.microsoft.office.officehubrow`).
- `payload_json` (String, Sensitive) The Android Enterprise managed configuration in Base64 encoded JSON format.
- `targeted_mobile_apps` (Set of String) Set of Android mobile app IDs that this configuration targets.

### Optional

- `connected_apps_enabled` (Boolean) Whether connected apps are enabled for this configuration.
- `description` (String) The optional description of the Android mobile app configuration
- `permission_actions` (Attributes Set) List of Android permissions and their corresponding actions.Specify permissions you want to override.If they are not chosen/specified explicitly, then the default behavior will apply. Learn more here: https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-android (see [below for nested schema](#nestedatt--permission_actions))
- `profile_applicability` (String) The profile applicability for this configuration. Possible values: `default`, `androidWorkProfile`, `androidDeviceOwner`. Defaults to `default`.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Android mobile app configuration.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `app_supports_oem_config` (Boolean) Whether the app supports OEM configuration. This is a computed value from the API.
- `id` (String) The unique identifier for this Android mobile app configuration
- `version` (Number) Version of the Android mobile app configuration.

<a id="nestedatt--permission_actions"></a>
### Nested Schema for `permission_actions`

Required:

- `action` (String) The action for this permission. Possible values: `prompt`, `autoGrant`, `autoDeny`
- `permission` (String) The Android permission (e.g., `android.permission.CAMERA`)


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using ID format: {id}
terraform import microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.example 00000000-0000-0000-0000-000000000000
```