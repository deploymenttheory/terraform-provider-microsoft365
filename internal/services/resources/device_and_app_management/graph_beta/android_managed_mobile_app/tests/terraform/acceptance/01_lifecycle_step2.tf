# Test 01 step 2: a different package targets a different app, so the resource is replaced.

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection" "test" {
  display_name = "acc-test-android-managed-mobile-app-lifecycle-${random_string.suffix.result}"
}

resource "microsoft365_graph_beta_device_and_app_management_android_managed_mobile_app" "test" {
  managed_app_protection_id = microsoft365_graph_beta_device_and_app_management_android_managed_app_protection.test.id
  mobile_app_identifier = {
    package_id = "com.example.acctest.renamed"
  }
}
