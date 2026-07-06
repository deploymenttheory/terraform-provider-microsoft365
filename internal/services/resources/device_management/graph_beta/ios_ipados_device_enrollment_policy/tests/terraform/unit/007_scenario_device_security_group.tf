resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "device_group" {
  name                         = "unit-test-ios-ade-device-group"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
  device_security_group        = "10000000-0000-0000-0000-000000000001"
}
