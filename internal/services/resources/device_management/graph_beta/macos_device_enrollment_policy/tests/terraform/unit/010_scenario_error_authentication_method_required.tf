resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "error_authentication_method" {
  name                         = "unit-test-macos-ade-error-authentication-method"
  requires_user_authentication = true
  await_device_configured      = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
