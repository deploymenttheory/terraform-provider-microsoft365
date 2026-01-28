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

resource "microsoft365_graph_beta_users_user" "dependency_owner" {
  display_name        = "acc-test-pc-owner-${random_string.app_suffix.result}"
  user_principal_name = "acc-test-pc-owner-${random_string.app_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-pc-owner-${random_string.app_suffix.result}"
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

# APP005: Public Client (Native/Mobile) Application Configuration
# Tests mobile/desktop application configuration with public client settings

resource "microsoft365_graph_beta_applications_application" "test_public_client" {
  display_name              = "acc-test-public-client-${random_string.app_suffix.result}"
  description               = "Public client acceptance test application"
  sign_in_audience          = "AzureADMyOrg"
  is_fallback_public_client = true

  public_client = {
    redirect_uris = [
      "http://localhost",
      "ms-appx-web://microsoft.aad.brokerplugin/acc-test-public-${random_string.app_suffix.result}"
    ]
  }

  required_resource_access = []

  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_owner.id
  ]

  prevent_duplicate_names = false
  hard_delete             = true
}


