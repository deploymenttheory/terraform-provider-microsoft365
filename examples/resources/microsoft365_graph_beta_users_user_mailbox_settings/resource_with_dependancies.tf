# Example 3: Complete workflow with user creation and license assignment
# This example demonstrates creating a user, assigning an Exchange Online license,
# and then configuring mailbox settings

# Step 1: Create the user
resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "Example User"
  user_principal_name = "example.user@yourdomain.com"
  mail_nickname       = "example.user"
  account_enabled     = true
  usage_location      = "US" # Required for license assignment

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Step 2: Assign a license that includes Exchange Online
# Note: Replace the sku_id with an actual SKU ID from your tenant
resource "microsoft365_graph_beta_users_user_license_assignment" "example_user_license" {
  user_id = microsoft365_graph_beta_users_user.example_user.id

  # Common SKU IDs that include Exchange Online:
  # - Microsoft 365 E3: 6fd2c87f-b296-42f0-b197-1e91e994b900
  # - Microsoft 365 E5: c7df2760-2c81-4ef7-b578-5b5392b571df
  # - Microsoft 365 Business Premium: f245ecc8-75af-4f8e-b61f-27d8114de5f3
  # - Exchange Online Plan 1: 4b9405b0-7788-4568-add1-99614e613b69
  # - Exchange Online Plan 2: 19ec0d23-8335-4cbd-94ac-6050e30712fa
  sku_id = "6fd2c87f-b296-42f0-b197-1e91e994b900" # Microsoft 365 E3

  # Optional: Disable specific service plans
  disabled_plans = []
}

# Step 3: Wait for mailbox provisioning
# Exchange Online mailboxes can take 1-2 minutes to provision after license assignment
resource "time_sleep" "wait_for_mailbox" {
  depends_on = [microsoft365_graph_beta_users_user_license_assignment.example_user_license]

  create_duration = "2m"
}

# Step 4: Configure mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "example_user_settings" {
  depends_on = [time_sleep.wait_for_mailbox]

  user_id                                   = microsoft365_graph_beta_users_user.example_user.id
  time_zone                                 = "Eastern Standard Time"
  date_format                               = "MM/dd/yyyy"
  time_format                               = "hh:mm tt"
  delegate_meeting_message_delivery_options = "sendToDelegateAndInformationToPrincipal"

  automatic_replies_setting = {
    status            = "disabled"
    external_audience = "none"
  }

  language = {
    locale = "en-US"
  }

  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "08:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Eastern Standard Time"
    }
  }
}

