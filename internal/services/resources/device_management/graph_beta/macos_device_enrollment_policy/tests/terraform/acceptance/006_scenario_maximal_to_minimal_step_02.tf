resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "downgrade_test" {
  name                         = "acc-test-macos-ade-downgrade-minimal"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = "IT Support"
  support_phone_number         = "+1-555-0100"

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
