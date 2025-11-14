resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_excel_maximal" {
  display_name         = "unit-test-android-managed-device-app-configuration-policy-microsoft-excel-maximal"
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

