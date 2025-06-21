# Minimal example with only required properties
resource "microsoft365_graph_beta_users_user" "minimal_example" {
  display_name        = "John Doe"
  account_enabled     = true
  user_principal_name = "john.doe@contoso.com"
  mail_nickname       = "johndoe"
  password_profile = {
    password                           = "SecurePassword123!"
    force_change_password_next_sign_in = true
  }
}

# Comprehensive example with most available properties
resource "microsoft365_graph_beta_users_user" "full_example" {
  # Required properties
  display_name        = "Jane Smith"
  account_enabled     = true
  user_principal_name = "jane.smith@contoso.com"
  mail_nickname       = "janesmith"

  # Password configuration
  password_profile = {
    password                                    = "VerySecurePassword456!"
    force_change_password_next_sign_in          = true
    force_change_password_next_sign_in_with_mfa = false
  }

  # Personal information
  given_name     = "Jane"
  surname        = "Smith"
  about_me       = "Product manager with 10+ years of experience in tech"
  job_title      = "Senior Product Manager"
  preferred_name = "Jane"

  # Organizational information
  department    = "Product Management"
  company_name  = "Contoso Ltd."
  employee_id   = "E12345"
  employee_type = "Full-Time"

  # Contact information
  usage_location  = "UK"
  city            = "London"
  country         = "United Kingdom"
  office_location = "London HQ, Floor 3"
  mobile_phone    = "+44 7700 900123"
  business_phones = ["+44 20 7946 0958", "+44 20 7946 0959"]
  fax_number      = "+44 20 7946 0957"

  # Address information
  street_address = "123 Oxford Street"
  postal_code    = "W1D 1DF"
  state          = "England"

  # Language and location preferences
  preferred_language = "en-GB"

  # Additional email addresses
  other_mails     = ["jane.smith.personal@example.com"]
  proxy_addresses = ["SMTP:jane.smith@contoso.com", "smtp:jane.s@contoso.com"]

  # IM addresses
  im_addresses = ["jane.smith@contoso.com"]

  # Show in address list
  show_in_address_list = true

  # Identity configuration
  identities = [
    {
      sign_in_type       = "emailAddress"
      issuer             = "contoso.onmicrosoft.com"
      issuer_assigned_id = "jane.smith@contoso.com"
    },
    {
      sign_in_type       = "userPrincipalName"
      issuer             = "contoso.onmicrosoft.com"
      issuer_assigned_id = "jane.smith@contoso.com"
    }
  ]
} 