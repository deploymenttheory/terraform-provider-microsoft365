resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "downgrade_test" {
  name        = "unit-test-ios-ade-downgrade"
  description = "Initial maximal configuration for downgrade testing"

  requires_user_authentication                       = true
  require_setup_assistant_with_modern_authentication = true

  locked_enrollment_enabled = true

  device_name_template = "{{DEVICETYPE}}-{{SERIAL}}"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  restore_disabled = true
  siri_disabled    = true

  role_scope_tag_ids                = ["0", "1", "2"]
  accessibility_appearance_disabled = true
  liquid_glass_disabled             = true
  device_features_tour_disabled     = true
}
