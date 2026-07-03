resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "minimal" {
  name                         = "acc-test-macos-ade-minimal"
  requires_user_authentication = false
  await_device_configured      = false
  support_department           = ""
  support_phone_number         = ""

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
