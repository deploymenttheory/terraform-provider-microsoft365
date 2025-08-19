
resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "push_maximal_android" {
  display_name           = "push maximal android"
  description            = "Complete configuration for unit testing with all features"
  platform_type          = "android"
  default_locale         = "en-US"
  branding_options       = ["none"]
  notification_templates = ["push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    }
  ]

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}