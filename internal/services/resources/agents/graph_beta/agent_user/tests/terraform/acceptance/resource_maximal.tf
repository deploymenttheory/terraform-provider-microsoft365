resource "random_string" "maximal_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "random_string" "dependency_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_agent_user" "dependency_user" {
  display_name        = "acc-test-dep-agent-user-${random_string.dependency_user_id.result}"
  user_principal_name = "acc-test-dep-agent-user-${random_string.dependency_user_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-dep-agent-user-${random_string.dependency_user_id.result}"
  account_enabled     = true
  identity_parent_id  = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}

resource "microsoft365_graph_beta_users_agent_user" "maximal" {
  account_enabled = true

  // Identity
  display_name        = "acc-test-agent-user-maximal-${random_string.maximal_user_id.result}"
  given_name          = "Maximal"
  surname             = "User"
  user_principal_name = "acc-test-agent-user-maximal-${random_string.maximal_user_id.result}@deploymenttheory.com"
  identity_parent_id  = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  preferred_language  = "en-US"
  password_policies   = "DisablePasswordExpiration"

  // Age and Consent (for minor users)
  age_group                  = "NotAdult"
  consent_provided_for_minor = "Granted"

  // Job Information
  job_title          = "Marketing Agent"
  company_name       = "Deployment Theory"
  department         = "Marketing"
  employee_id        = "1234567890"
  employee_type      = "full time"
  employee_hire_date = "2025-11-21T00:00:00Z"
  office_location    = "Building A"
  manager_id         = microsoft365_graph_beta_users_agent_user.dependency_user.id

  // Contact Information
  city            = "Redmond"
  state           = "WA"
  country         = "US"
  street_address  = "123 street"
  postal_code     = "98052"
  usage_location  = "US"
  business_phones = ["+1 425-555-0100"]
  mobile_phone    = "+1 425-555-0101"
  mail            = "acc-test-agent-user-maximal-${random_string.maximal_user_id.result}@deploymenttheory.com"
  fax_number      = "+1 425-555-0102"
  mail_nickname   = "acc-test-agent-user-maximal-${random_string.maximal_user_id.result}"
  other_mails     = ["acc-test-agent-user-maximal-${random_string.maximal_user_id.result}2.other@deploymenttheory.com"]
}
