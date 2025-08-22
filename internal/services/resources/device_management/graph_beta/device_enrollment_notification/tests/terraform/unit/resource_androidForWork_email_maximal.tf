
resource "microsoft365_graph_beta_device_management_device_enrollment_notification" "email_maximal_androidforwork" {
  display_name   = "email maximal androidForWork"
  description    = "Complete configuration for unit testing with all features"
  platform_type  = "androidForWork"
  default_locale = "en-US"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

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