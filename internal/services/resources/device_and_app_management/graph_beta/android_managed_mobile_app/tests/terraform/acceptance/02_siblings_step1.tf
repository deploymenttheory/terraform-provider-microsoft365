# Test 02 step 1: two apps on the same policy, created concurrently.

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection" "test" {
  display_name = "acc-test-android-managed-mobile-app-siblings-${random_string.suffix.result}"
}

resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "first" {
  managed_app_protection_id = microsoft365_graph_beta_device_and_app_management_android_managed_app_protection.test.id
  mobile_app_identifier = {
    package_id = "com.example.acctest.first"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "second" {
  managed_app_protection_id = microsoft365_graph_beta_device_and_app_management_android_managed_app_protection.test.id
  mobile_app_identifier = {
    package_id = "com.example.acctest.second"
  }
}
