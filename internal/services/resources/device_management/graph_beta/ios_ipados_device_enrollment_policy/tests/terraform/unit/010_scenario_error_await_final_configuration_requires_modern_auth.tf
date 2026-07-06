resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "error_await_final_configuration" {
  name                         = "unit-test-ios-ade-error-await-final-configuration"
  requires_user_authentication = true
  await_final_configuration    = true
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
