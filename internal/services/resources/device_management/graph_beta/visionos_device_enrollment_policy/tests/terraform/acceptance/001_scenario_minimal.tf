resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "minimal" {
  name                         = "acc-test-visionos-ade-minimal"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
