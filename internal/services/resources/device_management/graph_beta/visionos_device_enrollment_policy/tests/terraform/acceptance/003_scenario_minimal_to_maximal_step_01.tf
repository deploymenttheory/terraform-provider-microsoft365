resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "update_test" {
  name                         = "acc-test-visionos-ade-update"
  requires_user_authentication = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"
}
