# switch_a demotes itself, which is only valid once switch_b has been promoted - Graph has no
# unset action, so the depends_on orders switch_b's setDefaultProfile before switch_a's update.
resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "switch_a" {
  name                         = "unit-test-ios-ade-default-switch-a"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
  is_default_policy_assignment = false

  depends_on = [microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy.switch_b]
}

resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "switch_b" {
  name                         = "unit-test-ios-ade-default-switch-b"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
  is_default_policy_assignment = true
}
