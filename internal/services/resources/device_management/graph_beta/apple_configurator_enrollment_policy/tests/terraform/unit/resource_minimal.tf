resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "minimal" {
  display_name                                               = "Test Minimal Apple Configurator Enrollment Policy - Unique"
  description                                                = "Minimal apple configurator enrollment policy for unit testing"
  requires_user_authentication                               = false
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}