resource "microsoft365_graph_beta_users_user" "maximal_with_mail" {
  display_name        = "Maximal User With Mail"
  user_principal_name = "maximal.mail@deploymenttheory.com"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "WithMail"
  mail                = "maximal.mail@deploymenttheory.com"
  mail_nickname       = "maximal.mail"
  job_title           = "Senior Developer"
  department          = "Engineering"
  company_name        = "Contoso Ltd"
  office_location     = "Building A"
  city                = "Redmond"
  state               = "WA"
  country             = "US"
  postal_code         = "98052"
  usage_location      = "US"
  business_phones     = ["+1 425-555-0100"]
  mobile_phone        = "+1 425-555-0101"
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  identities = [
    {
      sign_in_type       = "emailAddress"
      issuer             = "DeploymentTheory.onmicrosoft.com"
      issuer_assigned_id = "maximal.mail@deploymenttheory.com"
    }
  ]
  other_mails          = ["maximal.mail.other@deploymenttheory.com"]
  show_in_address_list = true
}
