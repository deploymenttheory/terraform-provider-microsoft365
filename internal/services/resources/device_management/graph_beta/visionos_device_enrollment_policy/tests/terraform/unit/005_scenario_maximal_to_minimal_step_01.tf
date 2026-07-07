resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "downgrade_test" {
  name        = "unit-test-visionos-ade-downgrade"
  description = "Initial maximal configuration for downgrade testing"

  await_device_configured   = true
  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  passcode_disabled = true
  siri_disabled     = true

  role_scope_tag_ids = ["0", "1", "2"]
}
