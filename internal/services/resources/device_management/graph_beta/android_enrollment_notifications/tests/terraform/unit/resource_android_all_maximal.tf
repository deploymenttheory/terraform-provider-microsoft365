


resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "all_android" {
  display_name   = "Complete Test - all android"
  description    = "Complete configuration for unit testing with all features"
  platform_type  = "android"
  default_locale = "en-US"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
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
      group_id = "00000000-0000-0000-0000-000000000001"
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