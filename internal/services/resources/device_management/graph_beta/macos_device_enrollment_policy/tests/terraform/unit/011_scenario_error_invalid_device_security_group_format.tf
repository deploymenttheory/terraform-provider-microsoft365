resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "error_device_security_group" {
  name                         = "unit-test-macos-ade-error-device-security-group"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = ""
  support_phone_number         = ""
  device_security_group        = "not-a-guid"
}
