resource "microsoft365_graph_beta_applications_application" "test_maximal" {
  display_name     = "acc-test-app-maximal-${random_string.app_suffix.result}"
  description      = "Maximal test application"
  sign_in_audience = "AzureADMyOrg"

  # Application Identifier URIs are managed by the separate
  # microsoft365_graph_beta_applications_application_identifier_uri resource

  # Key credentials are managed by the separate
  # microsoft365_graph_beta_applications_application_certificate_credential resource

  # Password Credentials are managed by the separate
  # microsoft365_graph_beta_applications_application_password_credential resource

  group_membership_claims       = ["SecurityGroup"]
  notes                         = "This is a test application for acceptance testing"
  is_device_only_auth_supported = false
  is_fallback_public_client     = false
  service_management_reference  = "https://contoso.com/app-management"

  tags = [
    "terraform",
    "acceptance-test",
    "maximal"
  ]

  # API Configuration
  api = {
    accept_mapped_claims           = true
    requested_access_token_version = 2
  }

  # App Roles - Demonstrating all allowed_member_types combinations
  app_roles = [
    {
      id                   = random_uuid.app_role_id_1.result
      allowed_member_types = ["User"]
      description          = "App role assignable to users and groups"
      display_name         = "User Role"
      is_enabled           = true
      value                = "User.Role"
    },
    {
      id                   = random_uuid.app_role_id_2.result
      allowed_member_types = ["Application"]
      description          = "App role assignable to other applications (application permission)"
      display_name         = "Application Role"
      is_enabled           = false
      value                = "Application.Role"
    },
    {
      id                   = random_uuid.app_role_id_3.result
      allowed_member_types = ["User", "Application"]
      description          = "App role assignable to both users and applications"
      display_name         = "Combined Role"
      is_enabled           = false
      value                = "Combined.Role"
    }
  ]

  info = {
    marketing_url         = "https://contoso.com/marketing"
    privacy_statement_url = "https://contoso.com/privacy"
    support_url           = "https://contoso.com/support"
    terms_of_service_url  = "https://contoso.com/terms"
  }

  parental_control_settings = {
    countries_blocked_for_minors = ["US", "CA"]
    legal_age_group_rule         = "Allow"
  }

  public_client = {
    redirect_uris = [
      "http://localhost"
    ]
  }

  spa = {
    redirect_uris = [
      "https://contoso.com/spa-callback"
    ]
  }

  web = {
    home_page_url = "https://contoso.com"
    logout_url    = "https://contoso.com/logout"
    redirect_uris = [
      "https://contoso.com/callback"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
  }

  owner_user_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  prevent_duplicate_names = true
  hard_delete             = true
}


