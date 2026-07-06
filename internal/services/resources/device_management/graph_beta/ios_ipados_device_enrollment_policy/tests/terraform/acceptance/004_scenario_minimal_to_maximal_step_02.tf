resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "update_test" {
  name        = "acc-test-ios-ade-update-updated"
  description = "Updated to maximal configuration"

  requires_user_authentication                       = true
  require_setup_assistant_with_modern_authentication = true
  await_final_configuration                          = true

  locked_enrollment_enabled = true

  device_name_template = "{{DEVICETYPE}}-{{SERIAL}}"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  restore_disabled  = true
  apple_id_disabled = true
  siri_disabled     = true
  passcode_disabled = true

  role_scope_tag_ids = ["0", "1"]
}
