resource "microsoft365_graph_beta_users_user" "maximal_with_proxy" {
  display_name        = "Maximal User With Proxy"
  user_principal_name = "maximal.proxy@deploymenttheory.com"
  account_enabled     = true
  given_name          = "Maximal"
  surname             = "WithProxy"
  mail_nickname       = "maximal.proxy"
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
      issuer_assigned_id = "maximal.proxy@deploymenttheory.com"
    }
  ]
  other_mails          = ["maximal.proxy.other@deploymenttheory.com"]
  proxy_addresses      = ["SMTP:maximal.proxy@deploymenttheory.com"]
  show_in_address_list = true
}
