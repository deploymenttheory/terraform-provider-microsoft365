resource "microsoft365_graph_beta_users_user" "maximal" {
  account_enabled = true

  // Identity
  display_name        = "acc-test-user-maximal"
  given_name          = "Maximal"
  surname             = "User"
  user_principal_name = "acc-test-user-maximal@deploymenttheory.com"
  preferred_language  = "en-US"
  password_policies   = "DisablePasswordExpiration"

  // Age and Consent (for minor users)
  age_group                  = "NotAdult"
  consent_provided_for_minor = "Granted"

  // Job Information
  job_title          = "Senior Developer"
  company_name       = "Deployment Theory"
  department         = "Engineering"
  employee_id        = "1234567890"
  employee_type      = "full time"
  employee_hire_date = "2025-11-21T00:00:00Z"
  office_location    = "Building A"
  manager_id         = microsoft365_graph_beta_users_user.dependency_user.id

  // Contact Information
  city            = "Redmond"
  state           = "WA"
  country         = "US"
  street_address  = "123 street"
  postal_code     = "98052"
  usage_location  = "US"
  business_phones = ["+1 425-555-0100"]
  mobile_phone    = "+1 425-555-0101"
  mail            = "acc-test-user-maximal@deploymenttheory.com"
  fax_number      = "+1 425-555-0102"
  mail_nickname   = "acc-test-user-maximal"
  other_mails     = ["acc-test-user-maximal2.other@deploymenttheory.com"]

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}
