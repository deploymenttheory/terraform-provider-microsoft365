resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "switch_a" {
  name                         = "unit-test-visionos-ade-default-switch-a"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
  is_default_policy_assignment = true
}

resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "switch_b" {
  name                         = "unit-test-visionos-ade-default-switch-b"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
