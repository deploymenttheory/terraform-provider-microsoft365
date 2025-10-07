resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "enroll_without_user_affinity" {
  display_name                                                = "acc-test-apple-configurator-enrollment-policy-enroll-without-user-affinity"
  description                                                = "apple configurator enrollment policy without user affinity"
  requires_user_authentication                               = false
  enable_authentication_via_company_portal                  = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  
  timeouts = {
    create = "10s"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}