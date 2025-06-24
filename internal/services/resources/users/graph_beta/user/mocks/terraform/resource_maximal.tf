resource "microsoft365_graph_beta_users_user" "maximal" {
  display_name        = "Maximal User"
  user_principal_name = "maximal.user@contoso.com"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "User"
  mail                = "maximal.user@contoso.com"
  mail_nickname       = "maximal.user"
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
    password                          = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  identities = [
    {
      sign_in_type       = "emailAddress"
      issuer             = "contoso.com"
      issuer_assigned_id = "maximal.user@contoso.com"
    }
  ]
  other_mails     = ["maximal.user.other@contoso.com"]
  proxy_addresses = ["SMTP:maximal.user@contoso.com"]
  show_in_address_list = true
} 