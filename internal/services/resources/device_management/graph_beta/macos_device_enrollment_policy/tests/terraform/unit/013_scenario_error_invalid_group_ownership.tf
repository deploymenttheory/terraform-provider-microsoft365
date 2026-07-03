resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "error_group_ownership" {
  name                         = "unit-test-macos-ade-error-group-ownership"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = ""
  support_phone_number         = ""
  device_security_group        = "20000000-0000-0000-0000-000000000002"
}
