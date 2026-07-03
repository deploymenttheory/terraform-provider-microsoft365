resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "error_description" {
  name                         = "unit-test-macos-ade-error-description"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = ""
  support_phone_number         = ""
  description                  = format("%01600d", 0)
}
