resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "with_setup_assistant" {
  display_name                                                = "Test Setup Assistant Apple Configurator Enrollment Policy - Unique"
  description                                                = "Apple configurator enrollment policy with setup assistant authentication"
  requires_user_authentication                               = true
  enable_authentication_via_company_portal                  = false
  require_company_portal_on_setup_assistant_enrolled_devices = true

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}