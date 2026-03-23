resource "random_string" "lifecycle_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "lifecycle" {
  account_enabled     = false
  display_name        = "License Assignment Lifecycle Test User"
  user_principal_name = "license.test.lifecycle.${random_string.lifecycle_suffix.result}@deploymenttheory.com"
  mail_nickname       = "license.test.lifecycle.${random_string.lifecycle_suffix.result}"
  usage_location      = "GB"
  password_profile = {
    password                           = "SecureP@ssw0rd123!!!!"
    force_change_password_next_sign_in = true
  }
}

// Step 2: remove disabled_plans entirely.
// The provider must send disabledPlans=[] to the API so all plans become enabled.
// If the bug were present, the API would silently retain the plan from Step 1 and
// the next terraform plan would show unexpected drift.
resource "microsoft365_graph_beta_users_user_license_assignment" "lifecycle" {
  user_id = microsoft365_graph_beta_users_user.lifecycle.id
  sku_id  = "a403ebcc-fae0-4ca2-8c8c-7a907fd6c235" # Microsoft Fabric (Free) / POWER_BI_STANDARD

  depends_on = [
    microsoft365_graph_beta_users_user.lifecycle
  ]
}
