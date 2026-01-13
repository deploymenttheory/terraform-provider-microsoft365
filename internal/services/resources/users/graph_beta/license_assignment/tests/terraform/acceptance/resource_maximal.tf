resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "maximal" {
  account_enabled     = false
  display_name        = "License Assignment Test Maximal User"
  user_principal_name = "license.test.maximal.${random_string.maximal_suffix.result}@deploymenttheory.com"
  mail_nickname       = "license.test.maximal.${random_string.maximal_suffix.result}"
  usage_location      = "GB"
  password_profile = {
    password                           = "SecureP@ssw0rd123!!!!"
    force_change_password_next_sign_in = true
  }
}

resource "microsoft365_graph_beta_users_user_license_assignment" "maximal" {
  user_id = microsoft365_graph_beta_users_user.maximal.id
  sku_id  = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  disabled_plans = [
    "c948ea65-2053-4a5a-8a62-9eaaaf11b522" # PURVIEW_DISCOVERY
  ]

  depends_on = [
    microsoft365_graph_beta_users_user.maximal
  ]
}
