resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "update_test" {
  name        = "acc-test-macos-ade-update-updated"
  description = "Updated to maximal configuration"

  requires_user_authentication                               = true
  require_company_portal_on_setup_assistant_enrolled_devices = true
  await_device_configured                                    = true

  admin_account = {
    create_local_admin_account = true
    user_name                  = "localadmin"
    full_name                  = "Local Administrator"

    create_local_primary_account = true
    primary_account = {
      prefill_account_info = true
      user_name            = "primaryuser"
      full_name            = "Primary User"
    }
  }

  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  restore_disabled    = true
  apple_id_disabled   = true
  siri_disabled       = true
  file_vault_disabled = true

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
