resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "maximal" {
  name        = "acc-test-visionos-ade-maximal"
  description = "visionOS ADE enrollment policy exercising the full settings tree"

  await_device_configured   = true
  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  apple_id_disabled               = true
  apple_pay_disabled              = true
  diagnostics_disabled            = true
  get_started_screen_disabled     = false
  apple_intelligence_disabled     = false
  location_services_disabled      = false
  passcode_disabled               = true
  privacy_pane_disabled           = true
  screen_time_screen_disabled     = true
  siri_disabled                   = true
  software_update_screen_disabled = false
  terms_and_conditions_disabled   = false
  tips_screen_disabled            = true
  touch_id_disabled               = false

  role_scope_tag_ids = ["0"]
}
