resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "maximal" {
  name        = "acc-test-ios-ade-maximal"
  description = "iOS/iPadOS ADE enrollment policy exercising the full settings tree"

  requires_user_authentication                        = true
  require_setup_assistant_with_modern_authentication  = true
  await_final_configuration                           = true

  locked_enrollment_enabled = true

  device_name_template         = "{{DEVICETYPE}}-{{SERIAL}}"
  cellular_data_activation_url = "http://activation.carrier.net"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  passcode_disabled                         = true
  location_services_disabled                = false
  restore_disabled                          = true
  apple_id_disabled                         = true
  terms_and_conditions_disabled             = false
  touch_id_disabled                         = false
  apple_pay_disabled                        = true
  siri_disabled                             = true
  diagnostics_disabled                      = true
  privacy_pane_disabled                     = true
  restore_from_android_disabled             = true
  imessage_and_facetime_disabled            = true
  screen_time_screen_disabled               = true
  sim_setup_screen_disabled                 = true
  software_update_screen_disabled           = false
  watch_migration_screen_disabled           = true
  appearance_screen_disabled                = false
  device_to_device_migration_disabled       = true
  restore_completed_screen_disabled         = true
  software_update_completed_screen_disabled = true
  get_started_screen_disabled               = false
  action_button_screen_disabled             = true
  safety_screen_disabled                    = true
  terms_of_address_screen_disabled          = true
  apple_intelligence_disabled               = false
  lockdown_mode_disabled                    = true
  app_store_disabled                        = false
  camera_button_screen_disabled             = true
  multitasking_screen_disabled              = true
  os_showcase_screen_disabled               = true
  safety_and_handling_screen_disabled       = true
  web_content_filtering_disabled            = true

  role_scope_tag_ids = ["0"]
}
