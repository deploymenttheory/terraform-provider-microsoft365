resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "microsoft_outlook_maximal" {
  display_name         = "unit-test-android-managed-device-app-configuration-policy-microsoft-outlook-maximal"
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

