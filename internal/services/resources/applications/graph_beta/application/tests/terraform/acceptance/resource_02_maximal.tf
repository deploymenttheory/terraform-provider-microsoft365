# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "app_suffix" {
  length  = 8
  special = false
  upper   = false
}


# ==============================================================================
# Dependencies - Users for owners
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "dependency_owner_1" {
  display_name        = "acc-test-app-owner1-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-app-owner1-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-app-owner1-${random_string.app_suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_users_user" "dependency_owner_2" {
  display_name        = "acc-test-app-owner2-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-app-owner2-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-app-owner2-${random_string.app_suffix.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

# ==============================================================================
# Application
# ==============================================================================

# APP002: Maximal Application Configuration
# Tests application creation with all possible fields and nested configurations

resource "microsoft365_graph_beta_applications_application" "test_maximal" {
  display_name   = "acc-test-app-maximal-${random_string.app_suffix.result}"
  description    = "Maximal acceptance test application with all fields configured"
  sign_in_audience = "AzureADMyOrg"
  
  identifier_uris = [
    "https://deploymenttheory.com/acc-test-maximal-${random_string.app_suffix.result}"
  ]

  group_membership_claims = ["SecurityGroup"]
  notes                   = "This is a test application for acceptance testing"
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
    accept_mapped_claims          = true
    requested_access_token_version = 2
    oauth2_permission_scopes = []
    pre_authorized_applications = []
    known_client_applications = []
  }

  # App Roles
  app_roles = []

  # Informational URLs
  info = {
    marketing_url        = "https://contoso.com/marketing"
    privacy_statement_url = "https://contoso.com/privacy"
    support_url          = "https://contoso.com/support"
    terms_of_service_url = "https://contoso.com/terms"
  }

  # Key Credentials
  key_credentials = []

  # Password Credentials
  password_credentials = []

  # Optional Claims
  optional_claims = {
    access_token = []
    id_token     = []
    saml2_token  = []
  }

  # Parental Control Settings
  parental_control_settings = {
    countries_blocked_for_minors = ["US", "CA"]
    legal_age_group_rule         = "Allow"
  }

  # Public Client
  public_client = {
    redirect_uris = [
      "http://localhost"
    ]
  }

  # Required Resource Access
  required_resource_access = []

  # SPA Configuration
  spa = {
    redirect_uris = [
      "https://contoso.com/spa-callback"
    ]
  }

  # Web Configuration
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
    microsoft365_graph_beta_users_user.dependency_owner_1.id,
    microsoft365_graph_beta_users_user.dependency_owner_2.id
  ]

  prevent_duplicate_names = true
  hard_delete             = true
}


