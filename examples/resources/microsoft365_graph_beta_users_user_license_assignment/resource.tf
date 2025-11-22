# Example 1: Assign a single Office 365 E3 license to a user with disabled service plans
resource "microsoft365_graph_beta_users_user_license_assignment" "user_e3_license" {
  user_id = "john.doe@example.com" # Can be user ID (UUID) or UPN

  sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  # Optional: Disable specific service plans within this license
  disabled_plans = [
    "efb87545-963c-4e0d-99df-69c6916d9eb0" # Example: Microsoft Stream
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Example 2: Assign a license without disabling any service plans
resource "microsoft365_graph_beta_users_user_license_assignment" "user_powerbi_license" {
  user_id = "john.doe@example.com"
  sku_id  = "f30db892-07e9-47e9-837c-80727f46fd3d" # Power BI (free)
}

# Example 3: Assign multiple licenses to a single user
# Note: Create multiple resource instances, one per license
resource "microsoft365_graph_beta_users_user_license_assignment" "jane_e3" {
  user_id = "jane.smith@example.com"
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

resource "microsoft365_graph_beta_users_user_license_assignment" "jane_ems_e5" {
  user_id = "jane.smith@example.com"
  sku_id  = "b05e124f-c7cc-45a0-a6aa-8cf78c946968" # Enterprise Mobility + Security E5

  disabled_plans = [
    "8a256a2b-b617-496d-b51b-e76466e88db0" # Microsoft Defender for Cloud Apps
  ]
}

# Example 4: Assign Office 365 E5 license with multiple disabled plans
resource "microsoft365_graph_beta_users_user_license_assignment" "alice_e5" {
  user_id = "alice.wilson@example.com"
  sku_id  = "c7df2760-2c81-4ef7-b578-5b5392b571df" # Office 365 E5

  disabled_plans = [
    "57ff2da0-773e-42df-b2af-ffb7a2317929", # Teams
    "0feaeb32-d00e-4d66-bd5a-43b5b83db82c"  # Mya
  ]
}

# Example 5: Using a data source to get user ID dynamically
data "microsoft365_graph_beta_users_user" "target_user" {
  user_principal_name = "dynamic.user@example.com"
}

resource "microsoft365_graph_beta_users_user_license_assignment" "dynamic_user_license" {
  user_id = data.microsoft365_graph_beta_users_user.target_user.id
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 6: Using for_each to assign the same license to multiple users
variable "licensed_users" {
  type = set(string)
  default = [
    "user1@example.com",
    "user2@example.com",
    "user3@example.com"
  ]
}

resource "microsoft365_graph_beta_users_user_license_assignment" "bulk_e3_assignment" {
  for_each = var.licensed_users

  user_id = each.value
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3
}

# Example 7: Create user and assign license in one configuration
resource "microsoft365_graph_beta_users_user" "new_user" {
  user_principal_name = "new.employee@example.com"
  display_name        = "New Employee"
  mail_nickname       = "new.employee"
  account_enabled     = true
  usage_location      = "US"

  password_profile = {
    password                           = "TemporaryP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

resource "microsoft365_graph_beta_users_user_license_assignment" "new_user_license" {
  user_id = microsoft365_graph_beta_users_user.new_user.id
  sku_id  = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Office 365 E3

  depends_on = [microsoft365_graph_beta_users_user.new_user]
}

# Common Microsoft 365 SKU IDs for reference:
# - Office 365 E3: 6fd2c87f-b296-42f0-b197-1e91e994b900
# - Office 365 E5: c7df2760-2c81-4ef7-b578-5b5392b571df
# - Enterprise Mobility + Security E5: b05e124f-c7cc-45a0-a6aa-8cf78c946968
# - Microsoft 365 Business Premium: cbdc14ab-d96c-4c30-b9f4-6ada7cdc1d46
# - Power BI (free): f30db892-07e9-47e9-837c-80727f46fd3d

# Note: To remove a license, simply destroy the resource:
# terraform destroy -target=microsoft365_graph_beta_users_user_license_assignment.user_e3_license
