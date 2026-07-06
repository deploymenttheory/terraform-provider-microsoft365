resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "downgrade_test" {
  name        = "acc-test-ios-ade-downgrade"
  description = "Initial maximal configuration for downgrade testing"

  requires_user_authentication             = true
  enable_authentication_via_company_portal = true

  locked_enrollment_enabled = true

  device_name_template = "{{DEVICETYPE}}-{{SERIAL}}"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  restore_disabled = true
  siri_disabled    = true

  role_scope_tag_ids = ["0", "1", "2"]
}
