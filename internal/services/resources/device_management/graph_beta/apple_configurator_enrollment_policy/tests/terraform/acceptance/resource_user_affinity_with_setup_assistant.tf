resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "user_affinity_with_setup_assistant" {
  display_name                                               = "acc-test-apple-configurator-enrollment-policy-user-affinity-with-setup-assistant"
  description                                                = "apple configurator enrollment policy with user affinity via setup assistant"
  requires_user_authentication                               = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = true

  timeouts = {
    create = "10s"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}