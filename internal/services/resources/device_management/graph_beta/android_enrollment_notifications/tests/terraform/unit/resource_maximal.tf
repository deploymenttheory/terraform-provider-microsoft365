resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "complete" {
  display_name     = "Complete Test - Android Enrollment Notifications"
  description      = "Complete configuration for unit testing with all features"
  platform_type    = "androidForWork"
  default_locale   = "en-US"
  branding_options = "includeCompanyLogo"

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-US"
      subject          = "Device Enrollment"
      message_template = "Enroll your device now"
      is_default       = true
      template_type    = "push"
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