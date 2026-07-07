resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "minimal" {
  name                         = "unit-test-macos-ade-minimal"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
