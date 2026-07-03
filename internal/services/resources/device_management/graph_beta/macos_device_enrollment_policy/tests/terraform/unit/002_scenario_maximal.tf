resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "maximal" {
  name        = "unit-test-macos-ade-maximal"
  description = "macOS ADE enrollment policy exercising the full settings tree"

  requires_user_authentication                               = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = true

  await_device_configured = true

  admin_account = {
    create_local_admin_account = true
    user_name                  = "localadmin"
    full_name                  = "Local Administrator"
    hide_account               = true
    password_rotation_in_days  = 90

    create_local_primary_account = true
    primary_account = {
      prefill_account_info = true
      restrict_editing     = true
      user_name            = "primaryuser"
      full_name            = "Primary User"
    }
  }

  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  location_services_disabled                = false
  restore_disabled                          = true
  apple_id_disabled                         = true
  terms_and_conditions_disabled             = false
  touch_id_disabled                         = false
  apple_pay_disabled                        = true
  siri_disabled                             = true
  diagnostics_disabled                      = true
  file_vault_disabled                       = false
  icloud_diagnostics_disabled               = true
  icloud_storage_disabled                   = true
  display_tone_setup_disabled               = false
  screen_time_screen_disabled               = true
  privacy_pane_disabled                     = true
  accessibility_screen_disabled             = false
  auto_unlock_with_watch_disabled           = true
  lockdown_mode_disabled                    = true
  software_update_screen_disabled           = false
  software_update_completed_screen_disabled = true
  terms_of_address_screen_disabled          = true
  apple_intelligence_disabled               = false
  os_showcase_screen_disabled               = true
  app_store_disabled                        = false

  role_scope_tag_ids = ["0"]
}
