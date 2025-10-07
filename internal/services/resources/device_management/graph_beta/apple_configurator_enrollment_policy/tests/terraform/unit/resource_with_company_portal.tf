resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "with_company_portal" {
  display_name                                                = "Test Company Portal Apple Configurator Enrollment Policy - Unique"
  description                                                = "Apple configurator enrollment policy with company portal authentication"
  requires_user_authentication                               = false
  enable_authentication_via_company_portal                  = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}