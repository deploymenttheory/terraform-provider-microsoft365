resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "error_admin_account" {
  name                         = "unit-test-macos-ade-error-admin-account"
  requires_user_authentication = false
  await_device_configured      = true
  support_department           = ""
  support_phone_number         = ""
}
