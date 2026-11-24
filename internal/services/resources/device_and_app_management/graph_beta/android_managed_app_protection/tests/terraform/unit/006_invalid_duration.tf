resource "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection" "test_006" {
  display_name                      = "unit-test-android-managed-app-protection-invalid-duration"
  period_online_before_access_check = "PT1.5H"
}
