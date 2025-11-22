resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user_license_assignment" "dependancy" {
  user_id = microsoft365_graph_beta_users_user.dependancy.id
  sku_id  = "f30db892-07e9-47e9-837c-80727f46fd3d" # FLOW_FREE
}

resource "microsoft365_graph_beta_users_user" "dependancy" {
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
