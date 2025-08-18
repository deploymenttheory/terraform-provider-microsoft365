resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "example" {
  display_name   = "Android Enrollment Notifications"
  description    = "Android Enrollment Notification example"
  platform_type  = "androidForWork" // Options are: "androidForWork" , "android"
  default_locale = "en-US"          // This is the default locale for the notification messages.
  branding_options = [              // These branding optional apply to email notifications only.
    "includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    // Optional localized notification messages.
    // These need to match the notification templates defined.
    {
      locale           = "en-us" // should match the default_locale but be in lowercase
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-us" // should match the default_locale but be in lowercase
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    # {
    #   type     = "allLicensedUsersAssignmentTarget"
    # }
  ]

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}