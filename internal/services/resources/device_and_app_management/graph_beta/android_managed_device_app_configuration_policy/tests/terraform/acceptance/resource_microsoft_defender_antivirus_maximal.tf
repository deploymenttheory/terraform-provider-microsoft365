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

