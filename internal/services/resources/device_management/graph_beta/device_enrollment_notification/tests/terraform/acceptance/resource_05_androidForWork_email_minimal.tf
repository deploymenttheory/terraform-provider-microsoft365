# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-enrollment-notification-group-1-${random_string.suffix.result}"
  mail_nickname    = "acc-test-enrollment-notification-group-1-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group for device enrollment notification assignments"
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1
  ]
  create_duration = "15s"
}

# ==============================================================================
# Device Enrollment Notification
# ==============================================================================

resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "email_minimal_androidforwork" {
  display_name     = "email minimal androidForWork"
  description      = "Complete configuration for unit testing with all features"
  platform_type    = "androidForWork"
  default_locale   = "en-US"
  branding_options = ["none"]

  notification_templates = ["email"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    }
  ]

  role_scope_tag_ids = ["0"]

  depends_on = [time_sleep.wait_for_groups]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
