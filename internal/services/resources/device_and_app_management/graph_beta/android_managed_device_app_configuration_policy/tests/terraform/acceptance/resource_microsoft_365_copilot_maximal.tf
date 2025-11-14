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

