resource "microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy" "minimal" {
  display_name         = "unit-test-android-managed-device-app-configuration-policy-minimal"
  description          = "Unit test Android managed store app configuration"
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

