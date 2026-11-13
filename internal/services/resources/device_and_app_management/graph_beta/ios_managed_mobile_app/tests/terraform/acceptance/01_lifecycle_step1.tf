# Test 01 step 1: target one app on a new iOS app protection policy.

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_app_protection" "test" {
  display_name = "acc-test-ios-managed-mobile-app-lifecycle-${random_string.suffix.result}"
}

resource "microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app" "test" {
  managed_app_protection_id = microsoft365_graph_beta_device_and_app_management_ios_managed_app_protection.test.id
  mobile_app_identifier = {
    bundle_id = "com.example.acctest.lifecycle"
  }
}
