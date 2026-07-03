resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "downgrade_test" {
  name                         = "unit-test-macos-ade-downgrade-minimal"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = ""
  support_phone_number         = ""
}
