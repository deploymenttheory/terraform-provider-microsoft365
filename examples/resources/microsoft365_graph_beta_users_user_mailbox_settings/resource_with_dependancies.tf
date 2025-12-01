# Example 3: Complete workflow with user creation and license assignment
# This example demonstrates creating a user, assigning an Exchange Online license,
# and then configuring mailbox settings

# Step 1: Look up the license SKU using the licensing service plan reference datasource
# This ensures you always have the correct GUID without hardcoding it
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3" {
  string_id = "ENTERPRISEPACK" # Microsoft 365 E3

  # Alternative search options:
  # product_name = "Microsoft 365 E3"
  # guid = "6fd2c87f-b296-42f0-b197-1e91e994b900"
}

# Step 2: Create the user
resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "Example User"
  user_principal_name = "example.user@yourdomain.com"
  mail_nickname       = "example.user"
  account_enabled     = true
  usage_location      = "US" # Field is required for license assignment

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Step 3: Assign a license that includes Exchange Online
resource "microsoft365_graph_beta_users_user_license_assignment" "example_user_license" {
  user_id = microsoft365_graph_beta_users_user.example_user.id

  # Use the dynamically looked-up SKU ID from the datasource
  # This is more maintainable than hardcoding GUIDs and ensures accuracy
  sku_id = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].guid

  # Optional: Disable specific service plans
  disabled_plans = []
}

# Step 4: Wait for mailbox provisioning
# Exchange Online mailboxes can take 1-2 minutes to provision after license assignment
resource "time_sleep" "wait_for_mailbox_provisioning" {
  depends_on = [microsoft365_graph_beta_users_user_license_assignment.example_user_license]

  create_duration = "2m"
}

# Step 5: Configure mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "example_user_settings" {
  depends_on = [time_sleep.wait_for_mailbox_provisioning]

  user_id                                   = microsoft365_graph_beta_users_user.maximal_dependency_user.id
  time_zone                                 = "Greenwich Standard Time"
  delegate_meeting_message_delivery_options = "sendToDelegateOnly"

  automatic_replies_setting = {
    status            = "scheduled"
    external_audience = "all"

    scheduled_start_date_time = {
      date_time = "2030-03-14T07:00:00"
      time_zone = "UTC"
    }

    scheduled_end_date_time = {
      date_time = "2030-03-28T07:00:00"
      time_zone = "UTC"
    }

    internal_reply_message = "<html>\n<body>\n<p>I'm at our company's worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
    external_reply_message = "<html>\n<body>\n<p>I'm at the Deployment Theory worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
  }

  language = {
    locale = "en-US"
  }

  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "08:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Greenwich Standard Time"
    }
  }
}

