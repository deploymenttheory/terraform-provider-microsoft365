resource "microsoft365_graph_beta_device_management_ios_ipados_device_enrollment_policy" "unsupported_company_portal" {
  name                                     = "unit-test-ios-ade-company-portal-unsupported"
  requires_user_authentication             = true
  enable_authentication_via_company_portal = true
  support_department                       = "IT Support"
  support_phone_number                     = "+1-555-0100"
}
