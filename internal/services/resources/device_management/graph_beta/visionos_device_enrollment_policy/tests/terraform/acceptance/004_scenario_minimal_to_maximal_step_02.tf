resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "update_test" {
  name        = "acc-test-visionos-ade-update-updated"
  description = "Updated to maximal configuration"

  await_device_configured   = true
  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  passcode_disabled    = true
  apple_id_disabled    = true
  siri_disabled        = true
  tips_screen_disabled = true

  role_scope_tag_ids = ["0", "1"]
}
