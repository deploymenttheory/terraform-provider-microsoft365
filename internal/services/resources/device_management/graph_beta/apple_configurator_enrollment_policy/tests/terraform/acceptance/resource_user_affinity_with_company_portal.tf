resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "user_affinity_with_company_portal" {
  display_name                                                = "acc-test-apple-configurator-enrollment-policy-user-affinity-with-company-portal"
  description                                                = "apple configurator enrollment policy with user affinity via company portal"
  requires_user_authentication                               = false
  enable_authentication_via_company_portal                  = true
  require_company_portal_on_setup_assistant_enrolled_devices = false
  
  timeouts = {
    create = "10s"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}