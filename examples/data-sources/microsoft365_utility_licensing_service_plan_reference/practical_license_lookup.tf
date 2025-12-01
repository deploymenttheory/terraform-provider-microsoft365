# Practical Example: Dynamic license assignment using human-readable names
# This demonstrates how to avoid hardcoding GUIDs in your configurations

# Look up Microsoft 365 E3 (no Teams) by product name
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3_no_teams" {
  product_name = "Microsoft 365 E3 (no Teams)"
}

# Create a user
resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "Example User"
  user_principal_name = "example.user@yourdomain.com"
  mail_nickname       = "example.user"
  account_enabled     = true
  usage_location      = "US"

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Assign the license using the dynamically looked-up GUID
resource "microsoft365_graph_beta_users_user_license_assignment" "example_license" {
  user_id = microsoft365_graph_beta_users_user.example_user.id

  # No hardcoded GUID - always up-to-date from Microsoft's reference data
  sku_id = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].guid
}

# Output the license details for verification
output "assigned_license" {
  value = {
    product_name        = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].product_name
    sku_id              = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].guid
    service_plans_count = length(data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].service_plans_included)
  }
}

