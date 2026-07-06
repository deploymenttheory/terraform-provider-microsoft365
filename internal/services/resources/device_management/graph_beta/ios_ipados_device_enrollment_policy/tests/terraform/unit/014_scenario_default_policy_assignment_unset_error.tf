resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "default_assignment" {
  name                         = "unit-test-ios-ade-default-assignment"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
  is_default_policy_assignment = false
}
