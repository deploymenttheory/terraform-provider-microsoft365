resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "downgrade_test" {
  name        = "acc-test-macos-ade-downgrade"
  description = "Initial maximal configuration for downgrade testing"

  requires_user_authentication                               = true
  require_company_portal_on_setup_assistant_enrolled_devices = true
  await_device_configured                                    = true

  admin_account = {
    create_local_admin_account   = true
    user_name                    = "localadmin"
    full_name                    = "Local Administrator"
    create_local_primary_account = false
  }

  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  restore_disabled = true
  siri_disabled    = true

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
